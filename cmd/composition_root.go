package cmd

import (
	"context"
	kafkain "delivery/internal/adapters/in/kafka"
	"delivery/internal/adapters/out/grpc/geo"
	kafkaout "delivery/internal/adapters/out/kafka"
	outb "delivery/internal/adapters/out/outbox"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/application/eventhandlers"
	"delivery/internal/core/application/usecases/commands"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/jobs"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/outbox"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"log"
	"reflect"
	"sync"
)

type CompositionRoot struct {
	cfg           Config
	db            *pgxpool.Pool
	uow           ports.UnitOfWork
	mediatr       ddd.Mediatr
	eventRegistry outbox.EventRegistry

	closers []Closer
}

func NewCompositionRoot(cfg Config) *CompositionRoot {

	db := dbConnect(cfg)
	mediatr := ddd.NewMediatr()
	uow := createUnitOfWork(db, mediatr)
	eventRegistry := createEventRegistry()

	return &CompositionRoot{
		cfg:           cfg,
		db:            db,
		uow:           uow,
		mediatr:       mediatr,
		eventRegistry: eventRegistry,
	}
}

func dbConnect(cfg Config) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbSslMode)

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to create a database connection pool: %v", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return db
}

func createUnitOfWork(db *pgxpool.Pool, mediatr ddd.Mediatr) ports.UnitOfWork {
	uow, err := postgres.NewUnitOfWork(db, mediatr)
	if err != nil {
		log.Fatalf("Failed to create UnitOfWork: %v", err)
	}

	return uow
}

func createEventRegistry() outbox.EventRegistry {
	registry, err := outbox.NewEventRegistry()
	if err != nil {
		log.Fatalf("cannot create EventRegistry: %v", err)
	}

	err = registry.RegisterDomainEvent(reflect.TypeOf(order.CreatedDomainEvent{}))
	err = registry.RegisterDomainEvent(reflect.TypeOf(order.CompletedDomainEvent{}))

	if err != nil {
		log.Fatalf("cannot register domain event: %v", err)
	}

	return registry
}

func (cr *CompositionRoot) Db() *pgxpool.Pool {
	return cr.db
}

func (cr *CompositionRoot) UnitOfWork() ports.UnitOfWork {
	return cr.uow
}

func (cr *CompositionRoot) Mediatr() ddd.Mediatr {
	return cr.mediatr
}

func (cr *CompositionRoot) NewOrderDispatcher() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}

func (cr *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	cmdHandler, err := commands.NewCreateCourierCommandHandler(cr.uow)
	if err != nil {
		log.Fatalf("Failed to create CreateCourierCommandHandler: %v", err)
	}

	return cmdHandler
}

func (cr *CompositionRoot) NewAddStoragePlaceCommandHandler() commands.AddStoragePlaceCommandHandler {
	cmdHandler, err := commands.NewAddStoragePlaceCommandHandler(cr.uow)
	if err != nil {
		log.Fatalf("Failed to create AddStoragePlaceCommandHandler: %v", err)
	}

	return cmdHandler
}

func (cr *CompositionRoot) NewAssignOrderCommandHandler() commands.AssignOrderCommandHandler {
	cmdHandler, err := commands.NewAssignOrderCommandHandler(cr.uow, cr.NewOrderDispatcher())
	if err != nil {
		log.Fatalf("Failed to create AssignOrderCommandHandler: %v", err)
	}

	return cmdHandler
}

func (cr *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	cmdHandler, err := commands.NewCreateOrderCommandHandler(cr.uow, cr.NewGeoLocationService())
	if err != nil {
		log.Fatalf("Failed to create CreateOrderCommandHandler: %v", err)
	}

	return cmdHandler
}

func (cr *CompositionRoot) NewMoveCouriersCommandHandler() commands.MoveCouriersCommandHandler {
	cmdHandler, err := commands.NewMoveCouriersCommandHandler(cr.uow)
	if err != nil {
		log.Fatalf("Failed to create MoveCouriersCommandHandler: %v", err)
	}

	return cmdHandler
}

func (cr *CompositionRoot) NewAllCouriersQueryHandler() queries.AllCouriersQueryHandler {
	cmdHandler, err := queries.NewAllCouriersQueryHandler(cr.db)
	if err != nil {
		log.Fatalf("Failed to create AllCouriersQueryHandler: %v", err)
	}

	return cmdHandler
}

func (cr *CompositionRoot) NewIncompleteOrdersQueryHandler() queries.IncompleteOrdersQueryHandler {
	cmdHandler, err := queries.NewIncompleteOrdersQueryHandler(cr.db)
	if err != nil {
		log.Fatalf("Failed to create IncompleteOrdersQueryHandler: %v", err)
	}

	return cmdHandler
}

func (cr *CompositionRoot) NewAssignOrdersJob() cron.Job {
	job, err := jobs.NewAssignOrdersJob(cr.NewAssignOrderCommandHandler())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersJob: %v", err)
	}

	return job
}

func (cr *CompositionRoot) NewMoveCouriersJob() cron.Job {
	job, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewGeoLocationService() ports.GeoClient {

	return sync.OnceValue(func() ports.GeoClient {
		client, err := geo.NewClient(cr.cfg.GeoServiceGrpcHost)
		if err != nil {
			log.Fatalf("cannot create GeoClient: %v", err)
		}

		return client
	})()
}

func (cr *CompositionRoot) NewBasketConfirmedEventsConsumer() kafkain.BasketConfirmedEventsConsumer {
	return sync.OnceValue(func() kafkain.BasketConfirmedEventsConsumer {
		consumer, err := kafkain.NewBasketConfirmedEventsConsumer(
			[]string{cr.cfg.KafkaHost},
			cr.cfg.KafkaConsumerGroup,
			cr.cfg.KafkaBasketConfirmedTopic,
			cr.NewCreateOrderCommandHandler(),
		)
		if err != nil {
			log.Fatalf("failed to create BasketConfirmedEventsConsumer: %v", err)
		}
		cr.RegisterCloser(consumer)
		return consumer
	})()
}

func (cr *CompositionRoot) NewOrderChangedNotificationProducer() ports.NotificationProducer {
	return sync.OnceValue(func() ports.NotificationProducer {
		producer, err := kafkaout.NewOrderChangedNotificationProducer(
			[]string{cr.cfg.KafkaHost},
			cr.cfg.KafkaOrderChangedTopic,
		)

		if err != nil {
			log.Fatalf("failed to create OrderChangedNotificationProducer: %v", err)
		}

		cr.RegisterCloser(producer)
		return producer
	})()
}

func (cr *CompositionRoot) NewOrderEventsHandler(np ports.NotificationProducer) ddd.EventHandler {
	eventHandler, err := eventhandlers.NewOrderDomainEventsHandler(np)
	if err != nil {
		log.Fatalf("failed to create NotificationProducer: %v", err)
	}

	return eventHandler
}

func (cr *CompositionRoot) NewOutboxRepository() outb.OutboxRepository {
	ob, err := outb.NewRepository(cr.db)
	if err != nil {
		log.Fatalf("failed to create OutboxRepository: %v", err)
	}

	return ob
}

func (cr *CompositionRoot) NewOutboxJob() cron.Job {
	job, err := jobs.NewOutboxJob(cr.NewOutboxRepository(), cr.mediatr, cr.eventRegistry)
	if err != nil {
		log.Fatalf("cannot create OutboxJob: %v", err)
	}
	return job
}
