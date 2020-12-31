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
	Token        string
	Users        []string
	Workpackages []WorkPackage
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

// LeaveSession allows a user with the specified name to leave a
// session identified by the provided token
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

// RemoveSession deletes a session from the datastore
func (g GenjiDatastore) RemoveSession(token string) error {
	if len(token) != defaultTokenLength {
		return fmt.Errorf("Session token does not match desired length")
	}

	se, err := sessionExists(token)
	if !se {
		return fmt.Errorf("Specified session does not exist")
	}

	err = si.db.Exec("DELETE FROM sessions WHERE token = ?", token)

	return err
}

// AddWorkPackage adds a new work package to the specified session
// identified by the provided ID and with an optional summary
func (g GenjiDatastore) AddWorkPackage(token, id, summary string) error {
	if len(token) != defaultTokenLength {
		return fmt.Errorf("Session token does not match desired length")
	}
	if id == "" {
		return fmt.Errorf("ID should not be empty")
	}

	se, err := sessionExists(token)
	if !se {
		return fmt.Errorf("Specified session does not exist")
	}
	var wps []WorkPackage
	wps, err = getWorkpackagesFromSession(token)

	if !workpackageExists(wps, id) {
		wps = append(wps, WorkPackage{ID: id, Summary: summary})
	} else {
		return fmt.Errorf("Workpackage with ID: %s already part of session", id)
	}

	err = si.db.Exec("UPDATE sessions SET workpackages = ? WHERE token = ?", wps, token)

	return err
}

// RemoveWorkPackage removes a work package from the specified
// session where the work package is identified by the provided
// ID
func (g GenjiDatastore) RemoveWorkPackage(token, id string) error {
	if len(token) != defaultTokenLength {
		return fmt.Errorf("Session token does not match desired length")
	}
	if id == "" {
		return fmt.Errorf("ID should not be empty")
	}

	se, err := sessionExists(token)
	if !se {
		return fmt.Errorf("Specified session does not exist")
	}

	var wps []WorkPackage

	wps, err = getWorkpackagesFromSession(token)

	if err != nil {
		return fmt.Errorf("Unable to get workpackages from session")
	}

	wps, err = removeWorkpackage(wps, id)

	if err != nil {
		return fmt.Errorf("Unable to remove workpackage: %s from session", id)
	}

	err = si.db.Exec("UPDATE sessions SET workpackages = ? WHERE token = ?", wps, token)

	return err
}

func (g GenjiDatastore) AddEstimate(token, id string, effort, standardDeviation float64) error {
	if len(token) != defaultTokenLength {
		return fmt.Errorf("Session token does not match desired length")
	}
	if id == "" {
		return fmt.Errorf("ID should not be empty")
	}
	if effort < 0 {
		return fmt.Errorf("Effort < 0 not allowed")
	}
	if standardDeviation < 0 {
		return fmt.Errorf("Standard deviation < 0 not allowed")
	}

	se, err := sessionExists(token)
	if !se {
		return fmt.Errorf("Specified session does not exist")
	}

	var wps []WorkPackage

	wps, err = getWorkpackagesFromSession(token)

	if err != nil {
		return fmt.Errorf("Unable to get workpackages from session")
	}

	if !workpackageExists(wps, id) {
		return fmt.Errorf("Work package with ID: %s does not exist", id)
	}

	for _, elem := range wps {
		if elem.ID == id {
			elem.Effort = effort
			elem.StandardDeviation = standardDeviation
			break
		}
	}

	err = si.db.Exec("UPDATE sessions SET workpackages = ? WHERE token = ?", wps, token)

	return err
}

func (g GenjiDatastore) RemoveEstimate(token, id string) error {
	return nil
}

// GetUsers returns all users of a given session
func (g GenjiDatastore) GetUsers(token string) ([]string, error) {
	if len(token) != defaultTokenLength {
		return []string{}, fmt.Errorf("Session token does not match desired length")
	}

	se, err := sessionExists(token)
	if !se {
		return []string{}, fmt.Errorf("Specified session does not exist")
	}

	users, err := getUsersFromSession(token)

	return users, err
}

// GetWorkPackages returns all work packages of a given session
func (g GenjiDatastore) GetWorkPackages(token string) ([]WorkPackage, error) {
	if len(token) != defaultTokenLength {
		return []WorkPackage{}, fmt.Errorf("Session token does not match desired length")
	}

	se, err := sessionExists(token)
	if !se {
		return []WorkPackage{}, fmt.Errorf("Specified session does not exist")
	}

	workpackages, err := getWorkpackagesFromSession(token)

	return workpackages, err
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

func getWorkpackagesFromSession(t string) ([]WorkPackage, error) {
	var wps []WorkPackage

	res, err := si.db.Query("SELECT workpackages FROM sessions WHERE token = ?", t)

	if err != nil {
		return wps, err
	}

	defer res.Close()

	err = res.Iterate(func(d document.Document) error {
		err = document.Scan(d, &wps)
		if err != nil {
			return err
		}
		return nil
	})

	return wps, err
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

func workpackageExists(workpackages []WorkPackage, id string) bool {
	workpackageExists := false

	for _, elem := range workpackages {
		if elem.ID == id {
			workpackageExists = true
			break
		}
	}

	return workpackageExists
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

func removeWorkpackage(workpackages []WorkPackage, id string) ([]WorkPackage, error) {
	if workpackageExists(workpackages, id) {
		for i, e := range workpackages {
			if e.ID == id {
				workpackages = append(workpackages[:i], workpackages[i+1:]...)
				break
			}
		}
	} else {
		return workpackages, fmt.Errorf("Workpackage with ID: %s is not part of session", id)
	}

	return workpackages, nil
}
