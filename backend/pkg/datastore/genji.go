package datastore

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/sql/query"
	"sync"
)

// GenjiDB  abstracts the 3rd party genji deps
// for easier testing
type GenjiDB interface {
	Exec(q string, args ...interface{}) error
	Query(q string, args ...interface{}) (*query.Result, error)
	Update(fn func(tx *genji.Tx) error) error
}

// GenjiDatastore struct which holds the actual database
type GenjiDatastore struct {
	db GenjiDB
}

type session struct {
	Token string
	Users []string
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
	s := session{
		Token: st,
		Users: []string{},
	}
	err = si.db.Exec("INSERT INTO sessions VALUES ?", &s)
	if err != nil {
		return "", fmt.Errorf("Unable to store session token")
	}
	return st, nil
}

// JoinSession allows a user with the specified name to join a
// session identified by the given token
func (g GenjiDatastore) JoinSession(token, name string) error {
	if len(token) != defaultTokenLength {
		return fmt.Errorf("Session token does not match desired length")
	}
	if name == "" {
		return fmt.Errorf("User name should not be empty")
	}

	se, err := sessionExists(token)
	if !se {
		return fmt.Errorf("Specified session does not exist")
	}
	var u []string
	u, err = getUsersFromSession(token)

	if !userExists(u, name) {
		u = append(u, name)
	} else {
		return fmt.Errorf("User with name: %s already part of session", name)
	}

	err = si.db.Exec("UPDATE sessions SET users = ? WHERE token = ?", u, token)

	return err
}

func (g GenjiDatastore) LeaveSession(token, name string) error {
	if len(token) != defaultTokenLength {
		return fmt.Errorf("Session token does not match desired length")
	}
	if name == "" {
		return fmt.Errorf("User name should not be empty")
	}

	se, err := sessionExists(token)
	if !se {
		return fmt.Errorf("Specified session does not exist")
	}

	var u []string

	u, err = getUsersFromSession(token)

	if err != nil {
		return fmt.Errorf("Unable to get Users from session")
	}

	u, err = removeUser(u, name)

	if err != nil {
		return fmt.Errorf("Unable to remove user: %s from session", name)
	}

	err = si.db.Exec("UPDATE sessions SET users = ? WHERE token = ?", u, token)

	return err
}

func (g GenjiDatastore) CloseSession(token string) error {
	return nil
}

func (g GenjiDatastore) AddWorkPackage(token, id, summary string) error {
	return nil
}

func (g GenjiDatastore) RemoveWorkPackage(token, id string) error {
	return nil
}

func (g GenjiDatastore) AddEstimate(token, id string, effort, standardDeviation float64) error {
	return nil
}

func (g GenjiDatastore) RemoveEstimate(token, id string) error {
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

func sessionExists(t string) (bool, error) {
	var tokens []string
	sessionExists := false
	res, err := si.db.Query("SELECT token FROM sessions")
	defer res.Close()

	err = res.Iterate(func(d document.Document) error {
		var token string
		err = document.Scan(d, &token)
		if err != nil {
			return err
		}
		tokens = append(tokens, token)
		return nil
	})

	for _, elem := range tokens {
		if elem == t {
			sessionExists = true
			break
		}
	}

	return sessionExists, err
}

func getUsersFromSession(t string) ([]string, error) {
	var users []string

	res, err := si.db.Query("SELECT users FROM sessions WHERE token = ?", t)

	if err != nil {
		return users, err
	}

	defer res.Close()

	err = res.Iterate(func(d document.Document) error {
		err = document.Scan(d, &users)
		if err != nil {
			return err
		}
		return nil
	})

	return users, err
}

func userExists(users []string, user string) bool {
	userExists := false

	for _, elem := range users {
		if elem == user {
			userExists = true
			break
		}
	}

	return userExists
}

func removeUser(users []string, user string) ([]string, error) {
	if userExists(users, user) {
		for i, e := range users {
			if e == user {
				users = append(users[:i], users[i+1:]...)
				break
			}
		}
	} else {
		return users, fmt.Errorf("User with name: %s is not part of session", user)
	}

	return users, nil
}
