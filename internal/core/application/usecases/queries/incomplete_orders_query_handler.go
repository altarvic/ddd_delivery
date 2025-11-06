package queries

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncompleteOrder struct {
	OrderID   uuid.UUID `db:"id"`
	LocationX int       `db:"location_x"`
	LocationY int       `db:"location_y"`
}

type IncompleteOrdersResponse []*IncompleteOrder

type IncompleteOrdersQueryHandler interface {
	Handle(context.Context) (IncompleteOrdersResponse, error)
}

var _ IncompleteOrdersQueryHandler = &incompleteOrdersQueryHandler{}

type incompleteOrdersQueryHandler struct {
	db *pgxpool.Pool
}

func NewIncompleteOrdersQueryHandler(db *pgxpool.Pool) (IncompleteOrdersQueryHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	return &incompleteOrdersQueryHandler{db: db}, nil
}

func (cq *incompleteOrdersQueryHandler) Handle(ctx context.Context) (IncompleteOrdersResponse, error) {

	rows, err := cq.db.Query(ctx,
		fmt.Sprintf(`select id, location_x, location_y 
							from orders 
							where status != '%s'`, order.StatusCompleted))

	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[IncompleteOrder])
}
