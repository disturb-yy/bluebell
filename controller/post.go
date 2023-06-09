package controller

import (
	"strconv"

	"github.com/disturb-yy/bluebell/logic"
	"github.com/disturb-yy/bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子的处理函数
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数的校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("post with invalid param", zap.Error(err))
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

	// 2. 从 c 获取当前发请求的用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 3. 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取 url 的参数（帖子的id）
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 根据对应的 id 去数据库取出数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByID(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 按时间排序获取帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// @Summary 升级版帖子列表接口
// @Description 可按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
// 根据前端传来的参数(创建时间，分数)动态的获取帖子列表
// 1. 获取请求的Query string参数
// 2. 去 redis 查询 id 列表
// 3. 根据 id 去数据库查询帖子详细信息
func GetPostListHandler2(c *gin.Context) {
	// GET请求参数(query string)： /api/v1/posts2?page=10&order=time
	p := &models.ParamPostList{
		CommunityID: 0,
		Page:        1,
		Size:        10,
		Order:       models.OrderTime,
	}
	// 获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		// 根据请求的数据类型选择相应的方法（这里Query）去获取数据
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//page, size := getPageInfo(c)  // 使用from tag 获取， 不用自定义函数

	// 获取数据
	data, err := logic.GetPostListNew(p) // 更新：合二为一
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetCommunityPostListHandler 根据社区去查询帖子列表
//func GetCommunityPostListHandler(c *gin.Context) {
//	// 获取 Query string
//	p := &models.ParamPostList{
//		Page:        1,
//		Size:        10,
//		Order:       models.OrderTime,
//		CommunityID: 0,
//	}
//
//	// 获取请求中的数据
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid params", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	// 获取数据
//	data, err := logic.GetPostListNew(p)
//	if err != nil {
//		zap.L().Error("logic.GetCommunityPostList2() failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	ResponseSuccess(c, data)
//}
