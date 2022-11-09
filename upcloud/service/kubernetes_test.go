package service

import (
	"context"
	"fmt"
	"net/http"
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
	const clusterName = "go-sdk-test-ctx"
	const SSHKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIO3fnjc8UrsYDNU8365mL3lnOPQJg18V42Lt8U/8Sm+r testy_test"
	const networkCIDR = "176.16.1.0/24"

	// set when creating a private network for cluster
	network := ""
	// set when creating a cluster
	uuid := ""

	t.Cleanup(func() {
		recordWithContext(t, "delete_kubernetes_cluster", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			if len(uuid) > 0 {
				err := svc.DeleteKubernetesCluster(&request.DeleteKubernetesClusterRequest{UUID: uuid})
				require.NoError(t, err)

				err = waitForKubernetesClusterNotFound(rec, svc, uuid)
				require.NoError(t, err)
			}
		})
		recordWithContext(t, "delete_kubernetes_private_network", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			if len(network) > 0 {
				err := svc.DeleteNetwork(&request.DeleteNetworkRequest{UUID: network})
				require.NoError(t, err)
			}
		})
	})

	// this group is not to be run in parallel
	t.Run("Setup", func(t *testing.T) {
		recordWithContext(t, "create_kubernetes_private_network", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			t.Run("CreateKubernetesPrivateNetwork", func(t *testing.T) {
				n, err := svc.CreateNetwork(&request.CreateNetworkRequest{
					Name: "upcloud-go-sdk-test",
					Zone: zone,
					IPNetworks: []upcloud.IPNetwork{
						{
							Address: networkCIDR,
							DHCP:    upcloud.True,
							Family:  upcloud.IPAddressFamilyIPv4,
						},
					},
				})

				require.NoError(t, err)
				require.NotEmpty(t, n.UUID)

				network = n.UUID
			})
		})

		require.NotEmpty(t, network)

		recordWithContext(t, "create_kubernetes_cluster", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			t.Run("CreateKubernetesCluster", func(t *testing.T) {
				c, err := svc.CreateKubernetesCluster(&request.CreateKubernetesClusterRequest{
					Name:    clusterName,
					Network: network,
					NodeGroups: []upcloud.KubernetesNodeGroup{
						{
							Count: 2,
							Name:  "testgroup",
							Plan:  plan,
							Labels: []upcloud.Label{
								{Key: "managedBy", Value: "go-sdk"},
							},
							SSHKeys: []string{SSHKey},
						},
					},
					Zone: zone,
				})

				require.NoError(t, err)
				require.NotEmpty(t, c.UUID)

				uuid = c.UUID
			})
		})

		require.NotEmpty(t, uuid)

		recordWithContext(t, "wait_for_kubernetes_cluster_state", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			require.NotEmpty(t, uuid)

			t.Run("WaitForKubernetesClusterState", func(t *testing.T) {
				err := waitForKubernetesClusterState(rec, svc, uuid, upcloud.KubernetesClusterStateRunning)
				require.NoError(t, err)

				expected := upcloud.KubernetesClusterStateRunning

				c, err := svcContext.GetKubernetesCluster(ctx, &request.GetKubernetesClusterRequest{
					UUID: uuid,
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

		recordWithContext(t, "get_kubernetes_cluster", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			require.NotEmpty(t, uuid)

			expected := &upcloud.KubernetesCluster{
				Name:        clusterName,
				Network:     network,
				NetworkCIDR: networkCIDR,
				NodeGroups: []upcloud.KubernetesNodeGroup{
					{
						Count: 2,
						Name:  "testgroup",
						Plan:  plan,
						Labels: []upcloud.Label{
							{Key: "managedBy", Value: "go-sdk"},
						},
						SSHKeys:     []string{SSHKey},
						KubeletArgs: []upcloud.KubernetesKubeletArg{},
						Storage:     "01000000-0000-4000-8000-000160010100",
					},
				},
				State: upcloud.KubernetesClusterStateRunning,
				UUID:  uuid,
				Zone:  zone,
			}

			actual, err := svc.GetKubernetesCluster(&request.GetKubernetesClusterRequest{
				UUID: uuid,
			})

			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	})

	t.Run("GetKubernetesClusters", func(t *testing.T) {
		t.Parallel()

		recordWithContext(t, "get_kubernetes_clusters", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			require.NotEmpty(t, uuid)

			c, err := svc.GetKubernetesClusters(&request.GetKubernetesClustersRequest{})

			require.NoError(t, err)
			require.Len(t, c, 1)
		})
	})

	t.Run("GetKubernetesKubeconfig", func(t *testing.T) {
		t.Parallel()

		recordWithContext(t, "get_kubernetes_kubeconfig", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			require.NotEmpty(t, uuid)

			k, err := svc.GetKubernetesKubeconfig(&request.GetKubernetesKubeconfigRequest{
				UUID: uuid,
			})

			require.NoError(t, err)
			require.NotZero(t, k)
		})
	})

	t.Run("GetKubernetesVersions", func(t *testing.T) {
		t.Parallel()

		recordWithContext(t, "get_kubernetes_versions", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			v, err := svc.GetKubernetesVersions(&request.GetKubernetesVersionsRequest{})

			require.NoError(t, err)
			require.NotZero(t, v)
			require.NotZero(t, v[0])
		})
	})
}

func waitForKubernetesClusterNotFound(rec *recorder.Recorder, svc *Service, clusterUUID string) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}

	rec.AddPassthrough(func(h *http.Request) bool {
		return true
	})
	defer func() {
		rec.Passthroughs = nil
	}()

	attempts := 0
	sleepDuration := time.Second * 5

	for {
		attempts++

		_, err := svc.GetKubernetesCluster(&request.GetKubernetesClusterRequest{UUID: clusterUUID})
		if upcloudErr, ok := err.(*upcloud.Error); ok {
			if upcloudErr.ErrorCode == "NotFound" {
				break
			}
		}

		if err != nil {
			return err
		}
		time.Sleep(sleepDuration)

		if time.Duration(attempts)*sleepDuration >= waitTimeout {
			return fmt.Errorf("timeout %s reached", waitTimeout.String())
		}
	}
	// additional wait period to make sure the attached network is deletable
	time.Sleep(waitTimeout)

	return nil
}

func waitForKubernetesClusterState(rec *recorder.Recorder, svc *Service, clusterUUID string, desiredState upcloud.KubernetesClusterState) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}

	rec.AddPassthrough(func(h *http.Request) bool {
		return true
	})
	defer func() {
		rec.Passthroughs = nil
	}()

	_, err := svc.WaitForKubernetesClusterState(&request.WaitForKubernetesClusterStateRequest{
		UUID:         clusterUUID,
		DesiredState: desiredState,
		Timeout:      waitTimeout,
	})

	return err
}
