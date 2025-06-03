package controller

import (
	"github.com/taninchot-work/backend-challenge/internal/service"
	"net/http"
)

type ServerController interface {
	HealthCheck(w http.ResponseWriter, r *http.Request)
}

type serverControllerImpl struct {
	serverService service.ServerService
}

func NewServerController(serverService service.ServerService) ServerController {
	return &serverControllerImpl{
		serverService: serverService,
	}
}

func (s *serverControllerImpl) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := s.serverService.HealthCheck()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
