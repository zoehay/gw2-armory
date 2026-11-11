package routes

import (
	"os"
	"testing"
)

func TestCookieSecureFromEnv(t *testing.T) {
	tests := []struct {
		name  string
		value *string
		want  bool
	}{
		{"unset defaults to secure", nil, true},
		{"explicit false", ptr("false"), false},
		{"explicit true", ptr("true"), true},
		{"unparseable defaults to secure", ptr("nope"), true},
		{"empty defaults to secure", ptr(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setenv restores the original value after the test, including for the unset case
			t.Setenv("COOKIE_SECURE", "")
			if tt.value != nil {
				t.Setenv("COOKIE_SECURE", *tt.value)
			} else {
				os.Unsetenv("COOKIE_SECURE")
			}
			if got := cookieSecureFromEnv(); got != tt.want {
				t.Errorf("cookieSecureFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ptr(s string) *string { return &s }
