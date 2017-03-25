package main

import (
	// "database/sql"
	// "encoding/json"
	// "log"
	"time"
)

type Matchup struct {
	Id         int
	Name       string // ie: firstRound_g0
	ContenderA Contender
	APosts     []Post
	AVotes     int
	ContenderB Contender
	BPosts     []Post
	BVotes     int
	InProgress bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m *Matchup) DBTableName() string {
	return "matchups"
}

func (m *Matchup) Path() string {
	return "/matchups/"
}

// Places 32 Matchups into DB with InProgress set to true
func CreateFirstRoundMatchups() {

}

// Updates 32 Matchups in DB to InProgress false
func EndFirstRoundMatchups() {

}

// GetMatchup pulls any matchup from the DB and returns data needed to render matchupView
//
// name == firstRound_31
func GetMatchupData(name string) *Matchup {
	db := GetDBHandle()
	bracket, _ := GetHDMBracket(db, 1)

	var id string
	var contenderA 

	m := Matchup{

	}
}
