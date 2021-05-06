package domain

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

type Application interface {
	New(i interface{}) *App
}
