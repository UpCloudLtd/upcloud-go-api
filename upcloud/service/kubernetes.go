package service

import (
	"context"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

type Kubernetes interface {
	GetKubernetesClusters(ctx context.Context, r *request.GetKubernetesClustersRequest) ([]upcloud.KubernetesCluster, error)
	GetKubernetesCluster(ctx context.Context, r *request.GetKubernetesClusterRequest) (*upcloud.KubernetesCluster, error)
	CreateKubernetesCluster(ctx context.Context, r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error)
	DeleteKubernetesCluster(ctx context.Context, r *request.DeleteKubernetesClusterRequest) error
	GetKubernetesKubeconfig(ctx context.Context, r *request.GetKubernetesKubeconfigRequest) (string, error)
	GetKubernetesVersions(ctx context.Context, r *request.GetKubernetesVersionsRequest) ([]string, error)
	WaitForKubernetesClusterState(ctx context.Context, r *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error)
	GetKubernetesNodeGroups(ctx context.Context, r *request.GetKubernetesNodeGroupsRequest) ([]upcloud.KubernetesNodeGroup, error)
	GetKubernetesNodeGroup(ctx context.Context, r *request.GetKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error)
	CreateKubernetesNodeGroup(ctx context.Context, r *request.CreateKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error)
	ModifyKubernetesNodeGroup(ctx context.Context, r *request.ModifyKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error)
	DeleteKubernetesNodeGroup(ctx context.Context, r *request.DeleteKubernetesNodeGroupRequest) error
}

// GetKubernetesClusters retrieves a list of Kubernetes clusters (EXPERIMENTAL).
func (s *Service) GetKubernetesClusters(ctx context.Context, r *request.GetKubernetesClustersRequest) ([]upcloud.KubernetesCluster, error) {
	clusters := make([]upcloud.KubernetesCluster, 0)
	return clusters, s.get(ctx, r.RequestURL(), &clusters)
}

// GetKubernetesCluster retrieves details of a Kubernetes cluster (EXPERIMENTAL).
func (s *Service) GetKubernetesCluster(ctx context.Context, r *request.GetKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	kubernetesCluster := upcloud.KubernetesCluster{}
	return &kubernetesCluster, s.get(ctx, r.RequestURL(), &kubernetesCluster)
}

// CreateKubernetesCluster creates a new Kubernetes cluster (EXPERIMENTAL).
func (s *Service) CreateKubernetesCluster(ctx context.Context, r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	if r == nil || len(r.Network) == 0 {
		return nil, fmt.Errorf("bad request")
	}

	networkDetails, err := s.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: r.Network})

	if err != nil || networkDetails == nil || len(networkDetails.IPNetworks) == 0 {
		return nil, fmt.Errorf("invalid network: %w", err)
	}

	r.NetworkCIDR = networkDetails.IPNetworks[0].Address

	cluster := upcloud.KubernetesCluster{}

	return &cluster, s.create(ctx, r, &cluster)
}

// DeleteKubernetesCluster deletes an existing Kubernetes cluster (EXPERIMENTAL).
func (s *Service) DeleteKubernetesCluster(ctx context.Context, r *request.DeleteKubernetesClusterRequest) error {
	return s.delete(ctx, r)
}

// WaitForKubernetesClusterState (EXPERIMENTAL) blocks execution until the specified Kubernetes cluster has entered the
// specified state. If the state changes favorably, cluster details is returned. The method will give up
// after the specified timeout
func (s *Service) WaitForKubernetesClusterState(ctx context.Context, r *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error) {
	attempts := 0
	sleepDuration := time.Second * 5

	for {
		attempts++

		details, err := s.GetKubernetesCluster(ctx, &request.GetKubernetesClusterRequest{
			UUID: r.UUID,
		})
		if err != nil {
			return nil, err
		}

		if details.State == r.DesiredState {
			return details, nil
		}

		time.Sleep(sleepDuration)

		if time.Duration(attempts)*sleepDuration >= r.Timeout {
			return nil, fmt.Errorf("timeout reached while waiting for Kubernetes cluster to enter state \"%s\"", r.DesiredState)
		}
	}
}

// GetKubernetesKubeconfig retrieves kubeconfig of a Kubernetes cluster (EXPERIMENTAL).
func (s *Service) GetKubernetesKubeconfig(ctx context.Context, r *request.GetKubernetesKubeconfigRequest) (string, error) {
	// TODO: should timeout be part of GetKubernetesKubeconfigRequest ?
	const timeout time.Duration = 10 * time.Minute
	data := struct {
		Kubeconfig string `json:"kubeconfig"`
	}{}

	_, err := s.WaitForKubernetesClusterState(ctx, &request.WaitForKubernetesClusterStateRequest{
		DesiredState: upcloud.KubernetesClusterStateRunning,
		Timeout:      timeout,
		UUID:         r.UUID,
	})
	if err != nil {
		return "", err
	}

	err = s.get(ctx, r.RequestURL(), &data)
	return data.Kubeconfig, err
}

// GetKubernetesVersions retrieves a list of Kubernetes cluster versions (EXPERIMENTAL).
func (s *Service) GetKubernetesVersions(ctx context.Context, r *request.GetKubernetesVersionsRequest) ([]string, error) {
	versions := make([]string, 0)
	return versions, s.get(ctx, r.RequestURL(), &versions)
}

// GetKubernetesNodeGroups retrieves a list of Kubernetes cluster node groups.
func (s *Service) GetKubernetesNodeGroups(ctx context.Context, r *request.GetKubernetesNodeGroupsRequest) ([]upcloud.KubernetesNodeGroup, error) {
	ng := make([]upcloud.KubernetesNodeGroup, 0)
	return ng, s.get(ctx, r.RequestURL(), &ng)
}

// GetKubernetesNodeGroup retrieves details of a node group.
func (s *Service) GetKubernetesNodeGroup(ctx context.Context, r *request.GetKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error) {
	ng := upcloud.KubernetesNodeGroup{}
	return &ng, s.get(ctx, r.RequestURL(), &ng)
}

// CreateKubernetesNodeGroup creates a new node group.
func (s *Service) CreateKubernetesNodeGroup(ctx context.Context, r *request.CreateKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error) {
	ng := upcloud.KubernetesNodeGroup{}
	return &ng, s.create(ctx, r, &ng)
}

// ModifyKubernetesNodeGroup modifies an existing node group.
func (s *Service) ModifyKubernetesNodeGroup(ctx context.Context, r *request.ModifyKubernetesNodeGroupRequest) (*upcloud.KubernetesNodeGroup, error) {
	ng := upcloud.KubernetesNodeGroup{}
	return &ng, s.modify(ctx, r, &ng)
}

// DeleteKubernetesNodeGroup deletes an existing node group.
func (s *Service) DeleteKubernetesNodeGroup(ctx context.Context, r *request.DeleteKubernetesNodeGroupRequest) error {
	return s.delete(ctx, r)
}
