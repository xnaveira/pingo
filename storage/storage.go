package storage

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/xnaveira/pingo/model"
)

//Provisional storage module, not concurrently safe and in memory only

var matches model.Matches = []model.Match{}

//RepoMatchGetAll all the Matches
func RepoMatchGetAll() (model.Matches, error) {
	return matches, nil
}

//RepoMatchCreate stores a match
func RepoMatchCreate(m model.Match) error {
	matches = append(matches, m)
	return nil
}

//RepoMatchModify gets a new match and replaces an existing one
func RepoMatchModify(ID uuid.UUID, newmatch model.Match) (model.Match, error) {
	for i, m := range matches {
		if m.ID == ID {
			matches = append(matches[:i], matches[i+1:]...)
			matches = append(matches, newmatch)
			return newmatch, nil
		}
	}
	return model.Match{}, fmt.Errorf("Could not find Match with id of %d to modify", ID)
}

//RepoMatchGet retrieves a match from storage
func RepoMatchGet(ID uuid.UUID) (model.Match, error) {
	for _, m := range matches {
		if m.ID == ID {
			return m, nil
		}
	}
	return model.Match{}, fmt.Errorf("Could not find Match with id of %v", ID)
}

//RepoMatchDelete deletes a match from storage
func RepoMatchDelete(ID uuid.UUID) error {
	for i, m := range matches {
		if m.ID == ID {
			matches = append(matches[:i], matches[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Match with id of %d to delete", ID)
}
