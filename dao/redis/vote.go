package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// 投票分数计算
// 投一票加432分 86400/200 ——> 需要200张赞成票可以给你的帖子续一天在首页

/* 投票的几种情况
direction = 1 时，有两种情况
	1、之前没有投过票，现在投赞成票  --> 更新分数和投票记录， 差值的绝对值： 1 +432
	2、之前投反对票，现在投赞成票 --> 更新分数和投票记录  差值的绝对值： 2  +432 * 2
direction = 0 时，有两种情况
	1、之前投反对票，现在取消投票  --> 更新分数和投票记录  差值的绝对值： 1  +432
	2、之前投赞成票，现在取消投票  --> 更新分数和投票记录  差值的绝对值： 1  -432
direction = -1 时，有两种情况  --> 更新分数和投票记录
	1、之前没有投过票，现在投反对票  --> 更新分数和投票记录  差值的绝对值： 1  -432
	2、之前投赞成票，现在投反对票  --> 更新分数和投票记录  差值的绝对值： 2  -432*2

投票的限制：
	每个帖子自发表之日起，允许用户投票，超过一个星期就不允许投票
	1、到期之后将redis中保存的赞成票和反对票存储到mysql中
	2、到期之后删除那个 KeyPostVotedPrefix
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepested   = errors.New("不允许重复投票")
)

// CreatePost 在 redis 中创建帖子
func CreatePost(postID, communityID int64) (err error) {
	// 事务操作
	pipeline := rdb.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScore), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 把帖子 id 加到社区的 set 中
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ctx, cKey, postID)
	_, err = pipeline.Exec(ctx)
	return err
}

// VoteForPost 传入用户ID，帖子ID，投赞成、反对还是取消
func VoteForPost(userID, postID string, value float64) (err error) {
	// 1. 判断投票的限制
	// 使用ZScore，根据 postID 从有序集合 KeyPostTime 中得到其对应值
	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTime), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		// 帖子时间大于一周
		return ErrVoteTimeExpire
	}
	// 2. 更新分数+++
	// 查询当前用户给当前帖子之前的投票记录 -1 0 1
	ov := rdb.ZScore(ctx, getRedisKey(KeyPostVotedPrefix+postID), userID).Val()
	// 如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepested
	}
	// 存储帖子分数加减
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值
	// 2，3是要放在一个事务里面操作
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScore), op*diff*scorePerVote, postID)
	// 3. 记录用户为该帖子投过票
	if value == 0 {
		// 取消投票，等价于删除用户
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedPrefix+postID), postID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value, // 赞成票还是返回票
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)

	return
}
