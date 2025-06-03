package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/taninchot-work/backend-challenge/internal/core/util/jwt"
	"github.com/taninchot-work/backend-challenge/internal/dto"
	"github.com/taninchot-work/backend-challenge/internal/entity"
	"github.com/taninchot-work/backend-challenge/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (dto.UserGetMeResponse, error)
	GetUserList(ctx context.Context, req dto.UserListGetRequest) (dto.UserListGetResponse, error)
	RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserRegisterResponse, error)
	LoginUser(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error)
	UpdateUser(ctx context.Context, id string, req dto.UserUpdateRequest) (dto.UserUpdateResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepository: userRepository,
	}
}

func (s userServiceImpl) GetUserByID(ctx context.Context, id string) (dto.UserGetMeResponse, error) {
	user, err := s.userRepository.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("user not found with id:", id)
			return dto.UserGetMeResponse{}, fmt.Errorf("user with id %s not found", id)
		}
		log.Println("user get by id failed:", err)
		return dto.UserGetMeResponse{}, err
	}
	if user == (entity.User{}) {
		log.Println("user not found")
		return dto.UserGetMeResponse{}, fmt.Errorf("user with id %s not found", id)
	}
	return dto.UserGetMeResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s userServiceImpl) GetUserList(ctx context.Context, req dto.UserListGetRequest) (dto.UserListGetResponse, error) {
	offset := (req.Page - 1) * req.Limit

	users, err := s.userRepository.GetUserList(ctx, offset, req.Limit)
	if err != nil {
		log.Println("user list get failed:", err)
		return dto.UserListGetResponse{}, err
	}

	var userResponses []dto.UserListGetResponseItem
	for _, user := range users {
		userResponses = append(userResponses, dto.UserListGetResponseItem{
			ID:    user.ID.Hex(),
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return dto.UserListGetResponse{
		Users: userResponses,
		Page:  req.Page,
	}, nil
}

func (s userServiceImpl) RegisterUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserRegisterResponse, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserRegisterResponse{}, err
	}
	user := entity.User{
		ID:        bson.NewObjectID(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	createdUser, err := s.userRepository.SaveUser(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Println("user register failed due to duplicate email:", err)
			return dto.UserRegisterResponse{}, fmt.Errorf("email %s is already registered", req.Email)
		}
		return dto.UserRegisterResponse{}, err
	}
	accessToken, err := jwt.GenerateJwt(createdUser.ID.Hex())
	return dto.UserRegisterResponse{
		ID:          createdUser.ID.Hex(),
		Name:        createdUser.Name,
		Email:       createdUser.Email,
		AccessToken: accessToken,
	}, nil
}

func (s userServiceImpl) LoginUser(ctx context.Context, req dto.UserLoginRequest) (dto.UserLoginResponse, error) {
	// check if user exists
	user, err := s.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("user login failed user not found with email:", req.Email)
			return dto.UserLoginResponse{}, fmt.Errorf("invalid email or password")
		}
		log.Println("user login failed to get user by email:", err)
		return dto.UserLoginResponse{}, err
	}

	if user == (entity.User{}) {
		log.Println("user login user not found")
		return dto.UserLoginResponse{}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println("user login password mismatch:", err)
		return dto.UserLoginResponse{}, fmt.Errorf("invalid email or password")
	}
	accessToken, err := jwt.GenerateJwt(user.ID.Hex())
	if err != nil {
		log.Println("user login failed to generate access token:", err)
		return dto.UserLoginResponse{}, err
	}

	return dto.UserLoginResponse{
		ID:          user.ID.Hex(),
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil
}

func (s userServiceImpl) UpdateUser(ctx context.Context, id string, req dto.UserUpdateRequest) (dto.UserUpdateResponse, error) {
	user, err := s.userRepository.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("user update failed user not found with id:", id)
			return dto.UserUpdateResponse{}, fmt.Errorf("user with id %s not found", id)
		}
		log.Println("user update failed:", err)
		return dto.UserUpdateResponse{}, err
	}

	if user == (entity.User{}) {
		log.Println("user update user not found")
		return dto.UserUpdateResponse{}, fmt.Errorf("user with id %s not found", id)
	}

	user.Name = req.Name
	user.Email = req.Email

	updatedUser, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		log.Println("user update failed:", err)
		if mongo.IsDuplicateKeyError(err) {
			log.Println("user update failed due to duplicate email:", err)
			return dto.UserUpdateResponse{}, fmt.Errorf("email %s is already exists", req.Email)
		}
		return dto.UserUpdateResponse{}, err
	}

	return dto.UserUpdateResponse{
		ID:    updatedUser.ID.Hex(),
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}, nil
}

func (s userServiceImpl) DeleteUser(ctx context.Context, id string) error {
	user, err := s.userRepository.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("user delete failed user not found with id:", id)
			return fmt.Errorf("user with id %s not found", id)
		}
		log.Println("user delete failed:", err)
		return err
	}

	err = s.userRepository.DeleteUser(ctx, user.ID.Hex())
	if err != nil {
		log.Println("user delete failed:", err)
		return err
	}
	return nil
}
