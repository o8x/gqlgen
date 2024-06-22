package graph

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/99designs/gqlgen/_examples/sse/graph/model"
)

type TodoObserver struct {
	ID         string
	ListenFunc func(*model.TodoChangeMessage)
}

type TodoObserverManager struct {
	observers map[string]TodoObserver
}

func (t *TodoObserverManager) Register(o TodoObserver) {
	if _, ok := t.observers[o.ID]; !ok {
		t.observers[o.ID] = o
	}
}

func (t *TodoObserverManager) Unregister(o TodoObserver) {
	if _, ok := t.observers[o.ID]; !ok {
		delete(t.observers, o.ID)
	}
}

func (t *TodoObserverManager) Notify(msg *model.TodoChangeMessage) {
	for _, it := range t.observers {
		it.ListenFunc(msg)
	}
}

var observerManager = TodoObserverManager{
	observers: map[string]TodoObserver{},
}

var todos []*model.Todo

func createTodoHandler(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	input.Text = strings.TrimSpace(input.Text)

	if input.Text == "" {
		return nil, fmt.Errorf("text required")
	}

	newTodo := &model.Todo{
		ID:   uuid.NewString(),
		Text: input.Text,
		Done: input.Done,
	}

	todos = append(todos, newTodo)
	defer observerManager.Notify(&model.TodoChangeMessage{
		NewValue: newTodo,
		OldValue: nil,
		Action:   model.TodoChangeActionNew,
	})
	return newTodo, nil
}

func toggleTodoDoneHandler(ctx context.Context, id string) (*model.Todo, error) {
	for _, it := range todos {
		if it.ID == id {
			var todo = *it
			it.Done = !it.Done

			observerManager.Notify(&model.TodoChangeMessage{
				NewValue: it,
				OldValue: &todo,
				Action:   model.TodoChangeActionToggleDone,
			})
			return it, nil
		}
	}

	return nil, fmt.Errorf("todo item not found")
}

func getTodoListHandler(ctx context.Context) ([]*model.Todo, error) {
	return todos, nil
}

func onTodoChangeHandler(ctx context.Context) (<-chan *model.TodoChangeMessage, error) {
	ch := make(chan *model.TodoChangeMessage)

	go func() {
		defer close(ch)

		observer := TodoObserver{
			ID: uuid.NewString(),
			ListenFunc: func(msg *model.TodoChangeMessage) {
				ch <- msg
			},
		}

		observerManager.Register(observer)
		defer observerManager.Unregister(observer)

		<-ctx.Done()
		return
	}()

	return ch, nil
}
