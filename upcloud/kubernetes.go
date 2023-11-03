package upcloud

const (
	KubernetesClusterStatePending     KubernetesClusterState = "pending"
	KubernetesClusterStateRunning     KubernetesClusterState = "running"
	KubernetesClusterStateTerminating KubernetesClusterState = "terminating"
	KubernetesClusterStateTerminated  KubernetesClusterState = "terminated"
	KubernetesClusterStateFailed      KubernetesClusterState = "failed"
	KubernetesClusterStateUnknown     KubernetesClusterState = "unknown"

	KubernetesNodeGroupStatePending     KubernetesNodeGroupState = "pending"
	KubernetesNodeGroupStateRunning     KubernetesNodeGroupState = "running"
	KubernetesNodeGroupStateScalingUp   KubernetesNodeGroupState = "scaling-up"
	KubernetesNodeGroupStateScalingDown KubernetesNodeGroupState = "scaling-down"
	KubernetesNodeGroupStateTerminating KubernetesNodeGroupState = "terminating"
	KubernetesNodeGroupStateFailed      KubernetesNodeGroupState = "failed"
	KubernetesNodeGroupStateUnknown     KubernetesNodeGroupState = "unknown"

	KubernetesClusterTaintEffectNoExecute        KubernetesClusterTaintEffect = "NoExecute"
	KubernetesClusterTaintEffectNoSchedule       KubernetesClusterTaintEffect = "NoSchedule"
	KubernetesClusterTaintEffectPreferNoSchedule KubernetesClusterTaintEffect = "PreferNoSchedule"

	KubernetesNodeStateFailed      KubernetesNodeState = "failed"
	KubernetesNodeStatePending     KubernetesNodeState = "pending"
	KubernetesNodeStateRunning     KubernetesNodeState = "running"
	KubernetesNodeStateTerminating KubernetesNodeState = "terminating"
	KubernetesNodeStateUnknown     KubernetesNodeState = "unknown"
)

type (
	KubernetesClusterState       string
	KubernetesNodeGroupState     string
	KubernetesClusterType        string
	KubernetesClusterTaintEffect string
	KubernetesNodeState          string
)

type KubernetesCluster struct {
	ControlPlaneIPFilter []string               `json:"control_plane_ip_filter"`
	Name                 string                 `json:"name"`
	Network              string                 `json:"network"`
	NetworkCIDR          string                 `json:"network_cidr"`
	NodeGroups           []KubernetesNodeGroup  `json:"node_groups"`
	State                KubernetesClusterState `json:"state"`
	UUID                 string                 `json:"uuid"`
	Version              string                 `json:"version"`
	Zone                 string                 `json:"zone"`
	Plan                 string                 `json:"plan"`
	PrivateNodeGroups    bool                   `json:"private_node_groups"`
}

type KubernetesNodeGroup struct {
	AntiAffinity         bool                     `json:"anti_affinity,omitempty"`
	Count                int                      `json:"count,omitempty"`
	KubeletArgs          []KubernetesKubeletArg   `json:"kubelet_args,omitempty"`
	Labels               []Label                  `json:"labels,omitempty"`
	Name                 string                   `json:"name,omitempty"`
	Plan                 string                   `json:"plan,omitempty"`
	SSHKeys              []string                 `json:"ssh_keys,omitempty"`
	State                KubernetesNodeGroupState `json:"state,omitempty"`
	Storage              string                   `json:"storage,omitempty"`
	Taints               []KubernetesTaint        `json:"taints,omitempty"`
	UtilityNetworkAccess bool                     `json:"utility_network_access,omitempty"`
}

type KubernetesNodeGroupDetails struct {
	KubernetesNodeGroup

	Nodes []KubernetesNode `json:"nodes,omitempty"`
}

type KubernetesNode struct {
	UUID  string              `json:"uuid,omitempty"`
	Name  string              `json:"name,omitempty"`
	State KubernetesNodeState `json:"state,omitempty"`
}

type KubernetesKubeletArg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KubernetesTaint struct {
	Effect KubernetesClusterTaintEffect `json:"effect"`
	Key    string                       `json:"key"`
	Value  string                       `json:"value"`
}

type KubernetesPlan struct {
	Name         string `json:"name"`
	ServerNumber int    `json:"server_number"`
	MaxNodes     int    `json:"max_nodes"`
}

type KubernetesVersion struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}
