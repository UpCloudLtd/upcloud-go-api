package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBalancer(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancer", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create Load Balancer
		lb, err := createLoadBalancerAndNetwork(ctx, svc, "fi-hel1", "172.16.1.0/24")
		require.NoError(t, err)

		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()

		t.Logf("Created load balancer: %s", lb.Name)
		assert.Equal(t, upcloud.LoadBalancerMaintenanceDOWSunday, lb.MaintenanceDOW)
		assert.Equal(t, "20:01:01Z", lb.MaintenanceTime)

		// Modify Load Balancer
		t.Log("Modifying load balancer")

		newName := "new-name-for-lb"
		lb, err = svc.ModifyLoadBalancer(ctx, &request.ModifyLoadBalancerRequest{
			UUID:            lb.UUID,
			Name:            newName,
			MaintenanceDOW:  upcloud.LoadBalancerMaintenanceDOWMonday,
			MaintenanceTime: "00:01:01Z",
		})
		require.NoError(t, err)
		assert.Equal(t, newName, lb.Name)
		assert.Equal(t, upcloud.LoadBalancerMaintenanceDOWMonday, lb.MaintenanceDOW)
		assert.Equal(t, "00:01:01Z", lb.MaintenanceTime)
		t.Logf("Modified load balancer with UUID: %s", lb.UUID)

		// Get Load Balancer
		t.Logf("Get load balancer: %s", newName)
		lb, err = svc.GetLoadBalancer(ctx, &request.GetLoadBalancerRequest{
			UUID: lb.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, newName, lb.Name)

		// Get Load Balancers
		t.Log("Get load balancers")
		lbs, err := svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(lbs), 1)
	})
}

func TestLoadBalancerBackend(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancerbackend", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancerAndNetwork(ctx, svc, "fi-hel2", "172.16.2.0/24")
		require.NoError(t, err)
		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()

		t.Logf("Created load balancer for testing LB backend CRUD: %s", lb.Name)

		backend, err := createLoadBalancerBackend(ctx, svc, lb.UUID)
		require.NoError(t, err)
		assert.Equal(t, 30, backend.Properties.TimeoutServer)
		assert.Equal(t, upcloud.BoolPtr(true), backend.Properties.TLSEnabled)
		assert.Equal(t, upcloud.BoolPtr(true), backend.Properties.TLSVerify)
		assert.Equal(t, upcloud.BoolPtr(true), backend.Properties.TLSUseSystemCA)
		assert.Equal(t, upcloud.BoolPtr(false), backend.Properties.HTTP2Enabled)
		t.Logf("Created LB backend: %s", backend.Name)

		t.Logf("Modifying LB backend: %s", backend.Name)
		newName := "updatedName"
		backend, err = svc.ModifyLoadBalancerBackend(ctx, &request.ModifyLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        backend.Name,
			Backend: request.ModifyLoadBalancerBackend{
				Name: newName,
				Properties: &upcloud.LoadBalancerBackendProperties{
					HealthCheckType: upcloud.LoadBalancerHealthCheckTypeTCP,
				},
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, backend.Name, newName)
		assert.Equal(t, 30, backend.Properties.TimeoutServer)
		assert.Equal(t, upcloud.BoolPtr(true), backend.Properties.TLSEnabled)
		assert.Equal(t, upcloud.BoolPtr(true), backend.Properties.TLSVerify)
		assert.Equal(t, upcloud.BoolPtr(true), backend.Properties.TLSUseSystemCA)
		assert.Equal(t, upcloud.BoolPtr(false), backend.Properties.HTTP2Enabled)
		assert.Equal(t, upcloud.LoadBalancerHealthCheckType("tcp"), backend.Properties.HealthCheckType)
		t.Logf("Modified LB backend, new name is: %s", backend.Name)

		t.Logf("Get LB backend: %s", newName)
		backend, err = svc.GetLoadBalancerBackend(ctx, &request.GetLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        newName,
		})
		require.NoError(t, err)
		assert.EqualValues(t, backend.Name, newName)
		assert.Len(t, backend.Members, 1)

		t.Log("Get LB backend list")
		backends, err := svc.GetLoadBalancerBackends(ctx, &request.GetLoadBalancerBackendsRequest{
			ServiceUUID: lb.UUID,
		})
		require.NoError(t, err)
		assert.Len(t, backends, 1)
		assert.EqualValues(t, backends[0].Name, newName)
		assert.Len(t, backends[0].Members, 1)

		t.Logf("Deleting LB backend: %s", backend.Name)
		err = svc.DeleteLoadBalancerBackend(ctx, &request.DeleteLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        backend.Name,
		})
		require.NoError(t, err)
	})
}

// TestLoadBalancerCerticateBundlesAndBackendTLSConfigs tests certificate bundles and TLS configs
// Test:
//   - certificate bundle CRUD
//   - TLS config CRUD
//   - add TLS config to LB backend
//   - enable HTTP/2 for LB backend
//   - remove TLS config to LB backend
func TestLoadBalancerCerticateBundlesAndBackendTLSConfigs(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancercerticatebundlesandbackendtlsconfigs", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		net, err := createLoadBalancerAndPrivateNetwork(ctx, svc, "fi-hel1", "10.0.5.0/24")
		require.NoError(t, err)
		feName := "fe-1"
		beName := "be-1"
		lb, err := svc.CreateLoadBalancer(ctx, &request.CreateLoadBalancerRequest{
			Name:             fmt.Sprintf("go-api-lb-%s-%d", "fi-hel1", time.Now().Unix()),
			Zone:             "fi-hel1",
			Plan:             "development",
			NetworkUUID:      net.UUID,
			ConfiguredStatus: upcloud.LoadBalancerConfiguredStatusStarted,
			Frontends: []request.LoadBalancerFrontend{{
				Name:           feName,
				Mode:           upcloud.LoadBalancerModeHTTP,
				Port:           80,
				DefaultBackend: beName,
			}},
			Backends: []request.LoadBalancerBackend{{
				Name:     beName,
				Resolver: "ns-1",
				Members:  []request.LoadBalancerBackendMember{},
			}},
			Resolvers: []request.LoadBalancerResolver{{
				Name:         "ns-1",
				Nameservers:  []string{"10.1.1.100"},
				Retries:      10,
				Timeout:      10,
				TimeoutRetry: 10,
				CacheValid:   10,
				CacheInvalid: 10,
			}},
		})

		require.NoError(t, err)
		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()

		mc, err := svc.CreateLoadBalancerCertificateBundle(ctx, &request.CreateLoadBalancerCertificateBundleRequest{
			Type:        upcloud.LoadBalancerCertificateBundleTypeManual,
			Name:        "go-api-be-manual-certificate",
			Certificate: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZhekNDQTFPZ0F3SUJBZ0lVRzR1KzRDZmlHQ3pQSDk4dDA4QXh5VkE0QzVnd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdNQ2xOdmJXVXRVM1JoZEdVeElUQWZCZ05WQkFvTQpHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpEQWVGdzB5TWpBek1UY3hNekUxTURoYUZ3MHlNekF6Ck1UY3hNekUxTURoYU1FVXhDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWxMVk4wWVhSbE1TRXcKSHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXdnZ0lpTUEwR0NTcUdTSWIzRFFFQgpBUVVBQTRJQ0R3QXdnZ0lLQW9JQ0FRRHZaeG4vK3pUeVc0RVJ2S2t3V29tVXBpOG8ydEp6MWR2ZXIrREpySzNnCkNObFVvWXpSMjlDV3M3aks4MVhNc3ZtcUw1TXpUd1A3SHNtZDFxNjlGSStXY1BFMWFhYjk5MDlJQWsvR0dpSzIKelRsZU4zRVFRcFhuN3RueVB0WmFUOFkxM3lGSHBDNVJnUXpURThDUjlaaTJPTEV5eEdRMzZwQTYxOTBueFZnMgpTTGxhZk5HVFp0SnZOMS83cjltSmhFbGJyVUUram9lWEx3Tm9qSC9uWGs1Vy9Yd3paYm9JSHNTRlZZaksyemxnCm9xQzYrQXBvOXhGOW9ZN25sQWhRMEtLV3ZRVmJ3akdQbVZOMTdFVG9kSHNLSlpCb1h4RHNaVVRHQ0RESkNpbXoKVzY0YTc5bFdJeGl5T1E0LzdUbjJGaFBZMG9tSDVVYldDUHEyTW5YZWJrT2pnY3ZVWTROSXd6cFlWMFcyR0dHRwp3d25pOWZsbFlBTTlPRDNidlBYNU9hQVdOSlQ2cjZFYXNzaldsdjBUZUd4RStCWlorZzN4UHFIVEd6MndIekM1CjVhbkxEak0rNHZzQlZrZmtWM1NZN1c4M203NFZRK1FhM1dhTlh6aW5MMGtlRnh1cExYWThiS2hFelh6U0xLeisKQnI4UEdlR1JnYVNEZDFrcEZQZyt1ak44cXZnbzBSREk4SXFMUzd6YlhGb1FycDF4L2RXbTlTOERWRVhWb1VBMQpXUW5WdVdFQ29CUzRaZjQxZDA0cGZkQ3R0bk45ekhvc2d3WGJKOG0wVGZ2Zmt1aFZpdVZBTi9wK01wOVduUStICjExSEVuV3BTZk9oN1pQalR6anVBc2V2VmZWNGc0YTNrY3pNdjFycE5QelVVUHd0QXF4OTIzOXd3SVI5WTE1Y2wKT1FJREFRQUJvMU13VVRBZEJnTlZIUTRFRmdRVVhJcWxiajV1TVVGVC9qcU0ya1d2WVp0RE5rY3dId1lEVlIwagpCQmd3Rm9BVVhJcWxiajV1TVVGVC9qcU0ya1d2WVp0RE5rY3dEd1lEVlIwVEFRSC9CQVV3QXdFQi96QU5CZ2txCmhraUc5dzBCQVFzRkFBT0NBZ0VBQ00reGJiOW9rVUdPcWtRWmtndHh3eVF0em90WFRXNnlIYmdmeGhCd3d1disKc0ltT1QyOUw5WFVZblM5cTgrQktBSHVGa3JTRUlwR0ZVVWhDVWZTa2xlNmZuR0c1b05PWm1DdEszN3RoZlFVOQp2NEJvOHFCWkhqREh3azlWVHRhWm1BazBLYnhmaHVneVdWQ1ZsbURacm9TQ09pV0drVFZoc1hhS0RrYnc0RWwwCjJzY3lnYkFDdFZ4bkU1WjlmU0F3MU9QWXJZYUcySW5HTDQvMHVSZXo4aXl1UE9lNUNiL0RkNDl1eHFzR1FkM1IKQzdKNC9vWnB2b0V6UVJtakxib1FzQzkwU2ZqaFNpcGhHQlNiYUpCZGRsMDBrNVZzVXJxS1haU004cVFxVWZWLwpubEJtYjJOblVsa2RlOEtIczBQamhCaG8rLzdmaitMN21GYTJsNWpmdWlsdHVxdmgyWnladFJjd2didmJlaUxPCmZQSWlMQ2dTbnMwaitZMkVrS1drRUp6RXJQVm5sOTdaQktZclBaYmRYMFY5b2dvTC9qeEV5NzlsbzlKczI5djYKUkY2NmdvSlUwMkVKZTUwMmk3WHJzMzFZQ0tuSGd2ejUwTDZha0JpYWRSNmtrTXVXdkJ1d1l6MElaS1RMcXhqZAowOEdlUkJVeWFsUFZodGZKbzNNdXRuYUllL1pWVTdLQUl3S1Znb20zS09EY1RpWllQV3RWKzFnL0UvN3A1aGh2CkJERzFqcklRc1ZrZG4yNWZhNXNkNU9Qa1AvbDBRdXY1em16UEk3S1MrS2ZlWS92NHFBOTBtNGk2dkZORlRtbTAKSFNXV0JZTlR4blIxYjk2UElUcnRzOE15am9YTFg2QnUxVkZOSlByMkpnMDJMVlZvcTZSSWJlMVVvNjE5b2pBPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
			PrivateKey:  "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUpSUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQ1M4d2dna3JBZ0VBQW9JQ0FRRHZaeG4vK3pUeVc0RVIKdktrd1dvbVVwaThvMnRKejFkdmVyK0RKckszZ0NObFVvWXpSMjlDV3M3aks4MVhNc3ZtcUw1TXpUd1A3SHNtZAoxcTY5RkkrV2NQRTFhYWI5OTA5SUFrL0dHaUsyelRsZU4zRVFRcFhuN3RueVB0WmFUOFkxM3lGSHBDNVJnUXpUCkU4Q1I5WmkyT0xFeXhHUTM2cEE2MTkwbnhWZzJTTGxhZk5HVFp0SnZOMS83cjltSmhFbGJyVUUram9lWEx3Tm8KakgvblhrNVcvWHd6WmJvSUhzU0ZWWWpLMnpsZ29xQzYrQXBvOXhGOW9ZN25sQWhRMEtLV3ZRVmJ3akdQbVZOMQo3RVRvZEhzS0paQm9YeERzWlVUR0NEREpDaW16VzY0YTc5bFdJeGl5T1E0LzdUbjJGaFBZMG9tSDVVYldDUHEyCk1uWGVia09qZ2N2VVk0Tkl3enBZVjBXMkdHR0d3d25pOWZsbFlBTTlPRDNidlBYNU9hQVdOSlQ2cjZFYXNzalcKbHYwVGVHeEUrQlpaK2czeFBxSFRHejJ3SHpDNTVhbkxEak0rNHZzQlZrZmtWM1NZN1c4M203NFZRK1FhM1dhTgpYemluTDBrZUZ4dXBMWFk4YktoRXpYelNMS3orQnI4UEdlR1JnYVNEZDFrcEZQZyt1ak44cXZnbzBSREk4SXFMClM3emJYRm9RcnAxeC9kV205UzhEVkVYVm9VQTFXUW5WdVdFQ29CUzRaZjQxZDA0cGZkQ3R0bk45ekhvc2d3WGIKSjhtMFRmdmZrdWhWaXVWQU4vcCtNcDlXblErSDExSEVuV3BTZk9oN1pQalR6anVBc2V2VmZWNGc0YTNrY3pNdgoxcnBOUHpVVVB3dEFxeDkyMzl3d0lSOVkxNWNsT1FJREFRQUJBb0lDQVFDcUNtd2dNbmcvOEJoejFiSENRM3hYCkZkYUhTUzJUMHdHUllSRGpqZ0FPRVpyMERxN3IzQnFDLy9Jd1RMZlRaZ2dKQmpPaWpPd0I4TE01cGVPRkwxWngKZjVVRDRDQVpZUkJ4MEJxRFZjcjBWajM2R3B6MjlLUnZFV3JDTWptaitlZUtHZ3NVVEp3TmpnRGk1N091dUdlWQpmaG4yT2lJSXlWVmFSanF4NWV5cTJlcTFSOVMvd3BlVElSek9zdTlyU29la1V5SDFZZDBTMS9TdXpLU0lYS1orCkNSdXZrZ0NaaGVrRjMyUUMyY1VlUzBTb3FFY1VtUEJXY0dzRk4xTFV1K3ZQNzBBZ0ZZV0lQbHBXZHRQVzIrME0KbnZPNy9sSVI1amY4QkpOS0tDcklWMFVKb3ZTV3h1VGlxYjNpVUFnTUwxQTNnQXJwZUVOaEFRMjZYWXIweXhMRQpiMTRObWdnZzRUSktRMWp4aTZlRys5Nkp0WStteFdwaHJUNTVUT2s5blpEY2JTV29ibTgwcExJZEE0QTN0bjhJCjI5SXZ2dkdhVnh2NXFXaU5sL0sxWEZhK2JRRVY3K1AxT0RWSnV5VXEzdFg4dDlONVJ2c3RiUE5XNU8xQ01STGsKZExESVMwRFFKYytMa1ZGaGpZRXlJOWluQmtlRXV1ekI0L0k3M1JnTXpKMFZvZU1hSXEwdWpxbVF6KzN0ei9JUgp6VTFSN2FndEZmNHUvTS9jeVl2U0E4UEZVSis1Q25ldHFtWC91dzR0em10WmtLd3JnMVArblVpRTV4Z0ZLT1FZCjUyaU81aXRKWHo4aHJGdFpmbVMrcTI3eFRFWWlEdTFFSFg0Y1pLaFBQNVh3RzkyekN0ZUxJenc5QUh2TS9aRTcKNmI0OFMwQWR6T3dFN3VaVVhheEw5UUtDQVFFQTkrOVROV2RuMjVxK3lSWUlsc2x1NFI1N0xSQm1ucXFhaDljSQowZ1JGS0RZUFZTWFl0RlhVci9mK0psSmR4YUxEYjYwV0o4c2hWcHZ1MTlJTE0wNTR6M3ZYa0w0QjQwSWI1Z3JnCnlzbVZKVlpoYVdTNWU5RW02TEh2bXRoVWl0MzVDWnlLcTlkN2FwdzFOQ1VIUkFEWE9wSndBWkc4ckttRmtDeFYKTnFpVThnWi9LOXpPK0gra215bGVKZTkrZHhCR2RCZTczbnFsRTVPQVdJcjVLcFNpMU5sQWYxOG0yYnVUWkxiTQo2Rnd1MEc0SjUzd1ZaTFRhbWxVelNmdjFTRGRrMEZIbnpkdEJPMWI3OTNJeVpOQXc4eGlkMGVPVy90OE1QTm1RCmFXWnNqaWc4WEZVMlFaMXl5YUd6amoxNUpoVzl2WGN1SUpGc21rRzZiRS8wSm9ZZnd3S0NBUUVBOXpDNXNYN1IKQlN2MzlwVnFROXI4bDJXNVVqQk9iczFick05Q0NPQWFFZCtlY2xyZWZZRU5yMUVtL0x3Qk9GK1RuRnVnMjJZOQpKczliNjgrY0wzSEtweVJ5KzVVaGM1NzU0b3pDRklDemJzZHphM005TGVzVUtTVTBYYnlhNm5vc0UvWjNOWGpyCmpLQ1ROZEU2eFY4OTZ1Vm5FTXZpcDh4M1kvR0VBejh4U3lCOUFveGNtanVqZnVVdmlZbmN5SnBqUUoxOEhsZk8KMlloWWdTd3VxSFhISUswUXhGckJiTjIzZXZ5TUkrOVVSSVlnOWxtemFWR2hqczlHd0Y3UUQ0dFoyOXgwQVFOSwpBUFArMnI4ZUlLbjFucEs4VnFRdVNrNVdsZjg3L0xvV1Zha2xtRXkvdEJ0alpjNVBQanpMVmRTLy9QNWlkc2F1CkZuTnc2VDNmZWEwelV3S0NBUUVBMWsySDU2WXNzRFhPYUxOaDB5dmphalJWbGJzU2FGemdXei8wQU12dUZ2YTcKUkFjRmk4S1FwMVU4MlZUaWRyemNIc0JHWVRrRDVQKzliOUMvRzZiZFo4SU1yckI5b3poMk10MytOV29PUDRxdApnbEtzdktnbzhJTTBydXdFRDFBVVBVbVExejNYRUd4YTFHcVpJQjkxNmN1L2dxdThvS1dhcSthVjlUdThHb0toCkU0RzFhRGUwU09WMTJtWnJNbkRmNU9MSzRWK3pKZnVkdVdyT09nN2x2QUxZNi8rTDdqRmpFbStySjhEZU9neVQKQlFKTTM1SXZUYTBOT3dyTWxaSkQwb2lwUzFjVHlEM0VacnJQY2pJOXpUSGU0QmZQWVJmY1ZSQmM4YTIxY1I2NApKYnNGdmF0aEY0VnNWU3N2ZDByZGlWSGxqZ01GRTBSeTVjSXFMODVJendLQ0FRRUEzdE1nZ1QwRkpIbFhKQVBhCmIrS0drZDlUNkIrOWhDcEFPbzMyUTlQb0REYWRPUTVxdzQzRERVZkZNa3d6ZVdMR3lFcmN2UW56ay9tV0xnTFAKRXdHcm9YRzg2TWF0Q2ZIRDVoSG1uZDdLWU5FUVhVcmJXbm92aVV1TllmWXpXNnpYOFFMYXdPd0l3Wks2UU9nago1Mm1NZ2lOYS9nd2NmQkJYaTFOYUlpY2p3MG85QmtBSzljbE8vNE9QajVjajIvMDMvVFk1Zll5LzNONElraUNHCnlycW96dTdUVDMxVUlWUFlJdGhuWjdsRktDUVVzSjE1bWpYSXdkaGRPZW45K2hVdTRuOWVYczlkTlhDOVN1aS8KT3NpYXJlQXVRSmZ0Vm5RNW55c2VJeHFJS1oyNVV3blVRWUh5M3dIVDh4R1FaZ1hMTHo4TStXN3QzVFVoRWxBQgpGRWtxR3dLQ0FRRUFtT2VzMXFINlA0dDRsZ0VjK01Ubzc2VmJvYnhab29HQWdMWmc3N1J4c2hPaGxMdCtYZnIzCjFsOWgycFJ0eXFIdDRuMG9ob1A0VzNSN3VuQ05rNWl3SGNJSDVjNmhGQTYyelVEM1JjeEhJZERhdDBGL3RoWDgKTUpndDlsM0MrMVJwZFZlV0hlbEJOM0JlM0FtWEpXL1ZsL3lGTXVjWWxETnlXUFlPSmRuQ1BRZ0FGVFJnVnJlUQpiUjZCY29neUVRTVEzenRMWnNBRnRaZ25Sak05YkpLN1JjYzg5bGxaa1BuMHVKZjNKVWxMeVFFN0l2bEJsWi9tClZnUUhiRTkwQStZNzFpb1piQWh5TFcwTE9lTmhBS3NRNFJZbDlnb0N4dGp1ZnE2NnFTNGNzdGN6c2J5N083dFAKeXZkSXp2eEZRZmx4Yk8ra1ludDNkcFRIdUNuUkNIMFM0UT09Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K",
		})
		require.NoError(t, err)
		t.Logf("Created manual certificate bundle %s", mc.Name)
		assert.Equal(t, "go-api-be-manual-certificate", mc.Name)
		assert.Equal(t, upcloud.LoadBalancerCertificateBundleTypeManual, mc.Type)

		ac, err := svc.CreateLoadBalancerCertificateBundle(ctx, &request.CreateLoadBalancerCertificateBundleRequest{
			Type:        upcloud.LoadBalancerCertificateBundleTypeAuthority,
			Name:        "go-api-be-authority-certificate",
			Certificate: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZhekNDQTFPZ0F3SUJBZ0lVRzR1KzRDZmlHQ3pQSDk4dDA4QXh5VkE0QzVnd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdNQ2xOdmJXVXRVM1JoZEdVeElUQWZCZ05WQkFvTQpHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpEQWVGdzB5TWpBek1UY3hNekUxTURoYUZ3MHlNekF6Ck1UY3hNekUxTURoYU1FVXhDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWxMVk4wWVhSbE1TRXcKSHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXdnZ0lpTUEwR0NTcUdTSWIzRFFFQgpBUVVBQTRJQ0R3QXdnZ0lLQW9JQ0FRRHZaeG4vK3pUeVc0RVJ2S2t3V29tVXBpOG8ydEp6MWR2ZXIrREpySzNnCkNObFVvWXpSMjlDV3M3aks4MVhNc3ZtcUw1TXpUd1A3SHNtZDFxNjlGSStXY1BFMWFhYjk5MDlJQWsvR0dpSzIKelRsZU4zRVFRcFhuN3RueVB0WmFUOFkxM3lGSHBDNVJnUXpURThDUjlaaTJPTEV5eEdRMzZwQTYxOTBueFZnMgpTTGxhZk5HVFp0SnZOMS83cjltSmhFbGJyVUUram9lWEx3Tm9qSC9uWGs1Vy9Yd3paYm9JSHNTRlZZaksyemxnCm9xQzYrQXBvOXhGOW9ZN25sQWhRMEtLV3ZRVmJ3akdQbVZOMTdFVG9kSHNLSlpCb1h4RHNaVVRHQ0RESkNpbXoKVzY0YTc5bFdJeGl5T1E0LzdUbjJGaFBZMG9tSDVVYldDUHEyTW5YZWJrT2pnY3ZVWTROSXd6cFlWMFcyR0dHRwp3d25pOWZsbFlBTTlPRDNidlBYNU9hQVdOSlQ2cjZFYXNzaldsdjBUZUd4RStCWlorZzN4UHFIVEd6MndIekM1CjVhbkxEak0rNHZzQlZrZmtWM1NZN1c4M203NFZRK1FhM1dhTlh6aW5MMGtlRnh1cExYWThiS2hFelh6U0xLeisKQnI4UEdlR1JnYVNEZDFrcEZQZyt1ak44cXZnbzBSREk4SXFMUzd6YlhGb1FycDF4L2RXbTlTOERWRVhWb1VBMQpXUW5WdVdFQ29CUzRaZjQxZDA0cGZkQ3R0bk45ekhvc2d3WGJKOG0wVGZ2Zmt1aFZpdVZBTi9wK01wOVduUStICjExSEVuV3BTZk9oN1pQalR6anVBc2V2VmZWNGc0YTNrY3pNdjFycE5QelVVUHd0QXF4OTIzOXd3SVI5WTE1Y2wKT1FJREFRQUJvMU13VVRBZEJnTlZIUTRFRmdRVVhJcWxiajV1TVVGVC9qcU0ya1d2WVp0RE5rY3dId1lEVlIwagpCQmd3Rm9BVVhJcWxiajV1TVVGVC9qcU0ya1d2WVp0RE5rY3dEd1lEVlIwVEFRSC9CQVV3QXdFQi96QU5CZ2txCmhraUc5dzBCQVFzRkFBT0NBZ0VBQ00reGJiOW9rVUdPcWtRWmtndHh3eVF0em90WFRXNnlIYmdmeGhCd3d1disKc0ltT1QyOUw5WFVZblM5cTgrQktBSHVGa3JTRUlwR0ZVVWhDVWZTa2xlNmZuR0c1b05PWm1DdEszN3RoZlFVOQp2NEJvOHFCWkhqREh3azlWVHRhWm1BazBLYnhmaHVneVdWQ1ZsbURacm9TQ09pV0drVFZoc1hhS0RrYnc0RWwwCjJzY3lnYkFDdFZ4bkU1WjlmU0F3MU9QWXJZYUcySW5HTDQvMHVSZXo4aXl1UE9lNUNiL0RkNDl1eHFzR1FkM1IKQzdKNC9vWnB2b0V6UVJtakxib1FzQzkwU2ZqaFNpcGhHQlNiYUpCZGRsMDBrNVZzVXJxS1haU004cVFxVWZWLwpubEJtYjJOblVsa2RlOEtIczBQamhCaG8rLzdmaitMN21GYTJsNWpmdWlsdHVxdmgyWnladFJjd2didmJlaUxPCmZQSWlMQ2dTbnMwaitZMkVrS1drRUp6RXJQVm5sOTdaQktZclBaYmRYMFY5b2dvTC9qeEV5NzlsbzlKczI5djYKUkY2NmdvSlUwMkVKZTUwMmk3WHJzMzFZQ0tuSGd2ejUwTDZha0JpYWRSNmtrTXVXdkJ1d1l6MElaS1RMcXhqZAowOEdlUkJVeWFsUFZodGZKbzNNdXRuYUllL1pWVTdLQUl3S1Znb20zS09EY1RpWllQV3RWKzFnL0UvN3A1aGh2CkJERzFqcklRc1ZrZG4yNWZhNXNkNU9Qa1AvbDBRdXY1em16UEk3S1MrS2ZlWS92NHFBOTBtNGk2dkZORlRtbTAKSFNXV0JZTlR4blIxYjk2UElUcnRzOE15am9YTFg2QnUxVkZOSlByMkpnMDJMVlZvcTZSSWJlMVVvNjE5b2pBPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
		})
		require.NoError(t, err)
		t.Logf("Created authority certificate bundle %s", ac.Name)
		assert.Equal(t, "go-api-be-authority-certificate", ac.Name)
		assert.Equal(t, upcloud.LoadBalancerCertificateBundleTypeAuthority, ac.Type)

		mc, err = svc.ModifyLoadBalancerCertificateBundle(ctx, &request.ModifyLoadBalancerCertificateBundleRequest{
			UUID: mc.UUID,
			Name: "go-api-be-manual-certificate-edit",
		})
		assert.NoError(t, err)
		t.Logf("Modified manual certificate bundle %s", mc.Name)
		assert.Equal(t, "go-api-be-manual-certificate-edit", mc.Name)

		ac, err = svc.ModifyLoadBalancerCertificateBundle(ctx, &request.ModifyLoadBalancerCertificateBundleRequest{
			UUID: ac.UUID,
			Name: "go-api-be-authority-certificate-edit",
		})
		assert.NoError(t, err)
		t.Logf("Modified authority certificate bundle %s", ac.Name)
		assert.Equal(t, "go-api-be-authority-certificate-edit", ac.Name)

		certs, err := svc.GetLoadBalancerCertificateBundles(ctx, &request.GetLoadBalancerCertificateBundlesRequest{})
		assert.NoError(t, err)
		t.Logf("Fetched certificate bundle list (%d)", len(certs))
		assert.Len(t, certs, 2)

		certs, err = svc.GetLoadBalancerCertificateBundles(ctx, &request.GetLoadBalancerCertificateBundlesRequest{
			Page: &request.Page{
				Size:   1,
				Number: 0,
			},
		})
		assert.NoError(t, err)
		t.Logf("Fetched first item in certificate bundle list (%d)", len(certs))
		assert.Len(t, certs, 1)

		cert, err := svc.GetLoadBalancerCertificateBundle(ctx, &request.GetLoadBalancerCertificateBundleRequest{UUID: ac.UUID})
		assert.NoError(t, err)
		t.Logf("Fetched certificate bundle %s", cert.Name)
		assert.Equal(t, "go-api-be-authority-certificate-edit", cert.Name)
		assert.Equal(t, cert.Certificate, ac.Certificate)

		tls, err := svc.CreateLoadBalancerBackendTLSConfig(ctx, &request.CreateLoadBalancerBackendTLSConfigRequest{
			ServiceUUID: lb.UUID,
			BackendName: beName,
			Config: request.LoadBalancerBackendTLSConfig{
				Name:                  mc.Name,
				CertificateBundleUUID: mc.UUID,
			},
		})
		assert.NoError(t, err)
		t.Logf("Created new TLS config %s", tls.Name)
		assert.Equal(t, mc.Name, tls.Name)

		tls, err = svc.ModifyLoadBalancerBackendTLSConfig(ctx, &request.ModifyLoadBalancerBackendTLSConfigRequest{
			ServiceUUID: lb.UUID,
			BackendName: beName,
			Name:        tls.Name,
			Config: request.LoadBalancerBackendTLSConfig{
				Name:                  ac.Name,
				CertificateBundleUUID: ac.UUID,
			},
		})
		assert.NoError(t, err)
		t.Logf("Modified TLS config %s", tls.Name)
		assert.Equal(t, ac.Name, tls.Name)

		be, err := svc.ModifyLoadBalancerBackend(ctx, &request.ModifyLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        beName,
			Backend: request.ModifyLoadBalancerBackend{
				Name:     beName,
				Resolver: upcloud.StringPtr("ns-1"),
				Properties: &upcloud.LoadBalancerBackendProperties{
					TLSEnabled:   upcloud.BoolPtr(true),
					HTTP2Enabled: upcloud.BoolPtr(true),
				},
			},
		})
		assert.NoError(t, err)
		t.Logf("Modified Backend %s", be.Name)
		assert.Equal(t, upcloud.BoolPtr(true), be.Properties.TLSEnabled)
		assert.Equal(t, upcloud.BoolPtr(true), be.Properties.HTTP2Enabled)
		assert.Len(t, be.TLSConfigs, 1)

		configs, err := svc.GetLoadBalancerBackendTLSConfigs(ctx, &request.GetLoadBalancerBackendTLSConfigsRequest{
			ServiceUUID: lb.UUID,
			BackendName: beName,
		})
		assert.NoError(t, err)
		t.Logf("Fetched TLS Config list (%d)", len(configs))
		assert.Len(t, configs, 1)

		tls, err = svc.GetLoadBalancerBackendTLSConfig(ctx, &request.GetLoadBalancerBackendTLSConfigRequest{
			ServiceUUID: lb.UUID,
			BackendName: beName,
			Name:        ac.Name,
		})
		assert.NoError(t, err)
		t.Logf("Fetched TLS Config %s", tls.Name)
		assert.Equal(t, ac.UUID, tls.CertificateBundleUUID)

		err = svc.DeleteLoadBalancerBackendTLSConfig(ctx, &request.DeleteLoadBalancerBackendTLSConfigRequest{
			ServiceUUID: lb.UUID,
			BackendName: feName,
			Name:        tls.Name,
		})
		assert.Error(t, err)
		t.Logf("Backend TLS config %s deletion failed, like it should", tls.Name)

		be, err = svc.ModifyLoadBalancerBackend(ctx, &request.ModifyLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        beName,
			Backend: request.ModifyLoadBalancerBackend{
				Name:     beName,
				Resolver: upcloud.StringPtr("ns-1"),
				Properties: &upcloud.LoadBalancerBackendProperties{
					HTTP2Enabled: upcloud.BoolPtr(false),
				},
			},
		})
		assert.NoError(t, err)
		t.Logf("Modified Backend %s", be.Name)
		assert.Equal(t, upcloud.BoolPtr(false), be.Properties.HTTP2Enabled)
		assert.Len(t, be.TLSConfigs, 1)

		assert.NoError(t, svc.DeleteLoadBalancerBackendTLSConfig(ctx, &request.DeleteLoadBalancerBackendTLSConfigRequest{
			ServiceUUID: lb.UUID,
			BackendName: beName,
			Name:        tls.Name,
		}))
		t.Logf("Deleted TLS config %s", tls.Name)
		assert.NoError(t, svc.DeleteLoadBalancerCertificateBundle(ctx, &request.DeleteLoadBalancerCertificateBundleRequest{UUID: mc.UUID}))
		t.Logf("Deleted certificate bundle %s", mc.Name)
		assert.NoError(t, svc.DeleteLoadBalancerCertificateBundle(ctx, &request.DeleteLoadBalancerCertificateBundleRequest{UUID: ac.UUID}))
		t.Logf("Deleted certificate bundle %s", ac.Name)
	})
}

func TestLoadBalancerBackendMember(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancerbackendmember", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancerAndNetwork(ctx, svc, "nl-ams1", "172.16.3.0/24")
		require.NoError(t, err)
		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()
		t.Logf("Created load balancer for testing LB backend members CRUD: %s", lb.Name)

		backend, err := createLoadBalancerBackend(ctx, svc, lb.UUID)
		require.NoError(t, err)
		t.Logf("Created new backend %s for load balancer %s", backend.Name, lb.Name)

		name := "test_member"
		weight := 100
		maxSessions := 123
		enabled := true
		memberType := upcloud.LoadBalancerBackendMemberTypeStatic
		ip := "10.0.0.2"
		port := 80

		member, err := svc.CreateLoadBalancerBackendMember(ctx, &request.CreateLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Member: request.LoadBalancerBackendMember{
				Name:        name,
				Weight:      weight,
				MaxSessions: maxSessions,
				Enabled:     enabled,
				Type:        memberType,
				IP:          ip,
				Port:        port,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Name, name)
		assert.EqualValues(t, member.Weight, weight)
		assert.EqualValues(t, member.MaxSessions, maxSessions)
		assert.EqualValues(t, member.Enabled, enabled)
		assert.EqualValues(t, member.Type, memberType)
		assert.EqualValues(t, member.IP, ip)
		assert.EqualValues(t, member.Port, port)
		t.Logf("Created new load balancer backend member: %s", member.Name)

		newName := "test_member_TURBO"
		newWeight := 50
		newMaxSessions := 321
		member, err = svc.ModifyLoadBalancerBackendMember(ctx, &request.ModifyLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Name:        member.Name,
			Member: request.ModifyLoadBalancerBackendMember{
				Name:        newName,
				Weight:      upcloud.IntPtr(newWeight),
				MaxSessions: upcloud.IntPtr(newMaxSessions),
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Name, newName)
		assert.EqualValues(t, member.Weight, newWeight)
		assert.EqualValues(t, member.MaxSessions, newMaxSessions)
		t.Logf("Updated load balancers backend member name, weight and max sessions; new name: %s", member.Name)

		newIp := "231.231.231.231"
		newPort := 3003
		member, err = svc.ModifyLoadBalancerBackendMember(ctx, &request.ModifyLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Name:        member.Name,
			Member: request.ModifyLoadBalancerBackendMember{
				Enabled: upcloud.BoolPtr(true),
				IP:      upcloud.StringPtr(newIp),
				Port:    newPort,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Type, memberType)
		assert.EqualValues(t, member.IP, newIp)
		assert.EqualValues(t, member.Port, newPort)
		t.Logf("Updated load balancers backend member type, ip and port: %s", member.Name)

		t.Logf("Get load balancer backend member: %s", newName)
		member, err = svc.GetLoadBalancerBackendMember(ctx, &request.GetLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Name:        newName,
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Name, newName)
		assert.EqualValues(t, member.Enabled, true)
		assert.EqualValues(t, member.Weight, 50)
		assert.EqualValues(t, member.MaxSessions, 321)
		assert.EqualValues(t, member.Type, memberType)
		assert.EqualValues(t, member.IP, newIp)
		assert.EqualValues(t, member.Port, newPort)

		t.Log("Get load balancer backend members")
		members, err := svc.GetLoadBalancerBackendMembers(ctx, &request.GetLoadBalancerBackendMembersRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
		})

		require.NoError(t, err)
		assert.Len(t, members, 2)

		err = svc.DeleteLoadBalancerBackendMember(ctx, &request.DeleteLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Name:        member.Name,
		})
		require.NoError(t, err)
		t.Logf("Deleted load balancer backend member: %s", member.Name)
	})
}

func TestLoadBalancerResolver(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancerresolver", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancerAndNetwork(ctx, svc, "pl-waw1", "10.0.0.0/24")
		require.NoError(t, err)
		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()
		t.Logf("Created load balancer for testing LB resolvers CRUD: %s", lb.Name)

		name := "testname"
		nameServers := []string{"10.0.0.1", "10.0.0.2"}
		retries := 10
		timeout := 20
		timeoutRetry := 10
		cacheValid := 123
		cacheInvalid := 321

		t.Logf("Create resolver: %s", name)
		resolver, err := svc.CreateLoadBalancerResolver(ctx, &request.CreateLoadBalancerResolverRequest{
			ServiceUUID: lb.UUID,
			Resolver: request.LoadBalancerResolver{
				Name:         name,
				Nameservers:  nameServers,
				Retries:      retries,
				Timeout:      timeout,
				TimeoutRetry: timeoutRetry,
				CacheValid:   cacheValid,
				CacheInvalid: cacheInvalid,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, resolver.Name, name)
		assert.EqualValues(t, resolver.Retries, retries)
		assert.EqualValues(t, resolver.Timeout, timeout)
		assert.EqualValues(t, resolver.TimeoutRetry, timeoutRetry)
		assert.EqualValues(t, resolver.CacheValid, cacheValid)
		assert.EqualValues(t, resolver.CacheInvalid, cacheInvalid)
		assert.Len(t, resolver.Nameservers, 2)
		assert.Contains(t, resolver.Nameservers, "10.0.0.1")
		assert.Contains(t, resolver.Nameservers, "10.0.0.2")
		t.Logf("Created resolver %s for load balancer %s", resolver.Name, lb.Name)

		newName := "updated_testname"
		newNameServers := append(nameServers, "10.0.0.3")
		newRetries := 5
		newTimeout := 30

		t.Logf("Update resolver: %s", name)
		resolver, err = svc.ModifyLoadBalancerResolver(ctx, &request.ModifyLoadBalancerResolverRequest{
			ServiceUUID: lb.UUID,
			Name:        resolver.Name,
			Resolver: request.LoadBalancerResolver{
				Name:        newName,
				Nameservers: newNameServers,
				Retries:     newRetries,
				Timeout:     newTimeout,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, resolver.Name, newName)
		assert.EqualValues(t, resolver.Retries, newRetries)
		assert.EqualValues(t, resolver.Timeout, newTimeout)
		assert.Len(t, resolver.Nameservers, 3)
		assert.Contains(t, resolver.Nameservers, "10.0.0.3")
		t.Logf("Modified name, retries, timeout and nameservers for resolver %s", resolver.Name)

		newTimeoutRetry := 15
		newCacheValid := 124
		newCacheInvalid := 324
		resolver, err = svc.ModifyLoadBalancerResolver(ctx, &request.ModifyLoadBalancerResolverRequest{
			ServiceUUID: lb.UUID,
			Name:        resolver.Name,
			Resolver: request.LoadBalancerResolver{
				TimeoutRetry: newTimeoutRetry,
				CacheValid:   newCacheValid,
				CacheInvalid: newCacheInvalid,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, resolver.TimeoutRetry, newTimeoutRetry)
		assert.EqualValues(t, resolver.CacheValid, newCacheValid)
		assert.EqualValues(t, resolver.CacheInvalid, newCacheInvalid)
		t.Logf("Modified timeout_retry, cache_valid and cache_invalid for resolver %s", resolver.Name)

		t.Logf("Get resolver: %s", resolver.Name)
		resolver, err = svc.GetLoadBalancerResolver(ctx, &request.GetLoadBalancerResolverRequest{
			ServiceUUID: lb.UUID,
			Name:        resolver.Name,
		})

		require.NoError(t, err)
		assert.EqualValues(t, resolver.Name, newName)
		assert.Len(t, resolver.Nameservers, 3)
		assert.Contains(t, resolver.Nameservers, "10.0.0.1")
		assert.Contains(t, resolver.Nameservers, "10.0.0.2")
		assert.Contains(t, resolver.Nameservers, "10.0.0.3")

		t.Log("Get resolvers")
		resolvers, err := svc.GetLoadBalancerResolvers(ctx, &request.GetLoadBalancerResolversRequest{
			ServiceUUID: lb.UUID,
		})
		require.NoError(t, err)
		assert.Len(t, resolvers, 1)
		assert.EqualValues(t, resolvers[0].Name, resolver.Name)

		err = svc.DeleteLoadBalancerResolver(ctx, &request.DeleteLoadBalancerResolverRequest{
			ServiceUUID: lb.UUID,
			Name:        resolver.Name,
		})
		require.NoError(t, err)
		t.Logf("Deleted resolver %s for load balancer %s", resolver.Name, lb.Name)
	})
}

func TestGetLoadBalancerPlans(t *testing.T) {
	record(t, "getloadbalancerplans", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		plans, err := svc.GetLoadBalancerPlans(ctx, &request.GetLoadBalancerPlansRequest{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(plans), 2)
		plans, err = svc.GetLoadBalancerPlans(ctx, &request.GetLoadBalancerPlansRequest{Page: &request.Page{
			Size:   1,
			Number: 0,
		}})
		assert.NoError(t, err)
		assert.Len(t, plans, 1)
	})
}

func TestLoadBalancerFrontend(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancerfrontend", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancerAndNetwork(ctx, svc, "de-fra1", "10.0.3.0/24")
		require.NoError(t, err)
		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()
		t.Logf("Created LB for testing frontends: %s", lb.Name)
		be, err := createLoadBalancerBackend(ctx, svc, lb.UUID)
		require.NoError(t, err)
		t.Logf("Created backend %s for testing LB frontends", be.Name)

		fe, err := svc.CreateLoadBalancerFrontend(ctx, &request.CreateLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Frontend: request.LoadBalancerFrontend{
				Name:           "fe-1",
				Mode:           upcloud.LoadBalancerModeHTTP,
				Port:           443,
				DefaultBackend: be.Name,
				Rules:          []request.LoadBalancerFrontendRule{},
				TLSConfigs:     []request.LoadBalancerFrontendTLSConfig{},
				Properties: &upcloud.LoadBalancerFrontendProperties{
					TimeoutClient:        10,
					InboundProxyProtocol: upcloud.BoolPtr(true),
					HTTP2Enabled:         upcloud.BoolPtr(false),
				},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, "fe-1", fe.Name)
		assert.Equal(t, 10, fe.Properties.TimeoutClient)
		assert.Equal(t, upcloud.BoolPtr(true), fe.Properties.InboundProxyProtocol)
		assert.Equal(t, upcloud.BoolPtr(false), fe.Properties.HTTP2Enabled)
		t.Logf("Created frontend %s for load balancer %s", fe.Name, lb.Name)
		fe, err = svc.ModifyLoadBalancerFrontend(ctx, &request.ModifyLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Name:        fe.Name,
			Frontend: request.ModifyLoadBalancerFrontend{
				Name: "fe-2",
				Mode: upcloud.LoadBalancerModeTCP,
				Port: 80,
				Properties: &upcloud.LoadBalancerFrontendProperties{
					HTTP2Enabled: upcloud.BoolPtr(false),
				},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 10, fe.Properties.TimeoutClient)
		assert.Equal(t, upcloud.BoolPtr(true), fe.Properties.InboundProxyProtocol)
		assert.Equal(t, upcloud.BoolPtr(false), fe.Properties.HTTP2Enabled)
		t.Logf("Modified frontend %s", fe.Name)
		fe, err = svc.GetLoadBalancerFrontend(ctx, &request.GetLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Name:        fe.Name,
		})
		require.NoError(t, err)
		assert.Equal(t, "fe-2", fe.Name)
		assert.Equal(t, upcloud.LoadBalancerModeTCP, fe.Mode)
		assert.Equal(t, 80, fe.Port)
		assert.Equal(t, be.Name, fe.DefaultBackend)
		assert.Equal(t, 10, fe.Properties.TimeoutClient)
		assert.Equal(t, upcloud.BoolPtr(true), fe.Properties.InboundProxyProtocol)
		assert.Equal(t, upcloud.BoolPtr(false), fe.Properties.HTTP2Enabled)

		fes, err := svc.GetLoadBalancerFrontends(ctx, &request.GetLoadBalancerFrontendsRequest{ServiceUUID: lb.UUID})
		require.NoError(t, err)
		assert.Equal(t, *fe, fes[0])

		assert.NoError(t,
			svc.DeleteLoadBalancerFrontend(ctx, &request.DeleteLoadBalancerFrontendRequest{
				ServiceUUID: lb.UUID,
				Name:        fe.Name,
			}))
	})
}

func TestLoadBalancerFrontendRule(t *testing.T) {
	t.Parallel()
	const zone = "fi-hel2"
	record(t, "loadbalancerfrontendrule", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		net, err := createLoadBalancerAndPrivateNetwork(ctx, svc, zone, "10.0.1.0/24")
		require.NoError(t, err)
		lb, err := svc.CreateLoadBalancer(ctx, &request.CreateLoadBalancerRequest{
			Name:             fmt.Sprintf("go-api-lb-%s-%d", zone, time.Now().Unix()),
			Zone:             zone,
			Plan:             "development",
			NetworkUUID:      net.UUID,
			ConfiguredStatus: upcloud.LoadBalancerConfiguredStatusStarted,
			Frontends: []request.LoadBalancerFrontend{{
				Name:           "fe-1",
				Mode:           upcloud.LoadBalancerModeHTTP,
				Port:           80,
				DefaultBackend: "be-1",
			}},
			Backends: []request.LoadBalancerBackend{{
				Name:     "be-1",
				Resolver: "ns-1",
				Members:  []request.LoadBalancerBackendMember{},
			}},
			Resolvers: []request.LoadBalancerResolver{{
				Name:         "ns-1",
				Nameservers:  []string{"10.1.1.100"},
				Retries:      10,
				Timeout:      10,
				TimeoutRetry: 10,
				CacheValid:   10,
				CacheInvalid: 10,
			}},
		})

		require.NoError(t, err)
		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()
		rule, err := svc.CreateLoadBalancerFrontendRule(ctx, &request.CreateLoadBalancerFrontendRuleRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: lb.Frontends[0].Name,
			Rule: request.LoadBalancerFrontendRule{
				Name:     "rule-1",
				Priority: 10,
				Matchers: []upcloud.LoadBalancerMatcher{
					{
						Type: upcloud.LoadBalancerMatcherTypeSrcIP,
						SrcIP: &upcloud.LoadBalancerMatcherSourceIP{
							Value: "10.1.1.200",
						},
					},
					{
						Type:    upcloud.LoadBalancerMatcherTypeSrcIP,
						Inverse: upcloud.BoolPtr(true),
						SrcIP: &upcloud.LoadBalancerMatcherSourceIP{
							Value: "10.1.2.200",
						},
					},
				},
				Actions: []upcloud.LoadBalancerAction{{
					Type:      upcloud.LoadBalancerActionTypeTCPReject,
					TCPReject: &upcloud.LoadBalancerActionTCPReject{},
				}},
			},
		})
		require.NoError(t, err)
		t.Logf("Created frontend rule %s", rule.Name)
		assert.Len(t, rule.Actions, 1)
		assert.Len(t, rule.Matchers, 2)
		assert.Equal(t, upcloud.LoadBalancerActionTypeTCPReject, rule.Actions[0].Type)
		assert.Equal(t, upcloud.LoadBalancerMatcherTypeSrcIP, rule.Matchers[0].Type)
		assert.True(t, *rule.Matchers[1].Inverse)
		assert.Equal(t, "10.1.1.200", rule.Matchers[0].SrcIP.Value)
		assert.Equal(t, "rule-1", rule.Name)
		assert.Equal(t, 10, rule.Priority)

		rule, err = svc.ModifyLoadBalancerFrontendRule(ctx, &request.ModifyLoadBalancerFrontendRuleRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: lb.Frontends[0].Name,
			Name:         rule.Name,
			Rule: request.ModifyLoadBalancerFrontendRule{
				Name:     "rule",
				Priority: upcloud.IntPtr(20),
			},
		})
		require.NoError(t, err)
		t.Logf("Modified frontend rule %s", rule.Name)
		assert.Equal(t, "rule", rule.Name)
		assert.Equal(t, 20, rule.Priority)

		rule, err = svc.ReplaceLoadBalancerFrontendRule(ctx, &request.ReplaceLoadBalancerFrontendRuleRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: lb.Frontends[0].Name,
			Name:         rule.Name,
			Rule: request.LoadBalancerFrontendRule{
				Name:     "rule-1",
				Priority: 10,
				Matchers: []upcloud.LoadBalancerMatcher{{
					Type: upcloud.LoadBalancerMatcherTypeSrcIP,
					SrcIP: &upcloud.LoadBalancerMatcherSourceIP{
						Value: "10.1.1.201",
					},
				}},
				Actions: []upcloud.LoadBalancerAction{{
					Type: upcloud.LoadBalancerActionTypeHTTPReturn,
					HTTPReturn: &upcloud.LoadBalancerActionHTTPReturn{
						Status:      404,
						ContentType: "text/html",
						Payload:     "PGgxPmFwcGxlYmVlPC9oMT4K",
					},
				}},
			},
		})

		require.NoError(t, err)
		t.Logf("Replaced frontend rule %s", rule.Name)
		assert.Equal(t, "rule-1", rule.Name)
		assert.Equal(t, 10, rule.Priority)
		assert.Equal(t, upcloud.LoadBalancerActionTypeHTTPReturn, rule.Actions[0].Type)

		t.Logf("Get frontend %s rules", lb.Frontends[0].Name)
		rules, err := svc.GetLoadBalancerFrontendRules(ctx, &request.GetLoadBalancerFrontendRulesRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: lb.Frontends[0].Name,
		})
		require.NoError(t, err)
		assert.Len(t, rules, 1)

		t.Logf("Get frontend rule %s", rule.Name)
		rule, err = svc.GetLoadBalancerFrontendRule(ctx, &request.GetLoadBalancerFrontendRuleRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: lb.Frontends[0].Name,
			Name:         "rule-1",
		})
		require.NoError(t, err)
		require.NoError(t,
			svc.DeleteLoadBalancerFrontendRule(ctx, &request.DeleteLoadBalancerFrontendRuleRequest{
				ServiceUUID:  lb.UUID,
				FrontendName: lb.Frontends[0].Name,
				Name:         rule.Name,
			}))
		t.Logf("Deleted frontend rule %s", rule.Name)
	})
}

// TestLoadBalancerCerticateBundlesAndFrontendTLSConfigs tests certificate bundles and TLS configs
// Test:
//   - certificate bundle CRUD
//   - TLS config CRUD
//   - add TLS config to LB frontend
//   - enable HTTP/2 for LB frontend
//   - remove TLS config to LB frontend
func TestLoadBalancerCerticateBundlesAndFrontendTLSConfigs(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancercerticatebundlesandfrontendtlsconfigs", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		net, err := createLoadBalancerAndPrivateNetwork(ctx, svc, "fi-hel1", "10.0.4.0/24")
		require.NoError(t, err)
		feName := "fe-1"
		lb, err := svc.CreateLoadBalancer(ctx, &request.CreateLoadBalancerRequest{
			Name:             fmt.Sprintf("go-api-lb-%s-%d", "fi-hel1", time.Now().Unix()),
			Zone:             "fi-hel1",
			Plan:             "development",
			NetworkUUID:      net.UUID,
			ConfiguredStatus: upcloud.LoadBalancerConfiguredStatusStarted,
			Frontends: []request.LoadBalancerFrontend{{
				Name:           feName,
				Mode:           upcloud.LoadBalancerModeHTTP,
				Port:           80,
				DefaultBackend: "be-1",
			}},
			Backends: []request.LoadBalancerBackend{{
				Name:     "be-1",
				Resolver: "ns-1",
				Members:  []request.LoadBalancerBackendMember{},
			}},
			Resolvers: []request.LoadBalancerResolver{{
				Name:         "ns-1",
				Nameservers:  []string{"10.1.1.100"},
				Retries:      10,
				Timeout:      10,
				TimeoutRetry: 10,
				CacheValid:   10,
				CacheInvalid: 10,
			}},
		})

		require.NoError(t, err)
		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()

		mc, err := svc.CreateLoadBalancerCertificateBundle(ctx, &request.CreateLoadBalancerCertificateBundleRequest{
			Type:        upcloud.LoadBalancerCertificateBundleTypeManual,
			Name:        "go-api-fe-manual-certificate",
			Certificate: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZhekNDQTFPZ0F3SUJBZ0lVRzR1KzRDZmlHQ3pQSDk4dDA4QXh5VkE0QzVnd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1JURUxNQWtHQTFVRUJoTUNRVlV4RXpBUkJnTlZCQWdNQ2xOdmJXVXRVM1JoZEdVeElUQWZCZ05WQkFvTQpHRWx1ZEdWeWJtVjBJRmRwWkdkcGRITWdVSFI1SUV4MFpEQWVGdzB5TWpBek1UY3hNekUxTURoYUZ3MHlNekF6Ck1UY3hNekUxTURoYU1FVXhDekFKQmdOVkJBWVRBa0ZWTVJNd0VRWURWUVFJREFwVGIyMWxMVk4wWVhSbE1TRXcKSHdZRFZRUUtEQmhKYm5SbGNtNWxkQ0JYYVdSbmFYUnpJRkIwZVNCTWRHUXdnZ0lpTUEwR0NTcUdTSWIzRFFFQgpBUVVBQTRJQ0R3QXdnZ0lLQW9JQ0FRRHZaeG4vK3pUeVc0RVJ2S2t3V29tVXBpOG8ydEp6MWR2ZXIrREpySzNnCkNObFVvWXpSMjlDV3M3aks4MVhNc3ZtcUw1TXpUd1A3SHNtZDFxNjlGSStXY1BFMWFhYjk5MDlJQWsvR0dpSzIKelRsZU4zRVFRcFhuN3RueVB0WmFUOFkxM3lGSHBDNVJnUXpURThDUjlaaTJPTEV5eEdRMzZwQTYxOTBueFZnMgpTTGxhZk5HVFp0SnZOMS83cjltSmhFbGJyVUUram9lWEx3Tm9qSC9uWGs1Vy9Yd3paYm9JSHNTRlZZaksyemxnCm9xQzYrQXBvOXhGOW9ZN25sQWhRMEtLV3ZRVmJ3akdQbVZOMTdFVG9kSHNLSlpCb1h4RHNaVVRHQ0RESkNpbXoKVzY0YTc5bFdJeGl5T1E0LzdUbjJGaFBZMG9tSDVVYldDUHEyTW5YZWJrT2pnY3ZVWTROSXd6cFlWMFcyR0dHRwp3d25pOWZsbFlBTTlPRDNidlBYNU9hQVdOSlQ2cjZFYXNzaldsdjBUZUd4RStCWlorZzN4UHFIVEd6MndIekM1CjVhbkxEak0rNHZzQlZrZmtWM1NZN1c4M203NFZRK1FhM1dhTlh6aW5MMGtlRnh1cExYWThiS2hFelh6U0xLeisKQnI4UEdlR1JnYVNEZDFrcEZQZyt1ak44cXZnbzBSREk4SXFMUzd6YlhGb1FycDF4L2RXbTlTOERWRVhWb1VBMQpXUW5WdVdFQ29CUzRaZjQxZDA0cGZkQ3R0bk45ekhvc2d3WGJKOG0wVGZ2Zmt1aFZpdVZBTi9wK01wOVduUStICjExSEVuV3BTZk9oN1pQalR6anVBc2V2VmZWNGc0YTNrY3pNdjFycE5QelVVUHd0QXF4OTIzOXd3SVI5WTE1Y2wKT1FJREFRQUJvMU13VVRBZEJnTlZIUTRFRmdRVVhJcWxiajV1TVVGVC9qcU0ya1d2WVp0RE5rY3dId1lEVlIwagpCQmd3Rm9BVVhJcWxiajV1TVVGVC9qcU0ya1d2WVp0RE5rY3dEd1lEVlIwVEFRSC9CQVV3QXdFQi96QU5CZ2txCmhraUc5dzBCQVFzRkFBT0NBZ0VBQ00reGJiOW9rVUdPcWtRWmtndHh3eVF0em90WFRXNnlIYmdmeGhCd3d1disKc0ltT1QyOUw5WFVZblM5cTgrQktBSHVGa3JTRUlwR0ZVVWhDVWZTa2xlNmZuR0c1b05PWm1DdEszN3RoZlFVOQp2NEJvOHFCWkhqREh3azlWVHRhWm1BazBLYnhmaHVneVdWQ1ZsbURacm9TQ09pV0drVFZoc1hhS0RrYnc0RWwwCjJzY3lnYkFDdFZ4bkU1WjlmU0F3MU9QWXJZYUcySW5HTDQvMHVSZXo4aXl1UE9lNUNiL0RkNDl1eHFzR1FkM1IKQzdKNC9vWnB2b0V6UVJtakxib1FzQzkwU2ZqaFNpcGhHQlNiYUpCZGRsMDBrNVZzVXJxS1haU004cVFxVWZWLwpubEJtYjJOblVsa2RlOEtIczBQamhCaG8rLzdmaitMN21GYTJsNWpmdWlsdHVxdmgyWnladFJjd2didmJlaUxPCmZQSWlMQ2dTbnMwaitZMkVrS1drRUp6RXJQVm5sOTdaQktZclBaYmRYMFY5b2dvTC9qeEV5NzlsbzlKczI5djYKUkY2NmdvSlUwMkVKZTUwMmk3WHJzMzFZQ0tuSGd2ejUwTDZha0JpYWRSNmtrTXVXdkJ1d1l6MElaS1RMcXhqZAowOEdlUkJVeWFsUFZodGZKbzNNdXRuYUllL1pWVTdLQUl3S1Znb20zS09EY1RpWllQV3RWKzFnL0UvN3A1aGh2CkJERzFqcklRc1ZrZG4yNWZhNXNkNU9Qa1AvbDBRdXY1em16UEk3S1MrS2ZlWS92NHFBOTBtNGk2dkZORlRtbTAKSFNXV0JZTlR4blIxYjk2UElUcnRzOE15am9YTFg2QnUxVkZOSlByMkpnMDJMVlZvcTZSSWJlMVVvNjE5b2pBPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==",
			PrivateKey:  "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUpSUUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQ1M4d2dna3JBZ0VBQW9JQ0FRRHZaeG4vK3pUeVc0RVIKdktrd1dvbVVwaThvMnRKejFkdmVyK0RKckszZ0NObFVvWXpSMjlDV3M3aks4MVhNc3ZtcUw1TXpUd1A3SHNtZAoxcTY5RkkrV2NQRTFhYWI5OTA5SUFrL0dHaUsyelRsZU4zRVFRcFhuN3RueVB0WmFUOFkxM3lGSHBDNVJnUXpUCkU4Q1I5WmkyT0xFeXhHUTM2cEE2MTkwbnhWZzJTTGxhZk5HVFp0SnZOMS83cjltSmhFbGJyVUUram9lWEx3Tm8KakgvblhrNVcvWHd6WmJvSUhzU0ZWWWpLMnpsZ29xQzYrQXBvOXhGOW9ZN25sQWhRMEtLV3ZRVmJ3akdQbVZOMQo3RVRvZEhzS0paQm9YeERzWlVUR0NEREpDaW16VzY0YTc5bFdJeGl5T1E0LzdUbjJGaFBZMG9tSDVVYldDUHEyCk1uWGVia09qZ2N2VVk0Tkl3enBZVjBXMkdHR0d3d25pOWZsbFlBTTlPRDNidlBYNU9hQVdOSlQ2cjZFYXNzalcKbHYwVGVHeEUrQlpaK2czeFBxSFRHejJ3SHpDNTVhbkxEak0rNHZzQlZrZmtWM1NZN1c4M203NFZRK1FhM1dhTgpYemluTDBrZUZ4dXBMWFk4YktoRXpYelNMS3orQnI4UEdlR1JnYVNEZDFrcEZQZyt1ak44cXZnbzBSREk4SXFMClM3emJYRm9RcnAxeC9kV205UzhEVkVYVm9VQTFXUW5WdVdFQ29CUzRaZjQxZDA0cGZkQ3R0bk45ekhvc2d3WGIKSjhtMFRmdmZrdWhWaXVWQU4vcCtNcDlXblErSDExSEVuV3BTZk9oN1pQalR6anVBc2V2VmZWNGc0YTNrY3pNdgoxcnBOUHpVVVB3dEFxeDkyMzl3d0lSOVkxNWNsT1FJREFRQUJBb0lDQVFDcUNtd2dNbmcvOEJoejFiSENRM3hYCkZkYUhTUzJUMHdHUllSRGpqZ0FPRVpyMERxN3IzQnFDLy9Jd1RMZlRaZ2dKQmpPaWpPd0I4TE01cGVPRkwxWngKZjVVRDRDQVpZUkJ4MEJxRFZjcjBWajM2R3B6MjlLUnZFV3JDTWptaitlZUtHZ3NVVEp3TmpnRGk1N091dUdlWQpmaG4yT2lJSXlWVmFSanF4NWV5cTJlcTFSOVMvd3BlVElSek9zdTlyU29la1V5SDFZZDBTMS9TdXpLU0lYS1orCkNSdXZrZ0NaaGVrRjMyUUMyY1VlUzBTb3FFY1VtUEJXY0dzRk4xTFV1K3ZQNzBBZ0ZZV0lQbHBXZHRQVzIrME0KbnZPNy9sSVI1amY4QkpOS0tDcklWMFVKb3ZTV3h1VGlxYjNpVUFnTUwxQTNnQXJwZUVOaEFRMjZYWXIweXhMRQpiMTRObWdnZzRUSktRMWp4aTZlRys5Nkp0WStteFdwaHJUNTVUT2s5blpEY2JTV29ibTgwcExJZEE0QTN0bjhJCjI5SXZ2dkdhVnh2NXFXaU5sL0sxWEZhK2JRRVY3K1AxT0RWSnV5VXEzdFg4dDlONVJ2c3RiUE5XNU8xQ01STGsKZExESVMwRFFKYytMa1ZGaGpZRXlJOWluQmtlRXV1ekI0L0k3M1JnTXpKMFZvZU1hSXEwdWpxbVF6KzN0ei9JUgp6VTFSN2FndEZmNHUvTS9jeVl2U0E4UEZVSis1Q25ldHFtWC91dzR0em10WmtLd3JnMVArblVpRTV4Z0ZLT1FZCjUyaU81aXRKWHo4aHJGdFpmbVMrcTI3eFRFWWlEdTFFSFg0Y1pLaFBQNVh3RzkyekN0ZUxJenc5QUh2TS9aRTcKNmI0OFMwQWR6T3dFN3VaVVhheEw5UUtDQVFFQTkrOVROV2RuMjVxK3lSWUlsc2x1NFI1N0xSQm1ucXFhaDljSQowZ1JGS0RZUFZTWFl0RlhVci9mK0psSmR4YUxEYjYwV0o4c2hWcHZ1MTlJTE0wNTR6M3ZYa0w0QjQwSWI1Z3JnCnlzbVZKVlpoYVdTNWU5RW02TEh2bXRoVWl0MzVDWnlLcTlkN2FwdzFOQ1VIUkFEWE9wSndBWkc4ckttRmtDeFYKTnFpVThnWi9LOXpPK0gra215bGVKZTkrZHhCR2RCZTczbnFsRTVPQVdJcjVLcFNpMU5sQWYxOG0yYnVUWkxiTQo2Rnd1MEc0SjUzd1ZaTFRhbWxVelNmdjFTRGRrMEZIbnpkdEJPMWI3OTNJeVpOQXc4eGlkMGVPVy90OE1QTm1RCmFXWnNqaWc4WEZVMlFaMXl5YUd6amoxNUpoVzl2WGN1SUpGc21rRzZiRS8wSm9ZZnd3S0NBUUVBOXpDNXNYN1IKQlN2MzlwVnFROXI4bDJXNVVqQk9iczFick05Q0NPQWFFZCtlY2xyZWZZRU5yMUVtL0x3Qk9GK1RuRnVnMjJZOQpKczliNjgrY0wzSEtweVJ5KzVVaGM1NzU0b3pDRklDemJzZHphM005TGVzVUtTVTBYYnlhNm5vc0UvWjNOWGpyCmpLQ1ROZEU2eFY4OTZ1Vm5FTXZpcDh4M1kvR0VBejh4U3lCOUFveGNtanVqZnVVdmlZbmN5SnBqUUoxOEhsZk8KMlloWWdTd3VxSFhISUswUXhGckJiTjIzZXZ5TUkrOVVSSVlnOWxtemFWR2hqczlHd0Y3UUQ0dFoyOXgwQVFOSwpBUFArMnI4ZUlLbjFucEs4VnFRdVNrNVdsZjg3L0xvV1Zha2xtRXkvdEJ0alpjNVBQanpMVmRTLy9QNWlkc2F1CkZuTnc2VDNmZWEwelV3S0NBUUVBMWsySDU2WXNzRFhPYUxOaDB5dmphalJWbGJzU2FGemdXei8wQU12dUZ2YTcKUkFjRmk4S1FwMVU4MlZUaWRyemNIc0JHWVRrRDVQKzliOUMvRzZiZFo4SU1yckI5b3poMk10MytOV29PUDRxdApnbEtzdktnbzhJTTBydXdFRDFBVVBVbVExejNYRUd4YTFHcVpJQjkxNmN1L2dxdThvS1dhcSthVjlUdThHb0toCkU0RzFhRGUwU09WMTJtWnJNbkRmNU9MSzRWK3pKZnVkdVdyT09nN2x2QUxZNi8rTDdqRmpFbStySjhEZU9neVQKQlFKTTM1SXZUYTBOT3dyTWxaSkQwb2lwUzFjVHlEM0VacnJQY2pJOXpUSGU0QmZQWVJmY1ZSQmM4YTIxY1I2NApKYnNGdmF0aEY0VnNWU3N2ZDByZGlWSGxqZ01GRTBSeTVjSXFMODVJendLQ0FRRUEzdE1nZ1QwRkpIbFhKQVBhCmIrS0drZDlUNkIrOWhDcEFPbzMyUTlQb0REYWRPUTVxdzQzRERVZkZNa3d6ZVdMR3lFcmN2UW56ay9tV0xnTFAKRXdHcm9YRzg2TWF0Q2ZIRDVoSG1uZDdLWU5FUVhVcmJXbm92aVV1TllmWXpXNnpYOFFMYXdPd0l3Wks2UU9nago1Mm1NZ2lOYS9nd2NmQkJYaTFOYUlpY2p3MG85QmtBSzljbE8vNE9QajVjajIvMDMvVFk1Zll5LzNONElraUNHCnlycW96dTdUVDMxVUlWUFlJdGhuWjdsRktDUVVzSjE1bWpYSXdkaGRPZW45K2hVdTRuOWVYczlkTlhDOVN1aS8KT3NpYXJlQXVRSmZ0Vm5RNW55c2VJeHFJS1oyNVV3blVRWUh5M3dIVDh4R1FaZ1hMTHo4TStXN3QzVFVoRWxBQgpGRWtxR3dLQ0FRRUFtT2VzMXFINlA0dDRsZ0VjK01Ubzc2VmJvYnhab29HQWdMWmc3N1J4c2hPaGxMdCtYZnIzCjFsOWgycFJ0eXFIdDRuMG9ob1A0VzNSN3VuQ05rNWl3SGNJSDVjNmhGQTYyelVEM1JjeEhJZERhdDBGL3RoWDgKTUpndDlsM0MrMVJwZFZlV0hlbEJOM0JlM0FtWEpXL1ZsL3lGTXVjWWxETnlXUFlPSmRuQ1BRZ0FGVFJnVnJlUQpiUjZCY29neUVRTVEzenRMWnNBRnRaZ25Sak05YkpLN1JjYzg5bGxaa1BuMHVKZjNKVWxMeVFFN0l2bEJsWi9tClZnUUhiRTkwQStZNzFpb1piQWh5TFcwTE9lTmhBS3NRNFJZbDlnb0N4dGp1ZnE2NnFTNGNzdGN6c2J5N083dFAKeXZkSXp2eEZRZmx4Yk8ra1ludDNkcFRIdUNuUkNIMFM0UT09Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K",
		})
		require.NoError(t, err)
		t.Logf("Created manual certificate bundle %s", mc.Name)
		assert.Equal(t, "go-api-fe-manual-certificate", mc.Name)
		assert.Equal(t, upcloud.LoadBalancerCertificateBundleTypeManual, mc.Type)

		dc, err := svc.CreateLoadBalancerCertificateBundle(ctx, &request.CreateLoadBalancerCertificateBundleRequest{
			Type: upcloud.LoadBalancerCertificateBundleTypeDynamic,
			Name: "go-api-fe-dynamic-certificate",
			Hostnames: []string{
				"example.com",
				"app.example.com",
			},
			KeyType: "rsa",
		})
		require.NoError(t, err)
		t.Logf("Created dynamic certificate bundle %s", dc.Name)
		assert.Equal(t, "go-api-fe-dynamic-certificate", dc.Name)
		assert.Equal(t, upcloud.LoadBalancerCertificateBundleTypeDynamic, dc.Type)

		mc, err = svc.ModifyLoadBalancerCertificateBundle(ctx, &request.ModifyLoadBalancerCertificateBundleRequest{
			UUID: mc.UUID,
			Name: "go-api-fe-manual-certificate-edit",
		})
		assert.NoError(t, err)
		t.Logf("Modified manual certificate bundle %s", mc.Name)
		assert.Equal(t, "go-api-fe-manual-certificate-edit", mc.Name)

		dc, err = svc.ModifyLoadBalancerCertificateBundle(ctx, &request.ModifyLoadBalancerCertificateBundleRequest{
			UUID: dc.UUID,
			Name: "go-api-fe-dynamic-certificate-edit",
		})
		assert.NoError(t, err)
		t.Logf("Modified dynamic certificate bundle %s", dc.Name)
		assert.Equal(t, "go-api-fe-dynamic-certificate-edit", dc.Name)

		certs, err := svc.GetLoadBalancerCertificateBundles(ctx, &request.GetLoadBalancerCertificateBundlesRequest{})
		assert.NoError(t, err)
		t.Logf("Fetched certificate bundle list (%d)", len(certs))
		assert.Len(t, certs, 2)

		certs, err = svc.GetLoadBalancerCertificateBundles(ctx, &request.GetLoadBalancerCertificateBundlesRequest{
			Page: &request.Page{
				Size:   1,
				Number: 0,
			},
		})
		assert.NoError(t, err)
		t.Logf("Fetched first item in certificate bundle list (%d)", len(certs))
		assert.Len(t, certs, 1)

		cert, err := svc.GetLoadBalancerCertificateBundle(ctx, &request.GetLoadBalancerCertificateBundleRequest{UUID: dc.UUID})
		assert.NoError(t, err)
		t.Logf("Fetched certificate bundle %s", cert.Name)
		assert.Equal(t, "go-api-fe-dynamic-certificate-edit", cert.Name)
		assert.Len(t, cert.Hostnames, 2)

		tls, err := svc.CreateLoadBalancerFrontendTLSConfig(ctx, &request.CreateLoadBalancerFrontendTLSConfigRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: feName,
			Config: request.LoadBalancerFrontendTLSConfig{
				Name:                  mc.Name,
				CertificateBundleUUID: mc.UUID,
			},
		})
		assert.NoError(t, err)
		t.Logf("Created new TLS config %s", tls.Name)
		assert.Equal(t, mc.Name, tls.Name)

		tls, err = svc.ModifyLoadBalancerFrontendTLSConfig(ctx, &request.ModifyLoadBalancerFrontendTLSConfigRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: feName,
			Name:         tls.Name,
			Config: request.LoadBalancerFrontendTLSConfig{
				Name:                  dc.Name,
				CertificateBundleUUID: dc.UUID,
			},
		})
		assert.NoError(t, err)
		t.Logf("Modified TLS config %s", tls.Name)
		assert.Equal(t, dc.Name, tls.Name)

		fe, err := svc.ModifyLoadBalancerFrontend(ctx, &request.ModifyLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Name:        feName,
			Frontend: request.ModifyLoadBalancerFrontend{
				Name: feName,
				Mode: upcloud.LoadBalancerModeTCP,
				Port: 80,
				Properties: &upcloud.LoadBalancerFrontendProperties{
					InboundProxyProtocol: upcloud.BoolPtr(true),
					HTTP2Enabled:         upcloud.BoolPtr(true),
				},
			},
		})
		assert.NoError(t, err)
		t.Logf("Modified Frontend %s", fe.Name)
		assert.Equal(t, upcloud.BoolPtr(true), fe.Properties.InboundProxyProtocol)
		assert.Equal(t, upcloud.BoolPtr(true), fe.Properties.HTTP2Enabled)
		assert.Len(t, fe.TLSConfigs, 1)

		configs, err := svc.GetLoadBalancerFrontendTLSConfigs(ctx, &request.GetLoadBalancerFrontendTLSConfigsRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: feName,
		})
		assert.NoError(t, err)
		t.Logf("Fetched TLS Config list (%d)", len(configs))
		assert.Len(t, configs, 1)

		tls, err = svc.GetLoadBalancerFrontendTLSConfig(ctx, &request.GetLoadBalancerFrontendTLSConfigRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: feName,
			Name:         dc.Name,
		})
		assert.NoError(t, err)
		t.Logf("Fetched TLS Config %s", tls.Name)
		assert.Equal(t, dc.UUID, tls.CertificateBundleUUID)

		err = svc.DeleteLoadBalancerFrontendTLSConfig(ctx, &request.DeleteLoadBalancerFrontendTLSConfigRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: feName,
			Name:         tls.Name,
		})
		assert.Error(t, err)
		t.Logf("Frontend TLS config %s deletion failed, like it should", tls.Name)

		fe, err = svc.ModifyLoadBalancerFrontend(ctx, &request.ModifyLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Name:        feName,
			Frontend: request.ModifyLoadBalancerFrontend{
				Name: feName,
				Mode: upcloud.LoadBalancerModeTCP,
				Port: 80,
				Properties: &upcloud.LoadBalancerFrontendProperties{
					InboundProxyProtocol: upcloud.BoolPtr(true),
					HTTP2Enabled:         upcloud.BoolPtr(false),
				},
			},
		})
		assert.NoError(t, err)
		t.Logf("Modified Frontend %s", fe.Name)
		assert.Equal(t, upcloud.BoolPtr(true), fe.Properties.InboundProxyProtocol)
		assert.Equal(t, upcloud.BoolPtr(false), fe.Properties.HTTP2Enabled)
		assert.Len(t, fe.TLSConfigs, 1)

		assert.NoError(t, svc.DeleteLoadBalancerFrontendTLSConfig(ctx, &request.DeleteLoadBalancerFrontendTLSConfigRequest{
			ServiceUUID:  lb.UUID,
			FrontendName: feName,
			Name:         tls.Name,
		}))
		t.Logf("Deleted TLS config %s", tls.Name)
		assert.NoError(t, svc.DeleteLoadBalancerCertificateBundle(ctx, &request.DeleteLoadBalancerCertificateBundleRequest{UUID: mc.UUID}))
		t.Logf("Deleted certificate bundle %s", mc.Name)
		assert.NoError(t, svc.DeleteLoadBalancerCertificateBundle(ctx, &request.DeleteLoadBalancerCertificateBundleRequest{UUID: dc.UUID}))
		t.Logf("Deleted certificate bundle %s", dc.Name)
	})
}

func TestLoadBalancerPage(t *testing.T) {
	// do not run this test in parallel because it alters request.DefaultPage config which might cause unexpected results
	const zone = "fi-hel2"
	record(t, "loadbalancerpage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		net, err := createLoadBalancerAndPrivateNetwork(ctx, svc, zone, "172.16.0.0/24")
		require.NoError(t, err)
		lbs := make([]*upcloud.LoadBalancer, 0)
		for i := 0; i < 5; i++ {
			lb, err := createLoadBalancer(ctx, svc, net.UUID, zone)
			require.NoError(t, err)
			lbs = append(lbs, lb)
			if rec.Mode() != recorder.ModeReplaying {
				time.Sleep(1 * time.Second)
			}
		}

		// Test get-all feature by altering default page config.
		// We have 5 load balancers so this should create 3 requests and combine results
		// /load-balancer?limit=2&offset=0 (2)
		// /load-balancer?limit=2&offset=2 (2)
		// /load-balancer?limit=2&offset=4 (1) partial page
		tmp := *request.DefaultPage
		request.DefaultPage = &request.Page{
			Size:   2,
			Number: 1,
		}
		list, err := svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{})
		request.DefaultPage = &tmp // restore default
		assert.NoError(t, err)
		assert.Len(t, list, 5)
		assert.Equal(t, request.PageSizeMax, request.DefaultPage.Size)

		// test custom page
		list, err = svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{Page: &request.Page{
			Size:   2,
			Number: 0,
		}})
		assert.NoError(t, err)
		assert.Len(t, list, 2)

		// flush load balancers
		for _, lb := range lbs {
			if err := svc.DeleteLoadBalancer(ctx, &request.DeleteLoadBalancerRequest{UUID: lb.UUID}); err != nil {
				t.Logf("an error occurred when deleting LB '%s': %v", lb.Name, err)
				continue
			}
		}
		for _, lb := range lbs {
			if err := waitForLoadBalancerToShutdown(ctx, rec, svc, lb); err != nil {
				t.Log(err)
				continue
			}
		}
		if err := svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{UUID: net.UUID}); err != nil {
			t.Log(err)
		}

		list, err = svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{})
		assert.NoError(t, err)
		assert.Len(t, list, 0)
	})
}

func TestLoadBalancerNetwork(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancernetworks", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		const zone string = "pl-waw1"
		net, err := createLoadBalancerPrivateNetwork(ctx, svc, zone, "192.168.55.0/24")
		require.NoError(t, err)
		t.Log("Creating LB for testing networks")
		lb, err := svc.CreateLoadBalancer(ctx, &request.CreateLoadBalancerRequest{
			Name:             fmt.Sprintf("go-api-lb-networks-%s-%d", zone, time.Now().Unix()),
			Zone:             zone,
			Plan:             "development",
			ConfiguredStatus: "started",
			Networks: []request.LoadBalancerNetwork{
				{
					Name:   "public_net",
					Type:   upcloud.LoadBalancerNetworkTypePublic,
					Family: upcloud.LoadBalancerAddressFamilyIPv4,
				},
				{
					Name:   "private_net",
					Type:   upcloud.LoadBalancerNetworkTypePrivate,
					Family: upcloud.LoadBalancerAddressFamilyIPv4,
					UUID:   net.UUID,
				},
			},
			Frontends: []request.LoadBalancerFrontend{{
				Name:           "fe1_network_test",
				Mode:           upcloud.LoadBalancerModeHTTP,
				DefaultBackend: "be_network_test",
				Port:           80,
				Networks: []upcloud.LoadBalancerFrontendNetwork{
					{
						Name: "public_net",
					},
					{
						Name: "private_net",
					},
				},
			}},
			Backends: []request.LoadBalancerBackend{{
				Name:    "be_network_test",
				Members: make([]request.LoadBalancerBackendMember, 0),
			}},
			Resolvers: make([]request.LoadBalancerResolver, 0),
		})
		require.NoError(t, err)
		defer func() {
			err := cleanupLoadBalancer(ctx, rec, svc, lb)
			assert.NoError(t, err)
		}()

		t.Logf("Testing LB %s properties", lb.Name)
		assert.Len(t, lb.Networks, 2)
		assert.Len(t, lb.Frontends[0].Networks, 2)
		assert.Equal(t, upcloud.LoadBalancerNetworkTypePublic, lb.Networks[0].Type)
		assert.Equal(t, upcloud.LoadBalancerNetworkTypePrivate, lb.Networks[1].Type)
		assert.Equal(t, "public_net", lb.Frontends[0].Networks[0].Name)
		assert.Equal(t, "private_net", lb.Frontends[0].Networks[1].Name)

		t.Logf("Creating new LB %s frontend", lb.Name)
		fe, err := svc.CreateLoadBalancerFrontend(ctx, &request.CreateLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Frontend: request.LoadBalancerFrontend{
				Name:           "fe2_network_test",
				Mode:           upcloud.LoadBalancerModeHTTP,
				Port:           443,
				DefaultBackend: lb.Frontends[0].DefaultBackend,
				Networks: []upcloud.LoadBalancerFrontendNetwork{
					{
						Name: "private_net",
					},
				},
			},
		})
		require.NoError(t, err)
		assert.Len(t, fe.Networks, 1)

		t.Logf("Waiting LB %s to become online", lb.Name)
		if err = waitLoadBalancerOnline(ctx, t, rec, svc, lb.UUID); err != nil {
			t.Error(err)
			return
		}

		t.Logf("Getting LB %s details", lb.Name)
		lb, err = svc.GetLoadBalancer(ctx, &request.GetLoadBalancerRequest{UUID: lb.UUID})
		require.NoError(t, err)
		require.Len(t, lb.Nodes, 1)
		require.Len(t, lb.Frontends, 2)
		require.Len(t, lb.Nodes[0].Networks, 2)
		assert.Len(t, lb.Networks, 2)
		assert.Len(t, lb.Frontends[0].Networks, 2)
		assert.Len(t, lb.Frontends[1].Networks, 1)
		assert.Equal(t, "private_net", lb.Frontends[1].Networks[0].Name)

		t.Logf("Modifying LB %s network", lb.Name)
		lbNet, err := svc.ModifyLoadBalancerNetwork(ctx, &request.ModifyLoadBalancerNetworkRequest{
			ServiceUUID: lb.UUID,
			Name:        "private_net",
			Network: request.ModifyLoadBalancerNetwork{
				Name: "internal_net",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, "internal_net", lbNet.Name)

		t.Logf("Modifying LB %s network using context-less service", lb.Name)
		lbNet, err = svc.ModifyLoadBalancerNetwork(ctx, &request.ModifyLoadBalancerNetworkRequest{
			ServiceUUID: lb.UUID,
			Name:        "internal_net",
			Network: request.ModifyLoadBalancerNetwork{
				Name: "private_net",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, "private_net", lbNet.Name)
	})
}

func TestLoadBalancerLabels(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancerlabels", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create Load Balancers
		lb1, err := createLoadBalancerAndNetwork(ctx, svc, "fi-hel2", "172.16.12.0/24", upcloud.Label{Key: "zone", Value: "hel2"})
		require.NoError(t, err)
		t.Logf("Created load balancer: %s", lb1.Name)

		lb2, err := createLoadBalancerAndNetwork(ctx, svc, "fi-hel1", "172.16.14.0/24", upcloud.Label{Key: "zone", Value: "hel1"})
		require.NoError(t, err)
		t.Logf("Created load balancer: %s", lb2.Name)
		defer func() {
			if err := cleanupLoadBalancer(ctx, rec, svc, lb1); err != nil {
				t.Log(err)
			}
			if err := cleanupLoadBalancer(ctx, rec, svc, lb2); err != nil {
				t.Log(err)
			}
		}()

		// Get Load Balancers
		t.Log("Get load balancers labeled as 'zone'")
		lbs, err := svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{
			Filters: []request.QueryFilter{
				request.FilterLabelKey{Key: "zone"},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, len(lbs), 2)

		t.Log("Get load balancers labeled as 'zone' with limited page")
		lbs, err = svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{
			Page: &request.Page{
				Size: 1,
			},
			Filters: []request.QueryFilter{
				request.FilterLabelKey{Key: "zone"},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, len(lbs), 1)

		t.Log("Get load balancers labeled as 'zone=hel1'")
		lbs, err = svc.GetLoadBalancers(ctx, &request.GetLoadBalancersRequest{
			Filters: []request.QueryFilter{
				request.FilterLabel{
					Label: upcloud.Label{Key: "zone", Value: "hel1"},
				},
			},
		})
		require.NoError(t, err)
		if assert.Equal(t, len(lbs), 1) {
			assert.Equal(t, "fi-hel1", lbs[0].Zone)
		}

		lb1, err = svc.ModifyLoadBalancer(ctx, &request.ModifyLoadBalancerRequest{
			UUID:   lb1.UUID,
			Labels: &[]upcloud.Label{},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, len(lb1.Labels))
	})
}

func waitLoadBalancerOnline(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, UUID string) error {
	t.Helper()
	if rec.Mode() == recorder.ModeRecording {
		rec.AddPassthrough(func(h *http.Request) bool {
			return true
		})
		defer func() {
			rec.Passthroughs = nil
		}()
		for s := time.Now(); time.Since(s) < time.Minute*20; {
			lb, err := svc.GetLoadBalancer(ctx, &request.GetLoadBalancerRequest{UUID: UUID})
			if err != nil {
				return err
			}
			if lb.OperationalState == upcloud.LoadBalancerOperationalStateRunning {
				return nil
			}
			time.Sleep(time.Second * 2)
		}
	}
	return nil
}

func createLoadBalancerBackend(ctx context.Context, svc *Service, lbUUID string) (*upcloud.LoadBalancerBackend, error) {
	req := request.CreateLoadBalancerBackendRequest{
		ServiceUUID: lbUUID,
		Backend: request.LoadBalancerBackend{
			Name: fmt.Sprintf("go-api-lb-backend-%d", time.Now().Unix()),
			Properties: &upcloud.LoadBalancerBackendProperties{
				TimeoutServer:  30,
				TLSEnabled:     upcloud.BoolPtr(true),
				TLSUseSystemCA: upcloud.BoolPtr(true),
				TLSVerify:      upcloud.BoolPtr(true),
			},
			Members: []request.LoadBalancerBackendMember{
				{
					Name:        "default-lb-backend-member",
					Type:        "dynamic",
					Weight:      100,
					MaxSessions: 1000,
					Enabled:     true,
					Port:        8000,
					IP:          "196.123.123.123",
				},
			},
			TLSConfigs: []request.LoadBalancerBackendTLSConfig{},
		},
	}

	return svc.CreateLoadBalancerBackend(ctx, &req)
}

func cleanupLoadBalancer(ctx context.Context, rec *recorder.Recorder, svc *Service, lb *upcloud.LoadBalancer) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}
	err := deleteLoadBalancer(ctx, svc, lb)
	return err
}

func createLoadBalancer(ctx context.Context, svc *Service, networkUUID, zone string, label ...upcloud.Label) (*upcloud.LoadBalancer, error) {
	createLoadBalancerRequest := request.CreateLoadBalancerRequest{
		Name:             fmt.Sprintf("go-api-lb-%s-%d", zone, time.Now().Unix()),
		Zone:             zone,
		Plan:             "development",
		NetworkUUID:      networkUUID,
		ConfiguredStatus: "started",
		Frontends:        []request.LoadBalancerFrontend{},
		Backends:         []request.LoadBalancerBackend{},
		Resolvers:        []request.LoadBalancerResolver{},
		MaintenanceDOW:   upcloud.LoadBalancerMaintenanceDOWSunday,
		MaintenanceTime:  "20:01:01Z",
	}
	if len(label) > 0 {
		createLoadBalancerRequest.Labels = label
	}
	loadBalancerDetails, err := svc.CreateLoadBalancer(ctx, &createLoadBalancerRequest)
	if err != nil {
		return nil, err
	}

	return loadBalancerDetails, nil
}

func createLoadBalancerAndNetwork(ctx context.Context, svc *Service, zone, addr string, label ...upcloud.Label) (*upcloud.LoadBalancer, error) {
	n, err := createLoadBalancerAndPrivateNetwork(ctx, svc, zone, addr)
	if err != nil {
		return nil, err
	}
	return createLoadBalancer(ctx, svc, n.UUID, zone, label...)
}

func createLoadBalancerAndPrivateNetwork(ctx context.Context, svc *Service, zone, addr string) (*upcloud.Network, error) {
	return svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
		Name: fmt.Sprintf("go-api-lb-%d", time.Now().Unix()),
		Zone: zone,
		IPNetworks: []upcloud.IPNetwork{
			{
				Address: addr,
				DHCP:    upcloud.True,
				Family:  upcloud.IPAddressFamilyIPv4,
			},
		},
	})
}

func waitForLoadBalancerToShutdown(ctx context.Context, rec *recorder.Recorder, svc *Service, lb *upcloud.LoadBalancer) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}

	const maxRetries int = 100
	// wait delete request
	for i := 0; i <= maxRetries; i++ {
		_, err := svc.GetLoadBalancer(ctx, &request.GetLoadBalancerRequest{UUID: lb.UUID})
		if err != nil {
			if svcErr, ok := err.(*upcloud.Problem); ok && svcErr.Status == http.StatusNotFound {
				return nil
			}
		}
		time.Sleep(2 * time.Second)
	}
	return errors.New("max retries reached while waiting for load balancer instance to shutdown")
}

func waitLoadBalancerToShutdown(ctx context.Context, svc *Service, lb *upcloud.LoadBalancer) error {
	const maxRetries int = 100
	// wait delete request
	for i := 0; i <= maxRetries; i++ {
		_, err := svc.GetLoadBalancer(ctx, &request.GetLoadBalancerRequest{UUID: lb.UUID})
		if err != nil {
			if svcErr, ok := err.(*upcloud.Problem); ok && svcErr.Status == http.StatusNotFound {
				return nil
			}
		}
		time.Sleep(2 * time.Second)
	}
	return errors.New("max retries reached while waiting for load balancer instance to shutdown")
}

func deleteLoadBalancer(ctx context.Context, svc *Service, lb *upcloud.LoadBalancer) error {
	if err := svc.DeleteLoadBalancer(ctx, &request.DeleteLoadBalancerRequest{UUID: lb.UUID}); err != nil {
		return err
	}

	if err := waitLoadBalancerToShutdown(ctx, svc, lb); err != nil {
		return fmt.Errorf("unable to shutdown LB '%s' (%s) (check dangling networks)", lb.UUID, lb.Name)
	}

	var errs []error
	if lb.NetworkUUID != "" {
		if err := svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{UUID: lb.NetworkUUID}); err != nil {
			errs = append(errs, err)
		}
	}
	if len(lb.Networks) > 0 {
		for _, n := range lb.Networks {
			if n.Type == upcloud.LoadBalancerNetworkTypePrivate && n.UUID != "" && lb.NetworkUUID != n.UUID {
				if err := svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{UUID: n.UUID}); err != nil {
					errs = append(errs, err)
				}
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", errs)
	}
	return nil
}

func createLoadBalancerPrivateNetwork(ctx context.Context, svc *Service, zone, addr string) (*upcloud.Network, error) {
	return svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
		Name: fmt.Sprintf("go-api-lb-%d", time.Now().Unix()),
		Zone: zone,
		IPNetworks: []upcloud.IPNetwork{
			{
				Address: addr,
				DHCP:    upcloud.True,
				Family:  upcloud.IPAddressFamilyIPv4,
			},
		},
	})
}
