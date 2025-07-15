package gormrepository

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/VibeTeam/fitness-tracker-backend/workout/models"
	"github.com/VibeTeam/fitness-tracker-backend/workout/repository"
)

// newTestDB creates an in‑memory SQLite DB and migrates all workout models.
func newTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening DB: %v", err)
	}

	if err := db.AutoMigrate(
		&models.MuscleGroup{},
		&models.WorkoutType{},
		&models.WorkoutSession{},
		&models.WorkoutDetail{},
	); err != nil {
		t.Fatalf("migrating schema: %v", err)
	}

	return db
}

/*
CRUD path for MuscleGroup
*/
func TestMuscleGroupCRUD(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)

	var mgRepo repository.MuscleGroupRepository = NewMuscleGroupRepository(db)

	// CREATE
	orig := &models.MuscleGroup{Name: "Chest"}
	if err := mgRepo.Create(ctx, orig); err != nil {
		t.Fatalf("create: %v", err)
	}
	if orig.ID == 0 {
		t.Fatalf("create: expected auto‑ID")
	}

	// READ (single)
	got, err := mgRepo.GetByID(ctx, orig.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Name != "Chest" {
		t.Fatalf("get by id: wrong name %q", got.Name)
	}

	// UPDATE
	got.Name = "Upper Chest"
	if err := mgRepo.Update(ctx, got); err != nil {
		t.Fatalf("update: %v", err)
	}

	// LIST / COUNT
	list, err := mgRepo.List(ctx, 10, 0)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: want 1, got %d (err=%v)", len(list), err)
	}
	cnt, _ := mgRepo.Count(ctx)
	if cnt != 1 {
		t.Fatalf("count: want 1, got %d", cnt)
	}

	// DELETE
	if err := mgRepo.Delete(ctx, got.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if c, _ := mgRepo.Count(ctx); c != 0 {
		t.Fatalf("delete: record still present")
	}
}

/*
Basic path for WorkoutType (requires a MuscleGroup FK)
*/
func TestWorkoutTypeCreateAndGet(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)

	mgRepo := NewMuscleGroupRepository(db)
	wtRepo := NewWorkoutTypeRepository(db)

	// prerequisite muscle group
	mg := &models.MuscleGroup{Name: "Back"}
	if err := mgRepo.Create(ctx, mg); err != nil {
		t.Fatalf("create mg: %v", err)
	}

	// CREATE workout type
	wt := &models.WorkoutType{
		Name:          "Deadlift",
		MuscleGroupID: mg.ID,
	}
	if err := wtRepo.Create(ctx, wt); err != nil {
		t.Fatalf("create wt: %v", err)
	}

	// READ (preloaded muscle group)
	got, err := wtRepo.GetByID(ctx, wt.ID)
	if err != nil {
		t.Fatalf("get wt: %v", err)
	}
	if got.Name != "Deadlift" || got.MuscleGroup.ID != mg.ID {
		t.Fatalf("unexpected workout type %+v", got)
	}
}

/*
Smoke test for WorkoutSession & WorkoutDetail to ensure associations persist.
*/
func TestWorkoutSessionWithDetails(t *testing.T) {
	ctx := context.Background()
	db := newTestDB(t)

	mgRepo := NewMuscleGroupRepository(db)
	wtRepo := NewWorkoutTypeRepository(db)
	wsRepo := NewWorkoutSessionRepository(db)
	wdRepo := NewWorkoutDetailRepository(db)

	// setup prerequisite records
	mg := &models.MuscleGroup{Name: "Legs"}
	if err := mgRepo.Create(ctx, mg); err != nil {
		t.Fatalf("create mg: %v", err)
	}
	wt := &models.WorkoutType{Name: "Squat", MuscleGroupID: mg.ID}
	if err := wtRepo.Create(ctx, wt); err != nil {
		t.Fatalf("create wt: %v", err)
	}

	// CREATE session
	session := &models.WorkoutSession{
		WorkoutTypeID: wt.ID,
		UserID:        42,
		Datetime:      time.Now(),
	}
	if err := wsRepo.Create(ctx, session); err != nil {
		t.Fatalf("create session: %v", err)
	}

	// ADD a detail row
	detail := &models.WorkoutDetail{
		WorkoutSessionID: session.ID,
		DetailName:       "Reps",
		DetailValue:      "12",
	}
	if err := wdRepo.Create(ctx, detail); err != nil {
		t.Fatalf("create detail: %v", err)
	}

	// Round‑trip: fetch session preloaded with details
	stored, err := wsRepo.GetByID(ctx, session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if len(stored.Details) != 1 || stored.Details[0].DetailValue != "12" {
		t.Fatalf("details not persisted: %+v", stored.Details)
	}
}
