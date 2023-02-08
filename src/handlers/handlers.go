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
	Add(ctx context.Context, p *domain.Product) error
}

type uow interface {
	Products() Repo
	Rollback()
	Commit()
}

type AddBatchHandler struct {
	uow uow
}

func (h *AddBatchHandler) Handler(ctx context.Context, cmd CreateBatch) error {
	product := h.uow.Products().Get(ctx, cmd.sku)

	if product == nil {
		product = &domain.Product{
			Sku:     cmd.sku,
			Batches: make([]*domain.Batch, 1),
		}

	}
	product.Batches = append(product.Batches, &domain.Batch{Ref: cmd.ref, Sku: cmd.sku, Qty: cmd.qty, Eta: cmd.eta})

	if err := h.uow.Products().Add(ctx, product); err != nil {
		h.uow.Rollback()
		return err
	}

	h.uow.Commit()
	return nil
}

type Allocate struct {
	OrderId string
	Sku     string
	Qty     int
}

type AllocateHandler struct {
	uow uow
}

func (h *AllocateHandler) Handler(ctx context.Context, cmd Allocate) error {
	line := &domain.OrderLine{OrderId: cmd.OrderId, Sku: cmd.Sku, Qty: cmd.Qty}
	product := h.uow.Products().Get(ctx, cmd.Sku)

	if product == nil {
		return fmt.Errorf("Invalid sku %s", cmd.Sku)
	}

	product.Allocate(line)

	if err := h.uow.Products().Add(ctx, product); err != nil {
		h.uow.Rollback()
		return err
	}

	h.uow.Commit()
	return nil

}
