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

type repoGood struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	OwnerID   int       `db:"owner_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type repoColumnPatternGood struct {
	ID        string
	Name      string
	OwnerID   string
	CreatedAt string
	UpdatedAt string
}

const repoTableGood = "good"

var repoColumnGood = repoColumnPatternGood{
	ID:        "id",
	Name:      "name",
	OwnerID:   "owner_id",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

func (c *repoColumnPatternGood) columns() string {
	return strings.Join([]string{
		c.ID,
		c.Name,
		c.OwnerID,
		c.CreatedAt,
		c.UpdatedAt,
	}, ", ")
}

func (r *PostgresRepository) CreateGood(ctx context.Context, param barter.Good) (*barter.Good, common.Error) {
	insert := map[string]interface{}{
		repoColumnGood.Name:    param.Name,
		repoColumnGood.OwnerID: param.OwnerID,
	}

	// build SQL query
	query, args, err := r.pgsq.Insert(repoTableGood).
		SetMap(insert).
		Suffix(fmt.Sprintf("returning %s", repoColumnGood.columns())).
		ToSql()
	if err != nil {
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	var row repoGood
	if err = r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	good := barter.Good(row)
	return &good, nil
}

func (r *PostgresRepository) GetGoodByID(ctx context.Context, id int) (*barter.Good, common.Error) {
	where := sq.And{
		sq.Eq{repoColumnGood.ID: id},
	}

	// build SQL query
	query, args, err := r.pgsq.Select(repoColumnGood.columns()).
		From(repoTableGood).
		Where(where).
		Limit(1).
		ToSql()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, common.NewError(common.ErrorCodeResourceNotFound, err)
		}
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	var row repoGood
	if err = r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	good := barter.Good(row)
	return &good, nil
}

func (r *PostgresRepository) ListGoods(ctx context.Context) ([]barter.Good, common.Error) {
	return r.listGoods(ctx, r.db, sq.And{})
}

func (r *PostgresRepository) ListGoodsByOwner(ctx context.Context, ownerID int) ([]barter.Good, common.Error) {
	return r.listGoods(ctx, r.db, sq.And{
		sq.Eq{repoColumnGood.OwnerID: ownerID},
	})
}

func (r *PostgresRepository) listGoods(ctx context.Context, db sqlContextGetter, where sq.And) ([]barter.Good, common.Error) {
	// build SQL query
	query, args, err := r.pgsq.Select(repoColumnGood.columns()).
		From(repoTableGood).
		Where(where).
		OrderBy(fmt.Sprintf("%s desc", repoColumnGood.CreatedAt)).
		ToSql()
	if err != nil {
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	var rows []repoGood
	if err = db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	var goods []barter.Good
	for _, row := range rows {
		good := barter.Good(row)
		goods = append(goods, good)
	}

	return goods, nil
}

func (r *PostgresRepository) UpdateGood(ctx context.Context, good barter.Good) (*barter.Good, common.Error) {
	where := sq.And{
		sq.Eq{repoColumnGood.ID: good.ID},
	}

	update := map[string]interface{}{
		repoColumnGood.Name:      good.Name,
		repoColumnGood.OwnerID:   good.OwnerID,
		repoColumnGood.UpdatedAt: time.Now(),
	}

	// build SQL query
	query, args, err := r.pgsq.Update(repoTableGood).
		SetMap(update).
		Where(where).
		Suffix(fmt.Sprintf("returning %s", repoColumnGood.columns())).
		ToSql()
	if err != nil {
		return nil, common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	var row repoGood
	if err = r.db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}

	updatedGood := barter.Good(row)
	return &updatedGood, nil
}

func (r *PostgresRepository) DeleteGoodByID(ctx context.Context, id int) common.Error {
	where := sq.And{
		sq.Eq{repoColumnGood.ID: id},
	}

	// build SQL query
	query, args, err := r.pgsq.Delete(repoTableGood).
		Where(where).
		ToSql()
	if err != nil {
		return common.NewError(common.ErrorCodeInternalProcess, err)
	}

	// execute SQL query
	if _, err = r.db.ExecContext(ctx, query, args...); err != nil {
		return common.NewError(common.ErrorCodeRemoteProcess, err)
	}
	return nil
}
