package services

import (
	db "github.com/atsuyaourt/xyz-books/internal/db/sqlc"
)

type DefaultService struct {
	store db.Store
}

func NewDefaultService(store db.Store) (*DefaultService, error) {
	s := &DefaultService{
		store: store,
	}

	return s, nil
}
