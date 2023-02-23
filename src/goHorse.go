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

type FakeEvent struct{}

func (f *FakeEvent) Handle(e interface{}) error {
	fmt.Println("receive event", e)
	return nil
}

func bootstrap(uow *service_layer.UnitOfWork) *service_layer.MessageBus {
	bus := service_layer.NewMessageBus(uow)

	bus.RegisterCommandHandler("Allocate", handlers.NewAllocateHandler(uow))
	bus.RegisterCommandHandler("CreateBatch", handlers.NewAddBatchHandler(uow))
	bus.RegisterEventHandler("OutOfStock", &FakeEvent{})
	return bus
}

var options string = `Choose some action
1 - Create new Batch.
2 - Allocate products.
3 - Stop
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
			break
		default:
			continue
		}

		bus.HandlerCommand(context.Background(), cmd)

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
