package log

import (
	"sync"

	"github.com/hashicorp/go-multierror"
)

// MultiHandler sends the log output to multiple handlers concurrently.
type MultiHandler struct {
	handlers []Handler
}

func NewMultiHandler(handlers ...Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (b *MultiHandler) SetFormatter(f Formatter) {
	for _, h := range b.handlers {
		h.SetFormatter(f)
	}
}

func (b *MultiHandler) SetLevel(l Level) {
	for _, h := range b.handlers {
		h.SetLevel(l)
	}
}

func (b *MultiHandler) Handle(rec *Record) {
	wg := sync.WaitGroup{}
	wg.Add(len(b.handlers))
	for _, handler := range b.handlers {
		go func(handler Handler) {
			handler.Handle(rec)
			wg.Done()
		}(handler)
	}
	wg.Wait()
}

func (b *MultiHandler) Close() error {
	var result error
	var m sync.Mutex
	wg := sync.WaitGroup{}
	wg.Add(len(b.handlers))
	for _, handler := range b.handlers {
		go func(handler Handler) {
			err := handler.Close()
			m.Lock()
			if err != nil {
				result = multierror.Append(result, err)
			}
			m.Unlock()
			wg.Done()
		}(handler)
	}
	wg.Wait()
	return result
}
