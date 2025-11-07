package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type MoveCouriersCommandHandler interface {
	Handle(context.Context) error
}

var _ MoveCouriersCommandHandler = &moveCouriersCommandHandler{}

type moveCouriersCommandHandler struct {
	uow ports.UnitOfWork
}

func NewMoveCouriersCommandHandler(uow ports.UnitOfWork) (MoveCouriersCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &moveCouriersCommandHandler{
		uow: uow,
	}, nil
}

func (c *moveCouriersCommandHandler) Handle(ctx context.Context) error {

	return c.uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		assignedOrders, err := uowc.OrderRepository().GetAllInAssignedStatus(ctx)
		if err != nil {
			return err
		}

		if assignedOrders == nil || len(assignedOrders) == 0 {
			return nil // no assigned orders - no error here
		}

		for _, assignedOrder := range assignedOrders {

			cour, err := uowc.CourierRepository().Get(ctx, *assignedOrder.CourierId())
			if err != nil {
				return err
			}

			err = cour.Move(assignedOrder.Location())
			if err != nil {
				return err
			}

			if assignedOrder.Location().Equals(cour.Location()) {
				err = cour.CompleteOrder(assignedOrder)
				if err != nil {
					return err
				}

				err := assignedOrder.Complete()
				if err != nil {
					return err
				}
			}

			err = uowc.CourierRepository().Save(ctx, cour)
			if err != nil {
				return err
			}

			err = uowc.OrderRepository().Save(ctx, assignedOrder)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
