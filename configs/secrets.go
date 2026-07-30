package configs

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Known-insecure defaults that must never ship in networked release mode.
var knownWeakJWTSecrets = map[string]struct{}{
	"": {},
	"cilikube-secret-key-change-in-production": {},
	"change-me-jwt-secret-in-production":       {},
	"CHANGE_ME_JWT_SECRET_KEY":                {},
	"your_jwt_secret_key_change_in_production": {},
	"change-me-jwt-secret":                     {},
}

var knownWeakEncryptionKeys = map[string]struct{}{
	"": {},
	"change-me-in-production-32chars!!": {},
	"CHANGE_ME_32_BYTE_ENCRYPTION_KEY!!": {},
}

// applySecretEnvOverrides prefers process env over file values (K8s/Docker).
func applySecretEnvOverrides() {
	if GlobalConfig == nil {
		return
	}
	if v := strings.TrimSpace(os.Getenv("JWT_SECRET")); v != "" {
		GlobalConfig.JWT.SecretKey = v
	}
	if v := strings.TrimSpace(os.Getenv("CILIKUBE_JWT_SECRET")); v != "" {
		GlobalConfig.JWT.SecretKey = v
	}
	if v := strings.TrimSpace(os.Getenv("ENCRYPTION_KEY")); v != "" {
		GlobalConfig.Server.EncryptionKey = v
	}
	if v := strings.TrimSpace(os.Getenv("CILIKUBE_ENCRYPTION_KEY")); v != "" {
		GlobalConfig.Server.EncryptionKey = v
	}
}

func looksLikePlaceholderSecret(v string) bool {
	u := strings.ToUpper(strings.TrimSpace(v))
	return strings.HasPrefix(u, "CHANGE_ME") ||
		strings.HasPrefix(u, "REPLACE_") ||
		strings.HasPrefix(u, "YOUR_") ||
		strings.Contains(u, "CHANGE-IN-PRODUCTION") ||
		strings.Contains(u, "CHANGE_IN_PRODUCTION")
}

func isWeakJWTSecret(secret string) bool {
	s := strings.TrimSpace(secret)
	if _, ok := knownWeakJWTSecrets[s]; ok {
		return true
	}
	return len(s) < 16 || looksLikePlaceholderSecret(s)
}

func isWeakEncryptionKey(key string) bool {
	k := strings.TrimSpace(key)
	if _, ok := knownWeakEncryptionKeys[k]; ok {
		return true
	}
	return len(k) != 32 || looksLikePlaceholderSecret(k)
}

// ValidateSecrets refuses known/empty secrets in release mode for networked servers.
// Desktop sidecars generate random secrets; debug mode only warns.
func ValidateSecrets() error {
	if GlobalConfig == nil {
		return fmt.Errorf("configuration not loaded")
	}

	applySecretEnvOverrides()

	jwtWeak := isWeakJWTSecret(GlobalConfig.JWT.SecretKey)
	encWeak := isWeakEncryptionKey(GlobalConfig.Server.EncryptionKey)

	if IsDesktop() {
		if jwtWeak || encWeak {
			return fmt.Errorf("desktop config has weak JWT or encryption key; delete userData config.yaml to regenerate")
		}
		return nil
	}

	mode := strings.ToLower(strings.TrimSpace(GlobalConfig.Server.Mode))
	if mode != "release" {
		if jwtWeak {
			log.Printf("warning: JWT secret_key is empty or a known default — set a strong secret before production")
		}
		if encWeak {
			log.Printf("warning: server.encryptionKey must be exactly 32 non-default characters for production")
		}
		return nil
	}

	if jwtWeak {
		return fmt.Errorf("refusing to start in release mode: JWT secret_key is empty, too short, or a known default — set JWT_SECRET / jwt.secret_key")
	}
	if encWeak {
		return fmt.Errorf("refusing to start in release mode: server.encryptionKey must be exactly 32 characters and not a known default — set ENCRYPTION_KEY / server.encryptionKey")
	}
	return nil
}
