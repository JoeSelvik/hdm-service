package main

import (
	"database/sql"
	"fmt"
	fb "github.com/huandu/facebook"
)

type Post struct {
	Id          string `facebook:",required"`
	CreatedDate string
	Author      string
	TotalLikes  int
}

// CreatePostsTable creates the posts table if it does not exist
func CreatePostsTable(db *sql.DB) {
	sql_table := `
	CREATE TABLE IF NOT EXISTS posts(
		Id TEXT NOT NULL,
		CreatedDate DATETIME,
		Author TEXT,
		TotalLikes INT
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func populateTotalPosts(contenders []Contender, session *fb.Session) {
	// Get the group feed
	response, err := fb.Get(fmt.Sprintf("/%s/feed", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
		"feilds":       []string{"from", "created_time"},
	})
	handle_error("Error when getting feed", err, true)

	// Get the feed's paging object
	paging, err := response.Paging(session)
	handle_error("Error when generating the feed responses Paging object", err, true)

	var posts []Post
	count := 0

	// 25 posts per page
	for {
		results := paging.Data()

		// load data from each facebookPost into a Post struct
		for i := 0; i < len(results); i++ {
			var p Post
			facebookPost := fb.Result(results[i])

			id := facebookPost.Get("id")
			p.Id = id.(string)

			author := facebookPost.Get("from.name")
			p.Author = author.(string)

			createdDate := facebookPost.Get("created_time")
			p.CreatedDate = createdDate.(string)

			likesData := facebookPost.Get("likes.data")
			if likesData != nil {
				numLikes := facebookPost.Get("likes.data").([]interface{})
				p.TotalLikes = len(numLikes)
			} else {
				p.TotalLikes = 0
			}

			posts = append(posts, p)
			fmt.Println("Decoded post:", p)
		}

		count++
		fmt.Println("finished lap:", count)

		if count >= 1 {
			fmt.Println("found first 25 posts")
			break
		}

		noMore, err := paging.Next()
		handle_error("Error when accessing responses Next in loop", err, true)
		if noMore {
			fmt.Println("Reached the end of the feed!")
			break
		}
	}
	fmt.Println("number posts:", len(posts))
	fmt.Println("First post:", posts[1])
}
