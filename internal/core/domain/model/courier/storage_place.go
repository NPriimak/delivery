package courier

import (
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"strings"
)

const (
	MinVolume = 1
)

var (
	ErrCannotStoreOrderInThisStoragePlace = errors.New("cannot store order in this storage place")
	ErrOrderNotStoredInThisPlace          = errors.New("order is not stored in this place")
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderID     *uuid.UUID
}

func NewStoragePlace(
	name string,
	totalVolume int,
) (*StoragePlace, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}

	if totalVolume < MinVolume {
		return nil, errs.NewValueIsRequiredError("totalVolume")
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}, nil
}

func (s *StoragePlace) CanStore(volume int) (bool, error) {
	if volume < MinVolume {
		return false, errs.NewValueIsRequiredError("volume")
	}

	return volume <= s.totalVolume && !s.isOccupied(), nil
}

func (s *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	if orderID == uuid.Nil {
		return errs.NewValueIsRequiredError("orderID")
	}
	if volume < MinVolume {
		return errs.NewValueIsRequiredError("volume")
	}

	canStore, err := s.CanStore(volume)
	if err != nil {
		return err
	}

	if !canStore {
		return ErrCannotStoreOrderInThisStoragePlace
	}

	s.orderID = &orderID
	return nil
}

func (s *StoragePlace) Clear(orderID uuid.UUID) error {
	if orderID == uuid.Nil {
		return errs.NewValueIsRequiredError("orderID")
	}
	if !s.hasOrder(orderID) {
		return ErrOrderNotStoredInThisPlace
	}

	s.orderID = nil
	return nil
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderID
}

func (s *StoragePlace) isOccupied() bool {
	return s.orderID != nil
}

func (s *StoragePlace) hasOrder(orderID uuid.UUID) bool {
	return s.orderID != nil && *s.orderID == orderID
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
	if other == nil {
		return false
	}
	return s.id == other.id
}

// RestoreStoragePlace restore StoragePlace from db. DO NOT USE IN DOMAIN!
func RestoreStoragePlace(id uuid.UUID, name string, totalVolume int, orderID *uuid.UUID) *StoragePlace {
	return &StoragePlace{
		id:          id,
		name:        name,
		totalVolume: totalVolume,
		orderID:     orderID,
	}
}
