package log

import "time"

// Record contains all of the information about a single log message.
type Record struct {
	Message     string    // Formatted log message
	LoggerName  string    // Name of the logger module
	Level       Level     // Level of the record
	Time        time.Time // Time of the record (local time)
	Filename    string    // File name of the log call (absolute path)
	Line        int       // Line number in file
	ProcessID   int       // PID
	ProcessName string    // Name of the process
}
