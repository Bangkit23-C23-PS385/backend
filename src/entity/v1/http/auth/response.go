package auth

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type LoginResponse struct {
	UserResponse
	AccessToken string `json:"access_token"`
}
