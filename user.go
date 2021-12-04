package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrfyo/example-redis/result"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
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
		"name":     user.Name,
		"username": user.Username,
		"password": user.Password,
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
		Name:     m["name"],
		Username: m["username"],
		Password: m["password"],
	}
	return
}

func GetAllUser() (users []*User) {

	cmd := redisDB.Scan(ctx, 0, "user", 0)

	keys, _ := cmd.Val()

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

	if AnyEmptyStr(user.Name, user.Username, user.Password) {
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
