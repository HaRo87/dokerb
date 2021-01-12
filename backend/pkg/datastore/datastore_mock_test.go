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

func TestAddTaskNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddTask", "12345", "TEST01", "").Return(nil)

	err := ds.AddTask("12345", "TEST01", "")

	assert.NoError(t, err)
	m.MethodCalled("AddTask", "12345", "TEST01", "")
}

func TestAddTaskError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddTask", "12345", "TEST01", "").Return(fmt.Errorf("Some error"))

	err := ds.AddTask("12345", "TEST01", "")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("AddTask", "12345", "TEST01", "")
}

func TestRemoveTaskNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveTask", "12345", "TEST01").Return(nil)

	err := ds.RemoveTask("12345", "TEST01")

	assert.NoError(t, err)
	m.MethodCalled("RemoveTask", "12345", "TEST01")
}

func TestRemoveTaskError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveTask", "12345", "TEST01").Return(fmt.Errorf("Some error"))

	err := ds.RemoveTask("12345", "TEST01")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("RemoveTask", "12345", "TEST01")
}

func TestAddEstimateToTaskNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddEstimateToTask", "12345", "TEST01", 0.1, 0.2).Return(nil)

	err := ds.AddEstimateToTask("12345", "TEST01", 0.1, 0.2)

	assert.NoError(t, err)
	m.MethodCalled("AddEstimateToTask", "12345", "TEST01", 0.1, 0.2)
}

func TestAddEstimateToTaskError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddEstimateToTask", "12345", "TEST01", 0.1, 0.2).Return(fmt.Errorf("Some error"))

	err := ds.AddEstimateToTask("12345", "TEST01", 0.1, 0.2)

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("AddEstimateToTask", "12345", "TEST01", 0.1, 0.2)
}

func RemoveEstimateFromTask(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveEstimateFromTask", "12345", "TEST01").Return(nil)

	err := ds.RemoveEstimateFromTask("12345", "TEST01")

	assert.NoError(t, err)
	m.MethodCalled("RemoveEstimateFromTask", "12345", "TEST01")
}

func TestRemoveEstimateFromTaskError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveEstimateFromTask", "12345", "TEST01").Return(fmt.Errorf("Some error"))

	err := ds.RemoveEstimateFromTask("12345", "TEST01")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("RemoveEstimateFromTask", "12345", "TEST01")
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

func TestGetTasksNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetTasks", "12345").Return([]Task{Task{ID: "TEST01"}}, nil)

	res, err := ds.GetTasks("12345")

	assert.NoError(t, err)
	assert.Equal(t, "TEST01", res[0].ID)
	m.MethodCalled("GetTasks", "12345")
}

func TestGetTasksError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetTasks", "12345").Return([]Task{}, fmt.Errorf("Some error"))

	_, err := ds.GetTasks("12345")

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("GetTasks", "12345")
}

func TestAddEstimateNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddEstimate", "12345", Estimate{TaskID: "TEST01"}).Return(nil)

	err := ds.AddEstimate("12345", Estimate{TaskID: "TEST01"})

	assert.NoError(t, err)
	m.MethodCalled("AddEstimate", "12345", Estimate{TaskID: "TEST01"})
}

func TestAddEstimateError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("AddEstimate", "12345", Estimate{TaskID: "TEST01"}).Return(fmt.Errorf("Some error"))

	err := ds.AddEstimate("12345", Estimate{TaskID: "TEST01"})

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("AddEstimate", "12345", Estimate{TaskID: "TEST01"})
}

func TestRemoveEstimateNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveEstimate", "12345", Estimate{TaskID: "TEST01"}).Return(nil)

	err := ds.RemoveEstimate("12345", Estimate{TaskID: "TEST01"})

	assert.NoError(t, err)
	m.MethodCalled("RemoveEstimate", "12345", Estimate{TaskID: "TEST01"})
}

func TestRemoveEstimateError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("RemoveEstimate", "12345", Estimate{TaskID: "TEST01"}).Return(fmt.Errorf("Some error"))

	err := ds.RemoveEstimate("12345", Estimate{TaskID: "TEST01"})

	assert.Error(t, err)
	assert.Equal(t, "Some error", err.Error())
	m.MethodCalled("RemoveEstimate", "12345", Estimate{TaskID: "TEST01"})
}

func TestGetEstimatesNoError(t *testing.T) {
	var ds DataStore
	m := new(MockDatastore)
	ds = m

	m.On("GetEstimates", "12345").Return([]Estimate{Estimate{TaskID: "TEST01"}}, nil)

	res, err := ds.GetEstimates("12345")

	assert.NoError(t, err)
	assert.Equal(t, "TEST01", res[0].TaskID)
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
