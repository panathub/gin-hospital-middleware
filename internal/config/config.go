package config

import (
	"cmp"
	"os"
	"time"
)

type Config struct {
	Port              string
	DatabaseURL       string
	JWTSecret         string
	JWTExpiry         time.Duration
	HospitalABaseURL  string
}

func Load() Config {
	return Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/hospital_middleware?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiry:        24 * time.Hour,
		HospitalABaseURL: getEnv("HOSPITAL_A_BASE_URL", "https://hospital-a.api.co.th"),
	}
}

func getEnv(key, fallback string) string {
	return cmp.Or(os.Getenv(key), fallback)
}
