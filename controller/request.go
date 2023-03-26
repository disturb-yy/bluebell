package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrUSerNotLogin = errors.New("用户未登录")

// getCurrentUser 获取当前登录的用户ID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		return
	}
	userID, ok = uid.(int64)
	if !ok {
		err = ErrUSerNotLogin
		return
	}
	return
}

// getPageInfo 获取分页参数
func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("offset")
	sizeStr := c.Query("limit")

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
