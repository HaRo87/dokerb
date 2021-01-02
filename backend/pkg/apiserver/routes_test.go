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
	"net/http"
	"os"
	"testing"
)

type apiError struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type mock struct {
	Method      string
	Call        string
	ReturnValue interface{}
}

type testCase struct {
	Description          string
	Route                string
	Method               string
	Body                 map[string]string
	ExpectedError        bool
	ExpectedErrorMessage string
	ExpectedCode         int
	Mock                 []mock
}

var m *datastore.MockGenjiDB
var db *genji.DB
var td string

func setupTestCaseForRealDB(t *testing.T) func(t *testing.T) {
	td, _ = ioutil.TempDir("", "db-test")
	db, _ = genji.Open(td + "/my.db")
	db = db.WithContext(context.Background())

	return func(t *testing.T) {
		db.Close()
		os.RemoveAll(td)
		td = ""
	}
}

func setupTestCaseForMock(t *testing.T) func(t *testing.T) {
	m = new(datastore.MockGenjiDB)
	return func(t *testing.T) {
		m = nil
	}
}

func TestAPIRoutesForErrors(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)

	testCases := []testCase{
		{
			Description:          "Create a new Session fails",
			Route:                "/api/sessions",
			Method:               "POST",
			ExpectedError:        false,
			ExpectedErrorMessage: "Unable to store session token",
			ExpectedCode:         500,
			Mock: []mock{
				mock{
					Method:      "Exec",
					Call:        "INSERT INTO sessions VALUES ?",
					ReturnValue: fmt.Errorf("Ooops, something went wrong"),
				},
			},
		},
		{
			Description:          "Deleting a sessions fails due to wrong token length",
			Route:                "/api/sessions/12345",
			Method:               "DELETE",
			ExpectedError:        false,
			ExpectedErrorMessage: "Session token does not match desired length",
			ExpectedCode:         500,
		},
		{
			Description:          "Adding a user to a session fails due to wrong token length",
			Route:                "/api/sessions/12345/users/Tigger",
			Method:               "POST",
			ExpectedError:        false,
			ExpectedErrorMessage: "Session token does not match desired length",
			ExpectedCode:         500,
		},
		{
			Description:          "Getting users fails due to wrong token length",
			Route:                "/api/sessions/12345/users",
			Method:               "GET",
			ExpectedError:        false,
			ExpectedErrorMessage: "Session token does not match desired length",
			ExpectedCode:         500,
		},
		{
			Description:          "Deleting a user fails due to wrong token length",
			Route:                "/api/sessions/12345/users/Tigger",
			Method:               "DELETE",
			ExpectedError:        false,
			ExpectedErrorMessage: "Session token does not match desired length",
			ExpectedCode:         500,
		},
		{
			Description:          "Getting work packages fails due to wrong token length",
			Route:                "/api/sessions/12345/workpackages",
			Method:               "GET",
			ExpectedError:        false,
			ExpectedErrorMessage: "Session token does not match desired length",
			ExpectedCode:         500,
		},
		{
			Description: "Adding a work package to a session fails due to wrong token length",
			Route:       "/api/sessions/12345/workpackages",
			Method:      "POST",
			Body: map[string]string{
				"id": "TEST01",
			},
			ExpectedError:        false,
			ExpectedErrorMessage: "Session token does not match desired length",
			ExpectedCode:         500,
		},
		{
			Description:          "Deleting a work package fails due to wrong token length",
			Route:                "/api/sessions/12345/workpackages/TEST01",
			Method:               "DELETE",
			ExpectedError:        false,
			ExpectedErrorMessage: "Session token does not match desired length",
			ExpectedCode:         500,
		},
	}

	m.On("Exec", "CREATE TABLE sessions").Return(nil)

	app := NewServer(&Config{
		Static: static{Prefix: "/public", Path: "../../static"},
	}, m).Start()

	for _, test := range testCases {
		var req *http.Request

		if len(test.Body) > 0 {
			body, me := json.Marshal(test.Body)
			assert.NoError(t, me)
			req, _ = http.NewRequest(
				test.Method,
				test.Route,
				bytes.NewBuffer(body),
			)
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, _ = http.NewRequest(
				test.Method,
				test.Route,
				nil,
			)
		}

		if len(test.Mock) > 0 {
			for _, moc := range test.Mock {
				m.On(moc.Method, moc.Call).Return(moc.ReturnValue)
			}
		}

		res, err := app.Test(req, -1)

		assert.Equalf(t, test.ExpectedError, err != nil, test.Description)

		if test.ExpectedErrorMessage != "" {
			var ae apiError
			decoder := json.NewDecoder(res.Body)
			derr := decoder.Decode(&ae)
			assert.NoError(t, derr)
			assert.Equal(t, "error", ae.Message)
			assert.Equal(t, test.ExpectedErrorMessage, ae.Reason)
		}

		if test.ExpectedError {
			continue
		}

		assert.Equalf(t, test.ExpectedCode, res.StatusCode, test.Description)
	}
}
