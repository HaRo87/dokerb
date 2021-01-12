package datastore

// DataStore defines the common interface a datastore for
// the Doker backend must implement.
type DataStore interface {
	CreateSession() (string, error)
	JoinSession(token, name string) error
	LeaveSession(token, name string) error
	RemoveSession(token string) error
	AddTask(token, id, summary string) error
	RemoveTask(token, id string) error
	AddEstimateToTask(token, id string, effort, standardDeviation float64) error
	RemoveEstimateFromTask(token, id string) error
	GetUsers(token string) ([]string, error)
	GetTasks(token string) ([]Task, error)
	AddEstimate(token string, estimate Estimate) error
	RemoveEstimate(token string, estimate Estimate) error
	GetEstimates(token string) ([]Estimate, error)
}

// Task defines a single task
type Task struct {
	ID                string
	Summary           string
	Effort            float64
	StandardDeviation float64
}

// Estimate defines a user estimate for a specific
// task
type Estimate struct {
	TaskID         string
	UserName       string
	BestCase       float64
	MostLikelyCase float64
	WorstCase      float64
}
