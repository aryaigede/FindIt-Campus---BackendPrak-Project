package repository

import (
	"context"
	"errors"

	"findit-campus-backend/app/model"
	"findit-campus-backend/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClaimRepository interface {
	Create(ctx context.Context, claim *model.Claim) error
	FindByID(ctx context.Context, id int) (*model.Claim, error)
	FindByClaimantID(ctx context.Context, claimantID int, q helper.ListQuery) ([]model.Claim, int64, error)
	FindAll(ctx context.Context, status string, q helper.ListQuery) ([]model.Claim, int64, error)
	FindReportForClaim(ctx context.Context, reportID int) (*model.ClaimReportLookup, error)
}

type claimRepository struct {
	db *pgxpool.Pool
}

func NewClaimRepository(db *pgxpool.Pool) ClaimRepository {
	return &claimRepository{db: db}
}

func (r *claimRepository) Create(ctx context.Context, c *model.Claim) error {
	query := `
		INSERT INTO claims (report_id, claimant_id, proof_description, proof_image_url, status)
		VALUES ($1, $2, $3, $4, 'pending')
		RETURNING id, status, created_at, updated_at
	`
	err := r.db.QueryRow(
		ctx, query,
		c.ReportID, c.ClaimantID, c.ProofDescription, c.ProofImageURL,
	).Scan(&c.ID, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *claimRepository) FindByID(ctx context.Context, id int) (*model.Claim, error) {
	query := `
		SELECT c.id, c.report_id, c.claimant_id, c.reviewed_by, c.proof_description,
		       c.proof_image_url, c.status, c.review_note, c.reviewed_at, c.created_at, c.updated_at,
		       u.name AS claimant_name, rev.name AS reviewer_name, rep.title AS report_title
		FROM claims c
		JOIN users u ON c.claimant_id = u.id
		LEFT JOIN users rev ON c.reviewed_by = rev.id
		JOIN item_reports rep ON c.report_id = rep.id
		WHERE c.id = $1
	`
	var c model.Claim
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.ReportID, &c.ClaimantID, &c.ReviewedBy, &c.ProofDescription,
		&c.ProofImageURL, &c.Status, &c.ReviewNote, &c.ReviewedAt, &c.CreatedAt, &c.UpdatedAt,
		&c.ClaimantName, &c.ReviewerName, &c.ReportTitle,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *claimRepository) FindByClaimantID(ctx context.Context, claimantID int, q helper.ListQuery) ([]model.Claim, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM claims WHERE claimant_id = $1`
	if err := r.db.QueryRow(ctx, countQuery, claimantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT c.id, c.report_id, c.claimant_id, c.reviewed_by, c.proof_description,
		       c.proof_image_url, c.status, c.review_note, c.reviewed_at, c.created_at, c.updated_at,
		       u.name AS claimant_name, rev.name AS reviewer_name, rep.title AS report_title
		FROM claims c
		JOIN users u ON c.claimant_id = u.id
		LEFT JOIN users rev ON c.reviewed_by = rev.id
		JOIN item_reports rep ON c.report_id = rep.id
		WHERE c.claimant_id = $1
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, claimantID, q.Limit, q.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	claims := make([]model.Claim, 0)
	for rows.Next() {
		var c model.Claim
		if err := rows.Scan(
			&c.ID, &c.ReportID, &c.ClaimantID, &c.ReviewedBy, &c.ProofDescription,
			&c.ProofImageURL, &c.Status, &c.ReviewNote, &c.ReviewedAt, &c.CreatedAt, &c.UpdatedAt,
			&c.ClaimantName, &c.ReviewerName, &c.ReportTitle,
		); err != nil {
			return nil, 0, err
		}
		claims = append(claims, c)
	}

	return claims, total, nil
}

func (r *claimRepository) FindAll(ctx context.Context, status string, q helper.ListQuery) ([]model.Claim, int64, error) {
	var total int64
	countQuery := `
		SELECT COUNT(*)
		FROM claims
		WHERE ($1 = '' OR status = $1)
	`
	if err := r.db.QueryRow(ctx, countQuery, status).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT c.id, c.report_id, c.claimant_id, c.reviewed_by, c.proof_description,
		       c.proof_image_url, c.status, c.review_note, c.reviewed_at, c.created_at, c.updated_at,
		       u.name AS claimant_name, rev.name AS reviewer_name, rep.title AS report_title
		FROM claims c
		JOIN users u ON c.claimant_id = u.id
		LEFT JOIN users rev ON c.reviewed_by = rev.id
		JOIN item_reports rep ON c.report_id = rep.id
		WHERE ($1 = '' OR c.status = $1)
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, status, q.Limit, q.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	claims := make([]model.Claim, 0)
	for rows.Next() {
		var c model.Claim
		if err := rows.Scan(
			&c.ID, &c.ReportID, &c.ClaimantID, &c.ReviewedBy, &c.ProofDescription,
			&c.ProofImageURL, &c.Status, &c.ReviewNote, &c.ReviewedAt, &c.CreatedAt, &c.UpdatedAt,
			&c.ClaimantName, &c.ReviewerName, &c.ReportTitle,
		); err != nil {
			return nil, 0, err
		}
		claims = append(claims, c)
	}

	return claims, total, nil
}

func (r *claimRepository) FindReportForClaim(ctx context.Context, reportID int) (*model.ClaimReportLookup, error) {
	query := `
		SELECT id, user_id, report_type, status
		FROM item_reports
		WHERE id = $1
	`
	var rep model.ClaimReportLookup
	err := r.db.QueryRow(ctx, query, reportID).Scan(
		&rep.ID, &rep.UserID, &rep.ReportType, &rep.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rep, nil
}
