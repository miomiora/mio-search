package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"go.uber.org/zap"
	"mio-search/model"
)

const secret = "https://github.com/miomiora"

func InsertUser(u *model.User) (err error) {
	u.Password = encryptPassword(u.Password)
	err = db.Create(u).Error
	if err != nil {
		zap.L().Error("[dao mysql user] insert user error ", zap.Error(err))
	}
	return
}

func UserLogin(account, password string) (*model.User, error) {
	u := new(model.User)
	err := db.First(u, "account = ? and password = ?", account, encryptPassword(password)).Error
	if err != nil {
		zap.L().Error("[dao mysql user] user login error ", zap.Error(err))
		return nil, err
	}
	return u, err
}

func QueryUserByAccount(account string) (*model.User, error) {
	u := new(model.User)
	err := db.First(u, "account = ?", account).Error
	if err != nil {
		zap.L().Error("[dao mysql user] query user by account error ", zap.Error(err))
		return nil, err
	}
	return u, nil
}

func QueryUserByUserId(id int64) (*model.User, error) {
	u := new(model.User)
	err := db.First(u, "user_id = ?", id).Error
	if err != nil {
		zap.L().Error("[dao mysql user] query user by userId error ", zap.Error(err))
		return nil, err
	}
	return u, nil
}

func QueryUserList(page, size int) ([]*model.User, error) {
	var u []*model.User
	err := db.Limit(size).Offset(page - 1).Find(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func QueryUserListByText(p *model.SearchDTO) ([]*model.User, error) {
	var u []*model.User
	err := db.Limit(p.Size).Offset(p.Page-1).Where("account like ? or description like ?", "%"+p.Text+"%", "%"+p.Text+"%").Find(&u).Error
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CheckUserExist(account string) bool {
	if _, err := QueryUserByAccount(account); err != nil {
		return false
	}
	return true
}

// 密码加密
func encryptPassword(originPassword string) string {
	hash := md5.New()
	hash.Write([]byte(secret))
	hash.Write([]byte(originPassword))
	encryptString := hex.EncodeToString(hash.Sum(nil))
	return encryptString
}
