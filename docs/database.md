# Database Schema

This document outlines the database schema for the Makerble API, which uses PostgreSQL.

## Tables

### `patients` table

This table stores information about patients.

| Column      | Type    | Constraints           | Description                               |
| :---------- | :------ | :-------------------- | :---------------------------------------- |
| `id`        | SERIAL  | PRIMARY KEY           | Unique identifier for the patient         |
| `email`     | TEXT    | NOT NULL, UNIQUE      | Unique email address of the patient       |
| `name`      | TEXT    | NOT NULL              | Full name of the patient                  |
| `address`   | TEXT    | NOT NULL              | Residential address of the patient        |
| `bloodtype` | TEXT    |                       | Blood type of the patient (e.g., 'A+', 'B-') |
| `doctor_id` | INTEGER |                       | Foreign Key referencing `staffs.id` (the doctor assigned to the patient) |

### `staffs` table

This table stores information about staff members, including doctors.

| Column     | Type   | Constraints           | Description                               |
| :--------- | :----- | :-------------------- | :---------------------------------------- |
| `id`       | SERIAL | PRIMARY KEY           | Unique identifier for the staff member    |
| `email`    | TEXT   | NOT NULL, UNIQUE      | Unique email address of the staff member  |
| `password` | TEXT   | NOT NULL              | Hashed password for authentication        |
| `name`     | TEXT   | NOT NULL              | Full name of the staff member             |
| `role`     | TEXT   | NOT NULL              | Role of the staff member (e.g., 'doctor', 'admin', 'nurse') |

## Relationships

*   **Patients to Staffs (Many-to-One):** The `patients.doctor_id` column is intended to link a patient to a specific staff member (doctor). This is a logical foreign key, though not explicitly enforced with a database constraint in the current migration scripts.
