package testutil

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-hospital-middleware/internal/auth"
	"github.com/gin-hospital-middleware/internal/config"
	"github.com/gin-hospital-middleware/internal/handlers"
	"github.com/gin-hospital-middleware/internal/hospital"
	"github.com/gin-hospital-middleware/internal/models"
	"github.com/gin-hospital-middleware/internal/repository"
	"github.com/gin-hospital-middleware/internal/router"
	"github.com/gin-hospital-middleware/internal/service"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestRouter(t *testing.T, hospitalBaseURL string) (*gin.Engine, *gorm.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dsn := fmt.Sprintf("file:test_%s?mode=memory&cache=shared", uuid.New().String())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Staff{}, &models.Patient{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	staffRepo := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	tokenService := auth.NewTokenService("test-secret", config.Load().JWTExpiry)
	hospitalClient := hospital.NewClient(hospitalBaseURL)

	staffService := service.NewStaffService(staffRepo, tokenService)
	patientService := service.NewPatientService(patientRepo, hospitalClient)

	staffHandler := handlers.NewStaffHandler(staffService)
	patientHandler := handlers.NewPatientHandler(patientService)

	return router.Setup(tokenService, staffHandler, patientHandler), db
}
