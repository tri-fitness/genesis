package server

import (
	"context"
	"fmt"
	"net/http"
	"tri-fitness/genesis/config"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MuxConfiguration struct {
	PathPrefix string
	Handlers   []HandlerConfiguration
	Middleware []mux.MiddlewareFunc
}

type HandlerConfiguration struct {
	Path        string
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Methods     []string
}

type ServerParameters struct {
	fx.In

	Configuration    config.Configuration
	MuxConfiguration []MuxConfiguration `group:"muxConfiguration"`
	Logger           *zap.Logger
}

type Server struct {
	host       string
	port       int
	httpServer http.Server
	logger     *zap.Logger
}

func New(parameters ServerParameters) Server {

	serverConfig := parameters.Configuration.Server

	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port),
		Handler: newMux(parameters.Logger, parameters.MuxConfiguration),
	}

	s := Server{
		port:       serverConfig.Port,
		host:       serverConfig.Host,
		httpServer: httpServer,
		logger:     parameters.Logger,
	}

	return s
}

func newMux(
	logger *zap.Logger, configuration []MuxConfiguration) *mux.Router {

	r := mux.NewRouter()
	for _, m := range configuration {
		sr := r.PathPrefix(m.PathPrefix).Subrouter()
		for _, h := range m.Handlers {
			sr.HandleFunc(h.Path, h.HandlerFunc).Methods(h.Methods...)
			logger.Info(
				fmt.Sprintf("%s\t%s%s", h.Methods, m.PathPrefix, h.Path))
		}
		for _, middleware := range m.Middleware {
			sr.Use(middleware)
		}
	}

	return r
}

func (s *Server) Start() error {
	s.logger.Info(
		"Starting HTTP server",
		zap.String("host", s.host),
		zap.Int("port", s.port),
	)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Host() string {
	return s.host
}
