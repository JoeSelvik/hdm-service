package main

import (
	"database/sql"
	"fmt"
	"github.com/JoeSelvik/hdm-service/models"
	"log"
	"net/http"
	"time"
)

type PostController struct {
	config *Configuration
	db     *models.DB
	fh     Facebooker
}

// ServeHTTP routes incoming requests to the right service.
func (pc *PostController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := new(Contender)
	ServeResource(w, r, pc, c)
}

// Path returns the URL extension associated with the Contender resource.
func (pc *PostController) Path() string {
	return "/posts/"
}

// DBTableName returns the table name for Contenders.
func (pc *PostController) DBTableName() string {
	return "posts"
}

// Create writes a new post to the db for each given Resource.
func (pc *PostController) Create(m []Resource) ([]int, *ApplicationError) {
	// Create a slice of Contender pointers by asserting on a slice of Resources interfaces
	var posts []*Post
	for _, p := range m {
		posts = append(posts, p.(*Post))
	}

	// Create the SQL query
	q := fmt.Sprintf(`
	INSERT INTO %s (
		fb_id, fb_group_id,
		posted_date, author, likes,
		created_at, updated_at
	) values (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, pc.DBTableName())

	// Begin sql transaction
	tx, err := pc.db.Begin()
	if err != nil {
		msg := "Something is wrong with our database - we'll be back up soon!"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}
	defer tx.Rollback()

	var postIds []int
	for _, p := range posts {
		likes := sliceOfIntsToString(p.Likes)

		result, err := tx.Exec(q,
			p.FbId, p.FbGroupId,
			p.PostedDate, p.Author, likes)
		if err != nil {
			msg := fmt.Sprintf("Couldn't create post: %+v", p)
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		// Save each id to return
		id, err := result.LastInsertId()
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}
		postIds = append(postIds, int(id))
	}

	// Commit sql transaction
	if err = tx.Commit(); err != nil {
		msg := "Something is wrong with our database - we'll be back up soon!"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	return postIds, nil
}

// Read returns the post in the db for a given FbId.
func (pc *PostController) Read(fbId int) (Resource, *ApplicationError) {
	return nil, &ApplicationError{Code: http.StatusNotImplemented}
}

// Update writes the db column value for each variable post parameter.
//
// Writes TotalPosts, AvgLikesPerPost, TotalLikesReceived, TotalLikesGiven, PostsUsed, and UpdatedAt.
// todo: test when fb_id does not exist
func (pc *PostController) Update(m []Resource) *ApplicationError {
	msg := "No variable data to update on posts"
	return &ApplicationError{Msg: msg, Code: http.StatusNotImplemented}
}

// Destroy deletes any given id from the db.
func (pc *PostController) Destroy(ids []int) *ApplicationError {
	return &ApplicationError{Code: http.StatusNotImplemented}
}

// ReadCollection returns all posts in the db.
func (pc *PostController) ReadCollection() ([]Resource, *ApplicationError) {
	// Grab rows from table
	rows, err := pc.db.Query(fmt.Sprintf("SELECT * FROM %s", pc.DBTableName()))
	switch {
	case err == sql.ErrNoRows:
		log.Println("Contenders ReadCollection: no rows in table.")
		return []Resource{}, nil
	case err != nil:
		msg := "Something is wrong with our database - we'll be back up soon!"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}
	defer rows.Close()

	// Create a Contender from each row
	posts := make([]Resource, 0) // Container for the Resources we're about to return
	for rows.Next() {
		var fbId string
		var fbGroupId int
		var postedDate time.Time
		var author string
		var likesString string
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&fbId, &fbGroupId, &postedDate, &author, &likesString, &createdAt, &updatedAt)
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		// Split comma separated strings to slices of ints
		likes, err := stringOfIntsToSliceOfInts(likesString)
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		c := Post{
			FbId:       fbId,
			FbGroupId:  fbGroupId,
			PostedDate: postedDate,
			Author:     author,
			Likes:      likes,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		}

		posts = append(posts, &c)
	}
	return posts, nil
}

// /////
// Non API methods and helper functions
// todo: does this section belong?
// /////

func (pc *PostController) PopulatePostsTable() *ApplicationError {
	log.Println("Pulling posts from facebook and creating in db")

	// Get slice of post struct pointers from fb
	posts, aerr := pc.fh.PullPostsFromFb()
	if aerr != nil {
		return aerr
	}

	// Convert each contender struct ptr to Resource interface
	postResources := make([]Resource, len(posts))
	for i, v := range posts {
		postResources[i] = Resource(v)
	}

	// Populate Contenders table
	_, aerr = pc.Create(postResources)
	if aerr != nil {
		return aerr
	}

	return nil
}
