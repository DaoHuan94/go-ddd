package types

type RegisterArg struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type LoginArg struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshArg struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutArg struct {
	RefreshToken string `json:"refresh_token"`
}

