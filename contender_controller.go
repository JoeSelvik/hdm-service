package main

import (
	"database/sql"
	"fmt"
	fb "github.com/huandu/facebook"
)

// A ContenderController struct has a DB wrapper so that we can use all of our RESTful functions
//
// See scooby.db for a custom example
type ContenderController struct {
	db *sql.DB
}

// CreateContenderTable creates the contenders table if it does not exist
func CreateContenderTable(db *sql.DB) {
	sql_table := `
	CREATE TABLE IF NOT EXISTS contenders(
		Id TEXT NOT NULL,
		Name TEXT,
		TotalPosts INT,
		TotalLikesReceived INT,
		AvgLikesPerPost INT,
		TotalLikesGiven INT
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

// DB returns the database associated with the Contender resources
func (cc *ContenderController) DB() *sql.DB {
	return cc.db
}

func (cc *ContenderController) DBTableName() string {
	return "contenders"
}

func (cc *ContenderController) Path() string {
	return "/contenders/"
}

func (cc *ContenderController) Create(c *Contender) (*Contender, error) {
	// Put this new User into the db
	m, err := cc.DB().Create(uc, u)
	if err != nil {
		return nil, err
	}
}

// Returns a slice of Contenders for a given *Session
func populateContenders(session *fb.Session) []Contender {
	// response is a map[string]interface{}
	response, err := fb.Get(fmt.Sprintf("/%s/members", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
	})
	handle_error("Error when getting group members", err, true)

	// Create the paging object for /members response
	paging, err := response.Paging(session)
	handle_error("Error when generating the members responses Paging object", err, true)

	var contenders []Contender

	for {
		results := paging.Data()

		// map[administrator:false name:Jacob Glowacki id:1822807864675176]
		var c Contender

		for i := 0; i < len(results); i++ {
			results[i].Decode(&c)
			contenders = append(contenders, c)
		}

		noMore, err := paging.Next()
		handle_error("Error when accessing responses Next in loop:", err, true)
		if noMore {
			break
		}
	}

	return contenders
}
