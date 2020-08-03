package request

import (
	"encoding/json"
	"encoding/xml"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
	"github.com/stretchr/testify/assert"
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
		Zone:             "fi-hel2",
		Title:            "Integration test server #1",
		Hostname:         "debian.example.com",
		PasswordDelivery: PasswordDeliveryNone,
		StorageDevices: []upcloud.CreateServerStorageDevice{
			{
				Action:  upcloud.CreateServerStorageDeviceActionClone,
				Storage: "01000000-0000-4000-8000-000030060200",
				Title:   "disk1",
				Size:    30,
				Tier:    upcloud.StorageTierMaxIOPS,
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
		LoginUser: &LoginUser{
			CreatePassword: "no",
			SSHKeys: []string{
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCWf2MmpHweXCNUcW91PWZR5UqOkydBr1Gi1xDI16IW4JndMYkH9OF0sWvPz03kfY6NbcHY0bed1Q8BpAC//WfLltuvjrzk33IoFJZ2Ai+4fVdkevkf7pBeSvzdXSyKAT+suHrp/2Qu5hewIUdDCp+znkwyypIJ/C2hDphwbLR3QquOfn6KyKMPZC4my8dFvLxESI0UqeripaBHUGcvNG2LL563hXmWzUu/cyqCpg5IBzpj/ketg8m1KBO7U0dimIAczuxfHk3kp9bwOFquWA2vSFNuVkr8oavk36pHkU88qojYNEy/zUTINE0w6CE/EbDkQVDZEGgDtAkq4jL+4MPV negge@palinski",
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDJfx4OmD8D6mnPA0BPk2DVlbggEkMvB2cecSttauZuaYX7Vju6PvG+kXrUbTvO09oLQMoNYAk3RinqQLXo9eF7bzZIsgB4ZmKGau84kOpYjguhimkKtZiVTKF53G2pbnpiZUN9wfy3xK2mt/MkacjZ1Tp7lAgRGTfWDoTfQa88kzOJGNPWXd12HIvFtd/1KoS9vm5O0nDLV+5zSBLxEYNDmBlIGu1Y3XXle5ygL1BhfGvqOQnv/TdRZcrOgVGWHADvwEid91/+IycLNMc37uP7TdS6vOihFBMytfmFXAqt4+3AzYNmyc+R392RorFzobZ1UuEFm3gUod2Wvj8pY8d/ negge@palinski",
			},
		},
	}

	expectedJSON := `
	{
      "server": {
        "hostname": "debian.example.com",
        "ip_addresses": {
          "ip_address": [
            {
              "access": "private",
              "family": "IPv4"
            },
            {
              "access": "public",
              "family": "IPv4"
            },
            {
              "access": "public",
              "family": "IPv6"
            }
          ]
        },
        "login_user": {
          "create_password": "no",
          "ssh_keys": {
			  "ssh_key": [
                "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCWf2MmpHweXCNUcW91PWZR5UqOkydBr1Gi1xDI16IW4JndMYkH9OF0sWvPz03kfY6NbcHY0bed1Q8BpAC//WfLltuvjrzk33IoFJZ2Ai+4fVdkevkf7pBeSvzdXSyKAT+suHrp/2Qu5hewIUdDCp+znkwyypIJ/C2hDphwbLR3QquOfn6KyKMPZC4my8dFvLxESI0UqeripaBHUGcvNG2LL563hXmWzUu/cyqCpg5IBzpj/ketg8m1KBO7U0dimIAczuxfHk3kp9bwOFquWA2vSFNuVkr8oavk36pHkU88qojYNEy/zUTINE0w6CE/EbDkQVDZEGgDtAkq4jL+4MPV negge@palinski",
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDJfx4OmD8D6mnPA0BPk2DVlbggEkMvB2cecSttauZuaYX7Vju6PvG+kXrUbTvO09oLQMoNYAk3RinqQLXo9eF7bzZIsgB4ZmKGau84kOpYjguhimkKtZiVTKF53G2pbnpiZUN9wfy3xK2mt/MkacjZ1Tp7lAgRGTfWDoTfQa88kzOJGNPWXd12HIvFtd/1KoS9vm5O0nDLV+5zSBLxEYNDmBlIGu1Y3XXle5ygL1BhfGvqOQnv/TdRZcrOgVGWHADvwEid91/+IycLNMc37uP7TdS6vOihFBMytfmFXAqt4+3AzYNmyc+R392RorFzobZ1UuEFm3gUod2Wvj8pY8d/ negge@palinski"
			  ]
		  }
        },
        "password_delivery": "none",
        "storage_devices": {
          "storage_device": [
            {
              "action": "clone",
              "storage": "01000000-0000-4000-8000-000030060200",
              "title": "disk1",
              "size": 30,
              "tier": "maxiops"
            }
          ]
        },
        "title": "Integration test server #1",
        "zone": "fi-hel2"
      }
    }
	`
	actualJSON, err := json.Marshal(&request)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
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

// TestDeleteServerAndStoragesRequest tests that DeleteServerAndStoragesRequest objects behave correctly
func TestDeleteServerAndStoragesRequest(t *testing.T) {
	request := DeleteServerAndStoragesRequest{
		UUID: "foo",
	}

	assert.Equal(t, "/server/foo/?storages=1", request.RequestURL())
}

// TestTagServerRequest tests that TestTagServer behaves correctly
func TestTagServerRequest(t *testing.T) {
	// Test with multiple tags
	request := TagServerRequest{
		UUID: "foo",
		Tags: []string{
			"tag1",
			"tag2",
			"tag with spaces",
		},
	}

	assert.Equal(t, "/server/foo/tag/tag1,tag2,tag with spaces", request.RequestURL())

	// Test with single tag
	request = TagServerRequest{
		UUID: "foo",
		Tags: []string{
			"tag1",
		},
	}

	assert.Equal(t, "/server/foo/tag/tag1", request.RequestURL())
}

func TestUntagServerRequest(t *testing.T) {
	// Test with multiple tags
	request := UntagServerRequest{
		UUID: "foo",
		Tags: []string{
			"tag1",
			"tag2",
			"tag with spaces",
		},
	}

	assert.Equal(t, "/server/foo/untag/tag1,tag2,tag with spaces", request.RequestURL())

	// Test with single tag
	request = UntagServerRequest{
		UUID: "foo",
		Tags: []string{
			"tag1",
		},
	}

	assert.Equal(t, "/server/foo/untag/tag1", request.RequestURL())
}
