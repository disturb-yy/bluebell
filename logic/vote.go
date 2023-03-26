package logic

import (
	"strconv"

	"go.uber.org/zap"

	"github.com/disturb-yy/bluebell/dao/redis"
	"github.com/disturb-yy/bluebell/models"
)

// 投票功能：
// 1、用户投票的数据（参数放在controller）

// VoteForPost 处理用户为帖子投票的函数
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
