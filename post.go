package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

type Post struct {
	Id        string `facebook:",required"`
	CreatedAt time.Time
	UpdatedAt time.Time

	PostedDate string
	Author     string
	Likes      Likes
}

type Like struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Likes struct {
	Data []Like `json:"data"`
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
		Likes,
		CreatedAt,
		UpdatedAt
	) values (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	likes, err := json.Marshal(p.Likes)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(q, p.Id, p.PostedDate, p.Author, likes)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CreatePostsTable creates the posts table
//
// todo: should this check if post already exists and table
// and only print new entries?
func CreatePostsTable(startDate time.Time, db *sql.DB) error {
	q := `
	CREATE TABLE IF NOT EXISTS posts(
		Id TEXT NOT NULL,
		PostedDate DATETIME,
		Author TEXT,
		Likes BLOB,
		CreatedAt DATETIME,
		UpdatedAt DATETIME
	);
	`
	_, err := db.Exec(q)
	if err != nil {
		log.Println("Failed to CREATE posts table")
		return err
	}

	fbPosts, err := PullPostsFromFb(startDate)
	if err != nil {
		log.Println("Failed to get posts from fb:", err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to BEGIN txn:", err)
		return err
	}
	defer tx.Rollback()

	// Create each Post in DB
	for i := 0; i < len(fbPosts); i++ {
		_, err := fbPosts[i].CreatePost(tx)
		if err != nil {
			log.Printf("Failed to create post")
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

// GetHDMPosts returns a map of each Post in the DB indexed by Id
func GetHDMPosts(db *sql.DB) (map[string]Post, error) {
	rows, err := db.Query("SELECT * FROM posts")
	if err != nil {
		log.Println("Failed to query Posts table:", err)
		return nil, err
	}
	defer rows.Close()

	posts := make(map[string]Post)

	for rows.Next() {
		var id string
		var postedDate string // todo: should this be a time.Time?
		var author string
		var strLikes string // sqlite blob later to be unmarshalled
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&id, &postedDate, &author, &strLikes, &createdAt, &updatedAt)
		if err != nil {
			log.Println("Failed to scan post from posts table:", err)
			return nil, err
		}

		likes := Likes{}
		json.Unmarshal([]byte(strLikes), &likes)

		p := Post{
			Id:         id,
			PostedDate: postedDate,
			Author:     author,
			Likes:      likes,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		}
		posts[p.Id] = p
	}

	err = rows.Err()
	if err != nil {
		log.Println("Detected err from Posts row scan:", err)
		return nil, err
	}
	return posts, nil
}
