package service

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/k8s"
)

// ClusterService provides business logic around cluster management.
type ClusterService struct {
	k8sManager     *k8s.ClusterManager
	kubeconfigPath string
}

// NewClusterService creates a new ClusterService instance.
func NewClusterService(k8sManager *k8s.ClusterManager, kubeconfigPath string) *ClusterService {
	return &ClusterService{
		k8sManager:     k8sManager,
		kubeconfigPath: kubeconfigPath,
	}
}

// ListClusters returns a list of summary information for all managed clusters.
func (s *ClusterService) ListClusters() []models.ClusterListResponse {
	// The information structure returned by k8sManager is already suitable for the list page, we just convert it
	managerInfo := s.k8sManager.ListClusterInfo()
	response := make([]models.ClusterListResponse, len(managerInfo))
	for i, info := range managerInfo {
		response[i] = models.ClusterListResponse{
			ID:          info.ID, // Ensure k8s.ClusterInfoResponse has ID field
			Name:        info.Name,
			Server:      info.Server,
			Version:     info.Version,
			Status:      info.Status,
			Source:      info.Source,
			Environment: info.Environment,
		}
	}
	return response
}

// GetClusterByID gets detailed information for a single cluster.
func (s *ClusterService) GetClusterByID(id string) (*models.ClusterResponse, error) {
	cluster, err := s.k8sManager.GetClusterDetailFromDB(id)
	if err != nil {
		// If not in database, it might be a file-type cluster, we assemble a simple version from cache
		if info, ok := s.k8sManager.GetStatusFromCache(id); ok {
			return &models.ClusterResponse{
				ID:          info.ID,
				Name:        info.Name,
				Version:     info.Version,
				Status:      info.Status,
				Environment: info.Environment,
				Source:      info.Source,
			}, nil
		}
		return nil, fmt.Errorf("cluster ID '%s' not found: %w", id, err)
	}

	return &models.ClusterResponse{
		ID:          cluster.ID,
		Name:        cluster.Name,
		Provider:    cluster.Provider,
		Description: cluster.Description,
		Environment: cluster.Environment,
		Region:      cluster.Region,
		Version:     cluster.Version,
		Status:      cluster.Status,
		Labels:      cluster.Labels,
		CreatedAt:   cluster.CreatedAt,
		UpdatedAt:   cluster.UpdatedAt,
	}, nil
}

// CreateCluster handles the logic for creating a new cluster.
func (s *ClusterService) CreateCluster(req models.CreateClusterRequest) error {
	if k8s.IsShowcase() {
		return fmt.Errorf("showcase mode: importing real clusters is disabled")
	}
	// 1. Validate kubeconfig
	config, err := s.validateKubeconfig(req.KubeconfigData)
	if err != nil {
		return fmt.Errorf("invalid kubeconfig: %w", err)
	}

	// 2. Test connection
	if err := s.testConnection(config); err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	// 3. Decode and create cluster
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(req.KubeconfigData)
	if err != nil {
		return fmt.Errorf("kubeconfig data is not valid Base64 encoding: %w", err)
	}
	cluster := &store.Cluster{
		Name:           req.Name,
		KubeconfigData: kubeconfigBytes,
		Provider:       req.Provider,
		Description:    req.Description,
		Environment:    req.Environment,
		Region:         req.Region,
	}
	return s.k8sManager.AddDBCluster(cluster)
}

// UpdateCluster updates cluster information.
func (s *ClusterService) UpdateCluster(id string, req models.UpdateClusterRequest) error {
	if k8s.IsShowcase() {
		return fmt.Errorf("showcase mode: cluster updates are disabled")
	}
	return s.k8sManager.UpdateDBCluster(id, req)
}

// DeleteClusterByID handles the logic for deleting a cluster.
func (s *ClusterService) DeleteClusterByID(id string) error {
	if k8s.IsShowcase() {
		return fmt.Errorf("showcase mode: cluster deletion is disabled")
	}
	return s.k8sManager.RemoveDBClusterByID(id)
}

// SetActiveCluster handles the logic for switching the active cluster.
func (s *ClusterService) SetActiveCluster(id string) error {
	return s.k8sManager.SetActiveClusterByID(id)
}

// GetActiveClusterID gets the current active cluster ID
func (s *ClusterService) GetActiveClusterID() string {
	return s.k8sManager.GetActiveClusterID()
}

// validateKubeconfig validates the kubeconfig data
func (s *ClusterService) validateKubeconfig(kubeconfigData string) (*rest.Config, error) {
	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(kubeconfigData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode kubeconfig: %w", err)
	}

	// Parse kubeconfig
	config, err := clientcmd.RESTConfigFromKubeConfig(decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %w", err)
	}

	return config, nil
}

// testConnection tests the connection to the Kubernetes cluster.
// Must preserve client certificate auth (k3s kubeconfigs use cert/key, not bearer tokens).
func (s *ClusterService) testConnection(config *rest.Config) error {
	testConfig := rest.CopyConfig(config)
	if testConfig.Timeout == 0 {
		testConfig.Timeout = 15 * time.Second
	}

	clientset, err := kubernetes.NewForConfig(testConfig)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	if _, err = clientset.Discovery().ServerVersion(); err != nil {
		return fmt.Errorf("failed to connect to cluster: %w", err)
	}

	return nil
}

func (s *ClusterService) resolveLocalKubeconfigPath() (string, error) {
	path := s.kubeconfigPath
	if path == "" || path == "default" {
		if env := os.Getenv("KUBECONFIG"); env != "" {
			path = env
		} else if home := homedir.HomeDir(); home != "" {
			path = filepath.Join(home, ".kube", "config")
		}
	}
	if path == "" {
		return "", fmt.Errorf("no local kubeconfig path configured")
	}
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("local kubeconfig not readable at %s: %w", path, err)
	}
	return path, nil
}

// ListLocalKubeContexts lists contexts from the API host kubeconfig.
func (s *ClusterService) ListLocalKubeContexts() (*models.LocalKubeContextsResponse, error) {
	if k8s.IsShowcase() {
		return nil, fmt.Errorf("showcase mode: local kubeconfig access is disabled")
	}
	path, err := s.resolveLocalKubeconfigPath()
	if err != nil {
		return nil, err
	}
	apiCfg, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	existingByName := map[string]string{}
	existingByServer := map[string]string{}
	for _, info := range s.k8sManager.ListClusterInfo() {
		existingByName[info.Name] = info.Name
		if info.Server != "" {
			existingByServer[info.Server] = info.Name
		}
	}

	out := &models.LocalKubeContextsResponse{
		Path:     path,
		Contexts: make([]models.LocalKubeContext, 0, len(apiCfg.Contexts)),
	}
	for name, ctx := range apiCfg.Contexts {
		if ctx == nil {
			continue
		}
		server := ""
		if cl := apiCfg.Clusters[ctx.Cluster]; cl != nil {
			server = cl.Server
		}
		item := models.LocalKubeContext{
			Name:    name,
			Cluster: ctx.Cluster,
			Server:  server,
			User:    ctx.AuthInfo,
		}
		if existingByName[name] != "" {
			item.ConflictName = true
			item.ExistingName = existingByName[name]
		} else if existingByServer[server] != "" {
			item.ConflictServer = true
			item.ExistingName = existingByServer[server]
		}
		out.Contexts = append(out.Contexts, item)
	}
	return out, nil
}

// ImportLocalClusters imports selected local contexts into the DB as clusters.
func (s *ClusterService) ImportLocalClusters(req models.ImportLocalClustersRequest) (*models.ImportLocalClustersResult, error) {
	if k8s.IsShowcase() {
		return nil, fmt.Errorf("showcase mode: importing local kubeconfig is disabled")
	}
	path, err := s.resolveLocalKubeconfigPath()
	if err != nil {
		return nil, err
	}
	apiCfg, err := clientcmd.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	skipConflicts := true
	if req.SkipConflicts != nil {
		skipConflicts = *req.SkipConflicts
	}

	result := &models.ImportLocalClustersResult{
		Imported: []string{},
		Skipped:  map[string]string{},
		Failed:   map[string]string{},
	}

	existingByName := map[string]bool{}
	existingByServer := map[string]string{}
	for _, info := range s.k8sManager.ListClusterInfo() {
		existingByName[info.Name] = true
		if info.Server != "" {
			existingByServer[info.Server] = info.Name
		}
	}

	for _, ctxName := range req.Contexts {
		ctx := apiCfg.Contexts[ctxName]
		if ctx == nil {
			result.Failed[ctxName] = "context not found in local kubeconfig"
			continue
		}
		server := ""
		if cl := apiCfg.Clusters[ctx.Cluster]; cl != nil {
			server = cl.Server
		}
		if existingByName[ctxName] {
			if skipConflicts {
				result.Skipped[ctxName] = "name already registered"
				continue
			}
			result.Failed[ctxName] = "name already registered"
			continue
		}
		if server != "" {
			if other := existingByServer[server]; other != "" {
				if skipConflicts {
					result.Skipped[ctxName] = fmt.Sprintf("same server as existing cluster %q", other)
					continue
				}
			}
		}

		raw, err := extractContextKubeconfig(apiCfg, ctxName)
		if err != nil {
			result.Failed[ctxName] = err.Error()
			continue
		}
		b64 := base64.StdEncoding.EncodeToString(raw)
		if err := s.CreateCluster(models.CreateClusterRequest{
			Name:           ctxName,
			KubeconfigData: b64,
			Provider:       "kubernetes",
			Description:    fmt.Sprintf("Imported from %s context %s", path, ctxName),
			Environment:    "",
		}); err != nil {
			result.Failed[ctxName] = err.Error()
			continue
		}
		result.Imported = append(result.Imported, ctxName)
		existingByName[ctxName] = true
		if server != "" {
			existingByServer[server] = ctxName
		}
	}
	return result, nil
}

func extractContextKubeconfig(apiCfg *clientcmdapi.Config, contextName string) ([]byte, error) {
	ctx := apiCfg.Contexts[contextName]
	if ctx == nil {
		return nil, fmt.Errorf("context %q not found", contextName)
	}
	cluster := apiCfg.Clusters[ctx.Cluster]
	if cluster == nil {
		return nil, fmt.Errorf("cluster %q for context %q not found", ctx.Cluster, contextName)
	}
	auth := apiCfg.AuthInfos[ctx.AuthInfo]
	if auth == nil {
		return nil, fmt.Errorf("user %q for context %q not found", ctx.AuthInfo, contextName)
	}

	out := clientcmdapi.NewConfig()
	out.APIVersion = apiCfg.APIVersion
	out.Kind = apiCfg.Kind
	out.CurrentContext = contextName
	out.Contexts[contextName] = ctx.DeepCopy()
	out.Clusters[ctx.Cluster] = cluster.DeepCopy()
	out.AuthInfos[ctx.AuthInfo] = auth.DeepCopy()
	return clientcmd.Write(*out)
}
