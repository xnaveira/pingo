package model

import (
	"encoding/json"
	"testing"

	uuid "github.com/satori/go.uuid"
)

//Test
func TestMatchMarshalling(t *testing.T) {

	jsonMatch := `{"id":"00000000-0000-0000-0000-000000000000","playera":"Petter","playerb":"Xavier","games":[[{"a":0,"b":0},{"a":1,"b":0}],[{"a":0,"b":0},{"a":0,"b":1}]]}`

	m := new(Match)
	m.ID = uuid.Nil
	m.PlayerA = "Petter"
	m.PlayerB = "Xavier"
	m.Games = []Game{[]Result{Result{0, 0}, Result{1, 0}}, []Result{Result{0, 0}, Result{0, 1}}}

	j, _ := json.Marshal(m)

	if string(j) != jsonMatch {
		t.Errorf("Marshalling error got %s\nexpected %s", j, jsonMatch)
	}

}
