package infrastructure

import (
	"database/sql"
	"tri-fitness/genesis/domain"
)

type AccountQuery interface {
	Execute() ([]domain.Account, error)
}

type accountQuery struct {
	db   *sql.DB
	skip *int
	take *int
}
