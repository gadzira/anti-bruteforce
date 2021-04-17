package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gadzira/anti-bruteforce/internal/app"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	sep string = "#"
)

type Application interface{}

type Server struct {
	Router *mux.Router
	log    *zap.Logger
	app    *app.App
	// Application
}

func NewServer(l *zap.Logger, a *app.App) *Server {
	return &Server{
		app: a,
		log: l,
	}
}

// тут нужно создать экземпляр конфига

func (s *Server) Start(ctx context.Context, addr string) error {
	s.initializeRoutes()
	s.log.Info("Server is running on", zap.String("port", addr))
	err := http.ListenAndServe(addr, s.Router)
	if err != nil {
		s.log.Fatal("can't start server:" + err.Error())
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO: stop server
	return nil
}

func (s *Server) initializeRoutes() {
	// нужно прокинуть Storage в LoginHandler, а там вызвать метод AddNewBucket
	router := mux.NewRouter()
	router.Handle("/hello", HelloWorldHandler()).Methods("GET")
	router.Handle("/login", LoginHandler(s.app)).Methods("POST")
	router.Use(s.loggingMiddleware)
	s.Router = router
}

func HelloWorldHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//nolint:errcheck
		w.Write([]byte("Hello, World"))
	})
}

func LoginHandler(a *app.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log, pass, _ := r.BasicAuth()
		arg := []string{r.RemoteAddr, sep, log, sep, pass}
		builder := strings.Builder{}
		for _, i := range arg {
			builder.WriteString(i)
		}
		longString := builder.String()
		a.Storage.AddNewBucket(longString)
		// app.AddBucket(longString)

		// var b models.Bucket
		// b.Login = log
		// b.Password = pass
		// b.SourceIP = r.RemoteAddr

		// fmt.Println(b)

		//nolint:errcheck
		// w.Write([]byte("Created!"))

		payload, _ := json.Marshal("OK")
		w.Write([]byte(payload))

		// w.WriteHeader(http.StatusOK)
		// w.WriteHeader(http.StatusInternalServerError)
	})
}
