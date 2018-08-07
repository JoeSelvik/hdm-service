package main

import (
	"github.com/JoeSelvik/hdm-service/models"
	"testing"
)

// Create, ReadCollection, Update, Read, Destroy, Read contender.
func TestPostController_Create(t *testing.T) {
	config := NewConfig()
	db, err := models.OpenDB(config.DbTestPath)
	if err != nil {
		t.Fatal("Could not open test db")
	}
	pc := &PostController{db: db}

	// Create a post struct and convert it to a resource interface
	testAuthorId := 1225
	posts := []*Post{
		{
			AuthorFbId: testAuthorId,
			FbId:       "666_999"},
	}
	postResources := make([]Resource, len(posts))
	for i, v := range posts {
		postResources[i] = Resource(v)
	}

	pids, aerr := pc.Create(postResources)
	if aerr != nil {
		t.Fatal("Should not error when creating post")
	}
	if len(pids) != 1 {
		t.Fatal("Should only get back a single id")
	}

	// Read all posts
	resources, aerr := pc.ReadCollection()
	if aerr != nil {
		t.Fatal("Unable to read collection of posts")
	}
	var lookup *Post
	for _, c := range resources {
		if c.(*Post).AuthorFbId == testAuthorId {
			lookup = c.(*Post)
			break
		}
	}
	if lookup.AuthorFbId != testAuthorId {
		t.Fatal("Unable to find post in ReadCollection")
	}
}
