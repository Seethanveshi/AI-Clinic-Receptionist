package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppointmentRepository struct {
	DB *pgxpool.Pool
}

var ErrSlotAlreadyBooked = errors.New("slot already booked")

func NewAppointmentRepository(db *pgxpool.Pool) *AppointmentRepository {
	return &AppointmentRepository{DB: db}
}

func (r *AppointmentRepository) GetBookedSlots(
	ctx context.Context,
	doctorID uuid.UUID,
	date time.Time,
) ([]time.Time, error) {

	query := `
		SELECT start_time
		FROM appointments
		WHERE doctor_id = $1
		AND appointment_date = $2
		AND status = 'scheduled'
	`

	rows, err := r.DB.Query(ctx, query, doctorID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var booked []time.Time

	for rows.Next() {
		var t time.Time
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		booked = append(booked, t)
	}

	return booked, nil
}

func (r *AppointmentRepository) Create(
	ctx context.Context,
	doctorID uuid.UUID,
	patientName string,
	phone string,
	visitType string,
	date time.Time,
	start time.Time,
	end time.Time,
) error {

	query := `
		INSERT INTO appointments
		(doctor_id, patient_name, phone, visit_type, appointment_date, start_time, end_time)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := r.DB.Exec(ctx, query,
		doctorID,
		patientName,
		phone,
		visitType,
		date,
		start,
		end,
	)

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" &&
				pgErr.ConstraintName == "idx_doctor_slot" {
				return ErrSlotAlreadyBooked
			}
		}

		return err
	}

	return nil
}

func (r *AppointmentRepository) Cancel(
	ctx context.Context,
	doctorID uuid.UUID,
	phone string,
	date time.Time,
	startTime time.Time,
) error {

	query := `
		UPDATE appointments
		SET status = 'cancelled'
		WHERE doctor_id = $1
		AND phone = $2
		AND appointment_date = $3
		AND start_time = $4
		AND status = 'scheduled'
	`

	cmd, err := r.DB.Exec(
		ctx,
		query,
		doctorID,
		phone,
		date,
		startTime,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("appointment not found")
	}

	return nil
}