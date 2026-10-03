package service

import "testing"

func TestBuildNodeShellPod(t *testing.T) {
	pod, err := BuildNodeShellPod("worker-a", "")
	if err != nil {
		t.Fatal(err)
	}
	if pod.Spec.NodeName != "worker-a" {
		t.Fatalf("node %s", pod.Spec.NodeName)
	}
	if pod.Namespace != "kube-system" || !pod.Spec.HostPID {
		t.Fatalf("pod spec %+v", pod.Spec)
	}
	c := pod.Spec.Containers[0]
	if c.Image != "busybox:1.36" || c.Command[0] != "sleep" {
		t.Fatalf("container %+v", c)
	}
	if c.SecurityContext == nil || c.SecurityContext.Privileged == nil || !*c.SecurityContext.Privileged {
		t.Fatal("expected privileged")
	}
	if _, err := BuildNodeShellPod(" ", ""); err == nil {
		t.Fatal("expected empty node error")
	}
	if _, err := BuildNodeShellPod("n", "bad image"); err == nil {
		t.Fatal("expected bad image error")
	}
}
