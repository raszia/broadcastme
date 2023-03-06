# broadcastme

broadcastme is a Go package that provides a simple way to broadcast messages to multiple Go routines that are waiting to receive them.

The use of Go Generics in this package allows for the creation of a type-safe broadcast server that can handle any type of data, and provides flexibility to the client code. The package can be imported and used by other Go programs to easily implement a channel broadcasting feature in their applications.

## Features

+ Subscribe to a broadcast using a key.

+ Unsubscribe from a broadcast.

+ Add new broadcasts on the fly.

+ Supports Go's context.Context for graceful shutdown.

## Installation

To install, simply run:

```bash
go get github.com/raszia/broadcastme
```

## Usage

Import the package in your Go code with:

```go
import "github.com/raszia/broadcastme"
```

## Creating a Broadcast

```go
ch := make(chan string)
broadcastKey := "key1"
broadCastVal := "testVal"

broadcast := broadcastme.NewBroadcast(ch, broadCastKeyTest)
```

## Creating a Broadcast Server

Create a new broadcast server by calling `NewBroadcastServer()` with a context and a new broadcast. For example:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

server := broadcastme.NewBroadcastServerWithContext(ctx, broadcast)
```

## Subscribing to a Broadcast

Subscribe to a broadcast by calling `Subscribe()` on the server with a key. This will return a channel that can be used to listen for new broadcast messages. For example:

```go
channel := server.Subscribe(broadcastKey)
```

## Unsubscribing from a Broadcast

Unsubscribe from a broadcast by calling `Unsubscribe()` on the server with the channel that was returned by `Subscribe()`. For example:

```go
server.Unsubscribe(channel)
```

## Adding a New Broadcast to a created broadcast server

Add a new broadcast on the fly by calling `AddNewBroadcast()` on the server with a new broadcast. For example:

```go
ch2 := make(chan string)
server.AddNewBroadcast(broadcastme.NewBroadcast(ch2, "newkey"))
```

[examples](https://github.com/raszia/broadcastme/tree/main/example).

## License

Broadcastme is released under the MIT License. See `LICENSE` for more information.
