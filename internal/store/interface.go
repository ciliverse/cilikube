package store

// ClusterStore defines all methods required for interacting with cluster data persistent storage.
type ClusterStore interface {
	CreateCluster(cluster *Cluster) error
	GetClusterByID(id string) (*Cluster, error)
	GetClusterByName(name string) (*Cluster, error)
	GetAllClusters() ([]Cluster, error)
	UpdateCluster(cluster *Cluster) error
	DeleteClusterByName(name string) error
	DeleteClusterByID(id string) error
}
