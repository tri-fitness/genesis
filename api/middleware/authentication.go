package middleware

import (
	"net/http"
	"tri-fitness/genesis/api/response"

	"go.uber.org/zap"
)

type AuthenticationMiddleware struct {
	logger        *zap.Logger
	authenticator Authenticator
}

type AuthenticationMiddlewareParameters struct {
	Logger        *zap.Logger
	Authenticator Authenticator
}

func NewAuthenticationMiddleware(
	params AuthenticationMiddlewareParameters,
) AuthenticationMiddleware {
	return AuthenticationMiddleware{
		logger:        params.Logger,
		authenticator: params.Authenticator,
	}
}

func (m *AuthenticationMiddleware) Authenticate(
	h http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			// parse the header.
			rb := response.Builder(w)
			basic, realm, charset := "Basic", "api access", "utf8mb4"
			primary, secondary, ok := r.BasicAuth()
			if !ok {
				rb.
					Unauthorized().
					Challenge(basic, realm, charset).
					Respond()
				return
			}

			// authenticate.
			authenticated, err :=
				m.authenticator.Authenticate(Credentials{
					Primary:   primary,
					Secondary: secondary,
				})
			if err != nil {
				rb.InternalServerError().WithError(err).Respond()
			}
			if !authenticated {
				rb.
					Unauthorized().
					Challenge(basic, realm, charset).
					Respond()
				return
			}

			h.ServeHTTP(w, r)
		})
}
