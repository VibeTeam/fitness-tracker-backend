package models

import "time"

// MuscleGroup represents a primary muscle group targeted by a workout.
type MuscleGroup struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:text;not null"`
}

// WorkoutType represents a particular kind of workout (e.g., Bench Press) and the muscle group it trains.
type WorkoutType struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"type:text;not null"`
	MuscleGroupID uint   `gorm:"not null;index"`

	// Associations
	MuscleGroup *MuscleGroup `gorm:"foreignKey:MuscleGroupID"`
}

// WorkoutSession is a log entry for a completed workout instance performed by a user.
type WorkoutSession struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	WorkoutTypeID uint      `gorm:"not null;index"`
	UserID        uint      `gorm:"not null;index"`
	Datetime      time.Time `gorm:"not null"`

	// Associations
	WorkoutType *WorkoutType    `gorm:"foreignKey:WorkoutTypeID"`
	Details     []WorkoutDetail `gorm:"foreignKey:WorkoutSessionID"`
}

// WorkoutDetail stores arbitrary key-value data points for a workout session (e.g., reps, weight).
type WorkoutDetail struct {
	ID               uint   `gorm:"primaryKey;autoIncrement"`
	WorkoutSessionID uint   `gorm:"not null;index"`
	DetailName       string `gorm:"type:text;not null"`
	DetailValue      string `gorm:"type:text;not null"`
}
