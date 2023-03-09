package service

import "simpleTiktok/dao"

type FavoriteService interface {
	LikeAction(int64, string, string) error
	ListAction(int64) ([]dao.Video, error)
}
