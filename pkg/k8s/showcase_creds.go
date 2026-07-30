package k8s

// Public showcase login credentials — safe to display only when IsShowcase().
// These are NOT production secrets: the process talks to an in-memory fake cluster only.
const (
	ShowcaseAdminUser = "admin"
	ShowcaseAdminPass = "CiliKubeDemoAdmin2026!"
	ShowcaseGuestUser = "guest"
	ShowcaseGuestPass = "CiliKubeGuest2026!"
)

// ShowcaseAccount is returned by the public showcase info endpoint.
type ShowcaseAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Note     string `json:"note"`
}

// ShowcasePublicInfo is intentionally public in exhibit mode.
type ShowcasePublicInfo struct {
	Showcase bool              `json:"showcase"`
	Message  string            `json:"message"`
	Cluster  string            `json:"cluster"`
	Accounts []ShowcaseAccount `json:"accounts"`
}

// PublicShowcaseInfo returns demo credentials only when showcase mode is on.
// When showcase is off, Showcase=false and Accounts is empty (fail-closed).
func PublicShowcaseInfo() ShowcasePublicInfo {
	if !IsShowcase() {
		return ShowcasePublicInfo{Showcase: false}
	}
	return ShowcasePublicInfo{
		Showcase: true,
		Message:  "Multi-cluster showcase fleet: demo · prod-east · staging-lab (simulated)",
		Cluster:  ShowcaseClusterName,
		Accounts: []ShowcaseAccount{
			{
				Username: ShowcaseAdminUser,
				Password: ShowcaseAdminPass,
				Role:     "admin",
			},
			{
				Username: ShowcaseGuestUser,
				Password: ShowcaseGuestPass,
				Role:     "viewer",
			},
		},
	}
}
