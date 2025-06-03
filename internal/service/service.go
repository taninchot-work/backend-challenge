package service

import "github.com/taninchot-work/backend-challenge/internal/repository"

type Service struct {
	ServerService ServerService
	UserService   UserService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		ServerService: NewServerService(),
		UserService:   NewUserService(repository.UserRepository),
	}
}
