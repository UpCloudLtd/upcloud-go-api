package service

import (
	"context"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/require"
)

func TestKubernetesCtx(t *testing.T) {
	t.Parallel()

	const zone = "de-fra1"
	const plan = "K8S-2xCPU-4GB"

	// set when creating a private network for cluster
	network := ""
	// set when creating a cluster
	uuid := ""

	t.Cleanup(func() {
		recordWithContext(t, "deletekubernetescluster", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			if len(uuid) > 0 {
				err := svcContext.DeleteKubernetesCluster(ctx, &request.DeleteKubernetesClusterRequest{UUID: uuid})

				require.NoError(t, err)
			}
		})
		recordWithContext(t, "deletekubernetesprivatenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			if len(network) > 0 {
				err := svcContext.DeleteNetwork(ctx, &request.DeleteNetworkRequest{UUID: network})

				require.NoError(t, err)
			}
		})
	})

	// this group is not to be run in parallel
	t.Run("Setup", func(t *testing.T) {
		recordWithContext(t, "createkubernetesprivatenetwork", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			t.Run("CreateKubernetesPrivateNetwork", func(t *testing.T) {
				n, err := svcContext.CreateNetwork(ctx, &request.CreateNetworkRequest{
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

		recordWithContext(t, "createkubernetescluster", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			t.Run("CreateKubernetesCluster", func(t *testing.T) {
				c, err := svcContext.CreateKubernetesCluster(ctx, exampleCreateKubernetesClusterRequest(network, plan, zone))

				require.NoError(t, err)
				require.NotEmpty(t, c.UUID)

				uuid = c.UUID
			})
		})

		require.NotEmpty(t, uuid)

		recordWithContext(t, "waitforkubernetesclusterstate", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			require.NotEmpty(t, uuid)

			t.Run("WaitForKubernetesClusterState", func(t *testing.T) {
				expected := upcloud.KuberetesClusterStateReady

				c, err := svcContext.WaitForKubernetesClusterState(ctx, &request.WaitForKubernetesClusterStateRequest{
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

		recordWithContext(t, "getkubernetescluster", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			require.NotEmpty(t, uuid)

			expected := exampleKubernetesCluster(network, plan, uuid, zone)

			actual, err := svcContext.GetKubernetesCluster(ctx, &request.GetKubernetesClusterRequest{
				UUID: uuid,
			})

			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	})

	t.Run("GetKubernetesClusters", func(t *testing.T) {
		t.Parallel()

		recordWithContext(t, "getkubernetesclusters", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			require.NotEmpty(t, uuid)

			c, err := svcContext.GetKubernetesClusters(ctx, &request.GetKubernetesClustersRequest{})

			require.NoError(t, err)
			require.Len(t, c, 1)
		})
	})

	t.Run("GetKubernetesKubeconfig", func(t *testing.T) {
		t.Parallel()

		recordWithContext(t, "getkuberneteskubeconfig", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			require.NotEmpty(t, uuid)

			k, err := svcContext.GetKubernetesKubeconfig(ctx, &request.GetKubernetesKubeconfigRequest{
				UUID: uuid,
			})

			require.NoError(t, err)
			require.NotZero(t, k)
		})
	})

	t.Run("GetKubernetesPlans", func(t *testing.T) {
		t.Parallel()

		recordWithContext(t, "getkubernetesplans", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			p, err := svcContext.GetKubernetesPlans(ctx, &request.GetKubernetesPlansRequest{})

			require.NoError(t, err)
			require.NotZero(t, p)
		})
	})

	t.Run("GetKubernetesVersions", func(t *testing.T) {
		t.Parallel()

		recordWithContext(t, "getkubernetesversions", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
			v, err := svcContext.GetKubernetesVersions(ctx, &request.GetKubernetesVersionsRequest{})

			require.NoError(t, err)
			require.NotZero(t, v)
		})
	})
}
