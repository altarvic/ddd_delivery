package postgres

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"github.com/google/uuid"
	"testing"
)

func TestOrderRepository_Get(t *testing.T) {

	ctx, _, uow, err := setupTest(t, true)
	if err != nil {
		t.Fatal(err)
	}

	orderId := uuid.MustParse("0f9fd652-580d-4d44-8851-2e22daef93fe")

	var o *order.Order
	err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {
		o, err = uowc.OrderRepository().Get(ctx, orderId)
		return err
	})

	if err != nil {
		t.Fatal(err)
	}

	if o == nil {
		t.Fatalf("expected order %s exists", orderId.String())
	}

	// проверяем данные заказа
	if o.Id() != orderId ||
		o.CourierId() == nil ||
		*o.CourierId() != uuid.MustParse("54517cca-9ac1-4b49-a649-606aae75b621") ||
		o.Volume() != 11 ||
		o.Location().X() != 9 ||
		o.Location().Y() != 10 ||
		o.Status() != order.StatusAssigned {
		t.Fatal("wrong order data")
	}
}

func TestOrderRepository_GetFirstInCreatedStatus(t *testing.T) {

	ctx, _, uow, err := setupTest(t, true)
	if err != nil {
		t.Fatal(err)
	}

	var o *order.Order
	err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {
		o, err = uowc.OrderRepository().GetFirstInCreatedStatus(ctx)
		return err
	})

	if err != nil {
		t.Fatal(err)
	}

	if o == nil {
		t.Fatal("expected order exists")
	}

	// проверяем данные заказа
	if o.Status() != order.StatusCreated {
		t.Fatal("expected order has created status")
	}
}

func TestOrderRepository_GetAllInAssignedStatus(t *testing.T) {

	ctx, _, uow, err := setupTest(t, true)
	if err != nil {
		t.Fatal(err)
	}

	var assignedOrders []*order.Order
	err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {
		assignedOrders, err = uowc.OrderRepository().GetAllInAssignedStatus(ctx)
		return err
	})

	if err != nil {
		t.Fatal(err)
	}

	// проверяем данные
	if assignedOrders == nil || len(assignedOrders) != 17 {
		t.Fatal("expected 17 assigned orders")
	}

	_, found := find(assignedOrders, func(o *order.Order) bool {
		return o.Id().String() == "b5261984-e50b-465e-bd65-4b2c578ad130"
	})

	if !found {
		t.Fatal("expected order b5261984-e50b-465e-bd65-4b2c578ad130 found")
	}

	_, found = find(assignedOrders, func(o *order.Order) bool {
		return o.Id().String() == "43de8fd5-53d8-41d1-8b45-f533dd39f80a"
	})

	if found {
		t.Fatal("expected order 43de8fd5-53d8-41d1-8b45-f533dd39f80a not found")
	}

}
