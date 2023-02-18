package adapters

import (
	"context"
	"errors"
	"math/rand"
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

	p := a.products[sku]
	a.seen.Add(p)
	return p, nil

}

func (a *MemoryRepo) GetByBatchRef(cxt context.Context, batchRef string) (*domain.Product, error) {
	// TODO -> Refactor it to get by batchRef, it's random

	if len(a.products) == 0 {
		return nil, errors.New("There is no any product")
	}

	a.m.RLock()
	defer a.m.Unlock()

	randomIndex := rand.Intn(len(a.keys))
	key := a.keys[randomIndex]
	p := a.products[key]
	a.seen.Add(p)
	return p, nil

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
