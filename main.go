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

	// Parse the config
	config := NewConfig()

	log.Printf("Facebook Access Token:\t%s\n", config.FbAccessToken)
	log.Printf("Facebook Group Id:\t%d\n", config.FbGroupId)
	log.Printf("Database:\t%s\n", config.DbPath)
	log.Println()

	// Open the DB
	db, err := models.OpenDB(config.DbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the fb handle
	fh := FacebookHandle{config: config}

	// todo: add cli option to display these instead of running service
	//fetchContendersFromFb(&fh)
	//fetchPostsFromFb(&fh)

	// Register http handlers
	cc := &ContenderController{config: config, db: db, fh: &fh}
	http.Handle(cc.Path(), cc)

	pc := &PostController{config: config, db: db, fh: &fh}
	http.Handle(pc.Path(), pc)

	// Note - On 4/4/18 Facebook removed many of the api endpoints used.
	//        Commented out until an alternative solution is implemented

	// Create Contenders
	//aerr := cc.PopulateContendersTable()
	//if aerr != nil {
	//	panic(fmt.Sprintf("Could not populate contenders: %s\n%s", aerr, aerr))
	//}

	// Create Posts
	//aerr = pc.PopulatePostsTable()
	//if aerr != nil {
	//	panic(fmt.Sprintf("Could not populate posts: %s\n%s", aerr, aerr.Err))
	//}

	// Update contender's VDD
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

// fetchContendersFromFb is a convenience function to display Contenders fetched from facebook.
//
// Note - On 4/4/18 Facebook removed many of the api endpoints needed. This no longer works.
func fetchContendersFromFb(fh *FacebookHandle) {
	contenders, aerr := fh.PullContendersFromFb()
	if aerr != nil {
		panic(fmt.Sprintf("Couldn't get Facebook contenders: %s", aerr.Msg))
	}
	log.Printf("Found %d contenders\n", len(contenders))
}

// fetchPostsFromFb is a convenience function to display Posts fetched from facebook.
//
// Note - On 4/4/18 Facebook removed many of the api endpoints needed. This no longer works.
func fetchPostsFromFb(fh *FacebookHandle) {
	posts, aerr := fh.PullPostsFromFb()
	if aerr != nil {
		panic(fmt.Sprintf("Couldn't get Facebook posts: %s\n%s\n", aerr.Msg, aerr.Err))
	}
	log.Printf("Found %d posts\n", len(posts))
}
