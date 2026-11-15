package service

import (
	"strings"

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
