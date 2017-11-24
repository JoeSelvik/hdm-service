package main

import (
	// "database/sql"
	// "encoding/json"
	"fmt"
	"log"
	"time"
)

type Matchup struct {
	Id         int
	Name       string // ie: firstRound_g0
	ContenderA *Contender
	APosts     []string // slice of perm_urls
	AVotes     int
	ContenderB *Contender
	BPosts     []string // slice of perm_urls
	BVotes     int
	InProgress bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m *Matchup) serialize() {

}

func (m *Matchup) DBTableName() string {
	return "matchups"
}

func (m *Matchup) Path() string {
	return "/matchups/"
}

// Places 32 Matchups into DB with InProgress set to true
func CreateFirstRoundMatchups() {
	db := GetDBHandle(NewConfig())
	bracket, _ := GetHDMBracket(db, 1)

	matchups := make(map[string]*Matchup)

	for i := 0; i < len(bracket.Teams); i++ {
		name := bracket.Results.FirstRound[i][2].(string)
		log.Println("Creating Matchup: ", name)

		teams := bracket.Teams[i]
		// todo: bye's print No user with that name message, create a finished Matchup
		contenderA, _ := GetContenderByUsername(db, teams.ContenderAName)
		contenderB, _ := GetContenderByUsername(db, teams.ContenderBName)

		// Get five random posts for each contender
		// todo: mark these posts as used?
		// todo: need to write a GetHDMPostsByContenderUsername function
		// var aPosts []string

		m := Matchup{
			Name:       name,
			ContenderA: contenderA,
			APosts:     []string{"permalink_url1", "permalink_url2"},
			AVotes:     0,
			ContenderB: contenderB,
			BPosts:     []string{"permalink_url1", "permalink_url2"},
			BVotes:     0,
			InProgress: true,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		matchups[m.Name] = &m
	}

	for key, matchup := range matchups {
		log.Println(fmt.Sprintf("Matchup %s: %+v", key, matchup))
	}
}

// Updates 32 Matchups in DB to InProgress false
func EndFirstRoundMatchups() {

}

// todo: how to get all teamPairs in second round?
func CreateSecondRoundMatchups() {

}

// GetHDMMatchup pulls any matchup from the DB and returns data needed to render matchupView
//
// name == firstRound_31
// func GetHDMMatchup(name string) *Matchup {
// 	db := GetDBHandle()
// 	bracket, _ := GetHDMBracket(db, 1)

// 	var id string
// 	var contenderA Contender
// 	var aPosts string

// 	m := Matchup{}
// }
