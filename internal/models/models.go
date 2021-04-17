package models

import "time"

type Bucket struct {
	NС         int
	MС         int
	KС         int
	CreateTime time.Time
	TTL        float64
}
