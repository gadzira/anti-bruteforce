package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gadzira/anti-bruteforce/internal/database"
	"github.com/gadzira/anti-bruteforce/internal/domain"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Application interface{}

type Server struct {
	Router *mux.Router
	log    *zap.Logger
	app    *domain.App
}

func NewServer(l *zap.Logger, a *domain.App) *Server {
	return &Server{
		app: a,
		log: l,
	}
}

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
	router.Handle("/list_of_bucket", ListOfBucketHandler(s.app)).Methods("GET")
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

func ListOfBucketHandler(a *domain.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//nolint:errcheck
		bl := a.Storage.ShowBuckets()
		tm := make(map[string]int)
		for k, v := range bl {
			tm[k] = v.Limit
		}
		j, err := json.Marshal(tm)
		if err != nil {
			a.Logger.Error("can't marshal to jsom\n" + err.Error())
			return
		}
		_, err = w.Write(j)
		if err != nil {
			a.Logger.Error("can't write\n" + err.Error())
			return
		}
	})
}

func LoginHandler(a *domain.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log, pass, _ := r.BasicAuth()
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			a.Logger.Error("can't get ip from r.RemoteAddr\n" + err.Error())
			return
		}

		bl, err := a.DB.CheckInList(a.Ctx, ip, "black")
		if err != nil {
			a.Logger.Error("can't check black list in db: " + err.Error())
			return
		}

		wl, err := a.DB.CheckInList(a.Ctx, ip, "white")
		if err != nil {
			a.Logger.Error("can't check white list in db: " + err.Error())
			return
		}

		switch {
		case bl:
			_, err := w.Write([]byte(fmt.Sprintf("ok=%t", false)))
			if err != nil {
				a.Logger.Error("can't write\n" + err.Error())
				return
			}
		case wl:
			_, err := w.Write([]byte(fmt.Sprintf("ok=%t", true)))
			if err != nil {
				a.Logger.Error("can't write\n" + err.Error())
				return
			}
		default:
			cr := a.Storage.CheckRequest(log, pass, r.RemoteAddr)
			resultOfCheck := fmt.Sprintf("ok=%t", cr)
			_, err := w.Write([]byte(resultOfCheck))
			if err != nil {
				a.Logger.Error("can't write\n" + err.Error())
				return
			}
		}
	})
}

func ResetBucketHandler(a *domain.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		key := params.Get("key")
		a.Storage.ResetBucket(key)
		_, err := w.Write([]byte("done"))
		if err != nil {
			a.Logger.Error("can't write\n" + err.Error())
			return
		}
	})
}

func AddToListHandler(a *domain.App, list string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		e := database.Entry{
			IP:   params.Get("ip"),
			Mask: params.Get("mask"),
			List: list,
		}
		err := a.DB.AddToList(a.Ctx, &e)
		if err != nil {
			a.Logger.Error("can't add to list\n" + err.Error())
			return
		}
	})
}

func RemoveFromListHandler(a *domain.App, list string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		e := database.Entry{
			IP:   params.Get("ip"),
			Mask: params.Get("mask"),
			List: list,
		}
		err := a.DB.RemoveFromList(a.Ctx, &e)
		if err != nil {
			a.Logger.Error("can't remove from list\n" + err.Error())
			return
		}
	})
}
