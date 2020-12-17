package datastore

import (
	"fmt"
	"github.com/genjidb/genji"
	"golang.org/x/sys/unix"
	"os"
)

type GenjiDatastore struct {
	Name string
	db   genji.DB
}

func NewGenjiDatastore(name string) (DataStore, error) {
	ds := new(GenjiDatastore)

	if name == "" {
		return nil, fmt.Errorf("Empty name provided")
	}

	ds.Name = name

	return ds, nil
}

func (g GenjiDatastore) CreateStore(path string) error {
	if unix.Access(path, unix.W_OK) != nil {
		return fmt.Errorf("Provided path: %s not writable", path)
	}

	db, err := genji.Open(path + g.Name)
	g.db = *db

	return err
}

func (g GenjiDatastore) DeleteStore(path string) error {
	if path == "" {
		return fmt.Errorf("Empty path provided")
	}
	g.db.Close()
	os.Remove(path + g.Name)
	return nil
}

func (g GenjiDatastore) Close() error {
	return g.db.Close()
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
