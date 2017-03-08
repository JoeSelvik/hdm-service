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
	AccessToken     = "EAACEdEose0cBAAmg1hKFs34c3peBRSLrZCuep0kmWHzTaOTLPmdXBpLZCp8ixcRdpOBL1r90ChtWbcsOSBZCai8o9KglgkSxYNKKEl1ZCmKw5KE80oVEmX5xVkPAO6zkqaiU4ZBLWURv2zdcmXtEYVP5YTvradSRZCPTyvPdzrKlzhVJmuZB7mn"
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

func sampleBracketDataHandler(w http.ResponseWriter, r *http.Request) {
	teams := make([][]string, 4)
	teams[0] = []string{"joe", "matt"}
	teams[1] = []string{"tj", "cody"}
	teams[2] = []string{"george", "jim"}
	teams[3] = []string{"ted", "tim"}

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

	results := make([]interface{}, 3)
	results[0] = firstRound
	results[1] = secondRound
	results[2] = thirdRound

	// bracket := Bracket{teams, results}
	bracket := fullBracket()
	js, err := json.Marshal(bracket)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

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

	// err := CreateContenderTable(db)
	// if err != nil {
	// 	log.Println("Failed to create contenders table:", err)
	// 	return
	// }

	err := CreatePostsTable(GetStartTime(), db)
	if err != nil {
		log.Println("Failed to create posts table:", err)
		return
	}

	// bracket tables
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
	setupDatabase()
	// getFBData()

	// http.HandleFunc("/bracketData/", sampleBracketDataHandler)
	// http.ListenAndServe(":8080", nil)
}
