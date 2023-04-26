package model

import "gorm.io/gorm"

// User 用户
type User struct {
	gorm.Model

	UserId      int64   `json:"user_id" gorm:"not null"`
	Account     string  `json:"account" gorm:"not null"`
	Password    string  `json:"password" gorm:"not null"`
	Description *string `json:"description"`
	UserRole    bool    `json:"user_role" gorm:"default:0;size:false"`
}

// UserVO 登录成功返回给前端展示的用户数据
type UserVO struct {
	UserId      int64   `json:"user_id,string"`
	Account     string  `json:"account"`
	Token       *string `json:"token,omitempty"`
	Description *string `json:"description,omitempty"`
	UserRole    bool    `json:"user_role"`
}

// UserDTOLogin 用户登录所需要绑定的参数
type UserDTOLogin struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserDTORegister 用户注册所需要绑定的参数
type UserDTORegister struct {
	Account    string `json:"account" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}
