package upcloud_test

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalFirewallRules tests the FirewallRules and FirewallRule are unmarshaled correctly.
func TestUnmarshalFirewallRules(t *testing.T) {
	t.Parallel()
	originalJSON := `
      {
        "firewall_rules": {
          "firewall_rule": [
            {
              "action": "accept",
              "comment": "Allow HTTP from anywhere",
              "destination_address_end": "",
              "destination_address_start": "",
              "destination_port_end": "80",
              "destination_port_start": "80",
              "direction": "in",
              "family": "IPv4",
              "icmp_type": "",
              "position": "1",
              "protocol": "",
              "source_address_end": "",
              "source_address_start": "",
              "source_port_end": "",
              "source_port_start": ""
            },
            {
              "action": "accept",
              "comment": "Allow SSH from a specific network only",
              "destination_address_end": "",
              "destination_address_start": "",
              "destination_port_end": "22",
              "destination_port_start": "22",
              "direction": "in",
              "family": "IPv4",
              "icmp_type": "",
              "position": "2",
              "protocol": "tcp",
              "source_address_end": "192.168.1.255",
              "source_address_start": "192.168.1.1",
              "source_port_end": "",
              "source_port_start": ""
            },
            {
              "action": "accept",
              "comment": "Allow SSH over IPv6 from this range",
              "destination_address_end": "",
              "destination_address_start": "",
              "destination_port_end": "22",
              "destination_port_start": "22",
              "direction": "in",
              "family": "IPv6",
              "icmp_type": "",
              "position": "3",
              "protocol": "tcp",
              "source_address_end": "2a04:3540:1000:aaaa:bbbb:cccc:d001",
              "source_address_start": "2a04:3540:1000:aaaa:bbbb:cccc:d001",
              "source_port_end": "",
              "source_port_start": ""
            },
            {
              "action": "accept",
              "comment": "Allow ICMP echo request (ping)",
              "destination_address_end": "",
              "destination_address_start": "",
              "destination_port_end": "",
              "destination_port_start": "",
              "direction": "in",
              "family": "IPv4",
              "icmp_type": "8",
              "position": "4",
              "protocol": "icmp",
              "source_address_end": "",
              "source_address_start": "",
              "source_port_end": "",
              "source_port_start": ""
            },
            {
              "action": "drop",
              "comment": "",
              "destination_address_end": "",
              "destination_address_start": "",
              "destination_port_end": "",
              "destination_port_start": "",
              "direction": "in",
              "family": "",
              "icmp_type": "",
              "position": "5",
              "protocol": "",
              "source_address_end": "",
              "source_address_start": "",
              "source_port_end": "",
              "source_port_start": ""
            }
          ]
        }
      }
    `
	firewallRules := upcloud.FirewallRules{}
	err := json.Unmarshal([]byte(originalJSON), &firewallRules)
	assert.NoError(t, err)
	assert.Len(t, firewallRules.FirewallRules, 5)

	testData := []upcloud.FirewallRule{
		{
			Action:               upcloud.FirewallRuleActionAccept,
			Comment:              "Allow HTTP from anywhere",
			DestinationPortStart: "80",
			DestinationPortEnd:   "80",
			Direction:            upcloud.FirewallRuleDirectionIn,
			Family:               upcloud.IPAddressFamilyIPv4,
			Position:             1,
		},
		{
			Action:               upcloud.FirewallRuleActionAccept,
			Comment:              "Allow SSH from a specific network only",
			DestinationPortStart: "22",
			DestinationPortEnd:   "22",
			Direction:            upcloud.FirewallRuleDirectionIn,
			Family:               upcloud.IPAddressFamilyIPv4,
			Position:             2,
			Protocol:             upcloud.FirewallRuleProtocolTCP,
			SourceAddressStart:   "192.168.1.1",
			SourceAddressEnd:     "192.168.1.255",
		},
		{
			Action:               upcloud.FirewallRuleActionAccept,
			Comment:              "Allow SSH over IPv6 from this range",
			DestinationPortStart: "22",
			DestinationPortEnd:   "22",
			Direction:            upcloud.FirewallRuleDirectionIn,
			Family:               upcloud.IPAddressFamilyIPv6,
			Position:             3,
			Protocol:             upcloud.FirewallRuleProtocolTCP,
			SourceAddressStart:   "2a04:3540:1000:aaaa:bbbb:cccc:d001",
			SourceAddressEnd:     "2a04:3540:1000:aaaa:bbbb:cccc:d001",
		},
		{
			Action:    upcloud.FirewallRuleActionAccept,
			Comment:   "Allow ICMP echo request (ping)",
			Direction: upcloud.FirewallRuleDirectionIn,
			Family:    upcloud.IPAddressFamilyIPv4,
			ICMPType:  "8",
			Position:  4,
			Protocol:  upcloud.FirewallRuleProtocolICMP,
		},
		{
			Action:    upcloud.FirewallRuleActionDrop,
			Direction: upcloud.FirewallRuleDirectionIn,
			Position:  5,
		},
	}

	for i, expectedRule := range testData {
		actualRule := firewallRules.FirewallRules[i]
		assert.Equal(t, expectedRule, actualRule)
	}
}

// TestUnmarshalFirewallRule tests that FirewallRule is unmarshaled correctly on its own.
func TestUnmarshalFirewallRule(t *testing.T) {
	t.Parallel()
	originalJSON := `
      {
        "firewall_rule": {
            "action": "accept",
            "comment": "Allow HTTP from anywhere",
            "destination_address_end": "",
            "destination_address_start": "",
            "destination_port_end": "80",
            "destination_port_start": "80",
            "direction": "in",
            "family": "IPv4",
            "icmp_type": "",
            "position": "1",
            "protocol": "",
            "source_address_end": "",
            "source_address_start": "",
            "source_port_end": "",
            "source_port_start": ""
        }
      }
    `
	actualRule := upcloud.FirewallRule{}
	err := json.Unmarshal([]byte(originalJSON), &actualRule)
	assert.NoError(t, err)

	expectedRule := upcloud.FirewallRule{
		Action:               upcloud.FirewallRuleActionAccept,
		Comment:              "Allow HTTP from anywhere",
		DestinationPortStart: "80",
		DestinationPortEnd:   "80",
		Direction:            upcloud.FirewallRuleDirectionIn,
		Family:               upcloud.IPAddressFamilyIPv4,
		Position:             1,
	}

	assert.Equal(t, expectedRule, actualRule)
}
