package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/VibeTeam/fitness-tracker-backend/user/auth"
)

func TestTokenLifecycle(t *testing.T) {
	mgr := auth.NewManager("access‑key", "refresh‑key", time.Minute, 24*time.Hour)

	access, refresh, err := mgr.NewTokens(42)
	require.NoError(t, err)

	uid, err := mgr.ValidateAccessToken(access)
	require.NoError(t, err)
	require.Equal(t, int32(42), uid)

	_, _, err = mgr.RefreshTokens(refresh)
	require.NoError(t, err)
}
