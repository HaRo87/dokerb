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

func TestSingletonPatternWorks(t *testing.T) {
	m := new(MockGenjiDB)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	g1, err := NewGenjiDatastore(m)
	g2, err2 := NewGenjiDatastore(m)
	assert.NoError(t, err)
	assert.NoError(t, err2)

	assert.Equal(t, g1, g2)
}

func TestGenerateTokenWrongLength(t *testing.T) {
	_, err := generateToken(-20)
	assert.Errorf(t, err, "Invalid token length provided: -20, should be >= 20")
}

func TestGenerateTokenWithMinimalSuggestedLength(t *testing.T) {
	token, err := generateToken(20)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(token))
}

func TestGenerateTokenWithProperLength(t *testing.T) {
	token, err := generateToken(32)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(token))
}

func TestCorrectDBSetup(t *testing.T) {
	m := new(MockGenjiDB)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	_, err := NewGenjiDatastore(m)
	assert.NoError(t, err)

	m.MethodCalled("Exec", "CREATE TABLE sessions")
}

func TestCreateSessionFailsDueToExecError(t *testing.T) {
	m := new(MockGenjiDB)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	ds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	m.On("Exec", "INSERT INTO sessions (id) VALUES (?)").Return(fmt.Errorf("Ooops, something went wrong"))
	token, err2 := ds.CreateSession()
	assert.Empty(t, token)
	assert.Errorf(t, err2, "Unable to store session token")
	m.MethodCalled("Exec", "INSERT INTO sessions (id) VALUES (?)")
}

func TestCreateSessionSuccess(t *testing.T) {
	m := new(MockGenjiDB)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	ds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	m.On("Exec", "INSERT INTO sessions (id) VALUES (?)").Return(nil)
	token, err2 := ds.CreateSession()
	assert.Equal(t, 32, len(token))
	assert.NoError(t, err2)
	m.MethodCalled("Exec", "INSERT INTO sessions (id) VALUES (?)")
}
