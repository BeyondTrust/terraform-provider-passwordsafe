// Copyright 2025 BeyondTrust. All rights reserved.
package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestBuildProviderLogger_NoEnvVar_StderrOnly(t *testing.T) {
	const envVar = "TEST_PS_LOG_PATH_UNSET"
	_ = os.Unsetenv(envVar)

	logger := BuildProviderLogger(envVar)
	if logger == nil {
		t.Fatal("expected logger, got nil")
	}
	logger.Info("stderr-only")
	_ = logger.Sync()
}

func TestBuildProviderLogger_EnvVarSet_FileCreatedWith0600(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission semantics differ on Windows")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "provider.log")
	const envVar = "TEST_PS_LOG_PATH"
	t.Setenv(envVar, path)

	logger := BuildProviderLogger(envVar)
	if logger == nil {
		t.Fatal("expected logger, got nil")
	}
	logger.Info("file-sink")
	_ = logger.Sync()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("log file not created: %v", err)
	}
	if mode := info.Mode().Perm(); mode != 0o600 {
		t.Errorf("expected file mode 0600, got %o", mode)
	}
}

func TestBuildProviderLogger_ExistingFile_PermsTightened(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission semantics differ on Windows")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "existing.log")
	if err := os.WriteFile(path, []byte("preexisting\n"), 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}
	const envVar = "TEST_PS_LOG_PATH_EXISTING"
	t.Setenv(envVar, path)

	logger := BuildProviderLogger(envVar)
	if logger == nil {
		t.Fatal("expected logger, got nil")
	}
	logger.Info("existing-file")
	_ = logger.Sync()

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if mode := info.Mode().Perm(); mode != 0o600 {
		t.Errorf("expected file mode tightened to 0600, got %o", mode)
	}
}

func TestBuildProviderLogger_BadPath_FallsBackWithoutPanic(t *testing.T) {
	const envVar = "TEST_PS_LOG_PATH_BAD"
	t.Setenv(envVar, filepath.Join(string(os.PathSeparator), "this", "path", "does", "not", "exist", "provider.log"))

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()

	logger := BuildProviderLogger(envVar)
	if logger == nil {
		t.Fatal("expected logger, got nil")
	}
	logger.Info("fallback-ok")
	_ = logger.Sync()
}
