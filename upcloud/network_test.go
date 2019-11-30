package upcloud

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestUnmarshalIPAddresses tests that IPAddresses and IPAddress structs are unmarshaled correctly
func TestUnmarshalNetworks(t *testing.T) {
	originalXML := originalXML()

	networks := Networks{}
	err := xml.Unmarshal([]byte(originalXML), &networks)
	assert.Nil(t, err)
	assert.Len(t, networks.Networks, 3)

	firstNetwork := networks.Networks[0]
	assert.Equal(t, firstNetwork.Name, "Public 80.69.172.0/22")
	assert.Equal(t, firstNetwork.Type, NetworkAccessPublic)
	assert.Equal(t, firstNetwork.UUID, "03000000-0000-4000-8001-000000000000")
	assert.Equal(t, firstNetwork.Zone, "fi-hel1")
	assert.Len(t, firstNetwork.IPnetworks, 1)
	assert.Equal(t, firstNetwork.IPnetworks[0].Address, "80.69.172.0/22")
	assert.Equal(t, firstNetwork.IPnetworks[0].DHCP, "yes")
	assert.Equal(t, firstNetwork.IPnetworks[0].DHCPDefaultRoute, "yes")
	assert.Len(t, firstNetwork.IPnetworks[0].DHCPDNS, 2)
	assert.Equal(t, firstNetwork.IPnetworks[0].Family, IPAddressFamilyIPv4)
	assert.Equal(t, firstNetwork.IPnetworks[0].Gateway, "80.69.172.1")
	secondNetwork := networks.Networks[1]
	assert.Len(t, secondNetwork.Servers, 1)
	assert.Equal(t, secondNetwork.Type, NetworkAccessUtility)
	thirdNetwork := networks.Networks[2]
	assert.Len(t, thirdNetwork.Servers, 1)
	assert.Equal(t, thirdNetwork.Type, NetworkAccessPrivate)
}

func originalXML() string {

	return `<?xml version="1.0" encoding="utf-8"?>
<networks>
  <network>
    <ip_networks>
      <ip_network>
        <address>80.69.172.0/22</address>
        <dhcp>yes</dhcp>
        <dhcp_default_route>yes</dhcp_default_route>
        <dhcp_dns>94.237.127.9</dhcp_dns>
        <dhcp_dns>94.237.40.9</dhcp_dns>
        <family>IPv4</family>
        <gateway>80.69.172.1</gateway>
      </ip_network>
    </ip_networks>
    <name>Public 80.69.172.0/22</name>
    <type>public</type>
    <uuid>03000000-0000-4000-8001-000000000000</uuid>
    <zone>fi-hel1</zone>
  </network>
  <network>
    <ip_networks>
      <ip_network>
        <address>10.5.0.0/22</address>
        <dhcp>yes</dhcp>
        <dhcp_default_route>no</dhcp_default_route>
        <dhcp_routes>10.0.0.0/8</dhcp_routes>
        <family>IPv4</family>
        <gateway>10.5.0.1</gateway>
      </ip_network>
    </ip_networks>
    <name>Private 10.5.0.0/22</name>
    <servers>
      <server>
        <title>terraform.example.com (managed by terraform)</title>
        <uuid>00ff2dee-c437-4849-bb71-e875f79f15dd</uuid>
      </server>
    </servers>
    <type>utility</type>
    <uuid>03000000-0000-4000-8032-000000000000</uuid>
    <zone>nl-ams1</zone>
  </network>
    <network>
    <ip_networks>
      <ip_network>
        <address>10.0.0.0/24</address>
        <dhcp>yes</dhcp>
        <dhcp_default_route>no</dhcp_default_route>
        <family>IPv4</family>
        <gateway>10.0.0.1</gateway>
      </ip_network>
    </ip_networks>
    <name>Kubernetes Private</name>
    <servers>
      <server>
        <title>terraform.example.com (managed by terraform)</title>
        <uuid>00ff2dee-c437-4849-bb71-e875f79f15dd</uuid>
      </server>
    </servers>
    <type>private</type>
    <uuid>03b481d5-ff17-41e7-a9f1-6ea96bab35da</uuid>
    <zone>nl-ams1</zone>
  </network>
  </networks>`
}
