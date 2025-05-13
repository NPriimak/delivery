package shared

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type TxManager interface {
	Tx() *gorm.DB
	Db() *gorm.DB
	InTx() bool
	Track(agg ddd.AggregateRoot)
	ports.UnitOfWork
}

var _ ports.UnitOfWork = &txManager{}
var _ TxManager = &txManager{}

type txManager struct {
	tx                *gorm.DB
	db                *gorm.DB
	trackedAggregates []ddd.AggregateRoot
}

func NewTxManager(db *gorm.DB) (TxManager, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	tx := &txManager{
		db: db,
	}
	return tx, nil
}

func (u *txManager) Tx() *gorm.DB {
	return u.tx
}

func (u *txManager) Db() *gorm.DB {
	return u.db
}

func (u *txManager) InTx() bool {
	return u.tx != nil
}

func (u *txManager) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *txManager) Begin(ctx context.Context) {
	u.tx = u.db.WithContext(ctx).Begin()
}

func (u *txManager) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	committed := false
	defer func() {
		if !committed {
			if err := u.tx.WithContext(ctx).Rollback().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
				log.Error(err)
			}
			u.clearTx()
		}
	}()

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}
	committed = true
	u.clearTx()

	return nil
}

func (u *txManager) clearTx() {
	u.tx = nil
	u.trackedAggregates = nil
}
