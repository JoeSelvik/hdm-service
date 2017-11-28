/*
Dependancies
* github.com/mattn/go-sqlite3

*/
package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var (
	Config *Configuration
)

func HandleError(msg string, err error, exit bool) {
	if err != nil {
		fmt.Println(msg, ":", err)

		if exit {
			os.Exit(3)
		}
	}
}

// GetDBHandle opens a connection and returns an active handle to the sqlite db
func GetDBHandle(c *Configuration) *sql.DB {
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

func main() {
	log.Println("Welcome to the HerpDerp Madness service")
	log.Println()

	// Parse the config, define global Config variable.
	Config = NewConfig()

	// Print the config
	log.Printf("Facebook Access Token:\t%d\n", Config.FbAccessToken)
	log.Printf("Facebook Group Id:\t%d\n", Config.FbGroupId)
	log.Printf("Database:\t%s\n", Config.DbPath)
	log.Println()

	// Create db handle
	db := GetDBHandle(Config)

	con, err := PullContendersFromFb()
	if err != nil {
		panic("Couldn't get Facebook contenders")
	}
	log.Println(con)

	//posts, err := PullPostsFromFb(Config.StartTime)
	//if err != nil {
	//	panic("Couldn't get Facebook posts")
	//}
	//log.Println(len(posts))

	// Register http handlers
	cc := &ContenderController{db: db}
	http.Handle(cc.Path(), cc)

	// Register speak handle
	http.HandleFunc("/speak/", speakHandle)

	// Listen on port
	// todo: handle this with channels and check for errors?
	http.ListenAndServe(":8080", nil)
}

// speakHandle will return random dog quotes. Used to verify the system is up.
func speakHandle(w http.ResponseWriter, r *http.Request) {
	var speakQuotes = loadDogSounds()
	i := rand.Intn(len(speakQuotes))
	fmt.Fprintf(w, speakQuotes[i])
}

func loadDogSounds() []string {
	var dogSounds []string

	// Populate list of dog sounds.
	f, err := os.Open("dog_sounds.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		dogSounds = append(dogSounds, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return dogSounds
}
