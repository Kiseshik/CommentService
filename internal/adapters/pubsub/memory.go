package pubsub

import (
	"sync"

	"github.com/Kiseshik/CommentService.git/internal/core/domain"
)

type MemoryPubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan *domain.Comment
}

func NewInMemoryPubSub() *MemoryPubSub {
	return &MemoryPubSub{
		subscribers: make(map[string][]chan *domain.Comment),
	}
}

func (ps *MemoryPubSub) Subscribe(postID string) (<-chan *domain.Comment, func()) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan *domain.Comment, 10)
	ps.subscribers[postID] = append(ps.subscribers[postID], ch)
	unsubscribe := func() {
		ps.mu.Lock()
		defer ps.mu.Unlock()
		subs := ps.subscribers[postID]
		for i, sub := range subs {
			if sub == ch {
				ps.subscribers[postID] = append(subs[:i], subs[i+1:]...)
				close(ch)
				break
			}
		}
	}
	return ch, unsubscribe
}

func (ps *MemoryPubSub) Publish(postID string, comment *domain.Comment) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for _, ch := range ps.subscribers[postID] {
		select {
		case ch <- comment:
		default:
		}
	}
}
