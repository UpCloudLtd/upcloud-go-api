package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func unmarshalFixture[T any](t *testing.T, fixture string) T {
	t.Helper()
	var model T
	require.NoError(t, json.Unmarshal([]byte(fixture), &model))
	return model
}

func TestCuratedResponseModelsAcrossResources(t *testing.T) {
	t.Parallel()

	t.Run("account/details", func(t *testing.T) {
		model := unmarshalFixture[AccountDetails200](t, `{"account":{"username":"fixture-account","resource_limits":{}}}`)
		require.Equal(t, "fixture-account", model.Account.Username)
	})

	t.Run("account/billing", func(t *testing.T) {
		model := unmarshalFixture[AccountGetMonthlyBillingSummary200](t, `{"billing":{"currency":"USD","total_amount":12.5}}`)
		require.NotNil(t, model.Billing.Currency)
		require.Equal(t, "USD", *model.Billing.Currency)
	})

	t.Run("account-token/details", func(t *testing.T) {
		model := unmarshalFixture[AccountTokensGetToken200](t, `{"id":"01234567-89ab-4def-8123-456789abcdef","name":"deploy-token"}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "deploy-token", *model.Name)
	})

	t.Run("account-token/list", func(t *testing.T) {
		model := unmarshalFixture[AccountTokensListTokens200](t, `[{"id":"01234567-89ab-4def-8123-456789abcdef","name":"deploy-token"}]`)
		require.NotNil(t, model[0].Name)
		require.Equal(t, "deploy-token", *model[0].Name)
	})

	t.Run("audit-log/detail", func(t *testing.T) {
		model := unmarshalFixture[AuditLogsGetAuditLog200](t, `{"account_username":"fixture-user","action":"create","auth_type":"api","cloud_event_id":"01HQXYZ123456789ABCDEFGHJK","origin":"api","resource_id":"server-fixture","resource_type":"server","time":"2026-01-02T03:04:05Z","upcloud_correlation_id":"550e8400-e29b-41d4-a716-446655440000"}`)
		require.Equal(t, "fixture-user", model.AccountUsername)
	})

	t.Run("currency/list", func(t *testing.T) {
		model := unmarshalFixture[CurrencyListCurrencies200](t, `{"currencies":{"currency":[{"code":"USD","rate":1.25,"updated":"2026-01-02T03:04:05Z"}]}}`)
		require.Equal(t, "USD", model.Currencies.Currency[0].Code)
	})

	t.Run("managed-database/details", func(t *testing.T) {
		model := unmarshalFixture[GetDatabase200](t, `{"name":"analytics","type":"pg","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "analytics", *model.Name)
	})

	t.Run("managed-database/list", func(t *testing.T) {
		model := unmarshalFixture[ListDatabases200](t, `[{"name":"analytics","type":"pg","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}]`)
		require.NotNil(t, model[0].Name)
		require.Equal(t, "analytics", *model[0].Name)
	})

	t.Run("file-storage/details", func(t *testing.T) {
		model := unmarshalFixture[GetFileStorage200](t, `{"name":"archive-share","uuid":"123e4567-e89b-12d3-a456-426614174000","size_gib":250}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "archive-share", *model.Name)
	})

	t.Run("file-storage-share/details", func(t *testing.T) {
		model := unmarshalFixture[GetFileStorageShare200](t, `{"name":"read-only","path":"/public"}`)
		require.NotNil(t, model.Path)
		require.Equal(t, "/public", *model.Path)
	})

	t.Run("file-storage/list", func(t *testing.T) {
		model := unmarshalFixture[ListFileStorages200](t, `[{"name":"archive-share","size_gib":250}]`)
		require.NotNil(t, model[0].Name)
		require.Equal(t, "archive-share", *model[0].Name)
	})

	t.Run("firewall-ruleset/details", func(t *testing.T) {
		model := unmarshalFixture[GetFirewallRuleset200](t, `{"name":"edge-policy","enabled":true,"stateful":false}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "edge-policy", *model.Name)
	})

	t.Run("gateway-connection/details", func(t *testing.T) {
		model := unmarshalFixture[GetGatewayConnection200](t, `{"name":"site-link","type":"ipsec"}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "site-link", *model.Name)
	})

	t.Run("gateway-service/details", func(t *testing.T) {
		model := unmarshalFixture[GetGatewayService200](t, `{"name":"edge-gateway","plan":"development-gwaas-1c2g","configured_status":"started"}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "edge-gateway", *model.Name)
	})

	t.Run("gateway-tunnel/details", func(t *testing.T) {
		model := unmarshalFixture[GetGatewayTunnel200](t, `{"name":"site-tunnel","tunnel_healthy":true}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "site-tunnel", *model.Name)
	})

	t.Run("ip-address/details", func(t *testing.T) {
		model := unmarshalFixture[GetIPAddressDetails200](t, `{"ip_address":{"access":"public","address":"192.0.2.20"}}`)
		require.Equal(t, "public", string(model.IpAddress.Access))
	})

	t.Run("ip-address/list", func(t *testing.T) {
		model := unmarshalFixture[ListIPAddresses200](t, `{"ip_addresses":{"ip_address":[{"access":"public","address":"192.0.2.21"}]}}`)
		require.Equal(t, "public", string(model.IpAddresses.IpAddress[0].Access))
	})

	t.Run("event/list", func(t *testing.T) {
		model := unmarshalFixture[ListEvents200](t, `{"events":{"event":[{"server":"123e4567-e89b-12d3-a456-426614174000","time":"2026-01-02T03:04:05Z","type":"server_start"}]}}`)
		require.Equal(t, "server_start", string(model.Events.Event[0].Type))
	})

	t.Run("load-balancer/details", func(t *testing.T) {
		model := unmarshalFixture[GetLoadBalancer200](t, `{"name":"public-api","configured_status":"configured","created_at":"2026-01-02T03:04:05Z","plan":"medium","updated_at":"2026-01-02T03:04:05Z"}`)
		require.Equal(t, "public-api", model.Name)
	})

	t.Run("load-balancer-backend/details", func(t *testing.T) {
		model := unmarshalFixture[GetLoadBalancerBackend200](t, `{"created_at":"2026-01-02T03:04:05Z","members":[],"name":"api-backend","properties":{},"tls_configs":[],"updated_at":"2026-01-02T03:04:05Z"}`)
		require.Equal(t, "api-backend", model.Name)
	})

	t.Run("network/details", func(t *testing.T) {
		model := unmarshalFixture[GetNetworkDetails200](t, `{"network":{"name":"private-net"}}`)
		require.Equal(t, "private-net", model.Network.Name)
	})

	t.Run("network-peering/details", func(t *testing.T) {
		model := unmarshalFixture[GetNetworkPeering200](t, `{"name":"peer-link"}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "peer-link", *model.Name)
	})

	t.Run("object-storage/details", func(t *testing.T) {
		model := unmarshalFixture[GetObjectStorage200](t, `{"name":"archive-store","region":"fi-hel1","uuid":"123e4567-e89b-12d3-a456-426614174000"}`)
		require.NotNil(t, model.Name)
		require.Equal(t, "archive-store", *model.Name)
	})

	t.Run("object-storage/list", func(t *testing.T) {
		model := unmarshalFixture[ObjectStorage2ListObjectStorages200](t, `[{"name":"archive-store","region":"fi-hel1"}]`)
		require.NotNil(t, model[0].Name)
		require.Equal(t, "archive-store", *model[0].Name)
	})

	t.Run("price/list", func(t *testing.T) {
		model := unmarshalFixture[GetPrice200](t, `{"prices":{"currency":"USD","zone":[]}}`)
		require.Equal(t, "USD", model.Prices.Currency)
	})

	t.Run("router/details", func(t *testing.T) {
		model := unmarshalFixture[GetRouter200](t, `{"router":{"name":"edge-router"}}`)
		require.Equal(t, "edge-router", model.Router.Name)
	})

	t.Run("server/details", func(t *testing.T) {
		model := unmarshalFixture[GetServer200](t, `{"server":{"hostname":"web-1","core_number":"2","created":1760000000,"license":"0","memory_amount":"4096","state":"started","title":"web","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}}`)
		require.Equal(t, "web-1", model.Server.Hostname)
	})

	t.Run("server/list", func(t *testing.T) {
		model := unmarshalFixture[ListServers200](t, `{"servers":{"server":[{"core_number":"2","created":1760000000,"hostname":"web-1","license":"0","memory_amount":"4096","state":"started","title":"web","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}]}}`)
		require.NotNil(t, model.Servers.Server)
		require.Equal(t, "web-1", (*model.Servers.Server)[0].Hostname)
	})

	t.Run("server-group/details", func(t *testing.T) {
		model := unmarshalFixture[GetServerGroup200](t, `{"server_group":{"anti_affinity":"yes","anti_affinity_status":[],"labels":{"label":[]},"servers":{"server":[]},"title":"spread-group","uuid":"123e4567-e89b-12d3-a456-426614174000"}}`)
		require.Equal(t, "spread-group", model.ServerGroup.Title)
	})

	t.Run("storage/details", func(t *testing.T) {
		model := unmarshalFixture[GetStorageInfo200](t, `{"storage":{"access":"private","encrypted":"no","labels":[],"license":0,"servers":{"server":[]},"size":10,"state":"online"}}`)
		require.Equal(t, int64(10), model.Storage.Size)
	})

	t.Run("tag/details", func(t *testing.T) {
		model := unmarshalFixture[GetTag200](t, `{"tag":{"name":"QA","servers":{"server":[]}}}`)
		require.Equal(t, "QA", model.Tag.Name)
	})

	t.Run("plan/list", func(t *testing.T) {
		model := unmarshalFixture[ListPlans200](t, `{"plans":{"plan":[{"core_number":1,"current_offering":"yes","family":"general_purpose","memory_amount":1024,"name":"1xCPU-1GB"}]}}`)
		require.Equal(t, "1xCPU-1GB", model.Plans.Plan[0].Name)
	})

	t.Run("timezone/list", func(t *testing.T) {
		model := unmarshalFixture[ListTimezones200](t, `{"timezones":{"timezone":["UTC"]}}`)
		require.Equal(t, "UTC", model.Timezones.Timezone[0])
	})

	t.Run("zone/list", func(t *testing.T) {
		model := unmarshalFixture[ListZones200](t, `{"zones":{"zone":[{"description":"Helsinki","id":"fi-hel1","public":"yes"}]}}`)
		require.Equal(t, "fi-hel1", model.Zones.Zone[0].Id)
	})

	t.Run("partner-account/list", func(t *testing.T) {
		model := unmarshalFixture[ListPartnerAccounts200](t, `[{"username":"partner-user"}]`)
		require.Equal(t, "partner-user", model[0].Username)
	})

	t.Run("permission/list", func(t *testing.T) {
		model := unmarshalFixture[ListPermissions200](t, `{"permissions":{"permission":[{"target_identifier":"*","target_type":"server","user":"fixture-user"}]}}`)
		require.Equal(t, "fixture-user", model.Permissions.Permission[0].User)
	})
}
