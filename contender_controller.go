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
		TotalLikesGiven INT,
		CreatedAt DATETIME,
		UpdatedAt DATETIME
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}

	// var session = GetFBSession()
	// fbContenders = GetFBContenders(session)

	// for c in fbContenders
}

func UpdateContenderTable() {

}

func UpdateHDMContenderDependentData() {

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

// StoreContender places a Contender in the table and returns one
func CreateHDMContender(c *Contender, db *sql.DB) {
	sql_additem := `
	INSERT OR REPLACE INTO contenders(
		Id,
		Name,
		TotalPosts,
		TotalLikesReceived,
		AvgLikesPerPost,
		TotalLikesGiven,
		CreatedAt,
		UpdatedAt
	) values(?, ?, ?, ?, ?, ? CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err2 := stmt.Exec(c.Id, c.Name, c.TotalPosts, c.TotalLikesReceived, c.AvgLikePerPost, c.TotalLikesGiven)
	if err2 != nil {
		panic(err2)
	}
}

func GetHDMContenders() {

}

// Returns a slice of Contenders for a given *Session from a FB group
func GetFBContenders(session *fb.Session) []Contender {
	// response is a map[string]interface{}
	response, err := fb.Get(fmt.Sprintf("/%s/members", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
		"feilds":       []string{"name", "id", "picture", "context", "cover"},
	})
	handle_error("Error when getting group members", err, true)

	// Get the member's paging object
	paging, err := response.Paging(session)
	handle_error("Error when generating the members responses Paging object", err, true)

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
		handle_error("Error when accessing responses Next in loop:", err, true)
		if noMore {
			break
		}
	}

	fmt.Println("Number of Contenders:", len(contenders))
	return contenders
}

func (cc *ContenderController) updateIndependentData() {

}

func (cc *ContenderController) updateTotalPosts() {

}

func (cc *ContenderController) updateTotalLikesRx() {

}

func (cc *ContenderController) updateTotalLikesGiven() {

}
