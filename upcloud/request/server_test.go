package request

import (
	"encoding/xml"
	"github.com/jalle19/upcloud-go-sdk/upcloud"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestGetServerDetailsRequest tests that GetServerDetailsRequest objects behave correctly
func TestGetServerDetailsRequest(t *testing.T) {
	request := GetServerDetailsRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/server/foo", request.RequestURL())
}

// TestCreateServerRequest tests that CreateServerRequest objects behave correctly
func TestCreateServerRequest(t *testing.T) {
	request := CreateServerRequest{
		Zone:             "fi-hel1",
		Title:            "Integration test server #1",
		Hostname:         "debian.example.com",
		PasswordDelivery: PasswordDeliveryNone,
		StorageDevices: []upcloud.CreateServerStorageDevice{
			{
				Action:  upcloud.CreateServerStorageDeviceActionClone,
				Storage: "01000000-0000-4000-8000-000030060200",
				Title:   "disk1",
				Size:    30,
				Tier:    upcloud.CreateStorageDeviceTierMaxIOPS,
			},
		},
		IPAddresses: []CreateServerIPAddress{
			{
				Access: upcloud.IPAddressAccessPrivate,
				Family: upcloud.IPAddressFamilyIPv4,
			},
			{
				Access: upcloud.IPAddressAccessPublic,
				Family: upcloud.IPAddressFamilyIPv4,
			},
			{
				Access: upcloud.IPAddressAccessPublic,
				Family: upcloud.IPAddressFamilyIPv6,
			},
		},
	}

	expectedXML := "<server><hostname>debian.example.com</hostname><ip_addresses><ip_address><access>private</access><family>IPv4</family></ip_address><ip_address><access>public</access><family>IPv4</family></ip_address><ip_address><access>public</access><family>IPv6</family></ip_address></ip_addresses><password_delivery>none</password_delivery><storage_devices><storage_device><action>clone</action><storage>01000000-0000-4000-8000-000030060200</storage><title>disk1</title><size>30</size><tier>maxiops</tier></storage_device></storage_devices><title>Integration test server #1</title><zone>fi-hel1</zone></server>"
	actualXML, err := xml.Marshal(&request)
	assert.Nil(t, err)
	assert.Equal(t, expectedXML, string(actualXML))
	assert.Equal(t, "/server", request.RequestURL())
}

// TestStartServerRequest tests that StartServerRequest objects behave correctly
func TestStartServerRequest(t *testing.T) {
	request := StartServerRequest{
		UUID:    "foo",
		Timeout: time.Minute * 5,
	}

	assert.Equal(t, "/server/foo/start", request.RequestURL())
}

// TestStopServerRequest tests that StopServerRequest objects behave correctly
func TestStopServerRequest(t *testing.T) {
	request := StopServerRequest{
		UUID:     "foo",
		StopType: ServerStopTypeHard,
		Timeout:  time.Minute * 5,
	}

	expectedXML := "<stop_server><timeout>300</timeout><stop_type>hard</stop_type></stop_server>"
	actualXML, err := xml.Marshal(&request)
	assert.Nil(t, err)
	assert.Equal(t, expectedXML, string(actualXML))
	assert.Equal(t, "/server/foo/stop", request.RequestURL())
}

// TestRestartServerRequest tests that RestartServerRequest objects behave correctly
func TestRestartServerRequest(t *testing.T) {
	request := RestartServerRequest{
		UUID:          "foo",
		Timeout:       time.Minute * 5,
		StopType:      ServerStopTypeSoft,
		TimeoutAction: RestartTimeoutActionDestroy,
	}

	expectedXML := "<restart_server><timeout>300</timeout><stop_type>soft</stop_type><timeout_action>destroy</timeout_action></restart_server>"
	actualXML, err := xml.Marshal(&request)
	assert.Nil(t, err)
	assert.Equal(t, expectedXML, string(actualXML))
	assert.Equal(t, "/server/foo/restart", request.RequestURL())
}

// TestModifyServerRequest tests that ModifyServerRequest objects behave correctly
func TestModifyServerRequest(t *testing.T) {
	request := ModifyServerRequest{
		UUID:  "foo",
		Title: "Modified server",
	}

	expectedXML := "<server><title>Modified server</title></server>"
	actualXML, err := xml.Marshal(&request)
	assert.Nil(t, err)
	assert.Equal(t, expectedXML, string(actualXML))
	assert.Equal(t, "/server/foo", request.RequestURL())
}

// TestDeleteServerRequest tests that DeleteServerRequest objects behave correctly
func TestDeleteServerRequest(t *testing.T) {
	request := DeleteServerRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/server/foo", request.RequestURL())
}
