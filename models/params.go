package models

// 定义请求的参数结构体

const (
	OrderTime  = "time"  // 按照时间排序返回帖子列表
	OrderScore = "score" // 按照分数排序返回帖子列表
)

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 存储用户投票的功能
type ParamVoteData struct {
	// UserID int64  // 由于用户id可以从url处获取，因此可以不用定义
	PostID string `json:"post_id" binding:"required"` // 帖子id
	// oneof 代表使用 validator 使用的校验值只能是 0 -1 1
	Direction int8 `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票（1）还是反对票（-1）,取消投票（0）
}

// ParamPostList 获取帖子列表 query string 参数
type ParamPostList struct { // from tag 用于 Query string
	CommunityID int64  `json:"community_id" from:"community"` // 可以为空
	Page        int64  `json:"page" from:"page"`
	Size        int64  `json:"size" from:"size"`
	Order       string `json:"order" from:"order" example:"score"`
}
