package repository

import (
	"context"
	"errors"

	"findit-campus-backend/app/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("record not found")

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id int) (*model.User, error)
	GetRoleIDByName(ctx context.Context, roleName string) (int, error)
	Create(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT u.id, u.role_id, ro.name, u.name, u.email, u.password_hash,
		       u.student_id, u.phone, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles ro ON u.role_id = ro.id
		WHERE u.email = $1
	`
	var u model.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Email, &u.PasswordHash,
		&u.StudentID, &u.Phone, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	query := `
		SELECT u.id, u.role_id, ro.name, u.name, u.email, u.password_hash,
		       u.student_id, u.phone, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles ro ON u.role_id = ro.id
		WHERE u.id = $1
	`
	var u model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Email, &u.PasswordHash,
		&u.StudentID, &u.Phone, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetRoleIDByName(ctx context.Context, roleName string) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, `SELECT id FROM roles WHERE name = $1`, roleName).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return id, nil
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	query := `
		INSERT INTO users (role_id, name, email, password_hash, student_id, phone, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		ctx, query,
		u.RoleID, u.Name, u.Email, u.PasswordHash, u.StudentID, u.Phone, u.IsActive,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}