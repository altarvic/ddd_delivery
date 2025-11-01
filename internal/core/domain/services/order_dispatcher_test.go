package services

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
	"testing"
)

func TestOrderDispatcher_Dispatch(t *testing.T) {
	// create order dispatcher
	dispatcher := NewOrderDispatcher()

	// create 3 couriers
	loc, _ := kernel.NewLocation(1, 1)
	alice, _ := courier.NewCourier("Alice", 1, loc)

	loc, _ = kernel.NewLocation(7, 5)
	bob, _ := courier.NewCourier("Bob", 1, loc)

	loc, _ = kernel.NewLocation(3, 4)
	mallory, _ := courier.NewCourier("Mallory", 1, loc)

	couriers := []*courier.Courier{alice, bob, mallory}

	// create order that can't be taken (large volume)
	loc, _ = kernel.NewLocation(10, 10)
	o, _ := order.NewOrder(uuid.New(), loc, 100)

	// should be error
	_, err := dispatcher.Dispatch(o, couriers)
	if err == nil {
		t.Error("should be error")
	}

	// create regular order
	loc, _ = kernel.NewLocation(10, 10)
	o, _ = order.NewOrder(uuid.New(), loc, 8)

	// dispatch the order
	courier, err := dispatcher.Dispatch(o, couriers)
	if err != nil {
		t.Error(err)
	}

	// Bob had to take this order
	if courier.Id() != bob.Id() || o.CourierId() == nil || *o.CourierId() != bob.Id() {
		t.Fail()
	}

	// create another order
	loc, _ = kernel.NewLocation(10, 10)
	o, _ = order.NewOrder(uuid.New(), loc, 8)

	// dispatch the order
	courier, err = dispatcher.Dispatch(o, couriers)
	if err != nil {
		t.Error(err)
	}

	// Mallory had to take this order
	if courier.Id() != mallory.Id() || o.CourierId() == nil || *o.CourierId() != mallory.Id() {
		t.Fail()
	}

	// try to dispatch the same order
	courier, err = dispatcher.Dispatch(o, couriers)
	if err == nil {
		t.Error("order was already dispatched")
	}
}
