package service

// ProxyService struct no longer holds restConfig field
type ProxyService struct {
	// No longer need restConfig *rest.Config field
}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

//func (s *ProxyService) GetConfig() *rest.Config {
//	return s.restConfig
//}
