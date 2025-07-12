CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    bloodtype TEXT,
    doctor_id INTEGER
)