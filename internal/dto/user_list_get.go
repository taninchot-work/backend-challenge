package dto

type UserListGetRequest struct {
	Page  int `json:"page" validate:"required"`
	Limit int `json:"limit" validate:"required"`
}

type UserListGetResponse struct {
	Users []UserListGetResponseItem `json:"users"`
	Page  int                       `json:"page"`
}
type UserListGetResponseItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
