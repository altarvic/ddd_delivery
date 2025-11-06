package commands

import (
	"context"
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateCourierCommandHandler interface {
	Handle(context.Context, CreateCourierCommand) error
}

var _ CreateCourierCommandHandler = &createCourierCommandHandler{}

type createCourierCommandHandler struct {
	uow ports.UnitOfWork
}

func NewCreateCourierCommandHandler(uow ports.UnitOfWork) (CreateCourierCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &createCourierCommandHandler{
		uow: uow,
	}, nil
}

func (c *createCourierCommandHandler) Handle(ctx context.Context, cmd CreateCourierCommand) error {

	if !cmd.isValid {
		return errs.NewValueIsInvalidError("cmd")
	}

	cour, err := courier.NewCourier(cmd.name, cmd.speed, kernel.NewRandomLocation())
	if err != nil {
		return err
	}

	return c.uow.Do(ctx, func(ctx context.Context, uowc ports.UnitOfWorkComponents) error {
		return uowc.CourierRepository().Save(ctx, cour)
	})
}
