package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) LoggingMiddleware {
	return LoggingMiddleware{
		logger: logger,
	}
}

func (m *LoggingMiddleware) Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		m.logger.Info(fmt.Sprintf("%s - %s (%v)\n", r.Method, r.URL.Path, time.Since(start)))
	})
}
