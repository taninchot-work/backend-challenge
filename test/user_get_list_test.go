package test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/taninchot-work/backend-challenge/internal/dto"
	"github.com/taninchot-work/backend-challenge/internal/entity"
	mock_user_repository "github.com/taninchot-work/backend-challenge/internal/repository/mocks/user_repository_mock"
	"github.com/taninchot-work/backend-challenge/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

func TestGetUserListSuccess(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	req := dto.UserListGetRequest{
		Page:  1,
		Limit: 10,
	}
	offset := (req.Page - 1) * req.Limit
	objectID1, _ := bson.ObjectIDFromHex("683ecde861d005de5ec0907d")
	objectID2, _ := bson.ObjectIDFromHex("683ecde861d005de5ec0907e")
	usersEntity := []entity.User{
		{ID: objectID1, Name: "Test User 1", Email: "test1@gmail.com"},
		{ID: objectID2, Name: "Test User 2", Email: "test2@gmail.com"},
	}
	expectedResponse := dto.UserListGetResponse{
		Users: []dto.UserListGetResponseItem{
			{ID: "683ecde861d005de5ec0907d", Name: "Test User 1", Email: "test1@gmail.com"},
			{ID: "683ecde861d005de5ec0907e", Name: "Test User 2", Email: "test2@gmail.com"},
		},
		Page: req.Page,
	}

	mockUserRepository.On("GetUserList", ctx, offset, req.Limit).Return(usersEntity, nil)
	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.GetUserList(ctx, req)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
	mockUserRepository.AssertExpectations(t)
}

func TestGetUserListError(t *testing.T) {
	// Given
	ctx := context.Background()
	mockUserRepository := mock_user_repository.NewUserRepository(t)
	req := dto.UserListGetRequest{
		Page:  1,
		Limit: 10,
	}
	offset := (req.Page - 1) * req.Limit
	expectedError := errors.New("repository error")

	mockUserRepository.On("GetUserList", ctx, offset, req.Limit).Return([]entity.User{}, expectedError)
	userService := service.NewUserService(mockUserRepository)

	// When
	resp, err := userService.GetUserList(ctx, req)

	// Then
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dto.UserListGetResponse{}, resp)
	mockUserRepository.AssertExpectations(t)
}
