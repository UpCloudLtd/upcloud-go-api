package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalServerConfiguratons tests that ServerConfigurations and ServerConfiguration are unmarshaled correctly
func TestUnmarshalServerConfiguratons(t *testing.T) {
	originalJSON := `
{
    "server_sizes": {
      "server_size": [
        {
          "core_number": "1",
          "memory_amount": "512"
        },
        {
          "core_number": "1",
          "memory_amount": "768"
        },
        {
          "core_number": "10",
          "memory_amount": "65024"
        },
        {
          "core_number": "10",
          "memory_amount": "65536"
        }
      ]
    }
  }
`

	serverConfigurations := ServerConfigurations{}
	err := json.Unmarshal([]byte(originalJSON), &serverConfigurations)
	assert.Nil(t, err)
	assert.Len(t, serverConfigurations.ServerConfigurations, 4)

	testData := []ServerConfiguration{
		{
			CoreNumber:   1,
			MemoryAmount: 512,
		},
		{
			CoreNumber:   1,
			MemoryAmount: 768,
		},
		{
			CoreNumber:   10,
			MemoryAmount: 65024,
		},
		{
			CoreNumber:   10,
			MemoryAmount: 65536,
		},
	}

	for i, sc := range testData {
		configuration := serverConfigurations.ServerConfigurations[i]
		assert.Equal(t, sc.CoreNumber, configuration.CoreNumber)
		assert.Equal(t, sc.MemoryAmount, configuration.MemoryAmount)
	}
}

// TestUnmarshalServers tests that Servers and Server are unmarshaled correctly
func TestUnmarshalServers(t *testing.T) {
	originalJSON := `
        {
            "servers" : {
                "server" : [
                    {
                        "core_number" : "1",
                        "hostname": "foo",
                        "license": 0,
                        "memory_amount": "1024",
                        "plan": "1xCPU-1GB",
                        "progress": "95",
                        "state": "maintenance",
                        "tags": {
                            "tag": []
                        },
                        "title": "foo.example.com",
                        "uuid": "009114f1-cd89-4202-b057-5680eb6ba428",
                        "zone": "uk-lon1"
                    }
                ]
            }
        }
    `

	servers := Servers{}
	err := json.Unmarshal([]byte(originalJSON), &servers)
	assert.Nil(t, err)
	assert.Len(t, servers.Servers, 1)

	server := servers.Servers[0]
	assert.Equal(t, 1, server.CoreNumber)
	assert.Equal(t, "foo", server.Hostname)
	assert.Equal(t, 0.0, server.License)
	assert.Equal(t, 1024, server.MemoryAmount)
	assert.Equal(t, "1xCPU-1GB", server.Plan)
	assert.Equal(t, 95, server.Progress)
	assert.Equal(t, ServerStateMaintenance, server.State)
	assert.Empty(t, server.Tags)
	assert.Equal(t, "foo.example.com", server.Title)
	assert.Equal(t, "009114f1-cd89-4202-b057-5680eb6ba428", server.UUID)
	assert.Equal(t, "uk-lon1", server.Zone)
}

// TestUnmarshalServerDetails tests that ServerDetails objects are correctly unmarshaled
func TestUnmarshalServerDetails(t *testing.T) {
	originalJSON := `
      {
        "server": {
          "boot_order": "cdrom,disk",
          "core_number": "0",
          "firewall": "on",
          "host" : 7653311107,
          "hostname": "server1.example.com",
          "ip_addresses": {
            "ip_address": [
              {
                "access": "private",
                "address": "10.0.0.00",
                "family" : "IPv4"
              },
              {
                "access": "public",
                "address": "0.0.0.0",
                "family" : "IPv4"
              },
              {
                "access": "public",
                "address": "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
                "family" : "IPv6"
              }
            ]
          },
          "license": 0,
          "memory_amount": "2048",
          "metadata": "yes",
          "networking" : {
            "interfaces" : {
              "interface" : [
                {
                  "index" : 1,
                  "ip_addresses" : {
                    "ip_address" : [
                      {
                        "address" : "94.237.0.207",
                        "family" : "IPv4",
                        "floating" : "no"
                      }
                    ]
                  },
                  "mac" : "de:ff:ff:ff:66:89",
                  "network" : "037fcf2a-6745-45dd-867e-f9479ea8c044",
                  "type" : "public",
                  "bootable": "no"
                },
                {
                  "index" : 2,
                  "ip_addresses" : {
                    "ip_address" : [
                      {
                        "address" : "10.6.3.95",
                        "family" : "IPv4",
                        "floating" : "no"
                      }
                    ]
                  },
                  "mac" : "de:ff:ff:ff:ed:85",
                  "network" : "03000000-0000-4000-8045-000000000000",
                  "type" : "utility",
                  "bootable": "no"
                },
                {
                  "index" : 3,
                  "ip_addresses" : {
                    "ip_address" : [
                      {
                        "address" : "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
                        "family" : "IPv6",
                        "floating" : "no"
                      }
                    ]
                  },
                  "mac" : "de:ff:ff:ff:cc:20",
                  "network" : "03c93fd8-cc60-4849-91b8-6e404b228e2a",
                  "type" : "public",
                  "bootable": "no"
                }
              ]
            }
          },
          "nic_model": "virtio",
          "plan" : "1xCPU-2GB",
          "plan_ipv4_bytes": "3565675343",
          "plan_ipv6_bytes": "4534432",
          "simple_backup": "0100,dailies",
          "state": "started",
          "storage_devices": {
            "storage_device": [
              {
                "address": "virtio:0",
                "part_of_plan" : "yes",
                "storage": "012580a1-32a1-466e-a323-689ca16f2d43",
                "storage_size": 20,
                "storage_tier": "maxiops",
                "storage_title": "Storage for server1.example.com",
                "type": "disk",
                "boot_disk": "0"
              }
            ]
          },
          "tags" : {
             "tag" : [
                "DEV",
                "Ubuntu"
             ]
          },
          "timezone": "UTC",
          "title": "server1.example.com",
          "uuid": "0077fa3d-32db-4b09-9f5f-30d9e9afb565",
          "video_model": "cirrus",
          "remote_access_enabled": "yes",
          "remote_access_type": "vnc",
          "remote_access_host": "fi-hel1.vnc.upcloud.com",
          "remote_access_password": "aabbccdd",
          "remote_access_port": "3000",
          "zone": "fi-hel1"
        }
      }
    `

	serverDetails := ServerDetails{}
	err := json.Unmarshal([]byte(originalJSON), &serverDetails)
	assert.Nil(t, err)

	assert.Equal(t, "cdrom,disk", serverDetails.BootOrder)
	assert.Equal(t, "on", serverDetails.Firewall)
	assert.Len(t, serverDetails.IPAddresses, 3)
	assert.Equal(t, "virtio", serverDetails.NICModel)
	assert.Len(t, serverDetails.StorageDevices, 1)
	assert.Equal(t, "012580a1-32a1-466e-a323-689ca16f2d43", serverDetails.StorageDevices[0].UUID)
	assert.Equal(t, "virtio:0", serverDetails.StorageDevices[0].Address)
	assert.Equal(t, "yes", serverDetails.StorageDevices[0].PartOfPlan)
	assert.Equal(t, 20, serverDetails.StorageDevices[0].Size)
	assert.Equal(t, "maxiops", serverDetails.StorageDevices[0].Tier)
	assert.Equal(t, "Storage for server1.example.com", serverDetails.StorageDevices[0].Title)
	assert.Equal(t, StorageTypeDisk, serverDetails.StorageDevices[0].Type)
	assert.Equal(t, 0, serverDetails.StorageDevices[0].BootDisk)
	assert.Equal(t, "UTC", serverDetails.Timezone)
	assert.Equal(t, VideoModelCirrus, serverDetails.VideoModel)
	assert.True(t, serverDetails.RemoteAccessEnabled.Bool())
	assert.Equal(t, "fi-hel1.vnc.upcloud.com", serverDetails.RemoteAccessHost)
	assert.Equal(t, "aabbccdd", serverDetails.RemoteAccessPassword)
	assert.Equal(t, 3000, serverDetails.RemoteAccessPort)
	assert.Equal(t, RemoteAccessTypeVNC, serverDetails.RemoteAccessType)
	assert.Equal(t, "server1.example.com", serverDetails.Hostname)
	assert.Equal(t, "0100,dailies", serverDetails.SimpleBackup)
	assert.True(t, serverDetails.Metadata.Bool())

	networkingTestData := []ServerInterface{
		{
			Index: 1,
			IPAddresses: []IPAddress{
				{
					Address:  "94.237.0.207",
					Family:   IPAddressFamilyIPv4,
					Floating: False,
				},
			},
			MAC:      "de:ff:ff:ff:66:89",
			Network:  "037fcf2a-6745-45dd-867e-f9479ea8c044",
			Type:     NetworkTypePublic,
			Bootable: False,
		},
		{
			Index: 2,
			IPAddresses: []IPAddress{
				{
					Address:  "10.6.3.95",
					Family:   IPAddressFamilyIPv4,
					Floating: False,
				},
			},
			MAC:      "de:ff:ff:ff:ed:85",
			Network:  "03000000-0000-4000-8045-000000000000",
			Type:     NetworkTypeUtility,
			Bootable: False,
		},
		{
			Index: 3,
			IPAddresses: []IPAddress{
				{
					Address:  "xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx",
					Family:   "IPv6",
					Floating: False,
				},
			},
			MAC:      "de:ff:ff:ff:cc:20",
			Network:  "03c93fd8-cc60-4849-91b8-6e404b228e2a",
			Type:     NetworkTypePublic,
			Bootable: False,
		},
	}

	assert.Len(t, serverDetails.Networking.Interfaces, len(networkingTestData))

	for i, n := range networkingTestData {
		assert.Equal(t, n, serverDetails.Networking.Interfaces[i])
	}
}
