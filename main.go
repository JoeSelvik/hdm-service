package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
	"os"
)

// GetAccessToken returns the access token needed to make authenticated requests
//
// Generated at https://developers.facebook.com/tools/explorer/
func GetAccessToken() string {
	var accessToken = "EAACEdEose0cBACnb5IUNMBgzEUS9hYjzQZC7laNIeMelYZA6JSUZCVviBJdng8RI7WWycO9XL3ZAnu4hKzaEsVspJsZCCZCu34ra1atusWoeGucuNA74Ms9ZCvZAk4ol89zlev6aEtbneauw7yhL62MuQpDgMvZCShjpZCMMyI8sAgVQZDZD"
	return accessToken
}

// GetUserID returns the user's id associated with the access token provided for the app.
//
// May modify this in the future to just return the user map.
func GetUserID() string {
	var myAccessToken = GetAccessToken()

	res, err := fb.Get("/me", fb.Params{
		"access_token": myAccessToken,
	})
	if err != nil {
		fmt.Println("Error when accessing /me: ", err)
		os.Exit(3)
	}

	fmt.Println("User associated with access token: ", res)

	// TODO: is type assertion a bad idea here? Just handle the error from .Get?
	return res["id"].(string)
}

// GetGroupID returns the Herp Derp group_id
func GetGroupID() string {
	var groupID = "208678979226870"
	return groupID
}

type Contender struct {
	Id                 string `facebook:",required"`
	Name               string
	TotalPosts         int
	TotalLikesReceived int
	AvgLikePerPost     int
	TotalLikesGiven    int
}

type Post struct {
	Id          string `facebook:",required"`
	Title       string
	CreatedDate string
	From        string
	TotalLikes  int
	TotalShares int
}

// Returns a slice of Contenders for a given *Session
func CreateContenders(session *fb.Session) []Contender {
	// response is a map[string]interface{}
	response, err := fb.Get(fmt.Sprintf("/%s/members", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
	})
	if err != nil {
		fmt.Println("Error when getting group members:", err)
		os.Exit(3)
	}

	// Create the paging object for /members response
	paging, err := response.Paging(session)
	if err != nil {
		fmt.Println("Error when generating the members responses Paging object:", err)
		os.Exit(3)
	}

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
		if err != nil {
			fmt.Println("Error when accessing responses Next in loop:", err)
			os.Exit(3)
		}
		if noMore {
			break
		}
	}

	return contenders
}

func populateTotalPosts(contenders []Contender, session *fb.Session) {
	// Get the group feed
	response, err := fb.Get(fmt.Sprintf("/%s/feed", GetGroupID()), fb.Params{
		"access_token": GetAccessToken(),
		"feilds":       []string{"from", "created_time"},
	})
	if err != nil {
		fmt.Println("Error when getting feed:", err)
		os.Exit(3)
	}

	// Get the feed's paging object
	paging, err := response.Paging(session)
	if err != nil {
		fmt.Println("Error when generating the feed responses Paging object:", err)
		os.Exit(3)
	}

	var posts []Post
	count := 0

	// 25 posts per page
	for {
		results := paging.Data()

		// load data from each facebookPost into a Post struct
		for i := 0; i < len(results); i++ {
			var p Post
			facebookPost := fb.Result(results[i])

			from := facebookPost.Get("from.name")
			p.From = from.(string)

			createdDate := facebookPost.Get("created_time")
			p.CreatedDate = createdDate.(string)

			likesData := facebookPost.Get("likes.data")
			if likesData != nil {
				numLikes := facebookPost.Get("likes.data").([]interface{})
				p.TotalLikes = len(numLikes)
			} else {
				p.TotalLikes = 0
			}

			results[i].Decode(&p)
			posts = append(posts, p)
			fmt.Println("Decoded post:", p)
		}

		count++
		fmt.Println("finished lap:", count)

		if count >= 1 {
			fmt.Println("found first 25 posts")
			break
		}

		noMore, err := paging.Next()
		if err != nil {
			fmt.Println("Error when accessing responses Next in loop:", err)
			os.Exit(3)
		}
		if noMore {
			fmt.Println("Reached the end of the feed!")
			break
		}
	}
	fmt.Println("number posts:", len(posts))
	fmt.Println("First post:", posts[1])
}

func main() {
	var myAccessToken = GetAccessToken()
	// var herpDerpGroupID = GetGroupID()

	// "your-app-id", "your-app-secret", from 'development' app I made
	var globalApp = fb.New("756979584457445", "023c1d8f5e901c2111d7d136f5165b2a")
	session := globalApp.Session(myAccessToken)
	err := session.Validate()
	if err != nil {
		fmt.Println("Error validating session:", err)
		os.Exit(3)
	}

	contenders := CreateContenders(session)
	fmt.Println("number of members:", len(contenders))

	populateTotalPosts(contenders, session)

	// Extract post data for users
	// example API call on post gotten from feed
	// 208678979226870_1036221373139289?fields=attachments,comments,likes,from,description,created_time,name,picture,status_type,type,caption
}
