package storage

import (
	"fmt"
	"sync"
	"time"

	"github.com/gadzira/anti-bruteforce/internal/models"
	"go.uber.org/zap"
)

type StorageOfBuckets struct {
	bucket  map[string]*models.Bucket
	mu      sync.RWMutex
	log     *zap.Logger
	N, M, K int
	TTL     string
}

func New(n, m, k int, ttl string, l *zap.Logger) StorageOfBuckets {
	return StorageOfBuckets{
		bucket: map[string]*models.Bucket{},
		log:    l,
		N:      n,
		M:      m,
		K:      k,
		TTL:    ttl,
	}
}

func (s *StorageOfBuckets) CheckRequest(log, pass, ip string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var incomingParams = map[string]string{"log": log, "pass": pass, "ip": ip}

	s.GarbageCollector()

	for k, v := range incomingParams {
		switch k {
		case "log":
			updateableBucket, ok := s.bucket[v]
			if !ok {
				s.bucket[v] = createBucket(s.N, s.TTL)
			} else {
				// TODO: Add GC for all cases
				// If limit != 0 AND req in one minutes from bucket CreateTime - reduce limit
				if updateableBucket.Limit != 0 && inTimeSpan(updateableBucket.CreateTime, updateableBucket.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
					updateableBucket.Limit -= 1
				}
				if updateableBucket.Limit == 0 {
					return false, nil
				}
				// but req not in time span - we will await GC
			}
		case "pass":
			updateableBucket, ok := s.bucket[v]
			if !ok {
				s.bucket[v] = createBucket(s.M, s.TTL)
			} else {
				if updateableBucket.Limit != 0 && inTimeSpan(updateableBucket.CreateTime, updateableBucket.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
					updateableBucket.Limit -= 1
				} else {
					return false, nil
				}
			}
		case "ip":
			updateableBucket, ok := s.bucket[v]
			if !ok {
				s.bucket[v] = createBucket(s.K, s.TTL)
			} else {
				if updateableBucket.Limit != 0 && inTimeSpan(updateableBucket.CreateTime, updateableBucket.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
					updateableBucket.Limit -= 1
				} else {
					return false, nil
				}
			}
		default:
			s.log.Panic("unexpected case")
		}
	}
	// TODO: remove later
	for k, v := range s.bucket {
		fmt.Printf("KEY:%s\t VALUE:%v\n", k, v)
	}
	fmt.Println()
	return true, nil
}

func (s *StorageOfBuckets) ResetBucket(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.bucket, key)
}

// garbage collector which we deserve, actually
func (s *StorageOfBuckets) GarbageCollector() {
	if len(s.bucket) != 0 {
		for i, j := range s.bucket {
			ct := time.Now().UTC()
			ttl, _ := time.ParseDuration(j.TTL)
			itemForDelete, ok := s.bucket[i]
			if ok && !inTimeSpan(itemForDelete.CreateTime, itemForDelete.CreateTime.Add(ttl), ct) {
				delete(s.bucket, i)
			}
		}
	}
}

func createBucket(l int, ttl string) *models.Bucket {
	return &models.Bucket{
		Limit:      l - 1,
		TTL:        ttl,
		CreateTime: time.Now().UTC(),
	}
}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// func findEndTime(ct time.Time, s string) time.Time {
// ttl, _ := time.ParseDuration(s)
// return ct.Add(ttl)
// }
