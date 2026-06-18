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

	writers := []zapcore.WriteSyncer{zapcore.Lock(os.Stderr)}

	if path := os.Getenv(logFileEnvVar); path != "" {
		if file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600); err == nil {
			// OpenFile only applies perms on create; tighten existing files too via
			// the file descriptor (avoids the TOCTOU/symlink race of a path-based chmod).
			// If chmod fails, the file may still be group/world-readable, so close it
			// and fall back to stderr-only rather than leak debug output.
			if chmodErr := file.Chmod(0o600); chmodErr != nil {
				_ = file.Close()
			} else {
				writers = append(writers, zapcore.Lock(zapcore.AddSync(file)))
			}
		}
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), level)
	// Match the diagnostic context previously provided by zap.Config.Build():
	// caller info, stacktraces at error level, and error output to stderr.
	return zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.ErrorOutput(zapcore.Lock(os.Stderr)),
	)
}
