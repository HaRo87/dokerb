package datastore

// DataStore defines the common interface a datastore for
// the Doker backend must implement.
type DataStore interface {
	CreateStore(path string) error
	DeleteStore(path string) error
	Close() error
	CreateSession(timeout int64) (string, error)
	JoinSession(sessionHash, name string) error
	LeaveSession(sessionHash, name string) error
	AddWorkPackage(id, summary string) error
	RemoveWorkPackage(id string) error
	AddEstimate(id string, effort, standardDeviation float64) error
	RemoveEstimate(id string) error
}
