package model

import "time"

// ItemReturn merepresentasikan berita acara serah terima barang (tabel item_returns).
type ItemReturn struct {
	ID               int       `json:"id"`
	ReportID         int       `json:"report_id"`
	ClaimID          *int      `json:"claim_id"`
	HandedOverBy     int       `json:"handed_over_by"`
	HandedOverByName string    `json:"handed_over_by_name"`
	ReceivedBy       int       `json:"received_by"`
	ReceivedByName   string    `json:"received_by_name"`
	ReturnedAt       time.Time `json:"returned_at"`
	Notes            *string   `json:"notes,omitempty"`
}

// CreateReturnRequest adalah payload POST /api/v1/returns.
// claim_id wajib untuk laporan 'found' dan harus kosong untuk laporan 'lost'.
type CreateReturnRequest struct {
	ReportID int     `json:"report_id"`
	ClaimID  *int    `json:"claim_id"`
	Notes    *string `json:"notes"`
}

// ReportHistory merepresentasikan satu baris linimasa audit status laporan (tabel report_histories).
type ReportHistory struct {
	ID        int       `json:"id"`
	ReportID  int       `json:"report_id"`
	ActorID   int       `json:"actor_id"`
	ActorName string    `json:"actor_name"`
	OldStatus *string   `json:"old_status"`
	NewStatus string    `json:"new_status"`
	Note      *string   `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ReturnReportInfo adalah potret ringan item_reports yang dibutuhkan aturan serah terima.
type ReturnReportInfo struct {
	ID         int
	OwnerID    int
	ReportType string
	Status     string
}

// ReturnClaimInfo adalah potret ringan claims yang dibutuhkan aturan serah terima.
type ReturnClaimInfo struct {
	ID         int
	ReportID   int
	ClaimantID int
	Status     string
}

// CreateReturnParams adalah data tervalidasi yang disimpan dalam satu transaksi serah terima.
type CreateReturnParams struct {
	ReportID     int
	ClaimID      *int
	HandedOverBy int
	ReceivedBy   int
	OldStatus    string
	Notes        *string
}
