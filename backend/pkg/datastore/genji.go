package datastore

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/sql/query"
	dbestimate "github.com/haro87/dokerb/pkg/estimate"
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
	Tasks []Task
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
		if err.Error() != "table already exists" {
			return nil, fmt.Errorf("Unable to create sessions table")
		}
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

// AddTask adds a new task to the specified session
// identified by the provided ID and with an optional summary
func (g GenjiDatastore) AddTask(token, id, summary string) error {
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
	var tasks []Task
	tasks, err = getTasksFromSession(token)

	if !taskExists(tasks, id) {
		tasks = append(tasks, Task{ID: id, Summary: summary})
	} else {
		return fmt.Errorf("Task with ID: %s already part of session", id)
	}

	err = si.db.Exec("UPDATE sessions SET tasks = ? WHERE token = ?", tasks, token)

	return err
}

// RemoveTask removes a task from the specified
// session where the task is identified by the provided
// ID
func (g GenjiDatastore) RemoveTask(token, id string) error {
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

	var tasks []Task

	tasks, err = getTasksFromSession(token)

	if err != nil {
		return fmt.Errorf("Unable to get tasks from session")
	}

	tasks, err = removeTask(tasks, id)

	if err != nil {
		return fmt.Errorf("Unable to remove Task: %s from session", id)
	}

	err = si.db.Exec("UPDATE sessions SET tasks = ? WHERE token = ?", tasks, token)

	return err
}

// AddEstimateToTask adds provided effort and standard deviation estimates
// to the task specified by the given id assigned to a specific
// session identified by the given token
func (g GenjiDatastore) AddEstimateToTask(token, id string, effort, standardDeviation float64) error {
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

	var tasks []Task

	tasks, err = getTasksFromSession(token)

	if err != nil {
		return fmt.Errorf("Unable to get tasks from session")
	}

	if !taskExists(tasks, id) {
		return fmt.Errorf("Task with ID: %s does not exist", id)
	}

	for i, task := range tasks {
		if task.ID == id {
			task.Effort = effort
			task.StandardDeviation = standardDeviation
			tasks[i] = task
			break
		}
	}

	err = si.db.Exec("UPDATE sessions SET tasks = ? WHERE token = ?", tasks, token)

	return err
}

// RemoveEstimateFromTask removes the effort and standard deviation estimates from
// a specified task by the given id assigned to a specific
// session identified by the given token
func (g GenjiDatastore) RemoveEstimateFromTask(token, id string) error {
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

	var tasks []Task

	tasks, err = getTasksFromSession(token)

	if err != nil {
		return fmt.Errorf("Unable to get tasks from session")
	}

	if !taskExists(tasks, id) {
		return fmt.Errorf("Task with ID: %s does not exist", id)
	}

	for i, task := range tasks {
		if task.ID == id {
			task.Effort = 0.0
			task.StandardDeviation = 0.0
			tasks[i] = task
			break
		}
	}

	err = si.db.Exec("UPDATE sessions SET tasks = ? WHERE token = ?", tasks, token)

	return err
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

// GetTasks returns all tasks of a given session
func (g GenjiDatastore) GetTasks(token string) ([]Task, error) {
	if len(token) != defaultTokenLength {
		return []Task{}, fmt.Errorf("Session token does not match desired length")
	}

	se, err := sessionExists(token)
	if !se {
		return []Task{}, fmt.Errorf("Specified session does not exist")
	}

	tasks, err := getTasksFromSession(token)

	return tasks, err
}

// AddEstimate adds a new estimate to the specified session
func (g GenjiDatastore) AddEstimate(token string, estimate Estimate) error {
	if len(token) != defaultTokenLength {
		return fmt.Errorf("Session token does not match desired length")
	}

	if estimate.TaskID == "" {
		return fmt.Errorf("Task ID should not be empty")
	}

	if estimate.UserName == "" {
		return fmt.Errorf("User name should not be empty")
	}

	if _, e := dbestimate.NewDelphiEstimate(estimate.BestCase, estimate.MostLikelyCase, estimate.WorstCase); e != nil {
		return e
	}

	se, err := sessionExists(token)
	if !se {
		return fmt.Errorf("Specified session does not exist")
	}

	var est []Estimate

	est, err = getEstimatesFromSession(token)

	if estimateExists(est, estimate) {
		return fmt.Errorf("Specified estimate already exists")
	}

	var users []string

	users, err = getUsersFromSession(token)

	if !userExists(users, estimate.UserName) {
		return fmt.Errorf("User: %s is not part of session", estimate.UserName)
	}

	var tasks []Task

	tasks, err = getTasksFromSession(token)

	if !taskExists(tasks, estimate.TaskID) {
		return fmt.Errorf("Task with ID: %s is not part of session", estimate.TaskID)
	}

	est = append(est, estimate)

	err = si.db.Exec("UPDATE sessions SET estimates = ? WHERE token = ?", est, token)

	return err
}

// RemoveEstimate removes a existing estimate from the specified session
func (g GenjiDatastore) RemoveEstimate(token string, estimate Estimate) error {
	if len(token) != defaultTokenLength {
		return fmt.Errorf("Session token does not match desired length")
	}

	se, err := sessionExists(token)
	if !se {
		return fmt.Errorf("Specified session does not exist")
	}

	var est []Estimate

	est, err = getEstimatesFromSession(token)

	if err != nil {
		return err
	}

	est, err = removeEstimate(est, estimate)

	if err != nil {
		return err
	}

	err = si.db.Exec("UPDATE sessions SET estimates = ? WHERE token = ?", est, token)

	return err
}

// GetEstimates returns all estimates of a specified session
func (g GenjiDatastore) GetEstimates(token string) ([]Estimate, error) {
	if len(token) != defaultTokenLength {
		return []Estimate{}, fmt.Errorf("Session token does not match desired length")
	}

	se, err := sessionExists(token)
	if !se {
		return []Estimate{}, fmt.Errorf("Specified session does not exist")
	}

	est, err := getEstimatesFromSession(token)

	return est, err
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

func getTasksFromSession(t string) ([]Task, error) {
	var tasks []Task

	res, err := si.db.Query("SELECT tasks FROM sessions WHERE token = ?", t)

	if err != nil {
		return tasks, err
	}

	defer res.Close()

	err = res.Iterate(func(d document.Document) error {
		err = document.Scan(d, &tasks)
		if err != nil {
			return err
		}
		return nil
	})

	return tasks, err
}

func getEstimatesFromSession(t string) ([]Estimate, error) {
	var est []Estimate

	res, err := si.db.Query("SELECT estimates FROM sessions WHERE token = ?", t)

	if err != nil {
		return est, err
	}

	defer res.Close()

	err = res.Iterate(func(d document.Document) error {
		err = document.Scan(d, &est)
		if err != nil {
			return err
		}
		return nil
	})

	return est, err
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

func taskExists(tasks []Task, id string) bool {
	taskExists := false

	for _, task := range tasks {
		if task.ID == id {
			taskExists = true
			break
		}
	}

	return taskExists
}

func estimateExists(estimates []Estimate, estimate Estimate) bool {
	estimateExists := false

	for _, elem := range estimates {
		if elem.TaskID == estimate.TaskID && elem.UserName == estimate.UserName {
			estimateExists = true
			break
		}
	}

	return estimateExists
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

func removeTask(tasks []Task, id string) ([]Task, error) {
	if taskExists(tasks, id) {
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				break
			}
		}
	} else {
		return tasks, fmt.Errorf("Task with ID: %s is not part of session", id)
	}

	return tasks, nil
}

func removeEstimate(estimates []Estimate, estimate Estimate) ([]Estimate, error) {
	if estimateExists(estimates, estimate) {
		for i, e := range estimates {
			if e.TaskID == estimate.TaskID && e.UserName == estimate.UserName {
				estimates = append(estimates[:i], estimates[i+1:]...)
				break
			}
		}
	} else {
		return estimates, fmt.Errorf("Estimate with ID: %s and user name: %s is not part of session",
			estimate.TaskID,
			estimate.UserName)
	}

	return estimates, nil
}
