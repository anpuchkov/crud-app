package handler

import (
	"Api/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type Services struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Services {
	return &Services{db}
}

func (s *Services) Create(ctx context.Context, service domain.Service) error {
	_, err := s.db.Exec(ctx, "INSERT INTO service (service_name, description, service_time) VALUES ($1, $2, $3)",
		service.ServiceName, service.ServiceDescription, service.ServiceTime)

	return err
}

func (s *Services) GetByID(ctx context.Context, id int64) (domain.Service, error) {
	var service domain.Service
	err := s.db.QueryRow(ctx, "SELECT id, service_name, service_time, description from service WHERE id=$1", id).
		Scan(&service.Id, &service.ServiceName, &service.ServiceTime, &service.ServiceDescription)
	if errors.Is(err, sql.ErrNoRows) {
		return service, domain.ErrServiceNotFound
	}
	return service, err
}

func (s *Services) GetAll(ctx context.Context) ([]domain.Service, error) {
	services := make([]domain.Service, 0)
	rows, err := s.db.Query(ctx, "SELECT id, service_name, description, service_time FROM service")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var service domain.Service
		err = rows.Scan(&service.Id, &service.ServiceName, &service.ServiceDescription, &service.ServiceTime)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, rows.Err()
}

func (s *Services) Delete(ctx context.Context, id int64) error {
	_, err := s.db.Exec(ctx, "DELETE FROM service WHERE id=$1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrServiceNotFound
	}
	return err

}

func (s *Services) Update(ctx context.Context, id int64, input domain.UpdateService) error {
	setVals := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.ServiceName != nil {
		setVals = append(setVals, fmt.Sprintf("service_name=$%d", argId))
		args = append(args, *input.ServiceName)
		argId++
	}
	if input.ServiceDescription != nil {
		setVals = append(setVals, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.ServiceDescription)
		argId++
	}
	if input.ServiceTime != nil {
		setVals = append(setVals, fmt.Sprintf("service_time=$%d", argId))
		args = append(args, *input.ServiceTime)
		argId++
	}
	setQuery := strings.Join(setVals, ", ")
	query := fmt.Sprintf("UPDATE service SET %s WHERE id = $%d", setQuery, argId)

	args = append(args, id)
	_, err := s.db.Exec(ctx, query, args...)
	return err
}
