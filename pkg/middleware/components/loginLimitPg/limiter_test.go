package loginLimitPg

import "testing"

func TestLockAfterFailures(t *testing.T) {
	key := "test-lock-key"
	Reset(key)
	for i := 0; i < maxFailures; i++ {
		Failure(key)
	}
	allowed, wait := Allowed(key)
	if allowed {
		t.Fatalf("expected locked after %d failures", maxFailures)
	}
	if wait <= 0 {
		t.Fatalf("expected positive wait, got %d", wait)
	}
	Reset(key)
	if ok, _ := Allowed(key); !ok {
		t.Fatal("reset should unlock")
	}
}

func TestAllowedByDefault(t *testing.T) {
	if ok, _ := Allowed("unknown-key"); !ok {
		t.Fatal("unknown key must be allowed")
	}
}
