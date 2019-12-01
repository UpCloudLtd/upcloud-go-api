package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetNetworksRequest tests that GetNetworksRequest behaves correctly
func TestGetNetworksRequest(t *testing.T) {
	request := GetNetworksRequest{}

	assert.Equal(t, "/network", request.RequestURL())
}

// TestGetNetworksInZoneRequest tests that GetNetworksInZoneRequest behaves correctly
func TestGetNetworksInZoneRequest(t *testing.T) {
	request := GetNetworksInZoneRequest{Zone: "nl-ams1"}

	assert.Equal(t, "/network/?zone=nl-ams1", request.RequestURL())
}

// TestCreateSDNPrivateNetworkRequest tests that CreateSDNPrivateNetworkRequest behaves correctly
func TestCreateSDNPrivateNetworkRequest(t *testing.T) {
	request := CreateSDNPrivateNetworkRequest{}

	assert.Equal(t, "/network", request.RequestURL())
}

// TestGetNetworkDetailsRequest tests that GetNetworkDetailsRequest behaves correctly
func TestGetNetworkDetailsRequest(t *testing.T) {
	request := GetNetworkDetailsRequest{UUID: "test"}

	assert.Equal(t, "/network/test", request.RequestURL())
}

// TestModifyNetworkDetailsRequest tests that ModifyNetworkDetailsRequest behaves correctly
func TestModifyNetworkDetailsRequest(t *testing.T) {
	request := ModifyNetworkDetailsRequest{UUID: "test"}

	assert.Equal(t, "/network/test", request.RequestURL())
}

// TestDeleteNetworksRequest tests that DeleteNetworksRequest behaves correctly
func TestDeleteNetworkRequest(t *testing.T) {
	request := DeleteNetworkRequest{UUID: "test"}

	assert.Equal(t, "/network/test", request.RequestURL())
}

// TestCreateNetworkInterface tests that CreateNetworkInterface behaves correctly
func TestCreateNetworkInterface(t *testing.T) {
	request := CreateNetworkInterfaceRequest{ServerUUID: "test"}

	assert.Equal(t, "/server/test/networking/interface", request.RequestURL())
}

// TestListServerNetworks tests that ListServerNetworks behaves correctly
func TestListServerNetworks(t *testing.T) {
	request := ListServerNetworks{ServerUUID: "test"}

	assert.Equal(t, "/server/test/networking", request.RequestURL())
}

// TestModifyNetworkInterface tests that ModifyNetworkInterface behaves correctly
func TestModifyNetworkInterface(t *testing.T) {
	request := ModifyNetworkInterfaceRequest{ServerUUID: "test", Index: 1}

	assert.Equal(t, "/server/test/networking/interface/1", request.RequestURL())
}

// TestDeleteNetworkInterface tests that DeleteNetworkInterface behaves correctly
func TestDeleteNetworkInterface(t *testing.T) {
	request := DeleteNetworkInterfaceRequest{ServerUUID: "test", Index: 1}

	assert.Equal(t, "/server/test/networking/interface/1", request.RequestURL())
}
