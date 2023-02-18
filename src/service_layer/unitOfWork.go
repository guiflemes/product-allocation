package service_layer

import (
	"context"
	"product-allocation/src/adapters"
	"product-allocation/src/domain"
	"product-allocation/src/utils/collections"
)

type Repo interface {
	Add(ctx context.Context, p *domain.Product) error
	Get(ctx context.Context, sku string) (*domain.Product, error)
	GetByBatchRef(ctx context.Context, batchRef string) (*domain.Product, error)
	Seen() *collections.Set[*domain.Product]
}

type UnitOfWork struct {
	products   Repo
	EventQueue chan<- interface{}
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

func NewTestUow() *UnitOfWork {
	return &UnitOfWork{products: adapters.NewMemoryRepo()}
}
