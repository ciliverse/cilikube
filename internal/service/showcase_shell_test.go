package service

import "testing"

func TestShowcaseKubectlReadOnly(t *testing.T) {
	out, done := ShowcaseKubectl("kubectl get nodes")
	if done {
		t.Fatal("get nodes should keep the session open")
	}
	if !stringsContains(out, "demo-worker-1") || stringsContains(out, "172.") {
		t.Fatalf("nodes table = %q", out)
	}
	out, done = ShowcaseKubectl("exit")
	if !done || !stringsContains(out, "closed") {
		t.Fatalf("exit = %q done=%v", out, done)
	}
}

func TestShowcaseKubectlBlocksHostCommands(t *testing.T) {
	for _, line := range []string{"kubectl get pods; id", "ls /", "$(id)", "kubectl get pods | sh"} {
		out, done := ShowcaseKubectl(line)
		if done {
			t.Fatalf("%q closed the session", line)
		}
		if !stringsContains(out, "not available") {
			t.Fatalf("%q -> %q", line, out)
		}
	}
}

func TestShowcaseNodeShellDoesNotReadHost(t *testing.T) {
	out, _ := ShowcaseNodeShell("demo-worker-2", "hostname")
	if out != "demo-worker-2\r\n" {
		t.Fatalf("hostname = %q", out)
	}
	out, _ = ShowcaseNodeShell("demo-worker-2", "cat /etc/shadow")
	if !stringsContains(out, "not available") {
		t.Fatalf("cat = %q", out)
	}
}

func stringsContains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && (s == sub || len(s) > 0 && contains(s, sub)))
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
