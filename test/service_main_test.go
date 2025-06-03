package test

import (
	"flag"
	"github.com/taninchot-work/backend-challenge/internal/core/config"
	"os"
	"testing"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	flag.Parse()
	setup()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setup() {
	cfg = &config.Config{
		RestServer: config.RestServer{
			Port: 8080,
			Jwt: config.JwtConfig{
				Secret:   "SECRET_KEY",
				ExpireIn: 100000,
				Issuer:   "BACKEND_CHALLENGE",
			},
		},
	}
	config.SetConfig(cfg)
}
