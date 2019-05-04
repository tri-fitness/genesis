package middleware

import (
	"tri-fitness/genesis/application"
	"tri-fitness/genesis/domain"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Primary   string
	Secondary string
}

type Authenticator interface {
	Authenticate(Credentials) (bool, error)
}

type authenticator struct {
	accountService application.AccountService
}

func NewAuthenticator(accountService application.AccountService) Authenticator {
	return &authenticator{accountService: accountService}
}

func (b *authenticator) Authenticate(
	creds Credentials) (authenticated bool, err error) {

	// retrieve the account with the primary credential provided.
	var account domain.Account
	account, err =
		b.accountService.GetByPrimaryCredential(creds.Primary)
	if err != nil {
		return
	}

	// validate passwords match.
	err = bcrypt.CompareHashAndPassword(
		[]byte(account.SecondaryCredential), []byte(creds.Secondary))
	authenticated = err == nil
	return
}
