package models

import (
	"time"

	"github.com/google/uuid"
)

type AppointmentStatus string

const (
	StatusScheduled   AppointmentStatus = "scheduled"
	StatusCancelled   AppointmentStatus = "cancelled"
	StatusRescheduled AppointmentStatus = "rescheduled"
)

type Appointment struct {
	ID              uuid.UUID
	DoctorID        uuid.UUID
	PatientName     string
	Phone           string
	VisitType       string
	AppointmentDate time.Time
	StartTime       time.Time
	EndTime         time.Time
	Status          string
	CreatedAt       time.Time
}
