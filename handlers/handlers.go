package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/xnaveira/pingo/model"
	"github.com/xnaveira/pingo/storage"
)

type jsonMsg struct {
	Message string `json:"message"`
}

const bodySizeLimit int64 = 1048576

//Index returns a welcome message
func Index(w http.ResponseWriter, r *http.Request) {

	var welcomeMsg jsonMsg
	welcomeMsg.Message = "Welcome!"
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&welcomeMsg)
	if err != nil {
		panic(err)
	}
}

//MatchIndex show a list with all the matches
func MatchIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response, err := storage.RepoMatchGetAll()
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

//MatchCreate match creates a match and returns it
func MatchCreate(w http.ResponseWriter, r *http.Request) {
	var m model.Match
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, bodySizeLimit))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.Unmarshal(body, &m); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err = json.NewEncoder(w).Encode(model.Match{}); err != nil {
			panic(err)
		}
		return
	}

	m.ID = uuid.NewV4()
	if err := storage.RepoMatchCreate(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err = json.NewEncoder(w).Encode(model.Match{}); err != nil {
			panic(err)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		panic(err)
	}
}

//MatchModify modifies an existing match
func MatchModify(w http.ResponseWriter, r *http.Request) {
	var m model.Match
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, bodySizeLimit))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err = json.Unmarshal(body, &m); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err = json.NewEncoder(w).Encode(model.Match{}); err != nil {
			panic(err)
		}
		return
	}

	modifyID := strings.TrimPrefix(r.URL.Path, "/match/")
	modifyUUID, err := uuid.FromString(modifyID)
	if err != nil {
		panic(err)
	}
	if m, err = storage.RepoMatchModify(modifyUUID, m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(model.Match{}); err != nil {
			panic(err)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		panic(err)
	}
}

//MatchGet shows a match with the Id at the end of the path
func MatchGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fetchID := strings.TrimPrefix(r.URL.Path, "/match/")
	fetchUUID, err := uuid.FromString(string(fetchID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(model.Match{}); err != nil {
			panic(err)
		}
		return
	}
	response, err := storage.RepoMatchGet(fetchUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if err = json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

//MatchDelete deletes a match with the Id at the end of the path
func MatchDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	fetchID := strings.TrimPrefix(r.URL.Path, "/match/")
	fetchUUID, err := uuid.FromString(string(fetchID))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err = json.NewEncoder(w).Encode(model.Match{}); err != nil {
			panic(err)
		}
		return
	}
	err = storage.RepoMatchDelete(fetchUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)

}

//MatchDelete deletes a match with the Id at the end of the path

// json.Unmarshal(b, *m)
// Matches = append(Matches, *m)
// t.ID = uuid.NewV4()
// t.PlayerA = "Petter"
// t.PlayerB = "Xavier"
// //t.games = make([]Result,100)
// t.Games = []Result{Result{0, 0}}
// j, _ := json.Marshal(t)
// w.Write(j)
