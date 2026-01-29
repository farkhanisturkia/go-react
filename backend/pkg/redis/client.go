package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func Init() {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASS")
	db := 0

	Client = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis connection failed: %v", err))
	}

	fmt.Println("Redis connected")
}
