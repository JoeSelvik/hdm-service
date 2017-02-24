package main

type Contender struct {
	Id                 string `facebook:",required"`
	Name               string
	TotalPosts         int
	TotalLikesReceived int
	AvgLikePerPost     int
	TotalLikesGiven    int
}
