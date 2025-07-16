package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/VibeTeam/fitness-tracker-backend/user/auth"
	"github.com/VibeTeam/fitness-tracker-backend/user/handler"
	"github.com/VibeTeam/fitness-tracker-backend/user/models"
	"github.com/VibeTeam/fitness-tracker-backend/user/use_case"
)

/* ---------- minimal in‑memory UserRepository (implements Count) ------------- */

type memRepo struct {
	mu      sync.Mutex
	byEmail map[string]*models.User
	byID    map[uint]*models.User
	nextID  uint
}

func newMemRepo() *memRepo {
	return &memRepo{
		byEmail: make(map[string]*models.User),
		byID:    make(map[uint]*models.User),
		nextID:  1,
	}
}

func (r *memRepo) Create(_ context.Context, u *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.nextID
	r.nextID++
	cu := *u
	r.byEmail[u.Email] = &cu
	r.byID[u.ID] = &cu
	return nil
}

func (r *memRepo) GetByEmail(_ context.Context, email string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byEmail[email]
	if !ok {
		return nil, errors.New("not found")
	}
	cu := *u
	return &cu, nil
}

func (r *memRepo) GetByID(_ context.Context, id uint) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	cu := *u
	return &cu, nil
}

func (r *memRepo) List(_ context.Context, _, _ int) ([]*models.User, error) { return nil, nil }
func (r *memRepo) Update(_ context.Context, _ *models.User) error           { return nil }
func (r *memRepo) Delete(_ context.Context, _ uint) error                   { return nil }
func (r *memRepo) Count(_ context.Context) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.byID), nil
}

/* --------------------------------------------------------------------------- */

func TestLoginAndRefresh_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// backing store + service
	repo := newMemRepo()
	tm := auth.NewManager("access-key", "refresh-key", time.Minute, time.Hour)
	svc := use_case.NewAuthService(repo, tm)

	// register a real user so /auth/login can succeed
	_, _, err := svc.Register(context.Background(), "me@example.com", "pwd")
	require.NoError(t, err)

	// HTTP layer
	h := handler.NewAuthHandler(svc)
	router := gin.New()
	h.RegisterRoutes(router, func(c *gin.Context) {}) // no‑op auth

	/* -------- login -------- */
	loginBody := []byte(`{"email":"me@example.com","password":"pwd"}`)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var tokens map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &tokens))
	require.NotEmpty(t, tokens["access_token"])
	require.NotEmpty(t, tokens["refresh_token"])

	/* -------- refresh -------- */
	refBody := []byte(`{"refresh_token":"` + tokens["refresh_token"] + `"}`)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(refBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}
