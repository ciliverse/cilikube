package auth

import (
	"crypto/sha256"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Privacy IPs are never hardcoded — set CILIKUBE_PRIVACY_IPS (comma-separated)
// on the public showcase host only. Example:
//   Environment=CILIKUBE_PRIVACY_IPS=x.x.x.x
// Local/dev without the env var records real client IPs unchanged.

// Plausible public decoys (not the operator's network). Stable pick per real IP.
var privacyDecoyPool = []string{
	"47.98.211.64",
	"120.55.168.22",
	"39.105.120.88",
	"101.37.25.140",
	"112.124.58.201",
	"121.41.66.93",
}

func privacyIPSet() map[string]struct{} {
	set := make(map[string]struct{})
	for _, part := range strings.Split(os.Getenv("CILIKUBE_PRIVACY_IPS"), ",") {
		ip := strings.TrimSpace(part)
		if ip != "" {
			set[ip] = struct{}{}
		}
	}
	return set
}

// IsPrivacyIP reports whether an IP should be redacted from audit/session logs.
func IsPrivacyIP(ip string) bool {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return false
	}
	_, ok := privacyIPSet()[ip]
	return ok
}

// DecoyIP returns a stable fake address for a privacy IP (same real IP → same decoy).
func DecoyIP(real string) string {
	sum := sha256.Sum256([]byte("cilikube-privacy-v1|" + strings.TrimSpace(real)))
	return privacyDecoyPool[int(sum[0])%len(privacyDecoyPool)]
}

// MaskPrivacyIP replaces configured privacy IPs with a decoy; all other IPs pass through.
func MaskPrivacyIP(ip string) string {
	ip = strings.TrimSpace(ip)
	if !IsPrivacyIP(ip) {
		return ip
	}
	return DecoyIP(ip)
}

// AuditClientIP is ClientIP() with operator privacy masking applied (env-driven).
func AuditClientIP(c *gin.Context) string {
	if c == nil {
		return ""
	}
	return MaskPrivacyIP(c.ClientIP())
}

// ScrubPrivacyIPs rewrites historical audit/session rows that still contain privacy IPs.
func ScrubPrivacyIPs(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	for real := range privacyIPSet() {
		decoy := DecoyIP(real)
		if err := db.Exec(
			`UPDATE audit_logs SET ip_address = ?, details = REPLACE(COALESCE(details, ''), ?, ?) WHERE ip_address = ? OR details LIKE ?`,
			decoy, real, decoy, real, "%"+real+"%",
		).Error; err != nil {
			return err
		}
		if err := db.Exec(
			`UPDATE user_sessions SET ip_address = ? WHERE ip_address = ?`,
			decoy, real,
		).Error; err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "no such table") {
				return err
			}
		}
		if err := db.Exec(
			`UPDATE login_attempts SET ip_address = ? WHERE ip_address = ?`,
			decoy, real,
		).Error; err != nil {
			if !strings.Contains(strings.ToLower(err.Error()), "no such table") {
				return err
			}
		}
	}
	return nil
}
