package cmd

import "delivery/internal/core/domain/services"

type CompositionRoot struct {
	configs Config

	closers []Closer
}

func NewCompositionRoot(configs Config) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
	}
}

func (cr *CompositionRoot) NewOrderDispatcher() services.OrderDispatcher {
	return services.NewOrderDispatcher()
}
