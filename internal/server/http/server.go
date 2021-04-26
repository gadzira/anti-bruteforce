package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gadzira/anti-bruteforce/internal/app"
	"github.com/gadzira/anti-bruteforce/internal/db"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Application interface{}

type Server struct {
	Router *mux.Router
	log    *zap.Logger
	app    *app.App
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
	router := mux.NewRouter()
	router.Handle("/hello", HelloWorldHandler()).Methods("GET")
	router.Handle("/login", LoginHandler(s.app)).Methods("POST")
	router.Handle("/rst_bckt", ResetBucketHandler(s.app)).Methods("POST")
	router.Handle("/add_bl", AddToListHandler(s.app, "black")).Methods("POST")
	router.Handle("/add_wl", AddToListHandler(s.app, "white")).Methods("POST")
	router.Handle("/del_bl", RemoveFromListHandler(s.app, "black")).Methods("DELETE")
	router.Handle("/del_wl", RemoveFromListHandler(s.app, "white")).Methods("DELETE")
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
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			a.Logger.Error("can't get ip from r.RemoteAddr" + err.Error())
			return
		}

		bl, err := a.DB.CheckInList(a.Ctx, ip, "black")
		if err != nil {
			a.Logger.Error("can't check black list in db" + err.Error())
			return
		}

		wl, err := a.DB.CheckInList(a.Ctx, ip, "white")
		if err != nil {
			a.Logger.Error("can't check white list in db" + err.Error())
			return
		}

		if bl {
			w.Write([]byte(fmt.Sprintf("ok=%t", false)))
		} else if wl {
			w.Write([]byte(fmt.Sprintf("ok=%t", true)))
		} else {
			cr, _ := a.Storage.CheckRequest(log, pass, r.RemoteAddr)
			resultOfCheck := fmt.Sprintf("ok=%t", cr)
			w.Write([]byte(resultOfCheck))
		}
	})
}

func ResetBucketHandler(a *app.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		key := params.Get("key")
		a.Storage.ResetBucket(key)
	})
}

func AddToListHandler(a *app.App, list string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		e := db.Entry{
			IP:   params.Get("ip"),
			Mask: params.Get("mask"),
			List: list,
		}
		a.DB.AddToList(a.Ctx, &e)
	})
}

func RemoveFromListHandler(a *app.App, list string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		e := db.Entry{
			IP:   params.Get("ip"),
			Mask: params.Get("mask"),
			List: list,
		}
		a.DB.RemoveFromList(a.Ctx, &e)
	})
}
