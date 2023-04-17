package db

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// Redis Redis配置结构体
type Redis struct {
	Address     string
	Port     string
	Password string
	Database int
}

var rdb *redis.Client
var ctx = context.Background()

func RedisConn(config Redis) {
	addr := config.Address + ":" + config.Port
	password := ""
	if config.Password != "" {
		password = config.Password
	}
	database := 0
	if config.Database != 0 {
		database = config.Database
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       database, // use default DB
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