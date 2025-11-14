package postgres

import (
	"bytes"
	"context"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"math/rand"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func setupTest(t *testing.T, seedTestData bool) (context.Context, *pgxpool.Pool, ports.UnitOfWork, error) {
	ctx := context.Background()

	var scripts []string
	if seedTestData {
		scripts = []string{"seed_data.sql"}
	} else {
		scripts = []string{}
	}

	// run test container with Postgres DB
	postgresContainer, dsn, err := startPostgresContainer(ctx, scripts...)
	if err != nil {
		return nil, nil, nil, err
	}

	defer func() {
		if err != nil {
			_ = testcontainers.TerminateContainer(postgresContainer)
		} else {
			// cleanup after test finished
			testcontainers.CleanupContainer(t, postgresContainer)
		}
	}()

	// connect to DB
	db, err := connectDb(ctx, dsn)
	if err != nil {
		return nil, nil, nil, err
	}

	// create Mediatr
	mediatr := ddd.NewMediatr()

	// create UOW
	uow, err := NewUnitOfWork(db, mediatr)
	if err != nil {
		return nil, nil, nil, err
	}

	return ctx, db, uow, nil
}

func startPostgresContainer(ctx context.Context, additionalScripts ...string) (testcontainers.Container, string, error) {

	scripts := []string{filepath.Join(".", "testdata", "db-schema.sql")}
	for _, s := range additionalScripts {
		scripts = append(scripts, filepath.Join(".", "testdata", s))
	}

	postgresContainer, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithOrderedInitScripts(scripts...),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)

	if err != nil {
		return nil, "", err
	}

	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = testcontainers.TerminateContainer(postgresContainer)
		return nil, "", err

	}

	return postgresContainer, dsn, nil
}

func connectDb(ctx context.Context, dsn string) (*pgxpool.Pool, error) {

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createOrders(count int) []*order.Order {

	orders := make([]*order.Order, count)
	for i := range count {
		var volume int
		if i%2 == 0 {
			volume = rand.Intn(10) + 1
		} else {
			volume = rand.Intn(20) + 10
		}
		orders[i], _ = order.NewOrder(uuid.New(), kernel.NewRandomLocation(), volume)
	}

	sortById(orders)

	return orders
}

func createCouriers(count int) []*courier.Courier {

	couriers := make([]*courier.Courier, count)
	for i := range count {
		couriers[i], _ = courier.NewCourier(fmt.Sprintf("courier%d", i), rand.Intn(5)+1, kernel.NewRandomLocation())
		if rand.Intn(100) > 50 {
			_ = couriers[i].AddStoragePlace("trunk", 200)
		}
	}

	sortById(couriers)

	return couriers
}

func getRandomUnassignedOrder(orders []*order.Order) *order.Order {
	o := orders[rand.Intn(len(orders))]
	for o.Status() != order.StatusCreated {
		o = orders[rand.Intn(len(orders))]
	}

	return o
}

func equalUUIDs(u1 *uuid.UUID, u2 *uuid.UUID) bool {
	if u1 == nil && u2 == nil {
		return true
	}

	if u1 == nil || u2 == nil {
		return false
	}

	return *u1 == *u2
}

func equalOrders(dto orderDTO, o order.Order) bool {
	return dto.Id == o.Id() &&
		dto.LocationX == o.Location().X() &&
		dto.LocationY == o.Location().Y() &&
		dto.Volume == o.Volume() &&
		dto.Status == o.Status() &&
		equalUUIDs(dto.CourierId, o.CourierId())
}

func equalCouriers(dto courierDTO, c courier.Courier) bool {
	return dto.Id == c.Id() &&
		dto.LocationX == c.Location().X() &&
		dto.LocationY == c.Location().Y() &&
		dto.Name == c.Name() &&
		dto.Speed == c.Speed()
}

func equalStoragePlaces(dto storagePlaceDTO, sp courier.StoragePlace) bool {
	return dto.Id == sp.Id() &&
		dto.Name == sp.Name() &&
		dto.Volume == sp.TotalVolume() &&
		equalUUIDs(dto.OrderId, sp.OrderID())
}

type equalFunc[DTO, AggType any] func(DTO, AggType) bool

func compareRowsAndAggregates[DTO, AggType any](db *pgxpool.Pool, ctx context.Context, query string,
	aggs []*AggType, fn equalFunc[DTO, AggType]) error {

	rows, err := db.Query(ctx, query)
	if err != nil {
		return err
	}

	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[DTO])
	if err != nil {
		return err
	}

	if len(aggs) != len(dtos) {
		return errors.New("len mismatch")
	}

	for i, agg := range aggs {
		if !fn(dtos[i], *agg) {
			return errors.New("data mismatch")
		}
	}

	return nil
}

type hasId interface {
	Id() uuid.UUID
}

func sortById[S ~[]E, E hasId](sl S) {
	slices.SortFunc(sl, func(a E, b E) int {
		id1 := a.Id()
		id2 := b.Id()
		return bytes.Compare(id1[:], id2[:])
	})
}

func find[S []E, E any](sl S, predicate func(E) bool) (E, bool) {
	for _, element := range sl {
		if predicate(element) {
			return element, true
		}
	}

	var defValue E
	return defValue, false
}
