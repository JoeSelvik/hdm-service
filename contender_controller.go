package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ContenderController struct {
	db *sql.DB
}

// Path returns the URL extension associated with the Contender resource.
func (cc *ContenderController) Path() string {
	return "/contenders/"
}

// DBTableName returns the table name for Contenders.
func (cc *ContenderController) DBTableName() string {
	return "contenders"
}

func (cc *ContenderController) Create(c Contender) (int64, error) {
	// todo: %s and cc.DBTableName() instead?
	// todo: time.Now() instead of CURRENT_TIMESTAMP?
	q := `
	INSERT INTO contenders (
		fb_id, fb_group_id,
		name, total_posts, avg_likes_per_post, total_likes_received, total_likes_given, posts_used,
		created_at, updated_at
	) values (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	result, err := cc.db.Exec(q,
		c.FbId, c.FbGroupId,
		c.Name, c.TotalPosts, c.AvgLikesPerPost, c.TotalLikesReceived, c.TotalLikesGiven, c.PostsUsed)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ServeHTTP routes incoming requests to the right service.
func (cc *ContenderController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := new(Contender)
	ServeResource(w, r, cc, c)
}

// ReadCollection will display all the users. This might be restricted to Admin only later.
func (cc *ContenderController) ReadCollection(m Resource) (*[]Resource, error) {
	log.Println("Read collection: Contenders.")

	rows, err := cc.db.Query("SELECT * FROM contenders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contenders := make([]Resource, 0) // Container for the Resources we're about to return

	for rows.Next() {
		var id int
		var name string
		var totalPosts string // sqlite blob later to be unmarshalled
		var totalLikesReceived int
		var avgLikesPerPost int
		var totalLikesGiven int
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&id, &name, &totalPosts, &totalLikesReceived, &avgLikesPerPost, &totalLikesGiven, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		var posts []int
		json.Unmarshal([]byte(totalPosts), &posts)

		c := Contender{
			FbId:               id,
			Name:               name,
			TotalPosts:         posts,
			TotalLikesReceived: totalLikesReceived,
			AvgLikesPerPost:    avgLikesPerPost,
			TotalLikesGiven:    totalLikesGiven,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
		}

		// Make a new resource Value of type m.
		// todo: why does this work?
		//mType := reflect.TypeOf(m).Elem()
		//mVal := reflect.New(mType)
		//contenders = append(contenders, mVal.Interface().(Resource))

		contenders = append(contenders, &c)
	}

	return &contenders, nil
}
