package commands

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
)

type CreateOrderCommandHandler interface {
	Handle(context.Context, CreateOrderCommand) error
}

var _ CreateOrderCommandHandler = &createOrderCommandHandler{}

type createOrderCommandHandler struct {
	uow ports.UnitOfWork
	geo ports.GeoClient
}

func NewCreateOrderCommandHandler(uow ports.UnitOfWork, geo ports.GeoClient) (CreateOrderCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	if geo == nil {
		return nil, errs.NewValueIsRequiredError("geo")
	}

	return &createOrderCommandHandler{
		uow: uow,
		geo: geo,
	}, nil
}

func (c *createOrderCommandHandler) Handle(ctx context.Context, cmd CreateOrderCommand) error {

	if !cmd.isValid {
		return errs.NewValueIsInvalidError("cmd")
	}

	// Убедимся что заказ c заданным id не существует
	err := c.uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		ord, err := uowc.OrderRepository().Get(ctx, cmd.orderID)
		if err != nil {
			return err
		}

		if ord != nil {
			return errors.New("order already exists")
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Получаем координаты из geo сервиса
	loc, err := c.geo.GetGeolocation(ctx, cmd.street)
	if err != nil {
		return err
	}

	// Сохраним заказ в хранилище
	return c.uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		ord, err := order.NewOrder(cmd.orderID, loc, cmd.volume)
		if err != nil {
			return err
		}

		return uowc.OrderRepository().Save(ctx, ord)
	})
}
