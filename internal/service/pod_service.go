package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type PodService struct {
	client kubernetes.Interface
}

func NewPodService(client kubernetes.Interface) *PodService {
	return &PodService{client: client}
}

// ListNamespaces 列出所有命名空间
func (s *PodService) ListNamespaces() ([]string, error) {
	namespaceList, err := s.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	namespaces := make([]string, 0, len(namespaceList.Items))
	for _, ns := range namespaceList.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}

// 获取单个Pod
func (s *PodService) Get(namespace, name string) (*corev1.Pod, error) {
	return s.client.CoreV1().Pods(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Pod
func (s *PodService) Create(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {

	if pod.Namespace != "" && pod.Namespace != namespace {
		return nil, NewValidationError("pod namespace conflicts with path parameter")
	}

	return s.client.CoreV1().Pods(namespace).Create(
		context.TODO(),
		pod,
		metav1.CreateOptions{},
	)
}

// 更新Pod（包含冲突检测）
func (s *PodService) Update(namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	return s.client.CoreV1().Pods(namespace).Update(
		context.TODO(),
		pod,
		metav1.UpdateOptions{},
	)
}

// 删除Pod
func (s *PodService) Delete(namespace, name string) error {
	return s.client.CoreV1().Pods(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// 列表查询（支持分页和标签过滤）
func (s *PodService) List(namespace, selector string, limit int64) (*corev1.PodList, error) {
	return s.client.CoreV1().Pods(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch机制实现
func (s *PodService) Watch(namespace, selector string) (watch.Interface, error) {
	return s.client.CoreV1().Pods(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800),
		},
	)
}

// 辅助函数
func int64ptr(i int64) *int64 { return &i }

// 自定义错误类型
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(msg string) error {
	return &ValidationError{Message: msg}
}
