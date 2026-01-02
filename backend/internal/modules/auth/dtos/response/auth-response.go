package response

import "github.com/gianghp/statify/internal/modules/user/dtos/response"

type AuthResponse struct {
	User  response.UserDto `json:"user"`
	Token string           `json:"token"`
}
