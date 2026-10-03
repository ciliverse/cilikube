package auth

import (
	"encoding/base32"
	"testing"
	"time"
)

func TestTOTPRoundTrip(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatal(err)
	}
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(1_700_000_000, 0)
	code := totpAt(key, now.Unix()/30)
	if !ValidateTOTP(secret, code, now) {
		t.Fatalf("code %s rejected", code)
	}
}
