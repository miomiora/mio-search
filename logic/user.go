package logic

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"mio-search/dao/mysql"
	"mio-search/dao/redis"
	"mio-search/model"
	"mio-search/util"
)

var (
	ErrorUserExist = errors.New("用户已存在")
)

func UserLogin(p *model.UserDTOLogin) (*model.UserVO, error) {
	// 1、从数据库中校验用户名和密码是否正确
	user, err := mysql.UserLogin(p.Account, p.Password)
	if err != nil {
		return nil, err
	}

	// 2、登录成功，把 Token 存入 Redis 中
	token := uuid.NewString()
	err = redis.InsertTokenByUserId(token, user.UserId)
	if err != nil {
		return nil, err
	}

	// 3、返回用户数据给 controller
	return &model.UserVO{
		UserId:      user.UserId,
		Account:     user.Account,
		Token:       &token,
		Description: nil,
		UserRole:    user.UserRole,
	}, nil
}

func UserRegister(u *model.UserDTORegister) (err error) {
	// 判断用户存不存在
	if mysql.CheckUserExist(u.Account) {
		zap.L().Error("[logic user] check user exist ")
		return ErrorUserExist
	}
	// 生成userId
	userID := util.GenSnowflakeID()
	// 构造一个User实例
	user := &model.User{
		UserId:   userID,
		Account:  u.Account,
		Password: u.Password,
	}
	// 保存进数据库
	err = mysql.InsertUser(user)
	return
}

func GetUserList(p *model.SearchDTO) ([]*model.UserVO, error) {
	data, err := mysql.QueryUserListByText(p)
	if err != nil {
		return nil, err
	}

	var userList []*model.UserVO

	for _, value := range data {

		user := &model.UserVO{
			UserId:      value.UserId,
			Account:     value.Account,
			Token:       nil,
			Description: value.Description,
			UserRole:    value.UserRole,
		}

		userList = append(userList, user)
	}
	return userList, nil
}
