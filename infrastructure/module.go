package infrastructure

import (
	"database/sql"
	"time"
	"tri-fitness/genesis/config"
	"tri-fitness/genesis/domain"

	"github.com/freerware/work"
	"github.com/freerware/workfx"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SQLDataMapperResult struct {
	fx.Out

	SQLDataMappers map[work.TypeName]work.SQLDataMapper
}

type DBResult struct {
	fx.Out

	DB *sql.DB `name:"rwDB"`
}

type SQLDataMapperParameters struct {
	fx.In

	Logger *zap.Logger
}

var Module = fx.Options(
	fx.Provide(NewQueryFactory),
	fx.Provide(NewRepositoryFactory),
	fx.Provide(func(c config.Configuration) (DBResult, error) {
		time.Sleep(30000)
		dsn := c.Database.DSN()
		var current int
		retryCount := 3
		retryable := func() bool {
			return current < retryCount
		}
		var db *sql.DB
		var err error
		for retryable() {
			db, err = sql.Open("mysql", dsn)
			if err == nil {
				break
			}
			current = current + 1
		}
		if err != nil {
			return DBResult{}, err
		}
		err = db.Ping()
		if err != nil {
			return DBResult{}, err
		}
		return DBResult{
			DB: db,
		}, nil
	}),
	fx.Provide(func(parameters SQLDataMapperParameters) SQLDataMapperResult {
		dataMappers := make(map[work.TypeName]work.SQLDataMapper)
		accountTN := work.TypeNameOf(domain.Account{})
		dm := NewAccountDataMapper(AccountDataMapperParameters{
			Logger: parameters.Logger,
		})
		dataMappers[accountTN] = &dm
		result := SQLDataMapperResult{
			SQLDataMappers: dataMappers,
		}
		return result
	}),
	workfx.Modules.SQLUnit,
	fx.Provide(NewShortMessageService),
)
