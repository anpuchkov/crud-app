package service

import (
	"Api/domain"
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Repository interface {
	Create(ctx context.Context, service domain.Service) error
	GetByID(ctx context.Context, id int64) (domain.Service, error)
	GetAll(ctx context.Context) ([]domain.Service, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, service domain.UpdateService) error
}

type Services struct {
	repo Repository
}

func NewService(repo Repository) *Services {
	return &Services{
		repo: repo,
	}
}

func (s *Services) Create(ctx context.Context, service domain.Service) error {
	service.ServiceTime = pgtype.Timestamp{Time: time.Now()}
	return s.repo.Create(ctx, service)
}

func (s *Services) GetByID(ctx context.Context, id int64) (domain.Service, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Services) GetAll(ctx context.Context) ([]domain.Service, error) {
	return s.repo.GetAll(ctx)
}

func (s *Services) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Services) Update(ctx context.Context, id int64, input domain.UpdateService) error {
	return s.repo.Update(ctx, id, input)
}
