package service

type UserService interface {
	RegisterSrv(string, string) (int64, error)
	LoginSrv(string, string) (int64, error)
	BaseInfoSrv(int64, int64) (string, error)
}
