package service_layer

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type (
	Event          interface{}
	Command        interface{}
	CommandHandler interface {
		Handle(ctx context.Context, cmd Command) error
	}
	EventHandler interface{ Handle(event Event) error }
)

type messageBus struct {
	commandHandlers map[string]CommandHandler
	eventHandlers   map[string][]EventHandler
	errChan         chan error
	EventQueue      <-chan Event
	uow             UnitOfWork
	m               sync.RWMutex
}

func (b *messageBus) HandlerCommand(ctx context.Context, cmd Command) error {
	go b.handlerCommand(ctx, cmd)

	for {
		select {
		case err := <-b.errChan:

			if err != nil {
				return err
			}

			return nil

		case event := <-b.EventQueue:
			b.HandlerEvent(event)
		}

	}
}

func (b *messageBus) HandlerEvent(event Event) {
	b.m.RLocker()

	var wg sync.WaitGroup
	eventType := b.getType(event)
	handlers, ok := b.eventHandlers[eventType]

	if !ok {
		return
	}

	for _, handler := range handlers {
		wg.Add(1)

		go func(e Event, h EventHandler) {
			defer wg.Done()
			defer b.m.RUnlock()
			h.Handle(e)
		}(event, handler)
	}

	wg.Wait()

}

func (b *messageBus) handlerCommand(ctx context.Context, cmd interface{}) {
	b.m.RLocker()
	defer b.m.RUnlock()

	cmdType := b.getType(cmd)
	handler, ok := b.commandHandlers[cmdType]

	if ok {
		b.errChan <- handler.Handle(ctx, cmd)
		b.uow.CollectNewEvents()
		return
	}

	b.errChan <- fmt.Errorf("unable to find handler cmd %s", cmdType)

}

func (b *messageBus) RegisterCommandHandler(commandName string, handler CommandHandler) {
	b.m.Lock()
	defer b.m.Unlock()
	b.commandHandlers[commandName] = handler
}

func (b *messageBus) RegisterEventHandler(eventName string, handler EventHandler) {
	b.m.Lock()
	defer b.m.Unlock()

	handlers, ok := b.eventHandlers[eventName]

	if !ok {
		b.eventHandlers[eventName] = []EventHandler{handler}
		return
	}

	handlers = append(handlers, handler)

}

func (b *messageBus) getType(s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
