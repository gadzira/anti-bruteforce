package database_test

import (
	"context"
	"testing"

	"github.com/gadzira/anti-bruteforce/internal/database"
	"github.com/gadzira/anti-bruteforce/internal/logger"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const (
	dsn string = "host=localhost port=5432 user=postgres password=dbpass sslmode=disable"
)

type DataBaseSuite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
	e      *database.Entry
	logg   *zap.Logger
}

func (suite *DataBaseSuite) SetupTest() {
	suite.ctx, suite.cancel = context.WithCancel(context.Background())
	suite.e = &database.Entry{
		IP:   "172.0.0.1",
		Mask: "255.255.255.0",
		List: "black",
	}
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	suite.logg = l.InitLogger()
}

func (suite *DataBaseSuite) TestAddToListSuite() {
	defer suite.cancel()
	sql := database.New(suite.logg)
	err := sql.Connect(suite.ctx, dsn)
	if err != nil {
		suite.logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	err = sql.AddToList(suite.ctx, suite.e)
	suite.NoErrorf(err, "expected no error, but got %w", err)
}

func (suite *DataBaseSuite) TestCheckInListSuite() {
	defer suite.cancel()
	sql := database.New(suite.logg)
	err := sql.Connect(suite.ctx, dsn)
	if err != nil {
		suite.logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}
	ok, err := sql.CheckInList(suite.ctx, "172.0.0.1", "black")
	suite.True(ok, "expected true, but got %v", ok)
	suite.NoError(err, "expected no error, but got %w", err)
}

func (suite *DataBaseSuite) TestRemoveFromListSuite() {
	defer suite.cancel()
	sql := database.New(suite.logg)
	err := sql.Connect(suite.ctx, dsn)
	if err != nil {
		suite.logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}
	err = sql.RemoveFromList(suite.ctx, suite.e)
	suite.NoError(err, "expected no error, but got %w", err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DataBaseSuite))
}
