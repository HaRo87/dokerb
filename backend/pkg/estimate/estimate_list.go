package estimate

import (
	"fmt"
)

// List defines a list of estimates
type List struct {
	list []Estimator
}

// NewEstimateList creates a new estimate list
func NewEstimateList(list []Estimator) (*List, error) {
	if len(list) < 1 {
		return nil, fmt.Errorf("Cannot create a empty list")
	}
	ests := new(List)

	ests.list = list

	return ests, nil
}

func (e List) Len() int {
	return len(e.list)
}

func (e List) Less(i, j int) bool {
	return e.list[i].GetEffort() > e.list[j].GetEffort()
}

func (e List) Swap(i, j int) {
	e.list[i], e.list[j] = e.list[j], e.list[i]
}
