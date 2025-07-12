package main

import (
	"net/http"
	"rest-api-in-gin/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)


func (app *application) createPatient(c *gin.Context) {
	var patient database.Patient

	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staff := app.GetUserFromContext(c)
	
	if staff.Role != `Receptionist` {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to register patients"})
		return
	}

	err := app.models.Patients.Insert(&patient)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

func (app *application) getAllPatients(c *gin.Context) {
	patients, err := app.models.Patients.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
	}

	c.JSON(http.StatusOK, patients)
}

func (app *application) getPatient(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
	}

	patient, err := app.models.Patients.Get(id)
	if patient == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Patient"})
	}

	c.JSON(http.StatusOK, patient)
}

func (app *application) updatePatient(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
	}

	existingPatient, err := app.models.Patients.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingPatient == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}


	updatedPatient := &database.Patient{}
	if err := c.ShouldBindJSON(updatedPatient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPatient.Id = id
	if err := app.models.Patients.Update(updatedPatient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	c.JSON(http.StatusOK, updatedPatient)
}

func (app *application) deletePatient(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid patient ID"})
	}

	if err := app.models.Patients.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete event"})
	}

	c.JSON(http.StatusNoContent, nil )
}

func (app *application) getAttendeesForEvent(c *gin.Context) {
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event id"})
	// 	return
	// }

	// users, err := app.models.Attendees.GetAttendeesByEvent(id)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to to retreive attendees for events"})
	// 	return
	// }

	// c.JSON(http.StatusOK, users)
}

func (app *application) addDoctorToPatient(c *gin.Context) {
	patientId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event Id"})
		return
	}

	staffId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user Id"})
		return
	}

	patient, err := app.models.Patients.Get(patientId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive event"})
		return
	}
	if patient == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
	}

	staffToAdd, err := app.models.Staffs.Get(staffId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive user"})
		return
	}

	if staffToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	// user := app.GetUserFromContext(c)

	// if patient.Doctor_id != user.Id {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add an attendee"})
	// 	return
	// }

	// existingStaff, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userToAdd.Id)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive attendee"})
	// 	return
	// }
	// if existingStaff != nil {
	// 	c.JSON(http.StatusConflict, gin.H{"error": "Attendee already exists"})
	// 	return
	// }

	// attendee := database.Attendee{
	// 	EventId: event.Id,
	// 	UserId:  userToAdd.Id,
	// }

	// _, err = app.models.Attendees.Insert(&attendee)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add  attendee"})
	// 	return
	// }

	// c.JSON(http.StatusCreated, attendee)

}
