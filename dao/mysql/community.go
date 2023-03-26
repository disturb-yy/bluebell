package mysql

import (
	"database/sql"

	"github.com/disturb-yy/bluebell/models"
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	// 数据库查询语句
	sqlStr := `select community_id, community_name from community`
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows { // 空行无数据，要进行处理，但不返回错误
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityDetail 根据ID查询社区详情
func GetCommunityDetail(id int64) (community *models.CommunityDetail, err error) {
	// 数据库查询语句
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time 
				from community 
				where community_id = ?`
	if err = db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			// 空行无数据，要进行处理，但不返回错误
			zap.L().Warn("there is no community in db")
			err = ErrorInvalidID
		}
	}
	return
}
