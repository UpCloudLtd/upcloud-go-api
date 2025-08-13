package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const exampleKubernetesClusterJSON string = `{
	"control_plane_ip_filter": ["0.0.0.0/0"],
	"name": "upcloud-go-sdk-unit-test",
	"version": "1.27",
	"zone": "de-fra1",
	"uuid": "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
	"state": "running",
	"network": "03a98be3-7daa-443f-bb25-4bc6854b396c",
	"network_cidr": "172.16.0.0/24",
	"plan": "development",
	"private_node_groups": false,
	"node_groups": [
		{
			"name": "upcloud-go-sdk-unit-test",
			"plan": "K8S-2xCPU-4GB",
			"count": 1,
			"anti_affinity": true,
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
			"ssh_keys": ["somekey"],
			"utility_network_access": true
		},
		{
			"name": "upcloud-go-sdk-unit-test",
			"plan": "K8S-2xCPU-4GB",
			"count": 1,
			"anti_affinity": false,
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
			"ssh_keys": ["somekey"],
			"utility_network_access": true
		}
	]
}`

func TestKubernetes(t *testing.T) {
	t.Parallel()

	t.Run("KubernetesClusterUnMarshal", func(t *testing.T) {
		t.Parallel()

		expected := exampleKubernetesCluster()

		s := exampleKubernetesClusterJSON
		actual := KubernetesCluster{}
		err := json.Unmarshal([]byte(s), &actual)

		assert.NoError(t, err)
		assert.Equal(
			t,
			expected,
			actual,
		)
	})
}

func TestKubernetesNodeGroupDetails(t *testing.T) {
	t.Parallel()

	const nodeGroupDetailsJSON = `
	{
		"anti_affinity": true,
		"count": 2,
		"name": "grp-1",
		"plan": "1xCPU-1GB",
		"state": "running",
		"nodes": [
			{
				"name": "grp-1-7l7zj",
				"state": "running",
				"uuid": "00a02bfa-f565-40c9-b088-f2c7b8a75f97"
			},
			{
				"name": "grp-1-glkwv",
				"state": "terminating",
				"uuid": "00b56302-e211-40d9-83fa-177f0171e75a"
			}
		]
	}
	`
	got := KubernetesNodeGroupDetails{}
	err := json.Unmarshal([]byte(nodeGroupDetailsJSON), &got)
	want := KubernetesNodeGroupDetails{
		KubernetesNodeGroup: KubernetesNodeGroup{
			AntiAffinity: true,
			Count:        2,
			Name:         "grp-1",
			Plan:         "1xCPU-1GB",
			State:        KubernetesNodeGroupStateRunning,
		},
		Nodes: []KubernetesNode{
			{
				UUID:  "00a02bfa-f565-40c9-b088-f2c7b8a75f97",
				Name:  "grp-1-7l7zj",
				State: KubernetesNodeStateRunning,
			},
			{
				UUID:  "00b56302-e211-40d9-83fa-177f0171e75a",
				Name:  "grp-1-glkwv",
				State: KubernetesNodeStateTerminating,
			},
		},
	}
	require.NoError(t, err)
	require.Equal(t, want, got)
	// just to check that embedded KubernetesNodeGroup fields are directly available
	require.Equal(t, KubernetesNodeGroupStateRunning, got.State)
}

func TestKubernetesNodeGroupEncryptedCustomPlan(t *testing.T) {
	t.Parallel()

	p := []byte(`
	{
		"plan": "custom",
		"storage_encryption": "data-at-rest",
		"custom_plan": {
			"cores": 4,
			"memory": 2048,
			"storage_size": 30,
			"storage_tier": "hdd"
		}
	}
	`)
	got := KubernetesNodeGroup{}
	require.NoError(t, json.Unmarshal(p, &got))
	want := KubernetesNodeGroup{
		Plan:              "custom",
		StorageEncryption: StorageEncryptionDataAtRest,
		CustomPlan: &KubernetesNodeGroupCustomPlan{
			Cores:       4,
			Memory:      2048,
			StorageSize: 30,
			StorageTier: KubernetesStorageTierHDD,
		},
	}
	require.Equal(t, want, got)
}

func TestKubernetesNodeGroupCloudNativePlan(t *testing.T) {
	t.Parallel()

	p := []byte(`
	{
		"plan": "CLOUDNATIVE-16xCPU-192GB",
		"cloud_native_plan": {
			"storage_size": 50,
			"storage_tier": "maxiops"
		}
	}
	`)
	got := KubernetesNodeGroup{}
	require.NoError(t, json.Unmarshal(p, &got))
	want := KubernetesNodeGroup{
		Plan: "CLOUDNATIVE-16xCPU-192GB",
		CloudNativePlan: &KubernetesNodeGroupCloudNativePlan{
			StorageSize: 50,
			StorageTier: KubernetesStorageTierMaxIOPS,
		},
	}
	require.Equal(t, want, got)
}

func TestKubernetesNodeGroupGPUPlan(t *testing.T) {
	t.Parallel()

	p := []byte(`
	{
		"plan": "GPU-12xCPU-128GB-1xL40s",
		"gpu_plan": {
			"storage_size": 100,
			"storage_tier": "maxiops"
		}
	}
	`)
	got := KubernetesNodeGroup{}
	require.NoError(t, json.Unmarshal(p, &got))
	want := KubernetesNodeGroup{
		Plan: "GPU-12xCPU-128GB-1xL40s",
		GPUPlan: &KubernetesNodeGroupGPUPlan{
			StorageSize: 100,
			StorageTier: KubernetesStorageTierMaxIOPS,
		},
	}
	require.Equal(t, want, got)
}

func TestKubernetesStorageEncryption(t *testing.T) {
	t.Parallel()

	p := []byte(`
	{
		"storage_encryption": "data-at-rest"
	}
	`)
	got := KubernetesCluster{}
	require.NoError(t, json.Unmarshal(p, &got))
	want := KubernetesCluster{
		StorageEncryption: StorageEncryptionDataAtRest,
	}
	require.Equal(t, want, got)
}

func TestKubernetesAvailableVersions(t *testing.T) {
	t.Parallel()

	p := []byte(`
	{
		"versions": ["1.31"]
	}
	`)
	got := KubernetesClusterAvailableUpgrades{}
	require.NoError(t, json.Unmarshal(p, &got))
	want := KubernetesClusterAvailableUpgrades{
		Versions: []string{"1.31"},
	}
	require.Equal(t, want, got)
}

func exampleKubernetesCluster() KubernetesCluster {
	return KubernetesCluster{
		ControlPlaneIPFilter: []string{"0.0.0.0/0"},
		Name:                 "upcloud-go-sdk-unit-test",
		Network:              "03a98be3-7daa-443f-bb25-4bc6854b396c",
		NetworkCIDR:          "172.16.0.0/24",
		NodeGroups: []KubernetesNodeGroup{
			exampleKubernetesNodeGroup(true),
			exampleKubernetesNodeGroup(false),
		},
		State:             KubernetesClusterStateRunning,
		UUID:              "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
		Version:           "1.27",
		Zone:              "de-fra1",
		Plan:              "development",
		PrivateNodeGroups: false,
	}
}

func exampleKubernetesNodeGroup(antiAffinity bool) KubernetesNodeGroup {
	return KubernetesNodeGroup{
		Count: 1,
		Labels: []Label{
			{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
		Name:         "upcloud-go-sdk-unit-test",
		AntiAffinity: antiAffinity,
		Plan:         "K8S-2xCPU-4GB",
		KubeletArgs: []KubernetesKubeletArg{
			{
				Key:   "somekubeletkey",
				Value: "somekubeletvalue",
			},
		},
		Taints: []KubernetesTaint{
			{
				Effect: "NoExecute",
				Key:    "sometaintkey",
				Value:  "sometaintvalue",
			},
		},
		SSHKeys:              []string{"somekey"},
		UtilityNetworkAccess: true,
	}
}
