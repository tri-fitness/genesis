package resources

import (
	"tri-fitness/genesis/api/middleware"
	"tri-fitness/genesis/api/server"

	"github.com/gorilla/mux"
)

func (ar *AccountResource) MuxConfiguration() (authenticated, unauthenticated server.MuxConfiguration) {
	authMiddleware :=
		middleware.NewAuthenticationMiddleware(
			middleware.AuthenticationMiddlewareParameters{
				Logger:        ar.logger,
				Authenticator: ar.authenticator,
			})
	loggingMiddleware :=
		middleware.NewLoggingMiddleware(ar.logger)
	errorMiddleware :=
		middleware.NewErrorMiddleware(ar.logger)

	authenticated = server.MuxConfiguration{
		PathPrefix: "/accounts",
		Middleware: []mux.MiddlewareFunc{
			authMiddleware.Authenticate,
			//errorMiddleware.HandleError,
			loggingMiddleware.Log,
		},
		Handlers: []server.HandlerConfiguration{
			{
				Path:        "/{uuid}/",
				HandlerFunc: ar.Get,
				Methods:     []string{"GET"},
			},
			{
				Path:        "/{uuid}",
				HandlerFunc: ar.Get,
				Methods:     []string{"GET"},
			},
			{
				Path:        "/{uuid}",
				HandlerFunc: ar.Replace,
				Methods:     []string{"PUT"},
			},
			{
				Path:        "/{uuid}/",
				HandlerFunc: ar.Replace,
				Methods:     []string{"PUT"},
			},
			{
				Path:        "/{uuid}",
				HandlerFunc: ar.Delete,
				Methods:     []string{"DELETE"},
			},
			{
				Path:        "/{uuid}/",
				HandlerFunc: ar.Delete,
				Methods:     []string{"DELETE"},
			},
		},
	}

	unauthenticated = server.MuxConfiguration{
		PathPrefix: "/accounts",
		Middleware: []mux.MiddlewareFunc{
			loggingMiddleware.Log,
			errorMiddleware.HandleError,
		},
		Handlers: []server.HandlerConfiguration{
			{
				Path:        "",
				HandlerFunc: ar.CreateAndAppend,
				Methods:     []string{"POST"},
			},
			{
				Path:        "/",
				HandlerFunc: ar.CreateAndAppend,
				Methods:     []string{"POST"},
			},
		},
	}

	return
}
