package main

import (
	"log"
	"net/http"
	"rest-api-in-gin/internal/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type registerPatientRequest struct {
	Email 	string 	`json:"email" binding:"required,email"`
	Name 	string 	`json:"name" binding:"required,min=3"`
	Address string 	`json:"address" binding:"required,min=5"`
}

type registerStaffRequest struct {
	Email 		string 	`json:"email" binding:"required,email"`
	Password	string 	`json:"password" binding:"required,min=5"`
	Name 		string 	`json:"name" binding:"required,min=3"`
	Role 		string 	`json:"role" binding:"required"`
}

type loginRequest struct {
	Email 		string 	`json:"email" binding:"required,email"`
	Password	string 	`json:"password" binding:"required,min=5"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (app *application) login(c *gin.Context) {
	var auth loginRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingStaff, err := app.models.Staffs.GetByEmail(auth.Email)
	if existingStaff == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	log.Println(`Existing Staff: $1 `, existingStaff)

	err = bcrypt.CompareHashAndPassword([]byte(existingStaff.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"staffId": existingStaff.Id,
		"expr": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating the token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{Token: tokenString})

}

func (app *application) registerPatient(c *gin.Context) {
	var register registerPatientRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staff := app.GetUserFromContext(c)

	log.Println(`Staff Context: $1`, staff)
	
	if staff.Role != `Receptionist` {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to register patients"})
		return
	}

	patient := database.Patient{
		Email: register.Email,
		Name: register.Name,
		Address: register.Address,
	}

	err := app.models.Patients.Insert(&patient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Could not create patient"})
		return
	}

	c.JSON(http.StatusCreated, patient)

}

func (app *application) registerStaff(c *gin.Context) {
	var register registerStaffRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Something went wrong"})
		return
	}
	register.Password = string(hashedPassword)

	staff := database.Staff{
		Email: register.Email,
		Password: register.Password,
		Name: register.Name,
		Role: register.Role,
	}

	err = app.models.Staffs.Insert(&staff)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Could not create staff"})
		return
	}

	c.JSON(http.StatusCreated, staff)

}