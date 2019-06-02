package resources

import (
	"fmt"
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
	result.AuthenticatedMuxConfiguration,
		result.UnauthenticatedMuxConfiguration = r.MuxConfiguration()
	return result
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

	// retrieve the account uuid.
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

func (ar *AccountResource) CreateAndAppend(
	w http.ResponseWriter, request *http.Request) {
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
	vars := mux.Vars(request)
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

	uuid := u.Must(u.FromString(vars["uuid"]))
	if account.UUID != uuid {
		rb.BadRequest().WithError(fmt.Errorf("mismatching UUIDs")).Respond()
		return
	}

	// retrieve the account.
	a, err := ar.accountService.Get(uuid)
	if err != nil {
		rb.NotFound().WithError(err).Respond()
		return
	}

	if len(a.Confirmations) < len(account.Confirmations) {
		rb.BadRequest().
			WithError(fmt.Errorf("cannot remove confirmations")).Respond()
	}

	// upsert the account.
	err = ar.accountService.Put(account)
	if err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	rb.NoContent().Respond()
}

func (ar *AccountResource) Delete(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	rb := response.Builder(w)

	// retrieve the account uuid.
	uuid := u.Must(u.FromString(vars["uuid"]))
	_, err := ar.accountService.Get(uuid)
	if err != nil {
		rb.NotFound().WithError(err).Respond()
		return
	}

	// delete the account.
	err = ar.accountService.Remove(uuid)
	if err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	rb.NoContent().Respond()
}
