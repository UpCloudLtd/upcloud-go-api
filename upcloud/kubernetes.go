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
	KubernetesNodeGroupStateTerminating KubernetesNodeGroupState = "terminating"
	KubernetesNodeGroupStateFailed      KubernetesNodeGroupState = "failed"
	KubernetesNodeGroupStateUnknown     KubernetesNodeGroupState = "unknown"

	KubernetesClusterTaintEffectNoExecute        KubernetesClusterTaintEffect = "NoExecute"
	KubernetesClusterTaintEffectNoSchedule       KubernetesClusterTaintEffect = "NoSchedule"
	KubernetesClusterTaintEffectPreferNoSchedule KubernetesClusterTaintEffect = "PreferNoSchedule"
)

type (
	KubernetesClusterState       string
	KubernetesNodeGroupState     string
	KubernetesClusterType        string
	KubernetesClusterTaintEffect string
)

type KubernetesCluster struct {
	Name        string                 `json:"name"`
	Network     string                 `json:"network"`
	NetworkCIDR string                 `json:"network_cidr"`
	NodeGroups  []KubernetesNodeGroup  `json:"node_groups"`
	State       KubernetesClusterState `json:"state"`
	UUID        string                 `json:"uuid"`
	Zone        string                 `json:"zone"`
}

type KubernetesNodeGroup struct {
	AntiAffinity bool                     `json:"anti_affinity,omitempty"`
	Count        int                      `json:"count,omitempty"`
	KubeletArgs  []KubernetesKubeletArg   `json:"kubelet_args,omitempty"`
	Labels       []Label                  `json:"labels,omitempty"`
	Name         string                   `json:"name,omitempty"`
	Plan         string                   `json:"plan,omitempty"`
	SSHKeys      []string                 `json:"ssh_keys,omitempty"`
	State        KubernetesNodeGroupState `json:"state,omitempty"`
	Storage      string                   `json:"storage,omitempty"`
	Taints       []KubernetesTaint        `json:"taints,omitempty"`
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
