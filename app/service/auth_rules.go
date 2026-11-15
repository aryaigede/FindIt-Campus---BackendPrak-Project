package service

import (
	"strings"

	"findit-campus-backend/app/model"
)

func ValidateRegisterInput(req model.RegisterRequest) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.Name) == "" {
		errs["name"] = "name is required"
	}
	email := strings.TrimSpace(req.Email)
	if email == "" || !strings.Contains(email, "@") {
		errs["email"] = "valid email is required"
	}
	if len(req.Password) < 6 {
		errs["password"] = "password must be at least 6 characters"
	}

	return errs
}

func ValidateLoginInput(req model.LoginRequest) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.Email) == "" {
		errs["email"] = "email is required"
	}
	if req.Password == "" {
		errs["password"] = "password is required"
	}

	return errs
}