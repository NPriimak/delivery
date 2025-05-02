package courierrepo

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CourierRepository = &Repository{}

type Repository struct {
	uow ports.UnitOfWork
}

func NewRepository(uow ports.UnitOfWork) (*Repository, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}

	return &Repository{uow}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	if aggregate == nil {
		return errs.NewValueIsRequiredError("aggregate")
	}

	r.uow.Track(aggregate)
	dto := DomainToDTO(aggregate)

	isInTransaction := r.uow.InTx()
	if !isInTransaction {
		r.uow.Begin(ctx)
	}

	tx := r.uow.Tx()

	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Create(&dto).Error
	if err != nil {
		return err
	}

	if !isInTransaction {
		err := r.uow.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	if aggregate == nil {
		return errs.NewValueIsRequiredError("aggregate")
	}

	r.uow.Track(aggregate)

	dto := DomainToDTO(aggregate)
	isInTransaction := r.uow.InTx()
	if !isInTransaction {
		r.uow.Begin(ctx)
	}
	tx := r.uow.Tx()

	err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(&dto).Error
	if err != nil {
		return err
	}

	if !isInTransaction {
		err := r.uow.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	if ID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("ID")
	}

	dto := CourierDTO{}

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Find(&dto, ID)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("Courier by ID", ID)
	}

	aggregate := DtoToDomain(dto)
	return aggregate, nil
}

func (r *Repository) GetAllFree(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []CourierDTO

	tx := r.getTxOrDb()
	result := tx.WithContext(ctx).
		Preload(clause.Associations).
		Where(`NOT EXISTS (
            SELECT 1 FROM storage_places sp
            WHERE sp.courier_id = couriers.id AND sp.order_id IS NOT NULL
        )`).Find(&dtos)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errs.NewObjectNotFoundError("Free couriers", nil)
	}

	aggregates := make([]*courier.Courier, len(dtos))
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
