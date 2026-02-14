package service

type AdminService struct{}

func NewAdminService() *AdminService {
	a := &AdminService{}
	a.updRSAppListFromDB()

	return a
}
