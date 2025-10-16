package courirer

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"errors"
	"github.com/google/uuid"
	"math"
	"strings"
)

type Courier struct {
	id            uuid.UUID
	name          string
	speed         int
	location      kernel.Location
	storagePlaces []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("empty name")
	}

	if speed <= 0 {
		return nil, errors.New("speed <= 0")
	}

	if location.IsEmpty() {
		return nil, errors.New("location is empty")
	}

	bag, _ := NewStoragePlace("Bag", 10)

	return &Courier{
		id:            uuid.New(),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: []*StoragePlace{bag},
	}, nil
}

func (c *Courier) Id() uuid.UUID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) StoragePlaces() []*StoragePlace {
	return c.storagePlaces
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	s, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}

	c.storagePlaces = append(c.storagePlaces, s)
	return nil
}

func (c *Courier) CanTakeOrder(o *order.Order) (bool, error) {
	if !o.Status().IsValid() {
		return false, errors.New("invalid order status")
	}

	if o.Status() == order.StatusAssigned {
		return false, errors.New("order is already assigned")
	}

	if o.Status() == order.StatusCompleted {
		return false, errors.New("order is completed")
	}

	for _, place := range c.StoragePlaces() {
		if !place.IsOccupied() && place.CanStore(o.Volume()) {
			return true, nil
		}
	}

	return false, errors.New("storage place not found")
}

func (c *Courier) TakeOrder(o *order.Order) error {
	_, err := c.CanTakeOrder(o)

	if err != nil {
		return err
	}

	for _, place := range c.StoragePlaces() {
		if !place.IsOccupied() && place.CanStore(o.Volume()) {
			err = o.AssignCourier(c.Id())
			if err != nil {
				return err
			}

			err = place.Store(o.Id(), o.Volume())
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("storage place not found")
}

func (c *Courier) CompleteOrder(o *order.Order) error {

	var place *StoragePlace = nil
	for _, p := range c.StoragePlaces() {
		if p.OrderID() != nil && *p.OrderID() == o.Id() {
			place = p
			break
		}
	}

	if place == nil {
		return errors.New("non-owned order")
	}

	err := o.Complete()
	if err != nil {
		return err
	}

	place.Clear()
	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if target.IsEmpty() {
		return 0, errors.New("empty location")
	}

	distance, err := c.Location().DistanceTo(target)
	if err != nil {
		return 0, err
	}

	return roundFloat(float64(distance)/float64(c.speed), 3), nil
}

func (c *Courier) Move(target kernel.Location) error {
	if target.IsEmpty() {
		return errors.New("empty location")
	}

	if c.location.Equals(target) {
		return errors.New("on the target")
	}

	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}

	c.location = newLocation
	return nil
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

//
// func (c *Courier) Move(target kernel.Location) error {
//
// }
