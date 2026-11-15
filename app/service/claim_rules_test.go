package service_test

import (
	"testing"

	"findit-campus-backend/app/model"
	"findit-campus-backend/app/service"
)

func TestValidateCreateClaimInput(t *testing.T) {
	tests := []struct {
		name      string
		req       model.CreateClaimRequest
		expectErr bool
		errKey    string
	}{
		{
			name: "Valid input without image",
			req: model.CreateClaimRequest{
				ProofDescription: "Sticker Apple di sudut kanan laptop",
			},
			expectErr: false,
		},
		{
			name: "Valid input with image",
			req: model.CreateClaimRequest{
				ProofDescription: "Ada goresan kecil di bagian bawah",
				ProofImageURL:    func() *string { s := "https://example.com/proof.jpg"; return &s }(),
			},
			expectErr: false,
		},
		{
			name: "Empty proof description",
			req: model.CreateClaimRequest{
				ProofDescription: "   ",
			},
			expectErr: true,
			errKey:    "proof_description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := service.ValidateCreateClaimInput(tt.req)
			if tt.expectErr {
				if _, ok := errs[tt.errKey]; !ok {
					t.Errorf("expected error on key %s, got none: %v", tt.errKey, errs)
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("expected no errors, got: %v", errs)
				}
			}
		})
	}
}

func TestValidateReviewClaimInput(t *testing.T) {
	tests := []struct {
		name      string
		req       model.ReviewClaimRequest
		expectErr bool
		errKey    string
	}{
		{
			name: "Valid approved review",
			req: model.ReviewClaimRequest{
				Status:     "approved",
				ReviewNote: "Bukti fisik cocok dengan barang di pos satpam",
			},
			expectErr: false,
		},
		{
			name: "Valid rejected review",
			req: model.ReviewClaimRequest{
				Status:     "rejected",
				ReviewNote: "Ciri fisik berbeda",
			},
			expectErr: false,
		},
		{
			name: "Invalid status value",
			req: model.ReviewClaimRequest{
				Status:     "cancelled",
				ReviewNote: "Catatan valid",
			},
			expectErr: true,
			errKey:    "status",
		},
		{
			name: "Empty review note",
			req: model.ReviewClaimRequest{
				Status:     "approved",
				ReviewNote: "   ",
			},
			expectErr: true,
			errKey:    "review_note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := service.ValidateReviewClaimInput(tt.req)
			if tt.expectErr {
				if _, ok := errs[tt.errKey]; !ok {
					t.Errorf("expected error on key %s, got none: %v", tt.errKey, errs)
				}
			} else {
				if len(errs) > 0 {
					t.Errorf("expected no errors, got: %v", errs)
				}
			}
		})
	}
}

func TestValidateClaimEligibility(t *testing.T) {
	claimantID := 10

	t.Run("Eligible found open report", func(t *testing.T) {
		report := &model.ClaimReportLookup{
			ID:         1,
			UserID:     20, // different user
			ReportType: "found",
			Status:     "open",
		}
		err := service.ValidateClaimEligibility(report, claimantID)
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	})

	t.Run("Ineligible because report type is lost", func(t *testing.T) {
		report := &model.ClaimReportLookup{
			ID:         2,
			UserID:     20,
			ReportType: "lost",
			Status:     "open",
		}
		err := service.ValidateClaimEligibility(report, claimantID)
		if err != service.ErrReportNotEligibleType {
			t.Errorf("expected ErrReportNotEligibleType, got: %v", err)
		}
	})

	t.Run("Ineligible because report status is claimed", func(t *testing.T) {
		report := &model.ClaimReportLookup{
			ID:         3,
			UserID:     20,
			ReportType: "found",
			Status:     "claimed",
		}
		err := service.ValidateClaimEligibility(report, claimantID)
		if err != service.ErrReportNotOpen {
			t.Errorf("expected ErrReportNotOpen, got: %v", err)
		}
	})

	t.Run("Ineligible because claimant is report owner", func(t *testing.T) {
		report := &model.ClaimReportLookup{
			ID:         4,
			UserID:     claimantID, // same user
			ReportType: "found",
			Status:     "open",
		}
		err := service.ValidateClaimEligibility(report, claimantID)
		if err != service.ErrCannotClaimOwnReport {
			t.Errorf("expected ErrCannotClaimOwnReport, got: %v", err)
		}
	})
}
