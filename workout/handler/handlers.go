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

// -------- Muscle Groups --------

// MuscleGroupHandler handles CRUD operations for muscle groups.
type MuscleGroupHandler struct {
	repo repository.MuscleGroupRepository
}

func NewMuscleGroupHandler(repo repository.MuscleGroupRepository) *MuscleGroupHandler {
	return &MuscleGroupHandler{repo: repo}
}

func (h *MuscleGroupHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	mg := r.Group("/muscle-groups")
	mg.Use(auth)
	{
		mg.POST("", h.create)
		mg.GET("", h.list)
		mg.GET("/:id", h.getByID)
		mg.PUT("/:id", h.update)
		mg.DELETE("/:id", h.delete)
	}
}

type muscleGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

// create muscle group
// @Summary      Create muscle group
// @Tags         muscle-groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      muscleGroupRequest  true  "Muscle group"
// @Success      201      {object}  models.MuscleGroup
// @Failure      400      {object}  gin.H
// @Failure      500      {object}  gin.H
// @Router       /muscle-groups [post]
func (h *MuscleGroupHandler) create(c *gin.Context) {
	var req muscleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mg := &models.MuscleGroup{Name: req.Name}
	if err := h.repo.Create(c.Request.Context(), mg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mg)
}

// list muscle groups
// @Summary      List muscle groups
// @Tags         muscle-groups
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.MuscleGroup
// @Router       /muscle-groups [get]
func (h *MuscleGroupHandler) list(c *gin.Context) {
	groups, err := h.repo.List(c.Request.Context(), 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

// get muscle group
// @Summary      Get muscle group by ID
// @Tags         muscle-groups
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "MuscleGroup ID"
// @Success      200  {object}  models.MuscleGroup
// @Failure      400  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Router       /muscle-groups/{id} [get]
func (h *MuscleGroupHandler) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	mg, err := h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mg)
}

// update muscle group
// @Summary      Update muscle group
// @Tags         muscle-groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "MuscleGroup ID"
// @Param        payload  body      muscleGroupRequest true  "Update"
// @Success      200      {object}  models.MuscleGroup
// @Failure      400      {object}  gin.H
// @Failure      404      {object}  gin.H
// @Router       /muscle-groups/{id} [put]
func (h *MuscleGroupHandler) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	mg, err := h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var req muscleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mg.Name = req.Name
	if err := h.repo.Update(c.Request.Context(), mg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mg)
}

// delete muscle group
// @Summary      Delete muscle group
// @Tags         muscle-groups
// @Security     BearerAuth
// @Param        id   path      int  true  "MuscleGroup ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  gin.H
// @Router       /muscle-groups/{id} [delete]
func (h *MuscleGroupHandler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// -------- Workout Types --------

type WorkoutTypeHandler struct {
	repo repository.WorkoutTypeRepository
}

func NewWorkoutTypeHandler(repo repository.WorkoutTypeRepository) *WorkoutTypeHandler {
	return &WorkoutTypeHandler{repo: repo}
}

func (h *WorkoutTypeHandler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	wt := r.Group("/workout-types")
	wt.Use(auth)
	{
		wt.POST("", h.create)
		wt.GET("", h.list)
		wt.GET("/:id", h.getByID)
		wt.PUT("/:id", h.update)
		wt.DELETE("/:id", h.delete)
	}
}

type workoutTypeRequest struct {
	Name          string `json:"name" binding:"required"`
	MuscleGroupID uint   `json:"muscle_group_id" binding:"required"`
}

// create workout type
// @Summary      Create workout type
// @Tags         workout-types
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        payload  body      workoutTypeRequest  true  "Workout type"
// @Success      201      {object}  models.WorkoutType
// @Failure      400      {object}  gin.H
// @Failure      500      {object}  gin.H
// @Router       /workout-types [post]
func (h *WorkoutTypeHandler) create(c *gin.Context) {
	var req workoutTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wt := &models.WorkoutType{Name: req.Name, MuscleGroupID: req.MuscleGroupID}
	if err := h.repo.Create(c.Request.Context(), wt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// preload muscle group
	wt.MuscleGroup = &models.MuscleGroup{ID: req.MuscleGroupID}
	c.JSON(http.StatusCreated, wt)
}

// list workout types
// @Summary      List workout types
// @Tags         workout-types
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   models.WorkoutType
// @Router       /workout-types [get]
func (h *WorkoutTypeHandler) list(c *gin.Context) {
	types, err := h.repo.List(c.Request.Context(), 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

// get workout type
// @Summary      Get workout type by ID
// @Tags         workout-types
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "WorkoutType ID"
// @Success      200  {object}  models.WorkoutType
// @Failure      400  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Router       /workout-types/{id} [get]
func (h *WorkoutTypeHandler) getByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	wt, err := h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wt)
}

// update workout type
// @Summary      Update workout type
// @Tags         workout-types
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                 true  "WorkoutType ID"
// @Param        payload  body      workoutTypeRequest  true  "Update"
// @Success      200      {object}  models.WorkoutType
// @Failure      400      {object}  gin.H
// @Failure      404      {object}  gin.H
// @Router       /workout-types/{id} [put]
func (h *WorkoutTypeHandler) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	wt, err := h.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	var req workoutTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wt.Name = req.Name
	wt.MuscleGroupID = req.MuscleGroupID
	if err := h.repo.Update(c.Request.Context(), wt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wt)
}

// delete workout type
// @Summary      Delete workout type
// @Tags         workout-types
// @Security     BearerAuth
// @Param        id   path      int  true  "WorkoutType ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  gin.H
// @Router       /workout-types/{id} [delete]
func (h *WorkoutTypeHandler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.repo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// -------- Workout Sessions --------

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
	}
}

type workoutSessionRequest struct {
	WorkoutTypeID uint      `json:"workout_type_id" binding:"required"`
	Datetime      time.Time `json:"datetime"`
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
