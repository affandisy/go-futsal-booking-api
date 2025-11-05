package response

import (
	"go-futsal-booking-api/internal/domain"
	"time"
)

type UserResponse struct {
	ID        uint      `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Age:       user.Age,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
	}
}
