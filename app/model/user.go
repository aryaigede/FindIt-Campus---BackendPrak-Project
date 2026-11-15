package model

import "time"

type User struct {
	ID           int       `json:"id"`
	RoleID       int       `json:"role_id"`
	RoleName     string    `json:"role"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	StudentID    *string   `json:"student_id,omitempty"`
	Phone        *string   `json:"phone,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	StudentID *string `json:"student_id"`
	Phone     *string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}