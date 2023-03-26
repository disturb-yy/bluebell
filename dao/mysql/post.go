package mysql

import (
	"strings"

	"github.com/disturb-yy/bluebell/models"
	"github.com/jmoiron/sqlx"
)

func CreatePost(p *models.Post) (err error) {
	// 插入操作
	sqlStr := `insert into post 
			(post_id, title, content, author_id, community_id)
			VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return err
}

// GetPostByPid 根据 pid 在数据库查询帖子详情
func GetPostByPid(pid int64) (post *models.Post, err error) {
	// 0. 定义保存的变量
	post = new(models.Post)
	// 1. 数据库查询语句
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
	from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetPostList 获取帖子列表（按创建时间排序）
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	// 查询数据库
	// 首先想到的使用create_time进行排序，但由于帖子是递增的，因此可以使用 postid 排序
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
		from post order by create_time DESC limit ?, ?`
	// 返回数据
	posts = make([]*models.Post, 0, 2) // 不能写make([]*models.Post, 2)，因为这样不能用append
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

// GetPostListByIDs 根据给定的 id 列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	// sql 语句
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
			from post where post_id in (?) order by FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)

	err = db.Select(&postList, query, args...)
	return
}
