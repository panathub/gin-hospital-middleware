package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-hospital-middleware/internal/auth"
	"github.com/gin-hospital-middleware/internal/handlers"
	appmiddleware "github.com/gin-hospital-middleware/internal/middleware"
)

func Setup(
	tokenService *auth.TokenService,
	staffHandler *handlers.StaffHandler,
	patientHandler *handlers.PatientHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	staff := r.Group("/staff")
	{
		staff.POST("/create", staffHandler.Create)
		staff.POST("/login", staffHandler.Login)
	}

	patient := r.Group("/patient")
	patient.Use(appmiddleware.AuthRequired(tokenService))
	{
		patient.POST("/search", patientHandler.Search)
	}

	return r
}
