package service

import (
	"AIClinicReceptionist/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type AppointmentService struct {
	DoctorRepo      *repository.DoctorRepository
	ClinicRepo      *repository.ClinicHoursRepository
	AppointmentRepo *repository.AppointmentRepository
	Location        *time.Location
}

func NewAppointmentService(
	docRepo *repository.DoctorRepository,
	clinicRepo *repository.ClinicHoursRepository,
	appRepo *repository.AppointmentRepository,
	loc *time.Location,
) *AppointmentService {
	return &AppointmentService{
		DoctorRepo:      docRepo,
		ClinicRepo:      clinicRepo,
		AppointmentRepo: appRepo,
		Location:        loc,
	}
}

func (s *AppointmentService) CreateAppointment(
	ctx context.Context,
	doctorID uuid.UUID,
	patientName string,
	phone string,
	visitType string,
	date time.Time,
	startTimeStr string,
) error {

	// Convert date to clinic timezone
	date = date.In(s.Location)

	// Reject Sunday
	if date.Weekday() == time.Sunday {
		return errors.New("clinic closed on Sunday")
	}

	// Parse start time (HH:MM)
	startTimeParsed, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		return errors.New("invalid time format")
	}

	startTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		startTimeParsed.Hour(), startTimeParsed.Minute(),
		0, 0, s.Location,
	)

	now := time.Now().In(s.Location)

	if !startTime.After(now) {
		return errors.New("cannot book appointment in the past")
	}

	// Fetch consultation duration
	durationMinutes, err := s.DoctorRepo.GetByID(ctx, doctorID)
	if err != nil {
		return err
	}

	endTime := startTime.Add(time.Duration(durationMinutes) * time.Minute)

	// Fetch clinic hours
	clinicHours, err := s.ClinicRepo.GetByWeekday(
		ctx,
		doctorID,
		int(date.Weekday()),
	)
	if err != nil {
		return err
	}

	// Validate within working hours
	if !isWithinWorkingHours(startTime, endTime, clinicHours, date, s.Location) {
		return errors.New("selected time outside working hours")
	}

	// Insert appointment (DB will enforce uniqueness)
	err = s.AppointmentRepo.Create(
		ctx,
		doctorID,
		patientName,
		phone,
		visitType,
		date,
		startTime,
		endTime,
	)

	if err != nil {
		return err
	}

	return nil
}

func isWithinWorkingHours(
	start time.Time,
	end time.Time,
	ch *repository.ClinicHours,
	date time.Time,
	loc *time.Location,
) bool {

	workStart := time.Date(date.Year(), date.Month(), date.Day(),
		ch.StartTime.Hour(), ch.StartTime.Minute(), 0, 0, loc)

	workEnd := time.Date(date.Year(), date.Month(), date.Day(),
		ch.EndTime.Hour(), ch.EndTime.Minute(), 0, 0, loc)

	if start.Before(workStart) || end.After(workEnd) {
		return false
	}

	if ch.BreakStart != nil && ch.BreakEnd != nil {
		breakStart := time.Date(date.Year(), date.Month(), date.Day(),
			ch.BreakStart.Hour(), ch.BreakStart.Minute(), 0, 0, loc)

		breakEnd := time.Date(date.Year(), date.Month(), date.Day(),
			ch.BreakEnd.Hour(), ch.BreakEnd.Minute(), 0, 0, loc)

		if start.Before(breakEnd) && end.After(breakStart) {
			return false
		}
	}

	return true
}

func (s *AppointmentService) CancelAppointment(
	ctx context.Context,
	doctorID uuid.UUID,
	phone string,
	date time.Time,
	startTimeStr string,
) error {

	date = date.In(s.Location)

	startParsed, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		return errors.New("invalid time format")
	}

	startTime := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		startParsed.Hour(),
		startParsed.Minute(),
		0,
		0,
		s.Location,
	)

	now := time.Now().In(s.Location)

	if !startTime.After(now) {
		return errors.New("cannot cancel appointment in the past")
	}

	return s.AppointmentRepo.Cancel(
		ctx,
		doctorID,
		phone,
		date,
		startTime,
	)
}
