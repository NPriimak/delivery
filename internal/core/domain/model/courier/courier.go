package courier

import (
	"delivery/internal/core/domain/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"math"
	"strings"
)

const (
	MinSpeed = 1
)

var (
	ErrNoStoragePlace       = errors.New("no storage place")
	ErrOrderStorageNotFound = errors.New("order storage not found")
)

type Courier struct {
	id            uuid.UUID
	name          string
	speed         int
	location      kernel.Location
	storagePlaces []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if speed < MinSpeed {
		return nil, errs.NewValueIsRequiredError("speed")
	}

	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}

	storagePlaces := make([]*StoragePlace, 0)
	return &Courier{
		id:            uuid.New(),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: storagePlaces,
	}, nil
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	if strings.TrimSpace(name) == "" {
		return errs.NewValueIsRequiredError("name")
	}
	if volume < MinVolume {
		return errs.NewValueIsRequiredError("volume")
	}

	sp, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}

	c.storagePlaces = append(c.storagePlaces, sp)
	return nil
}

func (c *Courier) CanTakeOrder(order *order.Order) (bool, error) {
	if order == nil {
		return false, errs.NewValueIsRequiredError("order")
	}

	if c.storagePlaces == nil {
		return false, nil
	}

	freeStorage, err := c.findFirstFreeStorage(order.Volume())
	if err != nil {
		return false, err
	}

	if freeStorage == nil {
		return false, nil
	}

	return true, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}

	canTake, err := c.CanTakeOrder(order)
	if err != nil {
		return err
	}

	if !canTake {
		return ErrNoStoragePlace
	}

	freeStorage, err := c.findFirstFreeStorage(order.Volume())
	if err != nil {
		return err
	}

	if freeStorage == nil {
		return ErrNoStoragePlace
	}

	if err = freeStorage.Store(order.ID(), order.Volume()); err != nil {
		return err
	}
	if err = order.Assign(c.id); err != nil {
		return err
	}

	return nil
}

func (c *Courier) completeOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}

	storage, err := c.findOrderStorage(order.ID())
	if err != nil {
		return err
	}

	if storage == nil {
		return ErrOrderStorageNotFound
	}

	if err = storage.Clear(order.ID()); err != nil {
		return err
	}
	if err = order.Complete(); err != nil {
		return err
	}

	return nil
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (float64, error) {
	if target.IsEmpty() {
		return 0, errs.NewValueIsRequiredError("target")
	}
	distance, err := c.location.CountDistanceTo(target)
	if err != nil {
		return 0, err
	}

	time := float64(distance) / float64(c.speed)
	return time, err
}

func (c *Courier) Move(target kernel.Location) error {
	if target.IsEmpty() {
		return errs.NewValueIsRequiredError("target")
	}

	dx := float64(target.X()) - float64(c.location.X())
	dy := float64(target.Y()) - float64(c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := float64(c.location.X()) + dx
	newY := float64(c.location.Y()) + dy

	newLocation, err := kernel.NewLocation(uint8(newX), uint8(newY))
	if err != nil {
		return err
	}
	c.location = newLocation
	return nil
}

func (c *Courier) findOrderStorage(orderID uuid.UUID) (*StoragePlace, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("order")
	}

	for _, storagePlace := range c.storagePlaces {
		if storagePlace.isOccupied() && *storagePlace.OrderID() == orderID {
			return storagePlace, nil
		}
	}

	return nil, nil
}

func (c *Courier) findFirstFreeStorage(volume int) (*StoragePlace, error) {
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}
	for _, sp := range c.storagePlaces {
		canStore, err := sp.CanStore(volume)
		if err != nil {
			return nil, err
		}
		if canStore {
			return sp, nil
		}
	}
	return nil, nil
}

func (c *Courier) ID() uuid.UUID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) StoragePlaces() []*StoragePlace {
	return c.storagePlaces
}

func (c *Courier) Equals(other *Courier) bool {
	if other == nil {
		return false
	}

	return c.id == other.id
}
