package service

import (
	"errors"

	"findit-campus-backend/app/model"
	"findit-campus-backend/app/repository"
	"findit-campus-backend/helper"

	"github.com/gofiber/fiber/v2"
)

type ReturnHistoryService struct {
	repo repository.ReturnHistoryRepository
}

func NewReturnHistoryService(repo repository.ReturnHistoryRepository) *ReturnHistoryService {
	return &ReturnHistoryService{repo: repo}
}

// CreateReturn mencatat berita acara serah terima (security / admin).
func (s *ReturnHistoryService) CreateReturn(c *fiber.Ctx) error {
	actorID, ok := c.Locals("user_id").(int)
	if !ok {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req model.CreateReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}
	if errs := ValidateCreateReturnInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	report, err := s.repo.FindReportForReturn(c.Context(), req.ReportID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve report")
	}

	var claim *model.ReturnClaimInfo
	if req.ClaimID != nil {
		claim, err = s.repo.FindClaimForReturn(c.Context(), *req.ClaimID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve claim")
		}
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "claim not found")
		}
	}

	receiverID, ruleErr := ResolveReturnReceiver(*report, req.ClaimID, claim)
	if ruleErr != nil {
		return helper.Fail(c, ruleErr.Status, ruleErr.Message)
	}

	result, err := s.repo.CreateReturn(c.Context(), model.CreateReturnParams{
		ReportID:     report.ID,
		ClaimID:      req.ClaimID,
		HandedOverBy: actorID,
		ReceivedBy:   receiverID,
		OldStatus:    report.Status,
		Notes:        NormalizeReturnNotes(req.Notes),
	})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrReturnConflict):
			return helper.Fail(c, fiber.StatusConflict, "handover for this report or claim has already been recorded")
		case errors.Is(err, repository.ErrReportStatusChanged):
			return helper.Fail(c, fiber.StatusConflict, "report status changed, please retry")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to record handover")
	}

	return helper.Created(c, "handover recorded successfully", result)
}

// GetReturnByReport menampilkan detail serah terima berdasarkan ID laporan.
func (s *ReturnHistoryService) GetReturnByReport(c *fiber.Ctx) error {
	reportID, err := helper.ParamID(c, "report_id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid report_id")
	}

	result, err := s.repo.FindReturnByReportID(c.Context(), reportID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "handover record not found for this report")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve handover")
	}

	return helper.Success(c, "handover retrieved", result)
}

// ListHistories menampilkan linimasa audit perubahan status suatu laporan.
func (s *ReturnHistoryService) ListHistories(c *fiber.Ctx) error {
	reportID, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid report id")
	}

	exists, err := s.repo.ReportExists(c.Context(), reportID)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve report")
	}
	if !exists {
		return helper.Fail(c, fiber.StatusNotFound, "report not found")
	}

	histories, err := s.repo.ListHistoriesByReportID(c.Context(), reportID)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve report histories")
	}

	return helper.Success(c, "report histories retrieved", histories)
}
