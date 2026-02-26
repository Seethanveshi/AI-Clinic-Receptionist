CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE doctors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    specialization TEXT NOT NULL,
    consultation_duration_minutes INTEGER NOT NULL DEFAULT 20,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE clinic_hours (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    doctor_id UUID NOT NULL REFERENCES doctors(id) ON UPDATE CASCADE ON DELETE CASCADE,
    weekday INTEGER NOT NULL CHECK (weekday BETWEEN 0 AND 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    break_start TIME,
    break_end TIME
);

CREATE TYPE appointment_status AS ENUM (
    'scheduled',
    'cancelled',
    'rescheduled'
);

CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    doctor_id UUID NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    patient_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    visit_type TEXT,
    appointment_date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status appointment_status DEFAULT 'scheduled',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_doctor_slot_active
ON appointments (doctor_id, appointment_date, start_time)
WHERE status = 'scheduled';

CREATE INDEX idx_appointments_doctor_date
ON appointments (doctor_id, appointment_date);