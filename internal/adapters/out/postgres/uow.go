package postgres

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

var _ ports.UnitOfWork = &unitOfWork{}
var _ ports.UnitOfWorkComponents = &unitOfWorkComponents{}

type unitOfWork struct {
	db *pgxpool.Pool
}

type txKeyType struct{}

var txKey = txKeyType{}

type unitOfWorkComponents struct {
	tx pgx.Tx
}

func NewUnitOfWork(db *pgxpool.Pool) (ports.UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	uow := &unitOfWork{
		db: db,
	}

	return uow, nil
}

func (u *unitOfWork) Do(ctx context.Context, fn ports.UnitOfWorkDoFunc) error {

	tx := u.getCurrentTx(ctx)

	var err error
	if tx == nil { // create new transaction
		tx, err = u.db.Begin(ctx)
		if err != nil {
			return err
		}
	} else { // create nested transaction (savepoint)
		tx, err = tx.Begin(ctx)
		if err != nil {
			return err
		}
	}

	ctx = context.WithValue(ctx, txKey, tx)

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	err = fn(ctx, &unitOfWorkComponents{tx: tx})
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func (u *unitOfWork) getCurrentTx(ctx context.Context) pgx.Tx {
	val := ctx.Value(txKey)
	if tx, ok := val.(pgx.Tx); ok {
		return tx
	}

	return nil
}

func (uowc *unitOfWorkComponents) OrderRepository() ports.OrderRepository {
	return sync.OnceValue(
		func() ports.OrderRepository {
			repo, _ := NewOrderRepository(uowc.tx)
			return repo
		})()
}

func (uowc *unitOfWorkComponents) CourierRepository() ports.CourierRepository {
	return sync.OnceValue(
		func() ports.CourierRepository {
			repo, _ := NewCourierRepository(uowc.tx)
			return repo
		})()
}
