package request

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
)

const exampleCreateKubernetesClusterRequestJSON string = `{
	"network": "00000000-0000-0000-0000-000000000000",
	"node_groups": [
		{
			"count": 1,
			"labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			],
			"name": "name",
			"plan": "plan",
			"ssh_key": [
				"key",
				"key"
			]
		},
		{
			"count": 1,
			"labels": [
				{
					"key": "managedBy",
					"value": "upcloud-go-sdk-unit-test"
				}
			],
			"name": "name",
			"plan": "plan",
			"ssh_key": [
				"key",
				"key"
			]
		}
	],
	"title": "title",
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
			"labels=managedBy",
			"labels=managedBy=upcloud-go-sdk-unit-test",
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

	t.Run("GetKubernetesPlansRequest", func(t *testing.T) {
		t.Parallel()

		expected := fmt.Sprintf(
			"%s/plan",
			kubernetesClusterBasePath,
		)

		r := GetKubernetesPlansRequest{}
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
			"%s/version",
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

func exampleGetKubernetesClustersWithFiltersRequest() GetKubernetesClustersWithFiltersRequest {
	return GetKubernetesClustersWithFiltersRequest{
		Filters: []KubernetesFilter{
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
		Name:    "title",
		Network: "00000000-0000-0000-0000-000000000000",
		NodeGroups: []upcloud.KubernetesNodeGroup{
			exampleKubernetesNodeGroup(),
			exampleKubernetesNodeGroup(),
		},
		Zone: "zone",
	}
}

func exampleKubernetesNodeGroup() upcloud.KubernetesNodeGroup {
	return upcloud.KubernetesNodeGroup{
		Count: 1,
		Labels: []upcloud.Label{
			{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
		Name: "name",
		Plan: "plan",
		SSHKeys: []string{
			"key",
			"key",
		},
	}
}
