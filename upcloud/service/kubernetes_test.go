package service

import (
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/require"
)

func TestKubernetes(t *testing.T) {
	t.Parallel()

	const zone = "de-fra1"
	const plan = "K8S-2xCPU-4GB"

	// set when creating a private network for cluster
	network := ""
	// set when creating a cluster
	uuid := ""

	t.Cleanup(func() {
		record(t, "deletekubernetescluster", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			if len(uuid) > 0 {
				err := svc.DeleteKubernetesCluster(&request.DeleteKubernetesClusterRequest{UUID: uuid})

				require.NoError(t, err)
			}
		})
		record(t, "deletekubernetesprivatenetwork", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			if len(network) > 0 {
				err := svc.DeleteNetwork(&request.DeleteNetworkRequest{UUID: network})

				require.NoError(t, err)
			}
		})
	})

	// this group is not to be run in parallel
	t.Run("Setup", func(t *testing.T) {
		record(t, "createkubernetesprivatenetwork", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			t.Run("CreateKubernetesPrivateNetwork", func(t *testing.T) {
				n, err := svc.CreateNetwork(&request.CreateNetworkRequest{
					Name: "upcloud-go-sdk-integration-test",
					Zone: zone,
					IPNetworks: []upcloud.IPNetwork{
						{
							Address:          "172.16.0.0/24",
							DHCP:             upcloud.True,
							DHCPDefaultRoute: upcloud.False,
							DHCPDns: []string{
								"172.16.0.1",
								"172.16.0.2",
							},
							Family:  upcloud.IPAddressFamilyIPv4,
							Gateway: "172.16.0.1",
						},
					},
				})

				require.NoError(t, err)
				require.NotEmpty(t, n.UUID)

				network = n.UUID
			})
		})

		require.NotEmpty(t, network)

		record(t, "createkubernetescluster", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			t.Run("CreateKubernetesCluster", func(t *testing.T) {
				c, err := svc.CreateKubernetesCluster(exampleCreateKubernetesClusterRequest(network, plan, zone))

				require.NoError(t, err)
				require.NotEmpty(t, c.UUID)

				uuid = c.UUID
			})
		})

		require.NotEmpty(t, uuid)

		record(t, "waitforkubernetesclusterstate", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			t.Run("WaitForKubernetesClusterState", func(t *testing.T) {
				require.NotEmpty(t, uuid)

				expected := upcloud.KuberetesClusterStateReady

				c, err := svc.WaitForKubernetesClusterState(&request.WaitForKubernetesClusterStateRequest{
					DesiredState: upcloud.KuberetesClusterStateReady,
					Timeout:      time.Minute * 15,
					UUID:         uuid,
				})
				require.NotNil(t, c)
				actual := c.State

				require.NoError(t, err)
				require.Equal(t, expected, actual)
			})
		})
	})

	t.Run("GetKubernetesCluster", func(t *testing.T) {
		t.Parallel()

		record(t, "getkubernetescluster", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			require.NotEmpty(t, uuid)

			expected := exampleKubernetesCluster(network, plan, uuid, zone)

			actual, err := svc.GetKubernetesCluster(&request.GetKubernetesClusterRequest{
				UUID: uuid,
			})

			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	})

	t.Run("GetKubernetesClusters", func(t *testing.T) {
		require.NotEmpty(t, uuid)

		t.Parallel()

		record(t, "getkubernetesclusters", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			c, err := svc.GetKubernetesClusters(&request.GetKubernetesClustersRequest{})

			require.NoError(t, err)
			require.Len(t, c, 1)
		})
	})

	t.Run("GetKubernetesKubeconfig", func(t *testing.T) {
		require.NotEmpty(t, uuid)

		t.Parallel()

		record(t, "getkuberneteskubeconfig", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			k, err := svc.GetKubernetesKubeconfig(&request.GetKubernetesKubeconfigRequest{
				UUID: uuid,
			})

			require.NoError(t, err)
			require.NotZero(t, k)
		})
	})

	t.Run("GetKubernetesPlans", func(t *testing.T) {
		t.Parallel()

		record(t, "getkubernetesplans", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			p, err := svc.GetKubernetesPlans(&request.GetKubernetesPlansRequest{})

			require.NoError(t, err)
			require.NotZero(t, p)
		})
	})

	t.Run("GetKubernetesVersions", func(t *testing.T) {
		t.Parallel()

		record(t, "getkubernetesversions", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
			v, err := svc.GetKubernetesVersions(&request.GetKubernetesVersionsRequest{})

			require.NoError(t, err)
			require.NotZero(t, v)
		})
	})
}

func exampleKubernetesCluster(network, plan, uuid, zone string) *upcloud.KubernetesCluster {
	return &upcloud.KubernetesCluster{
		Name:    "upcloud-go-sdk-integration-test",
		Network: network,
		NodeGroups: []upcloud.KubernetesNodeGroup{
			exampleKubernetesNodeGroup(plan),
			exampleKubernetesNodeGroup(plan),
		},
		State: upcloud.KuberetesClusterStateReady,
		Type:  upcloud.KubernetesClusterTypeStandalone,
		UUID:  uuid,
		Zone:  zone,
	}
}

func exampleKubernetesNodeGroup(plan string) upcloud.KubernetesNodeGroup {
	return upcloud.KubernetesNodeGroup{
		Count: 1,
		Labels: []upcloud.Label{
			{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-integration-test",
			},
		},
		Name: "upcloud-go-sdk-integration-test",
		Plan: plan,
	}
}

func exampleCreateKubernetesClusterRequest(network, plan, zone string) *request.CreateKubernetesClusterRequest {
	return &request.CreateKubernetesClusterRequest{
		Name:    "upcloud-go-sdk-integration-test",
		Network: network,
		NodeGroups: []upcloud.KubernetesNodeGroup{
			exampleKubernetesNodeGroup(plan),
			exampleKubernetesNodeGroup(plan),
		},
		Storage: "01000000-0000-4000-8000-000160010100",
		Zone:    zone,
	}
}
