package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"testing"
)

func newValidLocation() kernel.Location {
	loc, _ := kernel.NewLocation(1, 1)
	return loc
}

func TestNewCourier(t *testing.T) {

	tests := []struct {
		testName    string
		courierName string
		speed       int
		location    kernel.Location
		expectError bool
	}{
		{
			testName:    "invalid courier name",
			courierName: "",
			speed:       1,
			location:    newValidLocation(),
			expectError: true,
		},
		{
			testName:    "invalid speed",
			courierName: "fullstop",
			speed:       0,
			location:    newValidLocation(),
			expectError: true,
		},
		{
			testName:    "invalid location",
			courierName: "lost",
			speed:       1,
			location:    kernel.Location{},
			expectError: true,
		},
		{
			testName:    "valid courier fields",
			courierName: "speedy",
			speed:       6,
			location:    newValidLocation(),
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			c, err := NewCourier(test.courierName, test.speed, test.location)
			if test.expectError {
				if err == nil {
					t.Fail()
				}
			} else {
				if err != nil {
					t.Error(err)
				}

				if c.Id() == uuid.Nil {
					t.Error("id is nil")
				}

				if c.Name() != test.courierName {
					t.Error("courier.Name")
				}

				if c.Speed() != test.speed {
					t.Error("courier.Speed")
				}

				if !c.Location().Equals(test.location) {
					t.Error("courier.Location")
				}

				if c.StoragePlaces() == nil || len(c.StoragePlaces()) != 1 {
					t.Error("courier.StoragePlaces")
				}

				if c.StoragePlaces()[0].Name() != "Bag" {
					t.Error("no Bag in courier.StoragePlaces")
				}

				if c.StoragePlaces()[0].TotalVolume() != 10 {
					t.Error("invalid Bag's volume in courier.StoragePlaces")
				}
			}
		})
	}
}

func TestCourier_AddStoragePlace(t *testing.T) {
	c, _ := NewCourier("vzuh", 8, newValidLocation())

	err := c.AddStoragePlace("", 500)
	if err == nil {
		t.Error("invalid storage name")
	}

	err = c.AddStoragePlace("trunk", 0)
	if err == nil {
		t.Error("invalid storage volume")
	}

	err = c.AddStoragePlace("trunk", 500)
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, place := range c.StoragePlaces() {
		if place.Name() == "trunk" && place.TotalVolume() == 500 {
			found = true
			break
		}
	}

	if !found {
		t.Error("storage place not found")
	}
}

func TestCourier_CanTakeOrder(t *testing.T) {
	c, _ := NewCourier("Slow", 2, newValidLocation())

	o, _ := order.NewOrder(uuid.New(), newValidLocation(), 200)
	if ok, _ := c.CanTakeOrder(o); ok {
		t.Error("can't take volume 200")
	}

	o, _ = order.NewOrder(uuid.New(), newValidLocation(), 10)
	if ok, _ := c.CanTakeOrder(o); !ok {
		t.Error("can take volume 10")
	}

	_ = c.AddStoragePlace("trunk", 500)
	o, _ = order.NewOrder(uuid.New(), newValidLocation(), 500)
	if ok, _ := c.CanTakeOrder(o); !ok {
		t.Error("can take volume 500")
	}

	o, _ = order.NewOrder(uuid.New(), newValidLocation(), 501)
	if ok, _ := c.CanTakeOrder(o); ok {
		t.Error("can't take volume 501")
	}
}

func TestCourier_TakeOrder(t *testing.T) {
	c, _ := NewCourier("Slow", 2, newValidLocation())

	o, _ := order.NewOrder(uuid.New(), newValidLocation(), 100)
	if err := c.TakeOrder(o); err == nil {
		t.Error("can't take volume 100")
	}

	o, _ = order.NewOrder(uuid.New(), newValidLocation(), 5)
	if err := c.TakeOrder(o); err != nil {
		t.Error("can take volume 5")
	}

	if o.Status() != order.StatusAssigned {
		t.Error("order status should be assigned")
	}

	if o.CourierId() == nil || *o.CourierId() != c.Id() {
		t.Error("order's courier id should be equal to courier's id")
	}
}

func TestCourier_CompleteOrder(t *testing.T) {
	c, _ := NewCourier("Slow", 2, newValidLocation())
	o, _ := order.NewOrder(uuid.New(), newValidLocation(), 8)

	if err := c.CompleteOrder(o); err == nil {
		t.Error("must not complete non-owned order")
	}

	_ = c.TakeOrder(o)

	if err := c.CompleteOrder(o); err != nil {
		t.Error(err)
	}

	if o.Status() != order.StatusCompleted {
		t.Error("order status must be completed")
	}

	if c.StoragePlaces()[0].IsOccupied() {
		t.Error("storage place must be empty")
	}
}

func TestCourier_CalculateTimeToLocation(t *testing.T) {

	tests := []struct {
		courierName      string
		courierLocationX int
		courierLocationY int
		courierSpeed     int
		targetLocationX  int
		targetLocationY  int
		expectedTime     float64
	}{
		{
			courierName:      "by foot",
			courierLocationX: 2,
			courierLocationY: 3,
			courierSpeed:     1,
			targetLocationX:  9,
			targetLocationY:  4,
			expectedTime:     8,
		},
		{
			courierName:      "bike",
			courierLocationX: 1,
			courierLocationY: 1,
			courierSpeed:     2,
			targetLocationX:  5,
			targetLocationY:  5,
			expectedTime:     4,
		},
		{
			courierName:      "moto",
			courierLocationX: 3,
			courierLocationY: 3,
			courierSpeed:     5,
			targetLocationX:  7,
			targetLocationY:  7,
			expectedTime:     1.6,
		},
		{
			courierName:      "car",
			courierLocationX: 3,
			courierLocationY: 3,
			courierSpeed:     9,
			targetLocationX:  9,
			targetLocationY:  9,
			expectedTime:     1.333,
		},
	}

	for _, test := range tests {
		t.Run(test.courierName, func(t *testing.T) {
			loc, _ := kernel.NewLocation(test.courierLocationX, test.courierLocationY)
			target, _ := kernel.NewLocation(test.targetLocationX, test.targetLocationY)
			c, _ := NewCourier(test.courierName, test.courierSpeed, loc)

			time, err := c.CalculateTimeToLocation(target)
			if err != nil {
				t.Error(err)
			}

			if time != test.expectedTime {
				t.Errorf("time: %f, expected: %f", time, test.expectedTime)
			}
		})
	}

}

func TestCourier_Move(t *testing.T) {

outerLoop:
	for range 10000 {

		loc := kernel.NewRandomLocation()
		speed := rand.Intn(9) + 1
		c, _ := NewCourier("test courier", speed, loc)

		target := kernel.NewRandomLocation()

		time, _ := c.CalculateTimeToLocation(target)
		steps := int(math.Ceil(time))

		for range steps {
			err := c.Move(target)
			if err != nil {
				t.Error(err)
				break outerLoop
			}
		}

		if !c.location.Equals(target) {
			t.Errorf("location: (%d, %d) target: (%d, %d)", c.Location().X(), c.Location().Y(), target.X(), target.Y())
			break
		}
	}
}
