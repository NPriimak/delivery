package kernel

import (
	"delivery/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_givenInvalidParams_whenCreateLocation_thenReturnError(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		x        uint8
		y        uint8
		expected error
	}{
		"x_is_zero":  {0, 9, errs.NewValueIsOutOfRangeError("x", 0, 1, 10)},
		"y_is_zero":  {9, 0, errs.NewValueIsOutOfRangeError("y", 0, 1, 10)},
		"x_too_much": {11, 9, errs.NewValueIsOutOfRangeError("x", 11, 1, 10)},
		"y_too_much": {9, 11, errs.NewValueIsOutOfRangeError("y", 11, 1, 10)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewLocation(test.x, test.y)
			assert.Errorf(t, err, test.expected.Error())
		})
	}
}

func Test_givenValidParams_whenCreateLocation_thenSuccess(t *testing.T) {
	var x uint8 = 5
	var y uint8 = 6

	res, err := NewLocation(x, y)

	assert.NoError(t, err)
	assert.Equal(t, x, res.X())
	assert.Equal(t, y, res.Y())
}

func Test_givenTwoEqualsLocation_whenCompareThem_thenReturnTrue(t *testing.T) {
	first, _ := NewLocation(4, 5)
	second, _ := NewLocation(4, 5)

	assert.True(t, first.Equals(second))
}

func Test_givenTwoNotEqualsLocation_whenCompareThem_thenReturnFalse(t *testing.T) {
	first, _ := NewLocation(4, 5)
	second, _ := NewLocation(4, 6)

	assert.False(t, first.Equals(second))
}

func Test_givenTwoValidLocations_thenCountDistanceTo_thenReturnCorrectDistance(t *testing.T) {
	first, _ := NewLocation(2, 6)
	second, _ := NewLocation(4, 9)
	third, _ := NewLocation(4, 9)
	thour, _ := NewLocation(10, 5)

	tests := map[string]struct {
		first    Location
		second   Location
		expected uint8
	}{
		"less_to_high":  {first, second, 5},
		"high_to_less":  {second, first, 5},
		"equals":        {second, third, 0},
		"one_more_time": {third, thour, 10},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(
				t,
				test.expected,
				test.first.DistanceTo(test.second),
			)
		})
	}
}

func Test_whenCreateRandomLocation_thenSuccess(t *testing.T) {
	assert.NotPanics(t, func() {
		CreateRandom()
	}, "expected not panic when creating random location")
}
