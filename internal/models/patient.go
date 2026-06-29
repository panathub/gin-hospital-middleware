package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Hospital     string    `gorm:"size:100;not null;index" json:"hospital"`
	FirstNameTH  string    `gorm:"size:100" json:"first_name_th"`
	MiddleNameTH string    `gorm:"size:100" json:"middle_name_th"`
	LastNameTH   string    `gorm:"size:100" json:"last_name_th"`
	FirstNameEN  string    `gorm:"size:100" json:"first_name_en"`
	MiddleNameEN string    `gorm:"size:100" json:"middle_name_en"`
	LastNameEN   string    `gorm:"size:100" json:"last_name_en"`
	DateOfBirth  string    `gorm:"size:20" json:"date_of_birth"`
	PatientHN    string    `gorm:"size:50;index" json:"patient_hn"`
	NationalID   string    `gorm:"size:20;index" json:"national_id"`
	PassportID   string    `gorm:"size:20;index" json:"passport_id"`
	PhoneNumber  string    `gorm:"size:30" json:"phone_number"`
	Email        string    `gorm:"size:255" json:"email"`
	Gender       string    `gorm:"size:1" json:"gender"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *Patient) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
