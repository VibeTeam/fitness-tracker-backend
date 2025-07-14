package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/VibeTeam/fitness-tracker-backend/user/auth"
)

// Auth returns a Gin middleware that validates Bearer JWT access tokens using the provided token manager.
// On success, the middleware stores the authenticated user ID under the key "user_id" in the Gin context.
func Auth(tokenMgr *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		userID, err := tokenMgr.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Store the authenticated user ID for downstream handlers
		c.Set("user_id", uint(userID))
		c.Next()
	}
}

// UserID retrieves the authenticated user ID from Gin context if present.
func UserID(c *gin.Context) (uint, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
