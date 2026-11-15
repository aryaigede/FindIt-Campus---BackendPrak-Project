package service

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"findit-campus-backend/app/model"
)

// ValidateCreateCategoryInput validates payload for creating a category.
func ValidateCreateCategoryInput(req model.CreateCategoryRequest) map[string]string {
	errs := make(map[string]string)
	name := strings.TrimSpace(req.Name)
	if name == "" {
		errs["name"] = "name is required"
	} else if len(name) > 100 {
		errs["name"] = "name cannot exceed 100 characters"
	}
	return errs
}

// ValidateUpdateCategoryInput validates payload for updating a category.
func ValidateUpdateCategoryInput(req model.UpdateCategoryRequest) map[string]string {
	errs := make(map[string]string)
	name := strings.TrimSpace(req.Name)
	if name == "" {
		errs["name"] = "name is required"
	} else if len(name) > 100 {
		errs["name"] = "name cannot exceed 100 characters"
	}
	return errs
}

// ValidateCreateLocationInput validates payload for creating a location.
func ValidateCreateLocationInput(req model.CreateLocationRequest) map[string]string {
	errs := make(map[string]string)
	name := strings.TrimSpace(req.Name)
	if name == "" {
		errs["name"] = "name is required"
	} else if len(name) > 120 {
		errs["name"] = "name cannot exceed 120 characters"
	}
	return errs
}

// ValidateUpdateLocationInput validates payload for updating a location.
func ValidateUpdateLocationInput(req model.UpdateLocationRequest) map[string]string {
	errs := make(map[string]string)
	name := strings.TrimSpace(req.Name)
	if name == "" {
		errs["name"] = "name is required"
	} else if len(name) > 120 {
		errs["name"] = "name cannot exceed 120 characters"
	}
	return errs
}

// ParseEventDate attempts to parse an event_date string in common supported formats.
func ParseEventDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	trimmed := strings.TrimSpace(dateStr)
	for _, f := range formats {
		if t, err := time.Parse(f, trimmed); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid date format")
}

// ValidateCreateReportInput validates payload for creating an item report.
func ValidateCreateReportInput(req model.CreateReportRequest) (map[string]string, *time.Time) {
	errs := make(map[string]string)

	if req.CategoryID <= 0 {
		errs["category_id"] = "category_id is required"
	}

	if req.LocationID <= 0 {
		errs["location_id"] = "location_id is required"
	}

	reportType := strings.TrimSpace(req.ReportType)
	if reportType == "" {
		errs["report_type"] = "report_type is required"
	} else if reportType != model.ReportTypeLost && reportType != model.ReportTypeFound {
		errs["report_type"] = "report_type must be either 'lost' or 'found'"
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		errs["title"] = "title is required"
	} else if len(title) > 180 {
		errs["title"] = "title cannot exceed 180 characters"
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		errs["description"] = "description is required"
	}

	if req.CustodyNote != nil && len(*req.CustodyNote) > 150 {
		errs["custody_note"] = "custody_note cannot exceed 150 characters"
	}

	var parsedDate *time.Time
	dateStr := strings.TrimSpace(req.EventDate)
	if dateStr == "" {
		errs["event_date"] = "event_date is required"
	} else {
		parsed, err := ParseEventDate(dateStr)
		if err != nil {
			errs["event_date"] = "valid event_date is required (format: YYYY-MM-DD or RFC3339)"
		} else if parsed.After(time.Now().Add(5 * time.Minute)) {
			errs["event_date"] = "event_date cannot be in the future"
		} else {
			parsedDate = &parsed
		}
	}

	return errs, parsedDate
}

// ValidateUpdateReportInput validates payload for updating an item report.
func ValidateUpdateReportInput(req model.UpdateReportRequest) (map[string]string, *time.Time) {
	errs := make(map[string]string)

	if req.CategoryID <= 0 {
		errs["category_id"] = "category_id is required"
	}

	if req.LocationID <= 0 {
		errs["location_id"] = "location_id is required"
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		errs["title"] = "title is required"
	} else if len(title) > 180 {
		errs["title"] = "title cannot exceed 180 characters"
	}

	description := strings.TrimSpace(req.Description)
	if description == "" {
		errs["description"] = "description is required"
	}

	if req.CustodyNote != nil && len(*req.CustodyNote) > 150 {
		errs["custody_note"] = "custody_note cannot exceed 150 characters"
	}

	var parsedDate *time.Time
	dateStr := strings.TrimSpace(req.EventDate)
	if dateStr == "" {
		errs["event_date"] = "event_date is required"
	} else {
		parsed, err := ParseEventDate(dateStr)
		if err != nil {
			errs["event_date"] = "valid event_date is required (format: YYYY-MM-DD or RFC3339)"
		} else if parsed.After(time.Now().Add(5 * time.Minute)) {
			errs["event_date"] = "event_date cannot be in the future"
		} else {
			parsedDate = &parsed
		}
	}

	return errs, parsedDate
}

// ValidateAddImageInput validates payload for adding an image URL.
func ValidateAddImageInput(req model.AddImageRequest) map[string]string {
	errs := make(map[string]string)
	rawURL := strings.TrimSpace(req.ImageURL)
	if rawURL == "" {
		errs["image_url"] = "image_url is required"
		return errs
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		errs["image_url"] = "image_url must be a valid HTTP or HTTPS URL"
	}

	return errs
}

// ValidateReportFilter validates query parameters for listing reports.
func ValidateReportFilter(reportType, status *string, page, limit int) map[string]string {
	errs := make(map[string]string)

	if reportType != nil {
		rt := strings.TrimSpace(*reportType)
		if rt != model.ReportTypeLost && rt != model.ReportTypeFound {
			errs["report_type"] = "report_type must be either 'lost' or 'found'"
		}
	}

	if status != nil {
		st := strings.TrimSpace(*status)
		validStatuses := map[string]bool{
			model.ReportStatusOpen:     true,
			model.ReportStatusClaimed:  true,
			model.ReportStatusReturned: true,
			model.ReportStatusClosed:   true,
		}
		if !validStatuses[st] {
			errs["status"] = "status must be one of 'open', 'claimed', 'returned', 'closed'"
		}
	}

	if page < 1 {
		errs["page"] = "page must be greater than or equal to 1"
	}

	if limit < 1 || limit > 100 {
		errs["limit"] = "limit must be between 1 and 100"
	}

	return errs
}

// CanUserModifyReport checks if user is allowed to edit, close, or add images to a report.
func CanUserModifyReport(reportUserID int, currentUserID int, currentUserRole string) bool {
	if currentUserRole == "admin" {
		return true
	}
	return reportUserID == currentUserID
}

// CanCloseReport checks if report can transition to closed.
func CanCloseReport(currentStatus string) bool {
	return currentStatus == model.ReportStatusOpen
}

// CanUpdateReport checks if report can be updated.
func CanUpdateReport(currentStatus string) bool {
	return currentStatus == model.ReportStatusOpen
}

// ValidateStatusTransition validates whether status transition is allowed.
func ValidateStatusTransition(currentStatus, newStatus string) bool {
	if currentStatus == newStatus {
		return false
	}
	switch currentStatus {
	case model.ReportStatusOpen:
		return newStatus == model.ReportStatusClaimed || newStatus == model.ReportStatusClosed
	case model.ReportStatusClaimed:
		return newStatus == model.ReportStatusReturned || newStatus == model.ReportStatusOpen
	default:
		return false
	}
}
