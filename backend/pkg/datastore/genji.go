package datastore

import (
	"fmt"
)

type GenjiDB interface {
	Exec(q string, args ...interface{}) error
}

type GenjiDatastore struct {
	db GenjiDB
}

func NewGenjiDatastore(db GenjiDB) (DataStore, error) {
	ds := new(GenjiDatastore)

	if db == nil {
		return nil, fmt.Errorf("Proper DB must be provided and not nil")
	}

	ds.db = db

	err := ds.db.Exec("CREATE TABLE sessions")

	if err != nil {
		return nil, fmt.Errorf("Unable to create sessions table")
	}

	return ds, nil
}

func (g GenjiDatastore) CreateSession(timeout int64) (string, error) {
	return "", nil
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
