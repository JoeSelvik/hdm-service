package main

import (
	// "database/sql"
	// "encoding/json"
	// "log"
	"time"
)

type Matchup struct {
	Id         int
	Name       string // ie: firstRound_g0
	ContenderA Contender
	APosts     []Post
	AVotes     int
	ContenderB Contender
	BPosts     []Post
	BVotes     int
	InProgress bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m *Matchup) DBTableName() string {
	return "matchups"
}

func (m *Matchup) Path() string {
	return "/matchups/"
}
