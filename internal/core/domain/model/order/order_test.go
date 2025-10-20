package order

import (
	"delivery/internal/core/domain/kernel"
	"github.com/google/uuid"
	"testing"
)

func newValidLocation() kernel.Location {
	loc, _ := kernel.NewLocation(1, 1)
	return loc
}

func TestNewOrder(t *testing.T) {

	loc := newValidLocation()

	tests := []struct {
		name        string
		orderId     uuid.UUID
		location    kernel.Location
		volume      int
		expectError bool
	}{
		{
			name:        "invalid id",
			orderId:     uuid.UUID{},
			volume:      10,
			location:    loc,
			expectError: true,
		},
		{
			name:        "invalid location",
			orderId:     uuid.New(),
			volume:      10,
			location:    kernel.Location{},
			expectError: true,
		},
		{
			name:        "invalid volume",
			orderId:     uuid.New(),
			volume:      0,
			location:    loc,
			expectError: true,
		},
		{
			name:        "valid order fields",
			orderId:     uuid.New(),
			volume:      10,
			location:    loc,
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o, err := NewOrder(test.orderId, test.location, test.volume)
			if test.expectError {
				if err == nil {
					t.Fail()
				}
			} else {
				if err != nil {
					t.Fail()
				}

				if o.Status() != StatusCreated {
					t.Fail()
				}

				if !o.Location().Equals(test.location) {
					t.Error("location")
				}

				if o.Volume() != test.volume {
					t.Error("volume")
				}

				if o.Id() != test.orderId {
					t.Error("orderId")
				}

				if o.CourierId() != nil {
					t.Error("courierId")
				}

				if !o.Location().Equals(test.location) {
					t.Error("location")
				}
			}
		})
	}
}

func TestOrder_Equals(t *testing.T) {
	loc := newValidLocation()
	volume := 10
	orderId := uuid.New()

	o1, _ := NewOrder(orderId, loc, volume)
	o2, _ := NewOrder(orderId, loc, volume)

	if !o1.Equals(o2) {
		t.Error("must be equal")
	}

	o2, _ = NewOrder(uuid.New(), loc, volume)

	if o1.Equals(o2) {
		t.Error("must not be equal")
	}

}

func TestOrder_AssignCourier(t *testing.T) {
	o, _ := NewOrder(uuid.New(), newValidLocation(), 10)

	err := o.AssignCourier(uuid.UUID{})
	if err == nil {
		t.Error("invalid courier Id")
	}

	courierId := uuid.New()
	err = o.AssignCourier(courierId)
	if err != nil {
		t.Error(err)
	}

	if o.Status() != StatusAssigned {
		t.Error("status != assigned")
	}

	if o.CourierId() == nil || *o.CourierId() != courierId {
		t.Error("invalid courierId")
	}

	err = o.AssignCourier(courierId)
	if err == nil {
		t.Error("courier already assigned")
	}
}

func TestOrder_Complete(t *testing.T) {
	o, _ := NewOrder(uuid.New(), newValidLocation(), 10)

	err := o.Complete()
	if err == nil {
		t.Error("no courier")
	}

	_ = o.AssignCourier(uuid.New())

	err = o.Complete()
	if err != nil {
		t.Error(err)
	}

	if o.Status() != StatusCompleted {
		t.Error("status != completed")
	}

	err = o.Complete()
	if err == nil {
		t.Error("already completed")
	}

}
