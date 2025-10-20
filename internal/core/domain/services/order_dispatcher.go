package services

import (
	"delivery/internal/core/domain/model/courirer"
	"delivery/internal/core/domain/model/order"
	"errors"
	"math"
)

type OrderDispatcher interface {
	Dispatch(*order.Order, []*courirer.Courier) (*courirer.Courier, error)
}

var _ OrderDispatcher = &orderDispatcher{}

type orderDispatcher struct {
}

func NewOrderDispatcher() OrderDispatcher {
	return &orderDispatcher{}
}

func (od *orderDispatcher) Dispatch(o *order.Order, couriers []*courirer.Courier) (*courirer.Courier, error) {
	if o.Status() != order.StatusCreated {
		return nil, errors.New("invalid order status")
	}

	fastestDeliveryTime := math.MaxFloat64
	fastestCourier := (*courirer.Courier)(nil)

	for _, courier := range couriers {

		if ok, _ := courier.CanTakeOrder(o); !ok {
			continue
		}

		deliveryTime, err := courier.CalculateTimeToLocation(o.Location())
		if err != nil {
			continue
		}

		if deliveryTime < fastestDeliveryTime {
			fastestDeliveryTime = deliveryTime
			fastestCourier = courier
		}
	}

	if fastestCourier == nil {
		return nil, errors.New("no matching courier")
	}

	if err := fastestCourier.TakeOrder(o); err != nil {
		return nil, err
	}

	return fastestCourier, nil
}
