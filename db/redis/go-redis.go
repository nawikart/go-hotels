package redis

import (
	// "fmt"
	"github.com/go-redis/redis"
	)

func ZScan2(k string, cs int, p string, count int64)  []string{
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	res, _, _ := client.ZScan(k, 0, p, count).Result()
	return res
}