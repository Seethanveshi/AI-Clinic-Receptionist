package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClinicHours struct {
	StartTime  time.Time
	EndTime    time.Time
	BreakStart *time.Time
	BreakEnd   *time.Time
}

type ClinicHoursRepository struct {
	DB *pgxpool.Pool
}

func NewClinicHoursRepository(db *pgxpool.Pool) *ClinicHoursRepository {
	return &ClinicHoursRepository{DB: db}
}

func (r *ClinicHoursRepository) GetByWeekday(
	ctx context.Context,
	doctorID uuid.UUID,
	weekday int,
) (*ClinicHours, error) {

	query := `
		SELECT start_time, end_time, break_start, break_end
		FROM clinic_hours
		WHERE doctor_id = $1
		AND weekday = $2
	`

	row := r.DB.QueryRow(ctx, query, doctorID, weekday)

	var ch ClinicHours

	err := row.Scan(
		&ch.StartTime,
		&ch.EndTime,
		&ch.BreakStart,
		&ch.BreakEnd,
	)
	if err != nil {
		return nil, err
	}

	return &ch, nil
}