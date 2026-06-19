// Copyright 2025 BeyondTrust. All rights reserved.
package utils

import "testing"

// FuzzValidateChangeFrequencyDays exercises ValidateChangeFrequencyDays with
// arbitrary input and asserts its documented invariants never break (and that
// it never panics).
func FuzzValidateChangeFrequencyDays(f *testing.F) {
	f.Add("xdays", 1)
	f.Add("xdays", 999)
	f.Add("xdays", 0)
	f.Add("xdays", 1000)
	f.Add("xdays", -5)
	f.Add("other", 0)
	f.Add("", 50)

	f.Fuzz(func(t *testing.T, changeFrequencyType string, changeFrequencyDays int) {
		err := ValidateChangeFrequencyDays(changeFrequencyType, changeFrequencyDays)

		if changeFrequencyType != "xdays" {
			if err != nil {
				t.Errorf("expected no error for type %q, got %v", changeFrequencyType, err)
			}
			return
		}

		inRange := changeFrequencyDays >= 1 && changeFrequencyDays <= 999
		if inRange && err != nil {
			t.Errorf("expected no error for xdays=%d (in range), got %v", changeFrequencyDays, err)
		}
		if !inRange && err == nil {
			t.Errorf("expected error for xdays=%d (out of range), got nil", changeFrequencyDays)
		}
	})
}
