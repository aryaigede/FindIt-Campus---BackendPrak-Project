package model

import "time"

// Constants for report types and statuses
const (
	ReportTypeLost  = "lost"
	ReportTypeFound = "found"

	ReportStatusOpen     = "open"
	ReportStatusClaimed  = "claimed"
	ReportStatusReturned = "returned"
	ReportStatusClosed   = "closed"
)

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

// ItemImage represents the item_images table entity.
type ItemImage struct {
	ID        int       `json:"id"`
	ReportID  int       `json:"report_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

// AddImageRequest represents the payload for adding an image to a report.
type AddImageRequest struct {
	ImageURL string `json:"image_url"`
}

// ItemReport represents the item_reports table entity with joined relation fields.
type ItemReport struct {
	ID           int         `json:"id"`
	UserID       int         `json:"user_id"`
	UserName     string      `json:"user_name,omitempty"`
	CategoryID   int         `json:"category_id"`
	CategoryName string      `json:"category_name,omitempty"`
	LocationID   int         `json:"location_id"`
	LocationName string      `json:"location_name,omitempty"`
	ReportType   string      `json:"report_type"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	CustodyNote  *string     `json:"custody_note,omitempty"`
	EventDate    time.Time   `json:"event_date"`
	Status       string      `json:"status"`
	Images       []ItemImage `json:"images,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// CreateReportRequest represents the payload for creating a new report.
type CreateReportRequest struct {
	CategoryID  int     `json:"category_id"`
	LocationID  int     `json:"location_id"`
	ReportType  string  `json:"report_type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CustodyNote *string `json:"custody_note"`
	EventDate   string  `json:"event_date"`
}

// UpdateReportRequest represents the payload for updating an existing report.
type UpdateReportRequest struct {
	CategoryID  int     `json:"category_id"`
	LocationID  int     `json:"location_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CustodyNote *string `json:"custody_note"`
	EventDate   string  `json:"event_date"`
}

// ReportFilter represents query parameters for filtering reports.
type ReportFilter struct {
	ReportType *string
	CategoryID *int
	LocationID *int
	Status     *string
	Search     string
	Page       int
	Limit      int
	Offset     int
}
