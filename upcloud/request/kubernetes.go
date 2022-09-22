package request

import (
	"fmt"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
)

const (
	kubernetesClusterBasePath string = "/kubernetes"
)

// GetKubernetesClustersRequest represents a request to list Kubernetes clusters
type GetKubernetesClustersRequest struct{}

func (r *GetKubernetesClustersRequest) RequestURL() string {
	return kubernetesClusterBasePath
}

type KubernetesFilter interface {
	ToQueryParam() string
}

// GetKubernetesClustersWithFiltersRequest represents a request to get all clusters
// using labels or label keys as filters.
// Using multiple filters returns only clusters that match all.
type GetKubernetesClustersWithFiltersRequest struct {
	Filters []KubernetesFilter
}

// RequestURL implements the Request interface.
func (r *GetKubernetesClustersWithFiltersRequest) RequestURL() string {
	if len(r.Filters) == 0 {
		return kubernetesClusterBasePath
	}

	params := ""
	for _, v := range r.Filters {
		if len(params) > 0 {
			params += "&"
		}
		params += v.ToQueryParam()
	}

	return fmt.Sprintf("%s?%s", kubernetesClusterBasePath, params)
}

// GetKubernetesClusterRequest represents a request to get a Kubernetes cluster details
type GetKubernetesClusterRequest struct {
	UUID string
}

func (r *GetKubernetesClusterRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", kubernetesClusterBasePath, r.UUID)
}

// CreateKubernetesClusterRequest represents a request to create a Kubernetes cluster
type CreateKubernetesClusterRequest struct {
	Name       string                        `json:"title"`
	Network    string                        `json:"network"`
	NodeGroups []upcloud.KubernetesNodeGroup `json:"node_groups"`
	Zone       string                        `json:"zone"`
}

func (r *CreateKubernetesClusterRequest) RequestURL() string {
	return kubernetesClusterBasePath
}

// DeleteKubernetesClusterRequest represents a request to delete a Kubernetes cluster
type DeleteKubernetesClusterRequest struct {
	UUID string `json:"-"`
}

func (r *DeleteKubernetesClusterRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", kubernetesClusterBasePath, r.UUID)
}

// WaitForKubernetesClusterStateRequest represents a request to wait for a Kubernetes cluster
// to enter a desired state
type WaitForKubernetesClusterStateRequest struct {
	DesiredState upcloud.KubernetesClusterState `json:"-"`
	Timeout      time.Duration                  `json:"-"`
	UUID         string                         `json:"-"`
}

func (r *WaitForKubernetesClusterStateRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s", kubernetesClusterBasePath, r.UUID)
}

// GetKubernetesKubeconfigRequest represents a request to get kubeconfig for a Kubernetes cluster
type GetKubernetesKubeconfigRequest struct {
	UUID string `json:"-"`
}

func (r *GetKubernetesKubeconfigRequest) RequestURL() string {
	return fmt.Sprintf("%s/%s/kubeconfig", kubernetesClusterBasePath, r.UUID)
}

// GetKubernetesPlansRequest represents a request to list available Kubernetes cluster plans
type GetKubernetesPlansRequest struct{}

func (r *GetKubernetesPlansRequest) RequestURL() string {
	return fmt.Sprintf("%s/plan", kubernetesClusterBasePath)
}

// GetKubernetesVersionsRequest represents a request to list available Kubernetes cluster versions
type GetKubernetesVersionsRequest struct{}

func (r *GetKubernetesVersionsRequest) RequestURL() string {
	return fmt.Sprintf("%s/version", kubernetesClusterBasePath)
}
