package controller

import (
	"errors"

	"github.com/disturb-yy/bluebell/dao/mysql"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/disturb-yy/bluebell/models"

	"github.com/disturb-yy/bluebell/logic"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Signup with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 翻译错误
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 使用bingding Tag 后，ShouldBindJSON会自动进行校验
	//// 手动对请求参数进行详细的业务规则校验
	//if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	// 请求参数有误，直接返回响应
	//	zap.L().Error("Signup with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}

	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			// 如果错误类型是mysql.ErrorUserExist
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 处理用户登录的函数
func LoginHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParam)
			return
		}
		// 翻译错误
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2. 业务处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("login.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNoExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   user.UserID, // id 大于 1<<53-1，而int64 最大值 1<<63-1
		"user_name": user.Username,
		"token":     user.Token,
	})
}
