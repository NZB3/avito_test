package authservice

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
	"project/internal/logger"
)

type authService struct {
	log    logger.Logger
	secret []byte
}

type claims struct {
	Admin  bool `json:"admin,omitempty"`
	UserID int  `json:"user_id"`
	jwt.StandardClaims
}

func New(log logger.Logger) *authService {
	s := os.Getenv("JWT_SECRET")
	return &authService{
		log:    log,
		secret: []byte(s),
	}
}

func (a *authService) Authenticate(ctx context.Context, tokenString string) (admin bool, err error) {
	const op = "authservice.Authenticate"
	var c claims
	token, err := jwt.ParseWithClaims(tokenString, &c, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})

	if err != nil {
		a.log.Errorf("%s Failed to parse token: %v", op, err)
		return false, err
	}

	if !token.Valid {
		return false, errors.New("invalid token")
	}

	return c.Admin, nil
}
