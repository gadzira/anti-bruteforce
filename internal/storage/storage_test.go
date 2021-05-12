package storage_test

import (
	"testing"

	"github.com/gadzira/anti-bruteforce/internal/logger"
	"github.com/gadzira/anti-bruteforce/internal/models"
	"github.com/gadzira/anti-bruteforce/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestCreateBucket(t *testing.T) {
	testStorage := make(map[string]*models.Bucket)
	testStorage["SomeLogin"] = storage.CreateBucket(10, "10m")
	testBucket := testStorage["SomeLogin"]
	require.Equal(t, 9, testBucket.Limit, "expected Limit 10, actually: %d", testBucket.Limit)
	require.Equal(t, "10m", testBucket.TTL, "expected TTL 10m, actually: %d", testBucket.TTL)
}

func TestResetBucket(t *testing.T) {
	logg := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	s := storage.New(10, 100, 1000, "12m", logg.InitLogger())
	s.Bucket["SomeLogin"] = storage.CreateBucket(10, "10m")
	s.ResetBucket("SomeLogin")
	require.Equal(t, 0, len(s.Bucket), "After reset expected valuee is 0 (zero), actually: %d", len(s.Bucket))
}
