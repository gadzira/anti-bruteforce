package internalhttp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gadzira/anti-bruteforce/internal/app"
	"github.com/gadzira/anti-bruteforce/internal/database"
	"github.com/gadzira/anti-bruteforce/internal/domain"
	"github.com/gadzira/anti-bruteforce/internal/logger"
	internalhttp "github.com/gadzira/anti-bruteforce/internal/server/http"
	"github.com/gadzira/anti-bruteforce/internal/storage"
	"go.uber.org/zap"
)

const (
	RemAddr = "127.0.0.1:8080"
)

func TestHelloWorldHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(internalhttp.HelloWorldHandler().ServeHTTP)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != "Hello, World" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "Hello, World")
	}
}

func TestLoginHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:8800/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("user", "p@$$w0rd")
	req.RemoteAddr = RemAddr
	rr := httptest.NewRecorder()
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	logg := l.InitLogger()
	bs := storage.New(10, 100, 1000, "5m", logg)
	sql := database.New(logg)
	err = sql.Connect(ctx, "host=localhost port=5432 user=postgres password=dbpass sslmode=disable")
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	na := &domain.App{
		Ctx:     ctx,
		Logger:  logg,
		Storage: &bs,
		DB:      &sql,
	}
	a := app.New(na)

	handler := http.HandlerFunc(internalhttp.LoginHandler(a).ServeHTTP)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Body.String() != `ok=true` {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), `ok=true`)
	}
}

func TestLoginHandlerNegative(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:8800/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("user", "p@$$w0rd")
	req.RemoteAddr = RemAddr
	rr := httptest.NewRecorder()
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	logg := l.InitLogger()
	bs := storage.New(1, 100, 1000, "5m", logg)
	sql := database.New(logg)
	err = sql.Connect(ctx, "host=localhost port=5432 user=postgres password=dbpass sslmode=disable")
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	na := &domain.App{
		Ctx:     ctx,
		Logger:  logg,
		Storage: &bs,
		DB:      &sql,
	}

	a := app.New(na)
	handler := http.HandlerFunc(internalhttp.LoginHandler(a).ServeHTTP)
	for i := 0; i < bs.N+1; i++ {
		handler.ServeHTTP(rr, req)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v\n want %v\n",
			status, http.StatusOK)
	}
	if rr.Body.String() != `ok=trueok=false` {
		t.Errorf("handler returned unexpected body: got %v\n want %v\n",
			rr.Body.String(), `ok=trueok=false`)
	}
}

func TestResetBucketHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l := logger.New("test.log", "INFO", 1024, 1, 1, false, false)
	logg := l.InitLogger()
	bs := storage.New(10, 100, 1000, "5m", logg)
	sql := database.New(logg)
	err := sql.Connect(ctx, "host=localhost port=5432 user=postgres password=dbpass sslmode=disable")
	if err != nil {
		logg.Fatal("can't connect to DB: %s\n", zap.String("err", err.Error()))
	}

	na := &domain.App{
		Ctx:     ctx,
		Logger:  logg,
		Storage: &bs,
		DB:      &sql,
	}

	a := app.New(na)
	req1, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:8800/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	req1.SetBasicAuth("user", "p@$$w0rd")
	req1.RemoteAddr = RemAddr
	rr1 := httptest.NewRecorder()

	loginHandler := http.HandlerFunc(internalhttp.LoginHandler(a).ServeHTTP)

	for i := 0; i < 3; i++ {
		loginHandler.ServeHTTP(rr1, req1)
		if status := rr1.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	}

	req2, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:8800/rst_bckt", nil)
	if err != nil {
		t.Fatal(err)
	}
	req2.SetBasicAuth("user", "p@$$w0rd")
	req2.RemoteAddr = RemAddr
	rr2 := httptest.NewRecorder()
	q := req2.URL.Query()
	q.Add("key", "user")
	req2.URL.RawQuery = q.Encode()

	resetHandler := http.HandlerFunc(internalhttp.ResetBucketHandler(a).ServeHTTP)
	resetHandler.ServeHTTP(rr2, req2)
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr2.Body.String() != `done` {
		t.Errorf("handler returned unexpected body: got %v want %v", rr2.Body.String(), `done`)
	}

	req3, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8800/list_of_bucket", nil)
	if err != nil {
		t.Fatal(err)
	}
	req3.SetBasicAuth("user", "p@$$w0rd")
	req3.RemoteAddr = RemAddr
	rr3 := httptest.NewRecorder()
	showBucketsHandler := http.HandlerFunc(internalhttp.ListOfBucketHandler(a).ServeHTTP)
	showBucketsHandler.ServeHTTP(rr3, req3)
	if status := rr3.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr3.Body.String() != `{"127.0.0.1:8080":997,"p@$$w0rd":97}` {
		t.Errorf("handler returned unexpected body: got %v want %v", rr3.Body.String(), `{"127.0.0.1:8080":997,"p@$$w0rd":97}`)
	}
}
