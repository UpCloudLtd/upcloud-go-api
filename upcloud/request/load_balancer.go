package request

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/google/uuid"
)

type GetLoadBalancersRequest struct{}

type GetLoadBalancerDetailsRequest struct {
	UUID string
}

type CreateLoadBalancerRequest struct {
	Name             string             `json:"name"`
	Plan             string             `json:"plan"`
	Zone             string             `json:"zone"`
	NetworkUuid      uuid.UUID          `json:"network_uuid"`
	ConfiguredStatus string             `json:"configured_status"`
	Frontends        []upcloud.Frontend `json:"frontends"`
	Backends         []upcloud.Backend  `json:"backends"`
	// Resolvers        []*upcloud.Resolver `json:"resolvers"` // TODO explore omit empty
}

type ModifyLoadBalancerRequest struct {
	UUID uuid.UUID `json:"-"`

	Name             string `json:"name,omitempty"`
	Plan             string `json:"plan,omitempty"`
	ConfiguredStatus string `json:"configured_status,omitempty"`
}

type DeleteLoadBalancerRequest struct {
	UUID uuid.UUID `json:"-"`
}

func (r *GetLoadBalancersRequest) RequestURL() string {
	return "/loadbalancer"
}

func (r *GetLoadBalancerDetailsRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s", r.UUID)
}

func (r *CreateLoadBalancerRequest) RequestURL() string {
	return "/loadbalancer"
}

func (r *ModifyLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s", r.UUID)
}

func (r *DeleteLoadBalancerRequest) RequestURL() string {
	return fmt.Sprintf("/loadbalancer/%s", r.UUID)
}
