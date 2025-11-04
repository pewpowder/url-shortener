package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
	1. Broadcaster is a struct that allow subscribe on some events, unsubscribe and publish some messages.
	2. The main idea is that we can add any number of channels and if we publish some message all of these channels will get it
	3. We don't save history of messages
*/

type broadcaster[T any] struct {
	mu   sync.RWMutex
	subs map[chan T]struct{}
}

func newBroadcaster[T any](subsSize int) *broadcaster[T] {
	return &broadcaster[T]{
		subs: make(map[chan T]struct{}, subsSize),
	}
}

func (b *broadcaster[T]) Subscribe(buffSize int) (<-chan T, func() bool) {
	ch := make(chan T, buffSize)

	b.subs[ch] = struct{}{}

	unsubscribe := func() bool {
		if _, ok := b.subs[ch]; ok {
			delete(b.subs, ch)
			close(ch)
			return true
		}

		return false
	}

	return ch, unsubscribe
}

/*
PS. Due to we cast our channels from chan T to <-chan T we can't find a chan for O(1) time in map.
So there are three ways to solve this problem:
1. Don't cast channels. But it's not a good approache because we allow our subscribers to write messages in chan
2. The second approach is looking for channel in array and narrow type before checking. But this aprroach also has downside.
Now we have searching time O(N)
3. The third approach is returning unsubscribe function from Subscribe method. This approach saves const searching time and we still have channels only for reading
*/

// func (b *broadcaster[T]) Unsubscribe(sub <-chan T) bool {
// 	b.mu.Lock()
// 	defer b.mu.Unlock()
// 	for ch := range b.subs {
// 		if (<-chan T)(ch) == sub {
// 			delete(b.subs, ch)
// 			close(ch)
// 			break
// 		}
// 	}
// 	return false
// }

func (b *broadcaster[T]) PrintSubsCount() {
	fmt.Println("Subs count is ", len(b.subs))
}

func (b *broadcaster[T]) Publish(message T) {
	for sub := range b.subs {
		sub <- message
	}
}

func RunPubSub() {
	b := newBroadcaster[int](10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s1, unsubscribeS1 := b.Subscribe(1)
	s2, unsubscribeS2 := b.Subscribe(1)

	go func() {
		defer unsubscribeS1()
		for {
			select {
			case v, ok := <-s1:
				if !ok {
					return
				}
				fmt.Println("s1:", v)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		defer unsubscribeS2()
		for {
			select {
			case v, ok := <-s2:
				if !ok {
					return
				}
				fmt.Println("s2:", v)
			case <-ctx.Done():
				return
			}
		}
	}()

	b.Publish(1)
	b.Publish(2)
	b.Publish(3)

	b.Publish(4)

	time.Sleep(time.Second * 6)
}
