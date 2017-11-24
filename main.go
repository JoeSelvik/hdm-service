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
	AccessToken     = "EAACEdEose0cBAG1A39ZCwsE50ZBpK0XgxD4lsKZARX8ymORZBhxZBxkZAvW9tELHJOVIszdBrwG52Jbz1aMpLBEUL8RCcHZC1qkZCskaGkDZBNxMPf8zwTXQOFetZCsyZCkZC5VAEZCDkwCsW7bIQVcHOXMQ6GxsUS4a9Rc97BV26INZBAGK146HZAUIKlNZCqejML5ND7cZD"
	HerpDerpGroupID = "208678979226870"
	GoTimeLayout    = "2006-01-02T15:04:05+0000"
)

func HandleError(msg string, err error, exit bool) {
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
	HandleError("Error validating session", err, true)

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
	HandleError("Error when accessing /me", err, true)

	fmt.Println("User associated with access token: ", res)

	// TODO: is type assertion a bad idea here? Just handle the error from .Get?
	return res["id"].(string)
}

// GetGroupID returns the Herp Derp group_id
func GetGroupID() string {
	var groupID = HerpDerpGroupID
	return groupID
}

// GetDBHandle opens a connection and returns an active handle to the sqlite db
func GetDBHandle(c *Config) *sql.DB {
	//var dbPath = "data.db"

	// sqlite setup and verification
	db, err := sql.Open("sqlite3", c.DbPath)
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

func bracketDataHandler(w http.ResponseWriter, r *http.Request) {
	// db := GetDBHandle()
	// bracket, _ := GetHDMBracket(db, 1)

	bracket, _ := GenerateInitialBracket()

	bracketJS := bracket.serialize()

	// bundle up JSBracket for transport!
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
	HandleError("Could not parse start time", err, true)
	return t
}

func main() {
	log.Println("Welcome to the HerpDerp Madness service")
	log.Println()

	// Parse the config
	config := NewConfig()

	// Print the config
	log.Printf("Facebook Access Token:\t%d\n", config.FbAccessToken)
	log.Printf("Facebook Group Id:\t%d\n", config.FbGroupId)
	log.Printf("Database:\t%s\n", config.DbPath)
	log.Println()

	// Create db handle
	db := GetDBHandle(config)

	// Register http handlers
	cc := &ContenderController{db: db}
	http.Handle(cc.Path(), cc)

	http.HandleFunc("/bracketData/", bracketDataHandler)
	http.ListenAndServe(":8080", nil)
}
