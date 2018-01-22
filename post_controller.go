package main

import (
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
	return nil, &ApplicationError{Code: http.StatusNotImplemented}
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
