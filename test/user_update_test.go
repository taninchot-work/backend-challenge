package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taninchot-work/backend-challenge/internal/dto"
	"github.com/taninchot-work/backend-challenge/internal/entity"
	mock_user_repository "github.com/taninchot-work/backend-challenge/internal/repository/mocks/user_repository_mock"
	"github.com/taninchot-work/backend-challenge/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"testing"
)

func TestUpdateUserSuccess(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	objectID, _ := bson.ObjectIDFromHex(userID)
	req := dto.UserUpdateRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}
	userEntity := entity.User{
		ID:    objectID,
		Name:  "Original Name",
		Email: "original@example.com",
	}
	updatedUserEntity := entity.User{
		ID:    objectID,
		Name:  req.Name,
		Email: req.Email,
	}

	expectedResponse := dto.UserUpdateResponse{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
	}

	mockUserRepository.On("GetUserById", ctx, userID).Return(userEntity, nil)
	mockUserRepository.On("UpdateUser", ctx, updatedUserEntity).Return(updatedUserEntity, nil)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.UpdateUser(ctx, userID, req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestUpdateUserFailUserNotFound(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "nonexistentuserid"
	req := dto.UserUpdateRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}
	expectedError := fmt.Errorf("user with id %s not found", userID)

	mockUserRepository.On("GetUserById", ctx, userID).Return(entity.User{}, mongo.ErrNoDocuments)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.UpdateUser(ctx, userID, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserUpdateResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestUpdateUserFailDuplicateEmail(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	objectID, _ := bson.ObjectIDFromHex(userID)
	req := dto.UserUpdateRequest{
		Name:  "Updated Name",
		Email: "duplicate@example.com",
	}
	userEntity := entity.User{
		ID:    objectID,
		Name:  "Original Name",
		Email: "original@example.com",
	}
	expectedError := fmt.Errorf("email %s is already exists", req.Email)
	duplicateKeyError := mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11000}}}

	mockUserRepository.On("GetUserById", ctx, userID).Return(userEntity, nil)
	mockUserRepository.On("UpdateUser", ctx, mock.Anything).Return(entity.User{}, duplicateKeyError)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.UpdateUser(ctx, userID, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserUpdateResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestUpdateUserFailUpdateUserError(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	objectID, _ := bson.ObjectIDFromHex(userID)
	req := dto.UserUpdateRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}
	userEntity := entity.User{
		ID:    objectID,
		Name:  "Original Name",
		Email: "original@example.com",
	}
	expectedError := errors.New("repository update error")

	mockUserRepository.On("GetUserById", ctx, userID).Return(userEntity, nil)
	mockUserRepository.On("UpdateUser", ctx, mock.Anything).Return(entity.User{}, expectedError)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.UpdateUser(ctx, userID, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserUpdateResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestUpdateUserFailGetUserByIdError(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	objectID, _ := bson.ObjectIDFromHex(userID)
	req := dto.UserUpdateRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}
	userEntity := entity.User{
		ID:    objectID,
		Name:  "Original Name",
		Email: "original@example.com",
	}
	expectedError := errors.New("repository get error")

	mockUserRepository.On("GetUserById", ctx, userID).Return(userEntity, expectedError)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.UpdateUser(ctx, userID, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserUpdateResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestUpdateUserFailGetUserByIdAndUserEmpty(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	userID := "683ecde861d005de5ec0907d"
	req := dto.UserUpdateRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
	}

	expectedError := fmt.Errorf("user with id %s not found", userID)

	mockUserRepository.On("GetUserById", ctx, userID).Return(entity.User{}, nil)

	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.UpdateUser(ctx, userID, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserUpdateResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}
