package main

import (
	"time"
)

type Post struct {
	FbId       string    `json:"fb_id" facebook:",required"`
	FbGroupId  int       `json:"fb_group_id"`
	PostedDate time.Time `json:"posted_date"`
	AuthorFbId int       `json:"author_fb_id"`
	Likes      []int     `json:"likes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"created_at"`
}

// SetCreatedAt will set the CreatedAt attribute of a User struct
func (p *Post) SetCreatedAt(t time.Time) {
	p.CreatedAt = t
}

// SetUpdatedAt will set the UpdatedAt attribute of a User struct
func (p *Post) SetUpdatedAt(t time.Time) {
	p.UpdatedAt = t
}
