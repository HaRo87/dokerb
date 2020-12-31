package estimate

import (
	"fmt"
)

// DelphiEstimate struct holds all initially needed values
// for calculating the effort and the standard deviation based
// on the Delphi method.
type DelphiEstimate struct {
	BestCase   float64
	MostLikely float64
	WorstCase  float64
}

// NewDelphiEstimate constructor returns a new DelphiEstimate
// object containing the provided values where b is for best effort,
// m is for most likely and w is for worst case.
func NewDelphiEstimate(b, m, w float64) (Estimator, error) {
	e := new(DelphiEstimate)

	if b < 0 {
		return nil, fmt.Errorf("Best case must be >= 0, provided: %g", b)
	}
	e.BestCase = b

	if m < 0 {
		return nil, fmt.Errorf("Most Likely must be >= 0, provided: %g", m)
	}
	if m < b {
		return nil, fmt.Errorf("Most Likely was smaller than Best Effort")
	}
	e.MostLikely = m

	if w < 0 {
		return nil, fmt.Errorf("Worst Case must be >= 0, provided: %g", w)
	}
	if w < m {
		return nil, fmt.Errorf("Worst Case was smaller than Most Likely")
	}
	e.WorstCase = w

	return e, nil
}

// GetEffort returns the calculated effort.
func (d DelphiEstimate) GetEffort() float64 {
	return (d.BestCase + 4*d.MostLikely + d.WorstCase) / 6
}

// GetStandardDeviation returns the calculated standard deviation.
func (d DelphiEstimate) GetStandardDeviation() float64 {
	return (d.WorstCase - d.BestCase) / 6
}
