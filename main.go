package main

import (
	"fmt"
	fb "github.com/huandu/facebook"
)

// GetAccessToken returns the access token needed to make authenticated requests
//
// Generated at https://developers.facebook.com/tools/explorer/
func GetAccessToken() string {
	var accessToken = "EAACEdEose0cBAIBSK5qPhnzdTZCfgAK5pM4MbseqiBIPOR5ZBT1hhrECg4vP06E1JK2aPs3schlDtVStzpLzoJmAvqiUtiDfdBobnR1ivx16tAifdOaziy1HNfqUJ9FBzfwGl8J2zU2o2ZC6QZCJiLKuTBr5jI335HvQojx6ugZDZD"
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

func main() {
	var myAccessToken = GetAccessToken()
	var _ = GetUserID()
	var herpDerpGroupID = GetGroupID()

	// Get list of all members
	res, err := fb.Get(fmt.Sprintf("/%s/members", herpDerpGroupID), fb.Params{
		// "fields":       "feed",
		"access_token": myAccessToken,
	})
	if err != nil {
		fmt.Println("Error when getting group members: ", err)
	}

	// START
	// figure out app and session
	// https://github.com/huandu/facebook#use-app-and-session

	// Create the paging object
	paging, _ := res.Paging(session)

	// Get current results
	results := paging.Data()
	fmt.Println("first list of members:", results)

	// get next page.
	noMore, err := paging.Next()
	fmt.Println("more results?:", noMore)

	results = paging.Data()
	fmt.Println("second list of members:", results)


	// // Get number of group members
	// fmt.Println("Herp Derp members length:", len(res["data"].([]interface{})))
	// for k, _ := range res {
	// 	fmt.Println(k)
	// }

	// Extract post data for users
	// example API call on post gotten from feed
	// 208678979226870_1036221373139289?fields=attachments,comments,likes,from,description,created_time,name,picture,status_type,type,caption
}
