package log

// Handler handles the output.
type Handler interface {
	SetFormatter(Formatter)
	SetLevel(Level)
	// Handle single log record.
	Handle(*Record)
	// Close the handler.
	Close() error
}
