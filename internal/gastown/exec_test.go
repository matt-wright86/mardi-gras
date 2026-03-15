package gastown

import (
	"testing"
	"time"
)

func TestSetCmdTimeoutDoubles(t *testing.T) {
	// Reset after test
	defer func() {
		timeoutLong = defaultTimeoutLong
		timeoutMedium = defaultTimeoutMedium
		timeoutShort = defaultTimeoutShort
	}()

	SetCmdTimeout(60) // double the 30s baseline
	if timeoutLong != 60*time.Second {
		t.Errorf("timeoutLong = %v, want 60s", timeoutLong)
	}
	if timeoutMedium != 30*time.Second {
		t.Errorf("timeoutMedium = %v, want 30s", timeoutMedium)
	}
	if timeoutShort != 10*time.Second {
		t.Errorf("timeoutShort = %v, want 10s", timeoutShort)
	}
}

func TestSetCmdTimeoutIgnoresZero(t *testing.T) {
	defer func() {
		timeoutLong = defaultTimeoutLong
		timeoutMedium = defaultTimeoutMedium
		timeoutShort = defaultTimeoutShort
	}()

	SetCmdTimeout(0)
	if timeoutLong != defaultTimeoutLong {
		t.Errorf("timeoutLong changed on zero input: %v", timeoutLong)
	}
}

func TestSetCmdTimeoutIgnoresNegative(t *testing.T) {
	defer func() {
		timeoutLong = defaultTimeoutLong
		timeoutMedium = defaultTimeoutMedium
		timeoutShort = defaultTimeoutShort
	}()

	SetCmdTimeout(-5)
	if timeoutLong != defaultTimeoutLong {
		t.Errorf("timeoutLong changed on negative input: %v", timeoutLong)
	}
}
