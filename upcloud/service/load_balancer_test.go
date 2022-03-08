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
			UUID: lbDetails.Uuid,
			Name: newName,
		})
		require.NoError(t, err)
		assert.Equal(t, newName, lbDetails.Name)
		t.Logf("Modified load balancer with UUID: %s", lbDetails.Uuid)

		// Delete Load Balancer
		t.Log("Deleting load balancer")
		err = deleteLoadBalancer(svc, lbDetails.Uuid)
		require.NoError(t, err)
		t.Logf("Deleted load balancer with UUID: %s", lbDetails.Uuid)
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

		backend, err := createLoadBalancerBackend(svc, lb.Uuid)
		require.NoError(t, err)
		t.Logf("Created LB backend: %s", backend.Name)

		t.Logf("Modifying LB backend: %s", backend.Name)
		newName := "updatedName"
		backend, err = svc.ModifyLoadBalancerBackend(&request.ModifyLoadBalancerBackendRequest{
			ServiceUUID:    lb.Uuid,
			BackendName:    backend.Name,
			NewBackendName: newName,
		})

		require.NoError(t, err)
		assert.EqualValues(t, backend.Name, newName)
		t.Logf("Modified LB backend, new name is: %s", backend.Name)

		t.Logf("Deleting LB backend: %s", backend.Name)
		err = svc.DeleteLoadBalancerBackend(&request.DeleteLoadBalancerBackendRequest{
			ServiceUUID: lb.Uuid,
			BackendName: backend.Name,
		})
		require.NoError(t, err)
	})
}

func TestGetLoadBalancerBackendDetails(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancerbackenddetails", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		backend, err := svc.GetLoadBalancerBackendDetails(&request.GetLoadBalancerBackendDetailsRequest{
			ServiceUUID: "0a988541-06e0-4252-b483-3ff9fdfdb5ae",
			BackendName: "updatedName",
		})
		require.NoError(t, err)
		assert.EqualValues(t, backend.Name, "updatedName")
		assert.Len(t, backend.Members, 1)

		mem := backend.Members[0]
		assert.EqualValues(t, mem.Port, 8000)
		assert.EqualValues(t, mem.Enabled, true)
		assert.EqualValues(t, mem.Ip, "196.123.123.123")
		assert.EqualValues(t, mem.MaxSessions, 1000)
		assert.EqualValues(t, mem.Weight, 100)
		assert.EqualValues(t, mem.Type, "dynamic")
		assert.EqualValues(t, mem.Name, "default-lb-backend-member")
	})
}

func TestGetLoadBalancerBackends(t *testing.T) {
	t.Parallel()

	record(t, "getloadbalancerbackends", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		backends, err := svc.GetLoadBalancerBackends(&request.GetLoadBalancerBackendsRequest{
			UUID: "0a988541-06e0-4252-b483-3ff9fdfdb5ae",
		})
		require.NoError(t, err)
		assert.Len(t, backends, 1)
		assert.EqualValues(t, backends[0].Name, "updatedName")
		assert.Len(t, backends[0].Members, 1)

		mem := backends[0].Members[0]
		assert.EqualValues(t, mem.Port, 8000)
		assert.EqualValues(t, mem.Enabled, true)
		assert.EqualValues(t, mem.Ip, "196.123.123.123")
		assert.EqualValues(t, mem.MaxSessions, 1000)
		assert.EqualValues(t, mem.Weight, 100)
		assert.EqualValues(t, mem.Type, "dynamic")
		assert.EqualValues(t, mem.Name, "default-lb-backend-member")
	})
}
