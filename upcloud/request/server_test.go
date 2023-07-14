package request

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/stretchr/testify/assert"
)

// TestGetServersWithFiltersRequest tests that GetServersWithFiltersRequest objects behave correctly
func TestGetServersWithFiltersRequest(t *testing.T) {
	request := GetServersWithFiltersRequest{
		Filters: []ServerFilter{
			FilterLabelKey{"onlyKey1"},
			FilterLabelKey{"onlyKey2"},
			FilterLabel{Label: upcloud.Label{
				Key:   "pairKey1",
				Value: "pairValue1",
			}},
			FilterLabel{Label: upcloud.Label{
				Key:   "pairKey2",
				Value: "pairValue2",
			}},
		},
	}

	assert.Equal(
		t,
		"/server/?label=onlyKey1&label=onlyKey2&label=pairKey1%3DpairValue1&label=pairKey2%3DpairValue2",
		request.RequestURL(),
	)
}

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
		ServerGroup:      "test-group",
		StorageDevices: []CreateServerStorageDevice{
			{
				Action:  CreateServerStorageDeviceActionClone,
				Storage: "01000000-0000-4000-8000-000020060100",
				Title:   "disk1",
				Size:    10,
				Tier:    upcloud.StorageTierMaxIOPS,
			},
		},
		SimpleBackup: "0430,monthlies",
		Labels: &upcloud.LabelSlice{
			upcloud.Label{
				Key:   "managedBy",
				Value: "upcloud-go-sdk-unit-test",
			},
		},
		Metadata: upcloud.True,
		NICModel: upcloud.NICModelVirtio,
		Networking: &CreateServerNetworking{
			Interfaces: []CreateServerInterface{
				{
					IPAddresses: []CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
					Type: upcloud.IPAddressAccessPublic,
				},
				{
					IPAddresses: []CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv4,
						},
					},
					Type: upcloud.IPAddressAccessUtility,
				},
				{
					IPAddresses: []CreateServerIPAddress{
						{
							Family: upcloud.IPAddressFamilyIPv6,
						},
					},
					Type: upcloud.IPAddressAccessPublic,
				},
			},
		},
		LoginUser: &LoginUser{
			CreatePassword: "no",
			SSHKeys: []string{
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCWf2MmpHweXCNUcW91PWZR5UqOkydBr1Gi1xDI16IW4JndMYkH9OF0sWvPz03kfY6NbcHY0bed1Q8BpAC//WfLltuvjrzk33IoFJZ2Ai+4fVdkevkf7pBeSvzdXSyKAT+suHrp/2Qu5hewIUdDCp+znkwyypIJ/C2hDphwbLR3QquOfn6KyKMPZC4my8dFvLxESI0UqeripaBHUGcvNG2LL563hXmWzUu/cyqCpg5IBzpj/ketg8m1KBO7U0dimIAczuxfHk3kp9bwOFquWA2vSFNuVkr8oavk36pHkU88qojYNEy/zUTINE0w6CE/EbDkQVDZEGgDtAkq4jL+4MPV negge@palinski",
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDJfx4OmD8D6mnPA0BPk2DVlbggEkMvB2cecSttauZuaYX7Vju6PvG+kXrUbTvO09oLQMoNYAk3RinqQLXo9eF7bzZIsgB4ZmKGau84kOpYjguhimkKtZiVTKF53G2pbnpiZUN9wfy3xK2mt/MkacjZ1Tp7lAgRGTfWDoTfQa88kzOJGNPWXd12HIvFtd/1KoS9vm5O0nDLV+5zSBLxEYNDmBlIGu1Y3XXle5ygL1BhfGvqOQnv/TdRZcrOgVGWHADvwEid91/+IycLNMc37uP7TdS6vOihFBMytfmFXAqt4+3AzYNmyc+R392RorFzobZ1UuEFm3gUod2Wvj8pY8d/ negge@palinski",
			},
		},
		RemoteAccessEnabled:  upcloud.True,
		RemoteAccessType:     upcloud.RemoteAccessTypeVNC,
		RemoteAccessPassword: "abcdefgh",
	}

	expectedJSON := `
	{
		"server":{
			 "hostname":"debian.example.com",
			 "login_user":{
					"create_password":"no",
					"ssh_keys":{
						 "ssh_key":[
								"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCWf2MmpHweXCNUcW91PWZR5UqOkydBr1Gi1xDI16IW4JndMYkH9OF0sWvPz03kfY6NbcHY0bed1Q8BpAC//WfLltuvjrzk33IoFJZ2Ai+4fVdkevkf7pBeSvzdXSyKAT+suHrp/2Qu5hewIUdDCp+znkwyypIJ/C2hDphwbLR3QquOfn6KyKMPZC4my8dFvLxESI0UqeripaBHUGcvNG2LL563hXmWzUu/cyqCpg5IBzpj/ketg8m1KBO7U0dimIAczuxfHk3kp9bwOFquWA2vSFNuVkr8oavk36pHkU88qojYNEy/zUTINE0w6CE/EbDkQVDZEGgDtAkq4jL+4MPV negge@palinski",
								"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDJfx4OmD8D6mnPA0BPk2DVlbggEkMvB2cecSttauZuaYX7Vju6PvG+kXrUbTvO09oLQMoNYAk3RinqQLXo9eF7bzZIsgB4ZmKGau84kOpYjguhimkKtZiVTKF53G2pbnpiZUN9wfy3xK2mt/MkacjZ1Tp7lAgRGTfWDoTfQa88kzOJGNPWXd12HIvFtd/1KoS9vm5O0nDLV+5zSBLxEYNDmBlIGu1Y3XXle5ygL1BhfGvqOQnv/TdRZcrOgVGWHADvwEid91/+IycLNMc37uP7TdS6vOihFBMytfmFXAqt4+3AzYNmyc+R392RorFzobZ1UuEFm3gUod2Wvj8pY8d/ negge@palinski"
						 ]
					}
			 },
			 "password_delivery":"none",
			 "server_group": "test-group",
			 "storage_devices":{
					"storage_device":[
						 {
								"action":"clone",
								"storage":"01000000-0000-4000-8000-000020060100",
								"title":"disk1",
								"size":10,
								"tier":"maxiops"
						 }
					]
			 },
			 "simple_backup":"0430,monthlies",
			 "labels":{
					"label":[
						 {
								"key":"managedBy",
								"value":"upcloud-go-sdk-unit-test"
						 }
					]
			 },
			 "metadata":"yes",
			 "nic_model":"virtio",
			 "networking":{
					"interfaces":{
						 "interface":[
								{
									 "ip_addresses":{
											"ip_address":[
												 {
														"family":"IPv4"
												 }
											]
									 },
									 "type":"public"
								},
								{
									 "ip_addresses":{
											"ip_address":[
												 {
														"family":"IPv4"
												 }
											]
									 },
									 "type":"utility"
								},
								{
									 "ip_addresses":{
											"ip_address":[
												 {
														"family":"IPv6"
												 }
											]
									 },
									 "type":"public"
								}
						 ]
					}
			 },
			 "title":"Integration test server #1",
			 "remote_access_enabled":"yes",
			 "remote_access_type":"vnc",
			 "remote_access_password":"abcdefgh",
			 "zone":"fi-hel2"
		}
 }
	`
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server", request.RequestURL())
}

// TestStartServerRequest_OmitValues tests that StartServerRequest objects behave correctly
// when Host and AvoidHost are not specified
func TestStartServerRequest_OmitValues(t *testing.T) {
	request := StartServerRequest{
		UUID: "foo",
	}

	expectedJSON := `
	  {
		  "server": {}
	  }
	`

	actualJSON, err := json.Marshal(request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))

	assert.Equal(t, "/server/foo/start", request.RequestURL())
}

// TestStartServerRequest_WithValues tests that StartServerRequest objects behave correctly
// when Host and AvoidHost are specified
func TestStartServerRequest_WithValues(t *testing.T) {
	request := StartServerRequest{
		UUID:      "foo",
		Host:      1010,
		AvoidHost: 1101,
	}

	expectedJSON := `
	  {
		  "server": {
			  "host": 1010,
			  "avoid_host": 1101 
		  }
	  }
	`

	actualJSON, err := json.Marshal(request)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))

	assert.Equal(t, "/server/foo/start", request.RequestURL())
}

// TestStopServerRequest tests that StopServerRequest objects behave correctly
func TestStopServerRequest(t *testing.T) {
	request := StopServerRequest{
		UUID:     "foo",
		StopType: ServerStopTypeSoft,
		Timeout:  time.Minute * 5,
	}

	expectedJSON := `
	  {
		"stop_server": {
		  "stop_type": "soft",
		  "timeout": "300"
		}
	  }
	`
	actualJSON, err := json.MarshalIndent(&request, "", "    ")
	assert.Nil(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/foo/stop", request.RequestURL())
}

// TestRestartServerRequest tests that RestartServerRequest objects behave correctly
func TestRestartServerRequest(t *testing.T) {
	request := RestartServerRequest{
		UUID:          "foo",
		Timeout:       time.Minute * 5,
		StopType:      ServerStopTypeSoft,
		TimeoutAction: RestartTimeoutActionDestroy,
		Host:          999,
	}

	expectedJSON := `
	  {
		"restart_server": {
		  "stop_type": "soft",
		  "timeout": "300",
		  "timeout_action": "destroy",
		  "host": 999
		}
	  }
	`
	actualJSON, err := json.Marshal(&request)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/foo/restart", request.RequestURL())
}

// TestRestartServerRequest_OmitHost tests that RestartServerRequest objects behave correctly
// when Host is omitted
func TestRestartServerRequest_OmitHost(t *testing.T) {
	request := RestartServerRequest{
		UUID:          "foo",
		Timeout:       time.Minute * 5,
		StopType:      ServerStopTypeSoft,
		TimeoutAction: RestartTimeoutActionDestroy,
	}

	expectedJSON := `
	  {
		"restart_server": {
		  "stop_type": "soft",
		  "timeout": "300",
		  "timeout_action": "destroy"
		}
	  }
	`
	actualJSON, err := json.Marshal(&request)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/foo/restart", request.RequestURL())
}

// TestModifyServerRequest tests that ModifyServerRequest objects behave correctly
func TestModifyServerRequest(t *testing.T) {
	request := ModifyServerRequest{
		UUID:         "foo",
		Title:        "Modified server",
		CoreNumber:   8,
		MemoryAmount: 16384,
		Plan:         "custom",
		Metadata:     upcloud.True,
	}

	expectedJSON := `
	  {
		"server" : {
          "title": "Modified server",
		  "core_number": "8",
		  "memory_amount": "16384",
		  "plan" : "custom",
		  "metadata": "yes"
		}
	  }
	`
	actualJSON, err := json.Marshal(&request)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/foo", request.RequestURL())
}

func TestModifyServerRequest_BooleanDefaults(t *testing.T) {
	request := ModifyServerRequest{
		UUID:         "foo",
		Title:        "Modified server",
		CoreNumber:   8,
		MemoryAmount: 16384,
		Plan:         "custom",
	}

	expectedJSON := `
	  {
		"server" : {
          "title": "Modified server",
		  "core_number": "8",
		  "memory_amount": "16384",
		  "plan" : "custom"
		}
	  }
	`
	actualJSON, err := json.Marshal(&request)
	assert.Nil(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
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

	request = DeleteServerAndStoragesRequest{
		UUID:    "bar",
		Backups: DeleteStorageBackupsModeDelete,
	}

	assert.Equal(t, "/server/bar/?storages=1&backups=delete", request.RequestURL())

	request = DeleteServerAndStoragesRequest{
		UUID:    "baz",
		Backups: DeleteStorageBackupsModeKeepLatest,
	}

	assert.Equal(t, "/server/baz/?storages=1&backups=keep_latest", request.RequestURL())

	request = DeleteServerAndStoragesRequest{
		UUID:    "fizz",
		Backups: DeleteStorageBackupsModeKeep,
	}

	assert.Equal(t, "/server/fizz/?storages=1&backups=keep", request.RequestURL())
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
