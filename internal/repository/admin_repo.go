package repository

import (
	"context"
	"whatsapp-store/internal/domain"
)

type adminRepo struct {
	q *Queries
}

// NewAdminRepository creates a new instance of the domain.AdminRepository
func NewAdminRepository(q *Queries) domain.AdminRepository {
	return &adminRepo{q: q}
}

func (r *adminRepo) GetByUsername(ctx context.Context, username string) (*domain.Admin, error) {
	admin, err := r.q.GetAdminByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &domain.Admin{
		ID:           admin.ID,
		Username:     admin.Username,
		PasswordHash: admin.PasswordHash,
	}, nil
}
