package log

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mattn/go-isatty"
)

// FileHandler is a handler implementation that writes the logging output to a *os.File.
// If given file is a tty, output will be colored.
type FileHandler struct {
	*BaseHandler
	f      *os.File
	m      sync.Mutex
	isatty bool
}

func NewFileHandler(f *os.File) *FileHandler {
	return &FileHandler{
		BaseHandler: NewBaseHandler(),
		f:           f,
		isatty:      isatty.IsTerminal(f.Fd()),
	}
}

func (h *FileHandler) Handle(rec *Record) {
	message := h.BaseHandler.FilterAndFormat(rec)
	if message == "" {
		return
	}
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	if h.isatty && LevelColors[rec.Level] != NOCOLOR {
		message = fmt.Sprintf("\033[%dm%s\033[0m", LevelColors[rec.Level], message)
	}
	h.m.Lock()
	fmt.Fprint(h.f, message)
	h.m.Unlock()
}

func (h *FileHandler) Close() error {
	return h.Close()
}

type Color int

// Colors for different log levels.
const (
	BLACK Color = (iota + 30)
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
	NOCOLOR = -1
)

var LevelColors = map[Level]Color{
	CRITICAL: MAGENTA,
	ERROR:    RED,
	WARNING:  YELLOW,
	NOTICE:   GREEN,
	INFO:     NOCOLOR,
	DEBUG:    BLUE,
}
