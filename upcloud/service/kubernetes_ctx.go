package service

import (
	"context"
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type KubernetesContext interface {
	GetKubernetesClusters(ctx context.Context, r *request.GetKubernetesClustersRequest) ([]upcloud.KubernetesCluster, error)
	GetKubernetesCluster(ctx context.Context, r *request.GetKubernetesClusterRequest) (*upcloud.KubernetesCluster, error)
	CreateKubernetesCluster(ctx context.Context, r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error)
	DeleteKubernetesCluster(ctx context.Context, r *request.DeleteKubernetesClusterRequest) error
	GetKubernetesKubeconfig(ctx context.Context, r *request.GetKubernetesKubeconfigRequest) (string, error)
	GetKubernetesPlans(ctx context.Context, r *request.GetKubernetesPlansRequest) ([]upcloud.KubernetesPlan, error)
	GetKubernetesVersions(ctx context.Context, r *request.GetKubernetesVersionsRequest) ([]string, error)
	WaitForKubernetesClusterState(ctx context.Context, r *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error)
}

var _ KubernetesContext = (*ServiceContext)(nil)

// GetKubernetesClusters retrieves a list of Kubernetes clusters (EXPERIMENTAL).
func (s *ServiceContext) GetKubernetesClusters(ctx context.Context, r *request.GetKubernetesClustersRequest) ([]upcloud.KubernetesCluster, error) {
	clusters := make([]upcloud.KubernetesCluster, 0)
	return clusters, s.get(ctx, r.RequestURL(), &clusters)
}

// GetKubernetesCluster retrieves details of a Kubernetes cluster (EXPERIMENTAL).
func (s *ServiceContext) GetKubernetesCluster(ctx context.Context, r *request.GetKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	kubernetesCluster := upcloud.KubernetesCluster{}
	return &kubernetesCluster, s.get(ctx, r.RequestURL(), &kubernetesCluster)
}

// CreateKubernetesCluster creates a new Kubernetes cluster (EXPERIMENTAL).
func (s *ServiceContext) CreateKubernetesCluster(ctx context.Context, r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	if r == nil || len(r.Network) == 0 {
		return nil, fmt.Errorf("bad request")
	}

	networkDetails, err := s.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: r.Network})

	if err != nil || networkDetails == nil || len(networkDetails.IPNetworks) == 0 {
		return nil, fmt.Errorf("invalid network: %w", err)
	}

	r.NetworkCIDR = networkDetails.IPNetworks[0].Address

	_, err = s.GetStorageDetails(ctx, &request.GetStorageDetailsRequest{UUID: r.Storage})
	if err != nil {
		return nil, fmt.Errorf("storage does not exist: %w", err)
	}

	cluster := upcloud.KubernetesCluster{}

	return &cluster, s.create(ctx, r, &cluster)
}

// DeleteKubernetesCluster deletes an existing Kubernetes cluster (EXPERIMENTAL).
func (s *ServiceContext) DeleteKubernetesCluster(ctx context.Context, r *request.DeleteKubernetesClusterRequest) error {
	return s.delete(ctx, r)
}

// WaitForKubernetesClusterState (EXPERIMENTAL) blocks execution until the specified Kubernetes cluster has entered the
// specified state. If the state changes favorably, cluster details is returned. The method will give up
// after the specified timeout
func (s *ServiceContext) WaitForKubernetesClusterState(ctx context.Context, r *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error) {
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

// GetKubernetesPlans retrieves a list of Kubernetes cluster plans (EXPERIMENTAL).
func (s *ServiceContext) GetKubernetesPlans(ctx context.Context, r *request.GetKubernetesPlansRequest) ([]upcloud.KubernetesPlan, error) {
	plans := make([]upcloud.KubernetesPlan, 0)
	return plans, s.get(ctx, r.RequestURL(), &plans)
}

// GetKubernetesKubeconfig retrieves kubeconfig of a Kubernetes cluster (EXPERIMENTAL).
func (s *ServiceContext) GetKubernetesKubeconfig(ctx context.Context, r *request.GetKubernetesKubeconfigRequest) (string, error) {
	data := struct {
		Kubeconfig string `json:"kubeconfig"`
	}{}

	_, err := s.WaitForKubernetesClusterState(ctx, &request.WaitForKubernetesClusterStateRequest{
		DesiredState: upcloud.KubernetesClusterStateRunning,
		Timeout:      s.client.GetTimeout(),
		UUID:         r.UUID,
	})
	if err != nil {
		return "", err
	}

	err = s.get(ctx, r.RequestURL(), &data)
	return data.Kubeconfig, err
}

// GetKubernetesVersions retrieves a list of Kubernetes cluster versions (EXPERIMENTAL).
func (s *ServiceContext) GetKubernetesVersions(ctx context.Context, r *request.GetKubernetesVersionsRequest) ([]string, error) {
	versions := make([]string, 0)
	return versions, s.get(ctx, r.RequestURL(), &versions)
}
