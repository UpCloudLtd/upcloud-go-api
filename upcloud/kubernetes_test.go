package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleKubernetesClusterJSON string = `{
	"title": "upcloud-go-sdk-unit-test",
	"zone": "de-fra1",
	"uuid": "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
	"state":"ready",
	"version": "v1.23.5",
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
			]
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
			]
		}
	],
	"storage": "01000000-0000-4000-8000-000160010100",
	"kubelet_args": null,
	"labels": null,
	"type": "standalone"
}`

const exampleKubenetesPlanJSON string = `{
		"description": "K8S-2xCPU-4GB",
		"name": "small"
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

	t.Run("KubernetesPlanUnMarshal", func(t *testing.T) {
		t.Parallel()

		expected := KubernetesPlan{
			Description: "K8S-2xCPU-4GB",
			Name:        "small",
		}

		s := exampleKubenetesPlanJSON
		actual := KubernetesPlan{}
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
		State:   KuberetesClusterStateReady,
		Storage: "01000000-0000-4000-8000-000160010100",
		Type:    KubernetesClusterTypeStandalone,
		UUID:    "0ddab8f4-97c0-4222-91ba-85a4fff7499b",
		Zone:    "de-fra1",
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
	}
}
