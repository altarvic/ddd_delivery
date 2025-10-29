package postgres

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ ports.OrderRepository = &orderRepository{}

type orderRepository struct {
	tx pgx.Tx
}

func NewOrderRepository(tx pgx.Tx) (ports.OrderRepository, error) {
	if tx == nil {
		return nil, errs.NewValueIsRequiredError("tx")
	}

	return &orderRepository{
		tx: tx,
	}, nil
}

func (or *orderRepository) Save(ctx context.Context, orders ...*order.Order) error {

	query := `insert into orders (id, courier_id, location_x, location_y, volume, status)
			  values ($1, $2, $3, $4, $5, $6)
			 	  on conflict (id)
					 do update set courier_id = EXCLUDED.courier_id,
								   location_x = EXCLUDED.location_x,
								   location_y = EXCLUDED.location_y,
								   volume     = EXCLUDED.volume,
								   status     = EXCLUDED.status;`

	for _, o := range orders {
		_, err := or.tx.Exec(ctx, query, o.Id(), o.CourierId(),
			o.Location().X(), o.Location().Y(), o.Volume(), o.Status())

		if err != nil {
			return err
		}
	}

	return nil
}

func (or *orderRepository) Get(ctx context.Context, id uuid.UUID) (*order.Order, error) {

	query := `select id, courier_id, location_x, location_y, volume, status
			  from orders
			  where id = $1`

	var dto = orderDTO{}
	err := or.tx.QueryRow(ctx, query, id).
		Scan(&dto.Id, &dto.CourierId, &dto.LocationX, &dto.LocationY, &dto.Volume, &dto.Status)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found (no error here)
		} else {
			return nil, err
		}
	}

	return dto.ToOrder(), nil
}

func (or *orderRepository) GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error) {

	query := fmt.Sprintf(`select id, courier_id, location_x, location_y, volume, status
						    	from orders
							    where status = '%s'
							    limit 1`, order.StatusCreated)

	var dto = orderDTO{}
	err := or.tx.QueryRow(ctx, query).
		Scan(&dto.Id, &dto.CourierId, &dto.LocationX, &dto.LocationY, &dto.Volume, &dto.Status)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found (no error here)
		} else {
			return nil, err
		}
	}

	return dto.ToOrder(), nil
}

func (or *orderRepository) GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error) {

	query := fmt.Sprintf(`select id, courier_id, location_x, location_y, volume, status
			                    from orders
			                    where status = '%s'
                                order by id`, order.StatusAssigned)

	rows, err := or.tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer rows.Close()

	orders := make([]*order.Order, 0, 100)
	for rows.Next() {

		var dto = orderDTO{}
		err = rows.Scan(&dto.Id, &dto.CourierId, &dto.LocationX, &dto.LocationY, &dto.Volume, &dto.Status)
		if err != nil {
			return nil, err
		}

		orders = append(orders, dto.ToOrder())
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
