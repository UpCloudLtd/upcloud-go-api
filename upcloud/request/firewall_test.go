package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/stretchr/testify/assert"
)

// TestGetFirewallRulesRequest tests that GetFirewallRulesRequest behaves correctly
func TestGetFirewallRulesRequest(t *testing.T) {
	request := GetFirewallRulesRequest{
		ServerUUID: "00798b85-efdc-41ca-8021-f6ef457b8531",
	}

	assert.Equal(t, "/server/00798b85-efdc-41ca-8021-f6ef457b8531/firewall_rule", request.RequestURL())
}

// TestGetFirewallRuleDetailsRequest tests that GetFirewallRuleDetailsRequest behaves correctly
func TestGetFirewallRuleDetailsRequest(t *testing.T) {
	request := GetFirewallRuleDetailsRequest{
		ServerUUID: "00798b85-efdc-41ca-8021-f6ef457b8531",
		Position:   1,
	}

	assert.Equal(t, "/server/00798b85-efdc-41ca-8021-f6ef457b8531/firewall_rule/1", request.RequestURL())
}

// TestCreateFirewallRuleRequest tests that CreateFirewallRuleRequest behaves correctly
func TestCreateFirewallRuleRequest(t *testing.T) {
	request := CreateFirewallRuleRequest{
		ServerUUID: "00798b85-efdc-41ca-8021-f6ef457b8531",
		FirewallRule: upcloud.FirewallRule{
			Direction:            upcloud.FirewallRuleDirectionIn,
			Action:               upcloud.FirewallRuleActionAccept,
			Family:               upcloud.IPAddressFamilyIPv4,
			Position:             1,
			Comment:              "Allow SSH from this network",
			DestinationPortStart: "22",
			DestinationPortEnd:   "22",
			SourceAddressStart:   "192.168.1.1",
			SourceAddressEnd:     "192.168.1.255",
			Protocol:             upcloud.FirewallRuleProtocolTCP,
		},
	}

	// Check the request URL
	assert.Equal(t, "/server/00798b85-efdc-41ca-8021-f6ef457b8531/firewall_rule", request.RequestURL())

	// Check marshaling
	actualJSON, err := json.Marshal(&request)
	assert.NoError(t, err)

	expectedJSON := `
	{
		"firewall_rule": {
		  "position": "1",
		  "direction": "in",
		  "family": "IPv4",
		  "protocol": "tcp",
		  "source_address_start": "192.168.1.1",
		  "source_address_end": "192.168.1.255",
		  "destination_port_start": "22",
		  "destination_port_end": "22",
		  "action": "accept",
		  "comment": "Allow SSH from this network"
		}
	  }
	`
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}

// TestDeleteFirewallRuleRequest tests that DeleteFirewallRuleRequest behaves correctly
func TestDeleteFirewallRuleRequest(t *testing.T) {
	request := DeleteFirewallRuleRequest{
		ServerUUID: "00798b85-efdc-41ca-8021-f6ef457b8531",
		Position:   1,
	}

	assert.Equal(t, "/server/00798b85-efdc-41ca-8021-f6ef457b8531/firewall_rule/1", request.RequestURL())
}

// TestCreateFirewallRulesRequest tests that CreateFirewallRulesRequest behaves correctly
func TestCreateFirewallRulesRequest(t *testing.T) {
	request := CreateFirewallRulesRequest{
		ServerUUID: "foo",
		FirewallRules: []upcloud.FirewallRule{
			{
				Direction:            upcloud.FirewallRuleDirectionIn,
				Family:               upcloud.IPAddressFamilyIPv4,
				Protocol:             upcloud.FirewallRuleProtocolTCP,
				DestinationPortStart: "22",
				DestinationPortEnd:   "22",
				Action:               upcloud.FirewallRuleActionAccept,
				Comment:              "Allow SSH to this network",
			},
			{
				Direction:            upcloud.FirewallRuleDirectionIn,
				Family:               upcloud.IPAddressFamilyIPv4,
				Protocol:             upcloud.FirewallRuleProtocolTCP,
				DestinationPortStart: "80",
				DestinationPortEnd:   "80",
				Action:               upcloud.FirewallRuleActionAccept,
				Comment:              "Allow HTTP to this network",
			},
		},
	}

	expectedJSON := `
	{
		"firewall_rules": {
		"firewall_rule": [
		{
			"direction": "in",
			"family": "IPv4",
			"protocol": "tcp",
			"destination_port_start": "22",
			"destination_port_end": "22",
			"action": "accept",
			"comment": "Allow SSH to this network"
		  },
		  {
			"direction": "in",
			"family": "IPv4",
			"protocol": "tcp",
			"destination_port_start": "80",
			"destination_port_end": "80",
			"action": "accept",
			"comment": "Allow HTTP to this network"
		  }
		]
	  }
	}
	`

	actualJSON, err := json.MarshalIndent(&request, "", "  ")
	assert.NoError(t, err)

	assert.JSONEq(t, expectedJSON, string(actualJSON))
	assert.Equal(t, "/server/foo/firewall_rule", request.RequestURL())
}
