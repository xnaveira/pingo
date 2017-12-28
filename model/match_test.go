package model

import (
	"encoding/json"
	"testing"

	uuid "github.com/satori/go.uuid"
)

//Test
func TestMatchMarshalling(t *testing.T) {

	jsonMatch := `{"id":"00000000-0000-0000-0000-000000000000","playera":"Petter","playerb":"Xavier","games":[{"a":0,"b":0}]}`

	m := new(Match)
	m.ID = uuid.Nil
	m.PlayerA = "Petter"
	m.PlayerB = "Xavier"
	m.Games = []Result{Result{0, 0}}

	j, _ := json.Marshal(m)

	if string(j) != jsonMatch {
		t.Errorf("Marshalling error got %s expected %s", j, jsonMatch)
	}

}
