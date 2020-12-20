package datastore

import (
	"github.com/stretchr/testify/mock"
)

// MockGenjiDB represents the mocked object
type MockGenjiDB struct {
	mock.Mock
}

// Exec implements the GenjiDB interface
func (m *MockGenjiDB) Exec(q string, args ...interface{}) error {
	arguments := m.Called(q)
	return arguments.Error(0)
}
