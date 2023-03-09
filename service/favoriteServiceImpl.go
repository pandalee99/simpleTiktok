package service

import (
	"simpleTiktok/dao"
	"strconv"
)

type FavoriteServiceImpl struct {
	Host string
}

// 点赞服务,用户给视频点赞/取消点赞
//
// actionType 1 点赞 2 取消点赞
func (f FavoriteServiceImpl) LikeAction(userId int64, videoIdStr string, actionTypeStr string) error {
	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil || actionType > 2 || actionType < 1 {
		return err
	}
	return dao.LikeAction(userId, videoIdStr, actionType == 1)
}

func (f FavoriteServiceImpl) ListAction(userId int64) ([]dao.Video, error) {
	return dao.LikeList(userId, f.Host)
}
