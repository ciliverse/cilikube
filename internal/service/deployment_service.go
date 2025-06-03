package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// DeploymentService 结构体不再持有 client 字段
type DeploymentService struct {
	// 不需要 client kubernetes.Interface 字段了
}

// NewDeploymentService 构造函数不再接收 kubernetes.Interface 参数
func NewDeploymentService() *DeploymentService {
	return &DeploymentService{}
}

// 获取单个Deployment
func (s *DeploymentService) Get(clientSet kubernetes.Interface, namespace, name string) (*appsv1.Deployment, error) {
	return clientSet.AppsV1().Deployments(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
}

// 创建Deployment
func (s *DeploymentService) Create(clientSet kubernetes.Interface, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {

	if deployment.Namespace != "" && deployment.Namespace != namespace {
		return nil, NewValidationError("deployment namespace conflicts with path parameter")
	}

	return clientSet.AppsV1().Deployments(namespace).Create(
		context.TODO(),
		deployment,
		metav1.CreateOptions{},
	)
}

// 更新Deployment（包含冲突检测）
func (s *DeploymentService) Update(clientSet kubernetes.Interface, namespace, name string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	// --- 严格校验 ---
	// 1. 名称必须匹配
	if deployment.Name != name {
		return nil, NewValidationError("YAML / JSON 中的 metadata.name ('" + deployment.Name + "') 与请求路径中的 name ('" + name + "') 不匹配")
	}
	// 2. 命名空间必须匹配（或在 YAML / JSON 中为空，此时使用路径参数）
	if deployment.Namespace == "" {
		deployment.Namespace = namespace // Set namespace from path if missing in YAML / JSON
	} else if deployment.Namespace != namespace {
		return nil, NewValidationError("YAML / JSON 中的 metadata.namespace ('" + deployment.Namespace + "') 与请求路径中的 namespace ('" + namespace + "') 不匹配")
	}
	// 3. Kind 和 APIVersion (可选但推荐)
	// ...

	return clientSet.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// 删除Deployment
func (s *DeploymentService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.AppsV1().Deployments(namespace).Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
}

// ListDeployments 列出所有Deployment
func (s *DeploymentService) List(clientSet kubernetes.Interface, namespace string) (*appsv1.DeploymentList, error) {
	return clientSet.AppsV1().Deployments(namespace).List(
		context.TODO(),
		metav1.ListOptions{},
	)
}

// ListDeploymentsByLabels 根据标签过滤列出Deployment
func (s *DeploymentService) ListByLabels(clientSet kubernetes.Interface, namespace, selector string) (*appsv1.DeploymentList, error) {
	return clientSet.AppsV1().Deployments(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: selector,
		},
	)
}

// WatchDeployments 实现Watch机制
func (s *DeploymentService) Watch(clientSet kubernetes.Interface, namespace, selector string) (watch.Interface, error) {
	return clientSet.AppsV1().Deployments(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector:  selector,
			Watch:          true,
			TimeoutSeconds: int64ptr(1800), // 30 minutes
		},
	)
}

// PatchDeployment 实现Patch机制
// func (s *DeploymentService) Patch(clientSet kubernetes.Interface, namespace, name string, patchData []byte) (*appsv1.Deployment, error) {
// 	return clientSet.AppsV1().Deployments(namespace).Patch(
// 		context.TODO(),
// 		name,
// 		types.StrategicMergePatchType,
// 		patchData,
// 		metav1.PatchOptions{},
// 	)
// }

// ReplaceDeployment 实现Replace机制
func (s *DeploymentService) Replace(clientSet kubernetes.Interface, namespace, name string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return clientSet.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// RollbackDeployment 实现Deployment回滚

// ScaleDeployment 实现Deployment扩缩容
func (s *DeploymentService) Scale(clientSet kubernetes.Interface, namespace, name string, replicas int32) (*appsv1.Deployment, error) {
	deployment, err := s.Get(clientSet, namespace, name)
	if err != nil {
		return nil, err
	}

	deployment.Spec.Replicas = &replicas
	return clientSet.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// PauseDeployment 实现Deployment暂停
func (s *DeploymentService) Pause(clientSet kubernetes.Interface, namespace, name string) (*appsv1.Deployment, error) {
	deployment, err := s.Get(clientSet, namespace, name)
	if err != nil {
		return nil, err
	}

	deployment.Spec.Paused = true
	return clientSet.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// ResumeDeployment 实现Deployment恢复
func (s *DeploymentService) Resume(clientSet kubernetes.Interface, namespace, name string) (*appsv1.Deployment, error) {
	deployment, err := s.Get(clientSet, namespace, name)
	if err != nil {
		return nil, err
	}

	deployment.Spec.Paused = false
	return clientSet.AppsV1().Deployments(namespace).Update(
		context.TODO(),
		deployment,
		metav1.UpdateOptions{},
	)
}

// PodList 实现获取Deployment关联的Pod列表查询（支持分页和标签过滤）
func (s *DeploymentService) PodList(clientSet kubernetes.Interface, namespace, deploymentName string, limit int64) (*corev1.PodList, error) {
	// 1. 获取 Deployment
	deployment, err := s.Get(clientSet, namespace, deploymentName)
	if err != nil {
		return nil, err
	}

	// 2. 获取 ReplicaSet（Deployment 控制器会创建 ReplicaSet）
	rsList, err := clientSet.AppsV1().ReplicaSets(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels).String(),
		},
	)
	if err != nil {
		return nil, err
	}

	// 3. 初始化一个空的 Selector
	allSelectors := labels.NewSelector()

	// 4. 过滤出 replicas > 0 的 ReplicaSet
	for _, rs := range rsList.Items {
		if rs.Spec.Replicas != nil && *rs.Spec.Replicas > 0 {
			// 5. 将 ReplicaSet 的 LabelSelector 转换为 Selector
			selector, err := metav1.LabelSelectorAsSelector(rs.Spec.Selector)
			if err != nil {
				return nil, err
			}

			// 6. 将 Selector 转换为 Requirements
			requirements, selectable := selector.Requirements()
			if !selectable {
				continue
			}

			// 7. 将 Requirements 添加到 allSelectors 中
			allSelectors = allSelectors.Add(requirements...)
		}
	}

	// 8. 如果没有活跃的 ReplicaSet，返回空列表
	if allSelectors.String() == "" {
		return &corev1.PodList{}, nil
	}

	// 9. 查询 Pod 列表
	return clientSet.CoreV1().Pods(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: allSelectors.String(),
			Limit:         limit,
		},
	)
}
