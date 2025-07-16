package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/VibeTeam/fitness-tracker-backend/user/auth"
)

// newRouter returns a Gin engine with the Auth middleware wired
// plus a simple GET /ping endpoint that echoes the user ID if present.
func newRouter(mgr *auth.Manager) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Auth(mgr))
	r.GET("/ping", func(c *gin.Context) {
		if uid, ok := UserID(c); ok {
			c.String(http.StatusOK, "uid=%d", uid)
			return
		}
		c.String(http.StatusOK, "no‑uid")
	})
	return r
}

func TestAuthMiddleware_Success(t *testing.T) {
	mgr := auth.NewManager("access‑key", "refresh‑key",
		time.Minute, time.Hour)
	access, _, err := mgr.NewTokens(42)
	require.NoError(t, err)

	r := newRouter(mgr)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "uid=42", w.Body.String())
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	mgr := auth.NewManager("a", "b", time.Minute, time.Hour)
	r := newRouter(mgr)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mgr := auth.NewManager("x", "y", time.Minute, time.Hour)
	r := newRouter(mgr)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
