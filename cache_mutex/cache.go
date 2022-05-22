package cacheMutex

import (
	"sync"
)

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{}
}

type Memo struct {
	f Func
	sync.Mutex
	cache map[string]*entry
}

func New(f Func) *Memo {
	return &Memo{f:	f, cache: make(map[string]*entry)}
}

func (m *Memo) Get(key string) (value interface{}, err error) {
	m.Lock()
	e := m.cache[key]
	if e == nil {
		e = &entry{ready: make(chan struct{})}
		m.cache[key] = e
		m.Unlock()

		e.res.value, e.res.err = m.f(key)
		close(e.ready)
	} else {
		m.Unlock()
		<-e.ready
	}

	return e.res.value, e.res.err
}
