package infrastructure

import (
	"database/sql"
	"errors"
	"strings"
	"tri-fitness/genesis/domain"

	"go.uber.org/fx"
)

type AccountDataMapperParameters struct {
	fx.In

	DB *sql.DB `name:"rwDB"`
}

type AccountDataMapper struct {
	db *sql.DB
}

func NewAccountDataMapper(
	parameters AccountDataMapperParameters) AccountDataMapper {

	return AccountDataMapper{db: parameters.DB}
}

func (dm *AccountDataMapper) Insert(accounts ...interface{}) error {
	as := []domain.Account{}
	for _, account := range accounts {
		var a domain.Account
		var ok bool
		if a, ok = account.(domain.Account); !ok {
			return errors.New("invalid type")
		}
		as = append(as, a)
	}
	return dm._Insert(as...)
}

func (dm *AccountDataMapper) _Insert(accounts ...domain.Account) error {

	sqlStatements := []string{}
	sqlArgs := []interface{}{}
	for _, a := range accounts {
		sql :=
			`INSERT INTO ACCOUNT
			(
				UUID,
				TYPE,
				GIVEN_NAME,
				SURNAME,
				BIO,
				EMAIL_ADDRESS,
				PHONE_NUMBER,
				PRIMARY_CREDENTIAL,
				SECONDARY_CREDENTIAL,
				CREATED_AT,
				UPDATED_AT,
				DELETED_AT
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		sqlStatements = append(sqlStatements, sql)
		sqlArgs =
			append(sqlArgs,
				a.UUID.String(), a.Type, a.GivenName, a.Surname, a.Bio,
				a.Email, a.Phone, a.PrimaryCredential,
				a.SecondaryCredential, a.CreatedAt, a.UpdatedAt, a.DeletedAt)
	}
	statement, err :=
		dm.db.Prepare(strings.Join(sqlStatements, ";"))
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(sqlArgs...)
	return err
}

func (dm *AccountDataMapper) Update(accounts ...interface{}) error {
	as := []domain.Account{}
	for _, account := range accounts {
		var a domain.Account
		var ok bool
		if a, ok = account.(domain.Account); !ok {
			return errors.New("invalid type")
		}
		as = append(as, a)
	}
	return dm._Update(as...)
}

func (dm *AccountDataMapper) _Update(accounts ...domain.Account) error {

	sqlStatements := []string{}
	sqlArgs := []interface{}{}
	for _, a := range accounts {
		sql :=
			`UPDATE
				ACCOUNT
			SET
				TYPE = ?,
				GIVEN_NAME = ?,
				SURNAME = ?,
				BIO = ?,
				EMAIL_ADDRESS = ?,
				PHONE_NUMBER = ?,
				PRIMARY_CREDENTIAL = ?,
				SECONDARY_CREDENTIAL = ?,
				CREATED_AT = ?,
				UPDATED_AT = ?,
				DELETED_AT = ?
			WHERE
				UUID = ?`
		sqlStatements = append(sqlStatements, sql)
		sqlArgs =
			append(sqlArgs,
				a.Type, a.GivenName, a.Surname, a.Bio, a.Email,
				a.Phone, a.PrimaryCredential, a.SecondaryCredential,
				a.CreatedAt, a.UpdatedAt, a.DeletedAt)
	}
	statement, err :=
		dm.db.Prepare(strings.Join(sqlStatements, ";"))
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(sqlArgs)
	return err
}

func (dm *AccountDataMapper) Delete(accounts ...interface{}) error {
	as := []domain.Account{}
	for _, account := range accounts {
		var a domain.Account
		var ok bool
		if a, ok = account.(domain.Account); !ok {
			return errors.New("invalid type")
		}
		as = append(as, a)
	}
	return dm._Delete(as...)
}

func (dm *AccountDataMapper) _Delete(accounts ...domain.Account) error {
	statement, err :=
		dm.db.Prepare("DELETE FROM ACCOUNT WHERE uuid IN (?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	uuids := []interface{}{}
	for _, account := range accounts {
		uuids = append(uuids, account.UUID.String())
	}
	_, err = statement.Exec(uuids...)
	return err
}
