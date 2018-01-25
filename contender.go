package main

import (
	"time"
)

type Contender struct {
	FbId               int       `json:"fb_id" facebook:",required"`
	FbGroupId          int       `json:"fb_group_id"`
	Name               string    `json:"name"`
	Posts              []string  `json:"posts"`
	AvgLikesPerPost    float64   `json:"avg_likes_per_post"` // todo: or float32?
	TotalLikesReceived int       `json:"total_likes_received"`
	TotalLikesGiven    int       `json:"total_likes_given"`
	PostsUsed          []string  `json:"posts_used"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// SetCreatedAt will set the CreatedAt attribute of a User struct
func (c *Contender) SetCreatedAt(t time.Time) {
	c.CreatedAt = t
}

// SetUpdatedAt will set the UpdatedAt attribute of a User struct
func (c *Contender) SetUpdatedAt(t time.Time) {
	c.UpdatedAt = t
}

// /////////////////
// old methods
// /////////////////

// Sort interface, http://stackoverflow.com/questions/19946992/sorting-a-map-of-structs-golang
type contenderSlice []*Contender

// Len is part of sort.Interface.
func (c contenderSlice) Len() int {
	return len(c)
}

// Swap is part of sort.Interface.
func (c contenderSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less is part of sort.Interface. Use AvgLikesPerPost as the value to sort by
func (c contenderSlice) Less(i, j int) bool {
	return c[i].AvgLikesPerPost > c[j].AvgLikesPerPost
}
