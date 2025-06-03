package dto

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required" minlength:"3" maxlength:"50"`
	Email    string `json:"email" validate:"required,email" minlength:"5" maxlength:"100"`
	Password string `json:"password" validate:"required" minlength:"8" maxlength:"50"`
}

type UserRegisterResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}
