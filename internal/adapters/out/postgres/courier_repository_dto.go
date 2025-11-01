package postgres

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"github.com/google/uuid"
)

type storagePlaceDTO struct {
	Id      uuid.UUID  `db:"id"`
	Name    string     `db:"name"`
	Volume  int        `db:"volume"`
	OrderId *uuid.UUID `db:"order_id"`
}

func (dto *storagePlaceDTO) ToStoragePlace() *courier.StoragePlace {
	return courier.RestoreStoragePlace(dto.Id, dto.Name, dto.Volume, dto.OrderId)
}

type courierDTO struct {
	Id            uuid.UUID         `db:"id"`
	Name          string            `db:"name"`
	Speed         int               `db:"speed"`
	LocationX     int               `db:"location_x"`
	LocationY     int               `db:"location_y"`
	StoragePlaces []storagePlaceDTO `db:"-"`
}

func (dto *courierDTO) ToCourier() *courier.Courier {
	loc := kernel.RestoreLocation(dto.LocationX, dto.LocationY)

	var storagePlaces []*courier.StoragePlace
	for _, spDTO := range dto.StoragePlaces {
		storagePlaces = append(storagePlaces, spDTO.ToStoragePlace())
	}

	return courier.RestoreCourier(dto.Id, dto.Name, dto.Speed, loc, storagePlaces)
}
