package infrastructure

import (
	"database/sql"
	"tri-fitness/genesis/domain"

	"github.com/go-sql-driver/mysql"
)

type findAccountByPrimaryCredential struct {
	accountQuery

	credential string
}

type FindAccountByPrimaryCredentialParameters struct {
	DB         *sql.DB
	Credential string
}

func NewFindAccountByPrimaryCredentialQuery(
	parameters FindAccountByPrimaryCredentialParameters,
) AccountQuery {
	return &findAccountByPrimaryCredential{
		accountQuery: accountQuery{
			db: parameters.DB,
		},
		credential: parameters.Credential,
	}
}

func (q *findAccountByPrimaryCredential) Execute() ([]domain.Account, error) {
	matches := []domain.Account{}
	statement, err :=
		q.db.Prepare(
			"SELECT * FROM ACCOUNT WHERE PRIMARY_CREDENTIAL = ?")
	if err != nil {
		return matches, err
	}
	defer statement.Close()

	rows, err := statement.Query(q.credential)
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
