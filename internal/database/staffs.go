package database

import (
	"context"
	"database/sql"
	"time"
)

type StaffModel struct {
	DB *sql.DB
}

type Staff struct {
	Id 			int		`json:"id"`	
	Name 		string	`json:"name"`
	Password 	string	`json:"-"`
	Email 		string	`json:"email"`
	Role 		string	`json:"role"`
}

func (m *StaffModel) Insert(staff *Staff) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO staffs (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id"

	return m.DB.QueryRowContext(ctx, query, staff.Name, staff.Email, staff.Password, staff.Role).Scan(&staff.Id)
}

func (m *StaffModel) getStaff(query string, args ...interface{}) (*Staff, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var staff Staff

	m.DB.QueryRowContext(ctx, query, args...).Scan(&staff.Id, &staff.Name, &staff.Password, &staff.Email, &staff.Role)
	// err := m.DB.QueryRowContext(ctx, query, args...).Scan(&staff.Id, &staff.Name, &staff.Email, &staff.Role)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	return &staff, nil
}

func (m *StaffModel) Get(id int) (*Staff, error) {
	query := "SELECT id, name, password, email, role FROM staffs WHERE id = $1"
	return m.getStaff(query, id)
}

func (m *StaffModel) GetByEmail(email string) (*Staff, error) {
	query := "SELECT id, name, password, email, role FROM staffs WHERE email = $1"
	return m.getStaff(query, email)
}

func (m *StaffModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM staffs WHERE id = $1"

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

