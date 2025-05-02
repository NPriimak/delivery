package cmd

import "delivery/internal/core/domain/services"

type CompositionRoot struct {
	configs Config
}

func NewCompositionRoot(c Config) CompositionRoot {
	app := CompositionRoot{
		configs: c,
	}
	return app
}

func (cr *CompositionRoot) NewOrderDispatcher() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}
