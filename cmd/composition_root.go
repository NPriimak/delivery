package cmd

import (
	"delivery/internal/core/domain/services"
	"gorm.io/gorm"
)

type CompositionRoot struct {
	configs Config
	gormDb  *gorm.DB

	closers []Closer
}

func NewCompositionRoot(c Config, gormDb *gorm.DB) CompositionRoot {
	app := CompositionRoot{
		configs: c,
	}
	return app
}

func (cr *CompositionRoot) NewOrderDispatcher() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}
