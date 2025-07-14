package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/VibeTeam/fitness-tracker-backend/llm/suggester"
	"github.com/VibeTeam/fitness-tracker-backend/shared/middleware"
	"github.com/VibeTeam/fitness-tracker-backend/workout/repository"
)

// SuggestHandler returns AI-based workout suggestions.
type SuggestHandler struct {
	sessionRepo repository.WorkoutSessionRepository
	suggester   *suggester.Suggester
}

type suggestionResponse struct {
	Suggestion string `json:"suggestion"`
}

// errorResponse is used for Swagger documentation of error payloads.
type errorResponse struct {
	Error string `json:"error"`
}

func NewSuggestHandler(repo repository.WorkoutSessionRepository, sg *suggester.Suggester) *SuggestHandler {
	return &SuggestHandler{sessionRepo: repo, suggester: sg}
}

func (h *SuggestHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	g := r.Group("/suggest-workout")
	g.Use(auth)
	g.GET("", h.suggest)
}

// Suggest workout
// @Summary      Suggest next workout
// @Tags         workout-suggestions
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  suggestionResponse
// @Failure      500  {object}  errorResponse
// @Router       /suggest-workout [get]
func (h *SuggestHandler) suggest(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	// fetch last 10 sessions
	sessions, err := h.sessionRepo.ListByUser(c.Request.Context(), uid, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(sessions) == 0 {
		c.JSON(http.StatusOK, suggestionResponse{Suggestion: "No history yet. Start with a full-body beginner routine."})
		return
	}
	var parts []string
	for _, s := range sessions {
		parts = append(parts, "Session "+strconv.Itoa(int(s.ID))+": "+s.WorkoutType.Name)
	}
	history := strings.Join(parts, "\n")
	suggestion, err := h.suggester.Suggest(history)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, suggestionResponse{Suggestion: suggestion})
}
