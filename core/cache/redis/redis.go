package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

func Init(opts ...Option) {
	options := NewOptions(opts...)
	rdb = redis.NewClient(&redis.Options{
		Addr:     options.Address + ":" + options.Port,
		Password: options.Password, // no password set
		DB:       options.Database, // use default DB
	})
	rdb.Conn(ctx)
	pong := rdb.Ping(ctx)
	if pong.Val() == "PONG" {
		log.Println("redis connect success")
	} else {
		log.Panicln("redis connect failure")
	}
}

func RedisCache(args ...interface{}) interface{} {
	res, _ := rdb.Do(ctx, args...).Result()
	return res
}

func RedisClose() {
	rdb.Close()
}
