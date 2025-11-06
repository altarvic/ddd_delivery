package commands

import (
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"strings"
)

type AddStoragePlaceCommand struct {
	courierID   uuid.UUID
	name        string
	totalVolume int
	isValid     bool
}

func NewAddStoragePlaceCommand(courierID uuid.UUID, name string, totalVolume int) (AddStoragePlaceCommand, error) {

	if courierID == uuid.Nil {
		return AddStoragePlaceCommand{}, errs.NewValueIsRequiredError("courierID")
	}

	if strings.TrimSpace(name) == "" {
		return AddStoragePlaceCommand{}, errs.NewValueIsRequiredError("name")
	}

	if totalVolume <= 0 {
		return AddStoragePlaceCommand{}, errors.New("totalVolume must be greater than 0")
	}

	return AddStoragePlaceCommand{
		courierID:   courierID,
		name:        name,
		totalVolume: totalVolume,
		isValid:     true,
	}, nil
}

func (a AddStoragePlaceCommand) CourierID() uuid.UUID {
	return a.courierID
}

func (a AddStoragePlaceCommand) Name() string {
	return a.name
}

func (a AddStoragePlaceCommand) TotalVolume() int {
	return a.totalVolume
}

func (a AddStoragePlaceCommand) IsValid() bool {
	return a.isValid
}
