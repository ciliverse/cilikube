package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/pkg/k8s"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// CRDService defines the interface for CRD operations
type CRDService interface {
	// CRD management
	ListCRDs(client *k8s.Client) (*models.CRDListResponse, error)
	GetCRD(client *k8s.Client, name string) (*models.CRDDetailResponse, error)

	// Custom resource management
	ListCustomResources(client *k8s.Client, group, version, plural, namespace string, limit int64, continueToken string) (*models.CustomResourceListResponse, error)
	GetCustomResource(client *k8s.Client, group, version, plural, namespace, name string) (*models.CustomResourceItem, error)
	CreateCustomResource(client *k8s.Client, group, version, plural, namespace string, resource *models.CustomResourceRequest) (*models.CustomResourceItem, error)
	UpdateCustomResource(client *k8s.Client, group, version, plural, namespace, name string, resource *models.CustomResourceRequest) (*models.CustomResourceItem, error)
	DeleteCustomResource(client *k8s.Client, group, version, plural, namespace, name string) error
}

type crdService struct{}

// NewCRDService creates a new CRD service instance
func NewCRDService() CRDService {
	return &crdService{}
}

// ListCRDs retrieves the list of CRDs
func (s *crdService) ListCRDs(client *k8s.Client) (*models.CRDListResponse, error) {
	apiExtClient, err := apiextensionsclientset.NewForConfig(client.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create apiextensions client: %w", err)
	}

	crdList, err := apiExtClient.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list CRDs: %w", err)
	}

	items := make([]models.CRDItem, 0, len(crdList.Items))
	for _, crd := range crdList.Items {
		item := s.convertCRDToItem(&crd)
		items = append(items, item)
	}

	// Sort by name
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return &models.CRDListResponse{
		Items: items,
		Total: len(items),
	}, nil
}

// GetCRD retrieves CRD details
func (s *crdService) GetCRD(client *k8s.Client, name string) (*models.CRDDetailResponse, error) {
	apiExtClient, err := apiextensionsclientset.NewForConfig(client.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create apiextensions client: %w", err)
	}

	crd, err := apiExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get CRD: %w", err)
	}

	detail := &models.CRDDetailResponse{
		CRDItem: s.convertCRDToItem(crd),
	}

	// Add version information
	for _, version := range crd.Spec.Versions {
		detail.Versions = append(detail.Versions, models.CRDVersion{
			Name:    version.Name,
			Served:  version.Served,
			Storage: version.Storage,
		})
	}

	// Add status conditions
	for _, condition := range crd.Status.Conditions {
		detail.Conditions = append(detail.Conditions, models.CRDCondition{
			Type:               string(condition.Type),
			Status:             string(condition.Status),
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		})
	}

	return detail, nil
}

// ListCustomResources retrieves the list of custom resources
func (s *crdService) ListCustomResources(client *k8s.Client, group, version, plural, namespace string, limit int64, continueToken string) (*models.CustomResourceListResponse, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}

	dynamicClient, err := dynamic.NewForConfig(client.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	listOptions := metav1.ListOptions{
		Limit:    limit,
		Continue: continueToken,
	}

	var list *unstructured.UnstructuredList
	if namespace != "" {
		list, err = dynamicClient.Resource(gvr).Namespace(namespace).List(context.TODO(), listOptions)
	} else {
		list, err = dynamicClient.Resource(gvr).List(context.TODO(), listOptions)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list custom resources: %w", err)
	}

	items := make([]models.CustomResourceItem, 0, len(list.Items))
	for _, item := range list.Items {
		crItem := s.convertUnstructuredToItem(&item)
		items = append(items, crItem)
	}

	return &models.CustomResourceListResponse{
		Items: items,
		Total: len(items),
	}, nil
}

// GetCustomResource retrieves custom resource details
func (s *crdService) GetCustomResource(client *k8s.Client, group, version, plural, namespace, name string) (*models.CustomResourceItem, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}

	dynamicClient, err := dynamic.NewForConfig(client.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	var obj *unstructured.Unstructured
	if namespace != "" {
		obj, err = dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	} else {
		obj, err = dynamicClient.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get custom resource: %w", err)
	}

	item := s.convertUnstructuredToItem(obj)
	return &item, nil
}

// CreateCustomResource creates a custom resource
func (s *crdService) CreateCustomResource(client *k8s.Client, group, version, plural, namespace string, resource *models.CustomResourceRequest) (*models.CustomResourceItem, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}

	dynamicClient, err := dynamic.NewForConfig(client.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Build unstructured object
	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": resource.APIVersion,
			"kind":       resource.Kind,
			"metadata":   resource.Metadata,
		},
	}

	if resource.Spec != nil {
		obj.Object["spec"] = resource.Spec
	}

	var created *unstructured.Unstructured
	if namespace != "" {
		created, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	} else {
		created, err = dynamicClient.Resource(gvr).Create(context.TODO(), obj, metav1.CreateOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create custom resource: %w", err)
	}

	item := s.convertUnstructuredToItem(created)
	return &item, nil
}

// UpdateCustomResource updates a custom resource
func (s *crdService) UpdateCustomResource(client *k8s.Client, group, version, plural, namespace, name string, resource *models.CustomResourceRequest) (*models.CustomResourceItem, error) {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}

	dynamicClient, err := dynamic.NewForConfig(client.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Get existing resource first
	var existing *unstructured.Unstructured
	if namespace != "" {
		existing, err = dynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	} else {
		existing, err = dynamicClient.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get existing custom resource: %w", err)
	}

	// Update fields
	if resource.Metadata != nil {
		existing.Object["metadata"] = resource.Metadata
	}
	if resource.Spec != nil {
		existing.Object["spec"] = resource.Spec
	}

	var updated *unstructured.Unstructured
	if namespace != "" {
		updated, err = dynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), existing, metav1.UpdateOptions{})
	} else {
		updated, err = dynamicClient.Resource(gvr).Update(context.TODO(), existing, metav1.UpdateOptions{})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update custom resource: %w", err)
	}

	item := s.convertUnstructuredToItem(updated)
	return &item, nil
}

// DeleteCustomResource deletes a custom resource
func (s *crdService) DeleteCustomResource(client *k8s.Client, group, version, plural, namespace, name string) error {
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: plural,
	}

	dynamicClient, err := dynamic.NewForConfig(client.Config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	if namespace != "" {
		err = dynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	} else {
		err = dynamicClient.Resource(gvr).Delete(context.TODO(), name, metav1.DeleteOptions{})
	}

	if err != nil {
		return fmt.Errorf("failed to delete custom resource: %w", err)
	}

	return nil
}

// convertCRDToItem converts CRD to CRDItem
func (s *crdService) convertCRDToItem(crd *apiextensionsv1.CustomResourceDefinition) models.CRDItem {
	var version string
	var categories, shortNames []string

	// Get storage version
	for _, v := range crd.Spec.Versions {
		if v.Storage {
			version = v.Name
			break
		}
	}
	if version == "" && len(crd.Spec.Versions) > 0 {
		version = crd.Spec.Versions[0].Name
	}

	// Get categories and shortNames
	if crd.Spec.Names.Categories != nil {
		categories = crd.Spec.Names.Categories
	}
	if crd.Spec.Names.ShortNames != nil {
		shortNames = crd.Spec.Names.ShortNames
	}

	return models.CRDItem{
		Name:        crd.Name,
		Group:       crd.Spec.Group,
		Version:     version,
		Kind:        crd.Spec.Names.Kind,
		Plural:      crd.Spec.Names.Plural,
		Singular:    crd.Spec.Names.Singular,
		Scope:       string(crd.Spec.Scope),
		Categories:  categories,
		ShortNames:  shortNames,
		CreatedAt:   crd.CreationTimestamp,
		Labels:      crd.Labels,
		Annotations: crd.Annotations,
	}
}

// convertUnstructuredToItem converts Unstructured to CustomResourceItem
func (s *crdService) convertUnstructuredToItem(obj *unstructured.Unstructured) models.CustomResourceItem {
	item := models.CustomResourceItem{
		Name:        obj.GetName(),
		Namespace:   obj.GetNamespace(),
		Kind:        obj.GetKind(),
		APIVersion:  obj.GetAPIVersion(),
		CreatedAt:   metav1.NewTime(obj.GetCreationTimestamp().Time),
		Labels:      obj.GetLabels(),
		Annotations: obj.GetAnnotations(),
	}

	// Get spec and status
	if spec, found, err := unstructured.NestedMap(obj.Object, "spec"); found && err == nil {
		item.Spec = spec
	}
	if status, found, err := unstructured.NestedMap(obj.Object, "status"); found && err == nil {
		item.Status = status
	}

	return item
}
