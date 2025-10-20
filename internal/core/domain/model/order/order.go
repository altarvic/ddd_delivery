package order

import (
	"delivery/internal/core/domain/kernel"
	"errors"
	"github.com/google/uuid"
)

type Order struct {
	id        uuid.UUID
	courierId *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func (o *Order) Id() uuid.UUID {
	return o.id
}

func (o *Order) CourierId() *uuid.UUID {
	return o.courierId
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func NewOrder(orderId uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderId == uuid.Nil {
		return nil, errors.New("empty orderId")
	}

	if location.IsEmpty() {
		return nil, errors.New("empty location")
	}

	if volume <= 0 {
		return nil, errors.New("volume <= 0")
	}

	return &Order{
		id:       orderId,
		location: location,
		volume:   volume,
		status:   StatusCreated,
	}, nil
}

func (o *Order) Equals(other *Order) bool {
	return other != nil && o.id == other.id
}

func (o *Order) AssignCourier(courierId uuid.UUID) error {
	if o.status != StatusCreated {
		return errors.New("courier already assigned")
	}

	if courierId == uuid.Nil {
		return errors.New("empty courierId")
	}

	o.courierId = &courierId
	o.status = StatusAssigned
	return nil
}

func (o *Order) Complete() error {
	if o.status == StatusCompleted {
		return errors.New("already completed")
	}

	if o.status == StatusCreated {
		return errors.New("order w/o assigned courier")
	}

	o.status = StatusCompleted
	return nil
}
