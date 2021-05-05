package app

import (
	"context"

	"github.com/gadzira/anti-bruteforce/internal/database"
	"github.com/gadzira/anti-bruteforce/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	Ctx     context.Context
	Logger  *zap.Logger
	Storage *storage.OfBuckets
	DB      *database.DataBase
}

func New(ctx context.Context, l *zap.Logger, s *storage.OfBuckets, db *database.DataBase) *App {
	return &App{
		Ctx:     ctx,
		Logger:  l,
		Storage: s,
		DB:      db,
	}
}
