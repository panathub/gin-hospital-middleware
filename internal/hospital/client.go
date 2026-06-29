package hospital

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-hospital-middleware/internal/models"
)

type PatientResponse struct {
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"`
	PatientHN    string `json:"patient_hn"`
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
}

type Client struct {
	httpClient *http.Client
	baseURLs   map[string]string
}

func NewClient(hospitalABaseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURLs: map[string]string{
			"hospital-a": hospitalABaseURL,
		},
	}
}

func (c *Client) SearchByID(hospital, id string) (*models.Patient, error) {
	baseURL, ok := c.baseURLs[hospital]
	if !ok {
		return nil, fmt.Errorf("unsupported hospital: %s", hospital)
	}

	endpoint, err := url.JoinPath(baseURL, "patient", "search", id)
	if err != nil {
		return nil, fmt.Errorf("build hospital API URL: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call hospital API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hospital API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp PatientResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &models.Patient{
		Hospital:     hospital,
		FirstNameTH:  apiResp.FirstNameTH,
		MiddleNameTH: apiResp.MiddleNameTH,
		LastNameTH:   apiResp.LastNameTH,
		FirstNameEN:  apiResp.FirstNameEN,
		MiddleNameEN: apiResp.MiddleNameEN,
		LastNameEN:   apiResp.LastNameEN,
		DateOfBirth:  apiResp.DateOfBirth,
		PatientHN:    apiResp.PatientHN,
		NationalID:   apiResp.NationalID,
		PassportID:   apiResp.PassportID,
		PhoneNumber:  apiResp.PhoneNumber,
		Email:        apiResp.Email,
		Gender:       apiResp.Gender,
	}, nil
}
