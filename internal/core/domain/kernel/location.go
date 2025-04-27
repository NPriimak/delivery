package kernel

import (
	"delivery/internal/pkg/errs"
	"math"
	"math/rand"
)

const (
	minX uint8 = 1
	maxX uint8 = 10
	minY uint8 = 1
	maxY uint8 = 10
)

// Location - Координата на доске, она состоит из X (горизонталь) и Y (вертикаль)
type Location struct {
	x     uint8
	y     uint8
	isSet bool
}

func NewLocation(x uint8, y uint8) (Location, error) {
	if !isXWithinAllowedRange(x) {
		return Location{}, errs.NewValueIsOutOfRangeError("x", x, minX, maxX)
	}
	if !isYWithinAllowedRange(y) {
		return Location{}, errs.NewValueIsOutOfRangeError("y", y, minY, maxY)
	}

	return Location{x, y, true}, nil
}

func isXWithinAllowedRange(x uint8) bool {
	return x >= minX && x <= maxX
}

func isYWithinAllowedRange(y uint8) bool {
	return y >= minY && y <= maxY
}

func (l Location) X() uint8 {
	return l.x
}

func (l Location) Y() uint8 {
	return l.y
}

func (l Location) Equals(x Location) bool {
	return l == x
}

func (l Location) CountDistanceTo(target Location) (uint8, error) {

	if target.IsEmpty() {
		return 0, errs.NewValueIsRequiredError("target")
	}

	dx := math.Abs(float64(int(l.x) - int(target.x)))
	dy := math.Abs(float64(int(l.y) - int(target.y)))
	return uint8(dx + dy), nil
}

func (l Location) IsEmpty() bool {
	return !l.isSet
}

func CreateRandomLocation() Location {
	return Location{
		x:     uint8(rand.Intn(10) + 1),
		y:     uint8(rand.Intn(10) + 1),
		isSet: true,
	}
}
