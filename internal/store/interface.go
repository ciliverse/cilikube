package store

// ClusterStore 定义了与集群数据持久化存储交互所需的所有方法。
type ClusterStore interface {
	CreateCluster(cluster *Cluster) error
	GetClusterByID(id string) (*Cluster, error)
	GetClusterByName(name string) (*Cluster, error)
	GetAllClusters() ([]Cluster, error)
	UpdateCluster(cluster *Cluster) error
	DeleteClusterByName(name string) error
	DeleteClusterByID(id string) error
}
