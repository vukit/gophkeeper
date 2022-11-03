package logger

import (
	"io"

	"github.com/rs/zerolog"
)

// Logger структура журнала сервера
type Logger struct {
	logger zerolog.Logger
}

// NewLogger возвращает журнал для сервера с выводом в io.Writer
func NewLogger(w io.Writer) *Logger {
	r := &Logger{}
	r.logger = zerolog.New(w).With().Timestamp().Logger()

	return r
}

// Fatal сохраняет сообщение в журнал с уровнем fatal
func (r *Logger) Fatal(message string) {
	r.logger.Fatal().Interface("message", message).Send()
}

// Warning сохраняет сообщение в журнал с уровнем warning
func (r *Logger) Warning(message string) {
	r.logger.Warn().Interface("message", message).Send()
}

// Info сохраняет сообщение в журнал с уровнем info
func (r *Logger) Info(message string) {
	r.logger.Info().Interface("message", message).Send()
}
