package storage

import (
	"fmt"
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

func (s *OfBuckets) CheckRequest(log, pass, ip string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.GarbageCollector()

	_, err := s.checkIncomingParameters(log, s.N)
	if err != nil {
		return false, err
	}

	_, err = s.checkIncomingParameters(pass, s.M)
	if err != nil {
		return false, err
	}

	_, err = s.checkIncomingParameters(ip, s.K)
	if err != nil {
		return false, err
	}

	// var incomingParams = map[string]string{"log": log, "pass": pass, "ip": ip}
	// for k, v := range incomingParams {
	// 	switch k {
	// 	case "log":
	// 		updateableBucket, ok := s.Bucket[v]
	// 		if !ok {
	// 			s.Bucket[v] = CreateBucket(s.N, s.TTL)
	// 		} else {
	// 			if updateableBucket.Limit == 0 {
	// 				return false, nil
	// 			}
	// 			if updateableBucket.Limit != 0 && inTimeSpan(updateableBucket.CreateTime, updateableBucket.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
	// 				updateableBucket.Limit--
	// 			}
	// 		}
	// 	case "pass":
	// 		updateableBucket, ok := s.Bucket[v]
	// 		if !ok {
	// 			s.Bucket[v] = CreateBucket(s.M, s.TTL)
	// 		} else {
	// 			if updateableBucket.Limit == 0 {
	// 				return false, nil
	// 			}
	// 			if updateableBucket.Limit != 0 && inTimeSpan(updateableBucket.CreateTime, updateableBucket.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
	// 				updateableBucket.Limit--
	// 			}
	// 		}
	// 	case "ip":
	// 		updateableBucket, ok := s.Bucket[v]
	// 		if !ok {
	// 			s.Bucket[v] = CreateBucket(s.K, s.TTL)
	// 		} else {
	// 			if updateableBucket.Limit == 0 {
	// 				return false, nil
	// 			}
	// 			if updateableBucket.Limit != 0 && inTimeSpan(updateableBucket.CreateTime, updateableBucket.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
	// 				updateableBucket.Limit--
	// 			}
	// 		}
	// 	default:
	// 		s.Log.Panic("unexpected case")
	// 	}
	// }

	// TODO: remove later
	for k, v := range s.Bucket {
		fmt.Printf("KEY:%s\t VALUE:%v\n", k, v)
	}
	fmt.Println()

	return true, nil
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

func (s *OfBuckets) checkIncomingParameters(key string, i int) (bool, error) {
	b, ok := s.Bucket[key]
	if !ok {
		s.Bucket[key] = CreateBucket(i, s.TTL)
		return true, nil
	} else {
		if b.Limit == 0 {
			return false, nil
		}
		if b.Limit != 0 && inTimeSpan(b.CreateTime, b.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
			b.Limit--
			// return true, nil
		}
	}
	return true, nil
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

// func (s *StorageOfBuckets) processIncommingParams(b *models.Bucket) (bool, error) {
// 	fmt.Println("Limit:", b.Limit)
// 	if b.Limit == 0 {
// 		return false, nil
// 	}
// 	if b.Limit != 0 && inTimeSpan(b.CreateTime, b.CreateTime.Add(time.Minute*1), time.Now().UTC()) {
// 		b.Limit -= 1
// 	}
// 	return true, nil
// }
