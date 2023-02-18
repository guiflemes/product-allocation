package main

import (
	"context"
	"fmt"
	"product-allocation/src/domain"
	"product-allocation/src/handlers"
	"product-allocation/src/service_layer"
	"time"
)

type FakeEvent struct{}

func (f *FakeEvent) Handler(e interface{}) {
	fmt.Println("receive event")
}

func bootstrap(uow *service_layer.UnitOfWork) *service_layer.MessageBus {
	bus := service_layer.NewMessageBus(uow)

	bus.RegisterCommandHandler("Allocate", handlers.NewAddBatchHandler(uow))
	bus.RegisterCommandHandler("CreateBatch", handlers.NewAddBatchHandler(uow))
	return bus
}

func GoHorse() {
	batch := &domain.CreateBatch{
		Ref: "ref1",
		Sku: "sku1",
		Qty: 2,
		Eta: time.Now(),
	}

	uow := service_layer.NewTestUow()
	bus := bootstrap(uow)
	bus.HandlerCommand(context.Background(), batch)

	fmt.Println(uow.Products().Seen())
}
