package commands

import (
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type CreateOrderCommand struct {
	orderID uuid.UUID
	street  string
	volume  int
	isValid bool
}

func NewCreateOrderCommand(orderID uuid.UUID, street string, volume int) (CreateOrderCommand, error) {

	if orderID == uuid.Nil {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("orderID")
	}

	if strings.TrimSpace(street) == "" {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("street")
	}

	if volume <= 0 {
		return CreateOrderCommand{}, errors.New("volume must be greater than 0")
	}

	return CreateOrderCommand{
		orderID: orderID,
		street:  street,
		volume:  volume,
		isValid: true}, nil
}

func (c CreateOrderCommand) OrderID() uuid.UUID {
	return c.orderID
}

func (c CreateOrderCommand) Street() string {
	return c.street
}

func (c CreateOrderCommand) Volume() int {
	return c.volume
}

func (c CreateOrderCommand) IsValid() bool {
	return c.isValid
}
