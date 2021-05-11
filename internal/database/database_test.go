package database_test

import (
	"context"
	"testing"

	"github.com/gadzira/anti-bruteforce/internal/database"
	"github.com/gadzira/anti-bruteforce/internal/logger"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type DataBaseSuite struct {
	suite.Suite
	e *database.Entry
}

func (suite *DataBaseSuite) SetupTest() {
	suite.e = &database.Entry{
		IP:   "172.0.0.1",
		Mask: "255.255.255.0",
		List: "black",
	}

}

func (suite *DataBaseSuite) TestAddToListSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	logg := l.InitLogger()
	sql := database.New(logg)
	err := sql.Connect(ctx, "host=localhost port=5432 user=postgres password=dbpass sslmode=disable")
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	err = sql.AddToList(ctx, suite.e)
	suite.NoErrorf(err, "expected no error, but got %w", err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DataBaseSuite))
}

func TestAddToList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	logg := l.InitLogger()
	sql := database.New(logg)
	err := sql.Connect(ctx, "host=localhost port=5432 user=postgres password=dbpass sslmode=disable")
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	e := &database.Entry{
		IP:   "172.0.0.1",
		Mask: "255.255.255.0",
		List: "black",
	}
	err = sql.AddToList(ctx, e)
	require.NoError(t, err, "expected no error, but got %w", err)
}

func TestCheckInList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	logg := l.InitLogger()
	sql := database.New(logg)
	err := sql.Connect(ctx, "host=localhost port=5432 user=postgres password=dbpass sslmode=disable")
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	ok, err := sql.CheckInList(ctx, "172.0.0.1", "black")
	require.True(t, ok, "expected true, but got %v", ok)
	require.NoError(t, err, "expected no error, but got %w", err)
}

func TestRemoveFromList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	logg := l.InitLogger()
	sql := database.New(logg)
	err := sql.Connect(ctx, "host=localhost port=5432 user=postgres password=dbpass sslmode=disable")
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}
	e := &database.Entry{
		IP:   "172.0.0.1",
		Mask: "255.255.255.0",
		List: "black",
	}
	err = sql.RemoveFromList(ctx, e)
	require.NoError(t, err, "expected no error, but got %w", err)
}

func TestCheckInListNegative(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	logg := l.InitLogger()
	sql := database.New(logg)
	err := sql.Connect(ctx, "host=localhost port=5432 user=postgres password=dbpass sslmode=disable")
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	ok, err := sql.CheckInList(ctx, "172.0.0.1", "black")
	require.False(t, ok, "expected false, but got %v", ok)
	require.NoError(t, err, "expected no error, but got %w", err)
}
