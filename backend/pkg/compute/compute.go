package compute

import (
	"fmt"
	"github.com/haro87/dokerb/pkg/datastore"
	"github.com/haro87/dokerb/pkg/estimate"
)

// CalculateAverageEstimate calculates the average estimate of all provided
// estimates matching a given work package ID
func CalculateAverageEstimate(estimates []datastore.Estimate, id string) (estimate.Estimator, error) {
	ests, err := extractEstimatesForWorkPackage(estimates, id)

	if err != nil {
		return nil, err
	}

	var b float64
	var m float64
	var w float64

	for _, est := range ests {
		b += est.BestCase
		m += est.MostLikelyCase
		w += est.WorstCase
	}

	b = b / float64(len(ests))
	m = m / float64(len(ests))
	w = w / float64(len(ests))

	var est estimate.Estimator

	est, err = estimate.NewDelphiEstimate(b, m, w)

	return est, err

}

func extractEstimatesForWorkPackage(estimates []datastore.Estimate, id string) ([]datastore.Estimate, error) {
	if id == "" {
		return []datastore.Estimate{}, fmt.Errorf("Work Package ID cannot be empty")
	}
	if len(estimates) < 1 {
		return []datastore.Estimate{}, fmt.Errorf("Not enough data to process")
	}
	var ests []datastore.Estimate

	for _, est := range estimates {
		if est.WorkPackageID == id {
			ests = append(ests, est)
		}
	}

	if len(ests) < 1 {
		return []datastore.Estimate{}, fmt.Errorf("Specified work package with ID: %s is not part of estimates", id)
	}

	return ests, nil
}
