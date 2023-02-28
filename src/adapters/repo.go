package adapters

import (
	"context"
	"errors"
	"product-allocation/src/domain"
	"product-allocation/src/utils/collections"
	"sync"
)

type MemoryRepo struct {
	seen          *collections.Set[*domain.Product]
	products      map[string]*domain.Product
	keys          []string
	sliceKeyIndex map[string]int
	m             sync.RWMutex
	batch_count   int
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		seen:          collections.NewSet[*domain.Product](),
		products:      make(map[string]*domain.Product),
		keys:          make([]string, 0),
		sliceKeyIndex: make(map[string]int),
	}
}

func (a *MemoryRepo) Seen() *collections.Set[*domain.Product] {
	return a.seen
}

func (a *MemoryRepo) Add(cxt context.Context, p *domain.Product) error {
	a.m.Lock()
	defer a.m.Unlock()

	a.seen.Add(p)
	a.products[p.Sku] = p
	a.keys = append(a.keys, p.Sku)
	index := len(a.keys) - 1
	a.sliceKeyIndex[p.Sku] = index
	return nil
}

func (a *MemoryRepo) Get(cxt context.Context, sku string) (*domain.Product, error) {
	a.m.RLock()
	defer a.m.RUnlock()

	p, ok := a.products[sku]

	if ok {
		a.seen.Add(p)
	}

	return p, nil

}

func (a *MemoryRepo) GetByBatchRef(cxt context.Context, batchRef string) (*domain.Product, error) {
	a.m.RLock()
	defer a.m.RUnlock()
	if len(a.products) == 0 {
		return nil, errors.New("There is no any product error")
	}

	for _, p := range a.products {
		for _, b := range p.Batches {
			if b.Ref == batchRef {
				a.seen.Add(p)
				return p, nil
			}
		}
	}

	return nil, errors.New("Product not found error")

}

// func AddOrUpdate(db *sql.DB, product Product) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	_, err = tx.Exec("UPDATE product SET name = ? WHERE id = ?", product.Name, product.ID)
// 	if err != nil {
// 		return err
// 	}

// 	for _, batch := range product.Batch {
// 		if batch.ID == 0 {
// 			res, err := tx.Exec("INSERT INTO batch (product_id, quantity) VALUES (?, ?)", product.ID, batch.Quantity)
// 			if err != nil {
// 				return err
// 			}
// 			batchID, err := res.LastInsertId()
// 			if err != nil {
// 				return err
// 			}
// 			batch.ID = int(batchID)
// 		} else {
// 			_, err := tx.Exec("UPDATE batch SET quantity = ? WHERE id = ?", batch.Quantity, batch.ID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return tx.Commit()
// }
