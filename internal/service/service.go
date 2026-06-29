package service

import (
	"errors"

	"github.com/gin-hospital-middleware/internal/auth"
	"github.com/gin-hospital-middleware/internal/hospital"
	"github.com/gin-hospital-middleware/internal/models"
	"github.com/gin-hospital-middleware/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrStaffExists       = errors.New("staff already exists")
	ErrInvalidCredential = errors.New("invalid username or password")
)

type StaffService struct {
	staffRepo    *repository.StaffRepository
	tokenService *auth.TokenService
}

func NewStaffService(staffRepo *repository.StaffRepository, tokenService *auth.TokenService) *StaffService {
	return &StaffService{staffRepo: staffRepo, tokenService: tokenService}
}

func (s *StaffService) Create(username, password, hospital string) (*models.Staff, error) {
	_, err := s.staffRepo.FindByUsernameAndHospital(username, hospital)
	if err == nil {
		return nil, ErrStaffExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	staff := &models.Staff{
		Username:     username,
		PasswordHash: string(hash),
		Hospital:     hospital,
	}
	if err := s.staffRepo.Create(staff); err != nil {
		return nil, err
	}
	return staff, nil
}

func (s *StaffService) Login(username, password, hospital string) (string, *models.Staff, error) {
	staff, err := s.staffRepo.FindByUsernameAndHospital(username, hospital)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrInvalidCredential
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredential
	}

	token, err := s.tokenService.Generate(staff.ID, staff.Username, staff.Hospital)
	if err != nil {
		return "", nil, err
	}
	return token, staff, nil
}

type PatientService struct {
	patientRepo    *repository.PatientRepository
	hospitalClient *hospital.Client
}

func NewPatientService(patientRepo *repository.PatientRepository, hospitalClient *hospital.Client) *PatientService {
	return &PatientService{patientRepo: patientRepo, hospitalClient: hospitalClient}
}

func (s *PatientService) Search(hospital string, filter repository.PatientSearchFilter) ([]models.Patient, error) {
	if filter.NationalID != "" {
		if patient, err := s.hospitalClient.SearchByID(hospital, filter.NationalID); err != nil {
			return nil, err
		} else if patient != nil {
			if _, err := s.patientRepo.Upsert(patient); err != nil {
				return nil, err
			}
		}
	}
	if filter.PassportID != "" {
		if patient, err := s.hospitalClient.SearchByID(hospital, filter.PassportID); err != nil {
			return nil, err
		} else if patient != nil {
			if _, err := s.patientRepo.Upsert(patient); err != nil {
				return nil, err
			}
		}
	}

	return s.patientRepo.Search(hospital, filter)
}
