package helper

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var once sync.Once

func GetRedisClient() *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Protocol: 2,
		})

		ctx := context.Background()
		_, err := redisClient.Ping(ctx).Result()

		if err != nil {
			log.Fatalf("Could not connect to Redis: %v", err)
		}

		fmt.Println("Connected to Redis!")

	})
	fmt.Println(redisClient)
	return redisClient
}

func ClearRedis() error {
	client := GetRedisClient()
	ctx := context.Background()

	err := client.FlushAll(ctx).Err()

	if err != nil {
		return fmt.Errorf("failed to clear Redis: %v", err)
	}
	fmt.Println("Redis cleared successfully!")
	return nil

}
