package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFirewallRulesContext performs the following actions with context:
//
// - creates a server
// - adds a firewall rule to the server
// - gets details about the firewall rule
// - deletes the firewall rule
func TestFirewallRulesContext(t *testing.T) {
	t.Parallel()
	recordWithContext(t, "firewallrules", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		// Create the server
		serverDetails, err := createServer(svc, "TestFirewallRules")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Create firewall rule
		t.Logf("Creating firewall rule #1 for server with UUID %s ...", serverDetails.UUID)
		_, err = svcContext.CreateFirewallRule(ctx, &request.CreateFirewallRuleRequest{
			ServerUUID: serverDetails.UUID,
			FirewallRule: upcloud.FirewallRule{
				Direction: upcloud.FirewallRuleDirectionIn,
				Action:    upcloud.FirewallRuleActionAccept,
				Family:    upcloud.IPAddressFamilyIPv4,
				Protocol:  upcloud.FirewallRuleProtocolTCP,
				Position:  1,
				Comment:   "This is the comment",
			},
		})
		require.NoError(t, err)
		t.Log("Firewall rule created")

		// Get list of firewall rules for this server
		firewallRules, err := svcContext.GetFirewallRules(ctx, &request.GetFirewallRulesRequest{
			ServerUUID: serverDetails.UUID,
		})
		require.NoError(t, err)
		assert.Len(t, firewallRules.FirewallRules, 1)
		assert.Equal(t, "This is the comment", firewallRules.FirewallRules[0].Comment)

		// Get details about the rule
		t.Log("Getting details about firewall rule #1 ...")
		firewallRule, err := svcContext.GetFirewallRuleDetails(ctx, &request.GetFirewallRuleDetailsRequest{
			ServerUUID: serverDetails.UUID,
			Position:   1,
		})
		require.NoError(t, err)
		assert.Equal(t, "This is the comment", firewallRule.Comment)
		t.Logf("Got firewall rule details, comment is %s", firewallRule.Comment)

		err = svcContext.CreateFirewallRules(ctx, &request.CreateFirewallRulesRequest{
			ServerUUID: serverDetails.UUID,
			FirewallRules: []upcloud.FirewallRule{
				{
					Direction:            upcloud.FirewallRuleDirectionIn,
					Action:               upcloud.FirewallRuleActionAccept,
					Family:               upcloud.IPAddressFamilyIPv4,
					Protocol:             upcloud.FirewallRuleProtocolTCP,
					DestinationPortStart: "80",
					DestinationPortEnd:   "80",
					Comment:              "This is a new comment 0",
				},
				{
					Direction:            upcloud.FirewallRuleDirectionIn,
					Action:               upcloud.FirewallRuleActionAccept,
					Family:               upcloud.IPAddressFamilyIPv4,
					Protocol:             upcloud.FirewallRuleProtocolTCP,
					DestinationPortStart: "22",
					DestinationPortEnd:   "22",
					Comment:              "This is a new comment 1",
				},
				{
					Direction: upcloud.FirewallRuleDirectionIn,
					Action:    upcloud.FirewallRuleActionDrop,
					Comment:   "This is a new comment 2",
				},
			},
		})
		require.NoError(t, err)

		firewallRulesPost, err := svcContext.GetFirewallRules(ctx, &request.GetFirewallRulesRequest{
			ServerUUID: serverDetails.UUID,
		})
		require.NoError(t, err)
		require.Len(t, firewallRulesPost.FirewallRules, 3)
		assert.Equal(t, firewallRulesPost.FirewallRules[2].Family, "")

		for i, rule := range firewallRulesPost.FirewallRules {
			assert.Equal(t, fmt.Sprintf("This is a new comment %d", i), rule.Comment)
		}

		// Delete the firewall rule
		t.Log("Deleting firewall rule #1 ...")
		err = svcContext.DeleteFirewallRule(ctx, &request.DeleteFirewallRuleRequest{
			ServerUUID: serverDetails.UUID,
			Position:   1,
		})
		require.NoError(t, err)
		t.Log("Firewall rule #1 deleted")

		firewallRulesPostDelete, err := svcContext.GetFirewallRules(ctx, &request.GetFirewallRulesRequest{
			ServerUUID: serverDetails.UUID,
		})
		require.NoError(t, err)
		assert.Len(t, firewallRulesPostDelete.FirewallRules, 2)
	})
}
