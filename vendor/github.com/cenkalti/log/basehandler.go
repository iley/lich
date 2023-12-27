package log

type BaseHandler struct {
	Level     Level
	Formatter Formatter
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{
		Level:     DefaultLevel,
		Formatter: DefaultFormatter,
	}
}

func (h *BaseHandler) SetLevel(l Level) {
	h.Level = l
}

func (h *BaseHandler) SetFormatter(f Formatter) {
	h.Formatter = f
}

func (h *BaseHandler) FilterAndFormat(rec *Record) string {
	if rec.Level > h.Level {
		return ""
	}
	return h.Formatter.Format(rec)
}
