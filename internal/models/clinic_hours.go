package models

import (
	"time"

	"github.com/google/uuid"
)

type ClinicHours struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	DoctorID   uuid.UUID
	Weekday    int
	StartTime  time.Time
	EndTime    time.Time
	BreakStart time.Time
	BreakEnd   time.Time
}
