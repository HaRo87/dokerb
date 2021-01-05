package datastore

import (
	"context"
	"fmt"
	"github.com/genjidb/genji"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

var m *MockGenjiDB
var db *genji.DB
var td string

func setupTestCaseForMock(t *testing.T) func(t *testing.T) {
	m = new(MockGenjiDB)
	return func(t *testing.T) {
		m = nil
	}
}

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

func TestNilDB(t *testing.T) {
	_, err := NewGenjiDatastore(nil)
	assert.Equal(t, "Proper DB must be provided and not nil", err.Error())
}

func TestErrorOnDBSetupCreateTable(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(fmt.Errorf("Ooops, something went wrong"))
	_, err := NewGenjiDatastore(m)
	assert.Equal(t, "Unable to create sessions table", err.Error())

	m.MethodCalled("Exec", "CREATE TABLE sessions")
}

func TestSingletonPatternWorks(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	g1, err := NewGenjiDatastore(m)
	g2, err2 := NewGenjiDatastore(m)
	assert.NoError(t, err)
	assert.NoError(t, err2)

	assert.Equal(t, g1, g2)
}

func TestGenerateTokenWrongLength(t *testing.T) {
	_, err := generateToken(-20)
	assert.Equal(t, "Invalid token length provided: -20, should be >= 20", err.Error())
}

func TestGenerateTokenWithMinimalSuggestedLength(t *testing.T) {
	token, err := generateToken(20)
	assert.NoError(t, err)
	assert.Equal(t, 20, len(token))
}

func TestGenerateTokenWithProperLength(t *testing.T) {
	token, err := generateToken(32)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(token))
}

func TestCorrectDBSetup(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	_, err := NewGenjiDatastore(m)
	assert.NoError(t, err)

	m.MethodCalled("Exec", "CREATE TABLE sessions")
}

func TestCorrectDBSetupWithGenji(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	_, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
}

func TestCreateSessionFailsDueToExecError(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	ds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	m.On("Exec", "INSERT INTO sessions VALUES ?").Return(fmt.Errorf("Ooops, something went wrong"))
	token, err2 := ds.CreateSession()
	assert.Empty(t, token)
	assert.Equal(t, "Unable to store session token", err2.Error())
	m.MethodCalled("Exec", "INSERT INTO sessions VALUES ?")
}

func TestCreateSessionSuccess(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	ds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	m.On("Exec", "INSERT INTO sessions VALUES ?").Return(nil)
	token, err2 := ds.CreateSession()
	assert.Equal(t, 32, len(token))
	assert.NoError(t, err2)
	m.MethodCalled("Exec", "INSERT INTO sessions VALUES ?")
}

func TestCreateSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	assert.Equal(t, 32, len(token))
}

func TestJoinSessionFailsDueToEmptyName(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.JoinSession("12345678901234567890123456789012", "")
	assert.Equal(t, "User name should not be empty", err2.Error())
}

func TestJoinSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.JoinSession("1234567890123456789012345678901212", "")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestJoinSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession("12345678901234567890123456789012", "Tigger")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestJoinSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
}

func TestJoinSessionErrorWhileTryingToAddUserTwiceWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
	err4 := gds.JoinSession(token, "Tigger")
	assert.Equal(t, "User with name: Tigger already part of session", err4.Error())
}

func TestRemoveUserFromEmptyList(t *testing.T) {
	l, err := removeUser([]string{}, "Tigger")
	assert.Equal(t, "User with name: Tigger is not part of session", err.Error())
	assert.Len(t, l, 0)
}

func TestRemoveUserFromListWithoutThatUserBeingPartOfThatList(t *testing.T) {
	users := []string{"Tigger", "Rabbit", "Piglet"}
	l, err := removeUser(users, "Winnie-the-Pooh")
	assert.Equal(t, "User with name: Winnie-the-Pooh is not part of session", err.Error())
	assert.Len(t, l, 3)
}

func TestRemoveUserSuccess(t *testing.T) {
	users := []string{"Tigger", "Rabbit", "Piglet"}
	l, err := removeUser(users, "Tigger")
	assert.NoError(t, err)
	assert.Len(t, l, 2)
	assert.NotContains(t, l, "Tigger")
}

func TestRemoveWorkPackageFromEmptyList(t *testing.T) {
	l, err := removeWorkpackage([]WorkPackage{}, "TEST01")
	assert.Equal(t, "Workpackage with ID: TEST01 is not part of session", err.Error())
	assert.Len(t, l, 0)
}

func TestRemoveWorkPackageFromListWithoutThatWorkPackageBeingPartOfThatList(t *testing.T) {
	wps := []WorkPackage{
		WorkPackage{
			ID: "TEST01",
		},
		WorkPackage{
			ID: "TEST02",
		},
		WorkPackage{
			ID: "TEST03",
		},
	}
	l, err := removeWorkpackage(wps, "TEST04")
	assert.Equal(t, "Workpackage with ID: TEST04 is not part of session", err.Error())
	assert.Len(t, l, 3)
}

func TestRemoveWorkPackageSuccess(t *testing.T) {
	wps := []WorkPackage{
		WorkPackage{
			ID: "TEST01",
		},
		WorkPackage{
			ID: "TEST02",
		},
		WorkPackage{
			ID: "TEST03",
		},
	}
	l, err := removeWorkpackage(wps, "TEST03")
	assert.NoError(t, err)
	assert.Len(t, l, 2)
	assert.NotContains(t, l, WorkPackage{ID: "TEST03"})
}

func TestRemoveEstimateFromEmptyList(t *testing.T) {
	l, err := removeEstimate([]Estimate{}, Estimate{WorkPackageID: "TEST01", UserName: "Tigger"})
	assert.Equal(t, "Estimate with ID: TEST01 and user name: Tigger is not part of session", err.Error())
	assert.Len(t, l, 0)
}

func TestRemoveEstimateFromListWithoutThatEstimateBeingPartOfThatListDueToIDAndUserName(t *testing.T) {
	est := []Estimate{
		Estimate{
			WorkPackageID: "TEST01",
			UserName:      "Tigger",
		},
		Estimate{
			WorkPackageID: "TEST02",
			UserName:      "Rabbit",
		},
		Estimate{
			WorkPackageID: "TEST03",
			UserName:      "Piglet",
		},
	}
	l, err := removeEstimate(est, Estimate{WorkPackageID: "TEST04", UserName: "Tigger"})
	assert.Equal(t, "Estimate with ID: TEST04 and user name: Tigger is not part of session", err.Error())
	assert.Len(t, l, 3)
}

func TestRemoveEstimateFromListWithoutThatEstimateBeingPartOfThatListDueToUserName(t *testing.T) {
	est := []Estimate{
		Estimate{
			WorkPackageID: "TEST01",
			UserName:      "Tigger",
		},
		Estimate{
			WorkPackageID: "TEST02",
			UserName:      "Rabbit",
		},
		Estimate{
			WorkPackageID: "TEST03",
			UserName:      "Piglet",
		},
	}
	l, err := removeEstimate(est, Estimate{WorkPackageID: "TEST01", UserName: "Piglet"})
	assert.Equal(t, "Estimate with ID: TEST01 and user name: Piglet is not part of session", err.Error())
	assert.Len(t, l, 3)
}

func TestRemoveEstimateSuccess(t *testing.T) {
	est := []Estimate{
		Estimate{
			WorkPackageID: "TEST01",
			UserName:      "Tigger",
		},
		Estimate{
			WorkPackageID: "TEST02",
			UserName:      "Rabbit",
		},
		Estimate{
			WorkPackageID: "TEST03",
			UserName:      "Piglet",
		},
	}
	l, err := removeEstimate(est, Estimate{WorkPackageID: "TEST01", UserName: "Tigger"})
	assert.NoError(t, err)
	assert.Len(t, l, 2)
	assert.NotContains(t, l, Estimate{WorkPackageID: "TEST01", UserName: "Tigger"})
}

func TestLeaveSessionFailsDueToEmptyName(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.LeaveSession("12345678901234567890123456789012", "")
	assert.Equal(t, "User name should not be empty", err2.Error())
}

func TestLeaveSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.LeaveSession("123456789012345678901234567890", "")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestLeaveSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.LeaveSession("12345678901234567890123456789012", "Tigger")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestLeaveSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
	err4 := gds.JoinSession(token, "Rabbit")
	assert.NoError(t, err4)
	err5 := gds.LeaveSession(token, "Tigger")
	assert.NoError(t, err5)
}

func TestRemoveSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveSession("123456789012345678901234567890")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestRemoveSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveSession("12345678901234567890123456789012")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestRemoveSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveSession(token)
	assert.NoError(t, err3)
}

func TestAddWorkPackageToSessionFailsDueToEmptyID(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddWorkPackage("12345678901234567890123456789012", "", "eat honey")
	assert.Equal(t, "ID should not be empty", err2.Error())
}

func TestAddWorkPackageToSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddWorkPackage("1234567890123456789012345678901212", "01", "eat honey")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestAddWorkPackageToSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage("12345678901234567890123456789012", "01", "eat honey")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestAddWorkPackageToSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage(token, "01", "eat honey")
	assert.NoError(t, err3)
}

func TestAddWorkPackageToSessionErrorWhileTryingToAddWorkPackageTwiceWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage(token, "01", "eat honey")
	assert.NoError(t, err3)
	err4 := gds.AddWorkPackage(token, "01", "eat honey")
	assert.Equal(t, "Workpackage with ID: 01 already part of session", err4.Error())
}

func TestRemoveWorkPackageFromSessionFailsDueToEmptyID(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveWorkPackage("12345678901234567890123456789012", "")
	assert.Equal(t, "ID should not be empty", err2.Error())
}

func TestRemoveWorkPackageFromSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveWorkPackage("1234567890123456789012345678901212", "01")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestRemoveWorkPackageFromSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveWorkPackage("12345678901234567890123456789012", "01")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestRemoveWorkPackageFromSessionFailsDueToNonExistingIDWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveWorkPackage(token, "01")
	assert.Equal(t, "Unable to remove workpackage: 01 from session", err3.Error())
}

func TestRemoveWorkPackageFromSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage(token, "01", "eat honey")
	assert.NoError(t, err3)
	err4 := gds.RemoveWorkPackage(token, "01")
	assert.NoError(t, err4)
}

func TestRemoveWorkPackageFromSessionErrorWhileTryingToRemoveWorkPackageTwiceWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage(token, "01", "eat honey")
	assert.NoError(t, err3)
	err4 := gds.AddWorkPackage(token, "02", "harvest honey")
	assert.NoError(t, err4)
	err5 := gds.RemoveWorkPackage(token, "01")
	assert.NoError(t, err5)
	err6 := gds.RemoveWorkPackage(token, "01")
	assert.Equal(t, "Unable to remove workpackage: 01 from session", err6.Error())
}

func TestGetUsersFromSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	_, err2 := gds.GetUsers("1234567890123456789012345678901212")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestGetUsersFromSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	_, err3 := gds.GetUsers("12345678901234567890123456789012")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestGetUsersFromSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
	users, err4 := gds.GetUsers(token)
	assert.NoError(t, err4)
	assert.Equal(t, []string{"Tigger"}, users)
}

func TestGetWorkPackagesFromSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	_, err2 := gds.GetWorkPackages("1234567890123456789012345678901212")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestGetWorkPackagesFromSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	_, err3 := gds.GetWorkPackages("12345678901234567890123456789012")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestGetWorkPackagesFromSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage(token, "01", "eat honey")
	assert.NoError(t, err3)
	wps, err4 := gds.GetWorkPackages(token)
	assert.NoError(t, err4)
	assert.Equal(t, "01", wps[0].ID)
	assert.Equal(t, "eat honey", wps[0].Summary)
}

func TestAddEstimateToWorkPackageFailsDueToEmptyID(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimateToWorkPackage("12345678901234567890123456789012", "", 0.0, 0.0)
	assert.Equal(t, "ID should not be empty", err2.Error())
}

func TestAddEstimateToWorkPackageFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimateToWorkPackage("1234567890123456789012345678901212", "01", 0.0, 0.0)
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestAddEstimateToWorkPackageFailsDueToIncorrectEffort(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimateToWorkPackage("12345678901234567890123456789012", "01", -0.1, 0.0)
	assert.Equal(t, "Effort < 0 not allowed", err2.Error())
}

func TestAddEstimateToWorkPackageFailsDueToIncorrectStandardDeviation(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimateToWorkPackage("12345678901234567890123456789012", "01", 0.1, -0.1)
	assert.Equal(t, "Standard deviation < 0 not allowed", err2.Error())
}

func TestAddEstimateToWorkPackageFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddEstimateToWorkPackage("12345678901234567890123456789012", "01", 0.0, 0.0)
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestAddEstimateToWorkPackageErrorWhileTryingToAddEstimateToNonExistingWorkPackageWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddEstimateToWorkPackage(token, "01", 1.5, 0.2)
	assert.Equal(t, "Work package with ID: 01 does not exist", err3.Error())
}

func TestAddEstimateToWorkPackageSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage(token, "01", "eat honey")
	assert.NoError(t, err3)
	err4 := gds.AddEstimateToWorkPackage(token, "01", 1.5, 0.2)
	assert.NoError(t, err4)
}

func TestRemoveEstimateFromWorkPackageFailsDueToEmptyID(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveEstimateFromWorkPackage("12345678901234567890123456789012", "")
	assert.Equal(t, "ID should not be empty", err2.Error())
}

func TestRemoveEstimateFromWorkPackageFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveEstimateFromWorkPackage("1234567890123456789012345678901212", "01")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestRemoveEstimateFromWorkPackageFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveEstimateFromWorkPackage("12345678901234567890123456789012", "01")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestRemoveEstimateFromWorkPackageErrorWhileTryingToRemoveEstimateFromNonExistingWorkPackageWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveEstimateFromWorkPackage(token, "01")
	assert.Equal(t, "Work package with ID: 01 does not exist", err3.Error())
}

func TestRemoveEstimateFromWorkPackageSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage(token, "01", "eat honey")
	assert.NoError(t, err3)
	err4 := gds.AddEstimateToWorkPackage(token, "01", 1.5, 0.2)
	assert.NoError(t, err4)
	wps, err5 := gds.GetWorkPackages(token)
	assert.NoError(t, err5)
	assert.Equal(t, 1.5, wps[0].Effort)
	assert.Equal(t, 0.2, wps[0].StandardDeviation)
	err6 := gds.RemoveEstimateFromWorkPackage(token, "01")
	assert.NoError(t, err6)
	wps2, err7 := gds.GetWorkPackages(token)
	assert.NoError(t, err7)
	assert.Equal(t, 0.0, wps2[0].Effort)
	assert.Equal(t, 0.0, wps2[0].StandardDeviation)
}

func TestAddEstimateToSessionFailsDueToEmptyWorkPackageID(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimate("12345678901234567890123456789012", Estimate{WorkPackageID: "", UserName: "Tigger"})
	assert.Equal(t, "Work Package ID should not be empty", err2.Error())
}

func TestAddEstimateToSessionFailsDueToEmptyUserName(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimate("12345678901234567890123456789012", Estimate{WorkPackageID: "TEST01", UserName: ""})
	assert.Equal(t, "User name should not be empty", err2.Error())
}

func TestAddEstimateToSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimate("1234567890123456789012345678901212", Estimate{WorkPackageID: "TEST01", UserName: "Tigger"})
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestAddEstimateToSessionFailsDueToWrongValueForBestCase(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimate("12345678901234567890123456789012", Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: -0.1})
	assert.Equal(t, "Best case must be >= 0, provided: -0.1", err2.Error())
}

func TestAddEstimateToSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddEstimate("12345678901234567890123456789012", Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestAddEstimateToSessionFailsDueToEstimateAlreadyExistingWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
	err4 := gds.AddWorkPackage(token, "TEST01", "")
	assert.NoError(t, err4)
	err5 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.NoError(t, err5)
	err6 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.Equal(t, "Specified estimate already exists", err6.Error())
}

func TestAddEstimateToSessionFailsDueToUserNotPartOfSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage(token, "TEST01", "")
	assert.NoError(t, err3)
	err4 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.Equal(t, "User: Tigger is not part of session", err4.Error())
}

func TestAddEstimateToSessionFailsDueToWorkPackageNotPartOfSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
	err4 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.Equal(t, "Work Package with ID: TEST01 is not part of session", err4.Error())
}

func TestGetEstimatesFromSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	_, err2 := gds.GetEstimates("1234567890123456789012345678901212")
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestGetEstimatesFromSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	_, err3 := gds.GetEstimates("12345678901234567890123456789012")
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestAddEstimateToSessionAndGetEstimatesSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
	err4 := gds.JoinSession(token, "Rabbit")
	assert.NoError(t, err4)
	err5 := gds.AddWorkPackage(token, "TEST01", "")
	assert.NoError(t, err5)
	err6 := gds.AddWorkPackage(token, "TEST02", "eat honey")
	assert.NoError(t, err6)
	err7 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.NoError(t, err7)
	err8 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Rabbit", BestCase: 0.5, MostLikelyCase: 1.5, WorstCase: 2.5})
	assert.NoError(t, err8)
	ests, err9 := gds.GetEstimates(token)
	assert.NoError(t, err9)
	assert.Equal(t, "Tigger", ests[0].UserName)
	assert.Equal(t, "TEST01", ests[0].WorkPackageID)
	assert.Equal(t, 0.1, ests[0].BestCase)
	assert.Equal(t, 0.5, ests[0].MostLikelyCase)
	assert.Equal(t, 1.0, ests[0].WorstCase)
	assert.Equal(t, "Rabbit", ests[1].UserName)
	assert.Equal(t, "TEST01", ests[1].WorkPackageID)
	assert.Equal(t, 0.5, ests[1].BestCase)
	assert.Equal(t, 1.5, ests[1].MostLikelyCase)
	assert.Equal(t, 2.5, ests[1].WorstCase)
}

func TestRemoveEstimateFromSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveEstimate("1234567890123456789012345678901212", Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.Equal(t, "Session token does not match desired length", err2.Error())
}

func TestRemoveEstimateFromSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveEstimate("12345678901234567890123456789012", Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.Equal(t, "Specified session does not exist", err3.Error())
}

func TestRemoveEstimateFromSessionFailsDueTOEstimateNotBeingPartOfSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
	err4 := gds.JoinSession(token, "Rabbit")
	assert.NoError(t, err4)
	err5 := gds.AddWorkPackage(token, "TEST01", "")
	assert.NoError(t, err5)
	err6 := gds.AddWorkPackage(token, "TEST02", "eat honey")
	assert.NoError(t, err6)
	err7 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.NoError(t, err7)
	err8 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Rabbit", BestCase: 0.5, MostLikelyCase: 1.5, WorstCase: 2.5})
	assert.NoError(t, err8)
	ests, err9 := gds.GetEstimates(token)
	assert.NoError(t, err9)
	assert.Len(t, ests, 2)
	err10 := gds.RemoveEstimate(token, Estimate{WorkPackageID: "TEST02", UserName: "Rabbit", BestCase: 0.5, MostLikelyCase: 1.5, WorstCase: 2.5})
	assert.Equal(t, "Estimate with ID: TEST02 and user name: Rabbit is not part of session", err10.Error())
	ests2, err11 := gds.GetEstimates(token)
	assert.NoError(t, err11)
	assert.Len(t, ests2, 2)
}

func TestRemoveEstimateFromSessionSuccessWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession(token, "Tigger")
	assert.NoError(t, err3)
	err4 := gds.JoinSession(token, "Rabbit")
	assert.NoError(t, err4)
	err5 := gds.AddWorkPackage(token, "TEST01", "")
	assert.NoError(t, err5)
	err6 := gds.AddWorkPackage(token, "TEST02", "eat honey")
	assert.NoError(t, err6)
	err7 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Tigger", BestCase: 0.1, MostLikelyCase: 0.5, WorstCase: 1.0})
	assert.NoError(t, err7)
	err8 := gds.AddEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Rabbit", BestCase: 0.5, MostLikelyCase: 1.5, WorstCase: 2.5})
	assert.NoError(t, err8)
	ests, err9 := gds.GetEstimates(token)
	assert.NoError(t, err9)
	assert.Len(t, ests, 2)
	err10 := gds.RemoveEstimate(token, Estimate{WorkPackageID: "TEST01", UserName: "Rabbit", BestCase: 0.5, MostLikelyCase: 1.5, WorstCase: 2.5})
	assert.NoError(t, err10)
	ests2, err11 := gds.GetEstimates(token)
	assert.NoError(t, err11)
	assert.Len(t, ests2, 1)
}
