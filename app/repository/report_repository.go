package repository

import (
	"context"
	"errors"

	"findit-campus-backend/app/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrConflict             = errors.New("record conflict")
	ErrForeignKeyConstraint = errors.New("foreign key constraint violation")
)

type ReportRepository interface {
	// Categories
	FindAllCategories(ctx context.Context) ([]model.Category, error)
	FindCategoryByID(ctx context.Context, id int) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id int) error

	// Locations
	FindAllLocations(ctx context.Context) ([]model.Location, error)
	FindLocationByID(ctx context.Context, id int) (*model.Location, error)
	CreateLocation(ctx context.Context, location *model.Location) error
	UpdateLocation(ctx context.Context, location *model.Location) error
	DeleteLocation(ctx context.Context, id int) error
}

type reportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) ReportRepository {
	return &reportRepository{db: db}
}

// Categories

func (r *reportRepository) FindAllCategories(ctx context.Context) ([]model.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY id ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]model.Category, 0)
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *reportRepository) FindCategoryByID(ctx context.Context, id int) (*model.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1
	`
	var c model.Category
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *reportRepository) CreateCategory(ctx context.Context, cat *model.Category) error {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, cat.Name, cat.Description).Scan(
		&cat.ID, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *reportRepository) UpdateCategory(ctx context.Context, cat *model.Category) error {
	query := `
		UPDATE categories
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query, cat.Name, cat.Description, cat.ID).Scan(&cat.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *reportRepository) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	ct, err := r.db.Exec(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrForeignKeyConstraint
		}
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Locations

func (r *reportRepository) FindAllLocations(ctx context.Context) ([]model.Location, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM locations
		ORDER BY id ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locations := make([]model.Location, 0)
	for rows.Next() {
		var l model.Location
		if err := rows.Scan(&l.ID, &l.Name, &l.Description, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *reportRepository) FindLocationByID(ctx context.Context, id int) (*model.Location, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM locations
		WHERE id = $1
	`
	var l model.Location
	err := r.db.QueryRow(ctx, query, id).Scan(
		&l.ID, &l.Name, &l.Description, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &l, nil
}

func (r *reportRepository) CreateLocation(ctx context.Context, loc *model.Location) error {
	query := `
		INSERT INTO locations (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, loc.Name, loc.Description).Scan(
		&loc.ID, &loc.CreatedAt, &loc.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *reportRepository) UpdateLocation(ctx context.Context, loc *model.Location) error {
	query := `
		UPDATE locations
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query, loc.Name, loc.Description, loc.ID).Scan(&loc.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrConflict
		}
		return err
	}
	return nil
}

func (r *reportRepository) DeleteLocation(ctx context.Context, id int) error {
	query := `DELETE FROM locations WHERE id = $1`
	ct, err := r.db.Exec(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrForeignKeyConstraint
		}
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
