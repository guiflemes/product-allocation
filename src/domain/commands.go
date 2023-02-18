package domain

import "time"

type Allocate struct {
	OrderId string
	Sku     string
	Qty     int
}

type CreateBatch struct {
	Ref string
	Sku string
	Qty int
	Eta time.Time
}

type ChangeBatchQuantity struct {
	Ref string
	Qty int
}
