package queries

import (
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetNotCompletedOrdersQuery struct {
	isSet bool
}

func NewGetNotCompletedOrdersQuery() (GetNotCompletedOrdersQuery, error) {
	return GetNotCompletedOrdersQuery{
		isSet: true,
	}, nil
}

func (q GetNotCompletedOrdersQuery) IsEmpty() bool {
	return !q.isSet
}

type GetNotCompletedOrdersResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	ID       uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Location LocationResponse `gorm:"embedded;embeddedPrefix:location_"`
}

func (OrderResponse) TableName() string {
	return "orders"
}

type GetNotCompletedOrdersQueryHandler interface {
	Handle(GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error)
}

type getNotCompletedOrdersQueryHandler struct {
	db *gorm.DB
}

func NewGetNotCompletedOrdersQueryHandler(db *gorm.DB) (GetNotCompletedOrdersQueryHandler, error) {
	if db == nil {
		return &getNotCompletedOrdersQueryHandler{}, errs.NewValueIsRequiredError("db")
	}
	return &getNotCompletedOrdersQueryHandler{db: db}, nil
}

func (q *getNotCompletedOrdersQueryHandler) Handle(query GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error) {
	if query.IsEmpty() {
		return GetNotCompletedOrdersResponse{}, errs.NewValueIsRequiredError("query")
	}

	var orders []OrderResponse
	result := q.db.Raw("SELECT id, courier_id, location_x, location_y, status FROM orders where status!=?",
		order.StatusCompleted).Scan(&orders)

	if result.Error != nil {
		return GetNotCompletedOrdersResponse{}, result.Error
	}

	return GetNotCompletedOrdersResponse{Orders: orders}, nil
}
