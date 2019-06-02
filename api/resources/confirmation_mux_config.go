package resources

import (
	"tri-fitness/genesis/api/middleware"
	"tri-fitness/genesis/api/server"

	"github.com/gorilla/mux"
)

func (cr *ConfirmationResource) MuxConfiguration() (authenticated server.MuxConfiguration) {
	authMiddleware :=
		middleware.NewAuthenticationMiddleware(
			middleware.AuthenticationMiddlewareParameters{
				Logger:        cr.logger,
				Authenticator: cr.authenticator,
			})
	loggingMiddleware :=
		middleware.NewLoggingMiddleware(cr.logger)
	//errorMiddleware :=
	//	middleware.NewErrorMiddleware(cr.logger)

	authenticated = server.MuxConfiguration{
		PathPrefix: "/accounts",
		Middleware: []mux.MiddlewareFunc{
			authMiddleware.Authenticate,
			//errorMiddleware.HandleError,
			loggingMiddleware.Log,
		},
		Handlers: []server.HandlerConfiguration{
			{
				Path:        "/{uuid}/confirmations/{id}/code",
				HandlerFunc: cr.Replace,
				Methods:     []string{"PUT"},
			},
			{
				Path:        "/{uuid}/confirmations/{id}/code/",
				HandlerFunc: cr.Replace,
				Methods:     []string{"PUT"},
			},
		},
	}

	return
}
