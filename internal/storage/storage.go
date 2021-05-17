package storage

import (
	"sync"
	"time"

	"github.com/gadzira/anti-bruteforce/internal/models"
	"go.uber.org/zap"
)

type OfBuckets struct {
	Bucket  map[string]*models.Bucket
	mu      sync.RWMutex
	Log     *zap.Logger
	N, M, K int
	TTL     string
}

func New(n, m, k int, ttl string, l *zap.Logger) OfBuckets {
	return OfBuckets{
		Bucket: map[string]*models.Bucket{},
		Log:    l,
		N:      n,
		M:      m,
		K:      k,
		TTL:    ttl,
	}
}

func (s *OfBuckets) CheckRequest(log, pass, ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.GarbageCollector()

	ok := s.checkIncomingParameters(log, s.N)
	if !ok {
		return ok
	}

	ok = s.checkIncomingParameters(pass, s.M)
	if !ok {
		return ok
	}

	ok = s.checkIncomingParameters(ip, s.K)
	if !ok {
		return ok
	}

	// TODO: remove later
	// for k, v := range s.Bucket {
	// 	fmt.Printf("KEY:%s\t VALUE:%v\n", k, v)
	// }
	// fmt.Println()

	return true
}

func (s *OfBuckets) ResetBucket(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Bucket, key)
}

func (s *OfBuckets) ShowBuckets() map[string]*models.Bucket {
	return s.Bucket
}

func (s *OfBuckets) GarbageCollector() {
	if len(s.Bucket) != 0 {
		for i, j := range s.Bucket {
			ct := time.Now().UTC()
			ttl, _ := time.ParseDuration(j.TTL)
			itemForDelete, ok := s.Bucket[i]
			if ok && !inTimeSpan(itemForDelete.CreateTime, itemForDelete.CreateTime.Add(ttl), ct) {
				delete(s.Bucket, i)
			}
		}
	}
}

func (s *OfBuckets) checkIncomingParameters(key string, i int) bool {
	b, ok := s.Bucket[key]
	if !ok {
		s.Bucket[key] = CreateBucket(i, s.TTL)
		return true
	}
	if b.Limit == 0 {
		return false
	}
	if b.Limit != 0 && inTimeSpan(b.CreateTime, b.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
		b.Limit--
	}

	return true
}

func CreateBucket(l int, ttl string) *models.Bucket {
	return &models.Bucket{
		Limit:      l - 1,
		TTL:        ttl,
		CreateTime: time.Now().UTC(),
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}
