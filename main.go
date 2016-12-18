package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
)

// GetAccessToken returns the access token needed to make authenticated requests
//
// Generated at https://developers.facebook.com/tools/explorer/
func GetAccessToken() string {
	var accessToken = "EAACEdEose0cBAMcn6ZCpQTQZBJ80Q1fhZBKC5ivKLSJ1ZAvnLUiax5SYt1DZBG8E7ZC8cjtHb5adrj1Y8apxtSHXWbZAZAkZCzZBNl8tZBNyYVjCnW8oZCEPLOMBunu4MJGL4agh7mop57rrDo55JAX1KMioE8S34pVuDcFVZAMKpwEry9AZDZD"
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
	Name               string
	Id                 string `facebook:",required"`
	TotalPosts         int
	TotalLikesReceived int
	AvgLikePerPost     int
	TotalLikesGiven    int
}

func main() {
	var myAccessToken = GetAccessToken()
	var herpDerpGroupID = GetGroupID()

	// "your-app-id", "your-app-secret", from 'development' app I made
	var globalApp = fb.New("756979584457445", "023c1d8f5e901c2111d7d136f5165b2a")
	session := globalApp.Session(myAccessToken)
	err := session.Validate()
	if err != nil {
		fmt.Println("Error validating session:", err)
	}

	// Get list of all members
	// response is a map[string]interface{}
	response, err := fb.Get(fmt.Sprintf("/%s/members", herpDerpGroupID), fb.Params{
		// "fields":       "feed",
		"access_token": myAccessToken,
	})
	if err != nil {
		fmt.Println("Error when getting group members:", err)
	}

	// Create the paging object for /members response
	paging, err := response.Paging(session)
	if err != nil {
		fmt.Println("Error when generating the responses Paging object:", err)
	}

	var contenders []Contender

	for {
		results := paging.Data()

		// map[administrator:false name:Jacob Glowacki id:1822807864675176]
		var c Contender

		for i := 0; i < len(results); i++ {
			results[i].Decode(&c)
			contenders = append(contenders, c)
			fmt.Println("Contender added", c)
		}

		noMore, err := paging.Next()
		if err != nil {
			fmt.Println("Error when accessing responses Next in loop:", err)
		}
		if noMore {
			break
		}
	}

	fmt.Println("number of members:", len(contenders))

	// Extract post data for users
	// example API call on post gotten from feed
	// 208678979226870_1036221373139289?fields=attachments,comments,likes,from,description,created_time,name,picture,status_type,type,caption
}
