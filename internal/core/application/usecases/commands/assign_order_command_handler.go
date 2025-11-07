package commands

import (
	"context"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AssignOrderCommandHandler interface {
	Handle(context.Context) error
}

var _ AssignOrderCommandHandler = &assignOrderCommandHandler{}

type assignOrderCommandHandler struct {
	uow ports.UnitOfWork
	d   services.OrderDispatcher
}

func NewAssignOrderCommandHandler(uow ports.UnitOfWork, d services.OrderDispatcher) (AssignOrderCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	if d == nil {
		return nil, errs.NewValueIsRequiredError("d")
	}

	return &assignOrderCommandHandler{
		uow: uow,
		d:   d,
	}, nil
}

func (c *assignOrderCommandHandler) Handle(ctx context.Context) error {

	return c.uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		couriers, err := uowc.CourierRepository().GetAllFree(ctx)
		if err != nil {
			return err
		}

		if couriers == nil || len(couriers) == 0 {
			return nil // No free courier available - no error here
		}

		ord, err := uowc.OrderRepository().GetFirstInCreatedStatus(ctx)
		if err != nil {
			return err
		}

		if ord == nil {
			return nil // No new orders available - no error here
		}

		cour, err := c.d.Dispatch(ord, couriers)
		if err != nil {
			return err
		}

		err = uowc.CourierRepository().Save(ctx, cour)
		if err != nil {
			return err
		}

		return uowc.OrderRepository().Save(ctx, ord)
	})
}
