package main

import (
	"bufio"
	"fmt"
	"github.com/JoeSelvik/hdm-service/models"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func main() {
	log.Println("Welcome to the HerpDerp Madness service")
	log.Println()

	// Parse the config, define global Config variable.
	config := NewConfig()

	// Print the config
	log.Printf("Facebook Access Token:\t%s\n", config.FbAccessToken)
	log.Printf("Facebook Group Id:\t%d\n", config.FbGroupId)
	log.Printf("Database:\t%s\n", config.DbPath)
	log.Println()

	// Open the DB
	db, err := models.OpenDB(config.DbPath)
	if err != nil {
		panic(err)
	}
	// todo: defer db.close()?

	// Create the fb handle
	fh := FacebookHandle{config: config}

	//// Pull fb contenders
	//contenders, aerr := PullContendersFromFb()
	//if aerr != nil {
	//	panic("Couldn't get Facebook contenders")
	//}
	//log.Println("found contenders:", len(contenders))

	//// Pull fb posts
	//posts, aerr := fh.PullPostsFromFb()
	//if aerr != nil {
	//	panic(fmt.Sprintf("Couldn't get Facebook posts: %s\n%s\n", aerr.Msg, aerr.Err))
	//}
	//log.Println("found posts:", len(posts))

	// Register http handlers
	cc := &ContenderController{config: config, db: db, fh: &fh}
	http.Handle(cc.Path(), cc)

	pc := &PostController{config: config, db: db, fh: &fh}
	http.Handle(pc.Path(), pc)

	//// Create Contenders
	//aerr := cc.PopulateContendersTable()
	//if aerr != nil {
	//	panic(fmt.Sprintf("Could not populate contenders: %s\n%s", aerr, aerr.Err))
	//}
	//
	//// Create Posts
	//aerr = pc.PopulatePostsTable()
	//if aerr != nil {
	//	panic(fmt.Sprintf("Could not populate posts: %s\n%s", aerr, aerr.Err))
	//}
	//
	//// Update contender's VDD
	//aerr = cc.UpdateContendersVariableDependentData(pc)
	//if aerr != nil {
	//	panic(fmt.Sprintf("Could not update contender's VDD: %s\n%s", aerr, aerr.Err))
	//}

	// Register speak handle
	http.HandleFunc("/speak/", speakHandle)

	// Listen on port
	// todo: handle this with channels and check for errors?
	http.ListenAndServe(":8080", nil)
}

// speakHandle prints random dog sounds to verify the system is alive.
func speakHandle(w http.ResponseWriter, r *http.Request) {
	var speakQuotes = loadDogSounds()
	i := rand.Intn(len(speakQuotes))
	fmt.Fprintf(w, speakQuotes[i])
}

// loadDogSounds returns a slice of dog sounds to choose from.
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
