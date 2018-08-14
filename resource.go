package main

import "time"

type Resource interface {
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}
