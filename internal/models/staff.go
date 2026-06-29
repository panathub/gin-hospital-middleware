package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Staff struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username     string    `gorm:"size:100;not null;uniqueIndex:idx_staff_username_hospital" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Hospital     string    `gorm:"size:100;not null;uniqueIndex:idx_staff_username_hospital" json:"hospital"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (s *Staff) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
