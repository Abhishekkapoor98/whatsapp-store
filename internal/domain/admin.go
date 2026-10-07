package domain

import (
	"context"
)

// Admin represents an administrator account.
type Admin struct {
	ID           int32  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Omitted from JSON responses for security
}

// AdminRepository defines database operations for admins.
type AdminRepository interface {
	GetByUsername(ctx context.Context, username string) (*Admin, error)
}

// AdminUsecase defines authentication logic.
type AdminUsecase interface {
	Login(ctx context.Context, username string, password string) (string, error) // Returns a JWT/Session token
}
