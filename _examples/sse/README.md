SSE example
=====

## Start The Application

```shell
go run . -v
```

## Mutations

The following operations can be run in GraphQL playground

- CreateTodo

```graphql
mutation {
    createTodo(input:{
        text: "new complete todo"
        done: true
    }) {
        id
        text
        done
    }
}
```

- ToggleTodoDone

may be executed multiple times

```graphql
mutation {
    toggleTodoDone(id: "ce1a1dea-c769-4d8e-a531-0458a304a5d2") {
        id
        text
        done
    }
}
```

### Subscription

```graphql
subscription {
    onTodoChange {
        action
        oldValue { id text done }
        newValue { id text done }
    }
}
```

The GraphQL playground does not support SSE yet. We can try out the subscription via curl.

```shell
curl --verbose -X POST \
    --url http://127.0.0.1:5001/graphql \
    -H "Accept: text/event-stream" \
    -H 'Content-Type: application/json' \
    --data '{"query":"subscription { onTodoChange { action oldValue { id text done } newValue { id text done } } }"}'
```

***important: Suppose this command is started before Mutations are executed***

The output is similar to the following

```shell
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:5001...
* Connected to 127.0.0.1 (127.0.0.1) port 5001
> POST /query HTTP/1.1
> Host: 127.0.0.1:5001
> User-Agent: curl/8.4.0
> Accept: text/event-stream
> Content-Type: application/json
> Content-Length: 104
>
< HTTP/1.1 200 OK
< Cache-Control: no-cache
< Connection: keep-alive
< Content-Type: text/event-stream
< Date: Sat, 22 Jun 2024 15:45:16 GMT
< Transfer-Encoding: chunked
<
:

event: next
data: {"data":{"onTodoChange":{"action":"New","oldValue":null,"newValue":{"id":"ce1a1dea-c769-4d8e-a531-0458a304a5d2","text":"new complete todo","done":true}}}}

event: next
data: {"data":{"onTodoChange":{"action":"ToggleDone","oldValue":{"id":"ce1a1dea-c769-4d8e-a531-0458a304a5d2","text":"new complete todo","done":true},"newValue":{"id":"ce1a1dea-c769-4d8e-a531-0458a304a5d2","text":"new complete todo","done":false}}}}

event: next
data: {"data":{"onTodoChange":{"action":"ToggleDone","oldValue":{"id":"ce1a1dea-c769-4d8e-a531-0458a304a5d2","text":"new complete todo","done":false},"newValue":{"id":"ce1a1dea-c769-4d8e-a531-0458a304a5d2","text":"new complete todo","done":true}}}}
```
