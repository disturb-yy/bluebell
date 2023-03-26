package mysql

import (
	"fmt"

	"github.com/disturb-yy/bluebell/settings"

	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// Init() 初始化数据库连接

func Init(cfg *settings.MySQLConfig) (err error) {
	// 从config.yaml加载配置
	// database source name
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		zap.L().Error("Init MySQL failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

func Close() {
	_ = db.Close()
}
