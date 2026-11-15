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

type ReportService struct {
	reportRepo repository.ReportRepository
	cfg        *config.Config
}

func NewReportService(reportRepo repository.ReportRepository, cfg *config.Config) *ReportService {
	return &ReportService{
		reportRepo: reportRepo,
		cfg:        cfg,
	}
}

// Categories

func (s *ReportService) ListCategories(c *fiber.Ctx) error {
	categories, err := s.reportRepo.FindAllCategories(c.Context())
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve categories")
	}
	return helper.Success(c, "categories retrieved successfully", categories)
}

func (s *ReportService) GetCategory(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid category id")
	}

	cat, err := s.reportRepo.FindCategoryByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "category not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve category")
	}

	return helper.Success(c, "category retrieved successfully", cat)
}

func (s *ReportService) CreateCategory(c *fiber.Ctx) error {
	var req model.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateCreateCategoryInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	name := strings.TrimSpace(req.Name)
	var desc *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		desc = &trimmed
	}

	cat := &model.Category{
		Name:        name,
		Description: desc,
	}

	if err := s.reportRepo.CreateCategory(c.Context(), cat); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return helper.Fail(c, fiber.StatusConflict, "category name already exists")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to create category")
	}

	return helper.Created(c, "category created successfully", cat)
}

func (s *ReportService) UpdateCategory(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid category id")
	}

	var req model.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateUpdateCategoryInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	name := strings.TrimSpace(req.Name)
	var desc *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		desc = &trimmed
	}

	cat := &model.Category{
		ID:          id,
		Name:        name,
		Description: desc,
	}

	if err := s.reportRepo.UpdateCategory(c.Context(), cat); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "category not found")
		}
		if errors.Is(err, repository.ErrConflict) {
			return helper.Fail(c, fiber.StatusConflict, "category name already exists")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to update category")
	}

	updatedCat, err := s.reportRepo.FindCategoryByID(c.Context(), id)
	if err == nil && updatedCat != nil {
		return helper.Success(c, "category updated successfully", updatedCat)
	}

	return helper.Success(c, "category updated successfully", cat)
}

func (s *ReportService) DeleteCategory(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid category id")
	}

	if err := s.reportRepo.DeleteCategory(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "category not found")
		}
		if errors.Is(err, repository.ErrForeignKeyConstraint) {
			return helper.Fail(c, fiber.StatusConflict, "category cannot be deleted because it is still in use by reports")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to delete category")
	}

	return helper.Success(c, "category deleted successfully", nil)
}

// Locations

func (s *ReportService) ListLocations(c *fiber.Ctx) error {
	locations, err := s.reportRepo.FindAllLocations(c.Context())
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve locations")
	}
	return helper.Success(c, "locations retrieved successfully", locations)
}

func (s *ReportService) GetLocation(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid location id")
	}

	loc, err := s.reportRepo.FindLocationByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "location not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve location")
	}

	return helper.Success(c, "location retrieved successfully", loc)
}

func (s *ReportService) CreateLocation(c *fiber.Ctx) error {
	var req model.CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateCreateLocationInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	name := strings.TrimSpace(req.Name)
	var desc *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		desc = &trimmed
	}

	loc := &model.Location{
		Name:        name,
		Description: desc,
	}

	if err := s.reportRepo.CreateLocation(c.Context(), loc); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			return helper.Fail(c, fiber.StatusConflict, "location name already exists")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to create location")
	}

	return helper.Created(c, "location created successfully", loc)
}

func (s *ReportService) UpdateLocation(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid location id")
	}

	var req model.UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateUpdateLocationInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	name := strings.TrimSpace(req.Name)
	var desc *string
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		desc = &trimmed
	}

	loc := &model.Location{
		ID:          id,
		Name:        name,
		Description: desc,
	}

	if err := s.reportRepo.UpdateLocation(c.Context(), loc); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "location not found")
		}
		if errors.Is(err, repository.ErrConflict) {
			return helper.Fail(c, fiber.StatusConflict, "location name already exists")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to update location")
	}

	updatedLoc, err := s.reportRepo.FindLocationByID(c.Context(), id)
	if err == nil && updatedLoc != nil {
		return helper.Success(c, "location updated successfully", updatedLoc)
	}

	return helper.Success(c, "location updated successfully", loc)
}

func (s *ReportService) DeleteLocation(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid location id")
	}

	if err := s.reportRepo.DeleteLocation(c.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "location not found")
		}
		if errors.Is(err, repository.ErrForeignKeyConstraint) {
			return helper.Fail(c, fiber.StatusConflict, "location cannot be deleted because it is still in use by reports")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to delete location")
	}

	return helper.Success(c, "location deleted successfully", nil)
}
