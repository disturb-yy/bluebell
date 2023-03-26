package redis

import (
	"strconv"
	"time"

	"github.com/disturb-yy/bluebell/models"
	"github.com/redis/go-redis/v9"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	// 2. 确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1
	// 3. ZRevRange 按分数从大到小的顺序查询指定数量的元素
	return rdb.ZRevRange(ctx, key, start, end).Result()
}

// GetPostIDs 获取帖子 ID
func GetPostIDs(p *models.ParamPostList) ([]string, error) {
	// 从 redis 获取 id
	// 1. 根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScore)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据 ids 查询每篇帖子的赞成票数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// 进行 redis 查询
	data = make([]int64, 0, len(ids))
	// 统计每个帖子投赞成票的数量，即 1 的数量
	//for _, id := range ids {
	//	// 生成 redis key
	//	key := getRedisKey(KeyPostVotedPrefix + id)
	//	// 查找 key 中分数是 1 的元素的数量 -> 统计每篇帖子赞成票的数量
	//	cnt := rdb.ZCount(ctx, key, "1", "1").Val()
	//	data = append(data, cnt)
	//}
	// 使用 pipeline 一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder  根据社区id获取对应帖子列表id
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 使用 Zinterstore 把分区的帖子set与帖子分数的zset 生成一个新的 zset
	// 针对新的zset， 按之前的逻辑取数据
	orderKey := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScore)
	}
	// 社区的key - community:id
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	// 利用缓存 key 减少 zinterstore 执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if rdb.Exists(ctx, key).Val() < 1 { // 判断 key 是否存在
		// 不存在，需要计算
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{cKey, orderKey}, // 社区id和 orderkey 取交集，存储在key中
			Aggregate: "MAX",
		}) // 计算交集
		pipeline.Expire(ctx, key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	// 存在的话就直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
