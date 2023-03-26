package controller

import (
	"github.com/disturb-yy/bluebell/logic"
	"github.com/disturb-yy/bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票

// PostVoteHandler 投票处理的函数
func PostVoteHandler(c *gin.Context) {
	// 1、 参数校验
	u := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(u); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam)
		}
		errData := removeTopStruct(errs.Translate(trans)) // 翻译并去除错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	// 获取当前的用户 id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 交给logic处理业务逻辑，即具体的投票业务逻辑
	if err = logic.VoteForPost(userID, u); err != nil {
		zap.L().Error("logic.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
