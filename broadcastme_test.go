package broadcastme_test

import (
	"broadcastme"
	"context"
	"sync"
	"testing"
	"time"
)

func TestBroadcastServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan string)
	broadCastKeyTest := "key1"
	broadCastVal1 := "testval1"

	broadCast := broadcastme.NewBroadcast(ch, broadCastKeyTest)

	broadcastServer := broadcastme.NewBroadcastServerWithContext(ctx, broadCast)

	listener1 := broadcastServer.Subscribe(broadCastKeyTest)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		name := "listener1"
		t.Logf("%s wait for recive...", name)
		msg1 := <-listener1.Listen()
		t.Logf("%s recived %v", name, msg1)
		if msg1 != broadCastVal1 {
			t.Errorf("%s expected %v got %v", name, broadCastVal1, msg1)
		}
	}()

	go func() {
		defer wg.Done()
		name := "listener2"
		listener2 := broadcastServer.Subscribe(broadCastKeyTest)
		t.Logf("%s wait for recive...", name)
		msg2 := <-listener2.Listen()
		t.Logf("%s recived %v", name, msg2)
		if msg2 != broadCastVal1 {
			t.Errorf("%s expected %v got %v", name, broadCastVal1, msg2)
		}
	}()

	time.Sleep(time.Millisecond * 10)
	t.Logf("sending %v", broadCastVal1)
	ch <- broadCastVal1

	go func() {
		defer wg.Done()
		name := "noReciver"
		t.Logf("%s wait for recive...", name)

		c := time.NewTicker(time.Millisecond * 20)
		select {
		case msg := <-listener1.Listen():
			if msg != "" {
				t.Errorf("%s should NOT recive anyting but got %v", name, msg)
			}
		case <-c.C:
		}
		t.Logf("%s didn't recive anyting it's ok", name)
	}()

	wg.Wait()

	broadcastServer.Unsubscribe(listener1)
	ch <- broadCastVal1

	msg4 := <-listener1.Listen()
	if msg4 != "" {
		t.Errorf("closed channel should not get somthing but it got %v", msg4)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		name := "listener with cancled context"
		listener2 := broadcastServer.Subscribe(broadCastKeyTest)
		<-listener2.Listen()
		msg2 := <-listener2.Listen()
		if msg2 != "" {
			t.Errorf("context is cancled %s should not get anyting but got %v", name, msg2)
		}
	}()

	cancel()
	wg.Wait()

	ch2 := make(chan string)
	broadCastKeyTest2 := "key2"
	broadCastVal2 := "testval2"

	broadCast2 := broadcastme.NewBroadcast(ch2, broadCastKeyTest2)
	broadcastServer.AddNewBroadcast(broadCast2)
	listener := broadcastServer.Subscribe(broadCastKeyTest2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		newMsg := <-listener.Listen()
		if newMsg != broadCastVal2 {
			t.Errorf("listener expected %v got %v", broadCastVal2, newMsg)
		}
	}()
	ch2 <- broadCastVal2
	wg.Wait()
}
