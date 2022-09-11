package service

type userService struct {
	config Config
}

func NewCommodityService(config Config) *userService {

	return &userService{
		config: config,
	}
}
