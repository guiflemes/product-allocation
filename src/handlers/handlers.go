package handlers

import (
	"context"
	"fmt"
	"product-allocation/src/domain"
	"product-allocation/src/service_layer"
)

type uow interface {
	Products() service_layer.Repo
	Rollback()
	Commit()
}

type AddBatchHandler struct {
	uow uow
}

func NewAddBatchHandler(uow uow) *AddBatchHandler {
	return &AddBatchHandler{uow: uow}
}

func (h *AddBatchHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*domain.CreateBatch)

	product, err := h.uow.Products().Get(ctx, cmd.Sku)

	if err != nil {
		return err
	}

	if product == nil {
		product = domain.NewProduct(cmd.Sku, 1)
	}

	product.Batches = append(product.Batches, domain.NewBatch(cmd.Ref, cmd.Sku, cmd.Qty, cmd.Eta))

	if err := h.uow.Products().Add(ctx, product); err != nil {
		h.uow.Rollback()
		return err
	}

	h.uow.Commit()
	return nil
}

type AllocateHandler struct {
	uow uow
}

func NewAllocateHandler(uow uow) *AllocateHandler {
	return &AllocateHandler{uow: uow}
}

func (h *AllocateHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*domain.Allocate)
	line := &domain.OrderLine{OrderId: cmd.OrderId, Sku: cmd.Sku, Qty: cmd.Qty}

	product, err := h.uow.Products().Get(ctx, cmd.Sku)

	if err != nil {
		return err
	}

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

type ChanceBatchQuantity struct {
	uow uow
}

func NewChanceBatchQuantity(uow uow) *ChanceBatchQuantity {
	return &ChanceBatchQuantity{uow: uow}
}

func (b *ChanceBatchQuantity) Handle(ctx context.Context, c interface{}) error {

	cmd := c.(*domain.ChangeBatchQuantity)
	product, err := b.uow.Products().GetByBatchRef(ctx, cmd.Ref)
	if err != nil {
		return err
	}

	product.ChangeBatchQuantity(cmd.Ref, cmd.Qty)
	b.uow.Commit()

	return nil
}
