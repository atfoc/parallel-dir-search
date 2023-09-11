package unboundedqueue

import "sync"

type UnboundedQueue struct {
	data []string
	cond *sync.Cond
}

func (q *UnboundedQueue) hasMore() bool {
	return len(q.data) > 0
}

func (q *UnboundedQueue) Pop() string {
	q.init()
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
	q.init()
	q.cond.L.Lock()
	q.data = append(q.data, name)
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

func (q *UnboundedQueue) PushAll(names []string) {
	q.init()
	q.cond.L.Lock()
	for _, el := range names {
		q.data = append(q.data, el)
	}
	q.cond.Broadcast()
	q.cond.L.Unlock()
}

func (q *UnboundedQueue) init() {
	if q.cond != nil {
		return
	}

	q.cond = sync.NewCond(&sync.Mutex{})
}
