package datastore

import (
	"github.com/stretchr/testify/mock"
)

// MockDatastore represents the mocked object
type MockDatastore struct {
	mock.Mock
}

// CreateSession implements the Datastore interface
func (m *MockDatastore) CreateSession() (string, error) {
	arguments := m.Called()
	return arguments.Get(0).(string), arguments.Error(1)
}

// JoinSession implements the Datastore interface
func (m *MockDatastore) JoinSession(t, n string) error {
	arguments := m.Called(t, n)
	return arguments.Error(0)
}

// LeaveSession implements the Datastore interface
func (m *MockDatastore) LeaveSession(t, n string) error {
	arguments := m.Called(t, n)
	return arguments.Error(0)
}

// RemoveSession implements the Datastore interface
func (m *MockDatastore) RemoveSession(t string) error {
	arguments := m.Called(t)
	return arguments.Error(0)
}

// AddWorkPackage implements the Datastore interface
func (m *MockDatastore) AddWorkPackage(t, id, s string) error {
	arguments := m.Called(t, id, s)
	return arguments.Error(0)
}

// RemoveWorkPackage implements the Datastore interface
func (m *MockDatastore) RemoveWorkPackage(t, id string) error {
	arguments := m.Called(t, id)
	return arguments.Error(0)
}

// AddEstimateToWorkPackage implements the Datastore interface
func (m *MockDatastore) AddEstimateToWorkPackage(t, id string, e, s float64) error {
	arguments := m.Called(t, id, e, s)
	return arguments.Error(0)
}

// RemoveEstimateFromWorkPackage implements the Datastore interface
func (m *MockDatastore) RemoveEstimateFromWorkPackage(t, id string) error {
	arguments := m.Called(t, id)
	return arguments.Error(0)
}

// GetUsers implements the Datastore interface
func (m *MockDatastore) GetUsers(t string) ([]string, error) {
	arguments := m.Called(t)
	return arguments.Get(0).([]string), arguments.Error(1)
}

// GetWorkPackages implements the Datastore interface
func (m *MockDatastore) GetWorkPackages(t string) ([]WorkPackage, error) {
	arguments := m.Called(t)
	return arguments.Get(0).([]WorkPackage), arguments.Error(1)
}

// AddEstimate implements the Datastore interface
func (m *MockDatastore) AddEstimate(t string, e Estimate) error {
	arguments := m.Called(t, e)
	return arguments.Error(0)
}

// RemoveEstimate implements the Datastore interface
func (m *MockDatastore) RemoveEstimate(t string, e Estimate) error {
	arguments := m.Called(t, e)
	return arguments.Error(0)
}

// GetEstimates implements the Datastore interface
func (m *MockDatastore) GetEstimates(t string) ([]Estimate, error) {
	arguments := m.Called(t)
	return arguments.Get(0).([]Estimate), arguments.Error(1)
}
