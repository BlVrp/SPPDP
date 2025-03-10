package logger

// Logger exposes functionality to write messages in stdout.
type Logger interface {
	// Info is used to send info message.
	Info(msg string)
	// Debug is used to send formatted debug message.
	Debug(msg string)
	// Warn is used to send formatted warning message.
	Warn(msg string)
	// Error is used to send formatted as error message.
	Error(msg string, err error)

	// InfoF is used to send formatted info message.
	InfoF(msg string, args ...any)
	// DebugF is used to send formatted debug message.
	DebugF(msg string, args ...any)
	// WarnF is used to send formatted warn message.
	WarnF(msg string, args ...any)
	// ErrorF is used to send formatted message about error.
	ErrorF(msg string, err error, args ...any)
}
