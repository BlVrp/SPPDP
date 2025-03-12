package zaplog

import (
	"fmt"

	"go.uber.org/zap"

	"one-help/internal/logger"
)

// ensures that zaplog implements logger.Logger.
var _ logger.Logger = (*zaplog)(nil)

// zaplog is an implementation of logger.Logger using zap.
type zaplog struct {
	client *zap.Logger
}

// NewLog is a constructor for a logger.Logger.
func NewLog() logger.Logger {
	// TODO: Make it configurable.
	return &zaplog{
		client: zap.NewExample(),
	}
}

// Info is used to send info message.
func (log zaplog) Info(msg string) {
	log.client.Info(msg)
}

// InfoF is used to send formatted info message.
// INFO: Uses fmt.Sprintf under the hood.
func (log zaplog) InfoF(msg string, args ...any) {
	log.client.Info(fmt.Sprintf(msg, args...))
}

// Debug is used to send debug message.
func (log zaplog) Debug(msg string) {
	log.client.Debug(msg)
}

// DebugF is used to send formatted debug message.
// INFO: Uses fmt.Sprintf under the hood.
func (log zaplog) DebugF(msg string, args ...any) {
	log.client.Debug(fmt.Sprintf(msg, args...))
}

// Warn is used to send warn message.
func (log zaplog) Warn(msg string) {
	log.client.Warn(msg)
}

// WarnF is used to send formatted warn message.
// INFO: Uses fmt.Sprintf under the hood.
func (log zaplog) WarnF(msg string, args ...any) {
	log.client.Warn(fmt.Sprintf(msg, args...))
}

// Error is used to send error message.
func (log zaplog) Error(msg string, err error) {
	log.client.Error(msg, zap.Error(err))
}

// ErrorF is used to send formatted message about error.
// INFO: Uses fmt.Sprintf under the hood.
func (log zaplog) ErrorF(msg string, err error, args ...any) {
	log.client.Error(fmt.Sprintf(msg, args...), zap.Error(err))
}
