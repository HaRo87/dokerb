package compute

import (
	"fmt"
	"github.com/haro87/dokerb/pkg/datastore"
	"github.com/haro87/dokerb/pkg/estimate"
	"sort"
)

// CalculateAverageEstimate calculates the average estimate of all provided
// estimates matching a given work package ID
func CalculateAverageEstimate(estimates []datastore.Estimate, id string) (estimate.Estimator, error) {
	ests, err := ExtractEstimatesForTask(estimates, id)

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

// GetUsersWithMaxDistanceBetweenEffort returns the two users, if they
// exist, who have the max distance between their effort estimates
func GetUsersWithMaxDistanceBetweenEffort(estimates []datastore.Estimate, id string) ([]string, error) {
	ests, err := ExtractEstimatesForTask(estimates, id)

	if err != nil {
		return nil, err
	}

	var list []estimate.UserEstimate

	for _, est := range ests {
		es, e := estimate.NewDelphiEstimate(est.BestCase, est.MostLikelyCase, est.WorstCase)
		if e != nil {
			return []string{}, e
		}
		list = append(list, estimate.UserEstimate{Name: est.UserName, Estimate: es})
	}

	l, le := estimate.NewEstimateList(list)

	if le != nil {
		return []string{}, le
	}

	sort.Sort(l)

	return []string{l.GetFirstUser(), l.GetLastUser()}, nil
}

// ExtractEstimatesForTask extracts all estimates for a specified
// work package ID
func ExtractEstimatesForTask(estimates []datastore.Estimate, id string) ([]datastore.Estimate, error) {
	if id == "" {
		return []datastore.Estimate{}, fmt.Errorf("Task ID cannot be empty")
	}
	if len(estimates) < 1 {
		return []datastore.Estimate{}, fmt.Errorf("Not enough data to process")
	}
	var ests []datastore.Estimate

	for _, est := range estimates {
		if est.TaskID == id {
			ests = append(ests, est)
		}
	}

	if len(ests) < 1 {
		return []datastore.Estimate{}, fmt.Errorf("Specified task with ID: %s is not part of estimates", id)
	}

	return ests, nil
}
