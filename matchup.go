// Note - this is pretty out of date. Needs to be updated to follow the patterns used for contenders and posts.

package main

import (
	"database/sql"
	"fmt"
	"github.com/JoeSelvik/hdm-service/models"
	"log"
	"time"
)

type Matchup struct {
	Id         int
	Name       string // ie: firstRound_g0
	ContenderA *models.Contender
	APosts     []string // slice of perm_urls
	AVotes     int
	ContenderB *models.Contender
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
	// todo: use contender_controller.db
	var db *sql.DB

	bracket, _ := GetHDMBracket(db, 1)

	matchups := make(map[string]*Matchup)

	for i := 0; i < len(bracket.Teams); i++ {
		name := bracket.Results.FirstRound[i][2].(string)
		log.Println("Creating Matchup: ", name)

		//teams := bracket.Teams[i]
		// todo: bye's print No user with that name message, create a finished Matchup
		//contenderA, _ := GetContenderByUsername(db, teams.ContenderAName)
		//contenderB, _ := GetContenderByUsername(db, teams.ContenderBName)

		// Get five random posts for each contender
		// todo: mark these posts as used?
		// todo: need to write a GetHDMPostsByContenderUsername function
		// var aPosts []string

		m := Matchup{
			Name:       name,
			ContenderA: nil,
			APosts:     []string{"permalink_url1", "permalink_url2"},
			AVotes:     0,
			ContenderB: nil,
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
