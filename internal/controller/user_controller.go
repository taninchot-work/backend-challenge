package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/taninchot-work/backend-challenge/internal/constant"
	"github.com/taninchot-work/backend-challenge/internal/core/util/json"
	"github.com/taninchot-work/backend-challenge/internal/dto"
	"github.com/taninchot-work/backend-challenge/internal/service"
	"net/http"
	"strings"
)

type UserController interface {
	GetMe(w http.ResponseWriter, r *http.Request)
	UserListGet(w http.ResponseWriter, r *http.Request)
	UserRegister(w http.ResponseWriter, r *http.Request)
	UserLogin(w http.ResponseWriter, r *http.Request)
	UserUpdate(w http.ResponseWriter, r *http.Request)
	UserDelete(w http.ResponseWriter, r *http.Request)
}

type userControllerImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userControllerImpl{
		userService: userService,
	}
}

func (c userControllerImpl) GetMe(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constant.CONTEXT_KEY_USER_ID).(string)
	if !ok || userId == "" {
		json.ResponseWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	response, err := c.userService.GetUserByID(r.Context(), userId)
	if err != nil {
		json.ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.ResponseWithSuccess(w, response)
	return
}

func (c userControllerImpl) UserListGet(w http.ResponseWriter, r *http.Request) {
	var req dto.UserListGetRequest

	if err := json.NewDecoder(r).Decode(&req); err != nil {
		json.ResponseWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// validate request
	if req.Page < 1 || req.Limit < 1 {
		req.Page = 1
		req.Limit = 20
	}
	response, err := c.userService.GetUserList(r.Context(), req)
	if err != nil {
		json.ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.ResponseWithSuccess(w, response)
	return
}

func (c userControllerImpl) UserRegister(w http.ResponseWriter, r *http.Request) {
	var req dto.UserRegisterRequest

	if err := json.NewDecoder(r).Decode(&req); err != nil {
		json.ResponseWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate request
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		json.ResponseWithError(w, "Validation failed: "+json.JoinErrors(validationErrors), http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(req.Email)

	response, err := c.userService.RegisterUser(r.Context(), req)
	if err != nil {
		json.ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.ResponseWithSuccess(w, response)
	return
}

func (c userControllerImpl) UserLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLoginRequest

	if err := json.NewDecoder(r).Decode(&req); err != nil {
		json.ResponseWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// validate request
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		json.ResponseWithError(w, "Validation failed: "+json.JoinErrors(validationErrors), http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(req.Email)

	response, err := c.userService.LoginUser(r.Context(), req)
	if err != nil {
		json.ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.ResponseWithSuccess(w, response)
	return
}

func (c userControllerImpl) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var req dto.UserUpdateRequest

	if err := json.NewDecoder(r).Decode(&req); err != nil {
		json.ResponseWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userId, ok := r.Context().Value(constant.CONTEXT_KEY_USER_ID).(string)
	if !ok || userId == "" {
		json.ResponseWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// validate request TODO: use a validation library
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Error())
		}
		json.ResponseWithError(w, "Validation failed: "+json.JoinErrors(validationErrors), http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(req.Email)

	response, err := c.userService.UpdateUser(r.Context(), userId, req)
	if err != nil {
		json.ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.ResponseWithSuccess(w, response)
	return
}

func (c userControllerImpl) UserDelete(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(constant.CONTEXT_KEY_USER_ID).(string)
	if !ok || userId == "" {
		json.ResponseWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := c.userService.DeleteUser(r.Context(), userId)
	if err != nil {
		json.ResponseWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.ResponseWithSuccess(w, "User deleted successfully")
	return
}
