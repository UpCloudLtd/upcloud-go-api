package service

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type Kubernetes interface {
	CreateKubernetesCluster(r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error)
	DeleteKubernetesCluster(r *request.DeleteKubernetesClusterRequest) error
	GetKubernetesCluster(r *request.GetKubernetesClusterRequest) (*upcloud.KubernetesCluster, error)
	GetKubernetesClusters(r *request.GetKubernetesClustersRequest) ([]upcloud.KubernetesCluster, error)
	GetKubernetesKubeconfig(r *request.GetKubernetesKubeconfigRequest) (string, error)
	GetKubernetesPlans(r *request.GetKubernetesPlansRequest) ([]upcloud.KubernetesPlan, error)
	GetKubernetesVersions(r *request.GetKubernetesVersionsRequest) ([]string, error)
	WaitForKubernetesClusterState(r *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error)
}

var _ Kubernetes = (*Service)(nil)

// GetKubernetesClusters retrieves a list of Kubernetes clusters (EXPERIMENTAL).
func (s *Service) GetKubernetesClusters(r *request.GetKubernetesClustersRequest) ([]upcloud.KubernetesCluster, error) {
	clusters := make([]upcloud.KubernetesCluster, 0)
	return clusters, s.get(r.RequestURL(), &clusters)
}

// GetKubernetesCluster retrieves details of a Kubernetes cluster (EXPERIMENTAL).
func (s *Service) GetKubernetesCluster(r *request.GetKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	kubernetesCluster := upcloud.KubernetesCluster{}
	return &kubernetesCluster, s.get(r.RequestURL(), &kubernetesCluster)
}

// CreateKubernetesCluster creates a new Kubernetes cluster (EXPERIMENTAL).
func (s *Service) CreateKubernetesCluster(r *request.CreateKubernetesClusterRequest) (*upcloud.KubernetesCluster, error) {
	if r == nil || len(r.Network) == 0 {
		return nil, fmt.Errorf("bad request")
	}

	networkDetails, err := s.GetNetworkDetails(&request.GetNetworkDetailsRequest{UUID: r.Network})

	if err != nil || networkDetails == nil || len(networkDetails.IPNetworks) == 0 {
		return nil, fmt.Errorf("invalid network: %w", err)
	}

	r.NetworkCIDR = networkDetails.IPNetworks[0].Address

	_, err = s.GetStorageDetails(&request.GetStorageDetailsRequest{UUID: r.Storage})
	if err != nil {
		return nil, fmt.Errorf("invalid storage template uuid: %w", err)
	}

	cluster := upcloud.KubernetesCluster{}

	err = s.create(r, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, err
}

// DeleteKubernetesCluster deletes an existing Kubernetes cluster (EXPERIMENTAL).
func (s *Service) DeleteKubernetesCluster(r *request.DeleteKubernetesClusterRequest) error {
	if err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL())); err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}

// WaitForKubernetesClusterState (EXPERIMENTAL) blocks execution until the specified Kubernetes cluster has entered the
// specified state. If the state changes favorably, cluster details is returned. The method will give up
// after the specified timeout
func (s *Service) WaitForKubernetesClusterState(r *request.WaitForKubernetesClusterStateRequest) (*upcloud.KubernetesCluster, error) {
	attempts := 0
	sleepDuration := time.Second * 5

	for {
		attempts++

		details, err := s.GetKubernetesCluster(&request.GetKubernetesClusterRequest{
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
func (s *Service) GetKubernetesPlans(r *request.GetKubernetesPlansRequest) ([]upcloud.KubernetesPlan, error) {
	plans := make([]upcloud.KubernetesPlan, 0)
	return plans, s.get(r.RequestURL(), &plans)
}

// GetKubernetesKubeconfig retrieves kubeconfig of a Kubernetes cluster (EXPERIMENTAL).
func (s *Service) GetKubernetesKubeconfig(r *request.GetKubernetesKubeconfigRequest) (string, error) {
	var kubeconfig string

	_, err := s.WaitForKubernetesClusterState(&request.WaitForKubernetesClusterStateRequest{
		DesiredState: upcloud.KuberetesClusterStateReady,
		Timeout:      s.client.GetTimeout(),
		UUID:         r.UUID,
	})
	if err != nil {
		return kubeconfig, err
	}

	return kubeconfig, s.get(r.RequestURL(), &kubeconfig)
}

// GetKubernetesVersions retrieves a list of Kubernetes cluster versions (EXPERIMENTAL).
func (s *Service) GetKubernetesVersions(r *request.GetKubernetesVersionsRequest) ([]string, error) {
	versions := make([]string, 0)
	return versions, s.get(r.RequestURL(), &versions)
}
