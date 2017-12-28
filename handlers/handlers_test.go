package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/xnaveira/pingo/model"
)

var testIndexOut = fmt.Sprintln(`{"message":"Welcome!"}`)

func TestIndex(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Index)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	//expected := `{"message":"Welcome!"}`
	if rr.Body.String() != testIndexOut {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), testIndexOut)
	}
}

var createdID uuid.UUID

func TestCreateMatch(t *testing.T) {

	type testMatch struct {
		match    string
		result   string
		httpcode int
	}

	var testMatches = []testMatch{
		testMatch{
			`{"playera":"Petter","playerb":"Xavier","games":[{"a":0,"b":0}]}`,
			`"playera":"Petter","playerb":"Xavier","games":[{"a":0,"b":0}]}`,
			http.StatusCreated},
		testMatch{
			`{"playea":"Petter","games":"games"}`,
			`{"id":"00000000-0000-0000-0000-000000000000","playera":"","playerb":"","games":null}`, //model.Match{}
			http.StatusUnprocessableEntity},
	}
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	for _, m := range testMatches {

		inputMatch := m.match
		expectedResult := m.result
		expectedCode := m.httpcode
		match := strings.NewReader(inputMatch)

		req, err := http.NewRequest("POST", "/", match)
		if err != nil {
			t.Fatal(err)
		}
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MatchCreate)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		var status int
		if status = rr.Code; status != expectedCode {
			t.Errorf("Testing %d handler returned wrong status code: got %v want %v",
				m.httpcode, status, expectedCode)
		}
		// Check the response body is what we expect.
		if !strings.Contains(rr.Body.String(), expectedResult) {
			t.Errorf("Testing %d handler returned unexpected body: got %v want %v",
				m.httpcode, rr.Body.String(), expectedResult)
		}
		//Get the id for the case created for the get match test
		if status == http.StatusCreated {
			var modelMatch model.Match
			_ = json.Unmarshal(rr.Body.Bytes(), &modelMatch)

			createdID = modelMatch.ID
		}

	}

}

func TestMatchModify(t *testing.T) {

	newUUID := uuid.NewV4()

	type testMatchModifyData struct {
		name          string
		oldmatchid    string
		newmatch      string
		expectedMatch string
		expectedCode  int
	}
	testData := []testMatchModifyData{
		testMatchModifyData{
			"Should modify it",
			createdID.String(),
			fmt.Sprintf(`{"id":"%s","playera":"Petter","playerb":"Xavier","games":[{"a":1,"b":0},{"a":0,"b":0}]}`, createdID.String()),
			fmt.Sprintf(`{"id":"%s","playera":"Petter","playerb":"Xavier","games":[{"a":1,"b":0},{"a":0,"b":0}]}`, createdID.String()),
			201,
		},
		testMatchModifyData{
			"Shouldn't modify it",
			newUUID.String(),
			fmt.Sprintf(`{"id":"%s","playera":"Petter","playerb":"Xavier","games":[{"a":1,"b":0},{"a":0,"b":0}]}`, newUUID.String()),
			`{"id":"00000000-0000-0000-0000-000000000000","playera":"","playerb":"","games":null}`, //model.Match{}
			500,
		},
		testMatchModifyData{
			"Modify unprocessable",
			createdID.String(),
			fmt.Sprintf(`{"id":"%s","plaxyera":"Toby","plaxyerb":"Xavier","games":1}`, createdID.String()),
			`{"id":"00000000-0000-0000-0000-000000000000","playera":"","playerb":"","games":null}`, //model.Match{}
			422,
		},
	}

	for _, td := range testData {
		req, err := http.NewRequest("PUT", fmt.Sprintf("/match/%s", td.oldmatchid), strings.NewReader(td.newmatch))
		if err != nil {
			panic(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MatchModify)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != td.expectedCode {
			t.Errorf("Testing `%s` handler returned wrong status code: got %d want %d",
				td.name, status, td.expectedCode)
		}
		if strings.Compare(strings.TrimSpace(rr.Body.String()), td.expectedMatch) != 0 {
			t.Errorf("Testing `%s` handler returned wrong match:\n got  %v want %v",
				td.name, rr.Body.String(), td.expectedMatch)
		}
	}
}

func TestMatchGet(t *testing.T) {

	type testMatchShowData struct {
		name   string
		input  string
		output int //http state
	}
	testData := []testMatchShowData{
		testMatchShowData{
			"Should find it",
			createdID.String(),
			http.StatusOK,
		},
		testMatchShowData{
			"Shouldn't find it",
			uuid.NewV4().String(),
			http.StatusNotFound,
		},
		testMatchShowData{
			"Show bad data",
			"1234",
			http.StatusBadRequest,
		},
	}

	for _, td := range testData {
		expectedCode := td.output
		req, err := http.NewRequest("GET", fmt.Sprintf("/match/%s", td.input), nil)
		if err != nil {
			panic(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MatchGet)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != expectedCode {
			t.Errorf("Testing `%s` handler returned wrong status code: got %d with %v want %d",
				td.name, status, td.input, expectedCode)
		}
	}
}

func TestMatchDelete(t *testing.T) {

	type testMatchDeleteData struct {
		name   string
		input  string
		output int //http state
	}
	testData := []testMatchDeleteData{
		testMatchDeleteData{
			"Should delete it",
			createdID.String(),
			http.StatusOK,
		},
		testMatchDeleteData{
			"Shouldn't delete it",
			uuid.NewV4().String(),
			http.StatusNotFound,
		},
		testMatchDeleteData{
			"Baddata",
			"1234",
			http.StatusBadRequest,
		},
	}

	for _, td := range testData {
		expectedCode := td.output
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/match/%s", td.input), nil)
		if err != nil {
			panic(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MatchDelete)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != expectedCode {
			t.Errorf("Testing `%s` handler returned wrong status code: got %d with %v want %d",
				td.name, status, td.input, expectedCode)
		}
	}
}
