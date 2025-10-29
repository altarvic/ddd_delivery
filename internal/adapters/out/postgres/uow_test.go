package postgres

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"errors"
	"testing"
)

func TestUnitOfWork_Do_WithCommit(t *testing.T) {

	ctx, db, uow, err := setupTest(t, false)
	if err != nil {
		t.Fatal(err)
	}

	// создаем курьеров и заказы
	couriers := createCouriers(25)
	orders := createOrders(100)
	dispatcher := services.NewOrderDispatcher()

	// случайным образом назначаем заказы курьерам (примерно 66% из них)
	for range len(couriers) - len(couriers)/3 {
		_, _ = dispatcher.Dispatch(getRandomUnassignedOrder(orders), couriers)
	}

	// сохраняем курьеров и заказы в БД
	err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		// курьеров
		err = uowc.CourierRepository().Save(ctx, couriers...)
		if err != nil {
			return err
		}

		// заказы
		err = uowc.OrderRepository().Save(ctx, orders...)
		if err != nil {
			return err
		}

		// коммитим транзакцию
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	// проверяем, что заказы были сохранены корректно
	query := `select id, 
					 courier_id, 
					 location_x, 
					 location_y, 
					 volume, 
					 status
		      from orders
		      order by id`

	err = compareRowsAndAggregates(db, ctx, query, orders, equalOrders)
	if err != nil {
		t.Fatal(err)
	}

	// проверяем, что курьеры были сохранены корректно
	query = `select id,
		 		    name,
					speed,
					location_x,
					location_y
			 from couriers
			 order by id`

	err = compareRowsAndAggregates(db, ctx, query, couriers, equalCouriers)
	if err != nil {
		t.Fatal(err)
	}

	// проверяем, что места хранения были сохранены корректно
	query = `select id,
		 		    name,
					volume,
					order_id
			 from storage_places
			 order by id`

	storagePlaces := make([]*courier.StoragePlace, 0, len(couriers)*2)
	for _, c := range couriers {
		storagePlaces = append(storagePlaces, c.StoragePlaces()...)
	}
	sortById(storagePlaces)

	err = compareRowsAndAggregates(db, ctx, query, storagePlaces, equalStoragePlaces)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnitOfWork_Do_WithRollback(t *testing.T) {

	ctx, db, uow, err := setupTest(t, false)
	if err != nil {
		t.Fatal(err)
	}

	// создаем курьеров и заказы
	couriers := createCouriers(2)
	orders := createOrders(10)

	// сохраняем курьеров и заказы в БД, но делаем Rollback
	err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		err = uowc.CourierRepository().Save(ctx, couriers...)
		if err != nil {
			return err
		}

		err = uowc.OrderRepository().Save(ctx, orders...)
		if err != nil {
			return err
		}

		// делаем Rollback
		return errors.New("undo")
	})

	if err == nil {
		t.Fatal("expected error")
	}

	// проверяем, что заказов и курьеров нет в БД
	query := `select 
				(select count(*) from couriers) +
				(select count(*) from orders)`

	var count int
	err = db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatal("expected empty tables")
	}
}

func TestUnitOfWork_Do_WithRollbackInNestedTransaction(t *testing.T) {

	ctx, db, uow, err := setupTest(t, false)
	if err != nil {
		t.Fatal(err)
	}

	couriers := createCouriers(3)
	orders := createOrders(10)

	err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		// сохраняем заказы во вложенной транзакции, но откатываем ее
		err = uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {
			err = uowc.OrderRepository().Save(ctx, orders...)
			if err != nil {
				return err
			}

			// отменяем
			return errors.New("undo orders")
		})

		if err == nil {
			return errors.New("expected error")
		}

		// сохраняем курьеров
		err = uowc.CourierRepository().Save(ctx, couriers...)
		if err != nil {
			return err
		}

		// commit
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	// проверяем заказы
	var count int
	err = db.QueryRow(ctx, "select count(*) from orders").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != 0 {
		t.Fatal("orders should be 0")
	}

	// проверяем курьеров
	err = db.QueryRow(ctx, "select count(*) from couriers").Scan(&count)
	if err != nil {
		t.Fatal(err)
	}

	if count != len(couriers) {
		t.Fatalf("couriers should be %d", len(couriers))
	}
}
