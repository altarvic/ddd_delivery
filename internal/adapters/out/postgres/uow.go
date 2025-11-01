package postgres

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	or ports.OrderRepository
	cr ports.CourierRepository
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
	if uowc.or != nil {
		return uowc.or
	}

	uowc.or, _ = NewOrderRepository(uowc.tx)
	return uowc.or
}

func (uowc *unitOfWorkComponents) CourierRepository() ports.CourierRepository {
	if uowc.cr != nil {
		return uowc.cr
	}

	uowc.cr, _ = NewCourierRepository(uowc.tx)
	return uowc.cr
}
