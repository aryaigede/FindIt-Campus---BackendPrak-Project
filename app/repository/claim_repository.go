package repository

import (
	"context"
	"errors"

	"findit-campus-backend/app/model"
	"findit-campus-backend/helper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrClaimConflict = errors.New("conflict: claim already exists or violates constraint")
var ErrAlreadyReviewed = errors.New("claim has already been reviewed")

type ClaimRepository interface {
	Create(ctx context.Context, claim *model.Claim) error
	FindByID(ctx context.Context, id int) (*model.Claim, error)
	FindByClaimantID(ctx context.Context, claimantID int, q helper.ListQuery) ([]model.Claim, int64, error)
	FindAll(ctx context.Context, status string, q helper.ListQuery) ([]model.Claim, int64, error)
	FindReportForClaim(ctx context.Context, reportID int) (*model.ClaimReportLookup, error)
	ReviewClaim(ctx context.Context, claimID int, reviewerID int, newStatus string, reviewNote string) (*model.Claim, error)
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrClaimConflict
		}
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

func (r *claimRepository) ReviewClaim(ctx context.Context, claimID int, reviewerID int, newStatus string, reviewNote string) (*model.Claim, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var (
		currentStatus string
		reportID      int
	)
	claimCheckQuery := `SELECT status, report_id FROM claims WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, claimCheckQuery, claimID).Scan(&currentStatus, &reportID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if currentStatus != model.ClaimStatusPending {
		return nil, ErrAlreadyReviewed
	}

	var updatedClaim model.Claim
	if newStatus == model.ClaimStatusApproved {
		var reportStatus string
		err = tx.QueryRow(ctx, `SELECT status FROM item_reports WHERE id = $1 FOR UPDATE`, reportID).Scan(&reportStatus)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		if reportStatus != "open" {
			return nil, ErrClaimConflict
		}

		updateClaimQuery := `
			UPDATE claims
			SET status = 'approved', reviewed_by = $2, review_note = $3, reviewed_at = NOW(), updated_at = NOW()
			WHERE id = $1
			RETURNING id, report_id, claimant_id, reviewed_by, proof_description, proof_image_url,
			          status, review_note, reviewed_at, created_at, updated_at
		`
		err = tx.QueryRow(ctx, updateClaimQuery, claimID, reviewerID, reviewNote).Scan(
			&updatedClaim.ID, &updatedClaim.ReportID, &updatedClaim.ClaimantID, &updatedClaim.ReviewedBy,
			&updatedClaim.ProofDescription, &updatedClaim.ProofImageURL, &updatedClaim.Status,
			&updatedClaim.ReviewNote, &updatedClaim.ReviewedAt, &updatedClaim.CreatedAt, &updatedClaim.UpdatedAt,
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return nil, ErrClaimConflict
			}
			return nil, err
		}

		autoRejectQuery := `
			UPDATE claims
			SET status = 'rejected', reviewed_by = $2, review_note = 'Auto-rejected: another claim was approved',
			    reviewed_at = NOW(), updated_at = NOW()
			WHERE report_id = $1 AND id != $3 AND status = 'pending'
		`
		if _, err := tx.Exec(ctx, autoRejectQuery, reportID, reviewerID, claimID); err != nil {
			return nil, err
		}

		updateReportQuery := `
			UPDATE item_reports
			SET status = 'claimed', updated_at = NOW()
			WHERE id = $1
		`
		if _, err := tx.Exec(ctx, updateReportQuery, reportID); err != nil {
			return nil, err
		}

		insertHistoryQuery := `
			INSERT INTO report_histories (report_id, actor_id, old_status, new_status, note, created_at)
			VALUES ($1, $2, $3, 'claimed', $4, NOW())
		`
		historyNote := "Claim approved"
		if reviewNote != "" {
			historyNote = "Claim approved: " + reviewNote
		}
		if _, err := tx.Exec(ctx, insertHistoryQuery, reportID, reviewerID, reportStatus, historyNote); err != nil {
			return nil, err
		}
	} else {
		updateClaimQuery := `
			UPDATE claims
			SET status = 'rejected', reviewed_by = $2, review_note = $3, reviewed_at = NOW(), updated_at = NOW()
			WHERE id = $1
			RETURNING id, report_id, claimant_id, reviewed_by, proof_description, proof_image_url,
			          status, review_note, reviewed_at, created_at, updated_at
		`
		err = tx.QueryRow(ctx, updateClaimQuery, claimID, reviewerID, reviewNote).Scan(
			&updatedClaim.ID, &updatedClaim.ReportID, &updatedClaim.ClaimantID, &updatedClaim.ReviewedBy,
			&updatedClaim.ProofDescription, &updatedClaim.ProofImageURL, &updatedClaim.Status,
			&updatedClaim.ReviewNote, &updatedClaim.ReviewedAt, &updatedClaim.CreatedAt, &updatedClaim.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrClaimConflict
		}
		return nil, err
	}

	_ = r.db.QueryRow(ctx, `
		SELECT u.name, rev.name, rep.title
		FROM claims c
		JOIN users u ON c.claimant_id = u.id
		LEFT JOIN users rev ON c.reviewed_by = rev.id
		JOIN item_reports rep ON c.report_id = rep.id
		WHERE c.id = $1
	`, updatedClaim.ID).Scan(&updatedClaim.ClaimantName, &updatedClaim.ReviewerName, &updatedClaim.ReportTitle)

	return &updatedClaim, nil
}
