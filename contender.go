package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	fb "github.com/huandu/facebook"
	"log"
	"time"
)

// todo: pointers or data as is?
type Contender struct {
	FbId               int `facebook:",required"`
	FbGroupId          int
	Name               string
	TotalPosts         []int
	AvgLikesPerPost    int
	TotalLikesReceived int
	TotalLikesGiven    int
	PostsUsed          []int
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// SetCreatedAt will set the CreatedAt attribute of a User struct
func (c *Contender) SetCreatedAt(t time.Time) {
	c.CreatedAt = t
}

// SetUpdatedAt will set the UpdatedAt attribute of a User struct
func (c *Contender) SetUpdatedAt(t time.Time) {
	c.UpdatedAt = t
}

// /////////////////
// old methods
// /////////////////

// Sort interface, http://stackoverflow.com/questions/19946992/sorting-a-map-of-structs-golang
type contenderSlice []*Contender

// Len is part of sort.Interface.
func (c contenderSlice) Len() int {
	return len(c)
}

// Swap is part of sort.Interface.
func (c contenderSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less is part of sort.Interface. Use AvgLikesPerPost as the value to sort by
func (c contenderSlice) Less(i, j int) bool {
	return c[i].AvgLikesPerPost > c[j].AvgLikesPerPost
}

// CreateContender places the Contender into the contenders table
func (c *Contender) CreateContender(tx *sql.Tx) (int64, error) {
	q := `
	INSERT INTO contenders (
		fb_id,
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

	result, err := tx.Exec(q, c.FbId, c.Name, posts, c.TotalLikesReceived, c.AvgLikesPerPost, c.TotalLikesGiven)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateContender updates the dependent data fields and updatedAt on the contender's db row
func (c *Contender) UpdateContender(tx *sql.Tx) (int64, error) {
	posts, err := json.Marshal(c.TotalPosts)
	if err != nil {
		return 0, err
	}

	q := `UPDATE contenders SET TotalPosts = ?, TotalLikesReceived = ?, AvgLikesPerPost = ?, TotalLikesGiven = ?, UpdatedAt = CURRENT_TIMESTAMP WHERE Id = ?`
	result, err := tx.Exec(q, posts, c.TotalLikesReceived, c.AvgLikesPerPost, c.TotalLikesGiven, c.FbId)
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

// UpdateHDMContenderDependentData creates a map of Contenders to update
//
// For each Post it updates the posting Contender's TotalPosts, TotalLikesReceived,
// and AvgLikesPerPost. Then for each like on the Post it updates those
// Contender's TotalLikesGiven. A transaction then updates each Contender
// in the map to be updated.
func UpdateHDMContenderDependentData() {
	db := GetDBHandle(NewConfig())
	posts, err := GetHDMPosts(db)
	HandleError("Could not get posts", err, true)

	// Initialize a map of Contenders to be updated
	contenders := make(map[string]Contender)

	// key, value: Post.Id, Post
	for _, p := range posts {
		// log.Println(fmt.Sprintf("Updating %s's post, %v", p.Author, p.Id))

		// Grab contender from db if it has not been updated yet
		var poster *Contender
		if val, ok := contenders[p.Author]; ok {
			poster = &val
		} else {
			poster, err = GetContenderByUsername(db, p.Author)
			if err != nil {
				// if post's author is no longer in the herp, skip it
				continue
			}
		}

		// Update Contender's data fields with Post data
		//poster.TotalPosts = append(poster.TotalPosts, p.Id)
		likesReceived := 0
		//for i := 0; i < len(poster.TotalPosts); i++ {
		//	likesReceived = len(posts[poster.TotalPosts[i]].Likes.Data) + likesReceived
		//}
		poster.TotalLikesReceived = likesReceived
		poster.AvgLikesPerPost = poster.TotalLikesReceived / len(poster.TotalPosts)
		contenders[poster.Name] = *poster

		// For each Post like, give a likes given to the contenders
		for j := 0; j < len(p.Likes.Data); j++ {
			var liker *Contender
			if val, ok := contenders[p.Likes.Data[j].Name]; ok {
				liker = &val
			} else {
				liker, _ = GetContenderByUsername(db, p.Likes.Data[j].Name)
			}

			// only update likes given for those in the herp
			if liker != nil {
				liker.TotalLikesGiven++
				contenders[liker.Name] = *liker
			}
		}
	}
	log.Println("Finished creating map of Contenders to update")

	// Update every Contender in db that was effected by Posts
	tx, err := db.Begin()
	HandleError("Failed to BEGIN txn", err, true)
	defer tx.Rollback()

	// key, value: Contender.Id, Contender
	for _, c := range contenders {
		_, err := c.UpdateContender(tx)
		if err != nil {
			HandleError("Could not update contender", err, true)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to COMMIT txn:", err)
	}
	log.Println("Updated Contender dependent data")
}

// GetContenderByUsername returns a pointer to the Contender witht he provided name.
func GetContenderByUsername(db *sql.DB, name string) (*Contender, error) {
	q := "SELECT * FROM contenders WHERE name = ?"
	var id int
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
		log.Printf(fmt.Sprintf("No user with that Name, %s.", name))
		return nil, err
	case err != nil:
		log.Fatal(fmt.Sprintf("Error getting %s from contenders table %s", name, err))
		return nil, err
	default:
		var posts []int
		json.Unmarshal([]byte(totalPosts), &posts)

		c := &Contender{
			FbId:               id,
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

// GetFBContenders returns a slice of Contenders for a given *Session from a FB group
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
			c.FbId = facebookContender.Get("id").(int) // todo: int cast needed?
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

	log.Println("Number of FB Contenders:", len(contenders))
	return contenders, nil
}
