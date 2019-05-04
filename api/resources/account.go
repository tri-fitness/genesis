package resources

import (
	"net/http"
	"tri-fitness/genesis/api/middleware"
	r "tri-fitness/genesis/api/representations"
	"tri-fitness/genesis/api/response"
	"tri-fitness/genesis/api/server"
	app "tri-fitness/genesis/application"

	u "github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AccountResourceResult struct {
	fx.Out

	AccountResource                 AccountResource
	AuthenticatedMuxConfiguration   server.MuxConfiguration `group:"muxConfiguration"`
	UnauthenticatedMuxConfiguration server.MuxConfiguration `group:"muxConfiguration"`
}

type AccountResourceParameters struct {
	fx.In

	AccountService          app.AccountService
	RepresentationFactories map[string]r.RepresentationFactory
	Logger                  *zap.Logger
	Authenticator           middleware.Authenticator
}

type AccountResource struct {
	accountService          app.AccountService
	representationFactories map[string]r.RepresentationFactory
	logger                  *zap.Logger
	authenticator           middleware.Authenticator
}

func NewAccountResource(
	parameters AccountResourceParameters,
) AccountResourceResult {
	r := AccountResource{
		accountService:          parameters.AccountService,
		representationFactories: parameters.RepresentationFactories,
		logger:                  parameters.Logger,
		authenticator:           parameters.Authenticator,
	}
	result := AccountResourceResult{
		AccountResource: r,
	}
	result.AuthenticatedMuxConfiguration, result.UnauthenticatedMuxConfiguration =
		r.MuxConfiguration()
	return result
}

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
			errorMiddleware.HandleError,
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

func (ar *AccountResource) Get(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	rb := response.Builder(w)

	// TODO(FREER) - don't use exact header value match here.
	rf, ok := ar.representationFactories[request.Header.Get("Accept")]
	if !ok {
		rb.NotAcceptable().Respond()
		return
	}

	// retrieve the account.
	uuid := u.Must(u.FromString(vars["uuid"]))
	account, err := ar.accountService.Get(uuid)
	if err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	// construct the representation.
	representation, err := rf.Account(account)
	if err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	// respond.
	rb.OK().Body(representation.AsBytes()).Respond()
}

func (ar *AccountResource) CreateAndAppend(w http.ResponseWriter, request *http.Request) {
	rb := response.Builder(w)

	// TODO(FREER) - don't use exact header value match here.
	rf, ok := ar.representationFactories[request.Header.Get("Accept")]
	if !ok {
		rb.NotAcceptable().Respond()
		return
	}

	account, err := rf.AccountEntityFromRequest(request)
	if err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	// create the account.
	err = ar.accountService.Create(account)
	if err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	uri, _ := request.URL.Parse("/" + account.UUID.String())
	rb.Created(*uri).Respond()
}

func (ar *AccountResource) Replace(w http.ResponseWriter, request *http.Request) {
	rb := response.Builder(w)
	rb.OK().Body([]byte("worked!"))
	rb.Respond()
}

func (ar *AccountResource) Delete(w http.ResponseWriter, request *http.Request) {
	rb := response.Builder(w)
	rb.OK().Body([]byte("worked!"))
	rb.Respond()
}
