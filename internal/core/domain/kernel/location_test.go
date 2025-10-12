package kernel

import (
	"fmt"
	"testing"
)

func TestNewLocation(t *testing.T) {

	tests := []struct {
		x             int
		y             int
		errIsExpected bool
	}{
		{
			x:             0,
			y:             0,
			errIsExpected: true,
		},
		{
			x:             1,
			y:             -1,
			errIsExpected: true,
		},
		{
			x:             11,
			y:             3,
			errIsExpected: true,
		},
		{
			x:             4,
			y:             3,
			errIsExpected: false,
		},
		{
			x:             -1,
			y:             5,
			errIsExpected: true,
		},
		{
			x:             10,
			y:             10,
			errIsExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d,%d", test.x, test.y), func(t *testing.T) {
			_, err := NewLocation(test.x, test.y)
			if test.errIsExpected {
				if err == nil {
					t.Fail()
				}
			} else {
				if err != nil {
					t.Fail()
				}
			}
		})
	}
}

func TestLocation_IsEmpty(t *testing.T) {
	l := Location{}
	if !l.IsEmpty() {
		t.Fail()
	}

	l, _ = NewLocation(1, 1)
	if l.IsEmpty() {
		t.Fail()
	}
}

func TestLocation_Equals(t *testing.T) {
	l1, _ := NewLocation(2, 3)
	l2, _ := NewLocation(2, 3)
	if !l1.Equals(l2) {
		t.Fail()
	}
}

func TestLocation_DistanceTo(t *testing.T) {

	tests := []struct {
		x1               int
		y1               int
		x2               int
		y2               int
		expectedDistance int
	}{
		{
			x1:               1,
			y1:               1,
			x2:               5,
			y2:               9,
			expectedDistance: 12,
		},
		{
			x1:               10,
			y1:               1,
			x2:               5,
			y2:               5,
			expectedDistance: 9,
		},
		{
			x1:               2,
			y1:               3,
			x2:               3,
			y2:               2,
			expectedDistance: 2,
		},
		{
			x1:               2,
			y1:               3,
			x2:               2,
			y2:               3,
			expectedDistance: 0,
		},
		{
			x1:               2,
			y1:               3,
			x2:               3,
			y2:               3,
			expectedDistance: 1,
		},
		{
			x1:               10,
			y1:               10,
			x2:               1,
			y2:               1,
			expectedDistance: 18,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("from %d,%d to %d,%d", test.x1, test.y1, test.x2, test.y2), func(t *testing.T) {
			l1, _ := NewLocation(test.x1, test.y1)
			l2, _ := NewLocation(test.x2, test.y2)

			d1, _ := l1.DistanceTo(l2)
			if d1 != test.expectedDistance {
				t.Fail()
			}

			d2, _ := l2.DistanceTo(l1)
			if d1 != d2 {
				t.Fail()
			}

		})
	}

}

func TestNewRandomLocation(t *testing.T) {
	for range 1000 {
		l := NewRandomLocation()
		if l.X() < minC || l.Y() < minC || l.X() > maxC || l.Y() > maxC {
			t.Fail()
		}
	}
}
