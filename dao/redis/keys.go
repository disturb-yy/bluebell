package redis

// redis key 设置 redis 的 key
// redis key 注意使用命名空间的方式，方便业务拆分

const (
	Prefix             = "bluebell:"   // 项目公用前缀，便于快速查找
	KeyPostTime        = "post:time"   // zset: 帖子及发帖时间
	KeyPostScore       = "post:score"  // zset: 帖子及投票的分数
	KeyPostVotedPrefix = "post:voted:" // zset: 记录用户及投票类型; 参数是Post ID
	KeyCommunitySetPF  = "community:"  // set; 保存每个分区帖子的id
)

// 给 reids key 加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
