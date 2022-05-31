package service

import (
	"context"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

// TestGetNetworksContext checks that network details are retrievable
// It:
//    - creates a server
//    - Gets all networks and verifies details are populated
//    - checks that at least one network has a server in.
func TestGetNetworksContext(t *testing.T) {
	recordWithContext(t, "getnetworks", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		_, err := createServer(svc, "TestGetNetworks")
		require.NoError(t, err)

		networks, err := svcContext.GetNetworks(ctx)
		require.NoError(t, err)

		assert.NotEmpty(t, networks.Networks)

		assert.NotEmpty(t, networks.Networks[0].IPNetworks)
		assert.NotEmpty(t, networks.Networks[0].Name)
		assert.NotEmpty(t, networks.Networks[0].Type)
		assert.NotEmpty(t, networks.Networks[0].UUID)
		assert.NotEmpty(t, networks.Networks[0].Zone)

		// Find a network with a server
		var found bool
		for _, n := range networks.Networks {
			if len(n.Servers) > 0 {
				found = true
				break
			}
		}
		assert.True(t, found)
	})
}

// TestGetNetworksInZoneContext checks that network details in a zone are retrievable
// It:
//    - creates a server
//    - Gets all networks in a zone and verifies details are populated
//    - checks that at least one network has a server in.
func TestGetNetworksInZoneContext(t *testing.T) {
	recordWithContext(t, "getnetworksinzone", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		_, err := createServer(svc, "TestGetNetworksInZone")
		require.NoError(t, err)

		networks, err := svcContext.GetNetworksInZone(ctx, &request.GetNetworksInZoneRequest{
			Zone: "fi-hel2",
		})
		require.NoError(t, err)

		assert.NotEmpty(t, networks)

		assert.NotEmpty(t, networks.Networks[0].IPNetworks)
		assert.NotEmpty(t, networks.Networks[0].Name)
		assert.NotEmpty(t, networks.Networks[0].Type)
		assert.NotEmpty(t, networks.Networks[0].UUID)

		// Find a network with a server
		var found bool
		var foundNetwork *upcloud.Network
		for i, n := range networks.Networks {
			if len(n.Servers) > 0 {
				foundNetwork = &networks.Networks[i]
				found = true
			}
			// Make sure all the networks are in the right zone.
			assert.Equal(t, "fi-hel2", n.Zone)
		}
		assert.True(t, found)
		require.NotNil(t, foundNetwork)

		network, err := svcContext.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
			UUID: foundNetwork.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, foundNetwork, network)
	})
}

// TestCreateModifyDeleteNetworkContext checks that the network functionality works.
// It:
//    - creates a network
//    - modifies the network
//    - creates a server
//    - stops the server
//    - creates a network interface in the network
//    - modifies the network interface
//    - deletes the network interface
//    - deletes the network
//    - verifies the network has been deleted.
func TestCreateModifyDeleteNetworkContext(t *testing.T) {
	recordWithContext(t, "createmodifydeletenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		network, err := svcContext.CreateNetwork(ctx, &request.CreateNetworkRequest{
			Name: "test private network (test)",
			Zone: "fi-hel2",
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "172.16.0.0/22",
					DHCP:             upcloud.True,
					DHCPDefaultRoute: upcloud.False,
					DHCPDns: []string{
						"172.16.0.10",
						"172.16.1.10",
					},
					Family:  upcloud.IPAddressFamilyIPv4,
					Gateway: "172.16.0.1",
				},
			},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, network.UUID)
		assert.Equal(t, "test private network (test)", network.Name)

		postModifyNetwork, err := svcContext.ModifyNetwork(ctx, &request.ModifyNetworkRequest{
			UUID: network.UUID,
			Name: "modified private network (test)",
		})
		require.NoError(t, err)
		assert.Equal(t, "modified private network (test)", postModifyNetwork.Name)

		serverDetails, err := createServer(svc, "TestCreateModifyDeleteNetwork")
		require.NoError(t, err)

		err = stopServer(svc, serverDetails.UUID)
		require.NoError(t, err)

		iface, err := svcContext.CreateNetworkInterface(ctx, &request.CreateNetworkInterfaceRequest{
			ServerUUID:  serverDetails.UUID,
			NetworkUUID: postModifyNetwork.UUID,
			Type:        postModifyNetwork.Type,
			IPAddresses: []request.CreateNetworkInterfaceIPAddress{
				{
					Family: upcloud.IPAddressFamilyIPv4,
				},
			},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, iface.IPAddresses)
		assert.NotEmpty(t, iface.IPAddresses[0].Address)

		modifyIface, err := svcContext.ModifyNetworkInterface(ctx, &request.ModifyNetworkInterfaceRequest{
			ServerUUID:   serverDetails.UUID,
			CurrentIndex: iface.Index,
			NewIndex:     iface.Index + 1,
		})
		require.NoError(t, err)
		assert.Equal(t, iface.Index+1, modifyIface.Index)

		err = svcContext.DeleteNetworkInterface(ctx, &request.DeleteNetworkInterfaceRequest{
			ServerUUID: serverDetails.UUID,
			Index:      modifyIface.Index,
		})
		require.NoError(t, err)

		err = svcContext.DeleteNetwork(ctx, &request.DeleteNetworkRequest{
			UUID: network.UUID,
		})
		require.NoError(t, err)

		networks, err := svcContext.GetNetworksInZone(ctx, &request.GetNetworksInZoneRequest{
			Zone: network.Zone,
		})
		require.NoError(t, err)

		var found bool
		for _, n := range networks.Networks {
			if n.UUID == network.UUID {
				found = true
			}
		}
		assert.False(t, found)
	})
}

// TestGetServerNetworksContext tests that the server networks retrieved via GetServerNetworks
// match those returned when creating the server.
func TestGetServerNetworksContext(t *testing.T) {
	recordWithContext(t, "getservernetworks", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		serverDetails, err := createServer(svc, "TestGetServerNetworks")
		require.NoError(t, err)

		networking, err := svcContext.GetServerNetworks(ctx, &request.GetServerNetworksRequest{
			ServerUUID: serverDetails.UUID,
		})
		require.NoError(t, err)

		sdNetworking := upcloud.Networking(serverDetails.Networking)
		assert.Equal(t, &sdNetworking, networking)
	})
}

// TestGetRoutersContext tests that some routers are returned when using GetRouters.
func TestGetRoutersContext(t *testing.T) {
	recordWithContext(t, "getrouters", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		routers, err := svcContext.GetRouters(ctx)
		require.NoError(t, err)

		assert.Greater(t, len(routers.Routers), 0)
		assert.Greater(t, len(routers.Routers[0].AttachedNetworks), 0)

		router, err := svcContext.GetRouterDetails(ctx, &request.GetRouterDetailsRequest{
			UUID: routers.Routers[0].UUID,
		})
		require.NoError(t, err)

		assert.Equal(t, &(routers.Routers[0]), router)
	})
}

// TestCreateModifyDeleteRouterContext tests router functionality:
// It:
//     - creates a router
//     - modifies a router
//     - retrieves all routers and ensures our new router is found
//     - deletes the router
//     - retrieves all routers and ensures our new router can't be found
func TestCreateModifyDeleteRouterContext(t *testing.T) {
	recordWithContext(t, "createmodifydeleterouter", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		router, err := svcContext.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "Testy McRouterface (test)",
		})
		require.NoError(t, err)

		assert.NotEmpty(t, router.UUID)
		assert.NotEmpty(t, router.Type)
		assert.Equal(t, "Testy McRouterface (test)", router.Name)

		modifiedRouter, err := svcContext.ModifyRouter(ctx, &request.ModifyRouterRequest{
			UUID: router.UUID,
			Name: "Modified name (test)",
		})
		require.NoError(t, err)

		assert.Equal(t, router.UUID, modifiedRouter.UUID)
		assert.Equal(t, router.Type, modifiedRouter.Type)
		assert.NotEqual(t, router.Name, modifiedRouter.Name)
		assert.Equal(t, "Modified name (test)", modifiedRouter.Name)

		routers, err := svcContext.GetRouters(ctx)
		require.NoError(t, err)

		var found bool
		for _, r := range routers.Routers {
			if r.UUID == modifiedRouter.UUID {
				found = true
				break
			}
		}
		assert.True(t, found)

		err = svcContext.DeleteRouter(ctx, &request.DeleteRouterRequest{
			UUID: modifiedRouter.UUID,
		})
		require.NoError(t, err)

		postDeleteRouters, err := svcContext.GetRouters(ctx)
		require.NoError(t, err)

		found = false
		for _, r := range postDeleteRouters.Routers {
			if r.UUID == modifiedRouter.UUID {
				found = true
				break
			}
		}
		assert.False(t, found)
	})
}

// TestCreateTwoNetwoksTwoServersAndARouterContext tests network, server and router functionality
// together.
// It:
//     - creates 2 new networks
//     - creates a router
//     - modifies the two networks to add the router
//     - creates 2 new servers
//     - adds network interfaces in each one server in each network
//     - verifies the network details in the interfaces is correct
//     - verifies the servers can be found in the network details
//     - detaches one of the routers and verifies it was detached
//     - deletes the servers, the routers and the networks
func TestCreateTwoNetworksTwoServersAndARouterContext(t *testing.T) {
	recordWithContext(t, "createtwonetworkstwoserversandarouter", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		network1, err := svcContext.CreateNetwork(ctx, &request.CreateNetworkRequest{
			Name: "test private network #1 (test)",
			Zone: "fi-hel2",
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "192.168.86.1/24",
					DHCP:             upcloud.True,
					DHCPDefaultRoute: upcloud.False,
					DHCPDns: []string{
						"192.168.86.10",
						"192.168.86.11",
					},
					Family:  upcloud.IPAddressFamilyIPv4,
					Gateway: "192.168.86.1",
				},
			},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, network1.UUID)

		network2, err := svcContext.CreateNetwork(ctx, &request.CreateNetworkRequest{
			Name: "test private network #2 (test)",
			Zone: "fi-hel2",
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "192.168.87.1/24",
					DHCP:             upcloud.True,
					DHCPDefaultRoute: upcloud.False,
					DHCPDns: []string{
						"192.168.87.10",
						"192.168.87.11",
					},
					Family:  upcloud.IPAddressFamilyIPv4,
					Gateway: "192.168.87.1",
				},
			},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, network2.UUID)

		router, err := svcContext.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "test router (test)",
		})
		assert.NoError(t, err)

		err = svcContext.AttachNetworkRouter(ctx, &request.AttachNetworkRouterRequest{
			RouterUUID:  router.UUID,
			NetworkUUID: network1.UUID,
		})
		require.NoError(t, err)
		network1, err = svcContext.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: network1.UUID})
		require.NoError(t, err)
		require.Equal(t, network1.Router, router.UUID)

		err = svcContext.AttachNetworkRouter(ctx, &request.AttachNetworkRouterRequest{
			RouterUUID:  router.UUID,
			NetworkUUID: network2.UUID,
		})
		require.NoError(t, err)
		network2, err = svcContext.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: network2.UUID})
		require.NoError(t, err)
		require.Equal(t, network2.Router, router.UUID)

		serverDetails1, err := createServer(svc, "TestCTNTR1")
		require.NoError(t, err)

		serverDetails2, err := createServer(svc, "TestCTNTR2")
		require.NoError(t, err)

		err = stopServer(svc, serverDetails1.UUID)
		require.NoError(t, err)

		err = stopServer(svc, serverDetails2.UUID)
		require.NoError(t, err)

		iface1, err := svcContext.CreateNetworkInterface(ctx, &request.CreateNetworkInterfaceRequest{
			ServerUUID:  serverDetails1.UUID,
			NetworkUUID: network1.UUID,
			Type:        network1.Type,
			IPAddresses: []request.CreateNetworkInterfaceIPAddress{
				{
					Family: upcloud.IPAddressFamilyIPv4,
				},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, network1.UUID, iface1.Network)

		iface2, err := svcContext.CreateNetworkInterface(ctx, &request.CreateNetworkInterfaceRequest{
			ServerUUID:  serverDetails2.UUID,
			NetworkUUID: network2.UUID,
			Type:        network2.Type,
			IPAddresses: []request.CreateNetworkInterfaceIPAddress{
				{
					Family: upcloud.IPAddressFamilyIPv4,
				},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, network2.UUID, iface2.Network)

		serverDetails1PostIface, err := svc.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: serverDetails1.UUID,
		})
		require.NoError(t, err)
		var found bool
		for _, iface := range serverDetails1PostIface.Networking.Interfaces {
			if iface.Network == network1.UUID {
				found = true
			}
		}
		assert.True(t, found)

		serverDetails2PostIface, err := svc.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: serverDetails2.UUID,
		})
		require.NoError(t, err)
		found = false
		for _, iface := range serverDetails2PostIface.Networking.Interfaces {
			if iface.Network == network2.UUID {
				found = true
			}
		}
		assert.True(t, found)

		network1Details, err := svcContext.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
			UUID: network1.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, router.UUID, network1Details.Router)
		found = false
		for _, server := range network1Details.Servers {
			if server.ServerUUID == serverDetails1.UUID {
				found = true
			}
		}
		assert.True(t, found)

		network2Details, err := svcContext.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
			UUID: network2.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, router.UUID, network2Details.Router)
		found = false
		for _, server := range network2Details.Servers {
			if server.ServerUUID == serverDetails2.UUID {
				found = true
			}
		}
		assert.True(t, found)

		// try detaching a router
		err = svcContext.DetachNetworkRouter(ctx, &request.DetachNetworkRouterRequest{NetworkUUID: network1.UUID})
		require.NoError(t, err)
		assert.Eventually(t, func() bool {
			details, err := svcContext.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
				UUID: network1.UUID,
			})
			require.NoError(t, err)
			return err == nil && details.Router == ""
		}, 15*time.Second, time.Second)

		err = deleteServer(svc, serverDetails1.UUID)
		require.NoError(t, err)

		err = deleteServer(svc, serverDetails2.UUID)
		require.NoError(t, err)

		err = svcContext.DeleteNetwork(ctx, &request.DeleteNetworkRequest{
			UUID: network1.UUID,
		})
		require.NoError(t, err)

		err = svcContext.DeleteNetwork(ctx, &request.DeleteNetworkRequest{
			UUID: network2.UUID,
		})
		require.NoError(t, err)

		err = svcContext.DeleteRouter(ctx, &request.DeleteRouterRequest{
			UUID: router.UUID,
		})
		require.NoError(t, err)
	})
}

func TestCreateNetworkAndServerContext(t *testing.T) {
	recordWithContext(t, "createnetworkandserver", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		network, err := svcContext.CreateNetwork(ctx, &request.CreateNetworkRequest{
			Name: "test_network_tcns (test)",
			Zone: "fi-hel2",
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "172.16.0.0/22",
					DHCP:             upcloud.True,
					DHCPDefaultRoute: upcloud.False,
					DHCPDns: []string{
						"172.16.0.10",
						"172.16.1.10",
					},
					Family:  upcloud.IPAddressFamilyIPv4,
					Gateway: "172.16.0.1",
				},
			},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, network.UUID)

		serverDetails, err := createServerWithNetwork(svc, "TestCreateNetworkAndServer", network.UUID)
		require.NoError(t, err)
		assert.NotEmpty(t, serverDetails.UUID)
		var found bool
		for _, iface := range serverDetails.Networking.Interfaces {
			if iface.Network == network.UUID {
				found = true
				break
			}
		}
		assert.True(t, found)
	})
}
