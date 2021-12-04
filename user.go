package main

import (
	"log"
	"strconv"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (User) TableName() string {
	return "user"
}

func (user *User) KeyName() string {
	return KeyGenerate(user.TableName(), strconv.Itoa(user.ID))
}

func (user *User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
	}
}

func CreateUser(user *User) (err error) {
	ID, err := NextID(user.TableName())
	if err != nil {
		return
	}
	user.ID = ID
	key := user.KeyName()
	intCmd := redisDB.HSet(ctx, key, user.ToMap())
	if err := intCmd.Err(); err != nil {
		log.Println(err)
	}
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
	user = new(User)
	key := KeyGenerate(user.TableName(), strconv.Itoa(id))
	cmd := redisDB.HGetAll(ctx, key)
	if err = cmd.Err(); err != nil {
		log.Println(err)
		return
	}
	m := cmd.Val()

	user.ID, err = strconv.Atoi(m["id"])
	user.Name = m["name"]
	return
}
