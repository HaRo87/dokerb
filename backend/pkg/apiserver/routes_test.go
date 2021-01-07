package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/genjidb/genji"
	"github.com/haro87/dokerb/pkg/datastore"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"testing"
)

type workPackage struct {
	ID                string  `json:"id"`
	Summary           string  `json:"summary"`
	Effort            float64 `json:"effort"`
	StandardDeviation float64 `json:"standarddeviation"`
}

type estimate struct {
	WorkPackageID  string  `json:"workpackageid"`
	UserName       string  `json:"username"`
	BestCase       float64 `json:"bestcase"`
	MostLikelyCase float64 `json:"mostlikelycase"`
	WorstCase      float64 `json:"worstcase"`
}

type apiResponse struct {
	Message      string        `json:"message"`
	Reason       string        `json:"reason"`
	Route        string        `json:"route"`
	Users        []string      `json:"users"`
	Workpackages []workPackage `json:"workpackages"`
	Estimates    []estimate    `json:"estimates"`
	Hint         string        `json:"hint"`
	Estimate     Estimate      `json:"estimate"`
}

var m *datastore.MockDatastore
var ds datastore.DataStore
var db *genji.DB
var td string
var tre *regexp.Regexp

const float64CompareThreshold = 0.001

func setupTestCaseForRealDatastore(t *testing.T) func(t *testing.T) {
	td, _ = ioutil.TempDir("", "db-test")
	db, _ = genji.Open(td + "/my.db")
	db = db.WithContext(context.Background())
	ds, _ = datastore.NewGenjiDatastore(db)
	tre, _ = regexp.Compile("/sessions/([\\d|\\w]*)")

	return func(t *testing.T) {
		db.Close()
		os.RemoveAll(td)
		td = ""
		ds = nil
	}
}

func setupTestCaseForMock(t *testing.T) func(t *testing.T) {
	tre, _ = regexp.Compile("/sessions/([\\d|\\w]*)")

	m = new(datastore.MockDatastore)
	return func(t *testing.T) {
		m = nil
	}
}

func TestCreateSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("CreateSession").Return("", fmt.Errorf("Unable to create session"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	m.MethodCalled("CreateSession")

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to create session", ar.Reason)
}

func TestCreateSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("CreateSession").Return("12345678901234567890abd456789012", nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	m.MethodCalled("CreateSession")

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	token := tre.FindStringSubmatch(ar.Route)[1]
	assert.Len(t, token, 32)
	assert.Equal(t, "/sessions/"+token, ar.Route)
}

func TestDeleteSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("RemoveSession", "12345").Return(fmt.Errorf("Unable to remove session"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	m.MethodCalled("RemoveSession", "12345")

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to remove session", ar.Reason)
}

func TestDeleteSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("RemoveSession", "12345").Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	m.MethodCalled("RemoveSession", "12345")

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
}

func TestAddUserToSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("JoinSession", "12345", "Tigger").Return(fmt.Errorf("Unable to add user"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"name": "Tigger",
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/users",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to add user", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestAddUserToSessionFailsDueToMissingHeader(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"name": "Tigger",
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/users",
		bytes.NewBuffer(body),
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unprocessable Entity", ar.Reason)
	assert.Equal(t, 400, res.StatusCode)
}

func TestAddUserToSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("JoinSession", "12345", "Tigger").Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"name": "Tigger",
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/users",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, "/sessions/12345/users/Tigger", ar.Route)
	assert.Equal(t, 200, res.StatusCode)
}

func TestGetUsersFromSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetUsers", "12345").Return([]string{}, fmt.Errorf("Unable to retrieve users"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/users",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to retrieve users", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetUsersFromSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetUsers", "12345").Return([]string{"Tigger", "Rabbit"}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/users",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, []string{"Tigger", "Rabbit"}, ar.Users)
	assert.Equal(t, 200, res.StatusCode)
}

func TestRemoveUserFromSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("LeaveSession", "12345", "Tigger").Return(fmt.Errorf("Unable to remove user from session"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345/users/Tigger",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to remove user from session", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestRemoveUserFromSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("LeaveSession", "12345", "Tigger").Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345/users/Tigger",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, 200, res.StatusCode)
}

func TestGetWorkPackagesFromSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetWorkPackages", "12345").Return([]datastore.WorkPackage{}, fmt.Errorf("Unable to retrieve work packages"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/workpackages",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to retrieve work packages", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetWorkPackagesFromSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetWorkPackages", "12345").Return([]datastore.WorkPackage{datastore.WorkPackage{ID: "TEST01"}}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/workpackages",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, "TEST01", ar.Workpackages[0].ID)
	assert.Equal(t, 200, res.StatusCode)
}

func TestAddWorkPackageToSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("AddWorkPackage", "12345", "TEST01", "eat honey").Return(fmt.Errorf("Unable to add work package"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"id":      "TEST01",
		"summary": "eat honey",
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/workpackages",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to add work package", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestAddWorkPackageToSessionFailsDueToMissingHeader(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("AddWorkPackage", "12345", "TEST01", "eat honey").Return(fmt.Errorf("Unable to add work package"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"id":      "TEST01",
		"summary": "eat honey",
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/workpackages",
		bytes.NewBuffer(body),
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unprocessable Entity", ar.Reason)
	assert.Equal(t, 400, res.StatusCode)
}

func TestAddWorkPackageToSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("AddWorkPackage", "12345", "TEST01", "eat honey").Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"id":      "TEST01",
		"summary": "eat honey",
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/workpackages",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, "/sessions/12345/workpackages/TEST01", ar.Route)
	assert.Equal(t, 200, res.StatusCode)
}

func TestRemoveWorkPackageFromSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("RemoveWorkPackage", "12345", "TEST01").Return(fmt.Errorf("Unable to remove work package from session"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345/workpackages/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to remove work package from session", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestRemoveWorkPackageFromSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("RemoveWorkPackage", "12345", "TEST01").Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345/workpackages/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, 200, res.StatusCode)
}

func TestUpdateWorkPackageEstimateOfSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("AddEstimateToWorkPackage", "12345", "TEST01", 1.2, 0.2).Return(fmt.Errorf("Unable to add estimate"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"effort":            1.2,
		"standarddeviation": 0.2,
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"PUT",
		"/api/sessions/12345/workpackages/TEST01",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to add estimate", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestUpdateWorkPackageEstimateOfSessionFailsDueToMissingHeader(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"effort":            1.2,
		"standarddeviation": 0.2,
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"PUT",
		"/api/sessions/12345/workpackages/TEST01",
		bytes.NewBuffer(body),
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unprocessable Entity", ar.Reason)
	assert.Equal(t, 400, res.StatusCode)
}

func TestUpdateWorkPackageEstimateOfSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("AddEstimateToWorkPackage", "12345", "TEST01", 1.2, 0.2).Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"effort":            1.2,
		"standarddeviation": 0.2,
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"PUT",
		"/api/sessions/12345/workpackages/TEST01",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, 200, res.StatusCode)
}

func TestResetEstimateOfWorkPackageFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("RemoveEstimateFromWorkPackage", "12345", "TEST01").Return(fmt.Errorf("Unable to reset estimate"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345/workpackages/TEST01/estimate",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to reset estimate", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestResetEstimateOfWorkPackageSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("RemoveEstimateFromWorkPackage", "12345", "TEST01").Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345/workpackages/TEST01/estimate",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, 200, res.StatusCode)
}

func TestAddUserEstimateToSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("AddEstimate", "12345", datastore.Estimate{
		WorkPackageID:  "TEST01",
		UserName:       "Tigger",
		BestCase:       0.5,
		MostLikelyCase: 1.5,
		WorstCase:      3.0}).Return(fmt.Errorf("Unable to add estimate"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"id":   "TEST01",
		"user": "Tigger",
		"b":    0.5,
		"m":    1.5,
		"w":    3.0,
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/estimates",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to add estimate", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestAddUserEstimateToSessionFailsDueToMissingHeader(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"id":   "TEST01",
		"user": "Tigger",
		"b":    0.5,
		"m":    1.5,
		"w":    3.0,
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/estimates",
		bytes.NewBuffer(body),
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unprocessable Entity", ar.Reason)
	assert.Equal(t, 400, res.StatusCode)
}

func TestAddUserEstimateToSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("AddEstimate", "12345", datastore.Estimate{
		WorkPackageID:  "TEST01",
		UserName:       "Tigger",
		BestCase:       0.5,
		MostLikelyCase: 1.5,
		WorstCase:      3.0}).Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	payloadf := map[string]interface{}{
		"id":   "TEST01",
		"user": "Tigger",
		"b":    0.5,
		"m":    1.5,
		"w":    3.0,
	}
	body, me := json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions/12345/estimates",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, "/sessions/12345/estimates/Tigger/TEST01", ar.Route)
	assert.Equal(t, 200, res.StatusCode)
}

func TestRemoveUserEstimateFromSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("RemoveEstimate", "12345", datastore.Estimate{
		WorkPackageID: "TEST01",
		UserName:      "Tigger"}).Return(fmt.Errorf("Unable to remove estimate"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345/estimates/Tigger/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to remove estimate", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestRemoveUserEstimateFromSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("RemoveEstimate", "12345", datastore.Estimate{
		WorkPackageID: "TEST01",
		UserName:      "Tigger"}).Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"DELETE",
		"/api/sessions/12345/estimates/Tigger/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, 200, res.StatusCode)
}

func TestGetUserEstimatesFromSessionFails(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{}, fmt.Errorf("Unable to retrieve estimates"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to retrieve estimates", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetUserEstimatesFromSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{datastore.Estimate{WorkPackageID: "TEST01"}}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, "TEST01", ar.Estimates[0].WorkPackageID)
	assert.Equal(t, 200, res.StatusCode)
}

func TestGetAverageEstimateForWorkPackageFromSessionFailsDueToErrorOnGetEstimates(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{}, fmt.Errorf("Unable to retrieve estimates"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to retrieve estimates", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetAverageEstimateForWorkPackageFromSessionFailsDueToNoEstimates(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Not enough data to process", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}
func TestGetAverageEstimateForWorkPackageFromSessionFailsDueToErrorOnGetUsers(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{datastore.Estimate{
		WorkPackageID: "TEST01",
		UserName:      "Tigger",
	}}, nil)

	m.On("GetUsers", "12345").Return([]string{}, fmt.Errorf("Unable to retrieve users"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to retrieve users", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetAverageEstimateForWorkPackageFromSessionFailsDueToInvalidEstimate(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{datastore.Estimate{
		WorkPackageID:  "TEST01",
		UserName:       "Tigger",
		BestCase:       0.5,
		MostLikelyCase: 0.2,
		WorstCase:      1.5,
	}}, nil)

	m.On("GetUsers", "12345").Return([]string{"Tigger"}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Most Likely was smaller than Best Effort", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetAverageEstimateForWorkPackageFromSessionSuccessWithAllUsersProvidedEstimates(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{datastore.Estimate{
		WorkPackageID:  "TEST01",
		UserName:       "Tigger",
		BestCase:       1.0,
		MostLikelyCase: 2.0,
		WorstCase:      4.0,
	},
		{
			WorkPackageID:  "TEST01",
			UserName:       "Rabbit",
			BestCase:       2.0,
			MostLikelyCase: 3.0,
			WorstCase:      5.0,
		},
	}, nil)

	m.On("GetUsers", "12345").Return([]string{"Tigger", "Rabbit"}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, []string{}, ar.Users)
	assert.True(t, math.Abs(2.666-ar.Estimate.Effort) <= float64CompareThreshold)
	assert.True(t, math.Abs(0.5-ar.Estimate.StandardDeviation) <= float64CompareThreshold)
	assert.Equal(t, 200, res.StatusCode)
}

func TestGetAverageEstimateForWorkPackageFromSessionSuccessWithNotAllUsersProvidedEstimates(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{datastore.Estimate{
		WorkPackageID:  "TEST01",
		UserName:       "Tigger",
		BestCase:       1.0,
		MostLikelyCase: 2.0,
		WorstCase:      4.0,
	},
		{
			WorkPackageID:  "TEST01",
			UserName:       "Rabbit",
			BestCase:       2.0,
			MostLikelyCase: 3.0,
			WorstCase:      5.0,
		},
	}, nil)

	m.On("GetUsers", "12345").Return([]string{"Tigger", "Rabbit", "Piglet"}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "warning", ar.Message)
	assert.Equal(t, "not all users did provide estimates", ar.Hint)
	assert.Equal(t, []string{"Piglet"}, ar.Users)
	assert.True(t, math.Abs(2.666-ar.Estimate.Effort) <= float64CompareThreshold)
	assert.True(t, math.Abs(0.5-ar.Estimate.StandardDeviation) <= float64CompareThreshold)
	assert.Equal(t, 200, res.StatusCode)
}

func TestGetUserWithMaxEstimateDistanceForWorkPackageFromSessionFailsDueToErrorOnGetEstimates(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{}, fmt.Errorf("Unable to retrieve estimates"))

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01/users/distance",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Unable to retrieve estimates", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetUserWithMaxEstimateDistanceForWorkPackageFromSessionFailsDueToNoEstimates(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01/users/distance",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Not enough data to process", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetUserWithMaxEstimateDistanceForWorkPackageFromSessionFailsDueToInvalidEstimate(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{datastore.Estimate{
		WorkPackageID:  "TEST01",
		UserName:       "Tigger",
		BestCase:       0.5,
		MostLikelyCase: 0.2,
		WorstCase:      1.5,
	}}, nil)

	m.On("GetUsers", "12345").Return([]string{"Tigger"}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01/users/distance",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "error", ar.Message)
	assert.Equal(t, "Most Likely was smaller than Best Effort", ar.Reason)
	assert.Equal(t, 500, res.StatusCode)
}

func TestGetUserWithMaxEstimateDistanceForWorkPackageFromSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	m.On("GetEstimates", "12345").Return([]datastore.Estimate{datastore.Estimate{
		WorkPackageID:  "TEST01",
		UserName:       "Tigger",
		BestCase:       1.0,
		MostLikelyCase: 2.0,
		WorstCase:      4.0,
	},
		{
			WorkPackageID:  "TEST01",
			UserName:       "Rabbit",
			BestCase:       2.0,
			MostLikelyCase: 3.0,
			WorstCase:      5.0,
		},
		{
			WorkPackageID:  "TEST01",
			UserName:       "Piglet",
			BestCase:       5.0,
			MostLikelyCase: 6.0,
			WorstCase:      7.0,
		},
	}, nil)

	m.On("GetUsers", "12345").Return([]string{"Tigger", "Rabbit", "Piglet"}, nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	req, _ := http.NewRequest(
		"GET",
		"/api/sessions/12345/estimates/TEST01/users/distance",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, []string{"Piglet", "Tigger"}, ar.Users)
	assert.Equal(t, 200, res.StatusCode)
}

func TestSmokeWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDatastore(t)
	defer setupAndTearDown(t)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, ds).Start()

	req, _ := http.NewRequest(
		"POST",
		"/api/sessions",
		nil,
	)

	res, err := app.Test(req, -1)

	assert.NoError(t, err)

	var ar apiResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	token := tre.FindStringSubmatch(ar.Route)[1]

	payload := map[string]string{
		"id": "TEST01",
	}
	body, me := json.Marshal(payload)

	assert.NoError(t, me)

	req, _ = http.NewRequest(
		"POST",
		"/api/sessions/"+token+"/workpackages",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)

	payload = map[string]string{
		"id":      "TEST02",
		"summary": "some test",
	}
	body, me = json.Marshal(payload)

	assert.NoError(t, me)

	req, _ = http.NewRequest(
		"POST",
		"/api/sessions/"+token+"/workpackages",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)

	payload = map[string]string{
		"name": "Tigger",
	}
	body, me = json.Marshal(payload)

	assert.NoError(t, me)

	req, _ = http.NewRequest(
		"POST",
		"/api/sessions/"+token+"/users",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)

	payload = map[string]string{
		"name": "Rabbit",
	}
	body, me = json.Marshal(payload)

	assert.NoError(t, me)

	req, _ = http.NewRequest(
		"POST",
		"/api/sessions/"+token+"/users",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)

	payloadf := map[string]interface{}{
		"id":   "TEST01",
		"user": "Tigger",
		"b":    0.5,
		"m":    1.0,
		"w":    2.0,
	}
	body, me = json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ = http.NewRequest(
		"POST",
		"/api/sessions/"+token+"/estimates",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)

	payloadf = map[string]interface{}{
		"id":   "TEST02",
		"user": "Tigger",
		"b":    0.2,
		"m":    1.2,
		"w":    1.5,
	}
	body, me = json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ = http.NewRequest(
		"POST",
		"/api/sessions/"+token+"/estimates",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)

	payloadf = map[string]interface{}{
		"id":   "TEST01",
		"user": "Rabbit",
		"b":    1.0,
		"m":    1.2,
		"w":    2.0,
	}
	body, me = json.Marshal(payloadf)

	assert.NoError(t, me)

	req, _ = http.NewRequest(
		"POST",
		"/api/sessions/"+token+"/estimates",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)

	req, _ = http.NewRequest(
		"GET",
		"/api/sessions/"+token+"/estimates",
		nil,
	)

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, "TEST01", ar.Estimates[0].WorkPackageID)
	assert.Equal(t, "Tigger", ar.Estimates[0].UserName)
	assert.Equal(t, 0.5, ar.Estimates[0].BestCase)
	assert.Equal(t, 1.0, ar.Estimates[0].MostLikelyCase)
	assert.Equal(t, 2.0, ar.Estimates[0].WorstCase)
	assert.Equal(t, "TEST02", ar.Estimates[1].WorkPackageID)
	assert.Equal(t, "Tigger", ar.Estimates[1].UserName)
	assert.Equal(t, "TEST01", ar.Estimates[2].WorkPackageID)
	assert.Equal(t, "Rabbit", ar.Estimates[2].UserName)
	assert.Len(t, ar.Estimates, 3)

	req, _ = http.NewRequest(
		"DELETE",
		"/api/sessions/"+token+"/estimates/Rabbit/TEST01",
		nil,
	)

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)

	req, _ = http.NewRequest(
		"GET",
		"/api/sessions/"+token+"/estimates",
		nil,
	)

	res, err = app.Test(req, -1)

	assert.NoError(t, err)

	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&ar)
	assert.NoError(t, err)
	assert.Equal(t, "ok", ar.Message)
	assert.Equal(t, "TEST01", ar.Estimates[0].WorkPackageID)
	assert.Equal(t, "Tigger", ar.Estimates[0].UserName)
	assert.Equal(t, 0.5, ar.Estimates[0].BestCase)
	assert.Equal(t, 1.0, ar.Estimates[0].MostLikelyCase)
	assert.Equal(t, 2.0, ar.Estimates[0].WorstCase)
	assert.Equal(t, "TEST02", ar.Estimates[1].WorkPackageID)
	assert.Equal(t, "Tigger", ar.Estimates[1].UserName)
	assert.Len(t, ar.Estimates, 2)
}
