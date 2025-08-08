package dto

type RegisterDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	AccessToken string `json:"access_token"`
}
