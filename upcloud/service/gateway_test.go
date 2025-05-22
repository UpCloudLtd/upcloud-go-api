package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

func TestGatewayPlans(t *testing.T) {
	t.Parallel()

	plansResponse := `
	[
		{
		  "name": "advanced",
		  "per_gateway_bandwidth_mbps": 10000,
		  "per_gateway_max_connections": 100000,
		  "server_number": 2,
		  "supported_features": [
			  "nat",
			  "vpn"
		  ],
		  "vpn_tunnel_amount": 10
		},
		{
		  "name": "production",
		  "per_gateway_bandwidth_mbps": 1000,
		  "per_gateway_max_connections": 50000,
		  "server_number": 2,
		  "supported_features": [
			  "nat",
			  "vpn"
		  ],
		  "vpn_tunnel_amount": 2
		}
	]
	`

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/gateway/plans", client.APIVersion), r.URL.Path)
		_, _ = fmt.Fprint(w, plansResponse)
	}))
	defer srv.Close()

	res, err := svc.GetGatewayPlans(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 2)

	firstPlan := res[0]
	secondPlan := res[1]

	assert.Equal(t, "advanced", firstPlan.Name)
	assert.Equal(t, 10000, firstPlan.PerGatewayBandwidthMbps)
	assert.Equal(t, 100000, firstPlan.PerGatewayMaxConnections)
	assert.Equal(t, 2, firstPlan.ServerNumber)
	assert.Len(t, firstPlan.SupportedFeatures, 2)
	assert.Equal(t, upcloud.GatewayFeatureNAT, firstPlan.SupportedFeatures[0])
	assert.Equal(t, upcloud.GatewayFeatureVPN, firstPlan.SupportedFeatures[1])
	assert.Equal(t, 10, firstPlan.VPNTunnelAmount)

	// Just check the name, no need to check all the properties again
	assert.Equal(t, "production", secondPlan.Name)
}

func TestNATGateway(t *testing.T) {
	t.Parallel()

	record(t, "gatewaynat", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "test-router",
		})
		if !assert.NoError(t, err) {
			return
		}
		gw, err := svc.CreateGateway(ctx, &request.CreateGatewayRequest{
			Name: "test",
			Zone: "pl-waw1",
			Features: []upcloud.GatewayFeature{
				upcloud.GatewayFeatureNAT,
			},
			Routers: []request.GatewayRouter{
				{
					UUID: router.UUID,
				},
			},
			ConfiguredStatus: upcloud.GatewayConfiguredStatusStarted,
		})
		if !assert.NoError(t, err) {
			return
		}

		if !assert.NoError(t, waitGatewayToStart(ctx, rec, svc, gw.UUID)) {
			return
		}
		gw, err = svc.GetGateway(ctx, &request.GetGatewayRequest{UUID: gw.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "test", gw.Name)
		assert.Equal(t, "pl-waw1", gw.Zone)
		if assert.GreaterOrEqual(t, 1, len(gw.Features)) {
			assert.Equal(t, upcloud.GatewayFeatureNAT, gw.Features[0])
		}
		if assert.Len(t, gw.Routers, 1) {
			assert.Equal(t, router.UUID, gw.Routers[0].UUID)
		}
		assert.Len(t, gw.Addresses, 1)
		gw, err = svc.ModifyGateway(ctx, &request.ModifyGatewayRequest{
			UUID:             gw.UUID,
			Name:             "new-name",
			ConfiguredStatus: upcloud.GatewayConfiguredStatusStopped,
		})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "new-name", gw.Name)
		assert.Equal(t, upcloud.GatewayConfiguredStatusStopped, gw.ConfiguredStatus)

		assert.NoError(t, svc.DeleteGateway(ctx, &request.DeleteGatewayRequest{UUID: gw.UUID}))

		if err := waitGatewayToDelete(ctx, rec, svc, gw.UUID); err != nil {
			t.Fatal(err)
		}
		if err := svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: router.UUID}); err != nil {
			t.Log(err)
		}
	})
}

func TestVPNGateway(t *testing.T) {
	t.Parallel()

	record(t, "gatewayvpn", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{Name: "test-router-vpn"})
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			err = svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: router.UUID})
			assert.NoError(t, err)
		}()

		plans, err := svc.GetGatewayPlans(ctx)
		if !assert.NoError(t, err) {
			return
		}
		if !assert.GreaterOrEqual(t, len(plans), 2, "plans response has less than 2 plans") {
			return
		}

		psk := "key123wouldkeepitthatwaybuthastobelonger"
		gw, err := svc.CreateGateway(ctx, &request.CreateGatewayRequest{
			Name: "test-vpn",
			Zone: "pl-waw1",
			Routers: []request.GatewayRouter{
				{
					UUID: router.UUID,
				},
			},
			Plan:             plans[0].Name,
			ConfiguredStatus: upcloud.GatewayConfiguredStatusStarted,
			Features: []upcloud.GatewayFeature{
				upcloud.GatewayFeatureVPN,
			},
			Addresses: []upcloud.GatewayAddress{
				{
					Name: "my-public-ip",
				},
			},
			Connections: []request.GatewayConnection{
				{
					Name: "example-connection",
					Type: upcloud.GatewayConnectionTypeIPSec,
					LocalRoutes: []upcloud.GatewayRoute{
						{
							Name:          "local-route",
							Type:          upcloud.GatewayRouteTypeStatic,
							StaticNetwork: "10.0.0.0/24",
						},
					},
					RemoteRoutes: []upcloud.GatewayRoute{
						{
							Name:          "remote-route",
							Type:          upcloud.GatewayRouteTypeStatic,
							StaticNetwork: "10.0.1.0/24",
						},
					},
					Tunnels: []request.GatewayTunnel{
						{
							Name: "example-tunnel",
							LocalAddress: upcloud.GatewayTunnelLocalAddress{
								Name: "my-public-ip",
							},
							RemoteAddress: upcloud.GatewayTunnelRemoteAddress{
								Address: "100.10.0.111",
							},
							IPSec: upcloud.GatewayTunnelIPSec{
								Authentication: upcloud.GatewayTunnelIPSecAuth{
									Authentication: upcloud.GatewayTunnelIPSecAuthTypePSK,
									PSK:            psk,
								},
							},
						},
					},
				},
			},
		})

		if !assert.NoError(t, err) {
			return
		}

		defer func() {
			err = svc.DeleteGateway(ctx, &request.DeleteGatewayRequest{UUID: gw.UUID})
			assert.NoError(t, err)

			err = waitGatewayToDelete(ctx, rec, svc, gw.UUID)
			assert.NoError(t, err)
		}()

		// Check plan
		assert.NotEmpty(t, gw.Plan)

		// Check addresses
		assert.Len(t, gw.Addresses, 1)
		assert.Equal(t, "my-public-ip", gw.Addresses[0].Name)

		// Check connections
		assert.Len(t, gw.Connections, 1)

		connection := gw.Connections[0]
		assert.NotEmpty(t, connection.UUID)
		assert.Equal(t, "example-connection", connection.Name)
		assert.Equal(t, upcloud.GatewayConnectionTypeIPSec, connection.Type)
		assert.Len(t, connection.LocalRoutes, 1)
		assert.Len(t, connection.RemoteRoutes, 1)
		assert.Len(t, connection.Tunnels, 1)

		// Check connection local routes
		localRoute := connection.LocalRoutes[0]
		assert.Equal(t, "local-route", localRoute.Name)
		assert.Equal(t, upcloud.GatewayRouteTypeStatic, localRoute.Type)
		assert.Equal(t, "10.0.0.0/24", localRoute.StaticNetwork)

		// Check connection remote routes
		remoteRoute := connection.RemoteRoutes[0]
		assert.Equal(t, "remote-route", remoteRoute.Name)
		assert.Equal(t, upcloud.GatewayRouteTypeStatic, remoteRoute.Type)
		assert.Equal(t, "10.0.1.0/24", remoteRoute.StaticNetwork)

		// Check connection tunnel
		tunnel := connection.Tunnels[0]
		assert.NotEmpty(t, tunnel.UUID)
		assert.Equal(t, "example-tunnel", tunnel.Name)
		assert.Equal(t, "my-public-ip", tunnel.LocalAddress.Name)
		assert.Equal(t, "100.10.0.111", tunnel.RemoteAddress.Address)
		assert.Equal(t, upcloud.GatewayTunnelIPSecAuthTypePSK, tunnel.IPSec.Authentication.Authentication)

		// Now check that we can modify just name, without affecting other fields
		gw, err = svc.ModifyGateway(ctx, &request.ModifyGatewayRequest{
			UUID: gw.UUID,
			Name: "updated",
		})

		assert.NoError(t, err)
		assert.Equal(t, "updated", gw.Name)
		assert.Len(t, gw.Addresses, 1)
		assert.Len(t, gw.Connections, 1)
		assert.Len(t, gw.Connections[0].Tunnels, 1)

		// Now let's see if we can modify other fields
		gw, err = svc.ModifyGateway(ctx, &request.ModifyGatewayRequest{
			UUID: gw.UUID,
			Connections: []request.GatewayConnection{
				{
					Name: "example-connection2",
					Type: upcloud.GatewayConnectionTypeIPSec,
					LocalRoutes: []upcloud.GatewayRoute{
						{
							Name:          "local-route2",
							Type:          upcloud.GatewayRouteTypeStatic,
							StaticNetwork: "11.0.0.0/24",
						},
					},
					RemoteRoutes: []upcloud.GatewayRoute{
						{
							Name:          "remote-route2",
							Type:          upcloud.GatewayRouteTypeStatic,
							StaticNetwork: "11.0.1.0/24",
						},
					},
					Tunnels: []request.GatewayTunnel{
						{
							Name: "example-tunnel2",
							LocalAddress: upcloud.GatewayTunnelLocalAddress{
								Name: "my-public-ip",
							},
							RemoteAddress: upcloud.GatewayTunnelRemoteAddress{
								Address: "200.10.0.111",
							},
							IPSec: upcloud.GatewayTunnelIPSec{
								Authentication: upcloud.GatewayTunnelIPSecAuth{
									Authentication: upcloud.GatewayTunnelIPSecAuthTypePSK,
									PSK:            psk,
								},
							},
						},
					},
				},
			},
		})

		// Check connections
		assert.Len(t, gw.Connections, 1)

		connection = gw.Connections[0]
		assert.NotEmpty(t, connection.UUID)
		assert.Equal(t, "example-connection2", connection.Name)
		assert.Equal(t, upcloud.GatewayConnectionTypeIPSec, connection.Type)
		assert.Len(t, connection.LocalRoutes, 1)
		assert.Len(t, connection.RemoteRoutes, 1)
		assert.Len(t, connection.Tunnels, 1)

		// Check connection local routes
		localRoute = connection.LocalRoutes[0]
		assert.Equal(t, "local-route2", localRoute.Name)
		assert.Equal(t, upcloud.GatewayRouteTypeStatic, localRoute.Type)
		assert.Equal(t, "11.0.0.0/24", localRoute.StaticNetwork)

		// Check connection remote routes
		remoteRoute = connection.RemoteRoutes[0]
		assert.Equal(t, "remote-route2", remoteRoute.Name)
		assert.Equal(t, upcloud.GatewayRouteTypeStatic, remoteRoute.Type)
		assert.Equal(t, "11.0.1.0/24", remoteRoute.StaticNetwork)

		// Check connection tunnel
		tunnel = connection.Tunnels[0]
		assert.NotEmpty(t, tunnel.UUID)
		assert.Equal(t, "example-tunnel2", tunnel.Name)
		assert.Equal(t, "my-public-ip", tunnel.LocalAddress.Name)
		assert.Equal(t, "200.10.0.111", tunnel.RemoteAddress.Address)
		assert.Equal(t, upcloud.GatewayTunnelIPSecAuthTypePSK, tunnel.IPSec.Authentication.Authentication)
	})
}

func TestVPNGatewayConnections(t *testing.T) {
	t.Parallel()

	record(t, "gatewayvpnconnections", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{Name: "test-router-vpn-conn"})
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			err = svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: router.UUID})
			assert.NoError(t, err)
		}()

		gw, err := svc.CreateGateway(ctx, &request.CreateGatewayRequest{
			Name: "test-vpn-conn",
			Zone: "pl-waw1",
			Routers: []request.GatewayRouter{
				{
					UUID: router.UUID,
				},
			},
			Plan:             "advanced",
			Features:         []upcloud.GatewayFeature{upcloud.GatewayFeatureVPN},
			ConfiguredStatus: upcloud.GatewayConfiguredStatusStarted,
			Addresses: []upcloud.GatewayAddress{
				{
					Name: "my-ip-address",
				},
			},
		})

		if !assert.NoError(t, err) {
			return
		}

		defer func() {
			err = svc.DeleteGateway(ctx, &request.DeleteGatewayRequest{UUID: gw.UUID})
			assert.NoError(t, err)

			err = waitGatewayToDelete(ctx, rec, svc, gw.UUID)
			assert.NoError(t, err)
		}()

		// Now try to create a connection
		conn, err := svc.CreateGatewayConnection(ctx, &request.CreateGatewayConnectionRequest{
			ServiceUUID: gw.UUID,
			Connection: request.GatewayConnection{
				Name: "added-connection",
				Type: upcloud.GatewayConnectionTypeIPSec,
				LocalRoutes: []upcloud.GatewayRoute{
					{
						Name:          "local-route",
						Type:          upcloud.GatewayRouteTypeStatic,
						StaticNetwork: "10.0.1.0/24",
					},
				},
				RemoteRoutes: []upcloud.GatewayRoute{
					{
						Name:          "remote-route",
						Type:          upcloud.GatewayRouteTypeStatic,
						StaticNetwork: "10.0.2.0/24",
					},
				},
				Tunnels: []request.GatewayTunnel{
					{
						Name: "added-tunnel",
						LocalAddress: upcloud.GatewayTunnelLocalAddress{
							Name: gw.Addresses[0].Name,
						},
						RemoteAddress: upcloud.GatewayTunnelRemoteAddress{
							Address: "100.10.0.111",
						},
						IPSec: upcloud.GatewayTunnelIPSec{
							Authentication: upcloud.GatewayTunnelIPSecAuth{
								Authentication: upcloud.GatewayTunnelIPSecAuthTypePSK,
								PSK:            "psk1234567890psk1234567890",
							},
						},
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, conn.UUID)
		assert.Equal(t, "added-connection", conn.Name)
		assert.Len(t, conn.LocalRoutes, 1)
		assert.Len(t, conn.RemoteRoutes, 1)
		assert.Len(t, conn.Tunnels, 1)
		assert.Equal(t, "local-route", conn.LocalRoutes[0].Name)
		assert.Equal(t, "remote-route", conn.RemoteRoutes[0].Name)
		assert.Equal(t, "added-tunnel", conn.Tunnels[0].Name)

		// Now check if listing and getting the details work
		connList, err := svc.GetGatewayConnections(ctx, &request.GetGatewayConnectionsRequest{ServiceUUID: gw.UUID})
		assert.NoError(t, err)
		assert.Len(t, connList, 1)

		connDetails, err := svc.GetGatewayConnection(ctx, &request.GetGatewayConnectionRequest{
			ServiceUUID: gw.UUID,
			UUID:        conn.UUID,
		})
		assert.NoError(t, err)
		assert.Equal(t, conn.Name, connDetails.Name)

		// Test modifying a connection
		conn, err = svc.ModifyGatewayConnection(ctx, &request.ModifyGatewayConnectionRequest{
			ServiceUUID: gw.UUID,
			UUID:        conn.UUID,
			Connection: request.ModifyGatewayConnection{
				RemoteRoutes: []upcloud.GatewayRoute{
					{
						Name:          "remote-route-updated",
						Type:          upcloud.GatewayRouteTypeStatic,
						StaticNetwork: "10.0.2.0/24",
					},
				},
				Tunnels: []request.GatewayTunnel{
					{
						Name: "added-tunnel-updated",
						LocalAddress: upcloud.GatewayTunnelLocalAddress{
							Name: gw.Addresses[0].Name,
						},
						RemoteAddress: upcloud.GatewayTunnelRemoteAddress{
							Address: "100.10.0.111",
						},
						IPSec: upcloud.GatewayTunnelIPSec{
							Authentication: upcloud.GatewayTunnelIPSecAuth{
								Authentication: upcloud.GatewayTunnelIPSecAuthTypePSK,
								PSK:            "psk1234567890psk1234567890",
							},
						},
					},
				},
			},
		})

		assert.NoError(t, err)
		assert.Equal(t, "local-route", conn.LocalRoutes[0].Name) // LocalRoutes should stay the same as before, they weren't modified
		assert.Equal(t, "remote-route-updated", conn.RemoteRoutes[0].Name)
		assert.Equal(t, "added-tunnel-updated", conn.Tunnels[0].Name)

		// Now delete the connection
		err = svc.DeleteGatewayConnection(ctx, &request.DeleteGatewayConnectionRequest{
			ServiceUUID: gw.UUID,
			UUID:        conn.UUID,
		})
		assert.NoError(t, err)
	})
}

func TestVPNGatewayConnectionTunnels(t *testing.T) {
	t.Parallel()

	record(t, "gatewayvpnconnectiontunnels", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{Name: "test-router-vpn-conn-tunnels"})
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			err = svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: router.UUID})
			assert.NoError(t, err)
		}()

		gw, err := svc.CreateGateway(ctx, &request.CreateGatewayRequest{
			Name: "test-vpn-conn-tunnels",
			Zone: "pl-waw1",
			Routers: []request.GatewayRouter{
				{
					UUID: router.UUID,
				},
			},
			Plan:             "advanced",
			Features:         []upcloud.GatewayFeature{upcloud.GatewayFeatureVPN},
			ConfiguredStatus: upcloud.GatewayConfiguredStatusStarted,
			Addresses: []upcloud.GatewayAddress{
				{
					Name: "my-ip-address",
				},
			},
			Connections: []request.GatewayConnection{
				{
					Name: "example-conn",
					Type: upcloud.GatewayConnectionTypeIPSec,
					LocalRoutes: []upcloud.GatewayRoute{
						{
							Name:          "local-route",
							Type:          upcloud.GatewayRouteTypeStatic,
							StaticNetwork: "10.0.1.0/24",
						},
					},
					RemoteRoutes: []upcloud.GatewayRoute{
						{
							Name:          "remote-route",
							Type:          upcloud.GatewayRouteTypeStatic,
							StaticNetwork: "10.0.2.0/24",
						},
					},
				},
			},
		})

		if !assert.NoError(t, err) {
			return
		}

		defer func() {
			err = svc.DeleteGateway(ctx, &request.DeleteGatewayRequest{UUID: gw.UUID})
			assert.NoError(t, err)

			err = waitGatewayToDelete(ctx, rec, svc, gw.UUID)
			assert.NoError(t, err)
		}()

		// First get tunnels to see if it handles empty response
		tunnels, err := svc.GetGatewayConnectionTunnels(ctx, &request.GetGatewayConnectionTunnelsRequest{
			ServiceUUID:    gw.UUID,
			ConnectionUUID: gw.Connections[0].UUID,
		})

		assert.NoError(t, err)
		assert.Len(t, tunnels, 0)

		// Now let's create a tunnel
		tunnel, err := svc.CreateGatewayConnectionTunnel(ctx, &request.CreateGatewayConnectionTunnelRequest{
			ServiceUUID:    gw.UUID,
			ConnectionUUID: gw.Connections[0].UUID,
			Tunnel: request.GatewayTunnel{
				Name: "added-tunnel",
				LocalAddress: upcloud.GatewayTunnelLocalAddress{
					Name: gw.Addresses[0].Name,
				},
				RemoteAddress: upcloud.GatewayTunnelRemoteAddress{
					Address: "100.10.0.111",
				},
				IPSec: upcloud.GatewayTunnelIPSec{
					Authentication: upcloud.GatewayTunnelIPSecAuth{
						Authentication: upcloud.GatewayTunnelIPSecAuthTypePSK,
						PSK:            "psk1234567890psk1234567890",
					},
				},
			},
		})

		assert.NoError(t, err)
		assert.Equal(t, "added-tunnel", tunnel.Name)
		assert.Equal(t, gw.Addresses[0].Name, tunnel.LocalAddress.Name)
		assert.Equal(t, "100.10.0.111", tunnel.RemoteAddress.Address)
		assert.Equal(t, upcloud.GatewayTunnelIPSecAuthTypePSK, tunnel.IPSec.Authentication.Authentication)

		// Check if listing tunnels work
		tunnels, err = svc.GetGatewayConnectionTunnels(ctx, &request.GetGatewayConnectionTunnelsRequest{
			ServiceUUID:    gw.UUID,
			ConnectionUUID: gw.Connections[0].UUID,
		})
		assert.NoError(t, err)
		assert.Len(t, tunnels, 1)

		// Check if getting details work
		tunnel, err = svc.GetGatewayConnectionTunnel(ctx, &request.GetGatewayConnectionTunnelRequest{
			ServiceUUID:    gw.UUID,
			ConnectionUUID: gw.Connections[0].UUID,
			UUID:           tunnels[0].UUID,
		})
		assert.NoError(t, err)
		assert.Equal(t, "100.10.0.111", tunnel.RemoteAddress.Address)

		// Try modifying tunnel
		updatedTunnelName := "updated-tunnel"
		tunnel, err = svc.ModifyGatewayConnectionTunnel(ctx, &request.ModifyGatewayConnectionTunnelRequest{
			ServiceUUID:    gw.UUID,
			ConnectionUUID: gw.Connections[0].UUID,
			UUID:           tunnel.UUID,
			Tunnel: request.ModifyGatewayTunnel{
				Name: updatedTunnelName,
				RemoteAddress: &upcloud.GatewayTunnelRemoteAddress{
					Address: "100.10.0.222",
				},
				IPSec: &upcloud.GatewayTunnelIPSec{
					Authentication: upcloud.GatewayTunnelIPSecAuth{
						PSK: "updatedsuperduperpsk1234566778",
					},
				},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, gw.Addresses[0].Name, tunnel.LocalAddress.Name)
		assert.Equal(t, "100.10.0.222", tunnel.RemoteAddress.Address)

		// Delete the tunnel
		err = svc.DeleteGatewayConnectionTunnel(ctx, &request.DeleteGatewayConnectionTunnelRequest{
			ServiceUUID:    gw.UUID,
			ConnectionUUID: gw.Connections[0].UUID,
			UUID:           tunnel.UUID,
		})
		assert.NoError(t, err)
	})
}

func waitGatewayToStart(ctx context.Context, rec *recorder.Recorder, svc *Service, UUID string) error {
	if !rec.IsRecording() {
		return nil
	}

	const timeout = 10 * time.Minute

	waitUntil := time.Now().Add(timeout)
	for {
		gw, err := svc.GetGateway(ctx, &request.GetGatewayRequest{UUID: UUID})
		if err != nil {
			return err
		}
		if gw.OperationalState == upcloud.GatewayOperationalStateRunning {
			return nil
		}
		if time.Now().After(waitUntil) {
			return fmt.Errorf("timeout %s reached", timeout.String())
		}
		time.Sleep(5 * time.Second)
	}
}

func waitGatewayToDelete(ctx context.Context, rec *recorder.Recorder, svc *Service, UUID string) error {
	if !rec.IsRecording() {
		return nil
	}

	const timeout = 10 * time.Minute

	waitUntil := time.Now().Add(timeout)
	for {
		_, err := svc.GetGateway(ctx, &request.GetGatewayRequest{UUID: UUID})
		if err != nil {
			var ucErr *upcloud.Problem
			if errors.As(err, &ucErr) && ucErr.Status == http.StatusNotFound {
				return nil
			}
			log.Printf("ERROR: %+v", err)
			return err
		}
		if time.Now().After(waitUntil) {
			return fmt.Errorf("timeout %s reached", timeout.String())
		}
		time.Sleep(5 * time.Second)
	}
}
