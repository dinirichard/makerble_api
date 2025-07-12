package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"rest-api-in-gin/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a router and inject a mock user for auth-required routes
func getRouterWithMockUser() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		// Mock a user and add it to the context with the key "user"
		// This simulates a user who has passed through the auth middleware
		mockStaff := &database.Staff{
			Id:   1,
			Role: "Receptionist", // Role required to create patients
		}
		c.Set("user", mockStaff)
		c.Next()
	})
	return router
}

func TestCreatePatientHandlerIntegration(t *testing.T) {
	// Ensure the testApp is initialized
	require.NotNil(t, testApp.models.Patients, "testApp is not initialized")

	// Reset database state before the test
	_, err := testApp.DB.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	// Setup router with a mock user having the 'Receptionist' role
	router := getRouterWithMockUser()
	router.POST("/patients", testApp.createPatient)

	// Patient data
	patientName := "Jane Doe"
	patientEmail := "jane.doe@example.com"
	patientAddress := "456 Oak Ave"
	patientPayload := []byte(`{"name":"` + patientName + `","email":"` + patientEmail + `","address":"` + patientAddress + `"}`)

	// Create request
	req, err := http.NewRequest("POST", "/patients", bytes.NewBuffer(patientPayload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Assert HTTP response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Assert response body
	var createdPatient database.Patient
	err = json.Unmarshal(rr.Body.Bytes(), &createdPatient)
	require.NoError(t, err)
	assert.Equal(t, patientName, createdPatient.Name)
	assert.Equal(t, patientEmail, createdPatient.Email)
	assert.NotZero(t, createdPatient.Id) // Check that an ID was assigned

	// Assert database state
	var dbPatient database.Patient
	err = testApp.DB.QueryRow("SELECT id, name, email, address FROM patients WHERE id = $1", createdPatient.Id).Scan(&dbPatient.Id, &dbPatient.Name, &dbPatient.Email, &dbPatient.Address)
	require.NoError(t, err)
	assert.Equal(t, createdPatient.Id, dbPatient.Id)
	assert.Equal(t, patientName, dbPatient.Name)
	assert.Equal(t, patientEmail, dbPatient.Email)
}

func TestGetPatientHandlerIntegration(t *testing.T) {
	// Ensure the testApp is initialized
	require.NotNil(t, testApp.models.Patients, "testApp is not initialized")

	// Reset database state and insert a test patient
	_, err := testApp.DB.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	patientToInsert := database.Patient{
		Name:    "John Smith",
		Email:   "john.smith@example.com",
		Address: "789 Pine St",
	}
	err = testApp.models.Patients.Insert(&patientToInsert)
	require.NoError(t, err)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/patients/:id", testApp.getPatient)

	// Create request
	req, err := http.NewRequest("GET", "/patients/"+strconv.Itoa(patientToInsert.Id), nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Assert HTTP response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert response body
	var fetchedPatient database.Patient
	err = json.Unmarshal(rr.Body.Bytes(), &fetchedPatient)
	require.NoError(t, err)
	assert.Equal(t, patientToInsert.Id, fetchedPatient.Id)
	assert.Equal(t, patientToInsert.Name, fetchedPatient.Name)
	assert.Equal(t, patientToInsert.Email, fetchedPatient.Email)
}

func TestUpdatePatientHandlerIntegration(t *testing.T) {
	// Reset and seed the database
	_, err := testApp.DB.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	patientToUpdate := database.Patient{Name: "Old Name", Email: "update@example.com", Address: "Old Address"}
	err = testApp.models.Patients.Insert(&patientToUpdate)
	require.NoError(t, err)

	// Setup router with mock user
	router := getRouterWithMockUser()
	router.PUT("/patients/:id", testApp.updatePatient)

	// New patient data
	updatedName := "New Name"
	updatedAddress := "New Address"
	// The email must be included to pass validation, even if it's the same as the old one.
	updatePayload := []byte(`{"name":"` + updatedName + `","email":"` + patientToUpdate.Email + `","address":"` + updatedAddress + `"}`)

	// Create request
	req, err := http.NewRequest("PUT", "/patients/"+strconv.Itoa(patientToUpdate.Id), bytes.NewBuffer(updatePayload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Assert HTTP response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert database state
	var dbPatient database.Patient
	err = testApp.DB.QueryRow("SELECT id, name, address FROM patients WHERE id = $1", patientToUpdate.Id).Scan(&dbPatient.Id, &dbPatient.Name, &dbPatient.Address)
	require.NoError(t, err)
	assert.Equal(t, updatedName, dbPatient.Name)
	assert.Equal(t, updatedAddress, dbPatient.Address)
}

func TestDeletePatientHandlerIntegration(t *testing.T) {
	// Reset and seed the database
	_, err := testApp.DB.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	patientToDelete := database.Patient{Name: "To Be Deleted", Email: "delete@example.com", Address: "Some Address"}
	err = testApp.models.Patients.Insert(&patientToDelete)
	require.NoError(t, err)

	// Setup router with mock user
	router := getRouterWithMockUser()
	router.DELETE("/patients/:id", testApp.deletePatient)

	// Create request
	req, err := http.NewRequest("DELETE", "/patients/"+strconv.Itoa(patientToDelete.Id), nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Assert HTTP response
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Assert database state
	var count int
	err = testApp.DB.QueryRow("SELECT COUNT(id) FROM patients WHERE id = $1", patientToDelete.Id).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Patient should have been deleted from the database")
}