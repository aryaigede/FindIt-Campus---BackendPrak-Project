package route_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"findit-campus-backend/app/model"
	"findit-campus-backend/app/repository"
	"findit-campus-backend/app/service"
	"findit-campus-backend/config"
	"findit-campus-backend/helper"
	"findit-campus-backend/middleware"
	"findit-campus-backend/route"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// MockClaimRepository simulates database interactions for testing endpoints
type MockClaimRepository struct {
	Reports     map[int]*model.ClaimReportLookup
	Claims      map[int]*model.Claim
	NextClaimID int
	Force23505  bool
}

func NewMockClaimRepository() *MockClaimRepository {
	return &MockClaimRepository{
		Reports:     make(map[int]*model.ClaimReportLookup),
		Claims:      make(map[int]*model.Claim),
		NextClaimID: 1,
	}
}

func (m *MockClaimRepository) Create(ctx context.Context, c *model.Claim) error {
	if m.Force23505 {
		return repository.ErrClaimConflict
	}
	// Check partial unique index: 1 pending per user per report
	for _, existing := range m.Claims {
		if existing.ReportID == c.ReportID && existing.ClaimantID == c.ClaimantID && existing.Status == model.ClaimStatusPending {
			return repository.ErrClaimConflict
		}
	}
	c.ID = m.NextClaimID
	m.NextClaimID++
	c.Status = model.ClaimStatusPending
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	m.Claims[c.ID] = c
	return nil
}

func (m *MockClaimRepository) FindByID(ctx context.Context, id int) (*model.Claim, error) {
	c, ok := m.Claims[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return c, nil
}

func (m *MockClaimRepository) FindByClaimantID(ctx context.Context, claimantID int, q helper.ListQuery) ([]model.Claim, int64, error) {
	var results []model.Claim
	for _, c := range m.Claims {
		if c.ClaimantID == claimantID {
			results = append(results, *c)
		}
	}
	return results, int64(len(results)), nil
}

func (m *MockClaimRepository) FindAll(ctx context.Context, status string, q helper.ListQuery) ([]model.Claim, int64, error) {
	var results []model.Claim
	for _, c := range m.Claims {
		if status == "" || c.Status == status {
			results = append(results, *c)
		}
	}
	return results, int64(len(results)), nil
}

func (m *MockClaimRepository) FindReportForClaim(ctx context.Context, reportID int) (*model.ClaimReportLookup, error) {
	r, ok := m.Reports[reportID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return r, nil
}

func (m *MockClaimRepository) ReviewClaim(ctx context.Context, claimID int, reviewerID int, newStatus string, reviewNote string) (*model.Claim, error) {
	c, ok := m.Claims[claimID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	if c.Status != model.ClaimStatusPending {
		return nil, repository.ErrAlreadyReviewed
	}

	if newStatus == model.ClaimStatusApproved {
		rep, ok := m.Reports[c.ReportID]
		if !ok {
			return nil, repository.ErrNotFound
		}
		if rep.Status != "open" {
			return nil, repository.ErrClaimConflict
		}
		rep.Status = "claimed"

		// auto-reject others
		for _, other := range m.Claims {
			if other.ReportID == c.ReportID && other.ID != c.ID && other.Status == model.ClaimStatusPending {
				other.Status = model.ClaimStatusRejected
				other.ReviewNote = func() *string { s := "Auto-rejected"; return &s }()
			}
		}
	}

	c.Status = newStatus
	c.ReviewedBy = &reviewerID
	c.ReviewNote = &reviewNote
	now := time.Now()
	c.ReviewedAt = &now
	c.UpdatedAt = now

	return c, nil
}

func generateTestToken(secret string, userID int, email string, role string) string {
	claims := middleware.JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func setupTestApp() (*fiber.App, *MockClaimRepository, *config.Config) {
	cfg := &config.Config{
		JWTSecret:       "test-secret-key-minimum-32-characters-long",
		JWTExpiresHours: 24,
	}

	mockRepo := NewMockClaimRepository()
	claimService := service.NewClaimService(mockRepo, cfg)

	app := fiber.New()
	api := app.Group("/api/v1")
	route.RegisterClaimRoutes(api, cfg, claimService)

	return app, mockRepo, cfg
}

func TestClaimEndpoints(t *testing.T) {
	app, mockRepo, cfg := setupTestApp()

	userToken := generateTestToken(cfg.JWTSecret, 10, "user@campus.ac.id", "user")
	securityToken := generateTestToken(cfg.JWTSecret, 2, "security@campus.ac.id", "security")
	adminToken := generateTestToken(cfg.JWTSecret, 1, "admin@campus.ac.id", "admin")

	// Seed sample report: found & open reported by user 99
	mockRepo.Reports[100] = &model.ClaimReportLookup{
		ID:         100,
		UserID:     99,
		ReportType: "found",
		Status:     "open",
	}

	// Seed sample report owned by user 10 (claimant)
	mockRepo.Reports[101] = &model.ClaimReportLookup{
		ID:         101,
		UserID:     10,
		ReportType: "found",
		Status:     "open",
	}

	// Seed sample report with type 'lost'
	mockRepo.Reports[102] = &model.ClaimReportLookup{
		ID:         102,
		UserID:     99,
		ReportType: "lost",
		Status:     "open",
	}

	t.Run("POST /api/v1/reports/:id/claims - Success submit claim", func(t *testing.T) {
		body, _ := json.Marshal(model.CreateClaimRequest{
			ProofDescription: "Stiker logo Google di pojok laptop ASUS",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reports/100/claims", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/reports/:id/claims - 409 Conflict duplicate pending claim", func(t *testing.T) {
		body, _ := json.Marshal(model.CreateClaimRequest{
			ProofDescription: "Mencoba submit klaim kedua saat masih pending",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reports/100/claims", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("expected 409 Conflict, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/reports/:id/claims - 403 Forbidden claiming own report", func(t *testing.T) {
		body, _ := json.Marshal(model.CreateClaimRequest{
			ProofDescription: "Klaim laporan sendiri",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reports/101/claims", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/v1/reports/:id/claims - 400 Bad Request report_type lost", func(t *testing.T) {
		body, _ := json.Marshal(model.CreateClaimRequest{
			ProofDescription: "Klaim laporan lost",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/reports/102/claims", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/claims/me - Success get my claims", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/claims/me", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/claims - 403 Forbidden for normal user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/claims", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/v1/claims - 200 OK for security role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/claims?status=pending", nil)
		req.Header.Set("Authorization", "Bearer "+securityToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("PATCH /api/v1/claims/:id/review - Success approved by security", func(t *testing.T) {
		body, _ := json.Marshal(model.ReviewClaimRequest{
			Status:     "approved",
			ReviewNote: "Bukti fisik diverifikasi di pos satpam",
		})
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/claims/1/review", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+securityToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("PATCH /api/v1/claims/:id/review - 409 Conflict cannot review already processed claim", func(t *testing.T) {
		body, _ := json.Marshal(model.ReviewClaimRequest{
			Status:     "rejected",
			ReviewNote: "Mencoba review ulang",
		})
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/claims/1/review", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+adminToken)

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("expected 409 Conflict, got %d", resp.StatusCode)
		}
	})
}
