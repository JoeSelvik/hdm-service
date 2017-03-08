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
func CreatePostsTable(startDate time.Time, db *sql.DB) error {
	q := `
	CREATE TABLE IF NOT EXISTS posts(
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
	fbPosts, err := GetFBPosts(startDate, session)
	if err != nil {
		log.Fatal("Failed to get posts from facebook")
		return err
	}

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

func GetFBPosts(startDate time.Time, session *fb.Session) ([]Post, error) {
	// Get the group feed
	response, err := fb.Get(fmt.Sprintf("/%s/feed", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
		"feilds":       []string{"from", "created_time"},
	})
	if err != nil {
		log.Fatal("Error requesting group feed")
		return nil, err
	}

	// Get the feed's paging object
	paging, err := response.Paging(session)
	if err != nil {
		log.Fatal("Error generating the feed response Paging object")
		return nil, err
	}

	var posts []Post
	count := 1

	// loop until a fb post's created_time is older than startDate
Loop:
	for {
		results := paging.Data()
		fmt.Println("Posts page ", count)

		// 25 posts per page, load data into a Post struct
		for i := 0; i < len(results); i++ {
			var p Post
			facebookPost := fb.Result(results[i]) // cast the var
			p.PostedDate = facebookPost.Get("created_time").(string)

			// stop when post reaches startDate
			t, err := time.Parse(GoTimeLayout, p.PostedDate)
			if err != nil {
				log.Fatal("Failed to parse post's postedDate")
				return nil, err
			}
			if t.Before(startDate) {
				log.Println("Reached a post before the startDate")
				break Loop
			}

			p.Id = facebookPost.Get("id").(string)
			p.Author = facebookPost.Get("from.name").(string)
			p.PostedDate = t.String()

			if facebookPost.Get("likes.data") != nil {
				numLikes := facebookPost.Get("likes.data").([]interface{})
				p.TotalLikes = len(numLikes)

				// for each like, give a likes given
				for j := 0; i < len(numLikes); j++ {
					c := GetContenderByUsername(GetDBHandle(), p.Author)

				}

			} else {
				p.TotalLikes = 0
			}

			posts = append(posts, p)
		}

		noMore, err := paging.Next()
		if err != nil {
			log.Fatal("Error accessing Response page's Next object")
			return nil, err
		}
		if noMore {
			fmt.Println("Reached the end of group feed")
			break Loop
		}
		count++
	}
	fmt.Println("Number posts:", len(posts))
	return posts, nil
}
