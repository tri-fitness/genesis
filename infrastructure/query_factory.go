package infrastructure

import (
	"database/sql"

	u "github.com/gofrs/uuid"
	"go.uber.org/fx"
)

type QueryFactory interface {
	FindAccountByUUID(u.UUID) AccountQuery
	FindAccountByPrimaryCredential(string) AccountQuery
}

type queryFactory struct {
	db *sql.DB
}

type QueryFactoryParameters struct {
	fx.In

	DB *sql.DB `name:"rwDB"`
}

func NewQueryFactory(parameters QueryFactoryParameters) QueryFactory {
	return &queryFactory{
		db: parameters.DB,
	}
}

func (f *queryFactory) FindAccountByUUID(uuid u.UUID) AccountQuery {
	return NewFindAccountByUUIDQuery(FindAccountByUUIDParameters{
		DB:          f.db,
		AccountUUID: uuid,
	})
}

func (f *queryFactory) FindAccountByPrimaryCredential(
	credential string) AccountQuery {

	return NewFindAccountByPrimaryCredentialQuery(
		FindAccountByPrimaryCredentialParameters{
			DB:         f.db,
			Credential: credential,
		})
}
