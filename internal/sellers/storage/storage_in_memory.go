package storage

import "app/internal/sellers"

// NewStorageInMemory returns a new StorageSellersInMemory
type ConfigStorageInMemory struct {
	db map[int]*sellers.Seller
}
func NewStorageInMemory(cfg *ConfigStorageInMemory) (st *StorageSellersInMemory) {
	// default values
	cfgDefaults := &ConfigStorageInMemory{
		db:     make(map[int]*sellers.Seller),
	}
	if cfg != nil {
		// use the values from the config
		cfgDefaults.db = cfg.db
	}
	
	// create storage
	st = &StorageSellersInMemory{db: cfgDefaults.db}
	return
}

// StorageSellersInMemory is a storage that holds the data in memory
type StorageSellersInMemory struct {
	// db
	db map[int]*sellers.Seller
}

// Read returns all sellers from the storage
func (st *StorageSellersInMemory) Read() (s map[int]*sellers.Seller, err error) {
	// return a copy of the db
	s = make(map[int]*sellers.Seller)
	for k, v := range st.db {
		s[k] = v
	}
	
	return
}

// Write writes all sellers to the storage
func (st *StorageSellersInMemory) Write(s map[int]*sellers.Seller) (err error) {
	// write all sellers to the db
	newS := make(map[int]*sellers.Seller)
	for k, v := range s {
		newS[k] = v
	}

	(*st).db = newS
	return
}