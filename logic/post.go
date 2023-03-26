package logic

import (
	"github.com/disturb-yy/bluebell/dao/mysql"
	"github.com/disturb-yy/bluebell/dao/redis"
	"github.com/disturb-yy/bluebell/models"
	"github.com/disturb-yy/bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成帖子ID
	p.ID = snowflake.GenID()

	// 2. 保存到数据库
	if err = mysql.CreatePost(p); err != nil {
		return
	}
	return redis.CreatePost(p.ID, p.CommunityID)
	// 3. 返回响应
}

// GetPostByID 根据帖子 id 查询帖子详情数据
func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	// 1. 根据 pid 查询数据库
	post, err := mysql.GetPostByPid(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByPid(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	// 2. 根据作者 id 查询作者id
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	// 3. 根据社区 id 查询社区详情
	community, err := mysql.GetCommunityDetail(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetail(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}

	// 4. 组合接口想用的数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}

	//  5. 返回查询到的数据
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	// 查询数据库
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	// 返回数据
	for _, post := range posts {
		// 根据作者 id 查询作者id
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// 根据社区 id 查询社区详情
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		// 44. 组合接口想用的数据
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2. 去 redis 查询 id 列表
	// 3. 根据 id 去数据库查询帖子详细信息
	ids, err := redis.GetPostIDs(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDs(p) return 0 data")
		return
	}
	// 根据 id 去 mysql 查询
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息
	for idx, post := range posts {
		// 根据作者 id 查询用户id和用户名
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// 根据社区 id 查询社区详情
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		// 44. 组合接口想用的数据
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VotesNum:        voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetCommunityPostList 按社区id获取帖子
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去 redis 查询 id

	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("GetCommunityPostList2 return 0 data")
		return
	}

	// 根据 id 去 mysql 查询
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息
	for idx, post := range posts {
		// 根据作者 id 查询作者id
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		// 根据社区 id 查询社区详情
		community, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		// 44. 组合接口想用的数据
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VotesNum:        voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostListNew 将两个查询帖子列表的接口合二为一的函数
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 获取数据
	if p.CommunityID == 0 {
		// 查所有帖子
		data, err = GetPostList2(p)
	} else {
		// 根据社区 id 查对应帖子
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
