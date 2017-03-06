package main

import (
	"database/sql"
	"fmt"
	fb "github.com/huandu/facebook"
	"log"
	"time"
)

type Post struct {
	Id         string `facebook:",required"`
	PostedDate string
	Author     string
	TotalLikes int
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func (p *Post) DBTableName() string {
	return "posts"
}

func (p *Post) Path() string {
	return "/posts/"
}

func (p *Post) CreatePost(tx *sql.Tx) (int64, error) {
	q := `
	INSERT INTO posts (
		Id,
		PostedDate,
		Author,
		TotalLikes,
		CreatedAt,
		UpdatedAt
	) values (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := tx.Exec(q, p.Id, p.PostedDate, p.Author, p.TotalLikes, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// asumes never goesdown?
func (p *Post) updateTotalLikes(n int) {

}

// CreatePostsTable creates the posts table if it does not exist
func CreatePostsTable(startDate string, db *sql.DB) error {
	q := `
	CREATE TABLE posts(
		Id TEXT NOT NULL,
		PostedDate DATETIME,
		Author TEXT,
		TotalLikes INT,
		CreatedAt DATETIME,
		UpdatedAt DATETIME
	);
	`

	// don't care about this result
	_, err := db.Exec(q)
	if err != nil {
		log.Println("Failed to CREATE posts table")
		return err
	}

	session := GetFBSession()
	fbPosts := GetFBPosts("blah", session)

	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to BEGIN txn:", err)
		return err
	}
	defer tx.Rollback()

	for i := 0; i < len(fbPosts); i++ {
		// should this check if post already exists?
		_, err := fbPosts[i].CreatePost(tx)
		if err != nil {
			return err
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Println("Failed to COMMIT txn:", err)
		return err
	}

	return nil
}

func GetFBPosts(startDate string, session *fb.Session) []Post {
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

	// loop until a fb post's created_time is older than startDate
	for {
		results := paging.Data()

		// 25 posts per page, load data into a Post struct
		for i := 0; i < len(results); i++ {
			var p Post
			facebookPost := fb.Result(results[i]) // cast the var
			p.PostedDate = facebookPost.Get("created_time").(string)

			// break if post is older than startDate
			// fb created_time str:      2017-03-04T13:05:20+0000
			// sqlite CURRENT_TIMESTAMP: 2017-03-06 15:36:17
			// layout := "Mon, 01/02/06, 03:04PM"
			// layout := "Mon Jan 2 15:04:05 MST 2006  (MST is GMT-0700)"
			value := "2017-03-04T13:05:20+0000"
			layout := "2006-01-02T15:04:05+0000"
			t, err := time.Parse(layout, p.PostedDate)
			if err != nil {
				log.Fatal("Failed to parse post's postedDate")
				return nil
			}
			t.
				p.Id = facebookPost.Get("id").(string)
			p.Author = facebookPost.Get("from.name").(string)

			if facebookPost.Get("likes.data") != nil {
				numLikes := facebookPost.Get("likes.data").([]interface{})
				p.TotalLikes = len(numLikes)

				// for each like, give a likes given
			} else {
				p.TotalLikes = 0
			}

			posts = append(posts, p)

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
	// fmt.Println("number posts:", len(posts))
	// fmt.Println("First post:", posts[1])
	return posts
}
