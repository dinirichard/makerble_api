package main

import (
	"rest-api-in-gin/internal/database"

	"github.com/gin-gonic/gin"
)

func (app *application) GetUserFromContext(c *gin.Context) *database.Staff {
	contextStaff, exist := c.Get("user")
	if !exist {
		return &database.Staff{}
	}

	staff, ok := contextStaff.(*database.Staff)
	if !ok {
		return &database.Staff{}
	}

	return staff
}