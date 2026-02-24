INSERT INTO doctors (name, specialization, consultation_duration_minutes)
VALUES ('Dr. Rahul', 'Cardiologist', 20)

INSERT INTO clinic_hours (doctor_id, weekday, start_time, end_time, break_start, break_end)
VALUES
('84fe826e-4695-433b-8903-fe42bf31b20c', 1, '10:00', '18:00', '13:00', '14:00'),
('84fe826e-4695-433b-8903-fe42bf31b20c', 2, '10:00', '18:00', '13:00', '14:00'),
('84fe826e-4695-433b-8903-fe42bf31b20c', 3, '10:00', '18:00', '13:00', '14:00'),
('84fe826e-4695-433b-8903-fe42bf31b20c', 4, '10:00', '18:00', '13:00', '14:00'),
('84fe826e-4695-433b-8903-fe42bf31b20c', 5, '10:00', '18:00', '13:00', '14:00'),
('84fe826e-4695-433b-8903-fe42bf31b20c', 6, '10:00', '18:00', '13:00', '14:00');