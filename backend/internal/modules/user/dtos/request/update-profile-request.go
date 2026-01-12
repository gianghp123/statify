package request

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
