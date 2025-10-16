package order

import "errors"

type Status string

const (
	StatusCreated   Status = "created"
	StatusAssigned  Status = "assigned"
	StatusCompleted Status = "completed"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	switch s {
	case StatusCreated, StatusAssigned, StatusCompleted:
		return true
	default:
		return false
	}
}

func StatusFromString(s string) (Status, error) {
	status := Status(s)
	if status.IsValid() {
		return status, nil
	}

	return status, errors.New("invalid status")
}
