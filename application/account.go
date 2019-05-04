package application

import (
	"fmt"
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
}

type AccountServiceParameters struct {
	fx.In

	Uniter            work.Uniter `name:"sqlWorkUniter"`
	RepositoryFactory infra.RepositoryFactory
	QueryFactory      infra.QueryFactory
}

func NewAccountService(
	parameters AccountServiceParameters) AccountService {
	return AccountService{
		uniter:            parameters.Uniter,
		repositoryFactory: parameters.RepositoryFactory,
		queryFactory:      parameters.QueryFactory,
	}
}

func (a *AccountService) Create(account domain.Account) error {
	unit, err := a.uniter.Unit()
	if err != nil {
		return err
	}
	repository := a.repositoryFactory.Account(unit)
	hashedCredential, err :=
		bcrypt.GenerateFromPassword(
			[]byte(account.SecondaryCredential), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.SecondaryCredential = string(hashedCredential)
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
