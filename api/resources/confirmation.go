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

type ConfirmationResourceResult struct {
	fx.Out

	ConfirmationResource          ConfirmationResource
	AuthenticatedMuxConfiguration server.MuxConfiguration `group:"muxConfiguration"`
}

type ConfirmationResourceParameters struct {
	fx.In

	AccountService          app.AccountService
	RepresentationFactories map[string]r.RepresentationFactory
	Logger                  *zap.Logger
	Authenticator           middleware.Authenticator
}

type ConfirmationResource struct {
	accountService          app.AccountService
	representationFactories map[string]r.RepresentationFactory
	logger                  *zap.Logger
	authenticator           middleware.Authenticator
}

func NewConfirmationResource(
	parameters ConfirmationResourceParameters,
) ConfirmationResourceResult {
	r := ConfirmationResource{
		accountService:          parameters.AccountService,
		representationFactories: parameters.RepresentationFactories,
		logger:                  parameters.Logger,
		authenticator:           parameters.Authenticator,
	}
	result := ConfirmationResourceResult{
		ConfirmationResource: r,
	}
	result.AuthenticatedMuxConfiguration = r.MuxConfiguration()
	return result
}

// PUT /accounts/{uuid}/confirmations/{id}/code
func (cr *ConfirmationResource) Replace(
	w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	rb := response.Builder(w)

	// TODO(FREER) - don't use exact header value match here.
	rf, ok := cr.representationFactories[request.Header.Get("Accept")]
	if !ok {
		rb.NotAcceptable().Respond()
		return
	}

	// retrieve the account uuid and confirmation id.
	uuid := u.Must(u.FromString(vars["uuid"]))
	// id, err := strconv.Atoi(vars["id"])
	//if err != nil {
	//	rb.InternalServerError().WithError(err).Respond()
	//	return
	//}

	// retrieve the account.
	account, err := cr.accountService.Get(uuid)
	if err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	code, err := rf.CodeEntityFromRequest(request)
	if err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	// confirm.
	if err = account.Confirm(code.Code); err != nil {
		rb.BadRequest().WithError(err).Respond()
		return
	}

	// update the account.
	if err = cr.accountService.Put(account); err != nil {
		rb.InternalServerError().WithError(err).Respond()
		return
	}

	rb.NoContent().Respond()
}
