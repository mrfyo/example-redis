package model

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	redisDB *redis.Client
	ctx     context.Context
)

func InitModel(db *redis.Client, c context.Context) {
	redisDB = db
	ctx = c
}

func nextID(name string) (ID int, err error) {
	id, err := redisDB.HIncrBy(ctx, "counter", name, 1).Result()
	if err != nil {
		return
	}
	ID = int(id)
	return
}
