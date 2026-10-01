package k8s

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// KubeconfigFromREST writes a one-context kubeconfig that kubectl can use.
// The bytes contain the cluster credential; callers must not log them.
func KubeconfigFromREST(cfg *rest.Config) ([]byte, error) {
	cluster := &clientcmdapi.Cluster{
		Server:                   cfg.Host,
		InsecureSkipTLSVerify:    cfg.Insecure,
		CertificateAuthorityData: cfg.CAData,
	}
	if cfg.CAFile != "" && len(cfg.CAData) == 0 {
		cluster.CertificateAuthority = cfg.CAFile
	}
	auth := &clientcmdapi.AuthInfo{
		Token:                 cfg.BearerToken,
		Username:              cfg.Username,
		Password:              cfg.Password,
		ClientCertificateData: cfg.CertData,
		ClientKeyData:         cfg.KeyData,
	}
	if cfg.CertFile != "" && len(cfg.CertData) == 0 {
		auth.ClientCertificate = cfg.CertFile
	}
	if cfg.KeyFile != "" && len(cfg.KeyData) == 0 {
		auth.ClientKey = cfg.KeyFile
	}
	raw := clientcmdapi.NewConfig()
	raw.Clusters["cilikube"] = cluster
	raw.AuthInfos["cilikube"] = auth
	raw.Contexts["cilikube"] = &clientcmdapi.Context{
		Cluster:  "cilikube",
		AuthInfo: "cilikube",
	}
	raw.CurrentContext = "cilikube"
	return clientcmd.Write(*raw)
}
