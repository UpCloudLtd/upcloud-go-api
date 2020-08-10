package service

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetNetworks checks that network details are retrievable
// It:
//    - creates a server
//    - Gets all networks and verifies details are populated
//    - checks that at least one network has a server in.
func TestGetNetworks(t *testing.T) {
	record(t, "getnetworks", func(t *testing.T, svc *Service) {
		_, err := createServer(svc, "TestGetNetworks")
		require.NoError(t, err)

		networks, err := svc.GetNetworks()
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

// TestGetNetworksInZone checks that network details in a zone are retrievable
// It:
//    - creates a server
//    - Gets all networks in a zone and verifies details are populated
//    - checks that at least one network has a server in.
func TestGetNetworksInZone(t *testing.T) {
	record(t, "getnetworksinzone", func(t *testing.T, svc *Service) {
		_, err := createServer(svc, "TestGetNetworksInZone")
		require.NoError(t, err)

		networks, err := svc.GetNetworksInZone(&request.GetNetworksInZoneRequest{
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
		for _, n := range networks.Networks {
			if len(n.Servers) > 0 {
				foundNetwork = &n
				found = true
			}
			// Make sure all the networks are in the right zone.
			assert.Equal(t, "fi-hel2", n.Zone)
		}
		assert.True(t, found)
		require.NotNil(t, foundNetwork)

		network, err := svc.GetNetworkDetails(&request.GetNetworkDetailsRequest{
			UUID: foundNetwork.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, foundNetwork, network)
	})
}

// TestCreateModifyDeleteNetwork checks that the network functionality works.
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
func TestCreateModifyDeleteNetwork(t *testing.T) {
	record(t, "createmodifydeletenetwork", func(t *testing.T, svc *Service) {
		network, err := svc.CreateNetwork(&request.CreateNetworkRequest{
			Name: "test private network (test)",
			Zone: "fi-hel2",
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "172.16.0.0/22",
					DHCP:             true,
					DHCPDefaultRoute: false,
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

		postModifyNetwork, err := svc.ModifyNetwork(&request.ModifyNetworkRequest{
			UUID: network.UUID,
			Name: "modified private network (test)",
		})
		require.NoError(t, err)
		assert.Equal(t, "modified private network (test)", postModifyNetwork.Name)

		serverDetails, err := createServer(svc, "TestCreateModifyDeleteNetwork")
		require.NoError(t, err)

		err = stopServer(svc, serverDetails.UUID)
		require.NoError(t, err)

		iface, err := svc.CreateNetworkInterface(&request.CreateNetworkInterfaceRequest{
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

		modifyIface, err := svc.ModifyNetworkInterface(&request.ModifyNetworkInterfaceRequest{
			ServerUUID:   serverDetails.UUID,
			CurrentIndex: iface.Index,
			NewIndex:     iface.Index + 1,
		})
		require.NoError(t, err)
		assert.Equal(t, iface.Index+1, modifyIface.Index)

		err = svc.DeleteNetworkInterface(&request.DeleteNetworkInterfaceRequest{
			ServerUUID: serverDetails.UUID,
			Index:      modifyIface.Index,
		})
		require.NoError(t, err)

		err = svc.DeleteNetwork(&request.DeleteNetworkRequest{
			UUID: network.UUID,
		})
		require.NoError(t, err)

		networks, err := svc.GetNetworksInZone(&request.GetNetworksInZoneRequest{
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

// TestGetServerNetworks tests that the server networks retrieved via GetServerNetworks
// match those returned when creating the server.
func TestGetServerNetworks(t *testing.T) {
	record(t, "getservernetworks", func(t *testing.T, svc *Service) {
		serverDetails, err := createServer(svc, "TestGetServerNetworks")
		require.NoError(t, err)

		networking, err := svc.GetServerNetworks(&request.GetServerNetworksRequest{
			ServerUUID: serverDetails.UUID,
		})
		require.NoError(t, err)

		sdNetworking := upcloud.Networking(serverDetails.Networking)
		assert.Equal(t, &sdNetworking, networking)
	})
}

// TestGetRouters tests that some routers are returned when using GetRouters.
func TestGetRouters(t *testing.T) {
	record(t, "getrouters", func(t *testing.T, svc *Service) {
		routers, err := svc.GetRouters()
		require.NoError(t, err)

		assert.Greater(t, len(routers.Routers), 0)
		assert.Greater(t, len(routers.Routers[0].AttachedNetworks), 0)

		router, err := svc.GetRouterDetails(&request.GetRouterDetailsRequest{
			UUID: routers.Routers[0].UUID,
		})
		require.NoError(t, err)

		assert.Equal(t, &(routers.Routers[0]), router)
	})
}

// TestCreateModifyDeleteRouter tests router functionality:
// It:
//     - creates a router
//     - modifies a router
//     - retrieves all routers and ensures our new router is found
//     - deletes the router
//     - retrieves all routers and ensures our new router can't be found
func TestCreateModifyDeleteRouter(t *testing.T) {
	record(t, "createmodifydeleterouter", func(t *testing.T, svc *Service) {
		router, err := svc.CreateRouter(&request.CreateRouterRequest{
			Name: "Testy McRouterface (test)",
		})
		require.NoError(t, err)

		assert.NotEmpty(t, router.UUID)
		assert.NotEmpty(t, router.Type)
		assert.Equal(t, "Testy McRouterface (test)", router.Name)

		modifiedRouter, err := svc.ModifyRouter(&request.ModifyRouterRequest{
			UUID: router.UUID,
			Name: "Modified name (test)",
		})
		require.NoError(t, err)

		assert.Equal(t, router.UUID, modifiedRouter.UUID)
		assert.Equal(t, router.Type, modifiedRouter.Type)
		assert.NotEqual(t, router.Name, modifiedRouter.Name)
		assert.Equal(t, "Modified name (test)", modifiedRouter.Name)

		routers, err := svc.GetRouters()
		require.NoError(t, err)

		var found bool
		for _, r := range routers.Routers {
			if r.UUID == modifiedRouter.UUID {
				found = true
				break
			}
		}
		assert.True(t, found)

		err = svc.DeleteRouter(&request.DeleteRouterRequest{
			UUID: modifiedRouter.UUID,
		})
		require.NoError(t, err)

		postDeleteRouters, err := svc.GetRouters()
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

// TestCreateTwoNetwoksTwoServersAndARouter tests network, server and router functionality
// together.
// It:
//     - creates 2 new networks
//     - creates a router
//     - modifies the two networks to add the router
//     - creates 2 new servers
//     - adds network interfaces in each one server in each network
//     - verifies the network details in the interfaces is correct
//     - verifies the servers can be found in the network details
//     - deletes the servers, the routers and the networks
func TestCreateTwoNetworksTwoServersAndARouter(t *testing.T) {
	record(t, "createtwonetworkstwoserversandarouter", func(t *testing.T, svc *Service) {
		network1, err := svc.CreateNetwork(&request.CreateNetworkRequest{
			Name: "test private network #1 (test)",
			Zone: "fi-hel2",
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "192.168.86.1/24",
					DHCP:             true,
					DHCPDefaultRoute: false,
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

		network2, err := svc.CreateNetwork(&request.CreateNetworkRequest{
			Name: "test private network #2 (test)",
			Zone: "fi-hel2",
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "192.168.87.1/24",
					DHCP:             true,
					DHCPDefaultRoute: false,
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

		router, err := svc.CreateRouter(&request.CreateRouterRequest{
			Name: "test router (test)",
		})
		assert.NoError(t, err)

		network1, err = svc.ModifyNetwork(&request.ModifyNetworkRequest{
			Router: router.UUID,
			UUID:   network1.UUID,
		})
		require.NoError(t, err)
		require.Equal(t, network1.Router, router.UUID)

		network2, err = svc.ModifyNetwork(&request.ModifyNetworkRequest{
			Router: router.UUID,
			UUID:   network2.UUID,
		})
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

		iface1, err := svc.CreateNetworkInterface(&request.CreateNetworkInterfaceRequest{
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

		iface2, err := svc.CreateNetworkInterface(&request.CreateNetworkInterfaceRequest{
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

		network1Details, err := svc.GetNetworkDetails(&request.GetNetworkDetailsRequest{
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

		network2Details, err := svc.GetNetworkDetails(&request.GetNetworkDetailsRequest{
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

		err = deleteServer(svc, serverDetails1.UUID)
		require.NoError(t, err)

		err = deleteServer(svc, serverDetails2.UUID)
		require.NoError(t, err)

		err = svc.DeleteNetwork(&request.DeleteNetworkRequest{
			UUID: network1.UUID,
		})
		require.NoError(t, err)

		err = svc.DeleteNetwork(&request.DeleteNetworkRequest{
			UUID: network2.UUID,
		})
		require.NoError(t, err)

		err = svc.DeleteRouter(&request.DeleteRouterRequest{
			UUID: router.UUID,
		})
		require.NoError(t, err)

	})
}
