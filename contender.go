package main

import "time"

type Contender struct {
	Id                 string `facebook:",required"`
	Name               string
	TotalPosts         []string
	TotalLikesReceived int
	AvgLikePerPost     int
	TotalLikesGiven    int
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
}
