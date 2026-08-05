package k8s

// SnapshotClients returns a copy of the current client map for background work.
func (cm *ClusterManager) SnapshotClients() map[string]*Client {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	out := make(map[string]*Client, len(cm.clients))
	for id, c := range cm.clients {
		out[id] = c
	}
	return out
}
