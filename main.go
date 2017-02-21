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

type Contender struct {
	Id                 string `facebook:",required"`
	Name               string
	TotalPosts         int
	TotalLikesReceived int
	AvgLikePerPost     int
	TotalLikesGiven    int
}

type Post struct {
	Id          string `facebook:",required"`
	CreatedDate string
	Author      string
	TotalLikes  int
}

// CreateContenderTable creates the contenders table if it does not exist
func CreateContenderTable(db *sql.DB) {
	sql_table := `
	CREATE TABLE IF NOT EXISTS contenders(
		Id TEXT NOT NULL,
		Name TEXT,
		TotalPosts INT,
		TotalLikesReceived INT,
		AvgLikesPerPost INT,
		TotalLikesGiven INT
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

// CreatePostsTable creates the posts table if it does not exist
func CreatePostsTable(db *sql.DB) {
	sql_table := `
	CREATE TABLE IF NOT EXISTS posts(
		Id TEXT NOT NULL,
		CreatedDate DATETIME,
		Author TEXT,
		TotalLikes INT
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

// Returns a slice of Contenders for a given *Session
func populateContenders(session *fb.Session) []Contender {
	// response is a map[string]interface{}
	response, err := fb.Get(fmt.Sprintf("/%s/members", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
	})
	handle_error("Error when getting group members", err, true)

	// Create the paging object for /members response
	paging, err := response.Paging(session)
	handle_error("Error when generating the members responses Paging object", err, true)

	var contenders []Contender

	for {
		results := paging.Data()

		// map[administrator:false name:Jacob Glowacki id:1822807864675176]
		var c Contender

		for i := 0; i < len(results); i++ {
			results[i].Decode(&c)
			contenders = append(contenders, c)
		}

		noMore, err := paging.Next()
		handle_error("Error when accessing responses Next in loop:", err, true)
		if noMore {
			break
		}
	}

	return contenders
}

func populateTotalPosts(contenders []Contender, session *fb.Session) {
	// Get the group feed
	response, err := fb.Get(fmt.Sprintf("/%s/feed", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
		"feilds":       []string{"from", "created_time"},
	})
	handle_error("Error when getting feed", err, true)

	// Get the feed's paging object
	paging, err := response.Paging(session)
	handle_error("Error when generating the feed responses Paging object", err, true)

	var posts []Post
	count := 0

	// 25 posts per page
	for {
		results := paging.Data()

		// load data from each facebookPost into a Post struct
		for i := 0; i < len(results); i++ {
			var p Post
			facebookPost := fb.Result(results[i])

			id := facebookPost.Get("id")
			p.Id = id.(string)

			author := facebookPost.Get("from.name")
			p.Author = author.(string)

			createdDate := facebookPost.Get("created_time")
			p.CreatedDate = createdDate.(string)

			likesData := facebookPost.Get("likes.data")
			if likesData != nil {
				numLikes := facebookPost.Get("likes.data").([]interface{})
				p.TotalLikes = len(numLikes)
			} else {
				p.TotalLikes = 0
			}

			posts = append(posts, p)
			fmt.Println("Decoded post:", p)
		}

		count++
		fmt.Println("finished lap:", count)

		if count >= 1 {
			fmt.Println("found first 25 posts")
			break
		}

		noMore, err := paging.Next()
		handle_error("Error when accessing responses Next in loop", err, true)
		if noMore {
			fmt.Println("Reached the end of the feed!")
			break
		}
	}
	fmt.Println("number posts:", len(posts))
	fmt.Println("First post:", posts[1])
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
