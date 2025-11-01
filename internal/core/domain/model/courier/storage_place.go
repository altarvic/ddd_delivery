package courier

import (
	"errors"
	"github.com/google/uuid"
	"strings"
)

// StoragePlace - место хранения курьера (рюкзак, багажник и т.п.),
type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name")
	}

	if totalVolume <= 0 {
		return nil, errors.New("totalVolume")
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
		orderID:     nil,
	}, nil
}

func (s *StoragePlace) Id() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
	return other != nil && s.id == other.id
}

func (s *StoragePlace) CanStore(volume int) bool {
	return volume <= s.totalVolume
}

func (s *StoragePlace) Store(orderId uuid.UUID, volume int) error {

	if s.orderID != nil {
		return errors.New("storage place already occupied")
	}

	if volume > s.totalVolume {
		return errors.New("volume must not exceed totalVolume")
	}

	s.orderID = &orderId

	return nil
}

func (s *StoragePlace) Clear() {
	s.orderID = nil
}

func (s *StoragePlace) IsOccupied() bool {
	return s.orderID != nil
}

// RestoreStoragePlace should be used ONLY inside Repository
func RestoreStoragePlace(id uuid.UUID, name string, totalVolume int, orderID *uuid.UUID) *StoragePlace {
	return &StoragePlace{
		id:          id,
		name:        name,
		totalVolume: totalVolume,
		orderID:     orderID,
	}
}
