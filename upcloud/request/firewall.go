package request

import "fmt"

/**
Represents a request for retrieving the firewall rules for a specific server
*/
type GetServerFirewallRulesRequest struct {
	UUID string
}

func (r *GetServerFirewallRulesRequest) RequestURL() string {
	return fmt.Sprintf("/server/%s/firewall_rule", r.UUID)
}
