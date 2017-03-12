package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	fb "github.com/huandu/facebook"
	"log"
	"time"
)

type Contender struct {
	Id                 string `facebook:",required"`
	Name               string
	TotalPosts         []string
	TotalLikesReceived int
	AvgLikesPerPost    int
	TotalLikesGiven    int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (c *Contender) DBTableName() string {
	return "contenders"
}

func (c *Contender) Path() string {
	return "/contenders/"
}

// CreateContender places the Contender into the contenders table
func (c *Contender) CreateContender(tx *sql.Tx) (int64, error) {
	q := `
	INSERT INTO contenders (
		Id,
		Name,
		TotalPosts,
		TotalLikesReceived,
		AvgLikesPerPost,
		TotalLikesGiven,
		CreatedAt,
		UpdatedAt
	) values (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	posts, err := json.Marshal(c.TotalPosts)
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(q, c.Id, c.Name, posts, c.TotalLikesReceived, c.AvgLikesPerPost, c.TotalLikesGiven)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Contender) UpdateContender(tx *sql.Tx) (int64, error) {
	posts, err := json.Marshal(c.TotalPosts)
	if err != nil {
		return 0, err
	}

	q := `UPDATE contenders SET TotalPosts = ?, TotalLikesReceived = ?, AvgLikesPerPost = ?, TotalLikesGiven = ?, UpdatedAt = CURRENT_TIMESTAMP WHERE Id = ?`
	result, err := tx.Exec(q, posts, c.TotalLikesReceived, c.AvgLikesPerPost, c.TotalLikesGiven, c.Id)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to update %s's row: %v", c.Name, err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Contender) updateIndependentData() {

}

func (c *Contender) updateTotalPosts(id string, tx *sql.Tx) {
	// get up to dat contender fromdb for posts data
	db := GetDBHandle()

	// Get current posts for contender from DB
	updatedContender, err := GetContenderByUsername(db, c.Name)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to get %s's TotalPosts's: %v", c.Name, err))
	}

	fmt.Printf("original contender posts: %+v\n", c.TotalPosts)
	fmt.Printf("updated contender: %+v\n", updatedContender.TotalPosts)
	newPosts := append(updatedContender.TotalPosts, id)
	fmt.Printf("new posts: %v\n\n", newPosts)

	posts, err := json.Marshal(newPosts)

	// c.TotalPosts = append(updatedContender.TotalPosts, id)

	// place back in db
	// todo: update updated time
	q := `UPDATE contenders SET TotalPosts = ? WHERE Name=?;`
	_, err = tx.Exec(q, posts, c.Name)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to increment %s's TotalLikesGiven: %v", c.Name, err))
	}
}

func (c *Contender) updateTotalLikesRx() {

}

// only incraments by one
func (c *Contender) updateTotalLikesGiven(tx *sql.Tx) error {
	c.TotalLikesGiven++
	// todo: update updated time
	q := `UPDATE contenders SET TotalLikesGiven = TotalLikesGiven + 1 WHERE Name=?;`
	_, err := tx.Exec(q, c.Name)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to increment %s's TotalLikesGiven: %v", c.Name, err))
		return err
	}
	return nil
}

// CreateContenderTable creates the contenders table if it does not exist
func CreateContenderTable(db *sql.DB) error {
	q := `
	CREATE TABLE IF NOT EXISTS contenders(
		Id TEXT NOT NULL,
		Name TEXT,
		TotalPosts BLOB,
		TotalLikesReceived INT,
		AvgLikesPerPost INT,
		TotalLikesGiven INT,
		CreatedAt DATETIME,
		UpdatedAt DATETIME
	);
	`
	_, err := db.Exec(q)
	if err != nil {
		log.Println("Failed to CREATE contenders table")
		return err
	}

	session := GetFBSession()
	fbContenders, err := GetFBContenders(session)
	if err != nil {
		log.Fatal("Failed to get members from facebook")
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to BEGIN txn:", err)
		return err
	}
	defer tx.Rollback()

	for i := 0; i < len(fbContenders); i++ {
		// should this check if contender already exists?
		_, err := fbContenders[i].CreateContender(tx)
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

func UpdateHDMContenderDependentData() {
	db := GetDBHandle()
	posts, err := GetHDMPosts(db)
	handle_error("Could not get posts", err, true)

	contenders := make(map[string]Contender)

	// key, value: Post.Id, Post
	for _, p := range posts {

		// Grab contender from db if it has not been updated yet
		var poster *Contender
		if val, ok := contenders[p.Author]; ok {
			poster = &val
		} else {
			poster, _ = GetContenderByUsername(db, p.Author)
		}

		poster.TotalPosts = append(poster.TotalPosts, p.Id)
		likesReceived := 0
		for i := 0; i < len(poster.TotalPosts); i++ {
			likesReceived = len(posts[poster.TotalPosts[i]].Likes.Data) + likesReceived
		}

		poster.TotalLikesReceived = likesReceived
		poster.AvgLikesPerPost = poster.TotalLikesReceived / len(poster.TotalPosts)

		// for each like, give a likes given
		for j := 0; j < len(p.Likes.Data); j++ {
			var liker *Contender
			if val, ok := contenders[p.Likes.Data[j].Name]; ok {
				liker = &val
			} else {
				liker, _ = GetContenderByUsername(db, p.Likes.Data[j].Name)
			}
			liker.TotalLikesGiven++
			contenders[liker.Name] = *liker
		}

		contenders[poster.Name] = *poster
	}

	tx, err := db.Begin()
	handle_error("Failed to BEGIN txn", err, true)
	defer tx.Rollback()

	// Update every contender in db that was updated
	// key, value: Contender.Id, Contender
	for _, c := range contenders {
		_, err := c.UpdateContender(tx)
		if err != nil {
			handle_error("Could not update contender", err, true)
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Println("Failed to COMMIT txn:", err)
	}
}

// func UpdateContenderTable() {
// 	session := GetFBSession()
// 	db := GetDBHandle()
// 	fbContenders := GetFBContenders(session)
// 	contenders := GetHDMContenders(db)

// 	for c
// }

// todo: should I return *Contender or Contender?
func GetContenderByUsername(db *sql.DB, name string) (*Contender, error) {
	q := "SELECT * FROM contenders WHERE name = ?"
	var id string
	var totalPosts string // sqlite blob later to be unmarshalled
	var totalLikesReceived int
	var avgLikesPerPost int
	var totalLikesGiven int
	var createdAt time.Time
	var updatedAt time.Time

	// todo: is this bad overwritting name?
	err := db.QueryRow(q, name).Scan(&id, &name, &totalPosts, &totalLikesReceived, &avgLikesPerPost, &totalLikesGiven, &createdAt, &updatedAt)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that Name.")
		return nil, err
	case err != nil:
		log.Fatal(fmt.Sprintf("Error getting %s from contenders table %s", name, err))
		return nil, err
	default:
		var posts []string
		json.Unmarshal([]byte(totalPosts), &posts)

		c := &Contender{
			Id:                 id,
			Name:               name,
			TotalPosts:         posts,
			TotalLikesReceived: totalLikesReceived,
			AvgLikesPerPost:    avgLikesPerPost,
			TotalLikesGiven:    totalLikesGiven,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
		}
		return c, nil
	}
}

func GetContenderById(db *sql.DB, id string) (*Contender, error) {
	q := "SELECT * FROM contenders WHERE id = ?"
	var name string
	var totalPosts string // sqlite blob later to be unmarshalled
	var totalLikesReceived int
	var avgLikesPerPost int
	var totalLikesGiven int
	var createdAt time.Time
	var updatedAt time.Time

	// todo: is this bad overwritting name?
	err := db.QueryRow(q, id).Scan(&id, &name, &totalPosts, &totalLikesReceived, &avgLikesPerPost, &totalLikesGiven, &createdAt, &updatedAt)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
		return nil, err
	case err != nil:
		log.Fatal(fmt.Sprintf("Error getting %s from contenders table %s", name, err))
		return nil, err
	default:
		var posts []string
		json.Unmarshal([]byte(totalPosts), &posts)

		c := &Contender{
			Id:                 id,
			Name:               name,
			TotalPosts:         posts,
			TotalLikesReceived: totalLikesReceived,
			AvgLikesPerPost:    avgLikesPerPost,
			TotalLikesGiven:    totalLikesGiven,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
		}
		return c, nil
	}
}

// GetHDMContenders returns a slice of Contenders from the contenders table
func GetHDMContenders(db *sql.DB) ([]Contender, error) {
	rows, err := db.Query("SELECT * FROM contenders")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var contenders []Contender

	for rows.Next() {
		var id string
		var name string
		var totalPosts string // sqlite blob later to be unmarshalled
		var totalLikesReceived int
		var avgLikesPerPost int
		var totalLikesGiven int
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&id, &name, &totalPosts, &totalLikesReceived, &avgLikesPerPost, &totalLikesGiven, &createdAt, &updatedAt)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to scan contender from table: %v", err))
			return nil, err
		}

		fmt.Printf("totalPosts: %s\n", totalPosts)

		var posts []string
		json.Unmarshal([]byte(totalPosts), &posts)
		fmt.Printf("posts: %s\n", posts)

		c := Contender{
			Id:                 id,
			Name:               name,
			TotalPosts:         posts,
			TotalLikesReceived: totalLikesReceived,
			AvgLikesPerPost:    avgLikesPerPost,
			TotalLikesGiven:    totalLikesGiven,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
		}
		contenders = append(contenders, c)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return contenders, nil
}

// Returns a slice of Contenders for a given *Session from a FB group
func GetFBContenders(session *fb.Session) ([]Contender, error) {
	// response is a map[string]interface{}
	response, err := fb.Get(fmt.Sprintf("/%s/members", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
		"feilds":       []string{"name", "id", "picture", "context", "cover"},
	})
	if err != nil {
		log.Fatal("Error requesting group members")
		return nil, err
	}

	// Get the member's paging object
	paging, err := response.Paging(session)
	if err != nil {
		log.Fatal("Error generating the member response Paging object")
		return nil, err
	}

	var contenders []Contender

	for {
		results := paging.Data()

		// map[administrator:false name:Jacob Glowacki id:1822807864675176]
		for i := 0; i < len(results); i++ {
			var c Contender
			facebookContender := fb.Result(results[i]) // cast the var
			c.Id = facebookContender.Get("id").(string)
			c.Name = facebookContender.Get("name").(string)

			contenders = append(contenders, c)
		}

		noMore, err := paging.Next()
		if err != nil {
			log.Fatal("Error accessing Response page's Next object")
			return nil, err
		}
		if noMore {
			break
		}
	}

	log.Println("Number of Contenders:", len(contenders))
	return contenders, nil
}
