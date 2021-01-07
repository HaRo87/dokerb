package datastore

import (
	"fmt"
	"github.com/genjidb/genji/sql/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecNoError(t *testing.T) {
	var ds GenjiDB
	m := new(MockGenjiDB)
	ds = m

	m.On("Exec", "CREATE TABLE sessions").Return(nil)

	err := ds.Exec("CREATE TABLE sessions")

	assert.NoError(t, err)
	m.MethodCalled("Exec", "CREATE TABLE sessions")
}

func TestExecError(t *testing.T) {
	var ds GenjiDB
	m := new(MockGenjiDB)
	ds = m

	m.On("Exec", "CREATE TABLE sessions").Return(fmt.Errorf("Some error"))

	err := ds.Exec("CREATE TABLE sessions")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("Exec", "CREATE TABLE sessions")
}

func TestQueryNoError(t *testing.T) {
	var ds GenjiDB
	m := new(MockGenjiDB)
	ds = m

	m.On("Query", "SELECT * FROM sessions").Return(new(query.Result), nil)

	_, err := ds.Query("SELECT * FROM sessions")

	assert.NoError(t, err)
	m.MethodCalled("Query", "SELECT * FROM sessions")
}

func TestQueryError(t *testing.T) {
	var ds GenjiDB
	m := new(MockGenjiDB)
	ds = m

	m.On("Query", "SELECT * FROM sessions").Return(new(query.Result), fmt.Errorf("Some error"))

	_, err := ds.Query("SELECT * FROM sessions")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("Query", "SELECT * FROM sessions")
}
