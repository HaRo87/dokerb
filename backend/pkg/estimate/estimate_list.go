package estimate

import (
	"fmt"
)

// UserEstimate represents a per user estimate
type UserEstimate struct {
	Name     string
	Estimate Estimator
}

// List defines a list of estimates
type List struct {
	list []UserEstimate
}

// NewEstimateList creates a new estimate list
func NewEstimateList(list []UserEstimate) (*List, error) {
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
	return e.list[i].Estimate.GetEffort() > e.list[j].Estimate.GetEffort()
}

func (e List) Swap(i, j int) {
	e.list[i], e.list[j] = e.list[j], e.list[i]
}

// GetFirstUser returns the first user in the list
func (e List) GetFirstUser() string {
	return e.list[0].Name
}

// GetLastUser returns the last user in the list if any
func (e List) GetLastUser() string {
	if len(e.list) > 1 {
		return e.list[len(e.list)-1].Name
	}
	return ""
}
