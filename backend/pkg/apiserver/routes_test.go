package apiserver

import (
	"context"
	"encoding/json"
	"github.com/genjidb/genji"
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

func TestSessionAPIRoutes(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)

	testCases := []struct {
		description          string
		route                string
		method               string
		expectedError        bool
		expectedErrorMessage string
		expectedCode         int
	}{
		{
			description:          "Create a new Session succeeds",
			route:                "/api/sessions",
			method:               "POST",
			expectedError:        false,
			expectedErrorMessage: "",
			expectedCode:         200,
		},
		{
			description:          "Deleting a sessions fails due to wrong token length",
			route:                "/api/sessions/12345",
			method:               "DELETE",
			expectedError:        false,
			expectedErrorMessage: "Session token does not match desired length",
			expectedCode:         500,
		},
	}

	app := NewServer(&Config{
		Database: database{Location: td + "/my.db"},
		Static:   static{Prefix: "/public", Path: "../../static"},
	}, db).Start()

	for _, test := range testCases {
		req, _ := http.NewRequest(
			test.method,
			test.route,
			nil,
		)

		res, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		if test.expectedErrorMessage != "" {
			var m apiError
			decoder := json.NewDecoder(res.Body)
			derr := decoder.Decode(&m)
			assert.NoError(t, derr)
			assert.Equal(t, "error", m.Message)
			assert.Equal(t, test.expectedErrorMessage, m.Reason)
		}

		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)
	}
}
