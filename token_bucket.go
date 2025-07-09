package main

import (
	"sync"
	"time"
)

// 令牌桶
type TokenBucket struct {
	capacity       int
	tokens         int
	rate           float64
	lastUpdateTime time.Time
	mutex          sync.Mutex
}

func NewTokenBucket(capacity int, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:       capacity,
		tokens:         0,
		rate:           rate,
		lastUpdateTime: time.Now(),
	}
}

func (tb *TokenBucket) GetToken() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastUpdateTime).Seconds()
	tokensToAdd := int(elapsed * tb.rate)

	// 确保不会超过容量上限
	if tokensToAdd > tb.capacity-tb.tokens {
		tokensToAdd = tb.capacity - tb.tokens
	}
	tb.tokens += tokensToAdd
	// 如果添加了令牌，更新上次更新时间
	if tokensToAdd > 0 {
		tb.lastUpdateTime = now
	}
	// 检查是否有足够的令牌可用，如果有，减少令牌数并返回 true，否则返回 false
	if tb.tokens > 0 {
		tb.tokens-- // 使用一个令牌
		return true // 成功获取令牌
	}
	return false // 桶中没有足够的令牌，获取失败

}
