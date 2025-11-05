package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"errors"
	"math"
)

type OrderDispatcher interface {
	Dispatch(*order.Order, []*courier.Courier) (*courier.Courier, error)
}

var _ OrderDispatcher = &orderDispatcher{}

type orderDispatcher struct {
}

func NewOrderDispatcher() OrderDispatcher {
	return &orderDispatcher{}
}

func (od *orderDispatcher) Dispatch(o *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if o.Status() != order.StatusCreated {
		return nil, errors.New("invalid order status")
	}

	fastestDeliveryTime := math.MaxFloat64
	fastestCourier := (*courier.Courier)(nil)

	for _, c := range couriers {

		if ok, _ := c.CanTakeOrder(o); !ok {
			continue
		}

		deliveryTime, err := c.CalculateTimeToLocation(o.Location())
		if err != nil {
			continue
		}

		if deliveryTime < fastestDeliveryTime {
			fastestDeliveryTime = deliveryTime
			fastestCourier = c
		}
	}

	if fastestCourier == nil {
		return nil, errors.New("no matching courier")
	}

	err := fastestCourier.TakeOrder(o)
	if err != nil {
		return nil, err
	}

	err = o.AssignCourier(fastestCourier.Id())
	if err != nil {
		return nil, err
	}

	return fastestCourier, nil
}
