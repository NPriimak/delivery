package orderrepo

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.OrderRepository = &Repository{}

type Repository struct {
	uow ports.UnitOfWork
}

func NewRepository(uow ports.UnitOfWork) (*Repository, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &Repository{
		uow: uow,
	}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *order.Order) error {
	r.uow.Track(aggregate)

	dto := DomainToDTO(aggregate)

	// Открыта ли транзакция?
	isInTransaction := r.uow.InTx()
	if !isInTransaction {
		r.uow.Begin(ctx)
	}
	tx := r.uow.Tx()

	// Вносим изменения
	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}

	// Если не было внешней в транзакции, то коммитим изменения
	if !isInTransaction {
		err := r.uow.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *order.Order) error {
	r.uow.Track(aggregate)

	dto := DomainToDTO(aggregate)

	// Открыта ли транзакция?
	isInTransaction := r.uow.InTx()
	if !isInTransaction {
		r.uow.Begin(ctx)
	}
	tx := r.uow.Tx()

	// Вносим изменения
	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}

	// Если не было внешней в транзакции, то коммитим изменения
	if !isInTransaction {
		err := r.uow.Commit(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*order.Order, error) {
	dto := OrderDTO{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, ID)
	if result.RowsAffected == 0 {
		return nil, nil
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error) {
	dto := OrderDTO{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Where("status = ?", order.StatusCreated).
		First(&dto)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errs.NewObjectNotFoundError("Free courier", nil)
		}
		return nil, result.Error
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error) {
	var dtos []OrderDTO

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Where("status = ?", order.StatusAssigned).
		Find(&dtos)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("Assigned orders", nil)
	}

	aggregates := make([]*order.Order, len(dtos))
	for i, dto := range dtos {
		aggregates[i] = DtoToDomain(dto)
	}

	return aggregates, nil
}

func (r *Repository) getTxOrDb() *gorm.DB {
	if tx := r.uow.Tx(); tx != nil {
		return tx
	}
	return r.uow.Db()
}
