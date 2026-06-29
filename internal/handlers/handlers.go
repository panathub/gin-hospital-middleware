package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-hospital-middleware/internal/auth"
	appmiddleware "github.com/gin-hospital-middleware/internal/middleware"
	"github.com/gin-hospital-middleware/internal/repository"
	"github.com/gin-hospital-middleware/internal/service"
)

type StaffHandler struct {
	staffService *service.StaffService
}

func NewStaffHandler(staffService *service.StaffService) *StaffHandler {
	return &StaffHandler{staffService: staffService}
}

type createStaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Hospital string `json:"hospital" binding:"required"`
}

type loginStaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

func (h *StaffHandler) Create(c *gin.Context) {
	var req createStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	staff, err := h.staffService.Create(req.Username, req.Password, req.Hospital)
	if err != nil {
		if errors.Is(err, service.ErrStaffExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "staff already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create staff"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       staff.ID,
		"username": staff.Username,
		"hospital": staff.Hospital,
	})
}

func (h *StaffHandler) Login(c *gin.Context) {
	var req loginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, staff, err := h.staffService.Login(req.Username, req.Password, req.Hospital)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredential) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": staff.Username,
		"hospital": staff.Hospital,
	})
}

type PatientHandler struct {
	patientService *service.PatientService
}

func NewPatientHandler(patientService *service.PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

type patientSearchRequest struct {
	NationalID  string `json:"national_id"`
	PassportID  string `json:"passport_id"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

func (h *PatientHandler) Search(c *gin.Context) {
	claimsValue, exists := c.Get(appmiddleware.StaffContextKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	claims := claimsValue.(*auth.Claims)

	var req patientSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := repository.PatientSearchFilter{
		NationalID:  req.NationalID,
		PassportID:  req.PassportID,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		DateOfBirth: req.DateOfBirth,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}

	patients, err := h.patientService.Search(claims.Hospital, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search patients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}
