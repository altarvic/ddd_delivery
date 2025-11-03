package commands

import (
	"delivery/internal/pkg/errs"
	"errors"
	"strings"
)

type CreateCourierCommand struct {
	name    string
	speed   int
	isValid bool
}

func NewCreateCourierCommand(name string, speed int) (CreateCourierCommand, error) {

	if strings.TrimSpace(name) == "" {
		return CreateCourierCommand{}, errs.NewValueIsRequiredError("name")
	}

	if speed <= 0 {
		return CreateCourierCommand{}, errors.New("speed must be greater than 0")
	}

	return CreateCourierCommand{
		name:    name,
		speed:   speed,
		isValid: true,
	}, nil
}

func (c CreateCourierCommand) Name() string {
	return c.name
}

func (c CreateCourierCommand) Speed() int {
	return c.speed
}

func (c CreateCourierCommand) IsValid() bool {
	return c.isValid
}
