package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleKubernetesClusterJSON string = `{
	"name": "upcloud-go-sdk-unit-test",
	"zone": "de-fra1",
	"uuid": "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
	"state": "running",
	"network": "03a98be3-7daa-443f-bb25-4bc6854b396c",
	"network_cidr": "172.16.0.0/24",
	"node_groups": [
		{
			"name": "upcloud-go-sdk-unit-test",
			"plan": "K8S-2xCPU-4GB",
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
			"ssh_keys": ["somekey"]
		},
		{
			"name": "upcloud-go-sdk-unit-test",
			"plan": "K8S-2xCPU-4GB",
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
			"ssh_keys": ["somekey"]
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

func exampleKubernetesCluster() KubernetesCluster {
	return KubernetesCluster{
		Name:        "upcloud-go-sdk-unit-test",
		Network:     "03a98be3-7daa-443f-bb25-4bc6854b396c",
		NetworkCIDR: "172.16.0.0/24",
		NodeGroups: []KubernetesNodeGroup{
			exampleKubernetesNodeGroup(),
			exampleKubernetesNodeGroup(),
		},
		State: KubernetesClusterStateRunning,
		UUID:  "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
		Zone:  "de-fra1",
	}
}

func exampleKubernetesNodeGroup() KubernetesNodeGroup {
	return KubernetesNodeGroup{
		Count: 1,
		Labels: []Label{
			{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
		Name: "upcloud-go-sdk-unit-test",
		Plan: "K8S-2xCPU-4GB",
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
		SSHKeys: []string{"somekey"},
	}
}
