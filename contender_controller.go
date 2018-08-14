package main

import (
	"database/sql"
	"fmt"
	"github.com/JoeSelvik/hdm-service/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ContenderController struct {
	config *Configuration
	db     *models.DB
	fh     Facebooker
}

// ServeHTTP routes incoming requests to the right service.
func (cc *ContenderController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := new(Contender)
	ServeResource(w, r, cc, c)
}

// Path returns the URL extension associated with the Contender resource.
func (cc *ContenderController) Path() string {
	return "/contenders/"
}

// DBTableName returns the table name for Contenders.
func (cc *ContenderController) DBTableName() string {
	return "contenders"
}

// Create writes a new contender to the db for each given Resource.
func (cc *ContenderController) Create(m []models.Resource) ([]int, *ApplicationError) {
	// create a slice of Contender pointers by asserting on a slice of Resources interfaces
	var contenders []*Contender
	for i := 0; i < len(m); i++ {
		c := m[i]
		contenders = append(contenders, c.(*Contender))
	}

	// Create the SQL query
	q := fmt.Sprintf(`
	INSERT INTO %s (
		fb_id, fb_group_id,
		name, posts, avg_likes_per_post, total_likes_received, total_likes_given, posts_used,
		created_at, updated_at
	) values (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, cc.DBTableName())

	// begin sql transaction
	tx, err := cc.db.Begin()
	if err != nil {
		msg := "Something is wrong with our database - we'll be back up soon!"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}
	defer tx.Rollback()

	// insert each Contender into contenders table
	var contenderIds []int
	for _, c := range contenders {
		posts := strings.Join(c.Posts[:], ", ")
		postsUsed := strings.Join(c.PostsUsed[:], ", ")

		result, err := tx.Exec(q,
			c.FbId, c.FbGroupId,
			c.Name, posts, c.AvgLikesPerPost, c.TotalLikesReceived, c.TotalLikesGiven, postsUsed)
		if err != nil {
			msg := fmt.Sprintf("Couldn't create contender: %+v", c)
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		// save each Id to return
		id, err := result.LastInsertId()
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}
		contenderIds = append(contenderIds, int(id))
	}

	// commit sql transaction
	if err = tx.Commit(); err != nil {
		msg := "Something is wrong with our database - we'll be back up soon!"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	return contenderIds, nil
}

// Read returns the contender in the db for a given FbId.
func (cc *ContenderController) Read(fbId int) (models.Resource, *ApplicationError) {
	// todo: find better way to shorten lines of code and reuse in ReadCollection
	var fbGroupId int
	var name string
	var totalPostsString string
	var avgLikesPerPost float64
	var totalLikesReceived int
	var totalLikesGiven int
	var postsUsedString string
	var createdAt time.Time
	var updatedAt time.Time

	// grab contender entry from table
	q := fmt.Sprintf("SELECT * FROM %s WHERE fb_id=%d", cc.DBTableName(), fbId)
	err := cc.db.QueryRow(q).Scan(&fbId, &fbGroupId, &name, &totalPostsString, &avgLikesPerPost, &totalLikesReceived,
		&totalLikesGiven, &postsUsedString, &createdAt, &updatedAt) // todo: okay to unscan into fbId arg?
	switch {
	case err == sql.ErrNoRows:
		msg := fmt.Sprintf("Couldn't find any resource with id: %d", fbId)
		return nil, &ApplicationError{Msg: msg, Code: http.StatusNotFound}
	case err != nil:
		msg := "Something is wrong with our database - we'll be back up soon!"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	// todo: find better way to abstract unloading strings of ints and creating individual contender and ReadCollection
	// split comma separated strings to slices of ints
	totalPosts := strings.Split(totalPostsString, ", ")
	postsUsed := strings.Split(postsUsedString, ", ")

	// create contender
	c := Contender{
		FbId:               fbId,
		FbGroupId:          fbGroupId,
		Name:               name,
		Posts:              totalPosts,
		AvgLikesPerPost:    float64(avgLikesPerPost),
		TotalLikesReceived: totalLikesReceived,
		TotalLikesGiven:    totalLikesGiven,
		PostsUsed:          postsUsed,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
	}

	return &c, nil
}

// Update writes the db column value for each variable Contender parameter.
//
// Writes Posts, AvgLikesPerPost, TotalLikesReceived, TotalLikesGiven, PostsUsed, and UpdatedAt.
func (cc *ContenderController) Update(m []models.Resource) *ApplicationError {
	// Create a slice of Contender pointers by asserting on a slice of Resources interfaces
	var contenders []*Contender
	for i := 0; i < len(m); i++ {
		c := m[i]
		contenders = append(contenders, c.(*Contender))
	}

	// begin sql transaction
	tx, err := cc.db.Begin()
	if err != nil {
		msg := "Something is wrong with our database - we'll be back up soon!"
		return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}
	defer tx.Rollback()

	// create the SQL query
	q := fmt.Sprintf(`
	UPDATE %s SET
		posts=?, avg_likes_per_post=?, total_likes_received=?, total_likes_given=?, posts_used=?,
		updated_at=CURRENT_TIMESTAMP
		WHERE fb_id=?
	`, cc.DBTableName())

	// iterate through each contender and update it in the db
	for _, c := range contenders {
		posts := strings.Join(c.Posts[:], ", ")
		postsUsed := strings.Join(c.PostsUsed[:], ", ")

		res, err := tx.Exec(q, posts, c.AvgLikesPerPost, c.TotalLikesReceived, c.TotalLikesGiven, postsUsed, c.FbId)
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		// Note - not really sure what this can error on
		numrows, err := res.RowsAffected()
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		// if more or less than one row is affected then we have a problem
		switch {
		case numrows == 0:
			msg := fmt.Sprintf("Couldn't find any resource to update with id: %d", c.FbId)
			return &ApplicationError{Msg: msg, Code: http.StatusNotFound}
		case numrows != 1:
			// this is really bad, should never see. May be an SQL injection attempt.
			msg := "Something is wrong with our database - we'll be back up soon!"
			return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}
	}

	// commit sql transaction
	if err = tx.Commit(); err != nil {
		msg := "Something is wrong with our database - we'll be back up soon!"
		return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	return nil
}

// Destroy deletes any given Id from the db.
func (cc *ContenderController) Destroy(ids []int) *ApplicationError {
	// begin sql transaction
	tx, err := cc.db.Begin()
	if err != nil {
		msg := "Something is wrong with our database - we'll be back up soon!"
		return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}
	defer tx.Rollback()

	// breate the SQL query
	q := fmt.Sprintf("DELETE FROM %s WHERE fb_id = $1;", cc.DBTableName())

	// iterate through each contender and update it in the db
	for _, v := range ids {
		// todo: a lot of repeated code from update's error handling
		res, err := tx.Exec(q, v)
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		// not really sure what this can error on
		numrows, err := res.RowsAffected()
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		// if more or less than one row is affected then we have a problem
		switch {
		case numrows == 0:
			msg := fmt.Sprintf("Couldn't find any resource to destroy with id: %d", v)
			return &ApplicationError{Msg: msg, Code: http.StatusNotFound}
		case numrows != 1:
			// this is really bad, should never see. May be an SQL injection attempt.
			msg := "Something is wrong with our database - we'll be back up soon!"
			return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}
	}

	// commit sql transaction
	if err = tx.Commit(); err != nil {
		msg := "Something is wrong with our database - we'll be back up soon!"
		return &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	return nil
}

// ReadCollection returns all Contenders in the db.
func (cc *ContenderController) ReadCollection() ([]models.Resource, *ApplicationError) {
	// grab rows from table
	rows, err := cc.db.Query(fmt.Sprintf("SELECT * FROM %s", cc.DBTableName()))
	switch {
	case err == sql.ErrNoRows:
		log.Println("Contenders ReadCollection: no rows in table.")
		return []models.Resource{}, nil
	case err != nil:
		msg := "Something is wrong with our database - we'll be back up soon!"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}
	defer rows.Close()

	// create a contender from each row
	contenders := make([]models.Resource, 0)
	for rows.Next() {
		var fbId int
		var fbGroupId int
		var name string
		var totalPostsString string
		var avgLikesPerPost float64
		var totalLikesReceived int
		var totalLikesGiven int
		var postsUsedString string
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&fbId, &fbGroupId, &name, &totalPostsString, &avgLikesPerPost,
			&totalLikesReceived, &totalLikesGiven, &postsUsedString, &createdAt, &updatedAt)
		if err != nil {
			msg := "Something is wrong with our database - we'll be back up soon!"
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}

		// split comma separated strings to slices of ints
		totalPosts := strings.Split(totalPostsString, ", ")
		postsUsed := strings.Split(postsUsedString, ", ")

		c := Contender{
			FbId:               fbId,
			FbGroupId:          fbGroupId,
			Name:               name,
			Posts:              totalPosts,
			AvgLikesPerPost:    float64(avgLikesPerPost),
			TotalLikesReceived: totalLikesReceived,
			TotalLikesGiven:    totalLikesGiven,
			PostsUsed:          postsUsed,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
		}

		contenders = append(contenders, &c)
	}
	return contenders, nil
}

// PopulateContendersTable pulls contenders via the facebook api and enters them into the contender table.
func (cc *ContenderController) PopulateContendersTable() *ApplicationError {
	log.Println("Pulling contenders from facebook and creating in db")

	// get slice of contender struct pointers from fb
	contenders, aerr := cc.fh.PullContendersFromFb()
	if aerr != nil {
		return aerr
	}

	// convert each contender struct ptr to Resource interface
	contenderResources := make([]models.Resource, len(contenders))
	for i, v := range contenders {
		contenderResources[i] = models.Resource(v)
	}

	// populate Contenders table
	_, aerr = cc.Create(contenderResources)
	if aerr != nil {
		return aerr
	}

	return nil
}

// UpdateContendersVariableDependentData updates all the contender fields that depend on post data.
func (cc *ContenderController) UpdateContendersVariableDependentData(pc *PostController) *ApplicationError {
	log.Println("Updating contender VDD in db")

	// get a slice of posts
	postResources, aerr := pc.ReadCollection()
	if aerr != nil {
		log.Printf("Failed pc.ReadCollection: %s\n", aerr.Msg)
		return aerr
	}
	var posts []*Post
	for _, p := range postResources {
		posts = append(posts, p.(*Post))
	}

	// create a map of contenders to update for each post
	contendersToUpdate := make(map[int]Contender)
	for _, p := range posts {
		// get contender who authored post
		var author *Contender
		if val, okay := contendersToUpdate[p.AuthorFbId]; okay {
			author = &val
		} else {
			authorResource, aerr := cc.Read(p.AuthorFbId)
			if aerr != nil {
				log.Printf("Failed to read author contender: %d", p.AuthorFbId)
				msg := "Something is wrong with our database - we'll be back up soon!"
				return &ApplicationError{Msg: msg, Code: http.StatusInternalServerError}
			}
			author = authorResource.(*Contender)
		}

		// update author's vd
		author.Posts = append(author.Posts, p.FbId)
		author.TotalLikesReceived = author.TotalLikesReceived + len(p.Likes)
		contendersToUpdate[author.FbId] = *author

		// for each like, update contender's likes given
		for _, l := range p.Likes {
			var liker *Contender
			if val, okay := contendersToUpdate[l]; okay {
				liker = &val
			} else {
				likerResource, aerr := cc.Read(l)
				if aerr != nil {
					log.Printf("Failed to read liker contender: %d", l)
					msg := "Something is wrong with our database - we'll be back up soon!"
					return &ApplicationError{Msg: msg, Code: http.StatusInternalServerError}
				}
				liker = likerResource.(*Contender)
			}
			liker.TotalLikesGiven = liker.TotalLikesGiven + 1
			contendersToUpdate[liker.FbId] = *liker
		}
	}

	// grab each contender from the map of to be updated, convert to ptr
	var contenders []*Contender
	for k := range contendersToUpdate {
		c := contendersToUpdate[k]
		contenders = append(contenders, &c)
	}

	// convert each contender struct ptr to Resource interface
	var contenderResources []models.Resource
	for _, v := range contenders {
		contenderResources = append(contenderResources, models.Resource(v))
	}

	aerr = cc.Update(contenderResources)
	if aerr != nil {
		log.Printf("Failed to update contenders: %s\n%s\n", aerr.Msg, aerr.Err)
		msg := "Something is wrong with our database - we'll be back up soon!"
		return &ApplicationError{Msg: msg, Code: http.StatusInternalServerError}
	}

	return nil
}

// stringOfIntsToSliceOfInts is a helper function that converts a string of ints to a slice of ints.
//
// todo: probably a better way to do this...
// "1, 2, 3" to []int{1, 2, 3}
// "1,2,3" will throw an error
// returns []int{} if given string is ""
func stringOfIntsToSliceOfInts(s string) ([]int, error) {
	stringSlice := strings.Split(s, ", ")
	var intSlice []int
	if stringSlice[0] != "" {
		intSlice = make([]int, len(stringSlice))
		for i, v := range stringSlice {
			s, err := strconv.Atoi(v)
			if err != nil {
				log.Println("Failed to convert string of ints to slice:", err)
				return nil, err
			}
			intSlice[i] = s
		}
	}
	return intSlice, nil
}

// sliceOfIntsToString is a helper function that converts a slice of post ids to a string of ids.
//
// todo: probably a better way to do this...
// todo: recover() from panic()? is this possible?
// https://stackoverflow.com/questions/25025467/catching-panics-in-golang
func sliceOfIntsToString(slicePostIds []int) string {
	stringPosts := fmt.Sprint(slicePostIds)                     // [1 2 3] to "[1 2 3]"
	splitStringPosts := strings.Split(stringPosts, " ")         // "[1 2 3]" to ["[1 2 3]"]
	joinedStringPosts := strings.Join(splitStringPosts, ", ")   // ["[1 2 3]"] to "[1, 2, 3]"
	trimmedStringPosts := strings.Trim(joinedStringPosts, "[]") // "[1, 2, 3]" to "1, 2, 3"
	return trimmedStringPosts
}
