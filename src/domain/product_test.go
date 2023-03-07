package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	today    time.Time = time.Now()
	tomorrow time.Time = today.AddDate(0, 0, 1)
	later    time.Time = tomorrow.AddDate(0, 0, 10)
)

func TestPrefersEarlierBatchs(t *testing.T) {
	assert := assert.New(t)

	earliest := NewBatch("speedy-batch", "MINIMALIST-SPOON", 100, today)
	medium := NewBatch("normal-batch", "MINIMALIST-SPOON", 100, tomorrow)
	latest := NewBatch("slow-batch", "MINIMALIST-SPOON", 100, later)
	product := NewProduct("MINIMALIST-SPOON", 1)
	product.Batches = []*Batch{earliest, medium, latest}
	line := &OrderLine{"order_1", "MINIMALIST-SPOON", 10}

	product.Allocate(line)

	assert.Equal(earliest.AvailableQuantity(), 90)
	assert.Equal(medium.AvailableQuantity(), 100)
	assert.Equal(latest.AvailableQuantity(), 100)

}

func TestReturnsAllocatedBatchRef(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		desc          string
		firstAddBatch *Batch
		nextBatchs    []*Batch
		line          *OrderLine
		product       *Product
	}

	for _, scenario := range []testCase{
		{
			desc:          "Allocation to Ref in-stock-batch-ref",
			firstAddBatch: NewBatch("in-stock-batch-ref", "MARVIN-DOG", 100, time.Now()),
			nextBatchs:    []*Batch{NewBatch("shipment-batch-ref", "MARVIN-DOG", 100, tomorrow)},
			line:          &OrderLine{OrderId: "order1", Sku: "MARVIN-DOG", Qty: 10},
			product:       NewProduct("MARVIN-DOG", 10),
		},
		{
			desc:          "Allocation to Ref in-stock-marvin-ref",
			firstAddBatch: NewBatch("in-stock-marvin-ref", "MARVIN-DOG", 100, time.Now()),
			nextBatchs:    []*Batch{NewBatch("shipment-batch-ref_1", "MARVIN-DOG", 100, tomorrow), NewBatch("shipment-batch-ref_2", "MARVIN-DOG", 100, tomorrow.AddDate(0, 0, 1))},
			line:          &OrderLine{OrderId: "order1", Sku: "MARVIN-DOG", Qty: 10},
			product:       NewProduct("MARVIN-DOG", 10),
		},
	} {
		t.Run(scenario.desc, func(t *testing.T) {
			batchs := scenario.nextBatchs
			batchs = append(batchs, scenario.firstAddBatch)
			scenario.product.Batches = batchs
			allocation := scenario.product.Allocate(scenario.line)
			assert.Equal(allocation, scenario.firstAddBatch.Ref)
		})
	}
}

func TestOutputsAllocatedEvent(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		desc     string
		line     *OrderLine
		product  *Product
		batchs   []*Batch
		expected *Allocated
	}

	for _, scenario := range []testCase{
		{
			desc:     "Allocated to Sku MARVIN-DOG",
			line:     &OrderLine{OrderId: "order1", Sku: "MARVIN-DOG", Qty: 10},
			product:  NewProduct("MARVIN-DOG", 10),
			batchs:   []*Batch{NewBatch("batch1", "MARVIN-DOG", 100, time.Now())},
			expected: &Allocated{OrderId: "order1", Sku: "MARVIN-DOG", Qty: 10, BatchRef: "batch1"},
		},
		{
			desc:     "Allocated to Sku MARVIN-HOUSE",
			line:     &OrderLine{OrderId: "order2", Sku: "MARVIN-HOUSE", Qty: 11},
			product:  NewProduct("MARVIN-HOUSE", 11),
			batchs:   []*Batch{NewBatch("batch2", "MARVIN-HOUSE", 200, time.Now())},
			expected: &Allocated{OrderId: "order2", Sku: "MARVIN-HOUSE", Qty: 11, BatchRef: "batch2"},
		},
		{
			desc:     "Allocated to Sku MARVIN-CAR",
			line:     &OrderLine{OrderId: "order3", Sku: "MARVIN-CAR", Qty: 15},
			product:  NewProduct("MARVIN-CAR", 15),
			batchs:   []*Batch{NewBatch("batch3", "MARVIN-CAR", 49, time.Now())},
			expected: &Allocated{OrderId: "order3", Sku: "MARVIN-CAR", Qty: 15, BatchRef: "batch3"},
		},
	} {
		t.Run(scenario.desc, func(t *testing.T) {
			scenario.product.Batches = scenario.batchs
			scenario.product.Allocate(scenario.line)
			assert.Equal(scenario.product.Events.Read()[len(scenario.product.Events.Read())-1], scenario.expected)
		})
	}
}

func TestRecordOutOfStockEventIfCannotAllocate(t *testing.T) {
	assert := assert.New(t)
	product := NewProduct("Sku1", 1)
	product.Batches = []*Batch{NewBatch("batch1", "Sku1", 10, time.Now())}
	product.Allocate(&OrderLine{OrderId: "order1", Sku: "Sku1", Qty: 10})

	allocation := product.Allocate(&OrderLine{OrderId: "order2", Sku: "Sku1", Qty: 1})

	assert.Equal(product.Events.Read()[len(product.Events.Read())-1], &OutOfStock{"Sku1"})
	assert.Equal("", allocation)
}

func TestIncrementsVersionNumber(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		previous int
		actual   int
	}

	for _, scenario := range []testCase{
		{
			previous: 1,
			actual:   2,
		},
		{
			previous: 7,
			actual:   8,
		},
	} {
		t.Run(fmt.Sprintf("From %d to %d", scenario.previous, scenario.actual), func(t *testing.T) {
			line := &OrderLine{"oref", "SCANDI-PEN", 10}
			product := NewProduct("SCANDI-PEN", scenario.previous)
			product.Batches = []*Batch{NewBatch("b1", "SCANDI-PEN", 100, time.Now())}
			product.Allocate(line)
			assert.Equal(product.VersionNumber, scenario.actual)
		})
	}
}
