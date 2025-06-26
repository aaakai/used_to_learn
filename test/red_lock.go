package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func RenneLock(client *redis.Client, key string, duration time.Duration, ctx context.Context, callback func()) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 续期锁的过期时间
			err := client.Expire(ctx, key, duration).Err()
			callback()
			if err != nil {
				fmt.Println("Error renewing lock:", err)
				return
			}
			fmt.Println("Lock renewed.")

		case <-ctx.Done():
			fmt.Println("Context cancelled, stopping lock renewal.")
			return
		}
	}
}

func test() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	key := "my_lock"
	duration := 10 * time.Second
	cancelCtx, cancel := context.WithCancel(ctx)
	RenneLock(client, key, duration, cancelCtx, func() {
		fmt.Println("Doing some work...", key)
	})
	defer cancel()
}
