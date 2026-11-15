package service

import (
	"errors"
	"math"
	"strconv"
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

// Reports

func (s *ReportService) CreateReport(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(int)
	if !ok || userID <= 0 {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	var req model.CreateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	errs, parsedDate := ValidateCreateReportInput(req)
	if len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	var custodyNote *string
	if req.CustodyNote != nil {
		trimmed := strings.TrimSpace(*req.CustodyNote)
		custodyNote = &trimmed
	}

	report := &model.ItemReport{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		LocationID:  req.LocationID,
		ReportType:  strings.ToLower(strings.TrimSpace(req.ReportType)),
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		CustodyNote: custodyNote,
		EventDate:   *parsedDate,
		Status:      model.ReportStatusOpen,
	}

	if err := s.reportRepo.CreateReport(c.Context(), report); err != nil {
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return helper.FailValidation(c, "validation failed", map[string]string{"category_id": "category does not exist"})
		}
		if errors.Is(err, repository.ErrLocationNotFound) {
			return helper.FailValidation(c, "validation failed", map[string]string{"location_id": "location does not exist"})
		}
		if errors.Is(err, repository.ErrForeignKeyConstraint) {
			return helper.Fail(c, fiber.StatusBadRequest, "invalid reference to category or location")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to create report")
	}

	createdReport, err := s.reportRepo.FindReportByID(c.Context(), report.ID)
	if err == nil && createdReport != nil {
		return helper.Created(c, "report created successfully", createdReport)
	}

	return helper.Created(c, "report created successfully", report)
}

func (s *ReportService) ListReports(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return helper.FailValidation(c, "validation failed", map[string]string{"page": "invalid page number"})
	}

	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return helper.FailValidation(c, "validation failed", map[string]string{"limit": "invalid limit number"})
	}

	var reportTypePtr *string
	if rt := strings.TrimSpace(c.Query("report_type")); rt != "" {
		lowerRt := strings.ToLower(rt)
		reportTypePtr = &lowerRt
	}

	var statusPtr *string
	if st := strings.TrimSpace(c.Query("status")); st != "" {
		lowerSt := strings.ToLower(st)
		statusPtr = &lowerSt
	}

	var categoryIDPtr *int
	if catStr := strings.TrimSpace(c.Query("category_id")); catStr != "" {
		catID, err := strconv.Atoi(catStr)
		if err != nil || catID <= 0 {
			return helper.FailValidation(c, "validation failed", map[string]string{"category_id": "invalid category_id"})
		}
		categoryIDPtr = &catID
	}

	var locationIDPtr *int
	if locStr := strings.TrimSpace(c.Query("location_id")); locStr != "" {
		locID, err := strconv.Atoi(locStr)
		if err != nil || locID <= 0 {
			return helper.FailValidation(c, "validation failed", map[string]string{"location_id": "invalid location_id"})
		}
		locationIDPtr = &locID
	}

	if errs := ValidateReportFilter(reportTypePtr, statusPtr, page, limit); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	filter := model.ReportFilter{
		ReportType: reportTypePtr,
		CategoryID: categoryIDPtr,
		LocationID: locationIDPtr,
		Status:     statusPtr,
		Search:     strings.TrimSpace(c.Query("search")),
		Page:       page,
		Limit:      limit,
		Offset:     (page - 1) * limit,
	}

	reports, totalItems, err := s.reportRepo.FindReports(c.Context(), filter)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve reports")
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	meta := helper.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return helper.SuccessList(c, "reports retrieved successfully", reports, meta)
}

func (s *ReportService) GetReport(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid report id")
	}

	report, err := s.reportRepo.FindReportByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to retrieve report")
	}

	return helper.Success(c, "report retrieved successfully", report)
}

func (s *ReportService) UpdateReport(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid report id")
	}

	userID, _ := c.Locals("user_id").(int)
	userRole, _ := c.Locals("user_role").(string)
	if userID <= 0 {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	existing, err := s.reportRepo.FindReportByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to check report")
	}

	if !CanUserModifyReport(existing.UserID, userID, userRole) {
		return helper.Fail(c, fiber.StatusForbidden, "you do not have permission to modify this report")
	}

	if !CanUpdateReport(existing.Status) {
		return helper.Fail(c, fiber.StatusConflict, "only open reports can be updated")
	}

	var req model.UpdateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	errs, parsedDate := ValidateUpdateReportInput(req)
	if len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	var custodyNote *string
	if req.CustodyNote != nil {
		trimmed := strings.TrimSpace(*req.CustodyNote)
		custodyNote = &trimmed
	}

	updateData := &model.ItemReport{
		ID:          id,
		CategoryID:  req.CategoryID,
		LocationID:  req.LocationID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		CustodyNote: custodyNote,
		EventDate:   *parsedDate,
	}

	if err := s.reportRepo.UpdateReport(c.Context(), updateData); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		if errors.Is(err, repository.ErrStatusNotOpen) {
			return helper.Fail(c, fiber.StatusConflict, "only open reports can be updated")
		}
		if errors.Is(err, repository.ErrCategoryNotFound) {
			return helper.FailValidation(c, "validation failed", map[string]string{"category_id": "category does not exist"})
		}
		if errors.Is(err, repository.ErrLocationNotFound) {
			return helper.FailValidation(c, "validation failed", map[string]string{"location_id": "location does not exist"})
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to update report")
	}

	updatedReport, err := s.reportRepo.FindReportByID(c.Context(), id)
	if err == nil && updatedReport != nil {
		return helper.Success(c, "report updated successfully", updatedReport)
	}

	return helper.Success(c, "report updated successfully", updateData)
}

func (s *ReportService) CloseReport(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid report id")
	}

	userID, _ := c.Locals("user_id").(int)
	userRole, _ := c.Locals("user_role").(string)
	if userID <= 0 {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	existing, err := s.reportRepo.FindReportByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to check report")
	}

	if !CanUserModifyReport(existing.UserID, userID, userRole) {
		return helper.Fail(c, fiber.StatusForbidden, "you do not have permission to close this report")
	}

	if !CanCloseReport(existing.Status) {
		return helper.Fail(c, fiber.StatusConflict, "only open reports can be closed")
	}

	var body struct {
		Note string `json:"note"`
	}
	_ = c.BodyParser(&body)
	note := strings.TrimSpace(body.Note)
	if note == "" {
		note = "Laporan ditutup oleh " + userRole
	}

	if err := s.reportRepo.CloseReport(c.Context(), id, userID, note); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		if errors.Is(err, repository.ErrStatusNotOpen) {
			return helper.Fail(c, fiber.StatusConflict, "only open reports can be closed")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to close report")
	}

	closedReport, _ := s.reportRepo.FindReportByID(c.Context(), id)
	return helper.Success(c, "report closed successfully", closedReport)
}

func (s *ReportService) AddReportImage(c *fiber.Ctx) error {
	id, err := helper.ParamID(c, "id")
	if err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid report id")
	}

	userID, _ := c.Locals("user_id").(int)
	userRole, _ := c.Locals("user_role").(string)
	if userID <= 0 {
		return helper.Fail(c, fiber.StatusUnauthorized, "unauthorized")
	}

	existing, err := s.reportRepo.FindReportByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to check report")
	}

	if !CanUserModifyReport(existing.UserID, userID, userRole) {
		return helper.Fail(c, fiber.StatusForbidden, "you do not have permission to add images to this report")
	}

	var req model.AddImageRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(c, fiber.StatusBadRequest, "invalid JSON payload")
	}

	if errs := ValidateAddImageInput(req); len(errs) > 0 {
		return helper.FailValidation(c, "validation failed", errs)
	}

	count, err := s.reportRepo.CountImagesByReportID(c.Context(), id)
	if err != nil {
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to check image count")
	}
	if count >= 5 {
		return helper.Fail(c, fiber.StatusBadRequest, "maximum 5 images allowed per report")
	}

	img := &model.ItemImage{
		ReportID: id,
		ImageURL: strings.TrimSpace(req.ImageURL),
	}

	if err := s.reportRepo.AddImage(c.Context(), img); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return helper.Fail(c, fiber.StatusNotFound, "report not found")
		}
		return helper.Fail(c, fiber.StatusInternalServerError, "failed to add image")
	}

	return helper.Created(c, "image added successfully", img)
}
