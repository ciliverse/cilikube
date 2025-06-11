package service

// ProxyService 结构体不再持有 restConfig 字段
type ProxyService struct {
	// 不需要 restConfig *rest.Config 字段了
}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

//func (s *ProxyService) GetConfig() *rest.Config {
//	return s.restConfig
//}
