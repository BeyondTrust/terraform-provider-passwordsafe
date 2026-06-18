// Copyright 2025 BeyondTrust. All rights reserved.
package provider_framework

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestEffectiveVerifyCA verifies that effectiveVerifyCA maps a types.Bool
// VerifyCA attribute to the expected effective bool. Null and Unknown values
// must both default to true so provider initialization preserves the safe
// "verify the CA" default when the user omits the attribute or it has not
// yet been resolved.
func TestEffectiveVerifyCA(t *testing.T) {
	tests := []struct {
		name     string
		input    types.Bool
		expected bool
	}{
		{
			name:     "Null defaults to true",
			input:    types.BoolNull(),
			expected: true,
		},
		{
			name:     "Unknown defaults to true",
			input:    types.BoolUnknown(),
			expected: true,
		},
		{
			name:     "Explicit false stays false",
			input:    types.BoolValue(false),
			expected: false,
		},
		{
			name:     "Explicit true stays true",
			input:    types.BoolValue(true),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := effectiveVerifyCA(tt.input)
			if got != tt.expected {
				t.Errorf("effectiveVerifyCA(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
