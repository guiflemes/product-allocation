package service_layer

import (
	"product-allocation/src/domain"
	"product-allocation/src/utils/collections"
)

type Repo interface {
	Add(p *domain.Product) error
	Get(sku string) (*domain.Product, error)
	GetByBatchRef(batchRef string) (*domain.Product, error)
	Seen() collections.Set[*domain.Product]
}

type UnitOfWork struct {
	products   Repo
	EventQueue chan<- Event
}

func (u *UnitOfWork) Products() Repo {
	return u.products
}

func (u *UnitOfWork) CollectNewEvents() {
	products := u.products.Seen()

	for e := range products.Iter() {
		u.EventQueue <- e
	}
}

func (u *UnitOfWork) Rollback() {}
func (u *UnitOfWork) Commit()   {}
