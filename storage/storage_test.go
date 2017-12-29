package storage

import (
	"fmt"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xnaveira/pingo/model"
)

var testMatches model.Matches
var testMatch model.Match
var testMatch2 model.Match
var testID uuid.UUID

//TODO: Table testing and indiviudal tests
func init() {
	testID = uuid.NewV4()
	testGame := []model.Result{model.Result{A: 0, B: 0}, model.Result{A: 0, B: 1}, model.Result{A: 0, B: 2}}
	testGame2 := []model.Result{model.Result{A: 1, B: 1}, model.Result{A: 0, B: 1}, model.Result{A: 0, B: 2}}
	testMatch = model.Match{ID: testID, PlayerA: "Petter", PlayerB: "Xavier", Games: []model.Game{testGame}}
	testMatch2 = model.Match{ID: testID, PlayerA: "Petter", PlayerB: "Xavier", Games: []model.Game{testGame2}}
	testMatches = append(testMatches, testMatch)
}

func TestStorageFunctions(t *testing.T) {
	var emptyMatches model.Matches = []model.Match{}
	m, err := RepoMatchGetAll()
	if !reflect.DeepEqual(m, emptyMatches) {
		t.Errorf("Error getting an empty list of matches: %v", m)
	}
	//Create
	err = RepoMatchCreate(testMatch)
	if err != nil {
		t.Errorf("Got an error creating match: %s", err.Error())
	}
	//Get
	findResult, err := RepoMatchGet(testID)
	if err != nil && !reflect.DeepEqual(findResult, model.Match{}) {
		t.Errorf("RepoFindMatch expected %+v and got %+v.", findResult, model.Match{})
	}
	if !reflect.DeepEqual(findResult, testMatch) {
		t.Errorf("RepoFindMatch expected %+v and got %+v", findResult, testMatch)
	}
	randomID := uuid.NewV4()
	findResult, err = RepoMatchGet(randomID)
	if !reflect.DeepEqual(err, fmt.Errorf("Could not find Match with id of %v", randomID)) {
		t.Errorf("Expected void match and error and got %+v and %s", findResult, err.Error())
	}
	//Modify
	modifyResult, err := RepoMatchModify(testID, testMatch2)
	if err != nil {
		t.Fatalf("Error testing RepoMAtchModify: %v", err)
	}
	findResult, err = RepoMatchGet(testID)
	if err != nil {
		t.Fatalf("Error testing RepoMAtchModify, get error: %v", err)
	}
	if !reflect.DeepEqual(modifyResult, findResult) {
		t.Errorf("RepoMatchModify expected %+v and got\n %+v", modifyResult, findResult)
	}

	modifyResult, err = RepoMatchModify(uuid.NewV4(), testMatch2)
	if !reflect.DeepEqual(modifyResult, model.Match{}) {
		t.Errorf("RepoMatchModify expected %+v and got\n %+v", modifyResult, findResult)
	}
	//Delete
	err = RepoMatchDelete(testID)
	if err != nil {
		t.Errorf("RepoDestroyMatch expected nil but got %s", err.Error())
	}
	err = RepoMatchDelete(randomID)
	if !reflect.DeepEqual(err, fmt.Errorf("Could not find Match with id of %d to delete", randomID)) {
		t.Errorf("RepoDestroyMatch expected error and got %s", err.Error())
	}

}
