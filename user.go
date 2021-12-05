package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mrfyo/example-redis/result"
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
// 领域方法
//

func (User) TableName() string {
	return "user"
}

func (user *User) KeyName() string {
	return KeyGenerate(user.TableName(), strconv.Itoa(user.ID))
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
	ID, err := NextID(user.TableName())
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

func UpdateUser(user *User) (err error) {
	key := user.KeyName()
	intCmd := redisDB.HSet(ctx, key, user.ToMap())
	if err := intCmd.Err(); err != nil {
		log.Println(err)
	}
	return
}

func RemoveUser(user *User) (err error) {
	key := user.KeyName()
	intCmd := redisDB.Del(ctx, key)
	if err = intCmd.Err(); err != nil {
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

func GetAllUser() (users []*User) {

	keys, err := redisDB.ZRangeByScore(ctx, UserRecordKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   "+inf",
		Count: 10,
	}).Result()
	if err != nil {
		return
	}

	for _, key := range keys {
		ID, err := ExtraID(key)
		if err != nil {
			continue
		}
		user, err := GetUser(ID)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return
}

//
// User API
//

func AddUser(c *gin.Context) {

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		result.Fail(c, 1, "form struct error")
		return
	}

	if AnyEmptyStr(user.Nickname, user.Username, user.Password) {
		result.Fail(c, 1, "form value error")
		return
	}
	if err := CreateUser(&user); err != nil {
		result.Fail(c, 1, "create fail")
		return
	}

	result.Success(c, user.ToMap())
}

func DeleteUser(c *gin.Context) {

	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ID <= 0 {
		result.Fail(c, 1, "Path Param error: id")
		return
	}

	user, err := GetUser(ID)
	if err != nil {
		result.Fail(c, 10, "user not exist")
		return
	}

	if err := RemoveUser(user); err != nil {
		result.Fail(c, 10, "Delete User Fail")
		return
	}

	result.Success(c, nil)
}

func ListUser(c *gin.Context) {

	users := GetAllUser()

	list := make([]map[string]interface{}, 0, len(users))

	for _, user := range users {
		list = append(list, user.ToMap())
	}

	result.Success(c, gin.H{
		"total": len(list),
		"items": list,
	})
}
