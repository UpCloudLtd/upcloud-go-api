package upcloud

const (
	KuberetesClusterStateConfiguring KubernetesClusterState = "configuring"
	KuberetesClusterStateReady       KubernetesClusterState = "ready"
	KubernetesClusterTypeStandalone  KubernetesClusterType  = "standalone"
)

type (
	KubernetesClusterState string
	KubernetesClusterType  string
)

type KubernetesCluster struct {
	Name        string                 `json:"title"`
	Network     string                 `json:"network"`
	NetworkCIDR string                 `json:"network_cidr"`
	NodeGroups  []KubernetesNodeGroup  `json:"node_groups"`
	State       KubernetesClusterState `json:"state"`
	Storage     string                 `json:"storage"`
	Type        KubernetesClusterType  `json:"type"`
	UUID        string                 `json:"uuid"`
	Zone        string                 `json:"zone"`
}

type KubernetesNodeGroup struct {
	Count   int      `json:"count,omitempty"`
	Labels  []Label  `json:"labels,omitempty"`
	Name    string   `json:"name,omitempty"`
	Plan    string   `json:"plan,omitempty"`
	SSHKeys []string `json:"ssh_key,omitempty"`
}

type KubernetesPlan struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
