package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"product-allocation/src/domain"
	"product-allocation/src/handlers"
	"product-allocation/src/service_layer"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FakeOutOfStockEvent struct{}

func (f *FakeOutOfStockEvent) Handle(e interface{}) error {
	event := e.(*domain.OutOfStock)
	fmt.Println("OutOfStock EventHandler executed", event.Sku)
	return nil
}

type FakeAllocatedEvent struct{}

func (f *FakeAllocatedEvent) Handle(e interface{}) error {
	event := e.(*domain.Allocated)
	fmt.Println("Allocated EventHandler executed", event.Sku)
	return nil
}

func bootstrap(uow *service_layer.UnitOfWork) *service_layer.MessageBus {
	bus := service_layer.NewMessageBus(uow)

	bus.RegisterCommandHandler("Allocate", handlers.NewAllocateHandler(uow))
	bus.RegisterCommandHandler("CreateBatch", handlers.NewAddBatchHandler(uow))
	bus.RegisterCommandHandler("ChangeBatchQuantity", handlers.NewChanceBatchQuantity(uow))
	bus.RegisterEventHandler("OutOfStock", &FakeOutOfStockEvent{})
	bus.RegisterEventHandler("Allocated", &FakeAllocatedEvent{})
	return bus
}

var options string = `Choose some action
1 - Create new Batch.
2 - Allocate products.
3 - Change Batch Qty.
4 - Stop
`

func GoHorse() {
	uow := service_layer.NewTestUow()
	bus := bootstrap(uow)

	var input int

	for {
		fmt.Println(options)
		fmt.Scanf("%d", &input)

		var cmd interface{}

		switch input {
		case 1:
			cmd = createBatchCmd()
		case 2:
			cmd = createAllocateCmd()
		case 3:
			cmd = createChangeBatchQuantityCmd()
		case 4:
			break
		default:
			continue
		}

		err := bus.HandlerCommand(context.Background(), cmd)

		func() {
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("command executed successfully")
		}()

	}

}

func createBatchCmd() *domain.CreateBatch {
	fmt.Println("CreateBatchCmd")
	batch := &domain.CreateBatch{}
	parseInputs(batch, []string{"Eta"})
	batch.Eta = time.Now()
	return batch
}

func createAllocateCmd() *domain.Allocate {
	fmt.Println("AllocateCmd")
	allocate := &domain.Allocate{}
	parseInputs(allocate, make([]string, 0))
	return allocate
}

func createChangeBatchQuantityCmd() *domain.ChangeBatchQuantity {
	fmt.Println("ChangeBatchQuantityCmd")
	changeBatchQty := &domain.ChangeBatchQuantity{}
	parseInputs(changeBatchQty, make([]string, 0))
	return changeBatchQty
}

func parseInputs(s interface{}, skip []string) {
	reader := bufio.NewReader(os.Stdin)

	v := func(s interface{}) reflect.Value {
		t := reflect.TypeOf(s)
		if t.Kind() == reflect.Ptr {
			return reflect.ValueOf(s).Elem()
		}
		return reflect.ValueOf(s)
	}(s)

	values := make([]interface{}, v.NumField())

	for i := range values {
		name := v.Type().Field(i).Name
		value := v.Field(i)
		fieldType := value.Type()

		toSkip := func() bool {
			for _, s := range skip {
				if name == s {
					return true
				}
			}
			return false
		}()

		if toSkip {
			continue
		}

		fmt.Println("Enter with a value to", name)
		input, err := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if err != nil {
			fmt.Println("An error occurred while reading input. Please try again", err)
		}

		v.Field(i).Set(reflect.ValueOf(func() interface{} {
			switch fieldType.Kind() {
			case reflect.Int:
				in, err := strconv.Atoi(input)

				if err != nil {
					panic(err)
				}

				return in
			default:
				return input
			}
		}()))

	}
}
