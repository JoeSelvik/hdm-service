package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetFbSession returns the pointer to a fb Session object.
//
// Panics if interacting with FB does not work.
// todo: learn how this works
func getFbSession() *fb.Session {
	// "your-app-id", "your-app-secret", from 'development' app I made
	var globalApp = fb.New("756979584457445", "023c1d8f5e901c2111d7d136f5165b2a")
	session := globalApp.Session(Config.FbAccessToken)
	err := session.Validate()
	if err != nil {
		panic(err)
	}

	return session
}

// PullContendersFromFb returns a slice of pointers to Contenders for a given *Session from a FB group
func PullContendersFromFb() ([]*Contender, *ApplicationError) {
	// response is a map[string]interface{}
	response, err := fb.Get(fmt.Sprintf("/%d/members", Config.FbGroupId), fb.Params{
		"access_token": Config.FbAccessToken,
		"fields":       []string{"name", "id"},
	})
	if err != nil {
		msg := "Failed to get group members from fb"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	// Get the member's paging object
	session := getFbSession()
	paging, err := response.Paging(session)
	if err != nil {
		msg := "Failed to page on the group members response"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	var contenders []*Contender
	for {
		results := paging.Data()

		// map[administrator:false name:Jacob Glowacki id:1822807864675176]
		for i := 0; i < len(results); i++ {
			var c Contender
			facebookContender := fb.Result(results[i]) // cast the var

			// Convert interface to it's real string value, then the string to an int.
			id, err := strconv.Atoi(facebookContender.Get("id").(string))
			if err != nil {
				msg := "Failed to convert fb contenders id to a string"
				return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
			}

			c.FbId = id
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

// PullPostsFromFb returns a slice of Posts from the Group feed up to a given date.
//
// todo: use start and end times
func PullPostsFromFb(startDate time.Time) ([]Post, *ApplicationError) {
	// Get the group feed
	response, err := fb.Get(fmt.Sprintf("/%d/feed", Config.FbGroupId), fb.Params{
		"access_token": Config.FbAccessToken,
		"fields":       []string{"from", "created_time", "likes"},
	})
	if err != nil {
		msg := "Failed to get group feed"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	// Get the feed's paging object
	session := getFbSession()
	paging, err := response.Paging(session)
	if err != nil {
		msg := "Failed to page on the group feed response"
		return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
	}

	var posts []Post
	count := 1

	// loop until a fb post's created_time is older than startDate
Loop:
	for {
		results := paging.Data()

		// 25 posts per page, load data into a Post struct
		for i := 0; i < len(results); i++ {
			var p Post
			facebookPost := fb.Result(results[i]) // cast the var

			// stop when post reaches startDate
			p.PostedDate = facebookPost.Get("created_time").(string)
			t, err := time.Parse(GoTimeLayout, p.PostedDate)
			if err != nil {
				msg := "Failed to parse a fb post's postedDate"
				return nil, &ApplicationError{Msg: msg, Err: err, Code: http.StatusInternalServerError}
			}
			if t.Before(startDate) {
				break Loop
			}

			p.Id = facebookPost.Get("id").(string)
			p.Author = facebookPost.Get("from.name").(string)
			p.PostedDate = t.String()

			// unload Likes data into a Like struct
			if facebookPost.Get("likes.data") != nil {
				var like_list []Like
				numLikes := facebookPost.Get("likes.data").([]interface{})
				for j := 0; j < len(numLikes); j++ {
					var l Like
					l.Id = numLikes[j].(map[string]interface{})["id"].(string)
					l.Name = numLikes[j].(map[string]interface{})["name"].(string)
					like_list = append(like_list, l)
				}
				p.Likes = Likes{Data: like_list}

			} else {
				p.Likes = Likes{Data: nil}
			}

			// save the new Post
			posts = append(posts, p)
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
		count++
	}
	return posts, nil
}
