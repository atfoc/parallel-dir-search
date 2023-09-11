package unboundedqueue

import (
	"errors"
	"sync"
)

type UnboundedQueue struct {
	data []string
	cond *sync.Cond
}

func (q *UnboundedQueue) hasMore() bool {
	return len(q.data) > 0
}

func (q *UnboundedQueue) Pop() string {
	q.failIfNotInitialized()
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for !q.hasMore() {
		q.cond.Wait()
	}
	result := q.data[0]
	q.data = q.data[1:]

	return result
}

func (q *UnboundedQueue) Push(name string) {
	q.failIfNotInitialized()
	q.cond.L.Lock()
	q.data = append(q.data, name)
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

func (q *UnboundedQueue) PushAll(names []string) {
	q.failIfNotInitialized()
	q.cond.L.Lock()
	for _, el := range names {
		q.data = append(q.data, el)
	}
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

func (q *UnboundedQueue) failIfNotInitialized() {
	if q.cond == nil {
		panic(errors.New("struct UnboundedQueue not initialized"))
	}
}

func New() *UnboundedQueue {
	return &UnboundedQueue{cond: sync.NewCond(&sync.Mutex{})}
}
