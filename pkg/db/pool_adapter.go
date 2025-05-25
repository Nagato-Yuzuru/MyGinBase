package db

import (
	"entgo.io/ent/dialect/sql"
	"errors"
	"github.com/jackc/pgx/v5/stdlib"
	"terraqt.io/colas/bedrock-go/pkg/errs"
)

func provideDriver(pool PGPool) (*sql.Driver, error) {
	adp := provideAdaptPool(pool)

	if adp == nil {
		return nil, errs.WrapCodeError(errs.ErrNotImplemented, errors.New("adp is nil"))
	}

	db := stdlib.OpenDBFromPool(adp.getStdPool())

	driver := sql.OpenDB("postgres", db)

	return driver, nil

}

func provideAdaptPool(pool PGPool) adaptPool {
	if pool, ok := pool.(*postgresPool); ok {
		return pool
	}
	return nil
}
