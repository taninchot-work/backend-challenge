package controller

import (
	"github.com/taninchot-work/backend-challenge/internal/core/middleware"
	"github.com/taninchot-work/backend-challenge/internal/service"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, svc *service.Service) {
	serverController := NewServerController(svc.ServerService)
	userController := NewUserController(svc.UserService)

	mux.HandleFunc("GET /health", serverController.HealthCheck)

	// user routes
	mux.HandleFunc("GET /api/v1/users/get/me", middleware.JwtMiddleware(userController.GetMe)) // protected route
	mux.HandleFunc("GET /api/v1/users/get/list", userController.UserListGet)
	mux.HandleFunc("POST /api/v1/users/register", userController.UserRegister)
	mux.HandleFunc("POST /api/v1/users/login", userController.UserLogin)
	mux.HandleFunc("POST /api/v1/users/update", middleware.JwtMiddleware(userController.UserUpdate)) // protected route
	mux.HandleFunc("POST /api/v1/users/delete", middleware.JwtMiddleware(userController.UserDelete)) // protected route

}
