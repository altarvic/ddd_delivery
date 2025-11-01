package postgres

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"github.com/google/uuid"
	"testing"
)

func TestCourierRepository_Get(t *testing.T) {

	ctx, _, uow, err := setupTest(t, true)
	if err != nil {
		t.Fatal(err)
	}

	courierId := uuid.MustParse("37886120-f01a-4b92-a3c6-931c1058bedf")

	var c *courier.Courier
	err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {
		c, err = uowc.CourierRepository().Get(ctx, courierId)
		return err
	})

	if err != nil {
		t.Fatal(err)
	}

	if c == nil {
		t.Fatalf("expected courier %s exists", courierId.String())
	}

	// проверяем данные курьера
	if c.Id() != courierId ||
		c.Name() != "courier14" ||
		c.Speed() != 4 ||
		c.Location().X() != 6 ||
		c.Location().Y() != 3 ||
		len(c.StoragePlaces()) != 2 {
		t.Fatal("wrong courier data")
	}

	sp := c.StoragePlaces()[0]
	if sp.Id() != uuid.MustParse("4e9e7f23-713a-49e6-941b-d2951713798a") ||
		sp.Name() != "Bag" ||
		sp.TotalVolume() != 10 ||
		sp.OrderID() != nil {
		t.Fatal("wrong storage place #1 data")
	}

	sp = c.StoragePlaces()[1]
	if sp.Id() != uuid.MustParse("a4719735-9cb7-4453-8157-4989a93e751d") ||
		sp.Name() != "trunk" ||
		sp.TotalVolume() != 200 ||
		sp.OrderID() != nil {
		t.Fatal("wrong storage place #2 data")
	}

}

func TestCourierRepository_GetAllFree(t *testing.T) {

	ctx, _, uow, err := setupTest(t, true)
	if err != nil {
		t.Fatal(err)
	}

	var freeCouriers []*courier.Courier
	err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {
		freeCouriers, err = uowc.CourierRepository().GetAllFree(ctx)
		return err
	})

	if err != nil {
		t.Fatal(err)
	}

	// проверяем данные
	if freeCouriers == nil || len(freeCouriers) != 12 {
		t.Fatal("expected 12 free couriers")
	}

	_, found := find(freeCouriers, func(c *courier.Courier) bool {
		return c.Id().String() == "ecc4c416-0fbb-44b2-9be8-4a988d01c431"
	})

	if !found {
		t.Fatal("expected courier ecc4c416-0fbb-44b2-9be8-4a988d01c431 found")
	}

	_, found = find(freeCouriers, func(c *courier.Courier) bool {
		return c.Id().String() == "0234c21c-e521-4f35-a5e3-e0af93c77bb8"
	})

	if found {
		t.Fatal("expected courier 0234c21c-e521-4f35-a5e3-e0af93c77bb8 not found")
	}

}
