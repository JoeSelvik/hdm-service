/*
Dependancies
* github.com/mattn/go-sqlite3

*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	fb "github.com/huandu/facebook"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	AccessToken     = "EAACEdEose0cBAMuvAHV8ecVbJ9gUVRwfkODHIEmIFsUIVpToUNgW7Wvg9t1ECH6SqIfZA32mccZCLItpbZCOJ87Qc5gb1KieAA8V1g4vbd2ZC3dLkCxzFCg0lj09Bl3BncE6wKJ5tzjuwDkIK8Bqn2Msq9npegchfFZCegqLZCzhQhTup3PdE5"
	HerpDerpGroupID = "208678979226870"
	GoTimeLayout    = "2006-01-02T15:04:05+0000"
)

func handle_error(msg string, err error, exit bool) {
	if err != nil {
		fmt.Println(msg, ":", err)

		if exit {
			os.Exit(3)
		}
	}
}

// GetAccessToken returns the access token needed to make authenticated requests
//
// Generated at https://developers.facebook.com/tools/explorer/
func GetAccessToken() string {
	var accessToken = AccessToken
	return accessToken
}

func GetFBSession() *fb.Session {
	// "your-app-id", "your-app-secret", from 'development' app I made
	var globalApp = fb.New("756979584457445", "023c1d8f5e901c2111d7d136f5165b2a")
	session := globalApp.Session(GetAccessToken())
	err := session.Validate()
	handle_error("Error validating session", err, true)

	return session
}

// GetUserID returns the user's id associated with the access token provided for the app.
//
// May modify this in the future to just return the user map.
func GetUserID() string {
	var myAccessToken = GetAccessToken()

	res, err := fb.Get("/me", fb.Params{
		"access_token": myAccessToken,
	})
	handle_error("Error when accessing /me", err, true)

	fmt.Println("User associated with access token: ", res)

	// TODO: is type assertion a bad idea here? Just handle the error from .Get?
	return res["id"].(string)
}

// GetGroupID returns the Herp Derp group_id
func GetGroupID() string {
	var groupID = HerpDerpGroupID
	return groupID
}

// GetDBHandle returns an active handle to the sqlite db
func GetDBHandle() *sql.DB {
	var dbPath = "hdm_dm.db"

	// sqlite setup and verification
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(fmt.Sprintf("Error when opening sqlite3: %s", err))
	}

	if db == nil {
		panic("db nil")
	}

	// sql.open may just validate its arguments without creating a connection to the database.
	// To verify that the data source name is valid, call Ping.
	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Error when pinging db: %s", err))
	}
	return db
}

func CreateSampleBracket() *Bracket {
	// Teams
	teams := make([]TeamPair, 4)
	teams[0] = TeamPair{"joe", "matt"}
	teams[1] = TeamPair{"jim", "mike"}
	teams[2] = TeamPair{"tim", "amy"}
	teams[3] = TeamPair{"tim", "amy"}

	// Each round of results
	firstRound := make([][]interface{}, 4)
	firstRound[0] = []interface{}{1, 0, "g1"}
	firstRound[1] = []interface{}{nil, nil, "g2"}
	firstRound[2] = []interface{}{nil, nil, "g3"}
	firstRound[3] = []interface{}{nil, nil, "g4"}

	secondRound := make([][]interface{}, 2)
	secondRound[0] = []interface{}{nil, nil, "g5"}
	secondRound[1] = []interface{}{nil, nil, "g6"}

	thirdRound := make([][]interface{}, 2)
	thirdRound[0] = []interface{}{nil, nil, "g7"}
	thirdRound[1] = []interface{}{nil, nil, "g8"}

	// Total results
	results := SixtyFourResults{}
	results.FirstRound = firstRound
	results.SecondRound = secondRound
	results.ThirdRound = thirdRound

	bracket := Bracket{666, teams, results, time.Now(), time.Now()}
	return &bracket
}

func sampleBracketDataHandler(w http.ResponseWriter, r *http.Request) {
	bracket := CreateSampleBracket()

	// Serialize a Bracket so jsQuery can understand it
	var bracketJS JSBracket

	// Teams needs to be a list of arrays
	// [["joe","matt"], ["jim","mike"]]
	teamJS := make([][]interface{}, 4)
	for i := 0; i < len(bracket.Teams); i++ {
		t := TeamPair{bracket.Teams[i].ContenderAName, bracket.Teams[i].ContenderBName}
		teamJS[i] = t.serialize()
	}
	bracketJS.Teams = teamJS

	// Results is a multi-dimension list
	// first list contains a list of main and a list of consolation results
	// winner list contains a list for each round
	// each round contains a list of each game
	// [
	// 	[ // main
	// 		[[1,0,"g1"],[null,null,"g2"],[null,null,"g3"],[null,null,"g4"]],  // round 1
	// 		[[null,null,"g5"],[null,null,"g6"]],  // round 2
	// 		[[null,null,"g7"],[null,null,"g8"]]  // round 3
	// 	]
	// ] // no consolation
	resultJS := make([]interface{}, 3)
	resultJS[0] = bracket.Results.FirstRound
	resultJS[1] = bracket.Results.SecondRound
	resultJS[2] = bracket.Results.ThirdRound
	bracketJS.Results = resultJS

	js, err := json.Marshal(bracketJS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// GetStartTime returns the time you want to get posts up until
func GetStartTime() time.Time {
	// fb created_time str:      2017-03-04T13:05:20+0000
	// sqlite CURRENT_TIMESTAMP: 2017-03-06 15:36:17
	// Golang template time      Mon, 01/02/06, 03:04PM
	// HDM golang template       Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700)
	value := "2016-01-01T00:00:00+0000"
	t, err := time.Parse(GoTimeLayout, value)
	handle_error("Could not parse start time", err, true)
	return t
}

func setupDatabase() {
	db := GetDBHandle()

	err := CreateContenderTable(db)
	if err != nil {
		log.Println("Failed to create contenders table:", err)
		return
	}

	err = CreatePostsTable(GetStartTime(), db)
	if err != nil {
		log.Println("Failed to create posts table:", err)
		return
	}

	// bracket tables
	// _ = CreateBracketsTable(db)
}

func getFBData() {
	var session = GetFBSession()
	// _, err := GetFBContenders(session)
	// handle_error("Error getting FBContenders", err, true)
	_, err := GetFBPosts(GetStartTime(), session)
	handle_error("Error getting FBPosts", err, true)

	// db := GetDBHandle()
	// contenders, err := GetHDMContenders(db)
	// handle_error("issue getting hdm contdenders", err, true)
	// fmt.Println("Number of Contenders:", len(contenders))
}

func main() {
	log.Println("hdm Madness")

	setupDatabase()
	UpdateHDMContenderDependentData()

	// CreateInitialTeams()

	// http.HandleFunc("/bracketData/", sampleBracketDataHandler)
	// http.ListenAndServe(":8080", nil)
}
