package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// func (app *application) createStaff(c *gin.Context) {
// 	var staff database.Staff

// 	if err := c.ShouldBindJSON(&staff); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := app.models.Staffs.Insert(&staff)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, staff)
// }

// func (app *application) getAllStaffs(c *gin.Context) {
// 	staffs, err := app.models.Staffs.GetAll()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff ID"})
// 	}

// 	c.JSON(http.StatusOK, staffs)
// }

func (app *application) getStaff(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff ID"})
	}

	staff, err := app.models.Staffs.Get(id)
	if staff == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Staff"})
	}

	c.JSON(http.StatusOK, staff)
}

// func (app *application) updateStaff(c *gin.Context) { 
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid staff ID"})
// 	}

// 	existingStaff, err := app.models.Staffs.Get(id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
// 		return
// 	}

// 	if existingStaff == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
// 	}

// 	updatedStaff := &database.Staff{}
// 	if err := c.ShouldBindJSON(updatedStaff); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updatedStaff.Id = id
// 	if err := app.models.Staffs.Update(updatedStaff); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, updatedStaff)
// }

func (app *application) deleteStaff(c *gin.Context) { 
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid staff ID"})
	}

	if err := app.models.Staffs.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to delete event"})
	}

	c.JSON(http.StatusNoContent, nil )
}