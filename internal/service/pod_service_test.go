// internal/service/pod_service_test.go
package service_test

import (
	"context"
	"testing"

	"github.com/ciliverse/cilikube/internal/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCreatePod(t *testing.T) {
	// 创建fake客户端
	client := fake.NewSimpleClientset()
	svc := service.NewPodService(client)

	// 测试用例
	t.Run("successful creation", func(t *testing.T) {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
			},
		}

		_, err := svc.Create("default", pod)
		if err != nil {
			t.Fatalf("创建失败: %v", err)
		}

		// 验证Pod是否存在
		_, err = client.CoreV1().Pods("default").Get(context.TODO(), "test-pod", metav1.GetOptions{})
		if err != nil {
			t.Fatalf("找不到已创建的Pod: %v", err)
		}
	})
}
