package datastore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyName(t *testing.T) {
	_, err := NewGenjiDatastore("")
	assert.Errorf(t, err, "Empty name provided")
}
