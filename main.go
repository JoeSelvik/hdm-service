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
	"net/http"
	"os"
)

const (
	AccessToken     = "EAACEdEose0cBAPQeEO3dZBZBWTjG66umHO0C9w8Gr9QoP0hUzHoWkmZB5IerAAaBzfKwdkKu1KjPZBKsg0lO9VlacXGmltFkVgbGHPuKtnDszUZBYwPoo1ZAZCUMeXWo3R8JC2TIE8JZBTcDVYA4EE8E2Wi728J86HSDy78YyANWk50bOwjMG0xz"
	HerpDerpGroupID = "208678979226870"
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

func bracketDataHandler(w http.ResponseWriter, r *http.Request) {
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

func setupDatabase() {
	// sqlite setup and verification
	db, err := sql.Open("sqlite3", "herp.db")
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

	CreateContenderTable(db)
	CreatePostsTable(db)
}

func getFBData() {
	// Facebook setup
	var myAccessToken = GetAccessToken()

	// "your-app-id", "your-app-secret", from 'development' app I made
	var globalApp = fb.New("756979584457445", "023c1d8f5e901c2111d7d136f5165b2a")
	session := globalApp.Session(myAccessToken)
	err := session.Validate()
	handle_error("Error validating session", err, true)

	contenders := populateContenders(session)
	fmt.Println("number of members:", len(contenders))

	populateTotalPosts(contenders, session)
}

func main() {
	// setupDatabase()
	// getFBData()
	http.HandleFunc("/bracketData/", bracketDataHandler)
	http.ListenAndServe(":8080", nil)
}
