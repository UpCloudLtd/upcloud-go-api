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
		lbDetails, err := createLoadBalancer(svc)
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
		lb, err := createLoadBalancer(svc)
		require.NoError(t, err)
		t.Logf("Created load balancer for testing LB backend CRUD: %s", lb.Name)

		backend, err := createLoadBalancerBackend(svc, lb.UUID)
		require.NoError(t, err)
		t.Logf("Created LB backend: %s", backend.Name)

		t.Logf("Modifying LB backend: %s", backend.Name)
		newName := "updatedName"
		backend, err = svc.ModifyLoadBalancerBackend(&request.ModifyLoadBalancerBackendRequest{
			ServiceUUID:    lb.UUID,
			BackendName:    backend.Name,
			NewBackendName: newName,
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
			UUID: "0abdab19-b11b-4f89-8f9f-961b524347b6",
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
		lb, err := createLoadBalancer(svc)
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
			ServiceUUID:       lb.UUID,
			BackendName:       backend.Name,
			MemberName:        name,
			MemberWeight:      weight,
			MemberMaxSessions: maxSessions,
			MemberEnabled:     enabled,
			MemberType:        memberType,
			MemberIP:          ip,
			MemberPort:        port,
			MemberServerUUID:  serverId,
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Name, name)
		assert.EqualValues(t, member.Weight, weight)
		assert.EqualValues(t, member.MaxSessions, maxSessions)
		assert.EqualValues(t, member.Enabled, enabled)
		assert.EqualValues(t, member.Type, memberType)
		assert.EqualValues(t, member.Ip, ip)
		assert.EqualValues(t, member.Port, port)
		t.Logf("Created new load balancer backend member: %s", member.Name)

		newName := "test_member_TURBO"
		newWeight := 50
		newMaxSessions := 321
		member, err = svc.ModifyLoadBalancerBackendMember(&request.ModifyLoadBalancerBackendMemberRequest{
			ServiceUUID:       lb.UUID,
			BackendName:       backend.Name,
			MemberName:        member.Name,
			NewMemberName:     newName,
			MemberWeight:      newWeight,
			MemberMaxSessions: newMaxSessions,
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
			MemberName:  member.Name,
			MemberType:  newType,
			MemberIP:    newIp,
			MemberPort:  newPort,
		})

		require.NoError(t, err)
		assert.EqualValues(t, member.Type, newType)
		assert.EqualValues(t, member.Ip, newIp)
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
		assert.EqualValues(t, member.Ip, "10.0.0.2")
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
		assert.EqualValues(t, member.Ip, "10.0.0.2")
		assert.EqualValues(t, member.Port, 80)
	})
}
