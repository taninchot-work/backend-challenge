package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taninchot-work/backend-challenge/internal/core/util/jwt"
	"github.com/taninchot-work/backend-challenge/internal/dto"
	"github.com/taninchot-work/backend-challenge/internal/entity"
	mock_user_repository "github.com/taninchot-work/backend-challenge/internal/repository/mocks/user_repository_mock"
	"github.com/taninchot-work/backend-challenge/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestLoginUserSuccess(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository)

	req := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userID := bson.NewObjectID()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	userEntity := entity.User{
		ID:       userID,
		Name:     "Test User",
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	expectedAccessToken, _ := jwt.GenerateJwt(userID.Hex())
	expectedResponse := dto.UserLoginResponse{
		ID:          userID.Hex(),
		Name:        userEntity.Name,
		Email:       userEntity.Email,
		AccessToken: expectedAccessToken,
	}

	mockUserRepository.On("GetUserByEmail", ctx, req.Email).Return(userEntity, nil)

	// When
	resp, err := userService.LoginUser(ctx, req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, resp.ID)
	assert.Equal(t, expectedResponse.Name, resp.Name)
	assert.Equal(t, expectedResponse.Email, resp.Email)
	assert.NotEmpty(t, resp.AccessToken)
	mockUserRepository.AssertExpectations(t)
}

func TestLoginUserFailGetUserByEmailMongoErrNoDocuments(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository)

	req := dto.UserLoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	expectedError := fmt.Errorf("invalid email or password")

	mockUserRepository.On("GetUserByEmail", ctx, req.Email).Return(entity.User{}, mongo.ErrNoDocuments)

	// When
	resp, err := userService.LoginUser(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserLoginResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestLoginUserFailGetUserByEmailGenericError(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository)

	req := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	expectedError := errors.New("some repository error")

	mockUserRepository.On("GetUserByEmail", ctx, req.Email).Return(entity.User{}, expectedError)

	// When
	resp, err := userService.LoginUser(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserLoginResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestLoginUserFailUserNotFoundByEmptyStruct(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository)

	req := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockUserRepository.On("GetUserByEmail", ctx, req.Email).Return(entity.User{}, nil)

	// When
	resp, err := userService.LoginUser(ctx, req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, dto.UserLoginResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestLoginUserFailPasswordMismatch(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userService := service.NewUserService(mockUserRepository)

	req := dto.UserLoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	userID := bson.NewObjectID()
	correctPassword := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	userEntity := entity.User{
		ID:       userID,
		Name:     "Test User",
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	expectedError := fmt.Errorf("invalid email or password")

	mockUserRepository.On("GetUserByEmail", ctx, req.Email).Return(userEntity, nil)

	// When
	resp, err := userService.LoginUser(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, dto.UserLoginResponse{}, resp)
	assert.Equal(t, expectedError, err)
	mockUserRepository.AssertExpectations(t)
}
