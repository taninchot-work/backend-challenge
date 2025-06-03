package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/taninchot-work/backend-challenge/internal/repository"
	"github.com/taninchot-work/backend-challenge/internal/service/background"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/taninchot-work/backend-challenge/internal/controller"
	"github.com/taninchot-work/backend-challenge/internal/core/config"
	"github.com/taninchot-work/backend-challenge/internal/core/db"
	"github.com/taninchot-work/backend-challenge/internal/core/middleware"
	"github.com/taninchot-work/backend-challenge/internal/service"
)

func main() {
	// initialize configuration
	config.InitConfig("config.yaml")

	// load configuration
	cfg := config.GetConfig()

	// initialize db
	err := db.InitializeMongoDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// initialize mux
	mux := http.NewServeMux()

	// initialize repository
	repositories := repository.NewRepository()

	// initialize services
	svc := service.NewService(repositories)

	// register routes
	controller.RegisterRoutes(mux, svc)

	// middlewares
	var handler http.Handler = mux
	handler = middleware.LoggingMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)

	// create configure http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.RestServer.Port),
		Handler: handler,
	}

	// setup signal handling
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// counting user
	go background.StartUserCountLogger(ctx, 10*time.Second)

	// start the server in a goroutine
	go func() {
		log.Printf("Starting REST Server on port %d\n", cfg.RestServer.Port)
		log.Printf("Local : http://localhost:%d\n", cfg.RestServer.Port)
		log.Println("waiting for request...")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("failed to served %s\n", err)
			stop()
		}
	}()

	// wait for the context to be canceled (i.e., SIGINT or SIGTERM)
	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server shutdown failed: %v\n", err)
	}

	log.Println("Server gracefully stopped")
}
