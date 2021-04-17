package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/gadzira/anti-bruteforce/internal/models"
	"go.uber.org/zap"
)

type StorageOfBuckets struct {
	mu     sync.RWMutex
	bucket map[string]models.Bucket
	log    *zap.Logger
	NC     int
	MC     int
	KC     int
	TTL    string
}

func New(N, M, K int, TTL string, l *zap.Logger) StorageOfBuckets {
	return StorageOfBuckets{
		bucket: make(map[string]models.Bucket),
		log:    l,
		NC:     N,
		MC:     M,
		KC:     K,
		TTL:    TTL,
	}
}

func (s *StorageOfBuckets) AddNewBucket(str string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ttl, _ := time.ParseDuration(s.TTL)
	d := ttl.Minutes()

	b := models.Bucket{
		NС:         s.NC,
		MС:         s.MC,
		KС:         s.KC,
		CreateTime: time.Now().Local(),
		TTL:        d,
	}
	s.bucket[str] = b
	s.log.Info("Added new bucket!")
	fmt.Println(s.bucket)
	// Откуда брать верхние значения счетчиков? конфиг или параметры?
	// Для meta-части должны братся из конфига
	// Для ключа - явно будут приходить снаружи
	return nil
}

func IfLogIsPresent() bool {
	return false
}
