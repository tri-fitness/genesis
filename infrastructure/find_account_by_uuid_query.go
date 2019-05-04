package infrastructure

import (
	"database/sql"
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

	matches := []domain.Account{}
	if err := q.db.Ping(); err != nil {
		return matches, err
	}
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
		a := domain.Account{}
		fields := []interface{}{
			&a.UUID,
			&a.Type,
			&a.GivenName,
			&a.Surname,
			&a.Bio,
			&a.Email,
			&a.Phone,
			&a.PrimaryCredential,
			&a.SecondaryCredential,
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
		matches = append(matches, a)
	}

	return matches, nil
}
