package ai

import "strings"

// isZh reports whether the client UI language prefers Chinese answers.
func isZh(lang string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(lang)), "zh")
}

// normalizeLang returns "zh" or "en".
func normalizeLang(lang string) string {
	if isZh(lang) {
		return "zh"
	}
	return "en"
}
