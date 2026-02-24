package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DoctorRepository struct {
	DB *pgxpool.Pool
}

func NewDoctorRepository(db *pgxpool.Pool) *DoctorRepository {
	return &DoctorRepository{DB: db}
}

func (r *DoctorRepository) GetByID(ctx context.Context, id uuid.UUID) (int, error) {
	var consultationDuration int

	query := `
		SELECT consultation_duration_minutes
		FROM doctors
		WHERE id = $1
	`

	err := r.DB.QueryRow(ctx, query, id).Scan(&consultationDuration)
	if err != nil {
		return 0, err
	}

	return consultationDuration, nil
}