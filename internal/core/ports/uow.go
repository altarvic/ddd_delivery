package ports

import "context"

type UnitOfWorkComponents interface {
	OrderRepository() OrderRepository
	CourierRepository() CourierRepository
}

type UnitOfWorkDoFunc = func(ctx context.Context, uowc UnitOfWorkComponents) error

type UnitOfWork interface {
	Do(ctx context.Context, fn UnitOfWorkDoFunc) error
}
