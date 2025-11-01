package cmd

import (
	"context"
	"delivery/internal/adapters/out/postgres"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type CompositionRoot struct {
	cfg Config
	uow ports.UnitOfWork

	closers []Closer
}

func NewCompositionRoot(cfg Config) *CompositionRoot {

	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.DbSslMode)

	db, err := pgxpool.New(context.Background(), dbConnStr)
	if err != nil {
		log.Fatal("Failed to create a database connection pool")
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	uow, err := postgres.NewUnitOfWork(db)
	if err != nil {
		log.Fatal("Failed to create Unit-of-Work object")
	}

	return &CompositionRoot{
		cfg: cfg,
		uow: uow,
	}
}

func (cr *CompositionRoot) NewOrderDispatcher() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}
