package repository

import (
	"app/internal/sellers"
	"app/internal/sellers/storage"
)

// NewRepositorySellersDefault creates a new repository
func NewRepositorySellersDefault(st storage.StorageSellers) *RepositorySellersDefault {
	return &RepositorySellersDefault{
		st: st,
	}
}

// RepositorySellersDefault is the default repository
type RepositorySellersDefault struct {
	st storage.StorageSellers
}

// GetById returns a seller by id
func (r *RepositorySellersDefault) GetById(id int) (s *sellers.Seller, err error) {
	sellers, err := r.st.Read()
	if err != nil {
		err = ErrRepositorySellersInternal
		return
	}

	s, ok := sellers[id]
	if !ok {
		err = ErrRepositorySellersNotFound
		return
	}

	return
}

// Save saves a seller
func (r *RepositorySellersDefault) Save(s *sellers.Seller) (err error) {
	// read the sellers
	sellers, err := r.st.Read()
	if err != nil {
		err = ErrRepositorySellersInternal
		return
	}

	// get the highest id
	var highestId int
	var cc int
	for k := range sellers {
		if cc == 0 {
			highestId = k
			cc++
		}
		if k > highestId {
			highestId = k
		}
	}

	// set the id
	s.ID = highestId + 1

	// add the seller
	sellers[s.ID] = s

	// write the sellers
	err = r.st.Write(sellers)
	if err != nil {
		err = ErrRepositorySellersInternal
		return
	}

	return
}