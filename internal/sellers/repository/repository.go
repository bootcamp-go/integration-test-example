package repository

import (
	"app/internal/sellers"
	"errors"
)

// RepositorySellers is an interface that defines the methods the repository needs to implement
type RepositorySellers interface {
	// GetById returns a seller by id
	GetById(id int) (s *sellers.Seller, err error)

	// Save saves a seller
	Save(s *sellers.Seller) (err error)
}

var (
	ErrRepositorySellersInternal = errors.New("repository error: internal error")
	ErrRepositorySellersNotFound = errors.New("repository error: seller not found")
)