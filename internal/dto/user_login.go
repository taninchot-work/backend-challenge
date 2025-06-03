package dto

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" minlength:"5" maxlength:"100"`
	Password string `json:"password" validate:"required" minlength:"8" maxlength:"50"`
}

type UserLoginResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}
