package app

import (
	"context"

	"github.com/gadzira/anti-bruteforce/internal/db"
	"github.com/gadzira/anti-bruteforce/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	Ctx     context.Context
	Logger  *zap.Logger
	Storage *storage.StorageOfBuckets
	DB      *db.DataBase
}

type Storage interface {
}

func New(ctx context.Context, l *zap.Logger, s *storage.StorageOfBuckets, db *db.DataBase) *App {
	return &App{
		Ctx:     ctx,
		Logger:  l,
		Storage: s,
		DB:      db,
	}
}
