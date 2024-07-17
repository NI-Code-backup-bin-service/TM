package dal

import (
	"testing"
)

func TestFormatLimitIdentifier(t *testing.T) {
	tests := []struct {
		name string
		limitType string
		expected string
	}{
		{"No change required", "NoChangeRequired", "nochangerequired"},
		{"Replace spaces with hyphens", "Replace spaces with hyphens", "replace-spaces-with-hyphens"},
		{"Ensure numbers and brackets are kept", "Numbers 1234 (and brackets)", "numbers-1234-(and-brackets)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatLimitIdentifier(tt.limitType)
			if err != nil {
				t.Errorf(err.Error())
			}
			if got != tt.expected {
				t.Errorf("formatLimitIdentifier(), expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}
