package service

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ciliverse/cilikube/pkg/k8s"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// HelmService wraps the helm CLI for release management.
type HelmService struct {
	clusterManager *k8s.ClusterManager
}

func NewHelmService(cm *k8s.ClusterManager) *HelmService {
	return &HelmService{clusterManager: cm}
}

type HelmRelease struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Revision   string `json:"revision"`
	Updated    string `json:"updated"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
}

type HelmInstallRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Chart     string `json:"chart"`
	Repo      string `json:"repo,omitempty"`
	Version   string `json:"version,omitempty"`
	Values    string `json:"values,omitempty"`
	CreateNS  bool   `json:"createNamespace,omitempty"`
}

type HelmUpgradeRequest struct {
	Chart   string `json:"chart"`
	Repo    string `json:"repo,omitempty"`
	Version string `json:"version,omitempty"`
	Values  string `json:"values,omitempty"`
}

func (s *HelmService) withKubeconfig(clusterID string, fn func(kubeconfigPath string) error) error {
	client, err := s.clientFor(clusterID)
	if err != nil {
		return err
	}
	cfg := client.Config
	apiCfg := clientcmdapi.NewConfig()
	clusterName := "cilikube"
	contextName := "cilikube"
	apiCfg.Clusters[clusterName] = &clientcmdapi.Cluster{
		Server:                   cfg.Host,
		CertificateAuthorityData: cfg.CAData,
		InsecureSkipTLSVerify:    cfg.Insecure,
	}
	authName := "cilikube-user"
	apiCfg.AuthInfos[authName] = &clientcmdapi.AuthInfo{
		Token:                 cfg.BearerToken,
		ClientCertificateData: cfg.CertData,
		ClientKeyData:         cfg.KeyData,
		Username:              cfg.Username,
		Password:              cfg.Password,
	}
	apiCfg.Contexts[contextName] = &clientcmdapi.Context{
		Cluster:  clusterName,
		AuthInfo: authName,
	}
	apiCfg.CurrentContext = contextName

	tmp, err := os.CreateTemp("", "cilikube-helm-*.kubeconfig")
	if err != nil {
		return err
	}
	path := tmp.Name()
	_ = tmp.Close()
	defer os.Remove(path)

	if err := clientcmd.WriteToFile(*apiCfg, path); err != nil {
		return err
	}
	return fn(path)
}

func (s *HelmService) clientFor(clusterID string) (*k8s.Client, error) {
	if clusterID != "" {
		return s.clusterManager.GetClientByID(clusterID)
	}
	return s.clusterManager.GetActiveClient()
}

func (s *HelmService) runHelm(kubeconfig string, args ...string) ([]byte, error) {
	if _, err := exec.LookPath("helm"); err != nil {
		return nil, fmt.Errorf("helm CLI not found on API host: %w", err)
	}
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	cmd := exec.Command("helm", full...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, fmt.Errorf("helm %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return out, nil
}

func (s *HelmService) ListReleases(clusterID, namespace string) ([]HelmRelease, error) {
	var releases []HelmRelease
	err := s.withKubeconfig(clusterID, func(kc string) error {
		args := []string{"list", "-o", "json"}
		if namespace != "" {
			args = append(args, "-n", namespace)
		} else {
			args = append(args, "-A")
		}
		out, err := s.runHelm(kc, args...)
		if err != nil {
			return err
		}
		if strings.TrimSpace(string(out)) == "" || string(out) == "null" {
			releases = []HelmRelease{}
			return nil
		}
		return json.Unmarshal(out, &releases)
	})
	return releases, err
}

func (s *HelmService) GetRelease(clusterID, namespace, name string) (map[string]any, error) {
	var result map[string]any
	err := s.withKubeconfig(clusterID, func(kc string) error {
		out, err := s.runHelm(kc, "get", "all", name, "-n", namespace, "-o", "json")
		if err != nil {
			// fallback: status
			out, err = s.runHelm(kc, "status", name, "-n", namespace, "-o", "json")
			if err != nil {
				return err
			}
		}
		return json.Unmarshal(out, &result)
	})
	return result, err
}

func (s *HelmService) Install(clusterID string, req HelmInstallRequest) (string, error) {
	if req.Name == "" || req.Chart == "" || req.Namespace == "" {
		return "", fmt.Errorf("name, chart, and namespace are required")
	}
	var output string
	err := s.withKubeconfig(clusterID, func(kc string) error {
		if req.Repo != "" {
			repoName := "cilikube-" + strings.ReplaceAll(req.Name, "/", "-")
			if _, err := s.runHelm(kc, "repo", "add", repoName, req.Repo); err != nil {
				// ignore if already exists
				if !strings.Contains(err.Error(), "already exists") {
					return err
				}
			}
			_, _ = s.runHelm(kc, "repo", "update")
		}
		args := []string{"install", req.Name, req.Chart, "-n", req.Namespace}
		if req.CreateNS {
			args = append(args, "--create-namespace")
		}
		if req.Version != "" {
			args = append(args, "--version", req.Version)
		}
		var valuesFile string
		if req.Values != "" {
			f, err := os.CreateTemp("", "cilikube-values-*.yaml")
			if err != nil {
				return err
			}
			valuesFile = f.Name()
			if _, err := f.WriteString(req.Values); err != nil {
				f.Close()
				os.Remove(valuesFile)
				return err
			}
			f.Close()
			defer os.Remove(valuesFile)
			args = append(args, "-f", valuesFile)
		}
		out, err := s.runHelm(kc, args...)
		output = string(out)
		return err
	})
	return output, err
}

func (s *HelmService) Upgrade(clusterID, namespace, name string, req HelmUpgradeRequest) (string, error) {
	if req.Chart == "" {
		return "", fmt.Errorf("chart is required")
	}
	var output string
	err := s.withKubeconfig(clusterID, func(kc string) error {
		args := []string{"upgrade", name, req.Chart, "-n", namespace}
		if req.Version != "" {
			args = append(args, "--version", req.Version)
		}
		if req.Values != "" {
			f, err := os.CreateTemp("", "cilikube-values-*.yaml")
			if err != nil {
				return err
			}
			path := f.Name()
			if _, err := f.WriteString(req.Values); err != nil {
				f.Close()
				os.Remove(path)
				return err
			}
			f.Close()
			defer os.Remove(path)
			args = append(args, "-f", path)
		}
		out, err := s.runHelm(kc, args...)
		output = string(out)
		return err
	})
	return output, err
}

func (s *HelmService) Rollback(clusterID, namespace, name, revision string) (string, error) {
	var output string
	err := s.withKubeconfig(clusterID, func(kc string) error {
		args := []string{"rollback", name, "-n", namespace}
		if revision != "" {
			args = append(args, revision)
		}
		out, err := s.runHelm(kc, args...)
		output = string(out)
		return err
	})
	return output, err
}

func (s *HelmService) Uninstall(clusterID, namespace, name string) (string, error) {
	var output string
	err := s.withKubeconfig(clusterID, func(kc string) error {
		out, err := s.runHelm(kc, "uninstall", name, "-n", namespace)
		output = string(out)
		return err
	})
	return output, err
}
