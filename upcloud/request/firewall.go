package request

import "fmt"

// GetServerFirewallRulesRequest represents a request for retrieving the firewall rules for a specific server
type GetFirewallRulesRequest struct {
	ServerUUID string
}

// RequestURL() implements the Request interface
func (r *GetFirewallRulesRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/firewall_rule", r.ServerUUID)
}
