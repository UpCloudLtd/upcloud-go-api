package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type FirewallContext interface {
	GetFirewallRules(ctx context.Context, r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error)
	GetFirewallRuleDetails(ctx context.Context, r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error)
	CreateFirewallRule(ctx context.Context, r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error)
	CreateFirewallRules(ctx context.Context, r *request.CreateFirewallRulesRequest) error
	DeleteFirewallRule(ctx context.Context, r *request.DeleteFirewallRuleRequest) error
}

// GetFirewallRules returns the firewall rules for the specified server
func (s *ServiceContext) GetFirewallRules(ctx context.Context, r *request.GetFirewallRulesRequest) (*upcloud.FirewallRules, error) {
	firewallRules := upcloud.FirewallRules{}
	return &firewallRules, s.get(ctx, r.RequestURL(), &firewallRules)
}

// GetFirewallRuleDetails returns extended details about the specified firewall rule
func (s *ServiceContext) GetFirewallRuleDetails(ctx context.Context, r *request.GetFirewallRuleDetailsRequest) (*upcloud.FirewallRule, error) {
	firewallRule := upcloud.FirewallRule{}
	return &firewallRule, s.get(ctx, r.RequestURL(), &firewallRule)
}

// CreateFirewallRule creates the firewall rule
func (s *ServiceContext) CreateFirewallRule(ctx context.Context, r *request.CreateFirewallRuleRequest) (*upcloud.FirewallRule, error) {
	firewallRule := upcloud.FirewallRule{}
	return &firewallRule, s.create(ctx, r, &firewallRule)
}

// CreateFirewallRules creates multiple firewall rules
func (s *ServiceContext) CreateFirewallRules(ctx context.Context, r *request.CreateFirewallRulesRequest) error {
	return s.replace(ctx, r, nil)
}

// DeleteFirewallRule deletes the specified firewall rule
func (s *ServiceContext) DeleteFirewallRule(ctx context.Context, r *request.DeleteFirewallRuleRequest) error {
	return s.delete(ctx, r)
}
