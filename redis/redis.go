package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lambda-lama/user-api/config"
	"github.com/redis/go-redis/v9"
)

func getConnection() *redis.Client {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.RedisPassword,
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println("ERR: ", err)
	}

	return client
}

func SetString(key string, value string, expiration time.Duration) error {
	client := getConnection()
	defer client.Close()

	_, err := client.Set(context.Background(), key, value, expiration).Result()
	if err != nil {
		fmt.Printf("Cache operation failed: %v\n", err)
		return err
	}

	return nil
}

// func Set(key string, value []byte, expiration time.Duration) {
// }

func Set(key string, value []byte, expiration time.Duration) {
	client := getConnection()
	defer client.Close()

	err := client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		fmt.Printf("Cache operation failed: %v\n", err)
	}
}

//	func Get(key string) (string, error) {
//	    return "", nil
//	}
func Get(key string) (string, error) {
	client := getConnection()
	defer client.Close()

	cache, err := client.Get(context.Background(), key).Result()
	if err == nil {
		return cache, nil
	} else if err != redis.Nil {
		fmt.Println(err)
		return "", err
	}
	return "", errors.New("not found")
}
