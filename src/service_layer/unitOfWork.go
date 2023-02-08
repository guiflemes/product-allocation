package service_layer

import (
	"product-allocation/src/domain"
	"product-allocation/src/utils/collections"
)

type Repo interface {
	Add(p *domain.Product)
	Get(sku string) (*domain.Product, error)
	GetByBatchRef(batchRef string) (*domain.Product, error)
	Seen() collections.Set[*domain.Product]
}

type UnitOfWork struct {
	repo       Repo
	EventQueue chan<- Event
}

func (u *UnitOfWork) CollectNewEvents() {
	products := u.repo.Seen()

	for e := range products.Iter() {
		u.EventQueue <- e
	}
}
