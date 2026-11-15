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
	// ErrReturnConflict dikembalikan saat unique constraint item_returns (SQLSTATE 23505) terlanggar.
	ErrReturnConflict = errors.New("item return already exists")
	// ErrReportStatusChanged dikembalikan saat status laporan berubah di tengah proses serah terima.
	ErrReportStatusChanged = errors.New("report status changed during return")
)

type ReturnHistoryRepository interface {
	FindReportForReturn(ctx context.Context, reportID int) (*model.ReturnReportInfo, error)
	FindClaimForReturn(ctx context.Context, claimID int) (*model.ReturnClaimInfo, error)
	CreateReturn(ctx context.Context, p model.CreateReturnParams) (*model.ItemReturn, error)
	FindReturnByReportID(ctx context.Context, reportID int) (*model.ItemReturn, error)
	ReportExists(ctx context.Context, reportID int) (bool, error)
	ListHistoriesByReportID(ctx context.Context, reportID int) ([]model.ReportHistory, error)
}

type returnHistoryRepository struct {
	db *pgxpool.Pool
}

func NewReturnHistoryRepository(db *pgxpool.Pool) ReturnHistoryRepository {
	return &returnHistoryRepository{db: db}
}

func (r *returnHistoryRepository) FindReportForReturn(ctx context.Context, reportID int) (*model.ReturnReportInfo, error) {
	var info model.ReturnReportInfo
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, report_type, status FROM item_reports WHERE id = $1`, reportID,
	).Scan(&info.ID, &info.OwnerID, &info.ReportType, &info.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &info, nil
}

func (r *returnHistoryRepository) FindClaimForReturn(ctx context.Context, claimID int) (*model.ReturnClaimInfo, error) {
	var info model.ReturnClaimInfo
	err := r.db.QueryRow(ctx,
		`SELECT id, report_id, claimant_id, status FROM claims WHERE id = $1`, claimID,
	).Scan(&info.ID, &info.ReportID, &info.ClaimantID, &info.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &info, nil
}

// CreateReturn menyimpan serah terima secara atomik:
// (1) insert item_returns, (2) ubah status laporan menjadi 'returned', (3) catat ke report_histories.
func (r *returnHistoryRepository) CreateReturn(ctx context.Context, p model.CreateReturnParams) (*model.ItemReturn, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // no-op setelah Commit berhasil

	var id int
	err = tx.QueryRow(ctx, `
		INSERT INTO item_returns (report_id, claim_id, handed_over_by, received_by, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, p.ReportID, p.ClaimID, p.HandedOverBy, p.ReceivedBy, p.Notes).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrReturnConflict
		}
		return nil, err
	}

	// Optimistic check: hanya berubah jika status masih sama seperti saat divalidasi.
	tag, err := tx.Exec(ctx, `
		UPDATE item_reports SET status = 'returned', updated_at = NOW()
		WHERE id = $1 AND status = $2
	`, p.ReportID, p.OldStatus)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() != 1 {
		return nil, ErrReportStatusChanged
	}

	if _, err = tx.Exec(ctx, `
		INSERT INTO report_histories (report_id, actor_id, old_status, new_status, note)
		VALUES ($1, $2, $3, 'returned', $4)
	`, p.ReportID, p.HandedOverBy, p.OldStatus, p.Notes); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.FindReturnByReportID(ctx, p.ReportID)
}

func (r *returnHistoryRepository) FindReturnByReportID(ctx context.Context, reportID int) (*model.ItemReturn, error) {
	query := `
		SELECT ir.id, ir.report_id, ir.claim_id,
		       ir.handed_over_by, hu.name,
		       ir.received_by, ru.name,
		       ir.returned_at, ir.notes
		FROM item_returns ir
		JOIN users hu ON hu.id = ir.handed_over_by
		JOIN users ru ON ru.id = ir.received_by
		WHERE ir.report_id = $1
	`
	var it model.ItemReturn
	err := r.db.QueryRow(ctx, query, reportID).Scan(
		&it.ID, &it.ReportID, &it.ClaimID,
		&it.HandedOverBy, &it.HandedOverByName,
		&it.ReceivedBy, &it.ReceivedByName,
		&it.ReturnedAt, &it.Notes,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &it, nil
}

func (r *returnHistoryRepository) ReportExists(ctx context.Context, reportID int) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM item_reports WHERE id = $1)`, reportID).Scan(&exists)
	return exists, err
}

func (r *returnHistoryRepository) ListHistoriesByReportID(ctx context.Context, reportID int) ([]model.ReportHistory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT h.id, h.report_id, h.actor_id, u.name, h.old_status, h.new_status, h.note, h.created_at
		FROM report_histories h
		JOIN users u ON u.id = h.actor_id
		WHERE h.report_id = $1
		ORDER BY h.created_at ASC, h.id ASC
	`, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	histories := make([]model.ReportHistory, 0)
	for rows.Next() {
		var h model.ReportHistory
		if err := rows.Scan(&h.ID, &h.ReportID, &h.ActorID, &h.ActorName, &h.OldStatus, &h.NewStatus, &h.Note, &h.CreatedAt); err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}
	return histories, rows.Err()
}
