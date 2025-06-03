package dto

type UserUpdateRequest struct {
	Name  string `json:"name" validate:"required" minlength:"3" maxlength:"50"`
	Email string `json:"email" validate:"required,email" minlength:"5" maxlength:"100"`
}

type UserUpdateResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
