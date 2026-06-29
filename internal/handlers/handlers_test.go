package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-hospital-middleware/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockHospitalServer() *httptest.Server {
	patients := map[string]map[string]string{
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

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/patient/search/")
		patient, ok := patients[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "patient not found"})
			return
		}
		_ = json.NewEncoder(w).Encode(patient)
	}))
}

func createStaff(t *testing.T, router http.Handler, username, password, hospital string) {
	t.Helper()
	body := map[string]string{
		"username": username,
		"password": password,
		"hospital": hospital,
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/create", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func loginStaff(t *testing.T, router http.Handler, username, password, hospital string) string {
	t.Helper()
	body := map[string]string{
		"username": username,
		"password": password,
		"hospital": hospital,
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	return resp["token"]
}

func TestStaffCreate_Positive(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)

	body := map[string]string{
		"username": "nurse01",
		"password": "secret12",
		"hospital": "hospital-a",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/create", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "nurse01", resp["username"])
	assert.Equal(t, "hospital-a", resp["hospital"])
}

func TestStaffCreate_Negative_Duplicate(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")

	body := map[string]string{
		"username": "nurse01",
		"password": "secret12",
		"hospital": "hospital-a",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/create", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestStaffCreate_Negative_MissingFields(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)

	body := map[string]string{"username": "nurse01"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/create", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStaffLogin_Positive(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")

	body := map[string]string{
		"username": "nurse01",
		"password": "secret12",
		"hospital": "hospital-a",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["token"])
	assert.Equal(t, "nurse01", resp["username"])
}

func TestStaffLogin_Negative_InvalidPassword(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")

	body := map[string]string{
		"username": "nurse01",
		"password": "wrongpass",
		"hospital": "hospital-a",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestStaffLogin_Negative_UnknownUser(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)

	body := map[string]string{
		"username": "unknown",
		"password": "secret12",
		"hospital": "hospital-a",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientSearch_Positive_ByNationalID(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")
	token := loginStaff(t, router, "nurse01", "secret12", "hospital-a")

	body := map[string]string{"national_id": "1234567890123"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Patients []map[string]any `json:"patients"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Patients, 1)
	assert.Equal(t, "1234567890123", resp.Patients[0]["national_id"])
	assert.Equal(t, "Somchai", resp.Patients[0]["first_name_en"])
}

func TestPatientSearch_Positive_ByFirstName(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, db := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")
	token := loginStaff(t, router, "nurse01", "secret12", "hospital-a")

	require.NoError(t, db.Exec(`INSERT INTO patients (
		id, hospital, first_name_th, middle_name_th, last_name_th,
		first_name_en, middle_name_en, last_name_en, date_of_birth,
		patient_hn, national_id, passport_id, phone_number, email, gender,
		created_at, updated_at
	) VALUES (
		'550e8400-e29b-41d4-a716-446655440000', 'hospital-a', '', '', '',
		'John', '', 'Doe', '1992-03-10',
		'HN-A-100', '9999999999999', '', '0800000000', 'john@example.com', 'M',
		datetime('now'), datetime('now')
	)`).Error)

	body := map[string]string{"first_name": "john"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Patients []map[string]any `json:"patients"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Patients, 1)
	assert.Equal(t, "John", resp.Patients[0]["first_name_en"])
}

func TestPatientSearch_Positive_EmptyBody(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, db := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")
	token := loginStaff(t, router, "nurse01", "secret12", "hospital-a")

	require.NoError(t, db.Exec(`INSERT INTO patients (
		id, hospital, first_name_th, middle_name_th, last_name_th,
		first_name_en, middle_name_en, last_name_en, date_of_birth,
		patient_hn, national_id, passport_id, phone_number, email, gender,
		created_at, updated_at
	) VALUES (
		'550e8400-e29b-41d4-a716-446655440002', 'hospital-a', '', '', '',
		'Alice', '', 'Brown', '1991-01-01',
		'HN-A-200', '2222222222222', '', '0800000002', 'alice@example.com', 'F',
		datetime('now'), datetime('now')
	)`).Error)

	req := httptest.NewRequest(http.MethodPost, "/patient/search", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Patients []map[string]any `json:"patients"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Patients, 1)
}

func TestPatientSearch_Negative_Unauthorized(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)

	body := map[string]string{"national_id": "1234567890123"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientSearch_Negative_InvalidToken(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)

	body := map[string]string{"national_id": "1234567890123"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientSearch_Negative_DifferentHospitalIsolation(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, db := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")
	token := loginStaff(t, router, "nurse01", "secret12", "hospital-a")

	require.NoError(t, db.Exec(`INSERT INTO patients (
		id, hospital, first_name_th, middle_name_th, last_name_th,
		first_name_en, middle_name_en, last_name_en, date_of_birth,
		patient_hn, national_id, passport_id, phone_number, email, gender,
		created_at, updated_at
	) VALUES (
		'550e8400-e29b-41d4-a716-446655440001', 'hospital-b', '', '', '',
		'Jane', '', 'Smith', '1990-01-01',
		'HN-B-001', '1111111111111', '', '0800000001', 'jane@example.com', 'F',
		datetime('now'), datetime('now')
	)`).Error)

	body := map[string]string{"first_name": "jane"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Patients []map[string]any `json:"patients"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Empty(t, resp.Patients)
}

func TestStaffCreate_Negative_ShortPassword(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)

	body := map[string]string{
		"username": "nurse01",
		"password": "short",
		"hospital": "hospital-a",
	}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/create", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStaffLogin_Negative_MissingFields(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)

	body := map[string]string{"username": "nurse01"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPatientSearch_Positive_ByPassportID(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")
	token := loginStaff(t, router, "nurse01", "secret12", "hospital-a")

	body := map[string]string{"passport_id": "AB1234567"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Patients []map[string]any `json:"patients"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp.Patients, 1)
	assert.Equal(t, "AB1234567", resp.Patients[0]["passport_id"])
	assert.Equal(t, "Somying", resp.Patients[0]["first_name_en"])
}

func TestPatientSearch_Negative_MalformedAuthHeader(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)

	body := map[string]string{"national_id": "1234567890123"}
	payload, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token not-bearer-format")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientSearch_Negative_InvalidJSON(t *testing.T) {
	mock := mockHospitalServer()
	defer mock.Close()

	router, _ := testutil.SetupTestRouter(t, mock.URL)
	createStaff(t, router, "nurse01", "secret12", "hospital-a")
	token := loginStaff(t, router, "nurse01", "secret12", "hospital-a")

	req := httptest.NewRequest(http.MethodPost, "/patient/search", strings.NewReader("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
