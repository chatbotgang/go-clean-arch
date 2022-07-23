package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	db   *sqlx.DB
	pgsq sq.StatementBuilderType
}

func NewPostgresRepository(ctx context.Context, db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
		// set the default placeholder as $ instead of ? because postgres uses $
		pgsq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// sqlContextGetter is an interface provided both by transaction and standard db connection
//type sqlContextGetter interface {
//	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
//	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
//	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
//}

//func (r *PostgresRepository) beginTx() (*sqlx.Tx, common.Error) {
//	tx, err := r.db.Beginx()
//	if err != nil {
//		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
//	}
//	return tx, nil
//}

// finishTx close an open transaction
// If error is provided, abort the transaction.
// If err is nil, commit the transaction.
//func (r *PostgresRepository) finishTx(err common.Error, tx *sqlx.Tx) common.Error {
//	if err != nil {
//		if rollbackErr := tx.Rollback(); rollbackErr != nil {
//			wrapError := multierror.Append(err, rollbackErr)
//			return common.NewError(common.ErrorCodeRemoteProcess, wrapError)
//		}
//
//		return err
//	} else {
//		if commitErr := tx.Commit(); commitErr != nil {
//			return common.NewError(common.ErrorCodeRemoteProcess, commitErr)
//		}
//
//		return nil
//	}
//}
