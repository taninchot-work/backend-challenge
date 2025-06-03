package repository

import (
	"github.com/taninchot-work/backend-challenge/internal/core/db"
)

type Repository struct {
	UserRepository UserRepository
}

func NewRepository() *Repository {
	mongoDatabase := db.GetDatabase()
	return &Repository{
		UserRepository: NewUserRepository(mongoDatabase.Collection("users")),
	}
}
