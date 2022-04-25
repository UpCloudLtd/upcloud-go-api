package service

import (
	"encoding/json"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type Firewall interface {
	GetFirewallRules(r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error)
	GetFirewallRuleDetails(r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error)
	CreateFirewallRule(r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error)
	CreateFirewallRules(r *request.CreateFirewallRulesRequest) error
	DeleteFirewallRule(r *request.DeleteFirewallRuleRequest) error
}

var _ Firewall = (*Service)(nil)

// GetFirewallRules returns the firewall rules for the specified server
func (s *Service) GetFirewallRules(r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error) {
	firewallRules := upcloud.FirewallRules{}
	response, err := s.basicGetRequest(r.RequestURL())

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &firewallRules)
	if err != nil {
		return nil, err
	}

	return &firewallRules, nil
}

// GetFirewallRuleDetails returns extended details about the specified firewall rule
func (s *Service) GetFirewallRuleDetails(r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error) {
	firewallRule := upcloud.FirewallRule{}
	response, err := s.basicGetRequest(r.RequestURL())

	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &firewallRule)
	if err != nil {
		return nil, err
	}

	return &firewallRule, nil
}

// CreateFirewallRule creates the firewall rule
func (s *Service) CreateFirewallRule(r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error) {
	firewallRule := upcloud.FirewallRule{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)

	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &firewallRule)
	if err != nil {
		return nil, err
	}

	return &firewallRule, nil
}

// CreateFirewallRules creates multiple firewall rules
func (s *Service) CreateFirewallRules(r *request.CreateFirewallRulesRequest) error {
	requestBody, _ := json.Marshal(r)
	_, err := s.client.PerformJSONPutRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)

	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}

// DeleteFirewallRule deletes the specified firewall rule
func (s *Service) DeleteFirewallRule(r *request.DeleteFirewallRuleRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))

	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}
