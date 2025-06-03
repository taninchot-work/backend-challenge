package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taninchot-work/backend-challenge/internal/entity"
	mock_user_repository "github.com/taninchot-work/backend-challenge/internal/repository/mocks/user_repository_mock"
	"github.com/taninchot-work/backend-challenge/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"testing"
)

func TestDeleteUserSuccess(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	objectID, _ := bson.ObjectIDFromHex(userID)
	userEntity := entity.User{
		ID:    objectID,
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockUserRepository.On("GetUserById", ctx, userID).Return(userEntity, nil)
	mockUserRepository.On("DeleteUser", ctx, userID).Return(nil)
	userService := service.NewUserService(mockUserRepository)

	// When
	err := userService.DeleteUser(ctx, userID)

	// Then
	assert.NoError(t, err)
	mockUserRepository.AssertExpectations(t)
}

func TestDeleteUserFailUserNotFound(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	expectedError := errors.New("user not found")

	mockUserRepository.On("GetUserById", ctx, userID).Return(entity.User{}, expectedError)
	userService := service.NewUserService(mockUserRepository)

	// When
	err := userService.DeleteUser(ctx, userID)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockUserRepository.AssertExpectations(t)
}

func TestDeleteUserFailGetUserByIdError(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	expectedError := fmt.Errorf("user with id %s not found", userID)

	mockUserRepository.On("GetUserById", ctx, userID).Return(entity.User{}, mongo.ErrNoDocuments)

	userService := service.NewUserService(mockUserRepository)

	// When
	err := userService.DeleteUser(ctx, userID)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockUserRepository.AssertExpectations(t)
}

func TestDeleteUserFailRepositoryDeleteError(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	objectID, _ := bson.ObjectIDFromHex(userID)
	userEntity := entity.User{
		ID:    objectID,
		Name:  "Test User",
		Email: "test@example.com",
	}
	expectedError := errors.New("repository delete error")

	mockUserRepository.On("GetUserById", ctx, userID).Return(userEntity, nil)
	mockUserRepository.On("DeleteUser", ctx, userID).Return(expectedError)

	userService := service.NewUserService(mockUserRepository)

	// When
	err := userService.DeleteUser(ctx, userID)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockUserRepository.AssertExpectations(t)
}
