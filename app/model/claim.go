package model

import "time"

// Status constants for Claim
const (
	ClaimStatusPending  = "pending"
	ClaimStatusApproved = "approved"
	ClaimStatusRejected = "rejected"
)

// Claim represents the claims table entity
type Claim struct {
	ID               int        `json:"id"`
	ReportID         int        `json:"report_id"`
	ClaimantID       int        `json:"claimant_id"`
	ReviewedBy       *int       `json:"reviewed_by,omitempty"`
	ProofDescription string     `json:"proof_description"`
	ProofImageURL    *string    `json:"proof_image_url,omitempty"`
	Status           string     `json:"status"`
	ReviewNote       *string    `json:"review_note,omitempty"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Optional joined fields for richer API responses
	ClaimantName *string `json:"claimant_name,omitempty"`
	ReviewerName *string `json:"reviewer_name,omitempty"`
	ReportTitle  *string `json:"report_title,omitempty"`
}

// CreateClaimRequest represents input payload for submitting a claim
type CreateClaimRequest struct {
	ProofDescription string  `json:"proof_description"`
	ProofImageURL    *string `json:"proof_image_url,omitempty"`
}

// ReviewClaimRequest represents input payload for reviewing a claim
type ReviewClaimRequest struct {
	Status     string `json:"status"`
	ReviewNote string `json:"review_note"`
}

// ClaimReportLookup represents minimal item_reports data needed for claim validation
type ClaimReportLookup struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	ReportType string `json:"report_type"`
	Status     string `json:"status"`
}
