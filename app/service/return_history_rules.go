package service

import (
	"net/http"
	"strings"

	"findit-campus-backend/app/model"
)

const maxReturnNotesLength = 500

// RuleError adalah pelanggaran aturan bisnis beserta status HTTP yang sesuai.
// Package ini sengaja tidak meng-import Fiber agar logika bisnis tetap murni.
type RuleError struct {
	Status  int
	Message string
}

func (e *RuleError) Error() string { return e.Message }

func ValidateCreateReturnInput(req model.CreateReturnRequest) map[string]string {
	errs := make(map[string]string)

	if req.ReportID <= 0 {
		errs["report_id"] = "report_id is required"
	}
	if req.ClaimID != nil && *req.ClaimID <= 0 {
		errs["claim_id"] = "claim_id must be a positive number"
	}
	if req.Notes != nil && len(strings.TrimSpace(*req.Notes)) > maxReturnNotesLength {
		errs["notes"] = "notes must be at most 500 characters"
	}

	return errs
}

// NormalizeReturnNotes merapikan catatan; catatan kosong menjadi nil (NULL di database).
func NormalizeReturnNotes(notes *string) *string {
	if notes == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*notes)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

// ResolveReturnReceiver menerapkan dua skenario serah terima dan menentukan penerima barang.
//
//	Skenario 1 (found): claim_id wajib, klaim 'approved' milik laporan ini,
//	                    status laporan 'claimed'. Penerima = pengklaim.
//	Skenario 2 (lost):  claim_id harus kosong, status laporan 'open'.
//	                    Penerima = pemilik laporan kehilangan.
func ResolveReturnReceiver(report model.ReturnReportInfo, claimID *int, claim *model.ReturnClaimInfo) (int, *RuleError) {
	switch report.ReportType {
	case "found":
		if claimID == nil {
			return 0, &RuleError{http.StatusBadRequest, "claim_id is required for a found report"}
		}
		if claim == nil || claim.ReportID != report.ID {
			return 0, &RuleError{http.StatusUnprocessableEntity, "claim does not belong to this report"}
		}
		if claim.Status != "approved" {
			return 0, &RuleError{http.StatusConflict, "claim must be approved before handover"}
		}
		if report.Status != "claimed" {
			return 0, &RuleError{http.StatusConflict, "found report must be in 'claimed' status"}
		}
		return claim.ClaimantID, nil

	case "lost":
		if claimID != nil {
			return 0, &RuleError{http.StatusBadRequest, "claim_id must be empty for a lost report"}
		}
		if report.Status != "open" {
			return 0, &RuleError{http.StatusConflict, "lost report must be in 'open' status"}
		}
		return report.OwnerID, nil
	}

	return 0, &RuleError{http.StatusUnprocessableEntity, "unsupported report type"}
}
