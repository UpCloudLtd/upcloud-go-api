package upcloud

import "encoding/json"

// Constants
const (
	FirewallRuleActionAccept = "accept"
	FirewallRuleActionReject = "reject"
	FirewallRuleActionDrop   = "drop"

	FirewallRuleDirectionIn  = "in"
	FirewallRuleDirectionOut = "out"

	FirewallRuleProtocolTCP  = "tcp"
	FirewallRuleProtocolUDP  = "udp"
	FirewallRuleProtocolICMP = "icmp"
)

// FirewallRules represents a list of firewall rules
type FirewallRules struct {
	FirewallRules []FirewallRule `xml:"firewall_rule" json:"firewall_rules"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *FirewallRules) UnmarshalJSON(b []byte) error {
	type localFirewallRule FirewallRule
	type firewallRuleWrapper struct {
		FirewallRules []localFirewallRule `json:"firewall_rule"`
	}

	v := struct {
		FirewallRules firewallRuleWrapper `json:"firewall_rules"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	for _, f := range v.FirewallRules.FirewallRules {
		s.FirewallRules = append(s.FirewallRules, FirewallRule(f))
	}

	return nil
}

// FirewallRule represents a single firewall rule. Note that most integer values are represented as strings
type FirewallRule struct {
	Action                  string `xml:"action" json:"action"`
	Comment                 string `xml:"comment,omitempty" json:"comment,omitempty"`
	DestinationAddressStart string `xml:"destination_address_start,omitempty" json:"destination_address_start,omitempty"`
	DestinationAddressEnd   string `xml:"destination_address_end,omitempty" json:"destination_address_end,omitempty"`
	DestinationPortStart    string `xml:"destination_port_start,omitempty" json:"destination_port_start,omitempty"`
	DestinationPortEnd      string `xml:"destination_port_end,omitempty" json:"destination_port_end,omitempty"`
	Direction               string `xml:"direction" json:"direction"`
	Family                  string `xml:"family" json:"family"`
	ICMPType                string `xml:"icmp_type,omitempty" json:"icmp_type,omitempty"`
	Position                int    `xml:"position" json:"position,string"`
	Protocol                string `xml:"protocol,omitempty" json:"protocol,omitempty"`
	SourceAddressStart      string `xml:"source_address_start,omitempty" json:"source_address_start,omitempty"`
	SourceAddressEnd        string `xml:"source_address_end,omitempty" json:"source_address_end,omitempty"`
	SourcePortStart         string `xml:"source_port_start,omitempty" json:"source_port_start,omitempty"`
	SourcePortEnd           string `xml:"source_port_end,omitempty" json:"source_port_end,omitempty"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (s *FirewallRule) UnmarshalJSON(b []byte) error {
	type localFirewallRule FirewallRule

	v := struct {
		FirewallRule localFirewallRule `json:"firewall_rule"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	(*s) = FirewallRule(v.FirewallRule)

	return nil
}
