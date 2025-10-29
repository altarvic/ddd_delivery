package kernel

import (
	"delivery/internal/pkg/errs"
	"errors"
	"math/rand"
)

const (
	minC = 1
	maxC = 10
)

var ErrLocationIsEmpty = errors.New("location is empty")

type Location struct {
	x int
	y int

	isSet bool
}

func NewLocation(x int, y int) (Location, error) {
	if x < minC || y < minC || x > maxC || y > maxC {
		return Location{}, errs.ErrValueIsOutOfRange
	}

	return Location{x: x, y: y, isSet: true}, nil
}

func NewRandomLocation() Location {
	return Location{
		x:     rand.Intn(maxC-minC+1) + minC,
		y:     rand.Intn(maxC-minC+1) + minC,
		isSet: true,
	}
}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) IsEmpty() bool {
	return !l.isSet
}

func (l Location) Equals(other Location) bool {
	return l == other
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (l Location) DistanceTo(target Location) (int, error) {
	if l.IsEmpty() || target.IsEmpty() {
		return 0, ErrLocationIsEmpty
	}

	return abs(l.x-target.x) + abs(l.y-target.y), nil
}

// RestoreLocation should be used ONLY inside Repository
func RestoreLocation(x int, y int) Location {
	return Location{
		x:     x,
		y:     y,
		isSet: true,
	}
}
