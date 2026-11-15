package service

import (
	"strings"
	"testing"
	"time"

	"findit-campus-backend/app/model"
)

func TestValidateCreateCategoryInput(t *testing.T) {
	tests := []struct {
		name    string
		req     model.CreateCategoryRequest
		wantErr bool
		field   string
	}{
		{
			name: "valid category",
			req: model.CreateCategoryRequest{
				Name: "Elektronik",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: model.CreateCategoryRequest{
				Name: "   ",
			},
			wantErr: true,
			field:   "name",
		},
		{
			name: "name exceeds 100 characters",
			req: model.CreateCategoryRequest{
				Name: strings.Repeat("A", 101),
			},
			wantErr: true,
			field:   "name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateCreateCategoryInput(tc.req)
			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantErr && len(errs) > 0 {
				t.Fatalf("expected no error, got %v", errs)
			}
			if tc.wantErr && tc.field != "" {
				if _, ok := errs[tc.field]; !ok {
					t.Errorf("expected error field %q, got errs %v", tc.field, errs)
				}
			}
		})
	}
}

func TestValidateCreateLocationInput(t *testing.T) {
	tests := []struct {
		name    string
		req     model.CreateLocationRequest
		wantErr bool
		field   string
	}{
		{
			name: "valid location",
			req: model.CreateLocationRequest{
				Name: "Gedung Kuliah Bersama Lantai 2",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: model.CreateLocationRequest{
				Name: "",
			},
			wantErr: true,
			field:   "name",
		},
		{
			name: "name exceeds 120 characters",
			req: model.CreateLocationRequest{
				Name: strings.Repeat("L", 121),
			},
			wantErr: true,
			field:   "name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateCreateLocationInput(tc.req)
			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantErr && len(errs) > 0 {
				t.Fatalf("expected no error, got %v", errs)
			}
			if tc.wantErr && tc.field != "" {
				if _, ok := errs[tc.field]; !ok {
					t.Errorf("expected error field %q, got errs %v", tc.field, errs)
				}
			}
		})
	}
}

func TestParseEventDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "RFC3339 formatted date",
			input:   "2026-11-15T10:00:00Z",
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DD formatted date",
			input:   "2026-11-14",
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DD HH:MM:SS formatted date",
			input:   "2026-11-14 09:30:00",
			wantErr: false,
		},
		{
			name:    "invalid date string",
			input:   "not-a-date",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseEventDate(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error parsing %q, got nil", tc.input)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected success parsing %q, got %v", tc.input, err)
			}
		})
	}
}

func TestValidateCreateReportInput(t *testing.T) {
	futureDate := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	pastDate := time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	longCustody := strings.Repeat("C", 151)

	tests := []struct {
		name    string
		req     model.CreateReportRequest
		wantErr bool
		field   string
	}{
		{
			name: "valid lost report",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  2,
				ReportType:  "lost",
				Title:       "Dompet Kulit Coklat",
				Description: "Berisi kartu mahasiswa dan KTP",
				EventDate:   pastDate,
			},
			wantErr: false,
		},
		{
			name: "valid found report with custody note",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  2,
				ReportType:  "found",
				Title:       "Kunci Motor Honda",
				Description: "Ditemukan di bangku depan kantin",
				CustodyNote: func() *string { s := "Dititipkan di pos satpam"; return &s }(),
				EventDate:   "2026-11-01",
			},
			wantErr: false,
		},
		{
			name: "missing category_id",
			req: model.CreateReportRequest{
				CategoryID:  0,
				LocationID:  1,
				ReportType:  "lost",
				Title:       "Laptop Asus ROG",
				Description: "Hilang di perpustakaan",
				EventDate:   pastDate,
			},
			wantErr: true,
			field:   "category_id",
		},
		{
			name: "missing location_id",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  -1,
				ReportType:  "lost",
				Title:       "Laptop Asus ROG",
				Description: "Hilang di perpustakaan",
				EventDate:   pastDate,
			},
			wantErr: true,
			field:   "location_id",
		},
		{
			name: "invalid report_type",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  1,
				ReportType:  "stolen",
				Title:       "Laptop Asus ROG",
				Description: "Hilang di perpustakaan",
				EventDate:   pastDate,
			},
			wantErr: true,
			field:   "report_type",
		},
		{
			name: "empty title",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  1,
				ReportType:  "lost",
				Title:       "   ",
				Description: "Hilang di perpustakaan",
				EventDate:   pastDate,
			},
			wantErr: true,
			field:   "title",
		},
		{
			name: "title exceeds 180 chars",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  1,
				ReportType:  "lost",
				Title:       strings.Repeat("T", 181),
				Description: "Hilang di perpustakaan",
				EventDate:   pastDate,
			},
			wantErr: true,
			field:   "title",
		},
		{
			name: "empty description",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  1,
				ReportType:  "lost",
				Title:       "Dompet Coklat",
				Description: " ",
				EventDate:   pastDate,
			},
			wantErr: true,
			field:   "description",
		},
		{
			name: "custody note exceeds 150 chars",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  1,
				ReportType:  "found",
				Title:       "Dompet Coklat",
				Description: "Ditemukan di taman",
				CustodyNote: &longCustody,
				EventDate:   pastDate,
			},
			wantErr: true,
			field:   "custody_note",
		},
		{
			name: "future event date",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  1,
				ReportType:  "lost",
				Title:       "Dompet Coklat",
				Description: "Baru akan hilang besok",
				EventDate:   futureDate,
			},
			wantErr: true,
			field:   "event_date",
		},
		{
			name: "invalid event date format",
			req: model.CreateReportRequest{
				CategoryID:  1,
				LocationID:  1,
				ReportType:  "lost",
				Title:       "Dompet Coklat",
				Description: "Tanggal salah",
				EventDate:   "32-13-2026",
			},
			wantErr: true,
			field:   "event_date",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs, parsedDate := ValidateCreateReportInput(tc.req)
			if tc.wantErr {
				if len(errs) == 0 {
					t.Fatalf("expected validation error, got none")
				}
				if tc.field != "" {
					if _, ok := errs[tc.field]; !ok {
						t.Errorf("expected error field %q, got %v", tc.field, errs)
					}
				}
			} else {
				if len(errs) > 0 {
					t.Fatalf("expected no validation error, got %v", errs)
				}
				if parsedDate == nil {
					t.Fatalf("expected parsedDate to not be nil")
				}
			}
		})
	}
}

func TestValidateAddImageInput(t *testing.T) {
	tests := []struct {
		name    string
		req     model.AddImageRequest
		wantErr bool
	}{
		{
			name: "valid https url",
			req: model.AddImageRequest{
				ImageURL: "https://example.com/images/laptop.jpg",
			},
			wantErr: false,
		},
		{
			name: "valid http url",
			req: model.AddImageRequest{
				ImageURL: "http://example.com/photos/item.png",
			},
			wantErr: false,
		},
		{
			name: "empty url",
			req: model.AddImageRequest{
				ImageURL: "  ",
			},
			wantErr: true,
		},
		{
			name: "ftp url scheme",
			req: model.AddImageRequest{
				ImageURL: "ftp://example.com/file.jpg",
			},
			wantErr: true,
		},
		{
			name: "relative path instead of url",
			req: model.AddImageRequest{
				ImageURL: "/uploads/image.png",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateAddImageInput(tc.req)
			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation error, got none")
			}
			if !tc.wantErr && len(errs) > 0 {
				t.Fatalf("expected no validation error, got %v", errs)
			}
		})
	}
}

func TestValidateReportFilter(t *testing.T) {
	validType := "lost"
	invalidType := "unknown"
	validStatus := "open"
	invalidStatus := "in_review"

	tests := []struct {
		name       string
		reportType *string
		status     *string
		page       int
		limit      int
		wantErr    bool
		field      string
	}{
		{
			name:       "valid filter with defaults",
			reportType: nil,
			status:     nil,
			page:       1,
			limit:      10,
			wantErr:    false,
		},
		{
			name:       "valid filter with type and status",
			reportType: &validType,
			status:     &validStatus,
			page:       2,
			limit:      50,
			wantErr:    false,
		},
		{
			name:       "invalid report_type",
			reportType: &invalidType,
			status:     nil,
			page:       1,
			limit:      10,
			wantErr:    true,
			field:      "report_type",
		},
		{
			name:       "invalid status",
			reportType: nil,
			status:     &invalidStatus,
			page:       1,
			limit:      10,
			wantErr:    true,
			field:      "status",
		},
		{
			name:       "page less than 1",
			reportType: nil,
			status:     nil,
			page:       0,
			limit:      10,
			wantErr:    true,
			field:      "page",
		},
		{
			name:       "limit exceeds 100",
			reportType: nil,
			status:     nil,
			page:       1,
			limit:      101,
			wantErr:    true,
			field:      "limit",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateReportFilter(tc.reportType, tc.status, tc.page, tc.limit)
			if tc.wantErr {
				if len(errs) == 0 {
					t.Fatalf("expected validation error, got none")
				}
				if tc.field != "" {
					if _, ok := errs[tc.field]; !ok {
						t.Errorf("expected error field %q, got %v", tc.field, errs)
					}
				}
			} else {
				if len(errs) > 0 {
					t.Fatalf("expected no validation error, got %v", errs)
				}
			}
		})
	}
}

func TestCanUserModifyReport(t *testing.T) {
	tests := []struct {
		name            string
		reportUserID    int
		currentUserID   int
		currentUserRole string
		expected        bool
	}{
		{
			name:            "owner can modify",
			reportUserID:    42,
			currentUserID:   42,
			currentUserRole: "user",
			expected:        true,
		},
		{
			name:            "admin can modify any report",
			reportUserID:    42,
			currentUserID:   99,
			currentUserRole: "admin",
			expected:        true,
		},
		{
			name:            "different normal user cannot modify",
			reportUserID:    42,
			currentUserID:   99,
			currentUserRole: "user",
			expected:        false,
		},
		{
			name:            "security officer cannot modify unless owner or admin",
			reportUserID:    42,
			currentUserID:   77,
			currentUserRole: "security",
			expected:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CanUserModifyReport(tc.reportUserID, tc.currentUserID, tc.currentUserRole)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestCanCloseAndCanUpdateReport(t *testing.T) {
	tests := []struct {
		status    string
		canClose  bool
		canUpdate bool
	}{
		{status: model.ReportStatusOpen, canClose: true, canUpdate: true},
		{status: model.ReportStatusClaimed, canClose: false, canUpdate: false},
		{status: model.ReportStatusReturned, canClose: false, canUpdate: false},
		{status: model.ReportStatusClosed, canClose: false, canUpdate: false},
	}

	for _, tc := range tests {
		t.Run(tc.status, func(t *testing.T) {
			if got := CanCloseReport(tc.status); got != tc.canClose {
				t.Errorf("CanCloseReport(%q) = %v; expected %v", tc.status, got, tc.canClose)
			}
			if got := CanUpdateReport(tc.status); got != tc.canUpdate {
				t.Errorf("CanUpdateReport(%q) = %v; expected %v", tc.status, got, tc.canUpdate)
			}
		})
	}
}

func TestValidateStatusTransition(t *testing.T) {
	tests := []struct {
		current  string
		newStat  string
		expected bool
	}{
		{current: model.ReportStatusOpen, newStat: model.ReportStatusClaimed, expected: true},
		{current: model.ReportStatusOpen, newStat: model.ReportStatusClosed, expected: true},
		{current: model.ReportStatusClaimed, newStat: model.ReportStatusReturned, expected: true},
		{current: model.ReportStatusClaimed, newStat: model.ReportStatusOpen, expected: true},
		{current: model.ReportStatusClosed, newStat: model.ReportStatusOpen, expected: false},
		{current: model.ReportStatusReturned, newStat: model.ReportStatusOpen, expected: false},
		{current: model.ReportStatusOpen, newStat: model.ReportStatusOpen, expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.current+"->"+tc.newStat, func(t *testing.T) {
			result := ValidateStatusTransition(tc.current, tc.newStat)
			if result != tc.expected {
				t.Errorf("ValidateStatusTransition(%q, %q) = %v; expected %v", tc.current, tc.newStat, result, tc.expected)
			}
		})
	}
}
