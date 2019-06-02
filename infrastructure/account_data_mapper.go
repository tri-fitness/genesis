package infrastructure

import (
	"database/sql"
	"errors"
	"strings"
	"tri-fitness/genesis/domain"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AccountDataMapperParameters struct {
	fx.In

	DB     *sql.DB `name:"rwDB"`
	Logger *zap.Logger
}

type AccountDataMapper struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewAccountDataMapper(
	parameters AccountDataMapperParameters) AccountDataMapper {

	return AccountDataMapper{db: parameters.DB, logger: parameters.Logger}
}

func (dm *AccountDataMapper) toAccount(accounts ...interface{}) ([]domain.Account, error) {
	accs := []domain.Account{}
	for _, account := range accounts {
		var acc domain.Account
		var ok bool
		if acc, ok = account.(domain.Account); !ok {
			return []domain.Account{}, errors.New("invalid type")
		}
		accs = append(accs, acc)
	}
	return accs, nil
}

func (dm *AccountDataMapper) Insert(tx *sql.Tx, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.insert(tx, accs...)
}

func (dm *AccountDataMapper) insertSQL(accounts ...domain.Account) (sql string, args []interface{}) {
	sql =
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
			TARGET_UUID,
			CREATED_AT,
			UPDATED_AT,
			DELETED_AT
		) VALUES `
	var vals []string
	for _, account := range accounts {
		vals = append(vals, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			account.UUID.String(),
			account.Type.String(),
			account.GivenName,
			account.Surname,
			account.Bio,
			account.Email,
			account.Phone,
			account.PrimaryCredential,
			account.SecondaryCredential,
			account.NotificationTargetUUID,
			account.CreatedAt,
			account.UpdatedAt,
			account.DeletedAt,
		)
	}
	sql = sql + strings.Join(vals, ", ") + ";"
	return
}

func (dm *AccountDataMapper) insertConfirmationsSQL(accounts ...domain.Account) (sql string, args []interface{}) {
	sql =
		`INSERT INTO CONFIRMATION
		(
			ID,
			ACCOUNT_UUID,
			TYPE,
			CODE,
			NOTIFICATION_UUID,
			EXPIRED_AT,
			CONFIRMED_AT,
			CREATED_AT,
			UPDATED_AT
		)
		VALUES `
	var vals []string
	for _, account := range accounts {
		for _, confirmation := range account.Confirmations {
			vals = append(vals, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
			args = append(args,
				confirmation.ID,
				account.UUID.String(),
				confirmation.Type.String(),
				confirmation.Code,
				confirmation.NotificationUUID,
				confirmation.ExpiredAt,
				confirmation.ConfirmedAt,
				confirmation.CreatedAt,
				confirmation.UpdatedAt,
			)
		}
	}
	sql = sql + strings.Join(vals, ", ") + ";"
	return
}

func (dm *AccountDataMapper) insert(tx *sql.Tx, accounts ...domain.Account) error {
	sql, args := dm.insertSQL(accounts...)
	err := dm.prepareAndExec(tx, sql, args)
	if err != nil {
		return err
	}
	sql, args = dm.insertConfirmationsSQL(accounts...)
	return dm.prepareAndExec(tx, sql, args)
}

func (dm *AccountDataMapper) Update(tx *sql.Tx, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.update(tx, accs...)
}

func (dm *AccountDataMapper) updateSQL(
	accounts ...domain.Account) (sql []string, args [][]interface{}) {
	for _, account := range accounts {
		s :=
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
				TARGET_UUID = ?,
				CREATED_AT = ?,
				UPDATED_AT = ?,
				DELETED_AT = ?
			WHERE
				UUID = ?`
		sql = append(sql, s)
		sArgs := []interface{}{
			account.Type.String(),
			account.GivenName,
			account.Surname,
			account.Bio,
			account.Email,
			account.Phone,
			account.PrimaryCredential,
			account.SecondaryCredential,
			account.NotificationTargetUUID,
			account.CreatedAt,
			account.UpdatedAt,
			account.DeletedAt,
			account.UUID.String(),
		}
		args = append(args, sArgs)

		// delete all confirmations.
		cSQL, cArgs := dm.deleteConfirmationSQL(accounts...)
		sql = append(sql, cSQL)
		args = append(args, cArgs)

		// re-insert all confirmations.
		cSQL, cArgs = dm.insertConfirmationsSQL(accounts...)
		sql = append(sql, cSQL)
		args = append(args, cArgs)
	}
	return
}

func (dm *AccountDataMapper) update(tx *sql.Tx, accounts ...domain.Account) error {
	sql, args := dm.updateSQL(accounts...)
	for idx, s := range sql {
		if err := dm.prepareAndExec(tx, s, args[idx]); err != nil {
			return err
		}
	}
	return nil
}

func (dm *AccountDataMapper) Delete(tx *sql.Tx, accounts ...interface{}) error {
	if len(accounts) == 0 {
		return nil
	}
	accs, err := dm.toAccount(accounts...)
	if err != nil {
		return err
	}
	return dm.delete(tx, accs...)
}

func (dm *AccountDataMapper) deleteSQL(accounts ...domain.Account) (sql string, args []interface{}) {
	sql = "DELETE FROM ACCOUNT WHERE UUID IN (?)"
	for _, account := range accounts {
		args = append(args, account.UUID.String())
	}
	return
}

func (dm *AccountDataMapper) deleteConfirmationSQL(accounts ...domain.Account) (sql string, args []interface{}) {
	sql = "DELETE FROM CONFIRMATION WHERE ID IN (?)"
	for _, account := range accounts {
		for _, confirmation := range account.Confirmations {
			args = append(args, confirmation.ID)
		}
	}
	return
}

func (dm *AccountDataMapper) delete(tx *sql.Tx, accounts ...domain.Account) error {
	sql, args := dm.deleteConfirmationSQL(accounts...)
	err := dm.prepareAndExec(tx, sql, args)
	if err != nil {
		return err
	}
	sql, args = dm.deleteSQL(accounts...)
	return dm.prepareAndExec(tx, sql, args)
}

func (dm *AccountDataMapper) prepareAndExec(
	tx *sql.Tx, sql string, args []interface{}) error {
	s, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer s.Close()
	_, err = s.Exec(args...)
	return err
}
