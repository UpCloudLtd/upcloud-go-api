package interfaces

import (
	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
)

type Host interface {
	GetHosts() (*upcloud.Hosts, error)
	GetHostDetails(r *request.GetHostDetailsRequest) (*upcloud.Host, error)
	ModifyHost(r *request.ModifyHostRequest) (*upcloud.Host, error)
}
