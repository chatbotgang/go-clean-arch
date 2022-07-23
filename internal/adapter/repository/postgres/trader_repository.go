package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"

	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/barter"
	"github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
)

type repoTrader struct {
	ID        int       `db:"id"`
	UID       string    `db:"uid"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const repoTableTrader = "trader"

type repoColumnPatternTrader struct {
	ID        string
	UID       string
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

var repoColumnTrader = repoColumnPatternTrader{
	ID:        "id",
	UID:       "uid",
	Email:     "email",
	Name:      "name",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

func (c *repoColumnPatternTrader) columns() string {
	return strings.Join([]string{
		c.ID,
		c.UID,
		c.Email,
		c.Name,
		c.CreatedAt,
		c.UpdatedAt,
	}, ", ")
}

func (r *PostgresRepository) CreateTrader(ctx context.Context, param barter.Trader) (*barter.Trader, common.Error) {
	insert := map[string]interface{}{
		repoColumnTrader.UID:   param.UID,
		repoColumnTrader.Email: param.Email,
		repoColumnTrader.Name:  param.Name,
	}
	// build SQL query
	query, args, err := r.pgsq.Insert(repoTableTrader).
		SetMap(insert).
		Suffix(fmt.Sprintf("returning %s", repoColumnTrader.columns())).
		ToSql()
	if err != nil {
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	row := repoTrader{}
	if err = r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	// map the query result back to domain model
	trader := barter.Trader(row)
	return &trader, nil
}

func (r *PostgresRepository) GetTraderByEmail(ctx context.Context, email string) (*barter.Trader, common.Error) {
	query, args, err := r.pgsq.Select(repoColumnTrader.columns()).
		From(repoTableTrader).
		Where(sq.Eq{repoColumnTrader.Email: email}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}
	row := repoTrader{}

	// get one row from result
	if err = r.db.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.NewError(common.ErrorCodeResourceNotFound, err, common.WithMsg("account is not found"))
		}
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	// map the query result back to domain model
	trader := barter.Trader(row)
	return &trader, nil
}
