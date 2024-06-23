package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/99designs/gqlgen/_examples/sse/graph/model"
	"github.com/99designs/gqlgen/client"
)

func TestSSE(t *testing.T) {
	c := client.New(NewServer())

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//	var resp any
	//	// BUG: the method has not been fully tested, has an unknown bug here and will always be blocked here. Instead of blocking in the Next method
	//	sse := c.SSE(context.Background(), `subscription { onTodoChange {
	//      action oldValue { id text done } newValue { id text done }
	//	}}`)
	//
	//	wg.Done()
	//	require.NoError(t, sse.Next(&resp))
	//}()
	//
	//// Ensure subscription is started
	//wg.Wait()

	var resp struct {
		CreateTodo     *model.Todo `json:"createTodo"`
		ToggleTodoDone *model.Todo `json:"toggleTodoDone"`
	}

	var todo1ID string
	var todo2ID string

	err := c.Post(`mutation { createTodo(input:{ text: "new complete todo" done: true }) { id text done } }`, &resp)
	require.NoError(t, err)
	require.True(t, resp.CreateTodo.Done)
	require.True(t, resp.CreateTodo.Text == "new complete todo")
	todo1ID = resp.CreateTodo.ID

	err = c.Post(`mutation { createTodo(input:{ text: "new todo" done: false }) { id text done } }`, &resp)
	require.NoError(t, err)
	require.False(t, resp.CreateTodo.Done)
	require.True(t, resp.CreateTodo.Text == "new todo")
	todo2ID = resp.CreateTodo.ID

	// Done True -> False
	err = c.Post(fmt.Sprintf(`mutation { toggleTodoDone(id: "%s") { id text done } }`, todo1ID), &resp)
	require.NoError(t, err)
	require.False(t, resp.ToggleTodoDone.Done)
	require.True(t, resp.ToggleTodoDone.ID == todo1ID)

	// Done False -> True
	err = c.Post(fmt.Sprintf(`mutation { toggleTodoDone(id: "%s") { id text done } }`, todo1ID), &resp)
	require.NoError(t, err)
	require.True(t, resp.ToggleTodoDone.Done)
	require.True(t, resp.ToggleTodoDone.ID == todo1ID)

	// Done True -> False
	err = c.Post(fmt.Sprintf(`mutation { toggleTodoDone(id: "%s") { id text done } }`, todo1ID), &resp)
	require.NoError(t, err)
	require.False(t, resp.ToggleTodoDone.Done)
	require.True(t, resp.ToggleTodoDone.ID == todo1ID)

	var todos struct {
		Todos []*model.Todo `json:"todos"`
	}
	err = c.Post("query { todos { id text done } }", &todos)
	require.NoError(t, err)
	require.NotNil(t, todos.Todos)
	require.True(t, len(todos.Todos) == 2)
	require.True(t, todos.Todos[0].ID == todo1ID)
	require.True(t, todos.Todos[1].ID == todo2ID)
}
