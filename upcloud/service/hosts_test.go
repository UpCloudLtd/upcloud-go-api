package service_test

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetModifyHosts tests host functionality works correctly.
// The test:
//   - Gets all available hosts
//   - Gets the details of a single host and compares to details from all hosts list
//   - Modifies a host
//   - Modifies the host back
func TestGetModifyHosts(t *testing.T) {
	//nolint:thelper // false positive, the function is not a helper
	record(t, "getmodifyhosts", func(t *testing.T, svc *service.Service) {
		hosts, err := svc.GetHosts()
		require.NoError(t, err)

		assert.NotEmpty(t, hosts.Hosts)

		assert.NotZero(t, hosts.Hosts[0].ID)
		assert.NotEmpty(t, hosts.Hosts[0].Description)
		assert.NotEmpty(t, hosts.Hosts[0].Stats)
		assert.NotEmpty(t, hosts.Hosts[0].Zone)

		host, err := svc.GetHostDetails(&request.GetHostDetailsRequest{
			ID: hosts.Hosts[0].ID,
		})
		require.NoError(t, err)

		assert.Equal(t, hosts.Hosts[0].ID, host.ID)
		assert.Equal(t, hosts.Hosts[0].Description, host.Description)
		assert.Equal(t, hosts.Hosts[0].WindowsEnabled, host.WindowsEnabled)
		assert.Equal(t, hosts.Hosts[0].Zone, host.Zone)

		oldDescription := host.Description

		modifiedHost, err := svc.ModifyHost(&request.ModifyHostRequest{
			ID:          host.ID,
			Description: oldDescription + "(modified)",
		})
		require.NoError(t, err)
		assert.Equal(t, oldDescription+"(modified)", modifiedHost.Description)
		assert.Equal(t, host.ID, modifiedHost.ID)

		_, err = svc.ModifyHost(&request.ModifyHostRequest{
			ID:          host.ID,
			Description: oldDescription,
		})
		require.NoError(t, err)
	})
}
