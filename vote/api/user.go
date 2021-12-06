package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mrfyo/example-redis/model"
	"github.com/mrfyo/example-redis/result"
	"github.com/mrfyo/example-redis/util"
)



func addUserHandler(c *gin.Context) {

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		result.Fail(c, 1, "form struct error")
		return
	}

	if util.AnyEmptyStr(user.Nickname, user.Username, user.Password) {
		result.Fail(c, 1, "form value error")
		return
	}
	
	if err := model.CreateUser(&user); err != nil {
		result.Fail(c, 1, "create fail")
		return
	}

	result.Success(c, user.ToMap())
}

func deleteUserHandler(c *gin.Context) {

	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil || ID <= 0 {
		result.Fail(c, 1, "Path Param error: id")
		return
	}

	user, err := model.GetUser(ID)
	if err != nil {
		result.Fail(c, 10, "user not exist")
		return
	}

	if err := model.RemoveUser(user); err != nil {
		result.Fail(c, 10, "Delete User Fail")
		return
	}

	result.Success(c, nil)
}

func listUserHandler(c *gin.Context) {
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		result.Fail(c, 2, "Query Param Error: offset")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if err != nil || limit <= 0 {
		result.Fail(c, 2, "Query Param Error: limit")
		return
	}

	users := model.GetAllUserByPage(offset, limit)
	total := model.CountUser()

	result.Success(c, gin.H{
		"total": total,
		"items": users,
	})
}
