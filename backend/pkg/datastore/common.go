package datastore

// DataStore defines the common interface a datastore for
// the Doker backend must implement.
type DataStore interface {
	CreateSession() (string, error)
	JoinSession(sessionHash, name string) error
	LeaveSession(sessionHash, name string) error
	AddWorkPackage(id, summary string) error
	RemoveWorkPackage(id string) error
	AddEstimate(id string, effort, standardDeviation float64) error
	RemoveEstimate(id string) error
}
