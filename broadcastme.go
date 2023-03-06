package broadcastme

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type BroadcastServer[T any, k comparable] interface {
	Subscribe(k) *channel[T, k]
	Unsubscribe(*channel[T, k])
	AddNewBroadcastWithContext(context.Context, *broadcast[T, k])
	AddNewBroadcast(*broadcast[T, k])
}

type broadcast[T any, k comparable] struct {
	source <-chan T
	key    k
}

type broadcastServer[T any, k comparable] struct {
	muRw      sync.RWMutex
	listeners map[k]*listener[T, k]
}

type listener[T any, k comparable] struct {
	muRw    sync.RWMutex
	key     k
	channel map[uuid.UUID]*channel[T, k]
}

type channel[T any, k comparable] struct {
	key     k
	channel chan T
	uuid    uuid.UUID
}

func newChannel[T any, k comparable](key k) *channel[T, k] {
	return &channel[T, k]{
		key:     key,
		channel: make(chan T),
		uuid:    uuid.New(),
	}
}

// Listen listen to a broadcast message with this method and get the messages.
func (ch *channel[T, k]) Listen() <-chan T {
	return ch.channel
}

// NewBroadcast creates a new broadcast message.
// You can send your message later on the source channel for those who have subscribed to this key.
func NewBroadcast[T any, k comparable](source <-chan T, key k) *broadcast[T, k] {
	return &broadcast[T, k]{
		source: source,
		key:    key,
	}
}

// NewBroadcastServerWithContext creates a new BroadcastServer with the given context and broadcast.
func NewBroadcastServerWithContext[T any, k comparable](ctx context.Context, b *broadcast[T, k]) BroadcastServer[T, k] {
	service := &broadcastServer[T, k]{
		muRw:      sync.RWMutex{},
		listeners: make(map[k]*listener[T, k]),
	}
	go service.serve(ctx, b)
	return service
}

// NewBroadcastServer creates a new BroadcastServer with the given broadcast.
// It is recommended to use NewBroadcastServerWithContext to create a BroadcastServer with a context.
func NewBroadcastServer[T any, k comparable](ctx context.Context, b *broadcast[T, k]) BroadcastServer[T, k] {
	return NewBroadcastServerWithContext(context.Background(), b)
}

// AddNewBroadcastWithContext add new broadcast to existing server with the given context and broadcast.
func (s *broadcastServer[T, k]) AddNewBroadcastWithContext(ctx context.Context, b *broadcast[T, k]) {
	go s.serve(ctx, b)
}

// AddNewBroadcastWithContext add new broadcast to existing server with broadcast.
// It is recommended to use
// AddNewBroadcastWithContext()
func (s *broadcastServer[T, k]) AddNewBroadcast(b *broadcast[T, k]) {
	s.AddNewBroadcastWithContext(context.Background(), b)
}

// Subscribe use this method to subscibe to a broadcast and get the channel
func (s *broadcastServer[T, k]) Subscribe(key k) *channel[T, k] {
	s.muRw.Lock()
	defer s.muRw.Unlock()

	newChan := newChannel[T](key)
	l, ok := s.listeners[key]
	if !ok {
		l = &listener[T, k]{
			muRw:    sync.RWMutex{},
			key:     key,
			channel: make(map[uuid.UUID]*channel[T, k]),
		}

	}
	s.listeners[key] = l
	l.addChannel(newChan)
	return newChan

}

// AddChannel adds a new channel to the listener.
func (l *listener[T, k]) addChannel(ch *channel[T, k]) {
	l.muRw.Lock()
	defer l.muRw.Unlock()
	l.channel[ch.uuid] = ch
}

// RemoveChannel removes a channel from the listener.
func (l *listener[T, k]) removeChannel(channelID uuid.UUID) {
	l.muRw.Lock()
	defer l.muRw.Unlock()
	delete(l.channel, channelID)
}

// Unsubscribe unsubscribe a channel.
func (s *broadcastServer[T, k]) Unsubscribe(channel *channel[T, k]) {
	s.muRw.RLock()
	defer s.muRw.RUnlock()
	s.listeners[channel.key].removeChannel(channel.uuid)
	close(channel.channel)
}

func (s *broadcastServer[T, k]) serve(ctx context.Context, b *broadcast[T, k]) {
	defer func() {

		listner, ok := s.listeners[b.key]
		if !ok { //no Subscrib yet
			return
		}
		for _, channel := range listner.channel {
			s.Unsubscribe(channel)
		}
		s.muRw.Lock()
		defer s.muRw.Unlock()
		delete(s.listeners, b.key)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-b.source:
			if !ok {
				return
			}
			s.muRw.RLock()
			listner, ok := s.listeners[b.key]
			if !ok { //no Subscribtion yet
				s.muRw.RUnlock()
				continue
			}
			s.muRw.RUnlock()
			listner.muRw.RLock()
			for _, ch := range listner.channel {
				select {
				case ch.channel <- val:
				case <-ctx.Done():
					listner.muRw.RUnlock()
					return
				default:
					continue
				}

			}
			listner.muRw.RUnlock()
		}

	}
}
