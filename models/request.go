package models

import "github.com/google/uuid"

type UserLoginRequest struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Name  string    `json:"name" db:"name"`
	Email string    `json:"email" db:"email"`
}
type UserRegisterRequest struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email" db:"email"`
	Role     string    `json:"role" db:"role"`
	UserType string    `json:"user_type" db:"user_type"`
}
