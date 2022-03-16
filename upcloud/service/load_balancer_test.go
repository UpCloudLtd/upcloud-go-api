package service

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBalancer(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancer", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create Load Balancer
		lb, err := createLoadBalancerAndNetwork(svc, "fi-hel1", "172.16.0.0/24")
		require.NoError(t, err)
		defer cleanupLoadBalancer(t, svc, lb)
		t.Logf("Created load balancer: %s", lb.Name)

		// Modify Load Balancer
		t.Log("Modifying load balancer")

		newName := "new-name-for-lb"
		lb, err = svc.ModifyLoadBalancer(&request.ModifyLoadBalancerRequest{
			UUID: lb.UUID,
			Name: newName,
		})
		require.NoError(t, err)
		assert.Equal(t, newName, lb.Name)
		t.Logf("Modified load balancer with UUID: %s", lb.UUID)

		// Get Load Balancer
		t.Logf("Get load balancer: %s", newName)
		lb, err = svc.GetLoadBalancer(&request.GetLoadBalancerRequest{
			UUID: lb.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, newName, lb.Name)

		// Get Load Balancers
		t.Log("Get load balancers")
		lbs, err := svc.GetLoadBalancers(&request.GetLoadBalancersRequest{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(lbs), 1)
	})
}

func TestLoadBalancerBackend(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancerbackend", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancerAndNetwork(svc, "fi-hel2", "172.16.2.0/24")
		require.NoError(t, err)
		defer cleanupLoadBalancer(t, svc, lb)
		t.Logf("Created load balancer for testing LB backend CRUD: %s", lb.Name)

		backend, err := createLoadBalancerBackend(svc, lb.UUID)
		require.NoError(t, err)
		t.Logf("Created LB backend: %s", backend.Name)

		t.Logf("Modifying LB backend: %s", backend.Name)
		newName := "updatedName"
		backend, err = svc.ModifyLoadBalancerBackend(&request.ModifyLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        backend.Name,
			Backend: request.ModifyLoadBalancerBackend{
				Name: newName,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, backend.Name, newName)
		t.Logf("Modified LB backend, new name is: %s", backend.Name)

		t.Logf("Get LB backend: %s", newName)
		backend, err = svc.GetLoadBalancerBackend(&request.GetLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        newName,
		})
		require.NoError(t, err)
		assert.EqualValues(t, backend.Name, newName)
		assert.Len(t, backend.Members, 1)

		t.Log("Get LB backend list")
		backends, err := svc.GetLoadBalancerBackends(&request.GetLoadBalancerBackendsRequest{
			ServiceUUID: lb.UUID,
		})
		require.NoError(t, err)
		assert.Len(t, backends, 1)
		assert.EqualValues(t, backends[0].Name, newName)
		assert.Len(t, backends[0].Members, 1)

		t.Logf("Deleting LB backend: %s", backend.Name)
		err = svc.DeleteLoadBalancerBackend(&request.DeleteLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        backend.Name,
		})
		require.NoError(t, err)
	})
}

func TestLoadBalancerBackendMember(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancerbackendmember", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancerAndNetwork(svc, "nl-ams1", "172.16.3.0/24")
		require.NoError(t, err)
		defer cleanupLoadBalancer(t, svc, lb)
		t.Logf("Created load balancer for testing LB backend members CRUD: %s", lb.Name)

		backend, err := createLoadBalancerBackend(svc, lb.UUID)
		require.NoError(t, err)
		t.Logf("Created new backend %s for load balancer %s", backend.Name, lb.Name)

		name := "test_member"
		weight := 100
		maxSessions := 123
		enabled := true
		memberType := upcloud.LoadBalancerBackendMemberTypeStatic
		ip := "10.0.0.2"
		port := 80

		member, err := svc.CreateLoadBalancerBackendMember(&request.CreateLoadBalancerBackendMemberRequest{
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
		member, err = svc.ModifyLoadBalancerBackendMember(&request.ModifyLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Name:        member.Name,
			Member: request.LoadBalancerBackendMember{
				Name:        newName,
				Weight:      newWeight,
				MaxSessions: newMaxSessions,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Name, newName)
		assert.EqualValues(t, member.Weight, newWeight)
		assert.EqualValues(t, member.MaxSessions, newMaxSessions)
		t.Logf("Updated load balancers backend member name, weight and max sessions; new name: %s", member.Name)

		newIp := "231.231.231.231"
		newPort := 3003
		member, err = svc.ModifyLoadBalancerBackendMember(&request.ModifyLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Name:        member.Name,
			Member: request.LoadBalancerBackendMember{
				IP:   newIp,
				Port: newPort,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Type, memberType)
		assert.EqualValues(t, member.IP, newIp)
		assert.EqualValues(t, member.Port, newPort)
		t.Logf("Updated load balancers backend member type, ip and port: %s", member.Name)

		t.Logf("Get load balancer backend member: %s", newName)
		member, err = svc.GetLoadBalancerBackendMember(&request.GetLoadBalancerBackendMemberRequest{
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
		members, err := svc.GetLoadBalancerBackendMembers(&request.GetLoadBalancerBackendMembersRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
		})

		require.NoError(t, err)
		assert.Len(t, members, 2)

		err = svc.DeleteLoadBalancerBackendMember(&request.DeleteLoadBalancerBackendMemberRequest{
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

	record(t, "loadbalancerresolver", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancerAndNetwork(svc, "pl-waw1", "10.0.0.0/24")
		require.NoError(t, err)
		defer cleanupLoadBalancer(t, svc, lb)
		t.Logf("Created load balancer for testing LB resolvers CRUD: %s", lb.Name)

		name := "testname"
		nameServers := []string{"10.0.0.1", "10.0.0.2"}
		retries := 10
		timeout := 20
		timeoutRetry := 10
		cacheValid := 123
		cacheInvalid := 321

		t.Logf("Create resolver: %s", name)
		resolver, err := svc.CreateLoadBalancerResolver(&request.CreateLoadBalancerResolverRequest{
			ServiceUUID: lb.UUID,
			Resolver: request.LoadBalancerResolver{
				Name:         name,
				Nameservers:  nameServers,
				Retries:      retries,
				Timeout:      timeout,
				TimeoutRetry: timeoutRetry,
				CacheValid:   cacheValid,
				CacheInvalid: cacheInvalid},
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
		resolver, err = svc.ModifyLoadBalancerResolver(&request.ModifyLoadBalancerRevolverRequest{
			ServiceUUID: lb.UUID,
			Name:        resolver.Name,
			Resolver: request.LoadBalancerResolver{
				Name:        newName,
				Nameservers: newNameServers,
				Retries:     newRetries,
				Timeout:     newTimeout},
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
		resolver, err = svc.ModifyLoadBalancerResolver(&request.ModifyLoadBalancerRevolverRequest{
			ServiceUUID: lb.UUID,
			Name:        resolver.Name,
			Resolver: request.LoadBalancerResolver{
				TimeoutRetry: newTimeoutRetry,
				CacheValid:   newCacheValid,
				CacheInvalid: newCacheInvalid},
		})

		require.NoError(t, err)
		assert.EqualValues(t, resolver.TimeoutRetry, newTimeoutRetry)
		assert.EqualValues(t, resolver.CacheValid, newCacheValid)
		assert.EqualValues(t, resolver.CacheInvalid, newCacheInvalid)
		t.Logf("Modified timeout_retry, cache_valid and cache_invalid for resolver %s", resolver.Name)

		t.Logf("Get resolver: %s", resolver.Name)
		resolver, err = svc.GetLoadBalancerResolver(&request.GetLoadBalancerResolverRequest{
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
		resolvers, err := svc.GetLoadBalancerResolvers(&request.GetLoadBalancerResolversRequest{
			ServiceUUID: lb.UUID,
		})
		require.NoError(t, err)
		assert.Len(t, resolvers, 1)
		assert.EqualValues(t, resolvers[0].Name, resolver.Name)

		err = svc.DeleteLoadBalancerResolver(&request.DeleteLoadBalancerResolverRequest{
			ServiceUUID: lb.UUID,
			Name:        resolver.Name,
		})
		require.NoError(t, err)
		t.Logf("Deleted resolver %s for load balancer %s", resolver.Name, lb.Name)
	})
}

func TestGetLoadBalancerPlans(t *testing.T) {
	record(t, "getloadbalancerplans", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		plans, err := svc.GetLoadBalancerPlans(&request.GetLoadBalancerPlansRequest{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(plans), 2)
	})
}

func TestLoadBalancerFrontend(t *testing.T) {
	t.Parallel()

	record(t, "loadbalancerfrontend", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancerAndNetwork(svc, "de-fra1", "10.0.0.1/24")
		require.NoError(t, err)
		defer cleanupLoadBalancer(t, svc, lb)
		t.Logf("Created LB for testing frontends: %s", lb.Name)
		be, err := createLoadBalancerBackend(svc, lb.UUID)
		require.NoError(t, err)
		t.Logf("Created backend %s for testing LB frontends", be.Name)

		fe, err := svc.CreateLoadBalancerFrontend(&request.CreateLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Frontend: request.LoadBalancerFrontend{
				Name:           "fe-1",
				Mode:           upcloud.LoadBalancerModeHTTP,
				Port:           443,
				DefaultBackend: be.Name,
				Rules:          []request.LoadBalancerFrontendRule{},
				TLSConfigs:     []request.LoadBalancerTLSConfig{},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, "fe-1", fe.Name)
		t.Logf("Created frontend %s for load balancer %s", fe.Name, lb.Name)
		fe, err = svc.ModifyLoadBalancerFrontend(&request.ModifyLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Name:        fe.Name,
			Frontend: request.ModifyLoadBalancerFrontend{
				Name: "fe-2",
				Mode: upcloud.LoadBalancerModeTCP,
				Port: 80},
		})
		require.NoError(t, err)
		t.Logf("Modified frontend %s", fe.Name)
		fe, err = svc.GetLoadBalancerFrontend(&request.GetLoadBalancerFrontendRequest{
			ServiceUUID: lb.UUID,
			Name:        fe.Name,
		})
		require.NoError(t, err)
		assert.Equal(t, "fe-2", fe.Name)
		assert.Equal(t, upcloud.LoadBalancerModeTCP, fe.Mode)
		assert.Equal(t, 80, fe.Port)
		assert.Equal(t, be.Name, fe.DefaultBackend)

		fes, err := svc.GetLoadBalancerFrontends(&request.GetLoadBalancerFrontendsRequest{ServiceUUID: lb.UUID})
		require.NoError(t, err)
		assert.Equal(t, *fe, fes[0])

		assert.NoError(t,
			svc.DeleteLoadBalancerFrontend(&request.DeleteLoadBalancerFrontendRequest{
				ServiceUUID: lb.UUID,
				Name:        fe.Name,
			}))
	})
}

func cleanupLoadBalancer(t *testing.T, svc *Service, lb *upcloud.LoadBalancer) {
	t.Logf("Cleanup LB: %s", lb.Name)
	if err := deleteLoadBalancer(svc, lb); err != nil {
		t.Log(err)
	}
}
