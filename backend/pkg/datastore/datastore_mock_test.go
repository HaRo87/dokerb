package datastore

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateSessionNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("CreateSession").Return("12345", nil)

	res, err := ds.CreateSession()

	assert.NoError(t, err)
	assert.Equal(t, "12345", res)
	m.MethodCalled("CreateSession")
}

func TestCreateSessionError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("CreateSession").Return("", fmt.Errorf("Some error"))

	_, err := ds.CreateSession()

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("CreateSession")
}

func TestJoinSessionNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("JoinSession", "12345", "Tigger").Return(nil)

	err := ds.JoinSession("12345", "Tigger")

	assert.NoError(t, err)
	m.MethodCalled("JoinSession", "12345", "Tigger")
}

func TestJoinSessionError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("JoinSession", "12345", "Tigger").Return(fmt.Errorf("Some error"))

	err := ds.JoinSession("12345", "Tigger")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("JoinSession", "12345", "Tigger")
}

func TestLeaveSessionNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("LeaveSession", "12345", "Tigger").Return(nil)

	err := ds.LeaveSession("12345", "Tigger")

	assert.NoError(t, err)
	m.MethodCalled("LeaveSession", "12345", "Tigger")
}

func TestLeaveSessionError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("LeaveSession", "12345", "Tigger").Return(fmt.Errorf("Some error"))

	err := ds.LeaveSession("12345", "Tigger")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("LeaveSession", "12345", "Tigger")
}

func TestRemoveSessionNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveSession", "12345").Return(nil)

	err := ds.RemoveSession("12345")

	assert.NoError(t, err)
	m.MethodCalled("RemoveSession", "12345")
}

func TestRemoveSessionError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveSession", "12345").Return(fmt.Errorf("Some error"))

	err := ds.RemoveSession("12345")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("RemoveSession", "12345")
}

func TestAddWorkPackageNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddWorkPackage", "12345", "TEST01", "").Return(nil)

	err := ds.AddWorkPackage("12345", "TEST01", "")

	assert.NoError(t, err)
	m.MethodCalled("AddWorkPackage", "12345", "TEST01", "")
}

func TestAddWorkPackageError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddWorkPackage", "12345", "TEST01", "").Return(fmt.Errorf("Some error"))

	err := ds.AddWorkPackage("12345", "TEST01", "")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("AddWorkPackage", "12345", "TEST01", "")
}

func TestRemoveWorkPackageNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveWorkPackage", "12345", "TEST01").Return(nil)

	err := ds.RemoveWorkPackage("12345", "TEST01")

	assert.NoError(t, err)
	m.MethodCalled("RemoveWorkPackage", "12345", "TEST01")
}

func TestRemoveWorkPackageError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveWorkPackage", "12345", "TEST01").Return(fmt.Errorf("Some error"))

	err := ds.RemoveWorkPackage("12345", "TEST01")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("RemoveWorkPackage", "12345", "TEST01")
}

func TestAddEstimateToWorkPackageNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddEstimateToWorkPackage", "12345", "TEST01", 0.1, 0.2).Return(nil)

	err := ds.AddEstimateToWorkPackage("12345", "TEST01", 0.1, 0.2)

	assert.NoError(t, err)
	m.MethodCalled("AddEstimateToWorkPackage", "12345", "TEST01", 0.1, 0.2)
}

func TestAddEstimateToWorkPackageError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddEstimateToWorkPackage", "12345", "TEST01", 0.1, 0.2).Return(fmt.Errorf("Some error"))

	err := ds.AddEstimateToWorkPackage("12345", "TEST01", 0.1, 0.2)

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("AddEstimateToWorkPackage", "12345", "TEST01", 0.1, 0.2)
}

func RemoveEstimateFromWorkPackage(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveEstimateFromWorkPackage", "12345", "TEST01").Return(nil)

	err := ds.RemoveEstimateFromWorkPackage("12345", "TEST01")

	assert.NoError(t, err)
	m.MethodCalled("RemoveEstimateFromWorkPackage", "12345", "TEST01")
}

func TestRemoveEstimateFromWorkPackageError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveEstimateFromWorkPackage", "12345", "TEST01").Return(fmt.Errorf("Some error"))

	err := ds.RemoveEstimateFromWorkPackage("12345", "TEST01")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("RemoveEstimateFromWorkPackage", "12345", "TEST01")
}

func TestGetUsersNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetUsers", "12345").Return([]string{"Tigger"}, nil)

	res, err := ds.GetUsers("12345")

	assert.NoError(t, err)
	assert.Equal(t, []string{"Tigger"}, res)
	m.MethodCalled("GetUsers", "12345")
}

func TestGetUsersError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetUsers", "12345").Return([]string{}, fmt.Errorf("Some error"))

	_, err := ds.GetUsers("12345")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("GetUsers", "12345")
}

func TestGetWorkPackagesNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetWorkPackages", "12345").Return([]WorkPackage{WorkPackage{ID: "TEST01"}}, nil)

	res, err := ds.GetWorkPackages("12345")

	assert.NoError(t, err)
	assert.Equal(t, "TEST01", res[0].ID)
	m.MethodCalled("GetWorkPackages", "12345")
}

func TestGetWorkPackagesError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetWorkPackages", "12345").Return([]WorkPackage{}, fmt.Errorf("Some error"))

	_, err := ds.GetWorkPackages("12345")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("GetWorkPackages", "12345")
}

func TestAddEstimateNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddEstimate", "12345", Estimate{WorkPackageID: "TEST01"}).Return(nil)

	err := ds.AddEstimate("12345", Estimate{WorkPackageID: "TEST01"})

	assert.NoError(t, err)
	m.MethodCalled("AddEstimate", "12345", Estimate{WorkPackageID: "TEST01"})
}

func TestAddEstimateError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddEstimate", "12345", Estimate{WorkPackageID: "TEST01"}).Return(fmt.Errorf("Some error"))

	err := ds.AddEstimate("12345", Estimate{WorkPackageID: "TEST01"})

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("AddEstimate", "12345", Estimate{WorkPackageID: "TEST01"})
}

func TestRemoveEstimateNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveEstimate", "12345", Estimate{WorkPackageID: "TEST01"}).Return(nil)

	err := ds.RemoveEstimate("12345", Estimate{WorkPackageID: "TEST01"})

	assert.NoError(t, err)
	m.MethodCalled("RemoveEstimate", "12345", Estimate{WorkPackageID: "TEST01"})
}

func TestRemoveEstimateError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveEstimate", "12345", Estimate{WorkPackageID: "TEST01"}).Return(fmt.Errorf("Some error"))

	err := ds.RemoveEstimate("12345", Estimate{WorkPackageID: "TEST01"})

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("RemoveEstimate", "12345", Estimate{WorkPackageID: "TEST01"})
}

func TestGetEstimatesNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetEstimates", "12345").Return([]Estimate{Estimate{WorkPackageID: "TEST01"}}, nil)

	res, err := ds.GetEstimates("12345")

	assert.NoError(t, err)
	assert.Equal(t, "TEST01", res[0].WorkPackageID)
	m.MethodCalled("GetEstimates", "12345")
}

func TestGetEstimatesError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetEstimates", "12345").Return([]Estimate{}, fmt.Errorf("Some error"))

	_, err := ds.GetEstimates("12345")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("GetEstimates", "12345")
}
