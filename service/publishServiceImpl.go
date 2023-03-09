package service

import "simpleTiktok/dao"

type PublishServiceImpl struct {
	Host string
}

func (p PublishServiceImpl) SavePOSTFile(content []byte, path, fileName, title string, userId int64) error {
	return dao.SavePOSTFile(content, path, fileName, title, userId)
}

func (p PublishServiceImpl) GetVideoOfUser(userId, jwtUserId int64) ([]dao.Video, error) {
	return dao.GetVideoOfUser(userId, jwtUserId, p.Host)
}

func (p PublishServiceImpl) GetFeed(latestTime int64, jwtUserId int64) ([]dao.Video, int64, error) {
	return dao.GetFeed(latestTime, jwtUserId, p.Host)
}
