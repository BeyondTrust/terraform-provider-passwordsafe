// Copyright 2025 BeyondTrust. All rights reserved.
// Package utils.
package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// BuildProviderLogger constructs a zap logger for a provider half.
//
// By default the logger writes to stderr only. Persisting debug-level client
// logs to a predictable, broadly-readable file in the current working directory
// is a security risk for a secrets-handling provider, so the on-disk sink is
// opt-in: when the environment variable named by logFileEnvVar holds a non-empty
// path, the logger additionally appends to that file. The file is opened (or
// created) with 0600 permissions so the debug output is not group/world-readable.
// If the file cannot be opened the logger falls back to stderr-only rather than
// failing provider startup.
func BuildProviderLogger(logFileEnvVar string) *zap.Logger {
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	level := zap.NewAtomicLevelAt(zap.DebugLevel)

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), level),
	}

	if path := os.Getenv(logFileEnvVar); path != "" {
		if file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600); err == nil {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(file), level))
		}
	}

	return zap.New(zapcore.NewTee(cores...))
}
