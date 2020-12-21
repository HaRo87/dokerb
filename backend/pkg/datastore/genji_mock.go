package datastore

import (
	"github.com/genjidb/genji"
	"github.com/genjidb/genji/sql/query"
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

// Query implements the GenjiDB interface
func (m *MockGenjiDB) Query(q string, args ...interface{}) (*query.Result, error) {
	arguments := m.Called(q)
	return arguments.Get(0).(*query.Result), arguments.Error(1)
}

func (m *MockGenjiDB) Update(fn func(tx *genji.Tx) error) error {
	arguments := m.Called(fn)
	return arguments.Error(0)
}
