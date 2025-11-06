package commands

import (
	"context"
	"delivery/internal/core/domain/kernel"
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
}

func NewCreateOrderCommandHandler(uow ports.UnitOfWork) (CreateOrderCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &createOrderCommandHandler{
		uow: uow,
	}, nil
}

func (c *createOrderCommandHandler) Handle(ctx context.Context, cmd CreateOrderCommand) error {

	if !cmd.isValid {
		return errs.NewValueIsInvalidError("cmd")
	}

	return c.uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		ord, err := uowc.OrderRepository().Get(ctx, cmd.orderID)
		if err != nil {
			return err
		}

		if ord != nil {
			return errors.New("order already exists")
		}

		// TODO: get location from cmd.street
		loc := kernel.NewRandomLocation()
		ord, err = order.NewOrder(cmd.orderID, loc, cmd.volume)
		if err != nil {
			return err
		}

		return uowc.OrderRepository().Save(ctx, ord)
	})
}
