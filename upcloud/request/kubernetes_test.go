package request

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

const exampleCreateKubernetesClusterRequestJSON string = `{
	"control_plane_ip_filter": null,
	"network": "00000000-0000-0000-0000-000000000000",
	"network_cidr": "172.16.0.1/24",
	"plan": "production",
	"private_node_groups": false,
	"node_groups": [
		{
			"count": 1,
			"labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			],
			"kubelet_args": [
				{
					"key": "somekubeletkey",
					"value": "somekubeletvalue"
				}
			],
			"taints": [
				{
					"effect": "NoExecute",
					"key": "sometaintkey",
					"value": "sometaintvalue"
				}
			],
			"name": "withExplicitUtilityTrue",
			"plan": "plan",
			"ssh_keys": [
				"key",
				"key"
			],
			"utility_network_access": true
		},
		{
			"count": 1,
			"labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			],
			"kubelet_args": [
				{
					"key": "somekubeletkey",
					"value": "somekubeletvalue"
				}
			],
			"taints": [
				{
					"effect": "NoExecute",
					"key": "sometaintkey",
					"value": "sometaintvalue"
				}
			],
			"name": "withExplicitUtilityFalse",
			"plan": "plan",
			"ssh_keys": [
				"key",
				"key"
			],
			"utility_network_access": false
		},
		{
			"count": 1,
			"labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			],
			"kubelet_args": [
				{
					"key": "somekubeletkey",
					"value": "somekubeletvalue"
				}
			],
			"taints": [
				{
					"effect": "NoExecute",
					"key": "sometaintkey",
					"value": "sometaintvalue"
				}
			],
			"name": "withoutExplicitUtility",
			"plan": "plan",
			"ssh_keys": [
				"key",
				"key"
			]
		}
	],
	"name": "title",
	"version": "1.26",
	"zone": "zone"
}`

func TestKubernetes(t *testing.T) {
	t.Parallel()

	t.Run("GetKubernetesClustersRequest", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, kubernetesClusterBasePath, (&GetKubernetesClustersRequest{}).RequestURL())
	})

	t.Run("GetKubernetesClustersRequestWithFilters", func(t *testing.T) {
		t.Parallel()

		expected := fmt.Sprintf(
			"%s?%s&%s",
			kubernetesClusterBasePath,
			"label=managedBy",
			"label=managedBy%3Dupcloud-go-sdk-unit-test",
		)

		r := exampleGetKubernetesClustersWithFiltersRequest()
		actual := r.RequestURL()

		assert.Equal(
			t,
			expected,
			actual,
		)
	})

	t.Run("GetKubernetesClusterRequest", func(t *testing.T) {
		t.Parallel()

		expected := fmt.Sprintf(
			"%s/%s",
			kubernetesClusterBasePath,
			"00000000-0000-0000-0000-000000000000",
		)

		r := GetKubernetesClusterRequest{
			UUID: "00000000-0000-0000-0000-000000000000",
		}
		actual := r.RequestURL()

		assert.Equal(
			t,
			expected,
			actual,
		)
	})

	t.Run("CreateKubernetesClusterRequestMarshal", func(t *testing.T) {
		t.Parallel()

		expected := exampleCreateKubernetesClusterRequestJSON

		r := exampleCreateKubernetesClusterRequest()
		b, err := json.Marshal(r)
		actual := string(b)

		assert.NoError(t, err)
		assert.JSONEq(
			t,
			expected,
			actual,
		)
	})

	t.Run("ModifyKubernetesClusterRequestMarshal", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			request  ModifyKubernetesClusterRequest
			expected string
		}{
			{
				request: ModifyKubernetesClusterRequest{
					ClusterUUID: "set-filter-omit-labels",
					Cluster: ModifyKubernetesCluster{
						ControlPlaneIPFilter: &[]string{"0.0.0.0/0"},
					},
				},
				expected: `{ "control_plane_ip_filter": ["0.0.0.0/0"] }`,
			},
			{
				request: ModifyKubernetesClusterRequest{
					ClusterUUID: "omit-filter-set-labels",
					Cluster: ModifyKubernetesCluster{
						Labels: &[]upcloud.Label{{Key: "tool", Value: "Go SDK"}},
					},
				},
				expected: `{ "labels": [{"key": "tool", "value": "Go SDK"}] }`,
			},
			{
				request: ModifyKubernetesClusterRequest{
					ClusterUUID: "clear-filter-clear-labels",
					Cluster: ModifyKubernetesCluster{
						ControlPlaneIPFilter: &[]string{},
						Labels:               &[]upcloud.Label{},
					},
				},
				expected: `{ "control_plane_ip_filter": [], "labels": [] }`,
			},
		}
		for _, test := range tests {
			assert.Equal(t, fmt.Sprintf("%s/%s", kubernetesClusterBasePath, test.request.ClusterUUID), test.request.RequestURL())
			test := test
			b, err := json.Marshal(&test.request)
			actual := string(b)

			assert.NoError(t, err)
			assert.JSONEq(
				t,
				test.expected,
				actual,
			)
		}
	})

	t.Run("DeleteKubernetesClusterRequest", func(t *testing.T) {
		t.Parallel()

		expected := fmt.Sprintf(
			"%s/%s",
			kubernetesClusterBasePath,
			"00000000-0000-0000-0000-000000000000",
		)

		r := DeleteKubernetesClusterRequest{
			UUID: "00000000-0000-0000-0000-000000000000",
		}
		actual := r.RequestURL()

		assert.Equal(
			t,
			expected,
			actual,
		)
	})

	t.Run("GetKubernetesKubeconfigRequest", func(t *testing.T) {
		t.Parallel()

		expected := fmt.Sprintf(
			"%s/%s/kubeconfig",
			kubernetesClusterBasePath,
			"00000000-0000-0000-0000-000000000000",
		)

		r := GetKubernetesKubeconfigRequest{
			UUID: "00000000-0000-0000-0000-000000000000",
		}
		actual := r.RequestURL()

		assert.Equal(
			t,
			expected,
			actual,
		)
	})

	t.Run("GetKubernetesVersionsRequest", func(t *testing.T) {
		t.Parallel()

		expected := fmt.Sprintf(
			"%s/versions",
			kubernetesClusterBasePath,
		)

		r := GetKubernetesVersionsRequest{}
		actual := r.RequestURL()

		assert.Equal(
			t,
			expected,
			actual,
		)
	})

	t.Run("WaitForKubernetesClusterStateRequest", func(t *testing.T) {
		t.Parallel()

		expected := fmt.Sprintf(
			"%s/%s",
			kubernetesClusterBasePath,
			"00000000-0000-0000-0000-000000000000",
		)

		r := WaitForKubernetesClusterStateRequest{
			UUID: "00000000-0000-0000-0000-000000000000",
		}
		actual := r.RequestURL()

		assert.Equal(
			t,
			expected,
			actual,
		)
	})
}

func TestGetKubernetesNodeGroupsRequest(t *testing.T) {
	t.Parallel()
	r := GetKubernetesNodeGroupsRequest{ClusterUUID: "id"}
	assert.Equal(t, fmt.Sprintf("%s/id/node-groups", kubernetesClusterBasePath), r.RequestURL())
}

func TestGetKubernetesNodeGroupRequest(t *testing.T) {
	t.Parallel()
	r := GetKubernetesNodeGroupRequest{ClusterUUID: "id", Name: "nid"}
	assert.Equal(t, fmt.Sprintf("%s/id/node-groups/nid", kubernetesClusterBasePath), r.RequestURL())
}

func TestCreateKubernetesNodeGroupRequest(t *testing.T) {
	t.Parallel()
	const expectedJSON string = `
	{
		"count": 4,
		"anti_affinity": true,
		"kubelet_args": [
		  {
			"key": "log-flush-frequency",
			"value": "5s"
		  }
		],
		"labels": [
		  {
			"key": "environment",
			"value": "development"
		  }
		],
		"name": "small",
		"plan": "K8S-2xCPU-4GB",
		"ssh_keys": [
		  "ssh-rsa AAAA.."
		],
		"storage": "01000000-0000-4000-8000-000160010100",
		"taints": [
		  {
			"effect": "NoSchedule",
			"key": "environment",
			"value": "development"
		  }
		]
	}
	`
	r := CreateKubernetesNodeGroupRequest{
		ClusterUUID: "id",
		NodeGroup: KubernetesNodeGroup{
			Count:        4,
			AntiAffinity: true,
			Labels: []upcloud.Label{
				{
					Key:   "environment",
					Value: "development",
				},
			},
			Name:    "small",
			Plan:    "K8S-2xCPU-4GB",
			SSHKeys: []string{"ssh-rsa AAAA.."},
			Storage: "01000000-0000-4000-8000-000160010100",
			KubeletArgs: []upcloud.KubernetesKubeletArg{
				{
					Key:   "log-flush-frequency",
					Value: "5s",
				},
			},
			Taints: []upcloud.KubernetesTaint{
				{
					Effect: "NoSchedule",
					Key:    "environment",
					Value:  "development",
				},
			},
		},
	}
	assert.Equal(t, fmt.Sprintf("%s/id/node-groups", kubernetesClusterBasePath), r.RequestURL())
	gotJS, err := json.Marshal(&r)
	if !assert.NoError(t, err) {
		return
	}
	assert.JSONEq(t, expectedJSON, string(gotJS))
}

func TestModifyKubernetesNodeGroupRequest(t *testing.T) {
	t.Parallel()
	const expectedJSON string = `
	{
		"count": 4
	}
	`
	r := ModifyKubernetesNodeGroupRequest{
		ClusterUUID: "id",
		Name:        "nid",
		NodeGroup: ModifyKubernetesNodeGroup{
			Count: 4,
		},
	}
	assert.Equal(t, fmt.Sprintf("%s/id/node-groups/nid", kubernetesClusterBasePath), r.RequestURL())
	gotJS, err := json.Marshal(&r)
	if !assert.NoError(t, err) {
		return
	}
	assert.JSONEq(t, expectedJSON, string(gotJS))
}

func TestDeleteKubernetesNodeGroupNodeRequest(t *testing.T) {
	r := DeleteKubernetesNodeGroupNodeRequest{
		ClusterUUID: "id",
		Name:        "nid",
		NodeName:    "name",
	}
	assert.Equal(t, fmt.Sprintf("%s/id/node-groups/nid/name", kubernetesClusterBasePath), r.RequestURL())
}

func TestDeleteKubernetesNodeGroupRequest(t *testing.T) {
	t.Parallel()
	r := DeleteKubernetesNodeGroupRequest{ClusterUUID: "id", Name: "nid"}
	assert.Equal(t, fmt.Sprintf("%s/id/node-groups/nid", kubernetesClusterBasePath), r.RequestURL())
}

func exampleGetKubernetesClustersWithFiltersRequest() GetKubernetesClustersWithFiltersRequest {
	return GetKubernetesClustersWithFiltersRequest{
		Filters: []QueryFilter{
			FilterLabelKey{"managedBy"},
			FilterLabel{Label: upcloud.Label{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			}},
		},
	}
}

func exampleCreateKubernetesClusterRequest() CreateKubernetesClusterRequest {
	return CreateKubernetesClusterRequest{
		Name:        "title",
		Network:     "00000000-0000-0000-0000-000000000000",
		NetworkCIDR: "172.16.0.1/24",
		NodeGroups: []KubernetesNodeGroup{
			exampleKubernetesNodeGroup("withExplicitUtilityTrue", upcloud.BoolPtr(true)),
			exampleKubernetesNodeGroup("withExplicitUtilityFalse", upcloud.BoolPtr(false)),
			exampleKubernetesNodeGroup("withoutExplicitUtility", nil),
		},
		Version:           "1.26",
		Zone:              "zone",
		Plan:              "production",
		PrivateNodeGroups: false,
	}
}

func exampleKubernetesNodeGroup(name string, utilityNetworkAccess *bool) KubernetesNodeGroup {
	ng := KubernetesNodeGroup{
		Count: 1,
		Labels: []upcloud.Label{
			{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
		KubeletArgs: []upcloud.KubernetesKubeletArg{
			{
				Key:   "somekubeletkey",
				Value: "somekubeletvalue",
			},
		},
		Taints: []upcloud.KubernetesTaint{
			{
				Effect: "NoExecute",
				Key:    "sometaintkey",
				Value:  "sometaintvalue",
			},
		},
		Name: name,
		Plan: "plan",
		SSHKeys: []string{
			"key",
			"key",
		},
	}

	if utilityNetworkAccess != nil {
		ng.UtilityNetworkAccess = utilityNetworkAccess
	}

	return ng
}
