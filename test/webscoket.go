package test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"log"
	"net/http"
	"sync"
	"time"
)

// 客户端连接信息
type Client struct {
	ID            string          // 连接ID
	AccountId     string          // 账号id, 一个账号可能有多个连接
	Socket        *websocket.Conn // 连接
	HeartbeatTime int64           // 前一次心跳时间
}

// 消息类型
const (
	MessageTypeHeartbeat = "heartbeat" // 心跳
	MessageTypeRegister  = "register"  // 注册

	HeartbeatCheckTime = 9  // 心跳检测几秒检测一次
	HeartbeatTime      = 20 // 心跳距离上一次的最大时间

	ChanBufferRegister   = 100 // 注册chan缓冲
	ChanBufferUnregister = 100 // 注销chan大小
)

// 客户端管理
type ClientManager struct {
	Clients  map[string]*Client  // 保存连接
	Accounts map[string][]string // 账号和连接关系,map的key是账号id即：AccountId，这里主要考虑到一个账号多个连接
	mu       *sync.Mutex
}

// 定义一个管理Manager
var Manager = ClientManager{
	Clients:  make(map[string]*Client),  // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Accounts: make(map[string][]string), // 账号和连接关系
	mu:       new(sync.Mutex),
}

var (
	RegisterChan   = make(chan *Client, ChanBufferRegister)   // 注册
	unregisterChan = make(chan *Client, ChanBufferUnregister) // 注销
)

// 封装回复消息
type ServiceMessage struct {
	Type    string                `json:"type"` // 类型
	Content ServiceMessageContent `json:"content"`
}
type ServiceMessageContent struct {
	Body     string `json:"body"`      // 主要数据
	MetaData string `json:"meta_data"` // 扩展数据
}

func CreateReplyMsg(t string, content ServiceMessageContent) []byte {
	replyMsg := ServiceMessage{
		Type:    t,
		Content: content,
	}
	msg, _ := json.Marshal(replyMsg)
	return msg
}

// 注册注销
func register() {
	for {
		select {
		case conn := <-RegisterChan: // 新注册，新连接
			// 加入连接,进行管理
			accountBind(conn)

			// 回复消息
			content := CreateReplyMsg(MessageTypeRegister, ServiceMessageContent{})
			_ = conn.Socket.WriteMessage(websocket.TextMessage, content)

		case conn := <-unregisterChan: // 注销，或者没有心跳
			// 关闭连接
			_ = conn.Socket.Close()

			// 删除Client
			unAccountBind(conn)
		}
	}
}

// 绑定账号
func accountBind(c *Client) {
	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	// 加入到连接
	Manager.Clients[c.ID] = c

	// 加入到绑定
	if _, ok := Manager.Accounts[c.AccountId]; ok { // 该账号已经有绑定，就追加一个绑定
		Manager.Accounts[c.AccountId] = append(Manager.Accounts[c.AccountId], c.ID)
	} else { // 没有就新增一个账号的绑定切片
		Manager.Accounts[c.AccountId] = []string{c.ID}
	}
}

// 解绑账号
func unAccountBind(c *Client) {
	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	// 取消连接
	delete(Manager.Clients, c.ID)

	// 取消绑定
	if len(Manager.Accounts[c.AccountId]) > 0 {
		for k, clientId := range Manager.Accounts[c.AccountId] {
			if clientId == c.ID { // 找到绑定客户端Id
				Manager.Accounts[c.AccountId] = append(Manager.Accounts[c.AccountId][:k], Manager.Accounts[c.AccountId][k+1:]...)
			}
		}
	}
}

// 维持心跳
func heartbeat() {
	for {
		// 获取所有的Clients
		Manager.mu.Lock()
		clients := make([]*Client, len(Manager.Clients))
		for _, c := range Manager.Clients {
			clients = append(clients, c)
		}
		Manager.mu.Unlock()

		for _, c := range clients {
			if time.Now().Unix()-c.HeartbeatTime > HeartbeatTime {
				unAccountBind(c)
			}
		}

		time.Sleep(time.Second * HeartbeatCheckTime)
	}
}

// 管理连接
func Start() {
	// 检查心跳
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		heartbeat()
	}()

	// 注册注销
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		register()
	}()
}

// 根据账号获取连接
func GetClient(accountId string) []*Client {
	clients := make([]*Client, 0)

	Manager.mu.Lock()
	defer Manager.mu.Unlock()

	if len(Manager.Accounts[accountId]) > 0 {
		for _, clientId := range Manager.Accounts[accountId] {
			if c, ok := Manager.Clients[clientId]; ok {
				clients = append(clients, c)
			}
		}
	}

	return clients
}

// 读取信息，即收到消息
func (c *Client) Read() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		// 读取消息
		_, body, err := c.Socket.ReadMessage()
		if err != nil {
			break
		}

		var msg struct {
			Type string `json:"type"`
		}
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Println(err)
			continue
		}

		if msg.Type == MessageTypeHeartbeat { // 维持心跳消息
			// 刷新连接时间
			c.HeartbeatTime = time.Now().Unix()

			// 回复心跳
			replyMsg := CreateReplyMsg(MessageTypeHeartbeat, ServiceMessageContent{})
			err = c.Socket.WriteMessage(websocket.TextMessage, replyMsg)
			if err != nil {
				log.Println(err)
			}
			continue
		}
	}
}

// 发送消息
func Send(accounts []string, message ServiceMessage) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for _, accountId := range accounts {
		// 获取连接id
		clients := GetClient(accountId)

		// 发送消息
		for _, c := range clients {
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}

	return nil
}

type MessageNotifyRequest struct {
	UserId string `form:"user_id"`
}

func MessageNotify(ctx *gin.Context) {
	// 获取参数
	var params MessageNotifyRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		log.Println(err)
		return
	}
	// TODO: 鉴权

	// 将http升级为websocket
	conn, err := (&websocket.Upgrader{
		// 1. 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil) // 升级
	if err != nil {
		log.Println(err)
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	// 创建一个实例连接
	ConnId := xid.New().String()
	client := &Client{
		ID:            ConnId, // 连接id
		AccountId:     fmt.Sprintf("%s", params.UserId),
		HeartbeatTime: time.Now().Unix(),
		Socket:        conn,
	}

	// 用户注册到用户连接管理
	RegisterChan <- client

	// 读取信息
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("MessageNotify read panic: %+v\n", r)
			}
		}()

		client.Read()
	}()
}
