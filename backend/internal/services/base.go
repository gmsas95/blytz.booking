package services

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrConflict     = errors.New("resource already exists")
	ErrBadRequest   = errors.New("invalid request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type BaseService struct {
	DB *gorm.DB
}

func NewBaseService(db *gorm.DB) *BaseService {
	return &BaseService{DB: db}
}
