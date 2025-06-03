// package service

// import (
// 	"context"

// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/watch"
// 	"k8s.io/client-go/kubernetes"
// )

// type PVCService struct {
// 	client kubernetes.Interface
// }

// func NewPVCService(client kubernetes.Interface) *PVCService {
// 	return &PVCService{client: client}
// }

// // 获取单个PVC
// func (s *PVCService) Get(namespace, name string) (*corev1.PersistentVolumeClaim, error) {
// 	return s.client.CoreV1().PersistentVolumeClaims(namespace).Get(
// 		context.TODO(),
// 		name,
// 		metav1.GetOptions{},
// 	)
// }

// // 创建PVC
// func (s *PVCService) Create(namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {

// 	if pvc.Namespace != "" && pvc.Namespace != namespace {
// 		return nil, NewValidationError("pvc namespace conflicts with path parameter")
// 	}

// 	return s.client.CoreV1().PersistentVolumeClaims(namespace).Create(
// 		context.TODO(),
// 		pvc,
// 		metav1.CreateOptions{},
// 	)
// }

// // 更新PVC
// func (s *PVCService) Update(namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
// 	return s.client.CoreV1().PersistentVolumeClaims(namespace).Update(
// 		context.TODO(),
// 		pvc,
// 		metav1.UpdateOptions{},
// 	)
// }

// // 删除PVC
// func (s *PVCService) Delete(namespace, name string) error {
// 	return s.client.CoreV1().PersistentVolumeClaims(namespace).Delete(
// 		context.TODO(),
// 		name,
// 		metav1.DeleteOptions{},
// 	)
// }

// // 列表查询（支持分页和标签过滤）
// func (s *PVCService) List(namespace, selector string, limit int64) (*corev1.PersistentVolumeClaimList, error) {
// 	return s.client.CoreV1().PersistentVolumeClaims(namespace).List(
// 		context.TODO(),
// 		metav1.ListOptions{
// 			LabelSelector: selector,
// 			Limit:         limit,
// 		},
// 	)
// }

// // Watch机制实现
//
//	func (s *PVCService) Watch(namespace, selector string) (watch.Interface, error) {
//		return s.client.CoreV1().PersistentVolumeClaims(namespace).Watch(
//			context.TODO(),
//			metav1.ListOptions{
//				LabelSelector:  selector,
//				Watch:          true,
//				TimeoutSeconds: int64ptr(1800),
//			},
//		)
//	}
package service

import (
	"context"
	"k8s.io/client-go/kubernetes"

	corev1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PVCService 结构体不再持有 client 字段
type PVCService struct {
	// 不需要 client kubernetes.Interface 字段了
}

// NewPVCService 构造函数不再接收 kubernetes.Interface 参数
func NewPVCService() *PVCService {
	return &PVCService{}
}

// Get retrieves a single PersistentVolumeClaim by namespace and name.
func (s *PVCService) Get(clientSet kubernetes.Interface, namespace, name string) (*corev1.PersistentVolumeClaim, error) {
	return clientSet.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// List retrieves a list of PersistentVolumeClaims in a specific namespace.
// Supports label selector filtering and limit.
func (s *PVCService) List(clientSet kubernetes.Interface, namespace, labelSelector string, limit int64) (*corev1.PersistentVolumeClaimList, error) {
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}
	if limit > 0 {
		listOptions.Limit = limit
	}

	return clientSet.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), listOptions)
}

// Create creates a new PersistentVolumeClaim.
func (s *PVCService) Create(clientSet kubernetes.Interface, namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	// Validate namespace consistency
	if pvc.Namespace != "" && pvc.Namespace != namespace {
		return nil, NewValidationError("PVC namespace conflicts with path parameter")
	}
	if pvc.Namespace == "" {
		pvc.Namespace = namespace // Ensure namespace is set
	}
	if pvc.Name == "" {
		return nil, NewValidationError("PVC name cannot be empty")
	}
	// Add more validation for spec if needed (e.g., required fields)

	return clientSet.CoreV1().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
}

// Update updates an existing PersistentVolumeClaim.
// Important Note: Most fields in pvc.Spec (storageClassName, accessModes, resources.requests.storage, volumeMode)
// are **immutable** after creation. Updates usually only apply to metadata (labels, annotations).
// The Kubernetes API server will reject spec changes for immutable fields.
// This service function allows the update call, but relies on the API server for enforcement.
// Consider adding validation here to prevent attempts to change immutable fields if desired.
func (s *PVCService) Update(clientSet kubernetes.Interface, namespace string, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	if pvc.Namespace != "" && pvc.Namespace != namespace {
		return nil, NewValidationError("PVC namespace conflicts with path parameter")
	}
	if pvc.Namespace == "" {
		pvc.Namespace = namespace
	}
	if pvc.Name == "" {
		return nil, NewValidationError("PVC name cannot be empty for update")
	}

	return clientSet.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), pvc, metav1.UpdateOptions{})
}

// Delete deletes a PersistentVolumeClaim by namespace and name.
func (s *PVCService) Delete(clientSet kubernetes.Interface, namespace, name string) error {
	return clientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// --- Error Handling (reuse or define locally) ---
// type ValidationError struct { Message string }
// func (e *ValidationError) Error() string { return e.Message }
// func NewValidationError(msg string) error { return &ValidationError{Message: msg} }
