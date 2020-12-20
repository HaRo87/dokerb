package datastore

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/genjidb/genji/sql/query"
	"sync"
)

// GenjiDB  abstracts the 3rd party genji deps
// for easier testing
type GenjiDB interface {
	Exec(q string, args ...interface{}) error
	Query(q string, args ...interface{}) (*query.Result, error)
}

// GenjiDatastore struct which holds the actual database
type GenjiDatastore struct {
	db GenjiDB
}

var lock = &sync.Mutex{}

var si *GenjiDatastore

const defaultTokenLength int = 32

// NewGenjiDatastore creates a new GenjiDatastore following
// the singleton design pattern
func NewGenjiDatastore(db GenjiDB) (DataStore, error) {
	if si == nil {
		lock.Lock()
		defer lock.Unlock()
		si = new(GenjiDatastore)
	}

	if db == nil {
		return nil, fmt.Errorf("Proper DB must be provided and not nil")
	}

	si.db = db

	err := si.db.Exec("CREATE TABLE sessions")

	if err != nil {
		return nil, fmt.Errorf("Unable to create sessions table")
	}

	return si, nil
}

// CreateSession creates a new session by generating a
// new session token and storing it in the datastore
func (g GenjiDatastore) CreateSession() (string, error) {
	st, err := generateToken(defaultTokenLength)
	if err != nil {
		return "", fmt.Errorf("Unable to create session token")
	}
	err = si.db.Exec("INSERT INTO sessions (id) VALUES (?)", st)
	if err != nil {
		return "", fmt.Errorf("Unable to store session token")
	}
	return st, nil
}

func (g GenjiDatastore) JoinSession(sessionHash, name string) error {
	return nil
}

func (g GenjiDatastore) LeaveSession(sessionHash, name string) error {
	return nil
}

func (g GenjiDatastore) AddWorkPackage(id, summary string) error {
	return nil
}

func (g GenjiDatastore) RemoveWorkPackage(id string) error {
	return nil
}

func (g GenjiDatastore) AddEstimate(id string, effort, standardDeviation float64) error {
	return nil
}

func (g GenjiDatastore) RemoveEstimate(id string) error {
	return nil
}

func generateToken(l int) (string, error) {
	if l <= 0 {
		return "", fmt.Errorf("Invalid token length provided: %d, should be >= 20", l)
	}
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Unable to generate token")
	}

	return hex.EncodeToString(b)[0:l], nil
}
