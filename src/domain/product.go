package domain

import (
	"product-allocation/src/utils/collections"
	"product-allocation/src/utils/math"
	"sort"
	"time"
)

type Product struct {
	Sku           string
	Batches       []Batch
	VersionNumber int
	Events        Events
}

func (p *Product) Allocate(line OrderLine) string {

	sort.Slice(p.Batches, func(i, j int) bool {
		return p.Batches[i].Eta.Nanosecond() > p.Batches[j].Eta.Nanosecond()
	})

	for _, b := range p.Batches {
		b.Allocate(line)
		p.VersionNumber += 1
		p.Events.append(
			Allocated{
				OrderId:  line.OrderId,
				Sku:      line.Sku,
				Qty:      line.Qty,
				BatchRef: b.Ref,
			},
		)

		return b.Ref
	}

	p.Events.append(OutOfStock{line.Sku})
	return ""

}

type OrderLine struct {
	OrderId string
	Sku     string
	Qty     int
}

type Batch struct {
	Ref               string
	Sku               string
	Qty               string
	Eta               time.Time
	purchasedQuantity int
	allocations       collections.Set[OrderLine]
}

func (b *Batch) Allocate(line OrderLine) {}

func (b *Batch) DeallocateOne() {}

func (b *Batch) AllocateQuantity() int {
	var slice []int
	for o := range b.allocations {
		slice = append(slice, o.Qty)
	}

	return math.Sum(slice)
}

func (b *Batch) AvailableQuantity() int {
	return b.purchasedQuantity - b.AllocateQuantity()
}

func (b *Batch) CanAllocate(line OrderLine) bool {
	return b.Sku == line.Sku && b.AllocateQuantity() >= line.Qty
}
