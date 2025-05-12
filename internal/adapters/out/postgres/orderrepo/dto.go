package orderrepo

import (
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
)

type OrderDTO struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey"`
	CourierID *uuid.UUID  `gorm:"type:uuid;index"`
	Location  LocationDTO `gorm:"embedded;embeddedPrefix:location_"`
	Volume    int
	Status    order.Status `gorm:"type:varchar(20)"`
}

type LocationDTO struct {
	X uint8
	Y uint8
}

func (OrderDTO) TableName() string {
	return "orders"
}
