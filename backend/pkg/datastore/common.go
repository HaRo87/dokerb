package datastore

// DataStore defines the common interface a datastore for
// the Doker backend must implement.
type DataStore interface {
	CreateSession() (string, error)
	JoinSession(token, name string) error
	LeaveSession(token, name string) error
	AddWorkPackage(token, id, summary string) error
	RemoveWorkPackage(token, id string) error
	AddEstimate(token, id string, effort, standardDeviation float64) error
	RemoveEstimate(token, id string) error
}
