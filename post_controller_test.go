package main

import (
	"github.com/JoeSelvik/hdm-service/models"
	"testing"
)

// TestPostController_Create tests Create and Read for posts.
//
// todo: expand once all of the post functions are implemented
func TestPostController_Create(t *testing.T) {
	config := NewConfig()
	db, err := models.OpenDB(config.DbTestPath)
	if err != nil {
		t.Fatal("Could not open test db")
	}
	pc := &PostController{db: db}
	originalPosts, postResources := testPosts()

	// create the post resource
	pids, aerr := pc.Create(postResources)
	if aerr != nil {
		t.Fatal("Should not error when creating post")
	}
	if len(pids) != 1 {
		t.Fatal("Should only get back a single id")
	}

	// read all posts
	resources, aerr := pc.ReadCollection()
	if aerr != nil {
		t.Fatal("Unable to read collection of posts")
	}
	var lookup *models.Post
	for _, c := range resources {
		if c.(*models.Post).AuthorFbId == originalPosts[0].AuthorFbId {
			lookup = c.(*models.Post)
			break
		}
	}
	if lookup.AuthorFbId != originalPosts[0].AuthorFbId {
		t.Fatal("Unable to find post in ReadCollection")
	}
}
