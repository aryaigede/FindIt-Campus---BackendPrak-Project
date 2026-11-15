package service

import (
	"errors"
	"strings"
	"time"

	"findit-campus-backend/app/model"
	"findit-campus-backend/app/repository"
	"findit-campus-backend/config"
	"findit-campus-backend/helper"
	"findit-campus-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateRegisterInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	existing, err := s.userRepo.FindByEmail(c.Context(), email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to verify email")
	}
	if existing != nil {
		return helper.Fail(c, fiber.StatusConflict, "email is already registered")
	}

	roleID, err := s.userRepo.GetRoleIDByName(c.Context(), "user")
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "default user role not configured")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to secure password")
	}

	user := &model.User{
		RoleID:       roleID,
		RoleName:     "user",
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: string(hashed),
		StudentID:    req.StudentID,
		Phone:        req.Phone,
		IsActive:     true,
	}

	if err := s.userRepo.Create(c.Context(), user); err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to register user")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to generate token")
	}

	return helper.Created(c, "registration successful", model.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateLoginInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.userRepo.FindByEmail(c.Context(), email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusUnauthorized, "invalid email or password")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to process login")
	}

	if !user.IsActive {
		return helper.Fail(c, fiber.StatusForbidden, "user account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return helper.Fail(c, fiber.StatusUnauthorized, "invalid email or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to generate token")
	}

	return helper.Success(c, "login successful", model.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	user, err := s.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "user not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve profile")
	}

	return helper.Success(c, "profile retrieved", user)
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.RoleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.JWTExpiresHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}