package postgres

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ ports.CourierRepository = &courierRepository{}

type courierRepository struct {
	tx pgx.Tx
}

func NewCourierRepository(tx pgx.Tx) (ports.CourierRepository, error) {
	if tx == nil {
		return nil, errs.NewValueIsRequiredError("tx")
	}

	return &courierRepository{
		tx: tx,
	}, nil
}

func (cr *courierRepository) Get(ctx context.Context, id uuid.UUID) (*courier.Courier, error) {

	query := `select c.id,
					 c.name,
					 c.speed,
					 c.location_x,
					 c.location_y,
					 sp.id,
					 sp.name,
					 sp.volume,
					 sp.order_id
			  from couriers c inner join storage_places sp on c.id = sp.courier_id
			  where c.id = $1
              order by c.id, sp.id`

	rows, err := cr.tx.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer rows.Close()

	cDTO := courierDTO{StoragePlaces: make([]storagePlaceDTO, 0, 10)}
	spDTO := storagePlaceDTO{}

	for rows.Next() {

		err = rows.Scan(&cDTO.Id, &cDTO.Name, &cDTO.Speed, &cDTO.LocationX, &cDTO.LocationY,
			&spDTO.Id, &spDTO.Name, &spDTO.Volume, &spDTO.OrderId)

		if err != nil {
			return nil, err
		}

		cDTO.StoragePlaces = append(cDTO.StoragePlaces, spDTO)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(cDTO.StoragePlaces) == 0 {
		return nil, nil // not found (no error here)
	}

	return cDTO.ToCourier(), nil
}

func (cr courierRepository) GetAllFree(ctx context.Context) ([]*courier.Courier, error) {

	query := `select c.id,
			  	     c.name,
			  	     c.speed,
			  	     c.location_x,
			  	     c.location_y,
			  	     sp.id,
			  	     sp.name,
			  	     sp.volume,
			  	     sp.order_id
			  from couriers c
			  		 inner join storage_places sp on c.id = sp.courier_id
			  where not exists (select null
			  				    from storage_places sp2
			  				    where sp2.courier_id = c.id
			  					  and sp2.order_id is not null)
			  order by c.id, sp.id`

	rows, err := cr.tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer rows.Close()

	courierDTOs := make([]courierDTO, 0, 100)
	currentCourier := (*courierDTO)(nil)

	for rows.Next() {

		cDTO := courierDTO{StoragePlaces: make([]storagePlaceDTO, 0, 10)}
		spDTO := storagePlaceDTO{}

		err = rows.Scan(&cDTO.Id, &cDTO.Name, &cDTO.Speed, &cDTO.LocationX, &cDTO.LocationY,
			&spDTO.Id, &spDTO.Name, &spDTO.Volume, &spDTO.OrderId)

		if err != nil {
			return nil, err
		}

		if currentCourier == nil || currentCourier.Id != cDTO.Id {
			courierDTOs = append(courierDTOs, cDTO)
			currentCourier = &courierDTOs[len(courierDTOs)-1]
		}

		currentCourier.StoragePlaces = append(currentCourier.StoragePlaces, spDTO)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if currentCourier == nil {
		return nil, nil // not found (no error here)
	}

	couriers := make([]*courier.Courier, 0, len(courierDTOs))
	for _, cDTO := range courierDTOs {
		couriers = append(couriers, cDTO.ToCourier())
	}

	return couriers, nil
}

func (cr *courierRepository) Save(ctx context.Context, couriers ...*courier.Courier) error {

	cQuery := `insert into couriers (id, name, speed, location_x, location_y)
	 		   values ($1, $2, $3, $4, $5)
			   on conflict (id)
				  do update set name       = EXCLUDED.name,
					    	    speed      = EXCLUDED.speed,
							    location_x = EXCLUDED.location_x,
							    location_y = EXCLUDED.location_y;`

	spQuery := `insert into storage_places (id, name, volume, order_id, courier_id)
				values ($1, $2, $3, $4, $5)
				on conflict (id)
				   do update set name       = EXCLUDED.name,
								 volume     = EXCLUDED.volume,
								 order_id   = EXCLUDED.order_id,
								 courier_id = EXCLUDED.courier_id;`

	for _, c := range couriers {

		_, err := cr.tx.Exec(ctx, cQuery, c.Id(), c.Name(), c.Speed(), c.Location().X(), c.Location().Y())
		if err != nil {
			return err
		}

		for _, sp := range c.StoragePlaces() {
			_, err = cr.tx.Exec(ctx, spQuery, sp.Id(), sp.Name(), sp.TotalVolume(), sp.OrderID(), c.Id())
			if err != nil {
				return err
			}
		}
	}

	return nil
}
