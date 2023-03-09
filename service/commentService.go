package service

import "simpleTiktok/dao"

type CommentService interface {
	CommentAction(int64, string, string, string, string) (dao.Comment, error)
	CommentList(string, int64) ([]dao.Comment, error)
}
