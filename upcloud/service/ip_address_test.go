package service

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetIPAddresses performs the following actions:
//   - creates a server
//   - retrieves all IP addresses
//   - compares the retrieved IP addresses with the created server's
//     ip addresses
func TestGetIPAddresses(t *testing.T) {
	record(t, "getipaddresses", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		serverDetails, err := createServerWithRecorder(rec, svc, "TestGetIPAddresses")
		require.NoError(t, err)
		assert.Greater(t, len(serverDetails.IPAddresses), 0)

		ipAddresses, err := svc.GetIPAddresses()
		require.NoError(t, err)
		var foundCount int
		for _, sip := range serverDetails.IPAddresses {
			for _, gip := range ipAddresses.IPAddresses {
				if sip.Address == gip.Address {
					foundCount++
					assert.Equal(t, sip.Access, gip.Access)
					assert.Equal(t, sip.Family, gip.Family)
					break
				}
			}
		}
		assert.Equal(t, len(serverDetails.IPAddresses), foundCount)

		for _, ip := range serverDetails.IPAddresses {
			require.NotEmpty(t, ip.Address)
			ipAddress, err := svc.GetIPAddressDetails(&request.GetIPAddressDetailsRequest{
				Address: ip.Address,
			})
			require.NoError(t, err)

			assert.Equal(t, ip.Address, ipAddress.Address)
			assert.Equal(t, ip.Access, ipAddress.Access)
			assert.Equal(t, ip.Family, ipAddress.Family)
		}
	})
}

// TestAttachModifyReleaseIPAddress performs the following actions
//
// - creates a server
// - assigns an additional IP address to it
// - modifies the PTR record of the IP address
// - deletes the IP address
func TestAttachModifyReleaseIPAddress(t *testing.T) {
	t.Parallel()

	record(t, "attachmodifyreleaseipaddress", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create the server
		serverDetails, err := createServerWithRecorder(rec, svc, "TestAttachModifyReleaseIPAddress")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", serverDetails.UUID)
		err = stopServerWithRecorder(rec, svc, serverDetails.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Assign an IP address
		t.Log("Assigning IP address to server ...")
		ipAddress, err := svc.AssignIPAddress(&request.AssignIPAddressRequest{
			Access:     upcloud.IPAddressAccessPublic,
			Family:     upcloud.IPAddressFamilyIPv6,
			ServerUUID: serverDetails.UUID,
		})
		require.NoError(t, err)
		t.Logf("Assigned IP address %s to server with UUID %s", ipAddress.Address, serverDetails.UUID)

		// Modify the PTR record
		t.Logf("Modifying PTR record for address %s ...", ipAddress.Address)
		ipAddress, err = svc.ModifyIPAddress(&request.ModifyIPAddressRequest{
			IPAddress: ipAddress.Address,
			PTRRecord: "such.pointer.example.com",
		})
		require.NoError(t, err)
		t.Logf("PTR record modified, new record is %s", ipAddress.PTRRecord)

		// Release the IP address
		t.Log("Releasing the IP address ...")
		err = svc.ReleaseIPAddress(&request.ReleaseIPAddressRequest{
			IPAddress: ipAddress.Address,
		})
		require.NoError(t, err)
		t.Log("The IP address is now released")
	})
}

func TestAttachModifyReleaseFloatingIPAddress(t *testing.T) {
	record(t, "attachmodifyreleasefloatingipaddress", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create the first server
		serverDetails1, err := createServerWithRecorder(rec, svc, "TestAttachModifyReleaseIPAddress1")
		require.NoError(t, err)
		t.Logf("Server 1 %s with UUID %s created", serverDetails1.Title, serverDetails1.UUID)

		// Create the second server
		serverDetails2, err := createServerWithRecorder(rec, svc, "TestAttachModifyReleaseIPAddress2")
		require.NoError(t, err)
		t.Logf("Server 2 %s with UUID %s created", serverDetails2.Title, serverDetails2.UUID)

		var mac string
		for _, ip := range serverDetails1.IPAddresses {
			if ip.Access == upcloud.IPAddressAccessPublic && ip.Family == upcloud.IPAddressFamilyIPv4 {
				ipDetails, err := svc.GetIPAddressDetails(&request.GetIPAddressDetailsRequest{
					Address: ip.Address,
				})
				require.NoError(t, err)
				mac = ipDetails.MAC
				break
			}
		}
		require.NotEmpty(t, mac)

		assignedIP, err := svc.AssignIPAddress(&request.AssignIPAddressRequest{
			Family:   upcloud.IPAddressFamilyIPv4,
			Floating: upcloud.True,
			MAC:      mac,
		})
		require.NoError(t, err)

		postAssignServerDetails1, err := svc.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: serverDetails1.UUID,
		})
		require.NoError(t, err)

		var found bool
		for _, inf := range postAssignServerDetails1.Networking.Interfaces {
			for _, ip := range inf.IPAddresses {
				if ip.Address == assignedIP.Address {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		assert.True(t, found)

		var mac2 string
		for _, ip := range serverDetails2.IPAddresses {
			if ip.Access == upcloud.IPAddressAccessPublic && ip.Family == upcloud.IPAddressFamilyIPv4 {
				ipDetails, err := svc.GetIPAddressDetails(&request.GetIPAddressDetailsRequest{
					Address: ip.Address,
				})
				require.NoError(t, err)
				mac2 = ipDetails.MAC
				break
			}
		}
		require.NotEmpty(t, mac2)

		_, err = svc.ModifyIPAddress(&request.ModifyIPAddressRequest{
			IPAddress: assignedIP.Address,
			MAC:       mac2,
		})
		require.NoError(t, err)

		postModifyServerDetails1, err := svc.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: serverDetails1.UUID,
		})
		require.NoError(t, err)

		found = false
		for _, inf := range postModifyServerDetails1.Networking.Interfaces {
			for _, ip := range inf.IPAddresses {
				if ip.Address == assignedIP.Address {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		assert.False(t, found)

		postModifyServerDetails2, err := svc.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: serverDetails2.UUID,
		})
		require.NoError(t, err)

		found = false
		for _, inf := range postModifyServerDetails2.Networking.Interfaces {
			for _, ip := range inf.IPAddresses {
				if ip.Address == assignedIP.Address {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		assert.True(t, found)

		// Unassign IP
		unassignIP, err := svc.ModifyIPAddress(&request.ModifyIPAddressRequest{
			IPAddress: assignedIP.Address,
		})
		require.NoError(t, err)
		assert.Empty(t, unassignIP.ServerUUID)
		assert.Empty(t, unassignIP.MAC)

		err = svc.ReleaseIPAddress(&request.ReleaseIPAddressRequest{
			IPAddress: assignedIP.Address,
		})
		require.NoError(t, err)
	})
}
