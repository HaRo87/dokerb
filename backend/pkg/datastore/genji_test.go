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
	assert.Errorf(t, err, "Proper DB must be provided and not nil")
}

func TestErrorOnDBSetupCreateTable(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(fmt.Errorf("Ooops, something went wrong"))
	_, err := NewGenjiDatastore(m)
	assert.Errorf(t, err, "Unable to create sessions table")

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
	assert.Errorf(t, err, "Invalid token length provided: -20, should be >= 20")
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
	assert.Errorf(t, err2, "Unable to store session token")
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
	assert.Errorf(t, err2, "User name should not be empty")
}

func TestJoinSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.JoinSession("1234567890123456789012345678901212", "")
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestJoinSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.JoinSession("12345678901234567890123456789012", "Tigger")
	assert.Errorf(t, err3, "Specified session does not exist")
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
	assert.Errorf(t, err4, "User with name: Tigger already part of session")
}

func TestRemoveUserFromEmptyList(t *testing.T) {
	l, err := removeUser([]string{}, "Tigger")
	assert.Errorf(t, err, "User with name: Tigger is not part of session")
	assert.Len(t, l, 0)
}

func TestRemoveUserFromListWithoutThatUserBeingPartOfThatList(t *testing.T) {
	users := []string{"Tigger", "Rabbit", "Piglet"}
	l, err := removeUser(users, "Winnie-the-Pooh")
	assert.Errorf(t, err, "User with name: Winnie-the-Pooh is not part of session")
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
	assert.Errorf(t, err, "Workpackage with ID: TEST01 is not part of session")
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
	assert.Errorf(t, err, "Workpackage with ID: TEST04 is not part of session")
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
	assert.Errorf(t, err, "Estimate with ID: TEST01 and user name: Tigger is not part of session")
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
	assert.Errorf(t, err, "Estimate with ID: TEST04 and user name: Tigger is not part of session")
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
	assert.Errorf(t, err, "Estimate with ID: TEST01 and user name: Piglet is not part of session")
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
	assert.Errorf(t, err2, "User name should not be empty")
}

func TestLeaveSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.LeaveSession("123456789012345678901234567890", "")
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestLeaveSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.LeaveSession("12345678901234567890123456789012", "Tigger")
	assert.Errorf(t, err3, "Specified session does not exist")
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
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestRemoveSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveSession("12345678901234567890123456789012")
	assert.Errorf(t, err3, "Specified session does not exist")
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
	assert.Errorf(t, err2, "ID should not be empty")
}

func TestAddWorkPackageToSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddWorkPackage("1234567890123456789012345678901212", "01", "eat honey")
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestAddWorkPackageToSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddWorkPackage("12345678901234567890123456789012", "01", "eat honey")
	assert.Errorf(t, err3, "Specified session does not exist")
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
	assert.Errorf(t, err4, "Workpackage with ID: 01 already part of session")
}

func TestRemoveWorkPackageFromSessionFailsDueToEmptyID(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveWorkPackage("12345678901234567890123456789012", "")
	assert.Errorf(t, err2, "ID should not be empty")
}

func TestRemoveWorkPackageFromSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveWorkPackage("1234567890123456789012345678901212", "01")
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestRemoveWorkPackageFromSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveWorkPackage("12345678901234567890123456789012", "01")
	assert.Errorf(t, err3, "Specified session does not exist")
}

func TestRemoveWorkPackageFromSessionFailsDueToNonExistingIDWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveWorkPackage(token, "01")
	assert.Errorf(t, err3, "Unable to remove workpackage: 01 from session")
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
	assert.Errorf(t, err6, "Unable to remove workpackage: 01 from session")
}

func TestGetUsersFromSessionFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	_, err2 := gds.GetUsers("1234567890123456789012345678901212")
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestGetUsersFromSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	_, err3 := gds.GetUsers("12345678901234567890123456789012")
	assert.Errorf(t, err3, "Specified session does not exist")
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
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestGetWorkPackagesFromSessionFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	_, err3 := gds.GetWorkPackages("12345678901234567890123456789012")
	assert.Errorf(t, err3, "Specified session does not exist")
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
	assert.Errorf(t, err2, "ID should not be empty")
}

func TestAddEstimateToWorkPackageFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimateToWorkPackage("1234567890123456789012345678901212", "01", 0.0, 0.0)
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestAddEstimateToWorkPackageFailsDueToIncorrectEffort(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimateToWorkPackage("12345678901234567890123456789012", "01", -0.1, 0.0)
	assert.Errorf(t, err2, "Effort < 0 not allowed")
}

func TestAddEstimateToWorkPackageFailsDueToIncorrectStandardDeviation(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.AddEstimateToWorkPackage("12345678901234567890123456789012", "01", 0.1, -0.1)
	assert.Errorf(t, err2, "Standard deviation < 0 not allowed")
}

func TestAddEstimateToWorkPackageFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddEstimateToWorkPackage("12345678901234567890123456789012", "01", 0.0, 0.0)
	assert.Errorf(t, err3, "Specified session does not exist")
}

func TestAddEstimateToWorkPackageErrorWhileTryingToAddEstimateToNonExistingWorkPackageWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.AddEstimateToWorkPackage(token, "01", 1.5, 0.2)
	assert.Errorf(t, err3, "Work package with ID: 01 does not exist")
}

func TestAddEstimateToSessionSuccessWithRealDB(t *testing.T) {
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
	assert.Errorf(t, err2, "ID should not be empty")
}

func TestRemoveEstimateFromWorkPackageFailsDueToWrongTokenLength(t *testing.T) {
	setupAndTearDown := setupTestCaseForMock(t)
	defer setupAndTearDown(t)
	m.On("Exec", "CREATE TABLE sessions").Return(nil)
	gds, err := NewGenjiDatastore(m)
	assert.NoError(t, err)
	err2 := gds.RemoveEstimateFromWorkPackage("1234567890123456789012345678901212", "01")
	assert.Errorf(t, err2, "Session token does not match desired length")
}

func TestRemoveEstimateFromWorkPackageFailsDueToNonExistingSessionWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	_, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveEstimateFromWorkPackage("12345678901234567890123456789012", "01")
	assert.Errorf(t, err3, "Specified session does not exist")
}

func TestRemoveEstimateFromWorkPackageErrorWhileTryingToRemoveEstimateFromNonExistingWorkPackageWithRealDB(t *testing.T) {
	setupAndTearDown := setupTestCaseForRealDB(t)
	defer setupAndTearDown(t)
	gds, err := NewGenjiDatastore(db)
	assert.NoError(t, err)
	token, err2 := gds.CreateSession()
	assert.NoError(t, err2)
	err3 := gds.RemoveEstimateFromWorkPackage(token, "01")
	assert.Errorf(t, err3, "Work package with ID: 01 does not exist")
}

func TestRemoveEstimateFromSessionSuccessWithRealDB(t *testing.T) {
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
