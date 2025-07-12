package database

import "database/sql"

type Models struct {
	Staffs StaffModel
	Patients PatientModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Staffs: StaffModel{DB: db},
		Patients: PatientModel{DB: db},
	}
}