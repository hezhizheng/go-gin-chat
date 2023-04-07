package safe

import "sync"

var Safety ThreadSafety

type ThreadSafety struct {
	mu sync.Mutex
}

func (receiver *ThreadSafety) Do(x func() interface{}) interface{} {
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	return x()
}
