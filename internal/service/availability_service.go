package service

import (
	"AIClinicReceptionist/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type AvailabilityService struct {
	DoctorRepo      *repository.DoctorRepository
	ClinicRepo      *repository.ClinicHoursRepository
	AppointmentRepo *repository.AppointmentRepository
	Location        *time.Location
}

func NewAvailabilityService(
	docRepo *repository.DoctorRepository,
	clinicRepo *repository.ClinicHoursRepository,
	appRepo *repository.AppointmentRepository,
	loc *time.Location,
) *AvailabilityService {
	return &AvailabilityService{
		DoctorRepo:      docRepo,
		ClinicRepo:      clinicRepo,
		AppointmentRepo: appRepo,
		Location:        loc,
	}
}

func (s *AvailabilityService) GetAvailableSlots(
	ctx context.Context,
	doctorID uuid.UUID,
	inputDate time.Time,
) ([]string, error) {

	// 1️⃣ Convert to clinic timezone
	date := inputDate.In(s.Location)

	// 2️⃣ Reject Sundays
	if date.Weekday() == time.Sunday {
		return []string{}, nil
	}

	// 3️⃣ Fetch consultation duration
	durationMinutes, err := s.DoctorRepo.GetByID(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Fetch clinic hours
	weekday := int(date.Weekday())
	clinicHours, err := s.ClinicRepo.GetByWeekday(ctx, doctorID, weekday)
	if err != nil {
		return nil, err
	}

	// 5️⃣ Generate all possible slots
	allSlots := generateSlots(
		date,
		clinicHours.StartTime,
		clinicHours.EndTime,
		clinicHours.BreakStart,
		clinicHours.BreakEnd,
		durationMinutes,
		s.Location,
	)

	// 6️⃣ Fetch booked slots
	booked, err := s.AppointmentRepo.GetBookedSlots(ctx, doctorID, date)
	if err != nil {
		return nil, err
	}

	// 7️⃣ Remove booked slots
	available := filterBookedSlots(allSlots, booked)

	return available, nil
}

func generateSlots(
	date time.Time,
	start time.Time,
	end time.Time,
	breakStart *time.Time,
	breakEnd *time.Time,
	duration int,
	loc *time.Location,
) []string {

	var slots []string

	startTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		start.Hour(), start.Minute(), 0, 0, loc,
	)

	endTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		end.Hour(), end.Minute(), 0, 0, loc,
	)

	now := time.Now().In(loc)

	for current := startTime; current.Add(time.Duration(duration)*time.Minute).Before(endTime) ||
		current.Add(time.Duration(duration)*time.Minute).Equal(endTime); current = current.Add(time.Duration(duration) * time.Minute) {

		if current.Before(now) {
			continue
		}

		// Exclude break time
		if breakStart != nil && breakEnd != nil {
			bs := time.Date(date.Year(), date.Month(), date.Day(),
				breakStart.Hour(), breakStart.Minute(), 0, 0, loc)

			be := time.Date(date.Year(), date.Month(), date.Day(),
				breakEnd.Hour(), breakEnd.Minute(), 0, 0, loc)

			if current.Before(be) && current.Add(time.Duration(duration)*time.Minute).After(bs) {
				continue
			}
		}

		slots = append(slots, current.Format("15:04"))
	}

	return slots
}

func filterBookedSlots(all []string, booked []time.Time) []string {

	bookedMap := make(map[string]bool)

	for _, b := range booked {
		bookedMap[b.Format("15:04")] = true
	}

	var result []string

	for _, slot := range all {
		if !bookedMap[slot] {
			result = append(result, slot)
		}
	}

	return result
}
