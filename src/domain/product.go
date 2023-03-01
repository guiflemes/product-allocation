package domain

import (
	"fmt"
	"product-allocation/src/utils/collections"
	"product-allocation/src/utils/math"
	"sort"
	"time"
)

type Product struct {
	Sku           string
	Batches       []*Batch
	VersionNumber int
	Events        *Events
}

func NewProduct(sku string, version int) *Product {
	return &Product{
		Sku:           sku,
		Batches:       make([]*Batch, 0),
		VersionNumber: version,
		Events:        &Events{items: make([]interface{}, 0)},
	}
}

func (p *Product) Allocate(line *OrderLine) string {

	sort.Slice(p.Batches, func(i, j int) bool {
		return p.Batches[i].AvailableQuantity() < p.Batches[j].AvailableQuantity()
	})

	for _, b := range p.Batches {
		if b.CanAllocate(line) {
			b.Allocate(line)
			p.VersionNumber++
			p.Events.Append(&Allocated{
				OrderId:  line.OrderId,
				Sku:      line.Sku,
				Qty:      line.Qty,
				BatchRef: b.Ref,
			})
			return b.Ref
		}

	}
	p.Events.Append(&OutOfStock{line.Sku})
	return ""

}

func (p *Product) ChangeBatchQuantity(ref string, qty int) {
	for _, b := range p.Batches {
		b.Qty = qty
		for b.AvailableQuantity() < 0 {
			line := b.DeallocateOne()
			p.Events.Append(&OutOfStock{line.Sku})
		}
	}
}

func (p *Product) String() string {
	return fmt.Sprintf("Products(Sku=%s, Batches=%v)", p.Sku, p.Batches)
}

type OrderLine struct {
	OrderId string
	Sku     string
	Qty     int
}

type Batch struct {
	Ref         string
	Sku         string
	Qty         int
	Eta         time.Time
	allocations collections.Set[*OrderLine]
}

func NewBatch(ref string, sku string, qty int, eta time.Time) *Batch {
	return &Batch{
		Ref:         ref,
		Sku:         sku,
		Qty:         qty,
		Eta:         eta,
		allocations: *collections.NewSet[*OrderLine](),
	}
}

func (b *Batch) Allocate(line *OrderLine) {

	if b.CanAllocate(line) {
		b.allocations.Add(line)
	}
}

func (b *Batch) DeallocateOne() *OrderLine {
	return b.allocations.Pop()
}

func (b *Batch) AllocateQuantity() int {
	var slice []int
	for o := range b.allocations.Iter() {
		slice = append(slice, o.Qty)
	}

	return math.Sum(slice)
}

func (b *Batch) getPurchasedQuantity() int {
	return b.Qty
}

func (b *Batch) AvailableQuantity() int {
	return b.getPurchasedQuantity() - b.AllocateQuantity()
}

func (b *Batch) CanAllocate(line *OrderLine) bool {
	return b.Sku == line.Sku && b.AvailableQuantity() >= line.Qty
}

func (b *Batch) String() string {
	return fmt.Sprintf("Batch(Ref=%s, Sku=%s, Qty=%d)", b.Ref, b.Sku, b.Qty)
}
