package configs

import "testing"

func TestIsWeakJWTSecret(t *testing.T) {
	if !isWeakJWTSecret("cilikube-secret-key-change-in-production") {
		t.Fatal("expected known default to be weak")
	}
	if !isWeakJWTSecret("short") {
		t.Fatal("expected short secret to be weak")
	}
	if isWeakJWTSecret("a-sufficiently-long-random-jwt-secret") {
		t.Fatal("expected strong secret to pass")
	}
}

func TestIsWeakEncryptionKey(t *testing.T) {
	if !isWeakEncryptionKey("change-me-in-production-32chars!!") {
		t.Fatal("expected known default to be weak")
	}
	if !isWeakEncryptionKey("too-short") {
		t.Fatal("expected wrong length to be weak")
	}
	if isWeakEncryptionKey("12345678901234567890123456789012") {
		t.Fatal("expected 32-char non-default key to pass")
	}
}
