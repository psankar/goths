package libgoths

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

const (
	LoginFailed   = "Invalid email or password"
	InternalError = "Internal server error"
)
