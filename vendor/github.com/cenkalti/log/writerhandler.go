package log

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

// WriterHandler is a handler implementation that writes the logging output to a io.Writer.
type WriterHandler struct {
	*BaseHandler
	w io.Writer
	m sync.Mutex
}

func NewWriterHandler(w io.Writer) *WriterHandler {
	return &WriterHandler{
		BaseHandler: NewBaseHandler(),
		w:           w,
	}
}

func (b *WriterHandler) Handle(rec *Record) {
	message := b.BaseHandler.FilterAndFormat(rec)
	if message == "" {
		return
	}
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	b.m.Lock()
	fmt.Fprint(b.w, message)
	b.m.Unlock()
}

func (b *WriterHandler) Close() error {
	if c, ok := b.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
