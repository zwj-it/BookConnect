package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const ContextUserIDKey = "userID"

var ErrorUserNoLogin = errors.New("用户未登录")

// getCurrentUser 获取当前登陆的用户id，常在controller调用
func getCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIDKey) //ok是一个bool值
	if !ok {
		err = ErrorUserNoLogin
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNoLogin
		return
	}
	return
}
func getPageInfo(c *gin.Context) (int64, int64) {
	//page是从哪开始，size是要查多少条
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var (
		page int64
		size int64
		err  error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
