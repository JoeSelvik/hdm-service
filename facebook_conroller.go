// Not unit tested because there are only two calls to the fb SDK library. The code is primarily
// parsing logic around the custom types that are returned from the dependency. Would
// have to implement fake interfaces for any functionality called on fb; paging and result.

package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Facebooker interface {
	PullContendersFromFb() ([]*Contender, *ApplicationError)
	PullPostsFromFb() ([]*Post, *ApplicationError)
}

type FacebookHandle struct {
	config *Configuration
}

// getFbSession returns the pointer to a fb Session object.
//
// Panics if interacting with FB fails.
func (fh FacebookHandle) getFbSession() *fb.Session {
	// "your-app-id", "your-app-secret", from 'development' app I made
	var globalApp = fb.New("756979584457445", "023c1d8f5e901c2111d7d136f5165b2a")
	session := globalApp.Session(fh.config.FbAccessToken)
	err := session.Validate()
	if err != nil {
		panic(err)
	}

	return session
}

// PullContendersFromFb returns a slice of pointers to Contenders from a FB group set in the configuration.
func (fh FacebookHandle) PullContendersFromFb() ([]*Contender, *ApplicationError) {
	// request members via fb graph api, response is a map[string]interface{} fb.Result
	response, err := fb.Get(fmt.Sprintf("/%d/members", fh.config.FbGroupId), fb.Params{
		"access_token": fh.config.FbAccessToken,
		"fields":       []string{"name", "id"},
	})
	if err != nil {
		msg := "Failed to get group members from fb"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	// get the member's paging object
	session := fh.getFbSession()
	paging, err := response.Paging(session)
	if err != nil {
		msg := "Failed to page on the group members response"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	// create each Contender from the facebook response
	var contenders []*Contender
	for {
		results := paging.Data()

		// ie - map[administrator:false name:Jack White id:6666666666666666]
		for i := 0; i < len(results); i++ {
			var c Contender
			facebookContender := fb.Result(results[i]) // cast the var

			// convert interface to its real string value, then the string to an int.
			id, err := strconv.Atoi(facebookContender.Get("id").(string))
			if err != nil {
				msg := "Failed to convert fb contenders id to a string"
				return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
			}

			c.FbId = id
			c.FbGroupId = fh.config.FbGroupId
			c.Name = facebookContender.Get("name").(string)
			contenders = append(contenders, &c)
		}

		noMore, err := paging.Next()
		if err != nil {
			msg := "Failed to get next paging object for members"
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}
		if noMore {
			break
		}
	}

	return contenders, nil
}

// PullPostsFromFb returns a slice of Post pointers from the configured group feed up to a given date.
func (fh FacebookHandle) PullPostsFromFb() ([]*Post, *ApplicationError) {
	// request group feed, response is a map[string]interface{} fb.Result
	response, err := fb.Get(fmt.Sprintf("/%d/feed", fh.config.FbGroupId), fb.Params{
		"access_token": fh.config.FbAccessToken,
		"fields":       []string{"from", "created_time", "likes"},
	})
	if err != nil {
		msg := "Failed to get group feed"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	// get the feed's paging object
	session := fh.getFbSession()
	paging, err := response.Paging(session)
	if err != nil {
		msg := "Failed to page on the group feed response"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	// create each Post from the facebook response
	var posts []*Post

	// loop until a fb post's created_time is older than config.StartTime
	// todo: clean up messy code
Loop:
	for {
		results := paging.Data()

		// 25 posts per page, load data into a Post struct
		for i := 0; i < len(results); i++ {
			var p Post
			facebookPost := fb.Result(results[i]) // cast the var

			// parse post's created_time
			t, err := time.Parse(GoTimeLayout, facebookPost.Get("created_time").(string))
			if err != nil {
				msg := "Failed to parse a fb post's postedDate"
				return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
			}

			// continue until post is after EndTime
			if t.After(fh.config.EndTime) {
				continue
			}
			// stop when post reaches startDate
			if t.Before(fh.config.StartTime) {
				break Loop
			}

			// convert interface to it's real string value, then the string to an int.
			id, err := strconv.Atoi(facebookPost.Get("from.id").(string))
			if err != nil {
				msg := fmt.Sprintf("Failed to convert fb contenders id to a string: %s", err)
				return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
			}

			p.FbId = facebookPost.Get("id").(string) // a post's id has an underscore, 111_222
			p.FbGroupId = fh.config.FbGroupId
			p.AuthorFbId = id
			p.PostedDate = t

			// extract fb_ids of contenders who liked post
			if facebookPost.Get("likes.data") != nil {
				postLikes := facebookPost.Get("likes.data").([]interface{})
				for _, l := range postLikes {
					// Convert interface to its real string value, then the string to an int.
					lid, err := strconv.Atoi(l.(map[string]interface{})["id"].(string))
					if err != nil {
						msg := "Failed to convert a posts liker id to an int"
						return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
					}
					p.Likes = append(p.Likes, lid)
				}
			} else {
				p.Likes = []int{}
			}

			// save the new Post
			posts = append(posts, &p)
		}

		noMore, err := paging.Next()
		if err != nil {
			msg := "Failed to get next paging object for posts"
			return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
		}
		if noMore {
			log.Println("Reached the end of group feed")
			break Loop
		}
	}
	return posts, nil
}
