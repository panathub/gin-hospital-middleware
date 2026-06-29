package auth_test

import (
	"testing"
	"time"

	"github.com/gin-hospital-middleware/internal/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenService_GenerateAndParse(t *testing.T) {
	svc := auth.NewTokenService("test-secret", time.Hour)
	staffID := uuid.New()

	token, err := svc.Generate(staffID, "nurse01", "hospital-a")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := svc.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, staffID, claims.StaffID)
	assert.Equal(t, "nurse01", claims.Username)
	assert.Equal(t, "hospital-a", claims.Hospital)
}

func TestTokenService_Parse_InvalidToken(t *testing.T) {
	svc := auth.NewTokenService("test-secret", time.Hour)
	_, err := svc.Parse("not-a-valid-token")
	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}
