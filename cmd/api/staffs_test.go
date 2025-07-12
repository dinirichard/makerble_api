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

func TestRegisterStaffHandlerIntegration(t *testing.T) {
	// Reset database state
	_, err := testApp.DB.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/staff/register", testApp.registerStaff)

	// Staff data
	staffName := "Dr. Alice"
	staffEmail := "alice@example.com"
	staffPassword := "password123"
	staffRole := "Doctor"
	staffPayload := []byte(`{"name":"` + staffName + `","email":"` + staffEmail + `","password":"` + staffPassword + `","role":"` + staffRole + `"}`)

	// Create request
	req, err := http.NewRequest("POST", "/auth/staff/register", bytes.NewBuffer(staffPayload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Assert HTTP response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Assert response body
	var createdStaff database.Staff
	err = json.Unmarshal(rr.Body.Bytes(), &createdStaff)
	require.NoError(t, err)
	assert.Equal(t, staffName, createdStaff.Name)
	assert.Equal(t, staffEmail, createdStaff.Email)
	assert.Equal(t, staffRole, createdStaff.Role)
	assert.NotZero(t, createdStaff.Id)

	// Assert database state
	var dbStaff database.Staff
	err = testApp.DB.QueryRow("SELECT id, name, email, role FROM staffs WHERE id = $1", createdStaff.Id).Scan(&dbStaff.Id, &dbStaff.Name, &dbStaff.Email, &dbStaff.Role)
	require.NoError(t, err)
	assert.Equal(t, createdStaff.Id, dbStaff.Id)
	assert.Equal(t, staffName, dbStaff.Name)
}

func TestGetStaffHandlerIntegration(t *testing.T) {
	// Reset database state and insert a test staff member
	_, err := testApp.DB.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	staffToInsert := database.Staff{
		Name:     "Dr. Bob",
		Email:    "bob@example.com",
		Password: "a-hashed-password", // In a real scenario, this would be properly hashed
		Role:     "Doctor",
	}
	err = testApp.models.Staffs.Insert(&staffToInsert)
	require.NoError(t, err)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/staffs/:id", testApp.getStaff)

	// Create request
	req, err := http.NewRequest("GET", "/staffs/"+strconv.Itoa(staffToInsert.Id), nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Assert HTTP response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert response body
	var fetchedStaff database.Staff
	err = json.Unmarshal(rr.Body.Bytes(), &fetchedStaff)
	require.NoError(t, err)
	assert.Equal(t, staffToInsert.Id, fetchedStaff.Id)
	assert.Equal(t, staffToInsert.Name, fetchedStaff.Name)
	assert.Equal(t, staffToInsert.Email, fetchedStaff.Email)
}

func TestDeleteStaffHandlerIntegration(t *testing.T) {
	// Reset and seed the database
	_, err := testApp.DB.Exec("TRUNCATE TABLE patients, staffs RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	staffToDelete := database.Staff{Name: "To Be Deleted", Email: "delete.staff@example.com", Password: "password", Role: "Temp"}
	err = testApp.models.Staffs.Insert(&staffToDelete)
	require.NoError(t, err)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/staffs/:id", testApp.deleteStaff)

	// Create request
	req, err := http.NewRequest("DELETE", "/staffs/"+strconv.Itoa(staffToDelete.Id), nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Assert HTTP response
	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Assert database state
	var count int
	err = testApp.DB.QueryRow("SELECT COUNT(id) FROM staffs WHERE id = $1", staffToDelete.Id).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Staff member should have been deleted")
}