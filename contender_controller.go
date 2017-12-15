package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContenderController struct {
	db *sql.DB
}

// ServeHTTP routes incoming requests to the right service.
func (cc *ContenderController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := new(Contender)
	ServeResource(w, r, cc, c)
}

// Path returns the URL extension associated with the Contender resource.
func (cc *ContenderController) Path() string {
	return "/contenders/"
}

// DBTableName returns the table name for Contenders.
func (cc *ContenderController) DBTableName() string {
	return "contenders"
}

func (cc *ContenderController) Create(m []Resource) ([]int, error) {
	// Create a slice of Contender pointers by asserting on a slice of Resources interfaces
	var contenders []*Contender
	for i := 0; i < len(m); i++ {
		c := m[i]
		contenders = append(contenders, c.(*Contender))
	}

	// Create the SQL query to use
	// todo: %s and cc.DBTableName() instead?
	// todo: time.Now() instead of CURRENT_TIMESTAMP?
	q := `
	INSERT INTO contenders (
		fb_id, fb_group_id,
		name, total_posts, avg_likes_per_post, total_likes_received, total_likes_given, posts_used,
		created_at, updated_at
	) values (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	// Begin sql transaction
	tx, err := cc.db.Begin()
	if err != nil {
		log.Println("Failed to begin txn:", err)
		return nil, err
	}
	defer tx.Rollback()

	// Insert each Contender into contenders table
	var contenderIds []int
	for i := 0; i < len(contenders); i++ {
		c := contenders[i]

		// todo: abstract this to a helper?
		// https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string/37533144
		posts := strings.Trim(strings.Join(strings.Split(fmt.Sprint(c.TotalPosts), " "), ","), "[]")
		postsUsed := strings.Trim(strings.Join(strings.Split(fmt.Sprint(c.PostsUsed), " "), ","), "[]")

		result, err := tx.Exec(q,
			c.FbId, c.FbGroupId,
			c.Name, posts, c.AvgLikesPerPost, c.TotalLikesReceived, c.TotalLikesGiven, postsUsed)
		if err != nil {
			log.Println("Failed to exec query when inserting contender:")
			fmt.Printf("%+v\n", c)
			log.Println("Error:", err)
			return nil, err
		}

		// Save each Id to return
		id, err := result.LastInsertId()
		if err != nil {
			log.Println("Failed to get LastInsertedId:", err)
			return nil, err
		}
		contenderIds = append(contenderIds, int(id))
	}

	// Commit sql transaction
	if err = tx.Commit(); err != nil {
		log.Println("Failed to Commit txn:", err)
		return nil, err
	}

	return contenderIds, nil
}

func (cc *ContenderController) Read(fbId int) (Resource, error) {
	log.Println("Read: Contender ", fbId)

	// todo: better way to shorten line of code and reuse in ReadCollection?
	var fbGroupId int
	var name string
	var totalPostsString string
	var avgLikesPerPost int
	var totalLikesReceived int
	var totalLikesGiven int
	var postsUsedString string
	var createdAt time.Time
	var updatedAt time.Time

	// Grab contender entry from table
	q := fmt.Sprintf("SELECT * FROM contenders WHERE fb_id=%d", fbId)
	err := cc.db.QueryRow(q).Scan(&fbId, &fbGroupId, &name, &totalPostsString, &avgLikesPerPost, &totalLikesReceived,
		&totalLikesGiven, &postsUsedString, &createdAt, &updatedAt) // todo: okay to unscan into fbId arg?
	switch {
	case err == sql.ErrNoRows:
		log.Println("Failed to find contender by id ", fbId) // 400-ish err
		return nil, err
	case err != nil:
		log.Println("Failed to query db:", err) // 500-ish err
		return nil, err
	}

	// todo: better way to abstract unloading strings of ints and creating individual contender (and ReadCollection)?
	// Split comma separated strings to slices of ints
	totalPosts, err := stringPostsToInts(totalPostsString)
	if err != nil {
		log.Println("Failed to convert total_posts to a slice of ints")
		return nil, err
	}
	postsUsed, err := stringPostsToInts(postsUsedString)
	if err != nil {
		log.Println("Failed to convert posts_used to a slice of ints")
		return nil, err
	}

	// Create Contender
	c := Contender{
		FbId:               fbId,
		FbGroupId:          fbGroupId,
		Name:               name,
		TotalPosts:         totalPosts,
		AvgLikesPerPost:    avgLikesPerPost,
		TotalLikesReceived: totalLikesReceived,
		TotalLikesGiven:    totalLikesGiven,
		PostsUsed:          postsUsed,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}

	return &c, nil
}

// ReadCollection will display all the users. This might be restricted to Admin only later.
func (cc *ContenderController) ReadCollection() ([]Resource, error) {
	log.Println("Read collection: Contenders")

	// Grab contender entries from table
	rows, err := cc.db.Query("SELECT * FROM contenders")
	if err != nil {
		log.Println("Failed to query db:", err)
		return nil, err
	}
	defer rows.Close()

	// Create a Contender from each row
	contenders := make([]Resource, 0) // Container for the Resources we're about to return
	for rows.Next() {
		var fbId int
		var fbGroupId int
		var name string
		var totalPostsString string
		var avgLikesPerPost int
		var totalLikesReceived int
		var totalLikesGiven int
		var postsUsedString string
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&fbId, &fbGroupId, &name, &totalPostsString, &avgLikesPerPost,
			&totalLikesReceived, &totalLikesGiven, &postsUsedString, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Failed to scan rows from db:", err)
			return nil, err
		}

		// Split comma separated strings to slices of ints
		totalPosts, err := stringPostsToInts(totalPostsString)
		if err != nil {
			log.Println("Failed to convert total_posts to a slice of ints")
			return nil, err
		}
		postsUsed, err := stringPostsToInts(postsUsedString)
		if err != nil {
			log.Println("Failed to convert posts_used to a slice of ints")
			return nil, err
		}

		c := Contender{
			FbId:               fbId,
			FbGroupId:          fbGroupId,
			Name:               name,
			TotalPosts:         totalPosts,
			AvgLikesPerPost:    avgLikesPerPost,
			TotalLikesReceived: totalLikesReceived,
			TotalLikesGiven:    totalLikesGiven,
			PostsUsed:          postsUsed,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
		}

		contenders = append(contenders, &c)
	}

	return contenders, nil
}

// /////
// Non API calls
// todo: does this belong?
// /////

func (cc *ContenderController) PopulateContendersTable() error {
	log.Println("Attempting to create Contenders")

	// Convert contender struct pointers into a slice of Resource interfaces
	contenders, err := PullContendersFromFb()
	contenderResources := make([]Resource, len(contenders))
	for i, v := range contenders {
		contenderResources[i] = Resource(v)
	}

	if err != nil {
		log.Println("Failed to get Contenders from fb:", err)
		return err
	}

	_, err = cc.Create(contenderResources)
	if err != nil {
		log.Println("Failed to create Contenders from FB:", err)
		return err
	}

	log.Println("Successfully created Contenders")
	return nil
}

// stringPostsToInts is a helper function that converts a string of ints to a slice of ints.
func stringPostsToInts(s string) ([]int, error) {
	stringSlice := strings.Split(s, ",")
	var intSlice []int
	if stringSlice[0] != "" {
		intSlice = make([]int, len(stringSlice))
		for i, v := range stringSlice {
			s, err := strconv.Atoi(v)
			if err != nil {
				log.Println("Failed to convert string of ints to slice:", err)
				return nil, err
			}
			intSlice[i] = s
		}
	}
	return intSlice, nil
}
