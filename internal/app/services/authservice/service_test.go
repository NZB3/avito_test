package authservice

import (
	"context"
	"os"
	"project/internal/logger"
	"testing"
)

func TestAuthService_Authenticate(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjo5OTk5OTk5OTk5OX0.8lYYtmn7uGSFgNwdBcJO7q5-0yZuC7YN_ek5ohHd6L4"
	os.Setenv("JWT_SECRET", "secret")
	as := New(logger.New())
	admin, err := as.Authenticate(context.Background(), jwt)
	if err != nil {
		t.Error(err)
	}

	t.Log(admin)
}
