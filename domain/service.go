package domain

import (
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrServiceNotFound = errors.New("service not found")
)

type Service struct {
	Id                 int64            `json:"id"`
	ServiceName        string           `json:"service_name"`
	ServiceDescription string           `json:"description"`
	ServiceTime        pgtype.Timestamp `json:"service_time"`
}

type UpdateService struct {
	ServiceName        *string           `json:"service_name"`
	ServiceDescription *string           `json:"description"`
	ServiceTime        *pgtype.Timestamp `json:"service_time"`
}
