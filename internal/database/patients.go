package database

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type PatientModel struct {
	DB *sql.DB
}

type Patient struct {
	Id 			int		`json:"id"`
	Email 		string 	`json:"email" binding:"required,email"`
	Name 		string 	`json:"name" binding:"required,min=3"`
	Address 	string 	`json:"address" binding:"required,min=5"`
	Bloodtype	string	`json:"bloodtype"`
	Doctor_id	int		`json:"doctor_id"`

}

func (m *PatientModel) Insert(patient *Patient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO patients (email, name, address ) VALUES ($1, $2, $3) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, patient.Email, patient.Name, patient.Address).Scan(&patient.Id)
}

// GetPatients returns all patients
//
//	@Summary		Returns all patients
//	@Description	Returns all patients
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]database.Patient
//	@Router			/api/v1/patients [get]
func (m *PatientModel) GetAll() ([]*Patient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM patients"

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	patients := []*Patient{}

	for rows.Next() {
		var patient Patient

		err := rows.Scan(&patient.Id, &patient.Name, &patient.Email, &patient.Address, &patient.Bloodtype)
		if err != nil {
			return nil, err
		}

		patients = append(patients, &patient)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil

}

// GetPatient returns a single patient
//
//	@Summary		Returns a single patient
//	@Description	Returns a single patient
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Patient ID"
//	@Success		200	{object}	database.Patient
//	@Router			/api/v1/patients/{id} [get]
func (m *PatientModel) Get(id int) (*Patient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM patients WHERE id = $1"

	var patient Patient

	m.DB.QueryRowContext(ctx, query, id).Scan(&patient.Id, &patient.Email, &patient.Name, &patient.Address, &patient.Bloodtype,  &patient.Doctor_id);
	log.Println(`Patient Email: $1`, patient )
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	return &patient, nil
}

// UpdatePatient updates an existing patient
//
//	@Summary		Updates an existing patient
//	@Description	Updates an existing patient
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Patient ID"
//	@Param			patient	body		database.Patient	true	"Patient"
//	@Success		200	{object}	database.Patient
//	@Router			/api/v1/patients/{id} [put]
//	@Security		BearerAuth
func (m *PatientModel) Update(patient *Patient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE patients SET name = $1, email = $2, address = $3, bloodtype = $4, doctor_id = $5 WHERE id = $6"

	_, err := m.DB.ExecContext(ctx, query, patient.Name, patient.Email, patient.Address, patient.Bloodtype, patient.Doctor_id, patient.Id)
	if err != nil {
		return err
	}

	return nil
}

// DeletePatient deletes an existing patient
//
//	@Summary		Deletes an existing patient
//	@Description	Deletes an existing patient
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Patient ID"
//	@Success		204
//	@Router			/api/v1/patients/{id} [delete]
//	@Security		BearerAuth
func (m *PatientModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM patients WHERE id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

