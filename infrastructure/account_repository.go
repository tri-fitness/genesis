package infrastructure

import (
	"errors"
	"tri-fitness/genesis/domain"

	"github.com/freerware/work"
	u "github.com/gofrs/uuid"
)

// AccountRepository represents a collection of all
// accounts within the application.
type AccountRepository interface {
	Get(u.UUID) (domain.Account, error)
	Put(domain.Account) error
	Remove(u.UUID) error
	Add(domain.Account) error
	Find(AccountQuery) ([]domain.Account, error)
	Size() (int, error)
}

type accountRepository struct {
	unit         work.Unit
	queryFactory QueryFactory
}

func NewAccountRepository(unit work.Unit, queryFactory QueryFactory) AccountRepository {
	return &accountRepository{unit: unit, queryFactory: queryFactory}
}

func (r *accountRepository) Find(query AccountQuery) ([]domain.Account, error) {
	return query.Execute()
}

func (r *accountRepository) Get(uuid u.UUID) (domain.Account, error) {
	query := r.queryFactory.FindAccountByUUID(uuid)
	matches, err := r.Find(query)
	if err != nil {
		return domain.Account{}, err
	}
	if len(matches) == 0 {
		return domain.Account{}, nil
	}
	return matches[0], nil
}

func (r *accountRepository) Put(account domain.Account) error {

	// check if the account exists.
	c, e := r.Get(account.UUID)
	if e != nil {
		return e
	}

	// if the account is not within the repository, add it.
	var empty u.UUID
	if c.UUID == empty {
		return r.Add(c)
	}

	// otherwise, replace the existing state.
	r.unit.Alter(account)
	return nil
}

func (r *accountRepository) Remove(uuid u.UUID) error {

	// check if the account exists.
	c, e := r.Get(uuid)
	if e != nil {
		return e
	}

	// if the account is not within the repository, throw an error.
	var empty u.UUID
	if c.UUID == empty {
		return errors.New("could not find the account")
	}

	// otherwise, remove the account.
	r.unit.Remove(c)
	return nil
}

func (r *accountRepository) Add(account domain.Account) error {

	// check if the account exists.
	c, e := r.Get(account.UUID)
	if e != nil {
		return e
	}

	// if the account is within the repository, throw an error.
	var empty u.UUID
	if c.UUID != empty {
		return errors.New("account already exists")
	}

	// otherwise, remove the account.
	r.unit.Add(account)
	return nil
}

func (r *accountRepository) Size() (int, error) {
	return 0, nil
}
