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
	Events        []interface{}
}

func NewProduct(sku string, version int) *Product {
	return &Product{
		Sku:           sku,
		Batches:       make([]*Batch, 0),
		VersionNumber: version,
		Events:        make([]interface{}, 0),
	}
}

func (p *Product) Allocate(line *OrderLine) {

	sort.Slice(p.Batches, func(i, j int) bool {
		return p.Batches[i].Eta.Nanosecond() > p.Batches[j].Eta.Nanosecond()
	})

	for _, b := range p.Batches {
		b.Allocate(line)
		p.VersionNumber += 1

		p.Events = append(p.Events, &Allocated{
			OrderId:  line.OrderId,
			Sku:      line.Sku,
			Qty:      line.Qty,
			BatchRef: b.Ref,
		})

	}

	p.Events = append(p.Events, &OutOfStock{line.Sku})

}

func (p *Product) ChangeBatchQuantity(ref string, qty int) {
	for _, b := range p.Batches {
		b.purchasedQuantity = qty
		for b.AvailableQuantity() < 0 {
			line := b.DeallocateOne()
			p.Events = append(p.Events, &Deallocate{line.OrderId, line.Sku, line.Qty})
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
	Ref               string
	Sku               string
	Qty               int
	Eta               time.Time
	purchasedQuantity int
	allocations       collections.Set[*OrderLine]
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

func (b *Batch) AvailableQuantity() int {
	return b.purchasedQuantity - b.AllocateQuantity()
}

func (b *Batch) CanAllocate(line *OrderLine) bool {
	return b.Sku == line.Sku && b.AllocateQuantity() >= line.Qty
}

func (b *Batch) String() string {
	return fmt.Sprintf("Batch(Ref=%s, Sku=%s, Qty=%d)", b.Ref, b.Sku, b.Qty)
}
