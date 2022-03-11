package service

import (
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBalancerCRUD(t *testing.T) {
	t.Parallel()

	record(t, "createmodifydeleteloadbalancer", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create Load Balancer
		lbDetails, err := createLoadBalancer(svc, "032d4c7f-61b5-4ea9-a2d6-d2357c3c9a88", "es-mad1")
		require.NoError(t, err)
		t.Logf("Created load balancer: %s", lbDetails.Name)

		// Modify Load Balancer
		t.Log("Modifying load balancer")

		newName := "new-name-for-lb"
		lbDetails, err = svc.ModifyLoadBalancer(&request.ModifyLoadBalancerRequest{
			UUID: lbDetails.UUID,
			Name: newName,
		})
		require.NoError(t, err)
		assert.Equal(t, newName, lbDetails.Name)
		t.Logf("Modified load balancer with UUID: %s", lbDetails.UUID)

		// Delete Load Balancer
		t.Log("Deleting load balancer")
		err = deleteLoadBalancer(svc, lbDetails.UUID)
		require.NoError(t, err)
		t.Logf("Deleted load balancer with UUID: %s", lbDetails.UUID)
	})
}

func TestGetLoadBalancerDetails(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancer", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		lbDetails, err := svc.GetLoadBalancerDetails(&request.GetLoadBalancerDetailsRequest{
			UUID: "0a70595d-f2ef-4933-8d2f-7b1931a62bc4",
		})
		require.NoError(t, err)
		fmt.Printf("lb: %v", lbDetails)
	})
}

func TestGetLoadBalancers(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancers", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		lbs, err := svc.GetLoadBalancers(&request.GetLoadBalancersRequest{})
		require.NoError(t, err)
		fmt.Printf("lbs: %v", lbs)
	})
}

func TestLoadBalancerBackendCRUD(t *testing.T) {
	t.Parallel()

	record(t, "createmodifydeleteloadbalancerbackend", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancer(svc, "032d4c7f-61b5-4ea9-a2d6-d2357c3c9a88", "es-mad1")
		require.NoError(t, err)
		t.Logf("Created load balancer for testing LB backend CRUD: %s", lb.Name)

		backend, err := createLoadBalancerBackend(svc, lb.UUID)
		require.NoError(t, err)
		t.Logf("Created LB backend: %s", backend.Name)

		t.Logf("Modifying LB backend: %s", backend.Name)
		newName := "updatedName"
		backend, err = svc.ModifyLoadBalancerBackend(&request.ModifyLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			Name:        backend.Name,
			Payload: request.ModifyLoadBalancerBackend{
				Name: newName,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, backend.Name, newName)
		t.Logf("Modified LB backend, new name is: %s", backend.Name)

		t.Logf("Deleting LB backend: %s", backend.Name)
		err = svc.DeleteLoadBalancerBackend(&request.DeleteLoadBalancerBackendRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
		})
		require.NoError(t, err)
	})
}

func TestGetLoadBalancerBackendDetails(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancerbackenddetails", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		backend, err := svc.GetLoadBalancerBackendDetails(&request.GetLoadBalancerBackendDetailsRequest{
			ServiceUUID: "0abdab19-b11b-4f89-8f9f-961b524347b6",
			BackendName: "updatedName",
		})
		require.NoError(t, err)
		assert.EqualValues(t, backend.Name, "updatedName")
		assert.Len(t, backend.Members, 0)
	})
}

func TestGetLoadBalancerBackends(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancerbackends", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		backends, err := svc.GetLoadBalancerBackends(&request.GetLoadBalancerBackendsRequest{
			ServiceUUID: "0abdab19-b11b-4f89-8f9f-961b524347b6",
		})
		require.NoError(t, err)
		assert.Len(t, backends, 1)
		assert.EqualValues(t, backends[0].Name, "updatedName")
		assert.Len(t, backends[0].Members, 0)
	})
}

func TestLoadBalancerBackendMemberCRUD(t *testing.T) {
	t.Parallel()

	record(t, "createmodifydeleteloadbalancerbackendmember", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancer(svc, "032d4c7f-61b5-4ea9-a2d6-d2357c3c9a88", "es-mad1")
		require.NoError(t, err)
		t.Logf("Created load balancer for testing LB backend members CRUD: %s", lb.Name)

		backend, err := createLoadBalancerBackend(svc, lb.UUID)
		require.NoError(t, err)
		t.Logf("Created new backend %s for load balancer %s", backend.Name, lb.Name)

		name := "test_member"
		weight := 100
		maxSessions := 123
		enabled := true
		memberType := "static"
		ip := "10.0.0.2"
		port := 80
		serverId := "0050febf-b881-4db1-85ce-4c92776a47e2"

		member, err := svc.CreateLoadBalancerBackendMember(&request.CreateLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Payload: request.CreateLoadBalancerBackendMember{
				Name:        name,
				Weight:      weight,
				MaxSessions: maxSessions,
				Enabled:     enabled,
				Type:        memberType,
				IP:          ip,
				Port:        port,
				ServerUUID:  serverId,
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
			Payload: request.ModifyLoadBalancerBackendMember{
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

		newType := "dynamic"
		newIp := "231.231.231.231"
		newPort := 3003
		member, err = svc.ModifyLoadBalancerBackendMember(&request.ModifyLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			Name:        member.Name,
			Payload: request.ModifyLoadBalancerBackendMember{
				Type: newType,
				IP:   newIp,
				Port: newPort,
			},
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Type, newType)
		assert.EqualValues(t, member.IP, newIp)
		assert.EqualValues(t, member.Port, newPort)
		t.Logf("Updated load balancers backend member type, ip and port: %s", member.Name)

		err = svc.DeleteLoadBalancerBackendMember(&request.DeleteLoadBalancerBackendMemberRequest{
			ServiceUUID: lb.UUID,
			BackendName: backend.Name,
			MemberName:  member.Name,
		})
		require.NoError(t, err)
		t.Logf("Deleted load balancer backend member: %s", member.Name)
	})
}

func TestGetLoadBalancerBackendMemberDetails(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancerbackendmemberdetails", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		member, err := svc.GetLoadBalancerBackendMemberDetails(&request.GetLoadBalancerBackendMemberDetailsRequest{
			ServiceUUID: "0a0cefa5-a6ef-41d1-b652-3106164f4164",
			BackendName: "go-test-lb-backend",
			MemberName:  "backend-1646814491350",
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Name, "backend-1646814491350")
		assert.EqualValues(t, member.Enabled, true)
		assert.EqualValues(t, member.Weight, 100)
		assert.EqualValues(t, member.MaxSessions, 1000)
		assert.EqualValues(t, member.Type, "static")
		assert.EqualValues(t, member.IP, "10.0.0.2")
		assert.EqualValues(t, member.Port, 80)
	})
}

func TestGetLoadBalancerBackendMembers(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancerbackendmembers", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		members, err := svc.GetLoadBalancerBackendMembers(&request.GetLoadBalancerBackendMembersRequest{
			ServiceUUID: "0a0cefa5-a6ef-41d1-b652-3106164f4164",
			BackendName: "go-test-lb-backend",
		})

		require.NoError(t, err)
		assert.Len(t, members, 1)

		member := members[0]
		assert.EqualValues(t, member.Name, "backend-1646814491350")
		assert.EqualValues(t, member.Enabled, true)
		assert.EqualValues(t, member.Weight, 100)
		assert.EqualValues(t, member.MaxSessions, 1000)
		assert.EqualValues(t, member.Type, "static")
		assert.EqualValues(t, member.IP, "10.0.0.2")
		assert.EqualValues(t, member.Port, 80)
	})
}

func TestLoadBalancerResolverCRUD(t *testing.T) {
	t.Parallel()

	record(t, "createmodifydeleteloadbalancerresolver", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		lb, err := createLoadBalancer(svc, "03716b3c-663f-46ab-a7a4-84221525030b", "pl-waw1")
		require.NoError(t, err)
		t.Logf("Created load balancer for testing LB resolvers CRUD: %s", lb.Name)

		name := "testname"
		nameServers := []string{"10.0.0.1", "10.0.0.2"}
		retries := 10
		timeout := 20
		timeoutRetry := 10
		cacheValid := 123
		cacheInvalid := 321

		resolver, err := svc.CreateLoadBalancerResolver(&request.CreateLoadBalancerResolverRequest{
			ServiceUUID:  lb.UUID,
			Name:         name,
			Nameservers:  nameServers,
			Retries:      retries,
			Timeout:      timeout,
			TimeoutRetry: timeoutRetry,
			CacheValid:   cacheValid,
			CacheInvalid: cacheInvalid,
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
		resolver, err = svc.ModifyLoadBalancerResolver(&request.ModifyLoadBalancerRevolverRequest{
			ServiceUUID:     lb.UUID,
			ResolverName:    resolver.Name,
			NewResolverName: newName,
			Nameservers:     newNameServers,
			Retries:         newRetries,
			Timeout:         newTimeout,
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
			ServiceUUID:  lb.UUID,
			ResolverName: resolver.Name,
			TimeoutRetry: newTimeoutRetry,
			CacheValid:   newCacheValid,
			CacheInvalid: newCacheInvalid,
		})

		require.NoError(t, err)
		assert.EqualValues(t, resolver.TimeoutRetry, newTimeoutRetry)
		assert.EqualValues(t, resolver.CacheValid, newCacheValid)
		assert.EqualValues(t, resolver.CacheInvalid, newCacheInvalid)
		t.Logf("Modified timeout_retry, cache_valid and cache_invalid for resolver %s", resolver.Name)

		err = svc.DeleteLoadBalancerResolver(&request.DeleteLoadBalancerResolverRequest{
			ServiceUUID:  lb.UUID,
			ResolverName: resolver.Name,
		})
		require.NoError(t, err)
		t.Logf("Deleted resolver %s for load balancer %s", resolver.Name, lb.Name)
	})
}

func TestGetLoadBalancerResolvers(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancerresolvers", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		resolvers, err := svc.GetLoadBalancerResolvers(&request.GetLoadBalancerResolversRequest{
			ServiceUUID: "0aceb3dc-e79e-4664-8224-627410be0e8f",
		})

		require.NoError(t, err)
		assert.Len(t, resolvers, 1)

		resolver := resolvers[0]
		assert.EqualValues(t, resolver.Name, "testname")
		assert.EqualValues(t, resolver.Retries, 10)
		assert.EqualValues(t, resolver.Timeout, 20)
		assert.EqualValues(t, resolver.TimeoutRetry, 10)
		assert.EqualValues(t, resolver.CacheValid, 123)
		assert.EqualValues(t, resolver.CacheInvalid, 321)
		assert.Len(t, resolver.Nameservers, 2)
		assert.Contains(t, resolver.Nameservers, "10.0.0.1")
		assert.Contains(t, resolver.Nameservers, "10.0.0.2")
	})
}

func TestGetLoadBalancerResolverDetails(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancerresolverdetails", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		resolver, err := svc.GetLoadBalancerResolverDetails(&request.GetLoadBalancerResolverDetailsRequest{
			ServiceUUID:  "0aceb3dc-e79e-4664-8224-627410be0e8f",
			ResolverName: "testname",
		})

		require.NoError(t, err)
		assert.EqualValues(t, resolver.Name, "testname")
		assert.EqualValues(t, resolver.Retries, 10)
		assert.EqualValues(t, resolver.Timeout, 20)
		assert.EqualValues(t, resolver.TimeoutRetry, 10)
		assert.EqualValues(t, resolver.CacheValid, 123)
		assert.EqualValues(t, resolver.CacheInvalid, 321)
		assert.Len(t, resolver.Nameservers, 2)
		assert.Contains(t, resolver.Nameservers, "10.0.0.1")
		assert.Contains(t, resolver.Nameservers, "10.0.0.2")
	})
}
