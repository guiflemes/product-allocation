package handlers

import (
	"context"
	"fmt"
	"product-allocation/src/domain"
	"time"
)

type CreateBatch struct {
	ref string
	sku string
	qty string
	eta time.Time
}

type Repo interface {
	Get(ctx context.Context, sku string) *domain.Product
	AddOrUpdate(context.Context, *domain.Product)
}

type AddHandler struct {
	repo Repo
}

func (h *AddHandler) handler(ctx context.Context, cmd CreateBatch) {
	product := h.repo.Get(ctx, cmd.sku)

	if product == nil {
		product = &domain.Product{
			Sku:     cmd.sku,
			Batches: make([]*domain.Batch, 1),
		}

	}
	product.Batches = append(product.Batches, &domain.Batch{Ref: cmd.ref, Sku: cmd.sku, Qty: cmd.qty, Eta: cmd.eta})

	h.repo.AddOrUpdate(ctx, product)
}

type Allocate struct {
	OrderId string
	Sku     string
	Qty     int
}

type AllocateHandler struct {
	repo Repo
}

func (h *AllocateHandler) handler(ctx context.Context, cmd Allocate) error {
	line := &domain.OrderLine{OrderId: cmd.OrderId, Sku: cmd.Sku, Qty: cmd.Qty}
	product := h.repo.Get(ctx, cmd.Sku)

	if product == nil {
		return fmt.Errorf("Invalid sku %s", cmd.Sku)
	}

	product.Allocate(line)
	h.repo.AddOrUpdate(ctx, product)
	return nil

}
