package test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/taninchot-work/backend-challenge/internal/dto"
	"github.com/taninchot-work/backend-challenge/internal/entity"
	mock_user_repository "github.com/taninchot-work/backend-challenge/internal/repository/mocks/user_repository_mock"
	"github.com/taninchot-work/backend-challenge/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"testing"
)

func TestGetUserByIdSuccess(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	expectedResponse := dto.UserGetMeResponse{
		ID:    "683ecde861d005de5ec0907d",
		Name:  "Test User",
		Email: "test@gmail.com",
	}
	objectId, err := bson.ObjectIDFromHex("683ecde861d005de5ec0907d")
	if err != nil {
		t.Fatalf("Failed to convert string to ObjectID: %v", err)
	}
	userEntity := entity.User{
		ID:    objectId,
		Name:  "Test User",
		Email: "test@gmail.com",
	}

	mockUserRepository.On("GetUserById", ctx, "683ecde861d005de5ec0907d").Return(userEntity, nil)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.GetUserByID(ctx, "683ecde861d005de5ec0907d")

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedResponse, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestGetUserByIdFailIdNotFound(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "nonexistentuserid"
	expectedError := fmt.Errorf("user with id %s not found", userID)

	mockUserRepository.On("GetUserById", ctx, userID).Return(entity.User{}, mongo.ErrNoDocuments)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.GetUserByID(ctx, userID)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserGetMeResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestGetUserByIdFailOnRepository(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	expectedError := fmt.Errorf("some thing went wrong")

	mockUserRepository.On("GetUserById", ctx, "683ecde861d005de5ec0907d").Return(entity.User{}, expectedError)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.GetUserByID(ctx, "683ecde861d005de5ec0907d")

	// Then
	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, dto.UserGetMeResponse{}, resp)
	assert.Equal(t, expectedError, err)
	mockUserRepository.AssertExpectations(t)
}

func TestGetUserByIdFailUserNotFound(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	expectedError := fmt.Errorf("user with id %s not found", "683ecde861d005de5ec0907d")

	mockUserRepository.On("GetUserById", ctx, "683ecde861d005de5ec0907d").Return(entity.User{}, nil)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.GetUserByID(ctx, "683ecde861d005de5ec0907d")

	// Then
	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, dto.UserGetMeResponse{}, resp)
	assert.Equal(t, expectedError, err)
	mockUserRepository.AssertExpectations(t)
}
