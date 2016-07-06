package upcloud

/**
Represents a list of firewall rules
*/
type FirewallRules struct {
	FirewallRules []FirewallRule `xml:"firewall_rule"`
}

/**
Represents a single firewall rule
*/
type FirewallRule struct {
	Action                  string `xml:"action"`
	DestinationAddressStart string `xml:"destination_address_start"`
	DestinationAddressEnd   string `xml:"destination_address_end"`
	DestinationPortStart    int    `xml:"destination_port_start"`
	DestinationPortEnd      int    `xml:"destination_port_end"`
	Direction               string `xml:"direction"`
	Family                  string `xml:"family"`
	ICMPType                int    `xml:"icmp_type"`
	Position                int    `xml:"position"`
	Protocol                string `xml:"protocol"`
	SourceAddressStart      string `xml:"source_address_start"`
	SourceAddressEnd        string `xml:"source_address_end"`
	SourcePortStart         int    `xml:"source_port_start"`
	SourcePortEnd           int    `xml:"source_port_end"`
}
