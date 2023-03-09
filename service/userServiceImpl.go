package service

import (
	"crypto/sha256"
	"fmt"
	"simpleTiktok/dao"
)

type UserServiceImpl struct{}

// 用户注册服务,返回用户id
func (u UserServiceImpl) RegisterSrv(userName, password string) (int64, error) {
	secretpwd := Encode(password)
	return dao.CreateUser(userName, secretpwd)
}

// sha256加密
func Encode(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%x", sum)
}

// 比较加密后的密码和数据库中存储的密码
func (u UserServiceImpl) LoginSrv(userName, password string) (int64, error) {
	secretpwd := Encode(password)
	return dao.CheckPassword(userName, secretpwd)
}

func (u UserServiceImpl) BaseInfoSrv(userId, queryedUserId int64) (dao.User, error) {
	return dao.GetUserInfo(userId, queryedUserId)
}
