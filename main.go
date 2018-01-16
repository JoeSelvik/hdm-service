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
	// todo: put into handles vs global?
	config := NewConfig()

	// Print the config
	log.Printf("Facebook Access Token:\t%s\n", config.FbAccessToken)
	log.Printf("Facebook Group Id:\t%d\n", config.FbGroupId)
	log.Printf("Database:\t%s\n", config.DbPath)
	log.Println()

	// Open the DB
	db, err := models.NewDB(config.DbPath)
	if err != nil {
		panic(err)
	}

	// Create the fb handle
	fh := FacebookHandle{config: config}

	//// Pull fb contenders
	//con, aerr := PullContendersFromFb()
	//if aerr != nil {
	//	panic("Couldn't get Facebook contenders")
	//}
	//log.Println("found contenders:", len(con))

	//// Pull fb posts
	//posts, aerr := PullPostsFromFb(Config.StartTime)
	//if aerr != nil {
	//	panic("Couldn't get Facebook posts")
	//}
	//log.Println("found posts:", len(posts))

	// Register http handlers
	cc := &ContenderController{config: config, db: db, fh: &fh}
	http.Handle(cc.Path(), cc)

	//// Create Contenders
	//aerr := cc.PopulateContendersTable()
	//if aerr != nil {
	//	panic(fmt.Sprintf("%s\n%s", aerr, aerr.Err))
	//}

	// Read Contenders
	contenders, aerr := cc.ReadCollection()
	if aerr != nil {
		panic(fmt.Sprintf("%s\n%s", aerr, aerr.Err))
	}
	log.Println("Number of contenders:", len(contenders))
	log.Printf("First Contender: %+v\n", contenders[0])

	//// Test update
	//cs := make([]Resource, 2) // allocates length 0 and capacity 2?
	//c0 := contenders[1].(*Contender)
	//c1 := contenders[2].(*Contender)
	//cs[0] = c0
	//cs[1] = c1
	//
	//log.Println(fmt.Sprintf("Contender0 pre modification\n%+v\n", c0))
	//log.Println(fmt.Sprintf("Contender1 pre modification\n%+v\n", c1))
	//for i := range cs {
	//	log.Println(fmt.Sprintf("Contenders pre modification\n%+v\n", cs[i]))
	//}
	//
	//c0.TotalPosts = []int{6, 6, 6}
	//c0.PostsUsed = []int{1, 2, 3}
	//c1.AvgLikesPerPost = 66
	//
	//log.Println(fmt.Sprintf("Contender0 post modification\n%+v\n", c0))
	//log.Println(fmt.Sprintf("Contender1 post modification\n%+v\n", c1))
	//for _, v := range cs {
	//	log.Println(fmt.Sprintf("Contenders post modification\n%+v\n", v))
	//}
	//
	//err = cc.Update(cs)
	//if err != nil {
	//	panic(err)
	//}

	//// Test read and destroy
	//c, aerr := cc.Read(10205178963326891)
	//if aerr != nil {
	//	panic(fmt.Sprintf("%s\n%s", aerr, aerr.Err))
	//}
	//log.Println(fmt.Sprintf("%+v", c))
	//
	//log.Println("Deleting")
	//cc.Destroy([]int{10205178963326891})
	//
	//log.Println("Reading deleted contender")
	//c, aerr = cc.Read(10205178963326891)
	//if aerr != nil {
	//	panic(fmt.Sprintf("%s\n%s", aerr, aerr.Err))
	//}
	//log.Println(fmt.Sprintf("%+v", c))

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
