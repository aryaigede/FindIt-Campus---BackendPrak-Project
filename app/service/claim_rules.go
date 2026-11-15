package service

import (
	"errors"
	"strings"

	"findit-campus-backend/app/model"
)

var (
	ErrReportNotEligibleType = errors.New("claims can only be submitted for found items")
	ErrReportNotOpen         = errors.New("claims can only be submitted for open reports")
	ErrCannotClaimOwnReport  = errors.New("users cannot claim their own reports")
)

// ValidateCreateClaimInput validates input fields for submitting a claim
func ValidateCreateClaimInput(req model.CreateClaimRequest) map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(req.ProofDescription) == "" {
		errs["proof_description"] = "proof description is required"
	}

	return errs
}

// ValidateReviewClaimInput validates review status and notes
func ValidateReviewClaimInput(req model.ReviewClaimRequest) map[string]string {
	errs := make(map[string]string)

	status := strings.TrimSpace(strings.ToLower(req.Status))
	if status != model.ClaimStatusApproved && status != model.ClaimStatusRejected {
		errs["status"] = "status must be either 'approved' or 'rejected'"
	}

	if strings.TrimSpace(req.ReviewNote) == "" {
		errs["review_note"] = "review note is required"
	}

	return errs
}

// ValidateClaimEligibility verifies that a report is open, found, and not reported by claimant
func ValidateClaimEligibility(report *model.ClaimReportLookup, claimantID int) error {
	if report == nil {
		return errors.New("report lookup is nil")
	}

	if report.ReportType != "found" {
		return ErrReportNotEligibleType
	}

	if report.Status != "open" {
		return ErrReportNotOpen
	}

	if report.UserID == claimantID {
		return ErrCannotClaimOwnReport
	}

	return nil
}
