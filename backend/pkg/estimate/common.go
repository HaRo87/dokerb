package estimate

// Estimator defines the common interface a estimator needs
// to implement.
type Estimator interface {
	GetEffort() float64
	GetStandardDeviation() float64
}
