package infrastructure

import (
	"database/sql"
	"fmt"
	"strings"
	"tri-fitness/genesis/domain"

	"github.com/go-sql-driver/mysql"
	u "github.com/gofrs/uuid"
)

type findAccountByUUID struct {
	accountQuery

	uuid u.UUID
}

type FindAccountByUUIDParameters struct {
	DB          *sql.DB
	AccountUUID u.UUID
}

func NewFindAccountByUUIDQuery(
	parameters FindAccountByUUIDParameters) AccountQuery {

	return &findAccountByUUID{
		accountQuery: accountQuery{
			db: parameters.DB,
		},
		uuid: parameters.AccountUUID,
	}
}

func (q *findAccountByUUID) Execute() ([]domain.Account, error) {

	// retrieve accounts.
	accounts, err := q.accounts()
	if err != nil {
		return []domain.Account{}, err
	}

	// retrieve account confirmations.
	uuids := []u.UUID{}
	for _, a := range accounts {
		uuids = append(uuids, a.UUID)
	}
	confirmationsByAccount, err := q.confirmationsByAccount(uuids)
	if err != nil {
		return []domain.Account{}, err
	}

	// match confirmations to accounts.
	for idx, a := range accounts {
		if confirmations, ok := confirmationsByAccount[a.UUID]; ok {
			a.Confirmations = append(a.Confirmations, confirmations...)
			accounts[idx] = a // range uses a copy.
		}
	}

	return accounts, nil
}

func (q findAccountByUUID) accounts() ([]domain.Account, error) {

	matches := []domain.Account{}
	statement, err :=
		q.db.Prepare("SELECT * FROM ACCOUNT WHERE UUID = ?")
	if err != nil {
		return matches, err
	}
	defer statement.Close()

	rows, err := statement.Query(q.uuid.String())
	if err != nil {
		return matches, err
	}
	defer rows.Close()

	for rows.Next() {
		var deletedAt mysql.NullTime
		var accountType string
		a := domain.Account{}
		fields := []interface{}{
			&a.UUID,
			&accountType,
			&a.GivenName,
			&a.Surname,
			&a.Bio,
			&a.Email,
			&a.Phone,
			&a.PrimaryCredential,
			&a.SecondaryCredential,
			&a.NotificationTargetUUID,
			&a.CreatedAt,
			&a.UpdatedAt,
			&deletedAt,
		}
		if err := rows.Scan(fields...); err != nil {
			return matches, err
		}
		if deletedAt.Valid {
			a.DeletedAt = &deletedAt.Time
		}
		t, err := domain.AccountTypeFromString(accountType)
		if err != nil {
			return matches, err
		}
		a.Type = t
		matches = append(matches, a)
	}

	return matches, nil
}

func (q findAccountByUUID) confirmationsByAccount(
	uuids []u.UUID,
) (map[u.UUID][]domain.Confirmation, error) {

	// gaurd.
	matches := make(map[u.UUID][]domain.Confirmation)
	if len(uuids) < 1 {
		return matches, nil
	}

	uStrs := []interface{}{}
	qParams := []string{}
	for _, uuid := range uuids {
		uStrs = append(uStrs, uuid.String())
		qParams = append(qParams, "?")
	}
	statement, err :=
		q.db.Prepare(fmt.Sprintf("SELECT * FROM CONFIRMATION WHERE ACCOUNT_UUID IN (%s)", strings.Join(qParams, ", ")))
	if err != nil {
		return matches, err
	}
	defer statement.Close()

	rows, err := statement.Query(uStrs...)
	if err != nil {
		return matches, err
	}
	defer rows.Close()

	for rows.Next() {
		var confirmedAt mysql.NullTime
		var confirmationType string
		var accountUUID u.UUID
		var notificationUUID sql.NullString
		c := domain.Confirmation{}
		fields := []interface{}{
			&c.ID,
			&accountUUID,
			&c.Code,
			&confirmationType,
			&notificationUUID,
			&c.ExpiredAt,
			&confirmedAt,
			&c.CreatedAt,
			&c.UpdatedAt,
		}
		if err := rows.Scan(fields...); err != nil {
			return matches, err
		}
		if confirmedAt.Valid {
			c.ConfirmedAt = &confirmedAt.Time
		}
		t, err := domain.NewConfirmationTypeFromString(confirmationType)
		if err != nil {
			return matches, err
		}
		c.Type = t
		if notificationUUID.Valid {
			uuid, err :=
				u.FromString(notificationUUID.String)
			if err != nil {
				return matches, err
			}
			c.NotificationUUID = &uuid
		}
		if _, ok := matches[accountUUID]; !ok {
			matches[accountUUID] = []domain.Confirmation{}
		}
		matches[accountUUID] = append(matches[accountUUID], c)
	}

	return matches, nil
}
