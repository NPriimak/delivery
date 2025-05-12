package queries

import (
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAllCouriersQuery struct {
	isSet bool
}

func NewGetAllCouriersQuery() (GetAllCouriersQuery, error) {
	return GetAllCouriersQuery{
		isSet: true,
	}, nil
}

func (q GetAllCouriersQuery) IsEmpty() bool {
	return !q.isSet
}

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string
	Location LocationResponse `gorm:"embedded;embeddedPrefix:location_"`
}

func (CourierResponse) TableName() string {
	return "couriers"
}

type GetAllCouriersQueryHandler interface {
	Handle(GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

type getAllCouriersQueryHandler struct {
	db *gorm.DB
}

func NewGetAllCouriersQueryHandler(db *gorm.DB) (GetAllCouriersQueryHandler, error) {
	if db == nil {
		return &getAllCouriersQueryHandler{}, errs.NewValueIsRequiredError("db")
	}
	return &getAllCouriersQueryHandler{db: db}, nil
}

func (q *getAllCouriersQueryHandler) Handle(query GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	if query.IsEmpty() {
		return GetAllCouriersResponse{}, errs.NewValueIsRequiredError("query")
	}

	var couriers []CourierResponse
	result := q.db.Raw("SELECT id,name, location_x, location_y FROM couriers").Scan(&couriers)

	if result.Error != nil {
		return GetAllCouriersResponse{}, result.Error
	}

	return GetAllCouriersResponse{Couriers: couriers}, nil
}
