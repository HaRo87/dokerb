package datastore

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNilDB(t *testing.T) {
	_, err := NewGenjiDatastore(nil)
	assert.Errorf(t, err, "Proper DB must be provided and not nil")
}

func TestErrorOnDBSetupCreateTable(t *testing.T) {
	m := new(MockGenjiDB)
	m.On("Exec", "CREATE TABLE sessions").Return(fmt.Errorf("Ooops, something went wrong"))
	_, err := NewGenjiDatastore(m)
	assert.Errorf(t, err, "Unable to create sessions table")

	m.MethodCalled("Exec", "CREATE TABLE sessions")
}

func TestCorrectDBSetup(t *testing.T) {
	m := new(MockGenjiDB)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	_, err := NewGenjiDatastore(m)
	assert.NoError(t, err)

	m.MethodCalled("Exec", "CREATE TABLE sessions")
}
