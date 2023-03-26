package logic

import (
	"github.com/disturb-yy/bluebell/dao/mysql"
	"github.com/disturb-yy/bluebell/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	// 查数据库 查找到所有的community，并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	// 根据给定的id查询对应的community
	return mysql.GetCommunityDetail(id)
}
