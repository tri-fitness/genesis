package middleware

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type ErrorMiddleware struct {
	logger *zap.Logger
}

func NewErrorMiddleware(l *zap.Logger) ErrorMiddleware {
	return ErrorMiddleware{logger: l}
}

func (m *ErrorMiddleware) HandleError(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				m.logger.Error(fmt.Sprintf("Recovered from panic (%v)\n", recovered))
				w.WriteHeader(500)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
