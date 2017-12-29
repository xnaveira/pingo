package model

import uuid "github.com/satori/go.uuid"

//Result is the result of a set
type Result struct {
	A int `json:"a"`
	B int `json:"b"`
}

//Game is a bunch of results
type Game []Result

//Match holds all the info about a match
type Match struct {
	ID      uuid.UUID `json:"id, omitempty"`
	PlayerA string    `json:"playera"`
	PlayerB string    `json:"playerb"`
	Games   []Game    `json:"games"`
}

//Matches is a slice of Match
type Matches []Match
