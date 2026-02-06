package project

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"icode.baidu.com/baidu/bjh-server-go/aigc-ui/library/gcp"
	"icode.baidu.com/baidu/bjh-server-go/aigc-ui/library/utils"
	"icode.baidu.com/baidu/gdp/extension/gtask"
	"icode.baidu.com/baidu/gdp/net/gaddr"

	"icode.baidu.com/baidu/gdp/logit"

	"icode.baidu.com/baidu/bjh-server-go/aigc-ui/errcode"
	"icode.baidu.com/baidu/bjh-server-go/aigc-ui/library/resource"
	"icode.baidu.com/baidu/bjh-server-go/aigc-ui/model/service/data/uriLimitConfig"
)

const (
	UriLimitRedisKeyPreFix = "aigc_uri_limit"
	RedisKeySecExpire      = 10
	RedisKeyMinExpire      = 70
	RedisKeyDayExpire      = 3600*24 + 10
)

type UriLimitInfo struct {
	Limit       *uriLimitConfig.UriLimitConfig
	UriKey      string
	BlackList   map[string]int
	IPBlackList map[string]int
}

type LimiterConfig struct {
	Limit       int    `json:"limit"`       // 窗口请求上限
	Window      int    `json:"window"`      // 窗口时间大小
	SmallWindow int    `json:"smallWindow"` // 小窗口时间大小
	LimitTime   int64  `json:"limit_time"`  // 封禁时间
	LimiterKey  string `json:"limiter_key"` // 限流封禁key
}

type UserLimiterInfo struct {
	Switch             int                      `json:"switch"`
	UserURLLimitConfig map[string]LimiterConfig `json:"user_url_limit_config"`
}

// 处理用户请求频控信息
func HandleUriLimit(ctx context.Context, uri string, uid int64, uip string) error {
	detailLimitInfo := getLimitInfo(ctx, uri)
	if detailLimitInfo == nil || detailLimitInfo.UriKey == "" {
		return nil
	}

	userLimiterInfo := getUserLimiterInfo(ctx, uri)

	uriKey := utils.GetMd5(detailLimitInfo.UriKey)

	// 处理下游服务熔断
	if detailLimitInfo.Limit.InterceptSwitch > 0 {
		return errcode.ErrServiceIntercept
	}
	logit.AddNotice(ctx, logit.String("UriLimitKey", uriKey))

	var (
		nowTime = time.Now()
		g       = &gtask.Group{
			Concurrent:    3,    // 控制并发上限
			AllowSomeFail: true, // 允许部分任务执行失败
		}
	)

	g.Go(func() error {
		// 用户纬度频控
		if err := checkUIDLimit(ctx, uriKey, uid, nowTime, detailLimitInfo); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		// ip 维度频控
		if err := checkIPLimit(ctx, uriKey, uip, nowTime, detailLimitInfo); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		// user 维度频控
		if userLimiterInfo == nil || userLimiterInfo.LimiterKey == "" {
			return nil
		}
		if err := checkUserLimiter(ctx, userLimiterInfo, uid); err != nil {
			return err
		}
		return nil
	})

	_, err := g.Wait()
	if err != nil {
		return err
	}

	// 接口整体频控判断(放在前置判断后)
	redisKey := fmt.Sprintf("%s%s_total_s_%d", UriLimitRedisKeyPreFix, uriKey, nowTime.Unix())
	if err := uriLimitRedisProcess(ctx, redisKey, "limit_uid_total_s", detailLimitInfo.Limit.LimitTotalS, time.Second*RedisKeySecExpire); err != nil {
		return errcode.ErrFreqControlLimit
	}
	return nil
}

// uidLimitHandle uid频控判断
func checkUIDLimit(ctx context.Context, uriKey string, uid int64, nowTime time.Time, detailLimitInfo *UriLimitInfo) error {
	if uid <= 0 {
		logit.AddNotice(ctx, logit.String("jump_id_limit", "1"))
		return nil
	}

	// 判断是否在黑名单（可以移动反作弊里面）
	if value, ok := detailLimitInfo.BlackList[fmt.Sprintf("%d", uid)]; ok && value == 1 {
		logit.AddWarning(ctx, logit.String("freq_reason", "black_uid"))
		return errcode.ErrUserBlack
	}

	var redisKey string
	// 秒级用户纬度频控判断
	redisKey = fmt.Sprintf("%s%s_s_%d_%d", UriLimitRedisKeyPreFix, uriKey, nowTime.Unix(), uid)
	if err := uriLimitRedisProcess(ctx, redisKey, "limit_uid_s", detailLimitInfo.Limit.LimitS, time.Second*RedisKeySecExpire); err != nil {
		return errcode.ErrFreqControlUIDLimit
	}

	// 分钟级用户纬度频控判断
	redisKey = fmt.Sprintf("%s%s_m_%d_%d", UriLimitRedisKeyPreFix, uriKey, nowTime.Unix()-int64(nowTime.Second()), uid)
	if err := uriLimitRedisProcess(ctx, redisKey, "limit_uid_m", detailLimitInfo.Limit.LimitM, time.Second*RedisKeyMinExpire); err != nil {
		return errcode.ErrFreqControlUIDLimit
	}

	// 天级用户纬度频控判断
	redisKey = fmt.Sprintf("%s%s_d_%s_%d", UriLimitRedisKeyPreFix, uriKey, nowTime.Format("20060102"), uid)
	if err := uriLimitRedisProcess(ctx, redisKey, "limit_uid_d", detailLimitInfo.Limit.LimitD, time.Second*RedisKeyDayExpire); err != nil {
		return errcode.ErrFreqControlUIDLimit
	}

	return nil
}

func checkUserLimiter(ctx context.Context, limitConfig *LimiterConfig, uid int64) error {
	if uid <= 0 {
		logit.AddNotice(ctx, logit.String("jump_id_limit", "1"))
		return nil
	}

	// 用户处于封禁使用状态
	status, err := resource.RedisClientAigc.Get(ctx, fmt.Sprintf("uid_user_limit_block_%d", uid)).Result()
	if err == nil && status == "1" {
		logit.AddNotice(ctx, logit.String("uid_user_limit_block", "1"))
		return errcode.ErrUserBlack
	}

	// 滑动窗口用户纬度频控判断
	client := resource.RedisClientAigc
	l, err := NewSlidingWindowLimiter(&client, limitConfig.Limit, limitConfig.Window)
	if err != nil {
		logit.AddError(ctx, logit.String("NewSlidingWindowLimiter err", err.Error()))
		return nil
	}
	if err := l.TryAcquire(ctx, fmt.Sprintf("%d", uid)); err != nil {
		// 存拦截记录
		value := fmt.Sprintf("user_limiter_%d_%d", uid, time.Now().Unix())
		resource.RedisClientAigc.SAdd(ctx, limitConfig.LimiterKey, value)

		// 封禁用户使用
		redisKey := fmt.Sprintf("uid_user_limit_block_%d", uid)
		limitTime := time.Duration(limitConfig.LimitTime) * time.Second
		resource.RedisClientAigc.Set(ctx, redisKey, "1", limitTime)

		return errcode.ErrFreqControlUIDLimit
	}
	return nil
}

// ipLimitHandle ip频控判断
func checkIPLimit(ctx context.Context, uriKey string, uip string, nowTime time.Time, detailLimitInfo *UriLimitInfo) error {
	if uip == "" || gaddr.IsInternalIP(net.ParseIP(uip)) {
		logit.AddNotice(ctx, logit.String("jump_ip_limit", uip))
		// 内网不拦截，防止内网网关拦截
		return nil
	}

	// 判断是否在黑名单（可以移动反作弊里面）
	if value, ok := detailLimitInfo.IPBlackList[uip]; ok && value == 1 {
		logit.AddWarning(ctx, logit.String("freq_reason", "black_uip"))
		return errcode.ErrIpBlack
	}
	if value, ok := detailLimitInfo.Limit.IPBlackList[uip]; ok && value == 1 {
		logit.AddWarning(ctx, logit.String("freq_reason", "limit_black_uip"))
		return errcode.ErrIpBlack
	}

	var redisKey string
	// 秒级IP纬度频控判断
	redisKey = fmt.Sprintf("%s%s_ip_s_%d_%s", UriLimitRedisKeyPreFix, uriKey, nowTime.Unix(), uip)
	if err := uriLimitRedisProcess(ctx, redisKey, "limit_ip_s", detailLimitInfo.Limit.IPLimitS, time.Second*RedisKeySecExpire); err != nil {
		return errcode.ErrFreqControlIPLimit
	}

	// 分钟级IP纬度频控判断
	redisKey = fmt.Sprintf("%s%s_ip_m_%d_%s", UriLimitRedisKeyPreFix, uriKey, nowTime.Unix()-int64(nowTime.Second()), uip)
	if err := uriLimitRedisProcess(ctx, redisKey, "limit_ip_m", detailLimitInfo.Limit.IPLimitM, time.Second*RedisKeyMinExpire); err != nil {
		return errcode.ErrFreqControlIPLimit
	}

	// 天级IP纬度频控判断
	redisKey = fmt.Sprintf("%s%s_ip_d_%s_%s", UriLimitRedisKeyPreFix, uriKey, nowTime.Format("20060102"), uip)
	if err := uriLimitRedisProcess(ctx, redisKey, "limit_ip_d", detailLimitInfo.Limit.IPLimitD, time.Second*RedisKeyDayExpire); err != nil {
		return errcode.ErrFreqControlUIDLimit
	}

	return nil
}

// uriLimitRedisProcess redis处理通用方法
func uriLimitRedisProcess(ctx context.Context, redisKey, reason string, limitNum int64, expire time.Duration) error {
	if limitNum <= 0 {
		return nil
	}

	incrRet, err := resource.RedisClientDujiaService.Incr(ctx, redisKey).Result()
	go func() {
		contextBg := logit.CopyAllFields(context.Background(), ctx)
		resource.RedisClientDujiaService.Expire(contextBg, redisKey, expire)
	}()
	if err == nil && incrRet > limitNum {
		logit.AddWarning(ctx, logit.String("freq_reason", reason))
		return errcode.ErrFreqControl
	}

	return nil
}

// getLimitInfo 获取限制详情
func getLimitInfo(ctx context.Context, uri string) *UriLimitInfo {
	// 获取gcp配置
	uriLimitInfo, err := uriLimitConfig.GetUriLimitConfigInfo(ctx)
	if err != nil {
		// 出错不影响正常流程
		logit.AddWarning(ctx, logit.AutoField("HandleUriLimit", err.Error()))
		return nil
	}

	// 获取uri频控信息(取消产品线、接口类型、版本后的判断，如果需要详细拆解配置可以使用正则配置)
	uriParts := strings.Split(uri, "/")
	if len(uriParts) >= 5 {
		uriParts = append(uriParts[:2], uriParts[5:]...)
	}
	checkUri := strings.Join(uriParts, "/")
	if item, ok := uriLimitInfo.UriLimitConfig[checkUri]; ok {
		return &UriLimitInfo{Limit: &item, BlackList: uriLimitInfo.BlackList, IPBlackList: uriLimitInfo.IPBlackList, UriKey: checkUri}
	}

	// 获取uri频控信息（正则配置）
	for uriKey, item := range uriLimitInfo.UriLimitConfig {
		if !item.NeedCheckReg {
			// 不需要正则判断
			continue
		}

		reg, err := regexp.Compile(uriKey)
		if reg == nil || err != nil {
			// 解释器初始化失败
			continue
		}
		result := reg.FindStringSubmatch(uri)
		if result != nil {
			// 匹配
			return &UriLimitInfo{Limit: &item, BlackList: uriLimitInfo.BlackList, IPBlackList: uriLimitInfo.IPBlackList, UriKey: checkUri}
		}
	}

	return nil
}

func getSlidingWindowLimiterConfig(ctx context.Context) (*UserLimiterInfo, error) {
	editConf, err := gcp.GetGcpData(ctx, "saas_user_limit_config", &UserLimiterInfo{}, map[string]any{})
	if err != nil {
		return nil, err
	}
	config, ok := editConf.(*UserLimiterInfo)
	if !ok {
		logit.AddWarning(ctx, logit.AutoField("get_saas_user_limit_config_failed", config))
		return nil, errors.New("get_saas_user_limit_config_failed")
	}
	return config, nil
}

func getUserLimiterInfo(ctx context.Context, url string) *LimiterConfig {
	// 获取gcp配置
	userLimitInfo, err := getSlidingWindowLimiterConfig(ctx)
	if err != nil {
		// 出错不影响正常流程
		logit.AddWarning(ctx, logit.AutoField("getUserLimiterInfo", err.Error()))
		return nil
	}
	if userLimitInfo != nil && userLimitInfo.Switch == 0 {
		return nil
	}
	logit.AddNotice(ctx, logit.AutoField("getUserLimiterInfo", userLimitInfo))
	// 获取url频控信息
	uriParts := strings.Split(url, "/")
	if len(uriParts) >= 5 {
		uriParts = append(uriParts[:2], uriParts[5:]...)
	}
	checkUri := strings.Join(uriParts, "/")
	if item, ok := userLimitInfo.UserURLLimitConfig[checkUri]; ok {
		return &item
	}
	return nil
}
