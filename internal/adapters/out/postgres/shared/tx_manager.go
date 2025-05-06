package shared

import (
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"gorm.io/gorm"
)

type TxManager interface {
	Tx() *gorm.DB
	Db() *gorm.DB
	InTx() bool
	Track(agg ddd.AggregateRoot)
	ports.UnitOfWork
}
