package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type AddStoragePlaceCommandHandler interface {
	Handle(context.Context, AddStoragePlaceCommand) error
}

var _ AddStoragePlaceCommandHandler = &addStoragePlaceCommandHandler{}

type addStoragePlaceCommandHandler struct {
	uow ports.UnitOfWork
}

func NewAddStoragePlaceCommandHandler(uow ports.UnitOfWork) (AddStoragePlaceCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &addStoragePlaceCommandHandler{
		uow: uow,
	}, nil
}

func (c *addStoragePlaceCommandHandler) Handle(ctx context.Context, cmd AddStoragePlaceCommand) error {

	if !cmd.isValid {
		return errs.NewValueIsInvalidError("cmd")
	}

	return c.uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {

		cour, err := uowc.CourierRepository().Get(ctx, cmd.courierID)
		if err != nil {
			return err
		}

		err = cour.AddStoragePlace(cmd.name, cmd.totalVolume)
		if err != nil {
			return err
		}

		return uowc.CourierRepository().Save(ctx, cour)
	})
}
