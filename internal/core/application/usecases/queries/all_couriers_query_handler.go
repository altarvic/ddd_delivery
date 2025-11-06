package queries

import (
	"context"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourierResponse struct {
	CourierID uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	LocationX int       `db:"location_x"`
	LocationY int       `db:"location_y"`
}

type AllCouriersResponse []*CourierResponse

type AllCouriersQueryHandler interface {
	Handle(context.Context) (AllCouriersResponse, error)
}

var _ AllCouriersQueryHandler = &allCouriersQueryHandler{}

type allCouriersQueryHandler struct {
	db *pgxpool.Pool
}

func NewAllCouriersQueryHandler(db *pgxpool.Pool) (AllCouriersQueryHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	return &allCouriersQueryHandler{db: db}, nil
}

func (aq *allCouriersQueryHandler) Handle(ctx context.Context) (AllCouriersResponse, error) {

	rows, err := aq.db.Query(ctx, "select id, name, location_x, location_y from couriers")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[CourierResponse])
}
