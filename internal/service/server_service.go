package service

type ServerService interface {
	HealthCheck() string
}

type serverServiceImpl struct{}

func (s *serverServiceImpl) HealthCheck() string {
	return "OK"
}

func NewServerService() ServerService {
	return &serverServiceImpl{}
}
