package kernel

import (
	"delivery/internal/pkg/errs"
	"math"
	"math/rand"
)

// Location - Координата на доске, она состоит из X (горизонталь) и Y (вертикаль)
type Location struct {
	x uint8
	y uint8
}

func NewLocation(x uint8, y uint8) (Location, error) {
	if x < 1 || x > 10 {
		return Location{}, errs.NewValueIsOutOfRangeError("x", x, 1, 10)
	}
	if y < 1 || y > 10 {
		return Location{}, errs.NewValueIsOutOfRangeError("y", y, 1, 10)
	}

	return Location{x, y}, nil
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

func (l Location) DistanceTo(target Location) uint8 {
	dx := math.Abs(float64(int(l.x) - int(target.x)))
	dy := math.Abs(float64(int(l.y) - int(target.y)))
	return uint8(dx + dy)
}

func CreateRandom() Location {
	return Location{
		x: uint8(rand.Intn(10) + 1),
		y: uint8(rand.Intn(10) + 1),
	}
}
