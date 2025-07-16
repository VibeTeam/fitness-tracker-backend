package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/VibeTeam/fitness-tracker-backend/workout/handler"
	"github.com/VibeTeam/fitness-tracker-backend/workout/models"
	"github.com/VibeTeam/fitness-tracker-backend/workout/repository/gormrepository"
)

// -----------------------------------------------------------------------------
// helpers
// -----------------------------------------------------------------------------

// testRouter creates an isolated Gin engine backed by an in‑mem SQLite DB and
// registers only the routes needed for these tests.
func testRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	// migrate the minimal set of tables we touch
	require.NoError(t, db.AutoMigrate(&models.MuscleGroup{}, &models.WorkoutType{},
		&models.WorkoutSession{}, &models.WorkoutDetail{}))

	// repositories
	mgRepo := gormrepository.NewMuscleGroupRepository(db)
	wtRepo := gormrepository.NewWorkoutTypeRepository(db)
	wsRepo := gormrepository.NewWorkoutSessionRepository(db)
	wdRepo := gormrepository.NewWorkoutDetailRepository(db)

	// handlers
	mgHandler := handler.NewMuscleGroupHandler(mgRepo)
	wtHandler := handler.NewWorkoutTypeHandler(wtRepo)
	wsHandler := handler.NewWorkoutSessionHandler(wsRepo, wdRepo)

	// stub auth: the handlers under test don’t inspect user ID, so no‑op is fine
	noAuth := func(c *gin.Context) { c.Next() }

	r := gin.New()
	mgHandler.RegisterRoutes(r, noAuth)
	wtHandler.RegisterRoutes(r, noAuth)
	wsHandler.RegisterRoutes(r, noAuth)

	return r, db
}

func asJSON(t *testing.T, v any) *bytes.Buffer {
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

// -----------------------------------------------------------------------------
// Muscle‑group happy‑path CRUD
// -----------------------------------------------------------------------------

func TestMuscleGroupLifecycle(t *testing.T) {
	r, _ := testRouter(t)

	// CREATE
	var mgResp models.MuscleGroup
	{
		reqBody := asJSON(t, map[string]any{"name": "Chest"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/muscle-groups", reqBody)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &mgResp))
		require.NotZero(t, mgResp.ID)
	}

	// LIST
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/muscle-groups", nil)
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var list []models.MuscleGroup
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &list))
		require.Len(t, list, 1)
		require.Equal(t, mgResp.ID, list[0].ID)
	}

	// DELETE
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete,
			fmt.Sprintf("/muscle-groups/%d", mgResp.ID), nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusNoContent, w.Code)
	}
}

// -----------------------------------------------------------------------------
// Workout‑type lifecycle (needs a parent muscle‑group)
// -----------------------------------------------------------------------------

func TestWorkoutTypeLifecycle(t *testing.T) {
	r, _ := testRouter(t)

	// prerequisite muscle‑group
	var mg models.MuscleGroup
	{
		reqBody := asJSON(t, map[string]any{"name": "Legs"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/muscle-groups", reqBody)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &mg))
	}

	// CREATE workout‑type
	var wt models.WorkoutType
	{
		reqBody := asJSON(t, map[string]any{
			"name":            "Squat",
			"muscle_group_id": mg.ID,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/workout-types", reqBody)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &wt))
		require.Equal(t, mg.ID, wt.MuscleGroupID)
	}

	// DELETE
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete,
			fmt.Sprintf("/workout-types/%d", wt.ID), nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusNoContent, w.Code)
	}
}

// -----------------------------------------------------------------------------
// Workout‑session with details – proves associations survive round‑trip
// -----------------------------------------------------------------------------

func TestWorkoutSessionWithDetails(t *testing.T) {
	r, _ := testRouter(t)

	// create supporting muscle‑group & workout‑type
	var mg models.MuscleGroup
	{
		reqBody := asJSON(t, map[string]any{"name": "Back"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/muscle-groups", reqBody)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &mg))
	}

	var wt models.WorkoutType
	{
		reqBody := asJSON(t, map[string]any{
			"name":            "Deadlift",
			"muscle_group_id": mg.ID,
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/workout-types", reqBody)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &wt))
	}

	// create workout‑session
	var ws models.WorkoutSession
	{
		reqBody := asJSON(t, map[string]any{
			"workout_type_id": wt.ID,
			"datetime":        time.Now().UTC(),
		})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/workout-sessions", reqBody)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &ws))
	}

	// add a detail row
	{
		reqBody := asJSON(t, map[string]any{
			"detail_name":     "Reps",
			"detail_value":    "12",
			"session_id":      ws.ID,
			"workout_type":    wt.ID,
			"workout_type_id": wt.ID,
		})
		w := httptest.NewRecorder()
		target := fmt.Sprintf("/workout-sessions/%d/details", ws.ID)
		req, _ := http.NewRequest(http.MethodPost, target, reqBody)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	// fetch the session and ensure detail is present
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet,
			fmt.Sprintf("/workout-sessions/%d", ws.ID), nil)
		r.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var stored models.WorkoutSession
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &stored))
		require.Len(t, stored.Details, 1)
		require.Equal(t, "Reps", stored.Details[0].DetailName)
	}
}
