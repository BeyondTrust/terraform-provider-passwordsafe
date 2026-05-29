// Copyright 2025 BeyondTrust. All rights reserved.
// Package utils.
package utils

import (
	"fmt"
	"sync"

	auth "github.com/BeyondTrust/go-client-library-passwordsafe/api/authentication"
	libentities "github.com/BeyondTrust/go-client-library-passwordsafe/api/entities"
)

var (
	sharedAuthMu       sync.Mutex
	sharedAuth         *auth.AuthenticationObj
	sharedSignAppinRsp libentities.SignAppinResponse
	sharedCacheKey     string
)

// InitSharedAuth lazily builds the shared AuthenticationObj on first call,
// performs the SignAppin handshake, and caches both for subsequent callers
// that pass the same cacheKey. Both muxed providers (framework + sdkv2) call
// this from their Configure; they receive identical provider-block config so
// they produce the same key and the second caller hits the cache.
//
// When cacheKey differs from the cached one (e.g., acceptance tests rotate
// mock-server URLs between cases) we sign the old session out and re-init.
// The signout against the prior (often dead) server is best-effort and any
// error is dropped — the cookie jar would be stale anyway.
//
// On init failure the cached state is cleared, so the next call retries
// rather than returning a half-built object.
func InitSharedAuth(cacheKey string, build func() (*auth.AuthenticationObj, error)) (*auth.AuthenticationObj, libentities.SignAppinResponse, error) {
	sharedAuthMu.Lock()
	defer sharedAuthMu.Unlock()

	if sharedAuth != nil && sharedCacheKey == cacheKey {
		return sharedAuth, sharedSignAppinRsp, nil
	}

	if sharedAuth != nil {
		_ = sharedAuth.SignOut()
	}
	sharedAuth = nil
	sharedSignAppinRsp = libentities.SignAppinResponse{}
	sharedCacheKey = ""

	authObj, err := build()
	if err != nil {
		return nil, libentities.SignAppinResponse{}, err
	}
	signAppin, err := authObj.GetPasswordSafeAuthentication()
	if err != nil {
		return nil, libentities.SignAppinResponse{}, err
	}

	sharedAuth = authObj
	sharedSignAppinRsp = signAppin
	sharedCacheKey = cacheKey
	return sharedAuth, sharedSignAppinRsp, nil
}

// ResetSharedAuthForTest clears the cached shared session without calling
// SignOut on the (often dead) test server. Acceptance tests rotate mock
// servers between cases; this lets them start each case from a clean state
// without paying for an HTTP round-trip against a closed listener.
//
// Test-only. Production code never calls this — production relies on the
// cacheKey path in InitSharedAuth to handle config changes.
func ResetSharedAuthForTest() {
	sharedAuthMu.Lock()
	defer sharedAuthMu.Unlock()
	sharedAuth = nil
	sharedSignAppinRsp = libentities.SignAppinResponse{}
	sharedCacheKey = ""
}

// ShutdownSharedAuth signs the shared session out of Password Safe and clears
// the cache so a later InitSharedAuth would rebuild from scratch. Safe to
// call multiple times — a no-op when no session is held.
//
// Called from main() after tf5server.Serve returns. Terraform may SIGKILL
// the plugin before Serve returns cleanly; in that case the server-side
// idle TTL handles cleanup.
func ShutdownSharedAuth() error {
	sharedAuthMu.Lock()
	defer sharedAuthMu.Unlock()

	if sharedAuth == nil {
		return nil
	}

	err := sharedAuth.SignOut()
	sharedAuth = nil
	sharedSignAppinRsp = libentities.SignAppinResponse{}
	sharedCacheKey = ""
	if err != nil {
		return fmt.Errorf("shared signout: %w", err)
	}
	return nil
}
