package storage

import (
	"app/internal/sellers"
	"errors"
)

type StorageSellers interface {
	// Read returns all sellers from the storage
	Read() (s map[int]*sellers.Seller, err error)

	// Write writes all sellers to the storage
	Write(s map[int]*sellers.Seller) error
}

var (
	ErrStorageSellersInternal = errors.New("storage sellers error: internal error")
)