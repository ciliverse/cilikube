package service

import (
	"context"
	"io"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Import net/url - Not directly used here, but might be needed elsewhere or was from previous iteration
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Used for Options
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme" // Required for Exec parameter encoding
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/yaml" // Preferred YAML library for K8s types
)

// PodService 结构体不再持有 client 和 config 字段
type PodService struct {
	// 不需要 client kubernetes.Interface 字段了
	// 不需要 config *rest.Config 字段了
}

// NewPodService 构造函数不再接收 kubernetes.Interface 和 *rest.Config 参数
func NewPodService() *PodService {
	return &PodService{}
}

// ListNamespaces 列出所有命名空间
func (s *PodService) ListNamespaces(clientSet kubernetes.Interface) ([]string, error) {
	namespaceList, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	namespaces := make([]string, 0, len(namespaceList.Items))
	for _, ns := range namespaceList.Items {
		namespaces = append(namespaces, ns.Name)
	}
	return namespaces, nil
}

// Get 获取单个Pod
func (s *PodService) Get(clientSet kubernetes.Interface, namespace, name string) (*corev1.Pod, error) {
	return clientSet.CoreV1().Pods(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// --- 添加缺失的 Create 方法 ---
// Create 创建 Pod（接收 Pod 对象, 用于 JSON 路径）
func (s *PodService) Create(clientSet kubernetes.Interface, namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	// 校验 Namespace (如果 Pod 对象中提供了 namespace)
	if pod.Namespace != "" && pod.Namespace != namespace {
		// 如果 Pod 对象中的 namespace 与 URL 路径参数不匹配，返回错误
		return nil, NewValidationError("pod namespace in body ('" + pod.Namespace + "') conflicts with path parameter ('" + namespace + "')")
	}
	// 确保最终创建时使用 URL 路径中的 namespace (或 YAML 中解析出的，如果调用 CreateFromYAML)
	pod.Namespace = namespace // Overwrite or set namespace from path parameter

	// 调用 Kubernetes API 创建 Pod
	return clientSet.CoreV1().Pods(namespace).Create(
		context.TODO(),
		pod, // 传递构造好的 Pod 对象
		metav1.CreateOptions{},
	)
}

// CreateFromYAML 创建 Pod (从 YAML)
func (s *PodService) CreateFromYAML(clientSet kubernetes.Interface, namespace string, yamlContent []byte) (*corev1.Pod, error) {
	var pod corev1.Pod
	err := yaml.Unmarshal(yamlContent, &pod)
	if err != nil {
		return nil, NewValidationError("无效的 YAML 格式: " + err.Error())
	}

	// 校验 Namespace
	if pod.Namespace != "" && pod.Namespace != namespace {
		return nil, NewValidationError("YAML 中的 namespace ('" + pod.Namespace + "') 与请求路径中的 namespace ('" + namespace + "') 不匹配")
	}
	// 如果 YAML 中未指定 namespace，则使用路径参数中的 namespace
	if pod.Namespace == "" {
		pod.Namespace = namespace
	}

	// 可选：校验 Kind 和 APIVersion
	// ...

	// 调用 Kubernetes API 创建 Pod (注意：这里仍然调用 K8s Client 的 Create)
	// 理论上也可以直接调用上面我们添加的 s.Create 方法，但直接调用 client 也一样
	return clientSet.CoreV1().Pods(pod.Namespace).Create(
		context.TODO(),
		&pod,
		metav1.CreateOptions{},
	)
}

// --- 添加缺失的 Update 方法 ---
// Update 更新 Pod（接收 Pod 对象, 用于 JSON 路径）
func (s *PodService) Update(clientSet kubernetes.Interface, namespace string, pod *corev1.Pod) (*corev1.Pod, error) {
	// 确保要更新的 Pod 对象的 Namespace 与 URL 路径参数一致
	if pod.Namespace != namespace {
		return nil, NewValidationError("pod namespace in body ('" + pod.Namespace + "') conflicts with path parameter ('" + namespace + "') during update")
	}

	// 调用 Kubernetes API 更新 Pod
	return clientSet.CoreV1().Pods(namespace).Update(
		context.TODO(),
		pod, // 传递构造好的、待更新的 Pod 对象
		metav1.UpdateOptions{},
	)
}

// UpdateFromYAML 更新Pod (从 YAML)
func (s *PodService) UpdateFromYAML(clientSet kubernetes.Interface, namespace, name string, yamlContent []byte) (*corev1.Pod, error) {
	var updatedPod corev1.Pod
	err := yaml.Unmarshal(yamlContent, &updatedPod)
	if err != nil {
		return nil, NewValidationError("无效的 YAML 格式: " + err.Error())
	}

	// --- 严格校验 ---
	// 1. 名称必须匹配
	if updatedPod.Name != name {
		return nil, NewValidationError("YAML 中的 metadata.name ('" + updatedPod.Name + "') 与请求路径中的 name ('" + name + "') 不匹配")
	}
	// 2. 命名空间必须匹配（或在 YAML 中为空，此时使用路径参数）
	if updatedPod.Namespace == "" {
		updatedPod.Namespace = namespace // Set namespace from path if missing in YAML
	} else if updatedPod.Namespace != namespace {
		return nil, NewValidationError("YAML 中的 metadata.namespace ('" + updatedPod.Namespace + "') 与请求路径中的 namespace ('" + namespace + "') 不匹配")
	}
	// 3. Kind 和 APIVersion (可选但推荐)
	// ...

	// 调用 Kubernetes API 更新 Pod (注意：这里仍然调用 K8s Client 的 Update)
	// 也可以调用上面我们添加的 s.Update 方法
	return clientSet.CoreV1().Pods(namespace).Update(
		context.TODO(),
		&updatedPod, // 传递反序列化后的 Pod 对象
		metav1.UpdateOptions{},
	)
}

// Delete 删除Pod
func (s *PodService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.CoreV1().Pods(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// List 列表查询（支持分页和标签过滤）
func (s *PodService) List(clientSet kubernetes.Interface, namespace, selector string, limit int64) (*corev1.PodList, error) {
	return clientSet.CoreV1().Pods(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
			Limit:         limit,
		},
	)
}

// Watch 机制实现
func (s *PodService) Watch(clientSet kubernetes.Interface, namespace, selector string) (watch.Interface, error) {
	return clientSet.CoreV1().Pods(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800), // 30 minutes
		},
	)
}

// GetPodLogs 获取 Pod 日志流
func (s *PodService) GetPodLogs(clientSet kubernetes.Interface, namespace, podName string, opts *corev1.PodLogOptions) (io.ReadCloser, error) {
	req := clientSet.CoreV1().Pods(namespace).GetLogs(podName, opts)
	return req.Stream(context.TODO())
}

// GetPodYAML 获取 Pod 的 YAML 定义
func (s *PodService) GetPodYAML(clientSet kubernetes.Interface, namespace, name string) ([]byte, error) {
	pod, err := s.Get(clientSet, namespace, name)
	if err != nil {
		return nil, err
	}

	// 清理 Kubernetes 添加的内部字段，使输出更干净（可选）
	pod.ObjectMeta.ManagedFields = nil
	pod.ObjectMeta.ResourceVersion = ""
	pod.ObjectMeta.UID = ""
	pod.ObjectMeta.SelfLink = ""
	pod.ObjectMeta.Generation = 0
	// pod.Status = corev1.PodStatus{} // Optionally clear status for pure definition

	yamlBytes, err := yaml.Marshal(pod)
	if err != nil {
		return nil, err
	}
	return yamlBytes, nil
}

// ExecOptions 定义 Exec 所需的选项
type ExecOptions struct {
	Namespace     string
	PodName       string
	ContainerName string
	Command       []string
	Stdin         io.Reader
	Stdout        io.Writer
	Stderr        io.Writer
	Tty           bool
}

// ExecIntoPod 在 Pod 容器内执行命令
func (s *PodService) ExecIntoPod(clientSet kubernetes.Interface, config *rest.Config, ctx context.Context, opts ExecOptions) error {
	req := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(opts.PodName).
		Namespace(opts.Namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Container: opts.ContainerName,
		Command:   opts.Command,
		Stdin:     opts.Stdin != nil,
		Stdout:    opts.Stdout != nil,
		Stderr:    opts.Stderr != nil,
		TTY:       opts.Tty,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  opts.Stdin,
		Stdout: opts.Stdout,
		Stderr: opts.Stderr,
		Tty:    opts.Tty,
		// TerminalSizeQueue: // TODO: Handle terminal resize notifications if TTY
	})
	return err // Return the error from StreamWithContext
}

// --- Helper Functions ---
func int64ptr(i int64) *int64 { return &i }

// 自定义错误类型
type ValidationError struct{ Message string }

func (e *ValidationError) Error() string  { return e.Message }
func NewValidationError(msg string) error { return &ValidationError{Message: msg} }
