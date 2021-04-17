package app

import (
	"github.com/gadzira/anti-bruteforce/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	logger  *zap.Logger
	Storage *storage.StorageOfBuckets
}

type Storage interface {
}

func New(l *zap.Logger, s *storage.StorageOfBuckets) *App {
	return &App{
		logger:  l,
		Storage: s,
	}
}

// func (a *App) AddNewBucket(s string) error {
// 	a.storage.AddNewBucket(s)
// 	return nil
// }
