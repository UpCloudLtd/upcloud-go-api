package request

import "fmt"

/**
GetServerFirewallRulesRequest represents a request for retrieving the firewall rules for a specific server
*/
type GetServerFirewallRulesRequest struct {
	UUID string
}

/**
RequestURL() implements the Request interface
*/
func (r *GetServerFirewallRulesRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/firewall_rule", r.UUID)
}
