package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"findit-campus-backend/app/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrConflict             = errors.New("record conflict")
	ErrForeignKeyConstraint = errors.New("foreign key constraint violation")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrLocationNotFound     = errors.New("location not found")
	ErrStatusNotOpen        = errors.New("report status is not open")
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

	// Reports
	CreateReport(ctx context.Context, report *model.ItemReport) error
	FindReports(ctx context.Context, filter model.ReportFilter) ([]model.ItemReport, int64, error)
	FindReportByID(ctx context.Context, id int) (*model.ItemReport, error)
	UpdateReport(ctx context.Context, report *model.ItemReport) error
	CloseReport(ctx context.Context, reportID int, actorID int, note string) error

	// Images
	AddImage(ctx context.Context, image *model.ItemImage) error
	CountImagesByReportID(ctx context.Context, reportID int) (int, error)
	FindImagesByReportID(ctx context.Context, reportID int) ([]model.ItemImage, error)
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

// Reports

func (r *reportRepository) CreateReport(ctx context.Context, report *model.ItemReport) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryReport := `
		INSERT INTO item_reports (
			user_id, category_id, location_id, report_type, title, description,
			custody_note, event_date, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'open')
		RETURNING id, status, created_at, updated_at
	`
	err = tx.QueryRow(ctx, queryReport,
		report.UserID, report.CategoryID, report.LocationID, report.ReportType,
		report.Title, report.Description, report.CustodyNote, report.EventDate,
	).Scan(&report.ID, &report.Status, &report.CreatedAt, &report.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			if strings.Contains(pgErr.ConstraintName, "category") || strings.Contains(pgErr.Detail, "category_id") {
				return ErrCategoryNotFound
			}
			if strings.Contains(pgErr.ConstraintName, "location") || strings.Contains(pgErr.Detail, "location_id") {
				return ErrLocationNotFound
			}
			return ErrForeignKeyConstraint
		}
		return err
	}

	queryHistory := `
		INSERT INTO report_histories (report_id, actor_id, old_status, new_status, note)
		VALUES ($1, $2, NULL, 'open', $3)
	`
	note := "Laporan dibuat"
	if _, err := tx.Exec(ctx, queryHistory, report.ID, report.UserID, note); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *reportRepository) FindReports(ctx context.Context, filter model.ReportFilter) ([]model.ItemReport, int64, error) {
	where := []string{"1=1"}
	var args []interface{}
	argIdx := 1

	if filter.ReportType != nil {
		where = append(where, fmt.Sprintf("r.report_type = $%d", argIdx))
		args = append(args, *filter.ReportType)
		argIdx++
	}
	if filter.CategoryID != nil {
		where = append(where, fmt.Sprintf("r.category_id = $%d", argIdx))
		args = append(args, *filter.CategoryID)
		argIdx++
	}
	if filter.LocationID != nil {
		where = append(where, fmt.Sprintf("r.location_id = $%d", argIdx))
		args = append(args, *filter.LocationID)
		argIdx++
	}
	if filter.Status != nil {
		where = append(where, fmt.Sprintf("r.status = $%d", argIdx))
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		where = append(where, fmt.Sprintf("(r.title ILIKE $%d OR r.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, searchPattern)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM item_reports r
		WHERE %s
	`, whereClause)

	var totalItems int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.user_id, u.name, r.category_id, c.name, r.location_id, l.name,
		       r.report_type, r.title, r.description, r.custody_note, r.event_date,
		       r.status, r.created_at, r.updated_at
		FROM item_reports r
		JOIN users u ON r.user_id = u.id
		JOIN categories c ON r.category_id = c.id
		JOIN locations l ON r.location_id = l.id
		WHERE %s
		ORDER BY r.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	queryArgs := append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	reports := make([]model.ItemReport, 0)
	for rows.Next() {
		var item model.ItemReport
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.UserName, &item.CategoryID, &item.CategoryName,
			&item.LocationID, &item.LocationName, &item.ReportType, &item.Title,
			&item.Description, &item.CustodyNote, &item.EventDate, &item.Status,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		reports = append(reports, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return reports, totalItems, nil
}

func (r *reportRepository) FindReportByID(ctx context.Context, id int) (*model.ItemReport, error) {
	query := `
		SELECT r.id, r.user_id, u.name, r.category_id, c.name, r.location_id, l.name,
		       r.report_type, r.title, r.description, r.custody_note, r.event_date,
		       r.status, r.created_at, r.updated_at
		FROM item_reports r
		JOIN users u ON r.user_id = u.id
		JOIN categories c ON r.category_id = c.id
		JOIN locations l ON r.location_id = l.id
		WHERE r.id = $1
	`
	var item model.ItemReport
	err := r.db.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.UserID, &item.UserName, &item.CategoryID, &item.CategoryName,
		&item.LocationID, &item.LocationName, &item.ReportType, &item.Title,
		&item.Description, &item.CustodyNote, &item.EventDate, &item.Status,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	images, err := r.FindImagesByReportID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Images = images

	return &item, nil
}

func (r *reportRepository) UpdateReport(ctx context.Context, report *model.ItemReport) error {
	var currentStatus string
	checkErr := r.db.QueryRow(ctx, "SELECT status FROM item_reports WHERE id = $1", report.ID).Scan(&currentStatus)
	if checkErr != nil {
		if errors.Is(checkErr, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return checkErr
	}

	if currentStatus != model.ReportStatusOpen {
		return ErrStatusNotOpen
	}

	query := `
		UPDATE item_reports
		SET category_id = $1, location_id = $2, title = $3, description = $4,
		    custody_note = $5, event_date = $6, updated_at = NOW()
		WHERE id = $7 AND status = 'open'
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		report.CategoryID, report.LocationID, report.Title, report.Description,
		report.CustodyNote, report.EventDate, report.ID,
	).Scan(&report.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			if strings.Contains(pgErr.ConstraintName, "category") || strings.Contains(pgErr.Detail, "category_id") {
				return ErrCategoryNotFound
			}
			if strings.Contains(pgErr.ConstraintName, "location") || strings.Contains(pgErr.Detail, "location_id") {
				return ErrLocationNotFound
			}
			return ErrForeignKeyConstraint
		}
		return err
	}
	return nil
}

func (r *reportRepository) CloseReport(ctx context.Context, reportID int, actorID int, note string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var currentStatus string
	err = tx.QueryRow(ctx, "SELECT status FROM item_reports WHERE id = $1 FOR UPDATE", reportID).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if currentStatus != model.ReportStatusOpen {
		return ErrStatusNotOpen
	}

	updateQuery := `
		UPDATE item_reports
		SET status = 'closed', updated_at = NOW()
		WHERE id = $1
	`
	if _, err := tx.Exec(ctx, updateQuery, reportID); err != nil {
		return err
	}

	if strings.TrimSpace(note) == "" {
		note = "Laporan ditutup"
	}

	historyQuery := `
		INSERT INTO report_histories (report_id, actor_id, old_status, new_status, note)
		VALUES ($1, $2, 'open', 'closed', $3)
	`
	if _, err := tx.Exec(ctx, historyQuery, reportID, actorID, note); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Images

func (r *reportRepository) AddImage(ctx context.Context, img *model.ItemImage) error {
	query := `
		INSERT INTO item_images (report_id, image_url)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query, img.ReportID, img.ImageURL).Scan(&img.ID, &img.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *reportRepository) CountImagesByReportID(ctx context.Context, reportID int) (int, error) {
	query := `SELECT COUNT(*) FROM item_images WHERE report_id = $1`
	var count int
	if err := r.db.QueryRow(ctx, query, reportID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *reportRepository) FindImagesByReportID(ctx context.Context, reportID int) ([]model.ItemImage, error) {
	query := `
		SELECT id, report_id, image_url, created_at
		FROM item_images
		WHERE report_id = $1
		ORDER BY id ASC
	`
	rows, err := r.db.Query(ctx, query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make([]model.ItemImage, 0)
	for rows.Next() {
		var img model.ItemImage
		if err := rows.Scan(&img.ID, &img.ReportID, &img.ImageURL, &img.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return images, nil
}
