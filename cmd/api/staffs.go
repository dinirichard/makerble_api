package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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