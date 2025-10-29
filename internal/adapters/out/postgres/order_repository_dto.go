package postgres

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
)

type orderDTO struct {
	Id        uuid.UUID    `db:"id"`
	CourierId *uuid.UUID   `db:"courier_id"`
	LocationX int          `db:"location_x"`
	LocationY int          `db:"location_y"`
	Volume    int          `db:"volume"`
	Status    order.Status `db:"status"`
}

func (dto *orderDTO) ToOrder() *order.Order {
	loc := kernel.RestoreLocation(dto.LocationX, dto.LocationY)
	return order.RestoreOrder(dto.Id, dto.CourierId, loc, dto.Volume, dto.Status)
}
