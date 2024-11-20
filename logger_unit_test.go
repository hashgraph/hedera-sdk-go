//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("test", LoggerLevelDebug)
	require.NotNil(t, logger)
	assert.Equal(t, LoggerLevelDebug, logger.level)
}

func TestLogger_SetLevel(t *testing.T) {
	logger := NewLogger("test", LoggerLevelDebug)
	logger.SetLevel(LoggerLevelInfo)
	assert.Equal(t, LoggerLevelInfo, logger.level)
}
func TestLogger_SetSilent(t *testing.T) {
	var buf bytes.Buffer
	writer := zerolog.ConsoleWriter{Out: &buf, TimeFormat: time.RFC3339}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	l := NewLogger("test", LoggerLevelTrace)
	l.SetLevel(LoggerLevelTrace)

	logger := zerolog.New(&writer)
	l.logger = &logger

	l.Debug("debug message", "debugKey", "debugValue")
	l.Info("info message", "infoKey", "infoValue")
	l.Warn("warn message", "warnKey", "warnValue")
	l.Error("error message", "errorKey", "errorValue")
	l.Trace("trace message", "traceKey", "traceValue")
	assert.Contains(t, buf.String(), "debug message")
	assert.Contains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "error message")
	assert.Contains(t, buf.String(), "trace message")

	// coverage only
	l.SetSilent(true)
	l.SetSilent(false)
}

func TestNewLoggerWithEnvironmentVariableSet(t *testing.T) {
	os.Setenv("HEDERA_SDK_GO_LOG_PRETTY", "1")
	logger := NewLogger("test", LoggerLevelDebug)
	require.NotNil(t, logger)
	assert.Equal(t, LoggerLevelDebug, logger.level)
	os.Unsetenv("HEDERA_SDK_GO_LOG_PRETTY")
}
