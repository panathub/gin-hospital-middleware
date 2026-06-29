package repository

import (
	"errors"
	"strings"

	"github.com/gin-hospital-middleware/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StaffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

func (r *StaffRepository) Create(staff *models.Staff) error {
	return r.db.Create(staff).Error
}

func (r *StaffRepository) FindByUsernameAndHospital(username, hospital string) (*models.Staff, error) {
	var staff models.Staff
	err := r.db.Where("username = ? AND hospital = ?", username, hospital).First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

type PatientSearchFilter struct {
	NationalID  string
	PassportID  string
	FirstName   string
	MiddleName  string
	LastName    string
	DateOfBirth string
	PhoneNumber string
	Email       string
}

type PatientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Upsert(patient *models.Patient) (*models.Patient, error) {
	var existing models.Patient
	query := r.db.Where("hospital = ?", patient.Hospital)
	switch {
	case patient.NationalID != "":
		query = query.Where("national_id = ?", patient.NationalID)
	case patient.PassportID != "":
		query = query.Where("passport_id = ?", patient.PassportID)
	case patient.PatientHN != "":
		query = query.Where("patient_hn = ?", patient.PatientHN)
	default:
		return nil, gorm.ErrRecordNotFound
	}

	err := query.First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := r.db.Create(patient).Error; err != nil {
			return nil, err
		}
		return patient, nil
	}
	if err != nil {
		return nil, err
	}

	existing.FirstNameTH = patient.FirstNameTH
	existing.MiddleNameTH = patient.MiddleNameTH
	existing.LastNameTH = patient.LastNameTH
	existing.FirstNameEN = patient.FirstNameEN
	existing.MiddleNameEN = patient.MiddleNameEN
	existing.LastNameEN = patient.LastNameEN
	existing.DateOfBirth = patient.DateOfBirth
	existing.PatientHN = patient.PatientHN
	existing.NationalID = patient.NationalID
	existing.PassportID = patient.PassportID
	existing.PhoneNumber = patient.PhoneNumber
	existing.Email = patient.Email
	existing.Gender = patient.Gender

	if err := r.db.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *PatientRepository) Search(hospital string, filter PatientSearchFilter) ([]models.Patient, error) {
	query := r.db.Where("hospital = ?", hospital)

	if filter.NationalID != "" {
		query = query.Where("national_id = ?", filter.NationalID)
	}
	if filter.PassportID != "" {
		query = query.Where("passport_id = ?", filter.PassportID)
	}
	if filter.DateOfBirth != "" {
		query = query.Where("date_of_birth = ?", filter.DateOfBirth)
	}
	if filter.PhoneNumber != "" {
		query = query.Where("phone_number = ?", filter.PhoneNumber)
	}
	if filter.Email != "" {
		query = query.Where("LOWER(email) = ?", strings.ToLower(filter.Email))
	}
	if filter.FirstName != "" {
		like := "%" + strings.ToLower(filter.FirstName) + "%"
		query = query.Where(
			"LOWER(first_name_th) LIKE ? OR LOWER(first_name_en) LIKE ?",
			like, like,
		)
	}
	if filter.MiddleName != "" {
		like := "%" + strings.ToLower(filter.MiddleName) + "%"
		query = query.Where(
			"LOWER(middle_name_th) LIKE ? OR LOWER(middle_name_en) LIKE ?",
			like, like,
		)
	}
	if filter.LastName != "" {
		like := "%" + strings.ToLower(filter.LastName) + "%"
		query = query.Where(
			"LOWER(last_name_th) LIKE ? OR LOWER(last_name_en) LIKE ?",
			like, like,
		)
	}

	var patients []models.Patient
	if err := query.Find(&patients).Error; err != nil {
		return nil, err
	}
	return patients, nil
}

func (r *PatientRepository) FindByID(id uuid.UUID) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.First(&patient, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}
