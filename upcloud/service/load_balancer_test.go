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
