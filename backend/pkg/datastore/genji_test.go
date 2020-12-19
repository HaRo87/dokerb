package datastore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNilDB(t *testing.T) {
	_, err := NewGenjiDatastore(nil)
	assert.Errorf(t, err, "Proper DB must be provided and not nil")
}

func TestCorrectDBSetup(t *testing.T) {
	m := new(MockGenjiDB)
	_, err := NewGenjiDatastore(m)
	assert.NoError(t, err)

	eq := `CREATE TABLE sessions`

	assert.Equal(t, m.CalledWith()[0], eq)
}
