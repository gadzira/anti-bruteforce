package internalhttp

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Local()
		s.log.Info("Handling request from",
			zap.String("remote IP", r.RemoteAddr),
			zap.Time("DateAndTime", time.Now().Local()),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("proto", r.Proto),
			zap.Int("status code", http.StatusOK),
			zap.Duration("latency", time.Since(start)))
		next.ServeHTTP(w, r)
	})
}
