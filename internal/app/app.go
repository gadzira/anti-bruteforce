package app

import (
	"context"

	"github.com/gadzira/anti-bruteforce/internal/database"
	"github.com/gadzira/anti-bruteforce/internal/domain"
	"github.com/gadzira/anti-bruteforce/internal/models"
	"github.com/gadzira/anti-bruteforce/internal/storage"
	"go.uber.org/zap"
)

type Application struct {
	Ctx     context.Context
	Logger  *zap.Logger
	DB      domain.DataBaseInterface
	storage domain.StorageInterface
}

func New(ctx context.Context, logger *zap.Logger, db database.DataBase, storage *storage.OfBuckets) *Application {
	return &Application{
		Ctx:     ctx,
		Logger:  logger,
		DB:      &db,
		storage: storage,
	}
}

func (a *Application) ShowBucket() map[string]*models.Bucket {
	return a.storage.ShowBuckets()
}

func (a *Application) CheckRequest(log, pass, ip string) bool {
	return a.storage.CheckRequest(log, pass, ip)
}

func (a *Application) ResetBucket(key string) {
	a.storage.ResetBucket(key)
}
