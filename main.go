package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func main() {
	fmt.Println("vim-go")

	var ctx = context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "pinode:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)

	value, err := client.Get(ctx, "key").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(value)
}
