package service_layer

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type (
	CommandHandler interface {
		Handle(ctx context.Context, cmd interface{}) error
	}
	EventHandler interface{ Handle(event interface{}) error }
)

type MessageBus struct {
	errChan         chan error
	eventQueue      <-chan interface{}
	commandHandlers map[string]CommandHandler
	eventHandlers   map[string][]EventHandler
	uow             *UnitOfWork
	m               sync.RWMutex
}

func NewMessageBus(uow *UnitOfWork) *MessageBus {
	errorChan := make(chan error)
	eventChan := make(chan interface{})
	commandHandlers := make(map[string]CommandHandler)
	eventHandlers := make(map[string][]EventHandler)

	uow.EventQueue = eventChan

	return &MessageBus{
		errChan:         errorChan,
		eventQueue:      eventChan,
		commandHandlers: commandHandlers,
		eventHandlers:   eventHandlers,
		uow:             uow,
	}
}

func (b *MessageBus) HandlerCommand(ctx context.Context, cmd interface{}) error {
	go b.handlerCommand(ctx, cmd)

	for {
		select {
		case err := <-b.errChan:
			fmt.Println("receiving err")

			if err != nil {
				return err
			}

			return nil

		case event := <-b.eventQueue:
			fmt.Println("receiving event")
			b.HandlerEvent(event)
		}

	}
}

func (b *MessageBus) HandlerEvent(event interface{}) {
	b.m.RLock()

	var wg sync.WaitGroup
	eventType := b.getType(event)
	handlers, ok := b.eventHandlers[eventType]

	if !ok {
		return
	}

	for _, handler := range handlers {
		wg.Add(1)

		go func(e interface{}, h EventHandler) {
			defer wg.Done()
			defer b.m.RUnlock()
			h.Handle(e)
		}(event, handler)
	}

	wg.Wait()

}

func (b *MessageBus) handlerCommand(ctx context.Context, cmd interface{}) {
	b.m.RLock()
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

func (b *MessageBus) RegisterCommandHandler(commandName string, handler CommandHandler) {
	b.m.Lock()
	defer b.m.Unlock()
	b.commandHandlers[commandName] = handler
}

func (b *MessageBus) RegisterEventHandler(eventName string, handler EventHandler) {
	b.m.Lock()
	defer b.m.Unlock()

	handlers, ok := b.eventHandlers[eventName]

	if !ok {
		b.eventHandlers[eventName] = []EventHandler{handler}
		return
	}

	handlers = append(handlers, handler)

}

func (b *MessageBus) getType(s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
