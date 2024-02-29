package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupNetworkPeeringTest(handler http.Handler) (*httptest.Server, *Service) {
	srv := httptest.NewServer(handler)
	return srv, New(client.New("user", "pass", client.WithBaseURL(srv.URL)))
}

func TestGetNetworkPeerings(t *testing.T) {
	t.Parallel()

	srv, svc := setupNetworkPeeringTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/network-peering"), r.URL.Path)
		fmt.Fprint(w, `
		{
			"network_peerings": {
			  "network_peering": [
				{
				  "configured_status": "active",
				  "name": "Peering A->B",
				  "network": {
					"ip_networks": {
					  "ip_network": [
						{
						  "address": "192.168.0.0/24",
						  "family": "IPv4"
						},
						{
						  "address": "fc02:c4f3::/64",
						  "family": "IPv6"
						}
					  ]
					},
					"uuid": "03126dc1-a69f-4bc2-8b24-e31c22d64712"
				  },
				  "peer_network": {
					"ip_networks": {
					  "ip_network": [
						{
						  "address": "192.168.99.0/24",
						  "family": "IPv4"
						},
						{
						  "address": "fc02:c4f3:99::/64",
						  "family": "IPv6"
						}
					  ]
					},
					"uuid": "03585987-bf7d-4544-8e9b-5a1b4d74a333"
				  },
				  "state": "active",
				  "uuid": "0f7984bc-5d72-4aaf-b587-90e6a8f32efc"
				}
			  ]
			}
		}
		`)
	}))
	defer srv.Close()
	p, err := svc.GetNetworkPeerings(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, p, 1)
	checkNetworkPeeringResponse(t, &p[0])
}

func TestGetNetworkPeering(t *testing.T) {
	t.Parallel()

	srv, svc := setupNetworkPeeringTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/network-peering/_UUID_"), r.URL.Path)
		fmt.Fprint(w, networkPeeringCommonResponse)
	}))
	defer srv.Close()
	p, err := svc.GetNetworkPeering(context.TODO(), &request.GetNetworkPeeringRequest{UUID: "_UUID_"})
	assert.NoError(t, err)
	assert.Equal(t, "03585987-bf7d-4544-8e9b-5a1b4d74a333", p.PeerNetwork.UUID)
	assert.Equal(t, "03126dc1-a69f-4bc2-8b24-e31c22d64712", p.Network.UUID)
	assert.Equal(t, "192.168.99.0/24", p.PeerNetwork.IPNetworks[0].Address)
	assert.Equal(t, "fc02:c4f3::/64", p.Network.IPNetworks[1].Address)
}

func TestCreateNetworkPeering(t *testing.T) {
	t.Parallel()

	srv, svc := setupNetworkPeeringTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/network-peering"), r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `
		{
			"network_peering": {
				"configured_status": "active",
				"name": "Peering A->B",
				"network": {
					"uuid": "03126dc1-a69f-4bc2-8b24-e31c22d64712"
				},
				"peer_network": {
					"uuid": "03585987-bf7d-4544-8e9b-5a1b4d74a333"
				},
				"labels": [
					{
						"key": "managedBy",
						"value": "upcloud-go-sdk-unit-test"
					}
			  	]
			}
		}`, string(body))
		fmt.Fprint(w, networkPeeringCommonResponse)
	}))
	defer srv.Close()

	p, err := svc.CreateNetworkPeering(context.TODO(), &request.CreateNetworkPeeringRequest{
		Name:             "Peering A->B",
		ConfiguredStatus: upcloud.NetworkPeeringConfiguredStatusActive,
		Network: request.NetworkPeeringNetwork{
			UUID: "03126dc1-a69f-4bc2-8b24-e31c22d64712",
		},
		PeerNetwork: request.NetworkPeeringNetwork{
			UUID: "03585987-bf7d-4544-8e9b-5a1b4d74a333",
		},
		Labels: []upcloud.Label{
			{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
	})
	if !assert.NoError(t, err) {
		return
	}
	checkNetworkPeeringResponse(t, p)
}

func TestModifyNetworkPeering(t *testing.T) {
	t.Parallel()

	srv, svc := setupNetworkPeeringTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/network-peering/_UUID_"), r.URL.Path)
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `
		{
			"network_peering": {
			  "configured_status": "disabled",
			  "name": "Peering A->B modified",
			  "labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			  ]
			}
		}
		`, string(body))
		fmt.Fprint(w, networkPeeringCommonResponse)
	}))
	defer srv.Close()

	p, err := svc.ModifyNetworkPeering(context.TODO(), &request.ModifyNetworkPeeringRequest{
		UUID: "_UUID_",
		NetworkPeering: request.ModifyNetworkPeering{
			Name:             "Peering A->B modified",
			ConfiguredStatus: upcloud.NetworkPeeringConfiguredStatusDisabled,
			Labels: &[]upcloud.Label{
				{
					Key:   "managedBy",
					Value: "upcloud-go-sdk-unit-test",
				},
			},
		},
	})
	if !assert.NoError(t, err) {
		return
	}
	checkNetworkPeeringResponse(t, p)
}

func TestDeleteNetworkPeering(t *testing.T) {
	t.Parallel()

	srv, svc := setupNetworkPeeringTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s%s", client.APIVersion, "/network-peering/_UUID_"), r.URL.Path)
	}))
	defer srv.Close()
	assert.NoError(t, svc.DeleteNetworkPeering(context.TODO(), &request.DeleteNetworkPeeringRequest{UUID: "_UUID_"}))
}

func checkNetworkPeeringResponse(t *testing.T, p *upcloud.NetworkPeering) {
	assert.Equal(t, "03585987-bf7d-4544-8e9b-5a1b4d74a333", p.PeerNetwork.UUID)
	assert.Equal(t, "03126dc1-a69f-4bc2-8b24-e31c22d64712", p.Network.UUID)
	assert.Equal(t, "192.168.99.0/24", p.PeerNetwork.IPNetworks[0].Address)
	assert.Equal(t, "fc02:c4f3::/64", p.Network.IPNetworks[1].Address)
}

const networkPeeringCommonResponse string = `
{
	"network_peering": {
		"state": "active",
		"uuid": "0f7984bc-5d72-4aaf-b587-90e6a8f32efc",
		"configured_status": "active",
		"name": "Peering A->B",
		"labels": [
			{
				"key": "managedBy",
				"value": "upcloud-go-sdk-unit-test"
			}
		],
		"network": {
			"ip_networks": {
				"ip_network": [
					{
						"address": "192.168.0.0/24",
						"family": "IPv4"
					},
					{
						"address": "fc02:c4f3::/64",
						"family": "IPv6"
					}
				]
			},
			"uuid": "03126dc1-a69f-4bc2-8b24-e31c22d64712"
		},
		"peer_network": {
			"ip_networks": {
			"ip_network": [
				{
					"address": "192.168.99.0/24",
					"family": "IPv4"
				},
				{
					"address": "fc02:c4f3:99::/64",
					"family": "IPv6"
				}
			]
			},
			"uuid": "03585987-bf7d-4544-8e9b-5a1b4d74a333"
		}
	}
}
`
