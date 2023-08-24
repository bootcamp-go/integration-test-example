package storage

import "app/internal/sellers"

func NewStorageSellersMock() *StorageSellersMock {
	return &StorageSellersMock{}
}

// StorageSellersMock is a mock implementation of the StorageSellers interface
type StorageSellersMock struct {
	ReadFunc func() (s map[int]*sellers.Seller, err error)
	WriteFunc func(s map[int]*sellers.Seller) (err error)
	// observers
	Calls struct {
		Read int
		Write int
	}
}

func (m *StorageSellersMock) Read() (s map[int]*sellers.Seller, err error) {
	// observers
	m.Calls.Read++

	// mock func
	s, err = m.ReadFunc()
	return
}

func (m *StorageSellersMock) Write(s map[int]*sellers.Seller) (err error) {
	// observers
	m.Calls.Write++

	// mock func
	err = m.WriteFunc(s)
	return
}