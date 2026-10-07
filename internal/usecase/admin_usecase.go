package usecase

import (
	"context"
	"errors"
	"time"
	"whatsapp-store/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type adminUsecase struct {
	adminRepo domain.AdminRepository
	jwtSecret string
}

// NewAdminUsecase initializes the admin business logic
func NewAdminUsecase(repo domain.AdminRepository, secret string) domain.AdminUsecase {
	return &adminUsecase{
		adminRepo: repo,
		jwtSecret: secret,
	}
}

func (u *adminUsecase) Login(ctx context.Context, username, password string) (string, error) {
	// 1. Fetch admin from database
	admin, err := u.adminRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// 2. Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// 3. Generate a 24-hour valid JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
