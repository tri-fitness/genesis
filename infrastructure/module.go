package infrastructure

import (
	"database/sql"
	"tri-fitness/genesis/config"
	"tri-fitness/genesis/domain"

	"github.com/freerware/work"
	"github.com/freerware/workfx"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
)

type DataMapperResult struct {
	fx.Out

	Inserters map[work.TypeName]work.Inserter
	Updaters  map[work.TypeName]work.Updater
	Deleters  map[work.TypeName]work.Deleter
}

type DBResult struct {
	fx.Out

	DB *sql.DB `name:"rwDB"`
}

type DataMapperParameters struct {
	fx.In

	DB *sql.DB `name:"rwDB"`
}

var Module = fx.Options(
	fx.Provide(NewQueryFactory),
	fx.Provide(NewRepositoryFactory),
	fx.Provide(func(c config.Configuration) (DBResult, error) {
		dsn := c.Database.DSN()
		db, err := sql.Open("mysql", dsn)
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
	fx.Provide(func(parameters DataMapperParameters) DataMapperResult {
		inserters := make(map[work.TypeName]work.Inserter)
		updaters := make(map[work.TypeName]work.Updater)
		deleters := make(map[work.TypeName]work.Deleter)
		accountTN := work.TypeNameOf(domain.Account{})
		dm := NewAccountDataMapper(AccountDataMapperParameters{
			DB: parameters.DB,
		})
		inserters[accountTN] = &dm
		updaters[accountTN] = &dm
		deleters[accountTN] = &dm
		result := DataMapperResult{
			Inserters: inserters,
			Updaters:  updaters,
			Deleters:  deleters,
		}
		return result
	}),
	workfx.Modules.SQLUnit,
)
