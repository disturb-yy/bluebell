package mysql

import (
	"testing"

	"github.com/disturb-yy/bluebell/settings"

	"github.com/disturb-yy/bluebell/models"
)

func init() {
	dbCfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "8357",
		DbName:       "bluebell",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := &models.Post{
		ID:          10,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err: %v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
