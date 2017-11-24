package main

import "time"

// todo: get and set id methods?
type Resource interface {
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}
