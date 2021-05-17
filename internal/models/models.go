package models

import "time"

type Bucket struct {
	Limit      int
	CreateTime time.Time
	TTL        string
}

type Result struct {
	ResultOfCheck string
}
