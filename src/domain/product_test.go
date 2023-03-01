package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
