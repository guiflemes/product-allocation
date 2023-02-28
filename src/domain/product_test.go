package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
			product.Batches = []*Batch{{Ref: "b1", Sku: "SCANDI-PEN", Qty: 100, Eta: time.Now()}}
			product.Allocate(line)
			assert.Equal(product.VersionNumber, scenario.actual)
		})
	}
}
