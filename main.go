package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
)

// TODO: modularize, ie ask for what kind of acess and for which user
func GetAccessToken() string {
	var accessToken = "CAACEdEose0cBANZBMXBYPhpegCUyx3B54uCJNFnOCyhnOAhqNteXwEwCTWZAYMW6HZCGLg4sHp83gVap1uiKKbz6Wjvh8iJf57LoZC1Tol3ZBZA3CKAWgERwxwANF4ghpObWnSFZAj8SszlRtUOJ6hfj4FZAIpYP7if4tBqWz7o2ByZCXquHGqIgWrP2BU1WcxpDB9iAtALp6Deizgb704i96"
	return accessToken
}

// GetUserID will return the user's id associated with the access token provided for the app.
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

// TODO: modularize, ie ask for which group
func GetGroupID() string {
	var groupID = "208678979226870"
	return groupID
}

func main() {
	var myAccessToken = GetAccessToken()
	var _ = GetUserID()
	var herpDerpGroupID = GetGroupID()

	// Start parsing through the feed!
	// Params: can specify which fields you want back
	res, err := fb.Get(fmt.Sprintf("/%s/feed", herpDerpGroupID), fb.Params{
		// "fields":       "feed",
		"access_token": myAccessToken,
	})
	if err != nil {
		fmt.Println("Error when getting group feed: ", err)
	}

	// grab 'data' for a map of posts
	fmt.Println("Feed length:", len(res))
	for k, _ := range res {
		fmt.Println(k)
	}

	// Extract post data for users
	// example API call on post gotten from feed
	// 208678979226870_1036221373139289?fields=attachments,comments,likes,from,description,created_time,name,picture,status_type,type,caption
}
