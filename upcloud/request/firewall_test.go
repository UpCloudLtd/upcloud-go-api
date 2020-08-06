package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
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
