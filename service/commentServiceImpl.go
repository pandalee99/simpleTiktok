package service

import (
	"errors"
	"simpleTiktok/dao"
	"strconv"
)

type CommentServiceImpl struct {
	Host string
}

// 评论服务
//
// actionTypeStr 1 发布评论 2 删除评论
func (c CommentServiceImpl) CommentAction(userId int64, videoIdStr string, actionTypeStr string, commentText string, commentIdStr string) (dao.Comment, error) {
	switch actionTypeStr {
	case "1": // 发布评论
		videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
		if err != nil {
			return dao.Comment{}, err
		}
		return dao.PublishComment(userId, videoId, commentText)
	case "2": // 删除评论
		commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			return dao.Comment{}, err
		}
		return dao.Comment{}, dao.DeleteComment(userId, commentId)
	default:
		return dao.Comment{}, errors.New("unknown action type")
	}
}

// 获取视频下所有评论服务
func (c CommentServiceImpl) CommentList(videoIdStr string, jwtUserId int64) ([]dao.Comment, error) {
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		return make([]dao.Comment, 0), err
	}
	return dao.GetCommentDescByPublishTime(videoId, jwtUserId)
}
