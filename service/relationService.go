package service

import "simpleTiktok/dao"

type RelationService interface {
	Action(int64, string, string) error
	MasterList(string, int64) ([]dao.User, error)
	FollowerList(string, int64) ([]dao.User, error)
}
