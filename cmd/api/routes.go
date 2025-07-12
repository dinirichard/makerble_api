package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1") 
	{
		v1.GET("/patients", app.getAllPatients)
		

		v1.GET("/staffs/:id", app.getStaff)
		v1.DELETE("/staffs/:id", app.deleteStaff)

		v1.POST("/auth/staff/register", app.registerStaff)
		v1.POST("/auth/login", app.login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/auth/register", app.registerPatient)

		authGroup.POST("/patients", app.createPatient)
		
		authGroup.GET("/patients/:id", app.getPatient)
		authGroup.PUT("/patients/:id", app.updatePatient)
		authGroup.DELETE("/patients/:id", app.deletePatient)
	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(302, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	return g
}