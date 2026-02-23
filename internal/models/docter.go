package models

import (
	"time"

	"github.com/google/uuid"
)

type Doctor struct {
	ID                          uuid.UUID `gorm:"primaryKey"`
	Name                        string
	Specialization              string
	ConsultationDurationMinutes int
	CreatedAt                   time.Time
}
