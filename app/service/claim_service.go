package service

import (
	"errors"
	"strings"

	"findit-campus-backend/app/model"
	"findit-campus-backend/app/repository"
	"findit-campus-backend/config"
	"findit-campus-backend/helper"

	"github.com/gofiber/fiber/v2"
)

type ClaimService struct {
	claimRepo repository.ClaimRepository
	cfg       *config.Config
}

func NewClaimService(claimRepo repository.ClaimRepository, cfg *config.Config) *ClaimService {
	return &ClaimService{claimRepo: claimRepo, cfg: cfg}
}

// SubmitClaim handles POST /api/v1/reports/:id/claims
func (s *ClaimService) SubmitClaim(c *fiber.Ctx) error {
	reportID, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid report id")
	}

	claimantID, ok := c.Locals("user_id").(int)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req model.CreateClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateCreateClaimInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	report, err := s.claimRepo.FindReportForClaim(c.Context(), reportID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to verify report")
	}

	if err := ValidateClaimEligibility(report, claimantID); err != nil {
		if errors.Is(err, ErrCannotClaimOwnReport) {
			return helper.Fail(c, fiber.StatusForbidden, err.Error())
		}
		return helper.Fail(c, fiber.StatusBadRequest, err.Error())
	}

	claim := &model.Claim{
		ReportID:         reportID,
		ClaimantID:       claimantID,
		ProofDescription: strings.TrimSpace(req.ProofDescription),
		ProofImageURL:    req.ProofImageURL,
	}

	if err := s.claimRepo.Create(c.Context(), claim); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return helper.Fail(c, fiber.StatusConflict, "a pending claim already exists for this report by the user")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to submit claim")
	}

	return helper.Created(c, "claim submitted successfully", claim)
}

// GetMyClaims handles GET /api/v1/claims/me
func (s *ClaimService) GetMyClaims(c *fiber.Ctx) error {
	claimantID, ok := c.Locals("user_id").(int)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	q := helper.ParseListQuery(c)
	claims, total, err := s.claimRepo.FindByClaimantID(c.Context(), claimantID, q)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve claims")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = int((total + int64(q.Limit) - 1) / int64(q.Limit))
	}

	meta := helper.PaginationMeta{
		Page:       q.Page,
		Limit:      q.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return helper.SuccessList(c, "claims retrieved successfully", claims, meta)
}

// ListClaims handles GET /api/v1/claims (security/admin)
func (s *ClaimService) ListClaims(c *fiber.Ctx) error {
	status := strings.TrimSpace(strings.ToLower(c.Query("status", "")))

	q := helper.ParseListQuery(c)
	claims, total, err := s.claimRepo.FindAll(c.Context(), status, q)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve claims")
	}

	totalPages := 0
	if q.Limit > 0 {
		totalPages = int((total + int64(q.Limit) - 1) / int64(q.Limit))
	}

	meta := helper.PaginationMeta{
		Page:       q.Page,
		Limit:      q.Limit,
		TotalItems: total,
		TotalPages: totalPages,
	}

	return helper.SuccessList(c, "claims retrieved successfully", claims, meta)
}

// ReviewClaim handles PATCH /api/v1/claims/:id/review (security/admin)
func (s *ClaimService) ReviewClaim(c *fiber.Ctx) error {
	claimID, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid claim id")
	}

	reviewerID, ok := c.Locals("user_id").(int)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req model.ReviewClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateReviewClaimInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	newStatus := strings.TrimSpace(strings.ToLower(req.Status))
	reviewNote := strings.TrimSpace(req.ReviewNote)

	updatedClaim, err := s.claimRepo.ReviewClaim(c.Context(), claimID, reviewerID, newStatus, reviewNote)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "claim not found")
		}
		if errors.Is(err, repository.ErrAlreadyReviewed) {
			return helper.Fail(c, fiber.StatusConflict, "claim has already been reviewed")
		}
		if errors.Is(err, repository.ErrConflict) {
			return helper.Fail(c, fiber.StatusConflict, "report is already claimed or another claim has been approved")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to review claim")
	}

	return helper.Success(c, "claim reviewed successfully", updatedClaim)
}
