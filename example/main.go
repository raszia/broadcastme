package main

import (
	"broadcastme"
	"context"
	"log"
	"sync"
	"time"
)

func main() {
	ch := make(chan string)
	broadcastKey := "key1"
	broadCastVal := "testVal"

	broadcast := broadcastme.NewBroadcast(ch, broadcastKey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := broadcastme.NewBroadcastServerWithContext(ctx, broadcast)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		log.Println("listener1 is waiting...")
		listener1 := server.Subscribe(broadcastKey)
		msg1 := <-listener1.Listen()
		log.Println("listener1 got:", msg1)
	}()

	go func() {
		defer wg.Done()
		log.Println("listener2 is waiting...")
		listener2 := server.Subscribe(broadcastKey)
		msg2 := <-listener2.Listen()
		log.Println("listener2 got:", msg2)
	}()

	time.Sleep(time.Second * 3)

	log.Println("sending data to broadcast channel:", broadCastVal)
	ch <- broadCastVal

	wg.Wait()
}
