package auth

import (
	"testing"
)

func TestMaskPrivacyIP_EnvOnly(t *testing.T) {
	t.Setenv("CILIKUBE_PRIVACY_IPS", "1.2.3.4")
	masked := MaskPrivacyIP("1.2.3.4")
	if masked == "1.2.3.4" {
		t.Fatal("env privacy IP should be masked")
	}
	if masked != DecoyIP("1.2.3.4") {
		t.Fatalf("mask should match stable decoy: got %s want %s", masked, DecoyIP("1.2.3.4"))
	}
	if MaskPrivacyIP("8.8.8.8") != "8.8.8.8" {
		t.Fatal("unlisted IP must stay real")
	}
}

func TestMaskPrivacyIP_NoEnvPassthrough(t *testing.T) {
	t.Setenv("CILIKUBE_PRIVACY_IPS", "")
	sample := "203.0.113.99"
	if MaskPrivacyIP(sample) != sample {
		t.Fatal("without CILIKUBE_PRIVACY_IPS, IPs must not be masked")
	}
}
