package main

import (
	"fmt"
	"github.com/JoeSelvik/hdm-service/models"
	"net/http"
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
	for i := 0; i < len(m); i++ {
		p := m[i]
		posts = append(posts, p.(*Post))
	}

	// Create the SQL query
	q := fmt.Sprintf(`
	INSERT INTO %s (
		fb_id, fb_group_id,
		posted_date, author, total_likes,
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
		result, err := tx.Exec(q,
			p.FbId, p.FbGroupId,
			p.PostedDate, posts, p.Author, p.Likes)
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

	return postIds, &ApplicationError{Code: http.StatusNotImplemented}
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
	return nil, &ApplicationError{Code: http.StatusNotImplemented}
}
