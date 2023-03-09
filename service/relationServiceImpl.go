package service

import (
	"errors"
	"fmt"
	"simpleTiktok/dao"
	"strconv"
)

type RelationServiceImpl struct{}

// 关注或者取关
//
// ActionType 1 关注 2 取关
func (r RelationServiceImpl) Action(userId int64, toUserIdStr string, actionTypeStr string) error {
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		return err
	}

	if userId == toUserId {
		return errors.New("cannot follow yourself")
	}

	// 检查关注者id是否存在
	_, err = dao.GetUserByUserId(toUserId)
	if err != nil {
		return fmt.Errorf("err in master id: %v", err)
	}

	// 进行关注的相关操作
	switch actionTypeStr {
	case "1":
		return dao.FollowMaster(userId, toUserId)
	case "2":
		return dao.UnFollowMaster(userId, toUserId)
	default:
		return errors.New("unknown action type")
	}
}

// 获取当前用户关注的所有用户
func (r RelationServiceImpl) MasterList(userId string, jwtUserId int64) ([]dao.User, error) {
	return dao.GetMasterList(userId, jwtUserId)
}

// 获取 关注当前用户 的 所有用户
func (r RelationServiceImpl) FollowerList(masterId string, jwtUserId int64) ([]dao.User, error) {
	return dao.GetFollowerList(masterId, jwtUserId)
}
