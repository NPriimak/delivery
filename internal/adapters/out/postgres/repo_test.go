package postgres

import (
	"context"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/adapters/out/postgres/shared"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/testcnts"
	"github.com/stretchr/testify/assert"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func setupTest(t *testing.T) (context.Context, *gorm.DB) {
	ctx := context.Background()
	postgresContainer, dsn, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		assert.NoError(t, err)
	}

	// Подключаемся к БД через Gorm
	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Авто миграция (создаём таблицу)
	err = db.AutoMigrate(&courierrepo.CourierDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&orderrepo.OrderDTO{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	assert.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	return ctx, db
}

func createTxManager(t *testing.T, db *gorm.DB) shared.TxManager {
	tx, err := shared.NewTxManager(db)
	assert.NoError(t, err)
	return tx
}

func createTestLocation(t *testing.T, x uint8, y uint8) kernel.Location {
	result, err := kernel.NewLocation(x, y)
	if err != nil {
		assert.NoError(t, err)
	}
	return result
}
