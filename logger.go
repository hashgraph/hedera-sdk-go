package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type LogLevel string

const (
	LoggerLevelTrace    LogLevel = "TRACE"
	LoggerLevelDebug    LogLevel = "DEBUG"
	LoggerLevelInfo     LogLevel = "INFO"
	LoggerLevelWarn     LogLevel = "WARN"
	LoggerLevelError    LogLevel = "ERROR"
	LoggerLevelDisabled LogLevel = "DISABLED"
)

type Logger interface {
	SetSilent(isSilent bool)
	SetLevel(level LogLevel)
	SubLoggerWithLevel(level LogLevel) Logger
	Error(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Trace(msg string, keysAndValues ...interface{})
}

type DefaultLogger struct {
	logger *zerolog.Logger
	level  LogLevel
}

func NewLogger(component string, level LogLevel) *DefaultLogger {
	var logger zerolog.Logger
	logger = zerolog.New(os.Stdout).With().Str("module", component).Timestamp().Logger()

	if os.Getenv("HEDERA_SDK_GO_LOG_PRETTY") != "" {
		// Pretty logging
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(output).With().Str("module", component).Timestamp().Logger()
	}

	logger = loggerForLevel(logger, level)
	return &DefaultLogger{
		logger: &logger,
		level:  level,
	}
}

func (l *DefaultLogger) SetSilent(isSilent bool) {
	if isSilent {
		logger := l.logger.Level(zerolog.Disabled)
		l.logger = &logger
	} else {
		l.SetLevel(l.level)
	}
}

func loggerForLevel(logger zerolog.Logger, level LogLevel) zerolog.Logger {
	switch level {
	case LoggerLevelTrace:
		return logger.Level(zerolog.TraceLevel)
	case LoggerLevelDebug:
		return logger.Level(zerolog.DebugLevel)
	case LoggerLevelInfo:
		return logger.Level(zerolog.InfoLevel)
	case LoggerLevelWarn:
		return logger.Level(zerolog.WarnLevel)
	case LoggerLevelError:
		return logger.Level(zerolog.ErrorLevel)
	default:
		return logger.Level(zerolog.Disabled)
	}
}

func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.level = level
	logger := loggerForLevel(*l.logger, level)
	l.logger = &logger
}

func (l *DefaultLogger) SubLoggerWithLevel(level LogLevel) Logger {
	l.level = level
	logger := loggerForLevel(*l.logger, level)

	return &DefaultLogger{
		logger: &logger,
		level:  level,
	}
}

func (l *DefaultLogger) Warn(msg string, keysAndValues ...interface{}) {
	addFields(l.logger.Warn(), msg, keysAndValues...)
}

func (l *DefaultLogger) Error(msg string, keysAndValues ...interface{}) {
	addFields(l.logger.Error(), msg, keysAndValues...)
}

func (l *DefaultLogger) Trace(msg string, keysAndValues ...interface{}) {
	addFields(l.logger.Trace(), msg, keysAndValues...)
}

func (l *DefaultLogger) Debug(msg string, keysAndValues ...interface{}) {
	addFields(l.logger.Debug(), msg, keysAndValues...)
}

func (l *DefaultLogger) Info(msg string, keysAndValues ...interface{}) {
	addFields(l.logger.Info(), msg, keysAndValues...)
}

func argsToFields(keysAndValues ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	for i := 0; i < len(keysAndValues); i += 2 {
		fields[fmt.Sprint(keysAndValues[i])] = keysAndValues[i+1]
	}
	return fields
}

func addFields(event *zerolog.Event, msg string, keysAndValues ...interface{}) {
	for key, value := range argsToFields(keysAndValues...) {
		switch v := value.(type) {
		case string:
			event.Str(key, v)
		case int64:
			event.Int64(key, v)
		case time.Duration:
			event.Dur(key, v)
		default:
			event.Str(key, fmt.Sprint(value))
		}
	}

	event.Msg(msg)
}
