package response

import "time"

type UserDto struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username" copier:"must,nopanic"`
	Email     string    `json:"email" copier:"must,nopanic"`
	Role      string    `json:"role" copier:"must,nopanic"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
