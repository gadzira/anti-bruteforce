package domain

import (
	"context"

	"github.com/gadzira/anti-bruteforce/internal/database"
	"github.com/gadzira/anti-bruteforce/internal/models"
)

type DataBaseInterface interface {
	Connect(ctx context.Context, dsn string) error
	Close() error
	AddToList(ctx context.Context, e *database.Entry) error
	RemoveFromList(ctx context.Context, e *database.Entry) error
	CheckInList(ctx context.Context, ip string, list string) (bool, error)
}

type StorageInterface interface {
	CheckRequest(log, pass, ip string) bool
	ResetBucket(key string)
	ShowBuckets() map[string]*models.Bucket
}
