package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"github.com/disturb-yy/bluebell/models"
)

/*  处理用户的相关操作
把每一次数据库操作封装成函数
待logic层根据业务需要进行调用
*/

const secret = "yangmi"

// Login 将用户登录的密码与数据库中的密码进行比对
func Login(user *models.User) (err error) {
	oPassword := user.Password // 保存用户登录密码
	// 执行查询语句
	sqlStr := `select user_id, username, password from user where username = ?`
	if err = db.Get(user, sqlStr, user.Username); err != nil {
		return
	}
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	// 判断密码是否正确
	passwd := encryptPassword(oPassword) // 对用户登录密码进行加密
	if passwd != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// InsertUser 向数据库插入用户一条新的数据
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	passwd := encryptPassword(user.Password)
	// 执行SQL入库
	sqlStr := `insert into user(user_id, username, password) values(?, ?, ?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, passwd)
	return
}

// CheckUserExist 检查指定用户名的用户是否已存在
func CheckUserExist(username string) (err error) {
	// 执行SQL查询
	sqlStr := "SELECT COUNT(user_id) FROM user WHERE username = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// GetUserById 根据 id 获取用户信息
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ? `
	err = db.Get(user, sqlStr, uid)
	return
}
