package logic

import (
	"errors"

	"github.com/disturb-yy/bluebell/pkg/jwt"
	"go.uber.org/zap"

	"github.com/disturb-yy/bluebell/dao/mysql"
	"github.com/disturb-yy/bluebell/models"
	"github.com/disturb-yy/bluebell/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 1. 判断用户是否存在
	if err = mysql.CheckUserExist(p.Username); err != nil {
		// 数据查询出错
		return
	}

	// 2. 生成用户Uid
	userID := snowflake.GenID()
	// 构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 3. 保存到数据库
	if err = mysql.InsertUser(user); err != nil {
		return errors.New("插入用户数据失败")
	}
	return
}

// Login 函数用户处理函数登录的业务
func Login(p *models.ParamLogin) (user *models.User, err error) {
	// 1. 用户是否已注册 不用检查，直接查询即可
	//if err = mysql.CheckUserExist(p.Username); err == nil {
	//	return errors.New("用户未注册")
	//}
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	// 2. 为已注册用户进行密码校验
	if err = mysql.Login(user); err != nil {
		// 登录失败
		return nil, err
	}
	// 登录，并生成JWT返回
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		zap.L().Error("jwt.GenToken(user.UserID, user.Username) failed", zap.Error(err))
		return
	}
	user.Token = token
	return
}
