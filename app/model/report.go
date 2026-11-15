package model

import "time"

// Category represents the categories table entity.
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents the payload for creating a category.
type CreateCategoryRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// UpdateCategoryRequest represents the payload for updating a category.
type UpdateCategoryRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// Location represents the locations table entity.
type Location struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateLocationRequest represents the payload for creating a location.
type CreateLocationRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// UpdateLocationRequest represents the payload for updating a location.
type UpdateLocationRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}
