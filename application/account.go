package application

import (
	"fmt"
	"time"
	"tri-fitness/genesis/domain"
	infra "tri-fitness/genesis/infrastructure"

	work "github.com/freerware/work"
	u "github.com/gofrs/uuid"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	uniter            work.Uniter
	repositoryFactory infra.RepositoryFactory
	queryFactory      infra.QueryFactory
	sms               infra.ShortMessageService
}

type AccountServiceParameters struct {
	fx.In

	Uniter            work.Uniter `name:"sqlWorkUniter"`
	RepositoryFactory infra.RepositoryFactory
	QueryFactory      infra.QueryFactory
	SMS               infra.ShortMessageService
}

func NewAccountService(
	parameters AccountServiceParameters) AccountService {
	return AccountService{
		uniter:            parameters.Uniter,
		repositoryFactory: parameters.RepositoryFactory,
		queryFactory:      parameters.QueryFactory,
		sms:               parameters.SMS,
	}
}

func (a *AccountService) Create(account domain.Account) error {
	// setup account for notifications.
	targetUUID, err := a.sms.RegisterAccount(account)
	if err != nil {
		return err
	}

	// save account.
	account.NotificationTargetUUID = targetUUID
	now := time.Now()
	account.UpdatedAt, account.CreatedAt = now, now
	unit, err := a.uniter.Unit()
	if err != nil {
		return err
	}
	repository := a.repositoryFactory.Account(unit)

	// hash password.
	hashedCredential, err :=
		a.hashCredential(account.SecondaryCredential)
	account.SecondaryCredential = string(hashedCredential)

	// create confirmations.
	code := 1234
	expiresAt := time.Now().Add(5 * time.Minute)
	pConfirmation := account.NewConfirmation(
		domain.ConfirmationTypePhoneNumber,
		code,
		expiresAt,
	)
	notificationUUID, err := a.sms.SendPhoneConfirmation(account, pConfirmation)
	pConfirmation.NotificationUUID = &notificationUUID
	account.SetConfirmation(pConfirmation)

	if err := repository.Add(account); err != nil {
		return err
	}
	return unit.Save()
}

func (a *AccountService) Put(account domain.Account) error {
	unit, err := a.uniter.Unit()
	if err != nil {
		return err
	}
	repository := a.repositoryFactory.Account(unit)
	if err := repository.Put(account); err != nil {
		return err
	}
	return unit.Save()
}

func (a *AccountService) Get(uuid u.UUID) (domain.Account, error) {
	var empty domain.Account
	unit, err := a.uniter.Unit()
	if err != nil {
		return empty, err
	}
	repository := a.repositoryFactory.Account(unit)
	account, err := repository.Get(uuid)
	if err != nil {
		return empty, err
	}

	//check if not found.
	var emptyUUID u.UUID
	if account.UUID == emptyUUID {
		return empty, fmt.Errorf("no account with UUID %q", uuid.String())
	}
	if err = unit.Save(); err != nil {
		return empty, err
	}
	return account, nil
}

func (a *AccountService) Remove(uuid u.UUID) error {
	unit, err := a.uniter.Unit()
	if err != nil {
		return err
	}
	repository := a.repositoryFactory.Account(unit)
	if err := repository.Remove(uuid); err != nil {
		return err
	}
	return unit.Save()
}

func (a *AccountService) GetByPrimaryCredential(
	credential string) (domain.Account, error) {

	var accounts []domain.Account
	var empty domain.Account
	unit, err := a.uniter.Unit()
	if err != nil {
		return empty, err
	}

	repository := a.repositoryFactory.Account(unit)
	query :=
		a.queryFactory.FindAccountByPrimaryCredential(credential)
	if accounts, err = repository.Find(query); err != nil {
		return empty, err
	}

	//check if not found.
	if len(accounts) == 0 {
		return empty, fmt.Errorf("no account with credential %q", credential)
	}
	if err = unit.Save(); err != nil {
		return empty, err
	}
	return accounts[0], nil
}

func (a *AccountService) hashCredential(
	credential string) (string, error) {

	hashedCredential, err :=
		bcrypt.GenerateFromPassword(
			[]byte(credential), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedCredential), nil
}
