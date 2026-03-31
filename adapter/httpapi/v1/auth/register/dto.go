package register

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

