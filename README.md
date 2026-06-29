# gin-hospital-middleware

Hospital Middleware API. The service lets hospital staff authenticate and search patient records from their own hospital, integrating with external Hospital Information Systems (HIS).

## Tech Stack

- Go 1.22
- Gin Framework
- PostgreSQL
- Docker & Docker Compose
- Nginx (reverse proxy)

## Project Structure

```
gin-hospital-middleware/
├── cmd/
│   ├── server/              # Main API server
│   └── mock-hospital-a/     # Mock Hospital A HIS API
├── internal/
│   ├── auth/                # JWT token service
│   ├── config/              # Environment configuration
│   ├── database/            # Database connection & migrations
│   ├── handlers/            # HTTP handlers
│   ├── hospital/            # External HIS API client
│   ├── middleware/          # Auth middleware
│   ├── models/              # Staff & Patient entities
│   ├── repository/          # Data access layer
│   ├── router/              # Route definitions
│   ├── service/             # Business logic
│   └── testutil/            # Test helpers
├── nginx/
│   └── nginx.conf           # Nginx reverse proxy config
├── docker-compose.yml
├── Dockerfile
└── Dockerfile.mock-hospital
```

## Entity-Relationship Diagram

```mermaid
erDiagram
    STAFF {
        uuid id PK
        string username
        string password_hash
        string hospital
        timestamp created_at
        timestamp updated_at
    }

    PATIENT {
        uuid id PK
        string hospital
        string first_name_th
        string middle_name_th
        string last_name_th
        string first_name_en
        string middle_name_en
        string last_name_en
        string date_of_birth
        string patient_hn
        string national_id
        string passport_id
        string phone_number
        string email
        string gender
        timestamp created_at
        timestamp updated_at
    }

    STAFF ||--o{ PATIENT : "searches within hospital"
```

**Constraints:**
- `username` + `hospital` is unique for staff (same username allowed across different hospitals).
- Each staff member can only search patients where `patient.hospital = staff.hospital`.
- Patient records are synced from external HIS APIs and cached locally.

## API Specification

### Health Check

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Service health check |

### Create Staff

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/staff/create` | No | Register a new hospital staff account |

**Request body:**
```json
{
  "username": "nurse01",
  "password": "secret12",
  "hospital": "hospital-a"
}
```

**Success (201):**
```json
{
  "id": "uuid",
  "username": "nurse01",
  "hospital": "hospital-a"
}
```

### Staff Login

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/staff/login` | No | Authenticate staff and receive JWT |

**Request body:**
```json
{
  "username": "nurse01",
  "password": "secret12",
  "hospital": "hospital-a"
}
```

**Success (200):**
```json
{
  "token": "jwt-token",
  "username": "nurse01",
  "hospital": "hospital-a"
}
```

### Search Patient

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/patient/search` | Bearer JWT | Search patients in staff's hospital |

**Headers:** `Authorization: Bearer <token>`

**Request body (all fields optional):**
```json
{
  "national_id": "1234567890123",
  "passport_id": "",
  "first_name": "Somchai",
  "middle_name": "",
  "last_name": "Jaidee",
  "date_of_birth": "1990-01-15",
  "phone_number": "0812345678",
  "email": "somchai@example.com"
}
```

**Success (200):**
```json
{
  "patients": [
    {
      "id": "uuid",
      "hospital": "hospital-a",
      "first_name_th": "สมชาย",
      "first_name_en": "Somchai",
      "national_id": "1234567890123",
      "gender": "M"
    }
  ]
}
```

When `national_id` or `passport_id` is provided, the middleware calls the external Hospital A API (`GET /patient/search/{id}`), caches the result, and returns matching patients scoped to the logged-in staff member's hospital.

## External Hospital A API (Mock)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/patient/search/{id}` | Search by `national_id` or `passport_id` |

Mock data is served by `cmd/mock-hospital-a` on port 9000 in Docker.

## Quick Start

### Run with Docker Compose

```bash
docker compose up --build
```

Services:
- **Nginx** → http://localhost:8088
- **API** (via Nginx) → http://localhost:8088
- **PostgreSQL** → `localhost:5432` (for pgAdmin / local clients)
- **Mock Hospital A** → internal Docker network only

### Example Flow

```bash
# Create staff
curl -X POST http://localhost:8088/staff/create \
  -H "Content-Type: application/json" \
  -d '{"username":"nurse01","password":"secret12","hospital":"hospital-a"}'

# Login
TOKEN=$(curl -s -X POST http://localhost:8088/staff/login \
  -H "Content-Type: application/json" \
  -d '{"username":"nurse01","password":"secret12","hospital":"hospital-a"}' \
  | jq -r '.token')

# Search patient by national ID
curl -X POST http://localhost:8088/patient/search \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"national_id":"1234567890123"}'
```

### Run Tests

```bash
go test ./...
```

### Run Locally (without Docker)

Requires PostgreSQL running locally.

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/hospital_middleware?sslmode=disable"
export HOSPITAL_A_BASE_URL="http://localhost:9000"

# Terminal 1: mock hospital API
go run ./cmd/mock-hospital-a

# Terminal 2: main API
go run ./cmd/server
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | API server port |
| `DATABASE_URL` | local postgres DSN | PostgreSQL connection string |
| `JWT_SECRET` | `dev-secret-change-in-production` | JWT signing secret |
| `HOSPITAL_A_BASE_URL` | `https://hospital-a.api.co.th` | Hospital A HIS base URL |
