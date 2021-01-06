package datastore

// DataStore defines the common interface a datastore for
// the Doker backend must implement.
type DataStore interface {
	CreateSession() (string, error)
	JoinSession(token, name string) error
	LeaveSession(token, name string) error
	RemoveSession(token string) error
	AddWorkPackage(token, id, summary string) error
	RemoveWorkPackage(token, id string) error
	AddEstimateToWorkPackage(token, id string, effort, standardDeviation float64) error
	RemoveEstimateFromWorkPackage(token, id string) error
	GetUsers(token string) ([]string, error)
	GetWorkPackages(token string) ([]WorkPackage, error)
	AddEstimate(token string, estimate Estimate) error
	RemoveEstimate(token string, estimate Estimate) error
	GetEstimates(token string) ([]Estimate, error)
}

// WorkPackage defines a single work package
type WorkPackage struct {
	ID                string
	Summary           string
	Effort            float64
	StandardDeviation float64
}

// Estimate defines a user estimate for a specific
// work package
type Estimate struct {
	WorkPackageID  string
	UserName       string
	BestCase       float64
	MostLikelyCase float64
	WorstCase      float64
}
