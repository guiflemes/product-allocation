package service_layer

import (
	"context"
	"fmt"
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

	for p := range products.Iter() {
		// TODO fix it to send events to EventQueue, nothing is sending to chanel, why?
		for _, e := range p.Events {
			fmt.Printf("sending %v to EventQueue\n", e)
			u.EventQueue <- e

		}

	}
}

func (u *UnitOfWork) Rollback() {}
func (u *UnitOfWork) Commit()   {}

func NewTestUow() *UnitOfWork {
	return &UnitOfWork{products: adapters.NewMemoryRepo()}
}
