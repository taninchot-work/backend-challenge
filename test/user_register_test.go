package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taninchot-work/backend-challenge/internal/core/util/jwt"
	"github.com/taninchot-work/backend-challenge/internal/dto"
	"github.com/taninchot-work/backend-challenge/internal/entity"
	mock_user_repository "github.com/taninchot-work/backend-challenge/internal/repository/mocks/user_repository_mock"
	"github.com/taninchot-work/backend-challenge/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
	"time"
)

func TestRegisterUserSuccess(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	req := dto.UserRegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	userID := bson.NewObjectID()
	userEntity := entity.User{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedAccessToken, err := jwt.GenerateJwt(userID.Hex())
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}
	expectedResponse := dto.UserRegisterResponse{
		ID:          userID.Hex(),
		Name:        req.Name,
		Email:       req.Email,
		AccessToken: expectedAccessToken,
	}

	mockUserRepository.On("SaveUser", ctx, mock.MatchedBy(func(u entity.User) bool {
		return u.Name == req.Name && u.Email == req.Email && bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) == nil
	})).Return(userEntity, nil)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.RegisterUser(ctx, req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, resp.ID)
	assert.Equal(t, expectedResponse.Name, resp.Name)
	assert.Equal(t, expectedResponse.Email, resp.Email)
	assert.NotEmpty(t, resp.AccessToken) // because it different each time
	mockUserRepository.AssertExpectations(t)
}

func TestRegisterHashPasswordError(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository)

	longPassword := strings.Repeat("a", 73)

	req := dto.UserRegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: longPassword,
	}

	expectedError := bcrypt.ErrPasswordTooLong

	// When
	resp, err := userService.RegisterUser(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserRegisterResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestRegisterUserDuplicateEmail(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	req := dto.UserRegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedError := fmt.Errorf("email %s is already registered", req.Email)
	duplicateKeyError := mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11000}}}

	mockUserRepository.On("SaveUser", ctx, mock.MatchedBy(func(u entity.User) bool {
		return u.Name == req.Name && u.Email == req.Email && bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) == nil
	})).Return(entity.User{}, duplicateKeyError)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.RegisterUser(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserRegisterResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestRegisterUserSaveUserError(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	req := dto.UserRegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedError := errors.New("repository save error")

	mockUserRepository.On("SaveUser", ctx, mock.MatchedBy(func(u entity.User) bool {
		return u.Name == req.Name && u.Email == req.Email && bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) == nil
	})).Return(entity.User{}, expectedError)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.RegisterUser(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserRegisterResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}
