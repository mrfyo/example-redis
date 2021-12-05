package model

import (
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mrfyo/example-redis/util"
)

var (
	UserRecordKey = "record:user"
)

type User struct {
	ID       int    `json:"id"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Time     int64  `json:"time"`
}

//
// 域方法
//

func (User) TableName() string {
	return "user"
}

func (user *User) KeyName() string {
	return util.KeyGenerate(user.TableName(), strconv.Itoa(user.ID))
}

func (user *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       user.ID,
		"nickname": user.Nickname,
		"username": user.Username,
		"password": user.Password,
		"time":     user.Time,
	}
}

//
// 持久层
//

func CreateUser(user *User) (err error) {
	ID, err := nextID(user.TableName())
	if err != nil {
		return
	}
	user.ID = ID
	user.Time = time.Now().Unix()

	_, err = redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := user.KeyName()
		err = redisDB.HSet(ctx, key, user.ToMap()).Err()

		err = redisDB.ZAdd(ctx, UserRecordKey, &redis.Z{
			Score:  float64(user.Time),
			Member: key,
		}).Err()

		return err
	})

	return
}

func RemoveUser(user *User) (err error) {
	redisDB.Pipelined(ctx, func(p redis.Pipeliner) error {
		key := user.KeyName()

		p.Del(ctx, key)
		p.ZRem(ctx, UserRecordKey, key)

		return nil
	})

	return
}

func UpdateUser(user *User) (err error) {
	key := user.KeyName()
	intCmd := redisDB.HSet(ctx, key, user.ToMap())
	if err := intCmd.Err(); err != nil {
		log.Println(err)
	}
	return
}

func GetUser(id int) (user *User, err error) {
	user = &User{
		ID: id,
	}
	cmd := redisDB.HGetAll(ctx, user.KeyName())
	if err = cmd.Err(); err != nil {
		log.Println(err)
		return
	}
	m := cmd.Val()

	user = &User{
		ID:       id,
		Nickname: m["nickname"],
		Username: m["username"],
		Password: m["password"],
	}
	return
}

func GetAllUserByPage(offset, limit int) (users []*User) {

	keys, err := redisDB.ZRevRangeByScore(ctx, UserRecordKey, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: int64(offset),
		Count:  int64(limit),
	}).Result()
	if err != nil {
		return
	}

	ids, err := util.BatchExtraID(keys)
	if err != nil {
		return
	}

	for _, ID := range ids {
		user, err := GetUser(ID)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return
}

func CountUser() (count int64) {
	i, err := redisDB.ZCount(ctx, UserRecordKey, "-inf", "+inf").Result()
	if err != nil {
		return 0
	}
	return i
}
