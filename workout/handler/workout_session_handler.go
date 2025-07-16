package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/VibeTeam/fitness-tracker-backend/shared/middleware"
	"github.com/VibeTeam/fitness-tracker-backend/workout/models"
	"github.com/VibeTeam/fitness-tracker-backend/workout/repository"
)

type WorkoutSessionHandler struct {
	repo       repository.WorkoutSessionRepository
	detailRepo repository.WorkoutDetailRepository
}

func NewWorkoutSessionHandler(repo repository.WorkoutSessionRepository, detailRepo repository.WorkoutDetailRepository) *WorkoutSessionHandler {
	return &WorkoutSessionHandler{repo: repo, detailRepo: detailRepo}
}

func (h *WorkoutSessionHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	ws := r.Group("/workout-sessions")
	ws.Use(auth)
	{
		ws.POST("", h.create)
		ws.GET("", h.list)
		ws.GET("/:id", h.getByID)
		ws.DELETE("/:id", h.delete)
		ws.POST("/:id/details", h.addDetail)
	}
}

type workoutSessionRequest struct {
	WorkoutTypeID uint      `json:"workout_type_id" binding:"required"`
	Datetime      time.Time `json:"datetime"`
}

// detail request DTO
type workoutDetailRequest struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// create session
// @Summary      Create workout session
// @Tags         workout-sessions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      workoutSessionRequest  true  "Session"
// @Success      201      {object}  models.WorkoutSession
// @Failure      400      {object}  gin.H
// @Failure      500      {object}  gin.H
// @Router       /workout-sessions [post]
func (h *WorkoutSessionHandler) create(c *gin.Context) {
	var req workoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	if req.Datetime.IsZero() {
		req.Datetime = time.Now()
	}
	session := &models.WorkoutSession{UserID: uid, WorkoutTypeID: req.WorkoutTypeID, Datetime: req.Datetime}
	if err := h.repo.Create(c.Request.Context(), session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, session)
}

// add detail
// @Summary      Add detail to workout session
// @Tags         workout-sessions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "WorkoutSession ID"
// @Param        payload  body      workoutDetailRequest true  "Detail"
// @Success      201      {object}  models.WorkoutDetail
// @Failure      400      {object}  gin.H
// @Failure      404      {object}  gin.H
// @Failure      500      {object}  gin.H
// @Router       /workout-sessions/{id}/details [post]
func (h *WorkoutSessionHandler) addDetail(c *gin.Context) {
	// parse session ID
	sid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// load session and verify ownership
	session, err := h.repo.GetByID(c.Request.Context(), uint(sid))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	uid, _ := middleware.UserID(c)
	if session.UserID != uid {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var req workoutDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detail := &models.WorkoutDetail{
		WorkoutSessionID: session.ID,
		DetailName:       req.Name,
		DetailValue:      req.Value,
	}
	if err := h.detailRepo.Create(c.Request.Context(), detail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, detail)
}

// list sessions for user
// @Summary      List workout sessions for user
// @Tags         workout-sessions
// @Security     BearerAuth
// @Produce      json
// @Param        limit   query     int  false  "Limit"
// @Param        offset  query     int  false  "Offset"
// @Success      200  {array}   models.WorkoutSession
// @Router       /workout-sessions [get]
func (h *WorkoutSessionHandler) list(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sessions, err := h.repo.ListByUser(c.Request.Context(), uid, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

// get session
// @Summary      Get workout session by ID
// @Tags         workout-sessions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "WorkoutSession ID"
// @Success      200  {object}  models.WorkoutSession
// @Failure      400  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Router       /workout-sessions/{id} [get]
func (h *WorkoutSessionHandler) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	session, err := h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	uid, _ := middleware.UserID(c)
	if session.UserID != uid {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

// delete session
// @Summary      Delete workout session
// @Tags         workout-sessions
// @Security     BearerAuth
// @Param        id   path      int  true  "WorkoutSession ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  gin.H
// @Router       /workout-sessions/{id} [delete]
func (h *WorkoutSessionHandler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	session, err := h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	uid, _ := middleware.UserID(c)
	if session.UserID != uid {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
