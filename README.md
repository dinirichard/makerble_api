# Makerble API

This is a Go-based API for managing patients and staff, with authentication and database migration capabilities.

## Features

-   Patient Management (CRUD operations)
-   Staff Management (Retrieve, Delete, Registration, Login)
-   Authentication (JWT-based)
-   Database Migrations (PostgreSQL)
-   Swagger/OpenAPI Documentation

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/dinirichard/makerble_api.git
cd makerble_api
```

### 2. Database Setup

Ensure your PostgreSQL server is running. The migration tool will create the `makerble` database if it doesn't exist.

To run migrations:

```bash
go run cmd/migrate/main.go up
```

To revert migrations:

```bash
go run cmd/migrate/main.go down
```

### 3. Run the API Server

```bash
go run ./cmd/api
```

The API server will start on `http://localhost:8080` (or the port configured in your environment).

### 4. Access API Documentation

Once the server is running, you can access the interactive API documentation (Swagger UI) at:

`http://localhost:8080/swagger/index.html`

## API Endpoints

### Public Endpoints

-   `GET /api/v1/patients` - Get all patients
-   `GET /api/v1/staffs/:id` - Get a staff member by ID
-   `DELETE /api/v1/staffs/:id` - Delete a staff member by ID
-   `POST /api/v1/auth/staff/register` - Register a new staff member
-   `POST /api/v1/auth/login` - Staff login

### Authenticated Endpoints (Requires JWT Token)

-   `POST /api/v1/auth/register` - Register a new patient
-   `POST /api/v1/patients` - Create a new patient
-   `GET /api/v1/patients/:id` - Get a patient by ID
-   `PUT /api/v1/patients/:id` - Update an existing patient
-   `DELETE /api/v1/patients/:id` - Delete an existing patient

## Database Schema

### `patients` table

| Column      | Type    | Description                                |
| :---------- | :------ | :----------------------------------------- |
| `id`        | SERIAL  | Primary Key, auto-incrementing             |
| `email`     | TEXT    | Unique email address of the patient        |
| `name`      | TEXT    | Name of the patient                        |
| `address`   | TEXT    | Address of the patient                     |
| `bloodtype` | TEXT    | Blood type of the patient (optional)       |
| `doctor_id` | INTEGER | ID of the associated doctor (staff member) |

### `staffs` table

| Column     | Type   | Description                               |
| :--------- | :----- | :---------------------------------------- |
| `id`       | SERIAL | Primary Key, auto-incrementing            |
| `email`    | TEXT   | Unique email address of the staff member  |
| `password` | TEXT   | Hashed password of the staff member       |
| `name`     | TEXT   | Name of the staff member                  |
| `role`     | TEXT   | Role of the staff member (e.g., "doctor") |

## Running Tests

```bash
go test -v ./cmd/api
```
