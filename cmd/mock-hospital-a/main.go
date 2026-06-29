package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var mockPatients = map[string]gin.H{
	"1234567890123": {
		"first_name_th":  "สมชาย",
		"middle_name_th": "",
		"last_name_th":   "ใจดี",
		"first_name_en":  "Somchai",
		"middle_name_en": "",
		"last_name_en":   "Jaidee",
		"date_of_birth":  "1990-01-15",
		"patient_hn":     "HN-A-001",
		"national_id":    "1234567890123",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "somchai@example.com",
		"gender":         "M",
	},
	"AB1234567": {
		"first_name_th":  "สมหญิง",
		"middle_name_th": "มณี",
		"last_name_th":   "รักษ์ดี",
		"first_name_en":  "Somying",
		"middle_name_en": "Manee",
		"last_name_en":   "Rakdee",
		"date_of_birth":  "1985-06-20",
		"patient_hn":     "HN-A-002",
		"national_id":    "",
		"passport_id":    "AB1234567",
		"phone_number":   "0898765432",
		"email":          "somying@example.com",
		"gender":         "F",
	},
}

func main() {
	r := gin.Default()

	r.GET("/patient/search/:id", func(c *gin.Context) {
		id := c.Param("id")
		patient, ok := mockPatients[id]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
			return
		}
		c.JSON(http.StatusOK, patient)
	})

	log.Println("mock hospital-a API listening on :9000")
	if err := r.Run(":9000"); err != nil {
		log.Fatalf("mock server failed: %v", err)
	}
}
