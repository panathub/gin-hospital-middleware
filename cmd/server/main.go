package main

import (
	"log"

	"github.com/gin-hospital-middleware/internal/auth"
	"github.com/gin-hospital-middleware/internal/config"
	"github.com/gin-hospital-middleware/internal/database"
	"github.com/gin-hospital-middleware/internal/handlers"
	"github.com/gin-hospital-middleware/internal/hospital"
	"github.com/gin-hospital-middleware/internal/repository"
	"github.com/gin-hospital-middleware/internal/router"
	"github.com/gin-hospital-middleware/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	staffRepo := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	tokenService := auth.NewTokenService(cfg.JWTSecret, cfg.JWTExpiry)
	hospitalClient := hospital.NewClient(cfg.HospitalABaseURL)

	staffService := service.NewStaffService(staffRepo, tokenService)
	patientService := service.NewPatientService(patientRepo, hospitalClient)

	staffHandler := handlers.NewStaffHandler(staffService)
	patientHandler := handlers.NewPatientHandler(patientService)

	engine := router.Setup(tokenService, staffHandler, patientHandler)
	log.Printf("server listening on :%s", cfg.Port)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
