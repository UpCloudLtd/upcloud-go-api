package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// TestGetNetworks checks that network details are retrievable
// It:
//   - creates a server
//   - Gets all networks and verifies details are populated
//   - checks that at least one network has a server in.
func TestGetNetworks(t *testing.T) {
	record(t, "getnetworks", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		_, err := createServer(ctx, rec, svc, "TestGetNetworks")
		require.NoError(t, err)

		networks, err := svc.GetNetworks(ctx)
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

// TestGetNetworksWithFilters checks if you can get a networks while filtering by labels
func TestGetNetworksWithFilters(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/network?label=env%3Dtest&label=managedBy%3Dupcloud"), r.URL.String())
		fmt.Fprint(w, `
		{
			"networks": {
				"network": [
					{
						"ip_networks": {
							"ip_network": [
								{
									"address": "123.54.123.0/22",
									"dhcp": "yes",
									"dhcp_default_route": "yes",
									"dhcp_dns": [
										"94.123.111.9",
										"94.123.112.9"
									],
									"family": "IPv4",
									"gateway": "123.54.123.1"
								}
							]
						},
						"labels": [
							{
								"key": "env",
								"value": "test"
							},
							{
								"key": "managedBy",
								"value": "upcloud"
							}
						],
						"name": "net1",
						"type": "public",
						"uuid": "uuid1",
						"zone": "fi-hel1"
					},
					{
						"ip_networks": {
							"ip_network": [
								{
									"address": "185.123.136.0/22",
									"dhcp": "yes",
									"dhcp_default_route": "yes",
									"dhcp_dns": [
										"94.123.127.9",
										"94.123.40.9"
									],
									"family": "IPv4",
									"gateway": "185.123.136.1"
								}
							]
						},
						"labels": [
							{
								"key": "env",
								"value": "test"
							},
							{
								"key": "managedBy",
								"value": "upcloud"
							}
						],
						"name": "net2",
						"type": "public",
						"uuid": "uuid2",
						"zone": "fi-hel1"
					}
				]
			}
		}
		`)
	}))

	defer srv.Close()

	filters := []request.QueryFilter{
		request.FilterLabel{Label: upcloud.Label{
			Key:   "env",
			Value: "test",
		}},
		request.FilterLabel{Label: upcloud.Label{
			Key:   "managedBy",
			Value: "upcloud",
		}},
	}

	res, err := svc.GetNetworks(context.Background(), filters...)
	assert.NoError(t, err)
	assert.Len(t, res.Networks, 2)
	assert.Equal(t, "uuid1", res.Networks[0].UUID)
	assert.Equal(t, "uuid2", res.Networks[1].UUID)
	assert.Equal(t, "env", res.Networks[0].Labels[0].Key)
	assert.Equal(t, "test", res.Networks[0].Labels[0].Value)
	assert.Equal(t, "managedBy", res.Networks[1].Labels[1].Key)
	assert.Equal(t, "upcloud", res.Networks[1].Labels[1].Value)
}

func TestGetNetworkDetails(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/network/_UUID_"), r.URL.Path)
		fmt.Fprint(w, `
			{
				"network": {
					"ip_networks": {
						"ip_network": [
							{
								"address": "172.16.2.0/24",
								"dhcp": "yes",
								"dhcp_default_route": "no",
								"family": "IPv4",
								"gateway": "172.16.2.1"
							}
						]
					},
					"labels": [
						{
							"key": "env",
							"value": "test"
						}
					],
					"name": "testnetwork",
					"type": "private",
					"uuid": "_UUID_",
					"zone": "de-fra1"
				}
			}
		`)
	}))

	defer srv.Close()

	network, err := svc.GetNetworkDetails(context.Background(), &request.GetNetworkDetailsRequest{UUID: "_UUID_"})
	require.NoError(t, err)
	assert.Equal(t, "testnetwork", network.Name)
	assert.Len(t, network.Labels, 1)
	assert.Equal(t, "env", network.Labels[0].Key)
	assert.Equal(t, "test", network.Labels[0].Value)
	assert.Equal(t, "de-fra1", network.Zone)
}

// TestGetNetworksInZone checks that network details in a zone are retrievable
// It:
//   - creates a server
//   - Gets all networks in a zone and verifies details are populated
//   - checks that at least one network has a server in.
func TestGetNetworksInZone(t *testing.T) {
	record(t, "getnetworksinzone", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		_, err := createServer(ctx, rec, svc, "TestGetNetworksInZone")
		require.NoError(t, err)

		networks, err := svc.GetNetworksInZone(ctx, &request.GetNetworksInZoneRequest{
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

		network, err := svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
			UUID: foundNetwork.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, foundNetwork, network)
	})
}

// TestCreateModifyDeleteNetwork checks that the network functionality works.
// It:
//   - creates a network
//   - modifies the network
//   - creates a server
//   - stops the server
//   - creates a network interface in the network
//   - modifies the network interface
//   - deletes the network interface
//   - deletes the network
//   - verifies the network has been deleted.
func TestCreateModifyDeleteNetwork(t *testing.T) {
	record(t, "createmodifydeletenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		network, err := svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
			Name: "test private network (test)",
			Zone: "fi-hel2",
			Labels: []upcloud.Label{
				{
					Key:   "env",
					Value: "test",
				},
			},
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "172.17.0.0/22",
					DHCP:             upcloud.True,
					DHCPDefaultRoute: upcloud.False,
					DHCPDns: []string{
						"172.17.0.10",
						"172.17.1.10",
					},
					DHCPRoutes: []string{
						"192.168.0.0/24",
						"192.168.100.100/32",
					},
					Family:  upcloud.IPAddressFamilyIPv4,
					Gateway: "172.17.0.1",
				},
			},
		})
		require.NoError(t, err)
		assert.NotEmpty(t, network.UUID)
		assert.Equal(t, "test private network (test)", network.Name)
		assert.Len(t, network.Labels, 1)
		assert.Equal(t, "env", network.Labels[0].Key)
		assert.Equal(t, "test", network.Labels[0].Value)
		assert.Equal(t, []string{"192.168.0.0/24", "192.168.100.100/32"}, network.IPNetworks[0].DHCPRoutes)

		postModifyNetwork, err := svc.ModifyNetwork(ctx, &request.ModifyNetworkRequest{
			UUID: network.UUID,
			Name: "modified private network (test)",
		})
		require.NoError(t, err)
		assert.Equal(t, "modified private network (test)", postModifyNetwork.Name)
		assert.Equal(t, network.IPNetworks, postModifyNetwork.IPNetworks)
		assert.Len(t, postModifyNetwork.Labels, 1) // Make sure labels are not deleted on simple update

		postModifyNetworkWithLabels, err := svc.ModifyNetwork(ctx, &request.ModifyNetworkRequest{
			UUID: network.UUID,
			Labels: &[]upcloud.Label{
				{
					Key:   "env",
					Value: "test",
				},
				{
					Key:   "managedBy",
					Value: "upcloud",
				},
			},
		})
		require.NoError(t, err)
		assert.Len(t, postModifyNetworkWithLabels.Labels, 2)
		assert.Equal(t, "env", postModifyNetworkWithLabels.Labels[0].Key)
		assert.Equal(t, "test", postModifyNetworkWithLabels.Labels[0].Value)
		assert.Equal(t, "managedBy", postModifyNetworkWithLabels.Labels[1].Key)
		assert.Equal(t, "upcloud", postModifyNetworkWithLabels.Labels[1].Value)

		serverDetails, err := createServer(ctx, rec, svc, "TestCreateModifyDeleteNetwork")
		require.NoError(t, err)

		err = stopServer(ctx, rec, svc, serverDetails.UUID)
		require.NoError(t, err)

		iface, err := svc.CreateNetworkInterface(ctx, &request.CreateNetworkInterfaceRequest{
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

		modifyIface, err := svc.ModifyNetworkInterface(ctx, &request.ModifyNetworkInterfaceRequest{
			ServerUUID:   serverDetails.UUID,
			CurrentIndex: iface.Index,
			NewIndex:     iface.Index + 1,
		})
		require.NoError(t, err)
		assert.Equal(t, iface.Index+1, modifyIface.Index)

		err = svc.DeleteNetworkInterface(ctx, &request.DeleteNetworkInterfaceRequest{
			ServerUUID: serverDetails.UUID,
			Index:      modifyIface.Index,
		})
		require.NoError(t, err)

		err = svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{
			UUID: network.UUID,
		})
		require.NoError(t, err)

		networks, err := svc.GetNetworksInZone(ctx, &request.GetNetworksInZoneRequest{
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
	record(t, "getservernetworks", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		serverDetails, err := createServer(ctx, rec, svc, "TestGetServerNetworks")
		require.NoError(t, err)

		networking, err := svc.GetServerNetworks(ctx, &request.GetServerNetworksRequest{
			ServerUUID: serverDetails.UUID,
		})
		require.NoError(t, err)

		sdNetworking := upcloud.Networking(serverDetails.Networking)
		assert.Equal(t, &sdNetworking, networking)
	})
}

// TestGetRouters tests that some routers are returned when using GetRouters.
func TestGetRouters(t *testing.T) {
	record(t, "getrouters", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		routers, err := svc.GetRouters(ctx)
		require.NoError(t, err)

		assert.Greater(t, len(routers.Routers), 0)
		assert.Greater(t, len(routers.Routers[0].AttachedNetworks), 0)

		router, err := svc.GetRouterDetails(ctx, &request.GetRouterDetailsRequest{
			UUID: routers.Routers[0].UUID,
		})
		require.NoError(t, err)

		assert.ElementsMatch(t, routers.Routers[0].AttachedNetworks, router.AttachedNetworks)
		assert.Equal(t, routers.Routers[0].Name, router.Name)
		assert.Equal(t, routers.Routers[0].Type, router.Type)
		assert.Equal(t, routers.Routers[0].UUID, router.UUID)
	})
}

// TestCreateModifyDeleteRouterContext tests router functionality:
// It:
//   - creates a router
//   - modifies a router
//   - retrieves all routers and ensures our new router is found
//   - deletes the router
//   - retrieves all routers and ensures our new router can't be found
func TestCreateModifyDeleteRouter(t *testing.T) {
	record(t, "createmodifydeleterouter", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "Testy McRouterface (test)",
			StaticRoutes: []upcloud.StaticRoute{
				{
					Name:    "static_route_0",
					Route:   "0.0.0.0/0",
					Nexthop: "10.0.0.100",
				},
			},
		})
		require.NoError(t, err)

		assert.NotEmpty(t, router.UUID)
		assert.NotEmpty(t, router.Type)
		assert.Equal(t, "Testy McRouterface (test)", router.Name)
		assert.Equal(t,
			[]upcloud.StaticRoute{{
				Name:    "static_route_0",
				Route:   "0.0.0.0/0",
				Nexthop: "10.0.0.100",
			}},
			router.StaticRoutes,
		)

		modifiedRouter, err := svc.ModifyRouter(ctx, &request.ModifyRouterRequest{
			UUID: router.UUID,
			Name: "Modified name (test)",
		})
		require.NoError(t, err)

		assert.Equal(t, router.UUID, modifiedRouter.UUID)
		assert.Equal(t, router.Type, modifiedRouter.Type)
		assert.Equal(t, router.StaticRoutes, modifiedRouter.StaticRoutes)
		assert.NotEqual(t, router.Name, modifiedRouter.Name)
		assert.Equal(t, "Modified name (test)", modifiedRouter.Name)

		routers, err := svc.GetRouters(ctx)
		require.NoError(t, err)

		var found bool
		for _, r := range routers.Routers {
			if r.UUID == modifiedRouter.UUID {
				found = true
				break
			}
		}
		assert.True(t, found)

		err = svc.DeleteRouter(ctx, &request.DeleteRouterRequest{
			UUID: modifiedRouter.UUID,
		})
		require.NoError(t, err)

		postDeleteRouters, err := svc.GetRouters(ctx)
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
//   - creates 2 new networks
//   - creates a router
//   - modifies the two networks to add the router
//   - creates 2 new servers
//   - adds network interfaces in each one server in each network
//   - verifies the network details in the interfaces is correct
//   - verifies the servers can be found in the network details
//   - detaches one of the routers and verifies it was detached
//   - deletes the servers, the routers and the networks
func TestCreateTwoNetworksTwoServersAndARouter(t *testing.T) {
	record(t, "createtwonetworkstwoserversandarouter", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		network1, err := svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
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

		network2, err := svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
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

		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "test router (test)",
		})
		assert.NoError(t, err)

		err = svc.AttachNetworkRouter(ctx, &request.AttachNetworkRouterRequest{
			RouterUUID:  router.UUID,
			NetworkUUID: network1.UUID,
		})
		require.NoError(t, err)
		network1, err = svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: network1.UUID})
		require.NoError(t, err)
		require.Equal(t, network1.Router, router.UUID)

		err = svc.AttachNetworkRouter(ctx, &request.AttachNetworkRouterRequest{
			RouterUUID:  router.UUID,
			NetworkUUID: network2.UUID,
		})
		require.NoError(t, err)
		network2, err = svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: network2.UUID})
		require.NoError(t, err)
		require.Equal(t, network2.Router, router.UUID)

		serverDetails1, err := createServer(ctx, rec, svc, "TestCTNTR1")
		require.NoError(t, err)

		serverDetails2, err := createServer(ctx, rec, svc, "TestCTNTR2")
		require.NoError(t, err)

		err = stopServer(ctx, rec, svc, serverDetails1.UUID)
		require.NoError(t, err)

		err = stopServer(ctx, rec, svc, serverDetails2.UUID)
		require.NoError(t, err)

		iface1, err := svc.CreateNetworkInterface(ctx, &request.CreateNetworkInterfaceRequest{
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

		iface2, err := svc.CreateNetworkInterface(ctx, &request.CreateNetworkInterfaceRequest{
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

		serverDetails1PostIface, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
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

		serverDetails2PostIface, err := svc.GetServerDetails(ctx, &request.GetServerDetailsRequest{
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

		network1Details, err := svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
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

		network2Details, err := svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
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
		err = svc.DetachNetworkRouter(ctx, &request.DetachNetworkRouterRequest{NetworkUUID: network1.UUID})
		require.NoError(t, err)

		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})

			assert.Eventually(t, func() bool {
				details, err := svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
					UUID: network1.UUID,
				})
				require.NoError(t, err)
				return err == nil && details.Router == ""
			}, 15*time.Second, time.Second)

			rec.Passthroughs = nil
		}
		details, err := svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{
			UUID: network1.UUID,
		})
		require.NoError(t, err)
		assert.Empty(t, details.Router)

		err = deleteServer(ctx, svc, serverDetails1.UUID)
		require.NoError(t, err)

		err = deleteServer(ctx, svc, serverDetails2.UUID)
		require.NoError(t, err)

		err = svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{
			UUID: network1.UUID,
		})
		require.NoError(t, err)

		err = svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{
			UUID: network2.UUID,
		})
		require.NoError(t, err)

		err = svc.DeleteRouter(ctx, &request.DeleteRouterRequest{
			UUID: router.UUID,
		})
		require.NoError(t, err)
	})
}

func TestCreateNetworkAndServer(t *testing.T) {
	record(t, "createnetworkandserver", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		network, err := svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
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

		serverDetails, err := createServerWithNetwork(ctx, rec, svc, "TestCreateNetworkAndServer", network.UUID)
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

func TestCreateNetworkAndServer_DHCPOptions(t *testing.T) {
	record(t, "dhcpnetworkconfigurations", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create a DHCP-enabled network WITH default route via DHCP
		netReq := &request.CreateNetworkRequest{
			Name: "sdk-test default-route yes",
			Zone: "fi-hel2",
			IPNetworks: []upcloud.IPNetwork{
				{
					Address:          "172.16.0.0/22",
					DHCP:             upcloud.True,
					DHCPDefaultRoute: upcloud.True,
					DHCPDns:          []string{"172.16.0.10", "172.16.1.10"},
					Family:           upcloud.IPAddressFamilyIPv4,
					Gateway:          "172.16.0.1",
				},
			},
		}
		network, err := svc.CreateNetwork(ctx, netReq)
		require.NoError(t, err)
		require.NotEmpty(t, network.UUID)

		// Create a server attached to this network
		srv, err := createServerWithNetwork(ctx, rec, svc, "TestCreateNetworkAndServerDHCP", network.UUID)
		require.NoError(t, err)
		require.NotEmpty(t, srv.UUID)

		// Re-fetch network details and assert DHCP default route = yes
		netDetails, err := svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: network.UUID})
		require.NoError(t, err)
		require.NotEmpty(t, netDetails.IPNetworks)
		ipNet := netDetails.IPNetworks[0]
		assert.Equal(t, upcloud.True, ipNet.DHCPDefaultRoute, "dhcp_default_route should be 'yes'")

		// Server should be attached to our network
		found := false
		for _, nic := range srv.Networking.Interfaces {
			if nic.Network == network.UUID {
				found = true
				break
			}
		}
		assert.True(t, found, "server should be attached to created network")

		t.Run("DHCP_AutoPopulation_AllFilters", func(t *testing.T) {
			modReq := &request.ModifyNetworkRequest{
				UUID: network.UUID,
				IPNetworks: []upcloud.IPNetwork{
					{
						Address:          "172.16.0.0/22",
						DHCP:             upcloud.True,
						DHCPDefaultRoute: upcloud.True,
						DHCPDns:          []string{"172.16.0.10", "172.16.1.10"},
						Family:           upcloud.IPAddressFamilyIPv4,
						Gateway:          "172.16.0.1",
						DHCPRoutesConfiguration: upcloud.DHCPRoutesConfiguration{
							EffectiveRoutesAutoPopulation: upcloud.EffectiveRoutesAutoPopulation{
								Enabled:             upcloud.True,
								FilterByDestination: &[]string{"172.16.0.0/22"},
								FilterByRouteType:   &[]upcloud.NetworkRouteType{"service"},
								ExcludeBySource:     &[]upcloud.NetworkRouteSource{"static-route"},
							},
						},
					},
				},
			}

			_, err := svc.ModifyNetwork(ctx, modReq)
			require.NoError(t, err, "ModifyNetwork with all filters should succeed")

			// Re-read and assert the config was applied
			netDetails, err := svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: network.UUID})
			require.NoError(t, err)
			require.NotEmpty(t, netDetails.IPNetworks)
			ipNet := netDetails.IPNetworks[0]

			require.Equal(t, upcloud.True,
				ipNet.DHCPRoutesConfiguration.EffectiveRoutesAutoPopulation.Enabled,
				"auto-population should be enabled")

			assert.ElementsMatch(t,
				[]string{"172.16.0.0/22"},
				*ipNet.DHCPRoutesConfiguration.EffectiveRoutesAutoPopulation.FilterByDestination,
			)

			assert.ElementsMatch(t,
				[]upcloud.NetworkRouteType{"service"},
				*ipNet.DHCPRoutesConfiguration.EffectiveRoutesAutoPopulation.FilterByRouteType,
			)

			assert.ElementsMatch(t,
				[]upcloud.NetworkRouteSource{"static-route"},
				*ipNet.DHCPRoutesConfiguration.EffectiveRoutesAutoPopulation.ExcludeBySource,
			)
		})

		t.Run("DHCP_AutoPopulation_Unset_Filters", func(t *testing.T) {
			modReq := &request.ModifyNetworkRequest{
				UUID: network.UUID,
				IPNetworks: []upcloud.IPNetwork{
					{
						Address:          "172.16.0.0/22",
						DHCP:             upcloud.True,
						DHCPDefaultRoute: upcloud.True,
						DHCPDns:          []string{"172.16.0.10", "172.16.1.10"},
						Family:           upcloud.IPAddressFamilyIPv4,
						Gateway:          "172.16.0.1",
						DHCPRoutesConfiguration: upcloud.DHCPRoutesConfiguration{
							EffectiveRoutesAutoPopulation: upcloud.EffectiveRoutesAutoPopulation{
								Enabled:             upcloud.True,
								FilterByDestination: &[]string{},
								ExcludeBySource:     &[]upcloud.NetworkRouteSource{"static-route"},
							},
						},
					},
				},
			}

			_, err := svc.ModifyNetwork(ctx, modReq)
			require.NoError(t, err, "ModifyNetwork with all filters should succeed")

			// Re-read and assert the config was applied
			netDetails, err := svc.GetNetworkDetails(ctx, &request.GetNetworkDetailsRequest{UUID: network.UUID})
			require.NoError(t, err)
			require.NotEmpty(t, netDetails.IPNetworks)
			ipNet := netDetails.IPNetworks[0]

			require.Equal(t, upcloud.True,
				ipNet.DHCPRoutesConfiguration.EffectiveRoutesAutoPopulation.Enabled,
				"auto-population should be enabled")

			assert.Nil(t, ipNet.DHCPRoutesConfiguration.EffectiveRoutesAutoPopulation.FilterByDestination)

			assert.ElementsMatch(t,
				[]upcloud.NetworkRouteType{"service"},
				*ipNet.DHCPRoutesConfiguration.EffectiveRoutesAutoPopulation.FilterByRouteType,
			)

			assert.ElementsMatch(t,
				[]upcloud.NetworkRouteSource{"static-route"},
				*ipNet.DHCPRoutesConfiguration.EffectiveRoutesAutoPopulation.ExcludeBySource,
			)
		})

		// Stop the server
		t.Logf("Stopping server with UUID %s ...", srv.UUID)
		err = stopServer(ctx, rec, svc, srv.UUID)
		require.NoError(t, err)
		t.Log("Server is now stopped")

		// Delete the server and storage
		t.Logf("Deleting the server with UUID %s, including storages...", srv.UUID)
		err = deleteServerAndStorages(ctx, svc, srv.UUID)
		require.NoError(t, err)
		t.Log("Server is now deleted")

		// Delete the network
		t.Logf("Deleting the network with UUID %s...", network.UUID)
		err = svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{
			UUID: network.UUID,
		})
		require.NoError(t, err)
		t.Log("Network is now deleted")
	})
}

// TestRoutersFilters tests router labels and filters
func TestRouterLabelsAndFilters(t *testing.T) {
	record(t, "routerlabelsandfilters", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		r1, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "test_router_labels_and_filters_1",
			Labels: []upcloud.Label{
				{
					Key:   "color",
					Value: "blue",
				},
			},
		})
		require.NoError(t, err)
		r2, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "test_router_labels_and_filters_1",
			Labels: []upcloud.Label{
				{
					Key:   "color",
					Value: "red",
				},
			},
		})
		require.NoError(t, err)
		routers, err := svc.GetRouters(ctx, request.FilterLabelKey{Key: "color"})
		require.NoError(t, err)
		assert.Equal(t, len(routers.Routers), 2)

		routers, err = svc.GetRouters(ctx, request.FilterLabel{
			Label: upcloud.Label{Key: "color", Value: "red"},
		})
		assert.NoError(t, err)
		assert.Equal(t, len(routers.Routers), 1)

		_, err = svc.ModifyRouter(ctx, &request.ModifyRouterRequest{
			UUID:   r2.UUID,
			Name:   r2.Name,
			Labels: &[]upcloud.Label{},
		})
		assert.NoError(t, err)

		routers, err = svc.GetRouters(ctx, request.FilterLabel{
			Label: upcloud.Label{Key: "color", Value: "red"},
		})
		assert.NoError(t, err)
		assert.Equal(t, len(routers.Routers), 0)

		require.NoError(t, svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: r1.UUID}))
		require.NoError(t, svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: r2.UUID}))
	})
}
