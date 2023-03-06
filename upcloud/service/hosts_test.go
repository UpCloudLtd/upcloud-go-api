package service

import (
	"context"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
)

// TestGetModifyHosts tests host functionality works correctly with context.
// The test:
//   - Gets all available hosts
//   - Gets the details of a single host and compares to details from all hosts list
//   - Modifies a host
//   - Modifies the host back
func TestGetModifyHosts(t *testing.T) {
	record(t, "getmodifyhosts", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		hosts, err := svc.GetHosts(ctx)
		require.NoError(t, err)

		assert.NotEmpty(t, hosts.Hosts)

		assert.NotZero(t, hosts.Hosts[0].ID)
		assert.NotEmpty(t, hosts.Hosts[0].Description)
		assert.NotEmpty(t, hosts.Hosts[0].Stats)
		assert.NotEmpty(t, hosts.Hosts[0].Zone)

		host, err := svc.GetHostDetails(ctx, &request.GetHostDetailsRequest{
			ID: hosts.Hosts[0].ID,
		})
		require.NoError(t, err)

		assert.Equal(t, hosts.Hosts[0].ID, host.ID)
		assert.Equal(t, hosts.Hosts[0].Description, host.Description)
		assert.Equal(t, hosts.Hosts[0].WindowsEnabled, host.WindowsEnabled)
		assert.Equal(t, hosts.Hosts[0].Zone, host.Zone)

		oldDescription := host.Description

		modifiedHost, err := svc.ModifyHost(ctx, &request.ModifyHostRequest{
			ID:          host.ID,
			Description: oldDescription + "(modified)",
		})
		require.NoError(t, err)
		assert.Equal(t, oldDescription+"(modified)", modifiedHost.Description)
		assert.Equal(t, host.ID, modifiedHost.ID)

		_, err = svc.ModifyHost(ctx, &request.ModifyHostRequest{
			ID:          host.ID,
			Description: oldDescription,
		})
		require.NoError(t, err)
	})
}
