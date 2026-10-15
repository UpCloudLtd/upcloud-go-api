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

	require.Equal(t, "fixture-account", unmarshalFixture[AccountDetails200](t, `{"account":{"username":"fixture-account","resource_limits":{}}}`).Account.Username)
	require.Equal(t, "USD", *unmarshalFixture[AccountGetMonthlyBillingSummary200](t, `{"billing":{"currency":"USD","total_amount":12.5}}`).Billing.Currency)
	require.Equal(t, "deploy-token", *unmarshalFixture[AccountTokensGetToken200](t, `{"id":"01234567-89ab-4def-8123-456789abcdef","name":"deploy-token"}`).Name)
	require.Equal(t, "deploy-token", *unmarshalFixture[AccountTokensListTokens200](t, `[{"id":"01234567-89ab-4def-8123-456789abcdef","name":"deploy-token"}]`)[0].Name)
	require.Equal(t, "fixture-user", unmarshalFixture[AuditLogsGetAuditLog200](t, `{"account_username":"fixture-user","action":"create","auth_type":"api","cloud_event_id":"01HQXYZ123456789ABCDEFGHJK","origin":"api","resource_id":"server-fixture","resource_type":"server","time":"2026-01-02T03:04:05Z","upcloud_correlation_id":"550e8400-e29b-41d4-a716-446655440000"}`).AccountUsername)
	require.Equal(t, "USD", unmarshalFixture[CurrencyListCurrencies200](t, `{"currencies":{"currency":[{"code":"USD","rate":1.25,"updated":"2026-01-02T03:04:05Z"}]}}`).Currencies.Currency[0].Code)
	require.Equal(t, "analytics", *unmarshalFixture[GetDatabase200](t, `{"name":"analytics","type":"pg","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}`).Name)
	require.Equal(t, "analytics", *unmarshalFixture[ListDatabases200](t, `[{"name":"analytics","type":"pg","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}]`)[0].Name)
	require.Equal(t, "archive-share", *unmarshalFixture[GetFileStorage200](t, `{"name":"archive-share","uuid":"123e4567-e89b-12d3-a456-426614174000","size_gib":250}`).Name)
	require.Equal(t, "/public", *unmarshalFixture[GetFileStorageShare200](t, `{"name":"read-only","path":"/public"}`).Path)
	require.Equal(t, "archive-share", *unmarshalFixture[ListFileStorages200](t, `[{"name":"archive-share","size_gib":250}]`)[0].Name)
	require.Equal(t, "edge-policy", *unmarshalFixture[GetFirewallRuleset200](t, `{"name":"edge-policy","enabled":true,"stateful":false}`).Name)
	require.Equal(t, "site-link", *unmarshalFixture[GetGatewayConnection200](t, `{"name":"site-link","type":"ipsec"}`).Name)
	require.Equal(t, "edge-gateway", *unmarshalFixture[GetGatewayService200](t, `{"name":"edge-gateway","plan":"development-gwaas-1c2g","configured_status":"started"}`).Name)
	require.Equal(t, "site-tunnel", *unmarshalFixture[GetGatewayTunnel200](t, `{"name":"site-tunnel","tunnel_healthy":true}`).Name)
	require.Equal(t, "public", string(unmarshalFixture[GetIPAddressDetails200](t, `{"ip_address":{"access":"public","address":"192.0.2.20"}}`).IpAddress.Access))
	require.Equal(t, "public", string(unmarshalFixture[ListIPAddresses200](t, `{"ip_addresses":{"ip_address":[{"access":"public","address":"192.0.2.21"}]}}`).IpAddresses.IpAddress[0].Access))
	require.Equal(t, "server_start", string(unmarshalFixture[ListEvents200](t, `{"events":{"event":[{"server":"123e4567-e89b-12d3-a456-426614174000","time":"2026-01-02T03:04:05Z","type":"server_start"}]}}`).Events.Event[0].Type))
	require.Equal(t, "public-api", unmarshalFixture[GetLoadBalancer200](t, `{"name":"public-api","configured_status":"configured","created_at":"2026-01-02T03:04:05Z","plan":"medium","updated_at":"2026-01-02T03:04:05Z"}`).Name)
	require.Equal(t, "api-backend", unmarshalFixture[GetLoadBalancerBackend200](t, `{"created_at":"2026-01-02T03:04:05Z","members":[],"name":"api-backend","properties":{},"tls_configs":[],"updated_at":"2026-01-02T03:04:05Z"}`).Name)
	require.Equal(t, "private-net", unmarshalFixture[GetNetworkDetails200](t, `{"network":{"name":"private-net"}}`).Network.Name)
	require.Equal(t, "peer-link", *unmarshalFixture[GetNetworkPeering200](t, `{"name":"peer-link"}`).Name)
	require.Equal(t, "archive-store", *unmarshalFixture[GetObjectStorage200](t, `{"name":"archive-store","region":"fi-hel1","uuid":"123e4567-e89b-12d3-a456-426614174000"}`).Name)
	require.Equal(t, "archive-store", *unmarshalFixture[ObjectStorage2ListObjectStorages200](t, `[{"name":"archive-store","region":"fi-hel1"}]`)[0].Name)
	require.Equal(t, "USD", unmarshalFixture[GetPrice200](t, `{"prices":{"currency":"USD","zone":[]}}`).Prices.Currency)
	require.Equal(t, "edge-router", unmarshalFixture[GetRouter200](t, `{"router":{"name":"edge-router"}}`).Router.Name)
	require.Equal(t, "web-1", unmarshalFixture[GetServer200](t, `{"server":{"hostname":"web-1","core_number":"2","created":1760000000,"license":"0","memory_amount":"4096","state":"started","title":"web","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}}`).Server.Hostname)
	require.Equal(t, "web-1", (*unmarshalFixture[ListServers200](t, `{"servers":{"server":[{"core_number":"2","created":1760000000,"hostname":"web-1","license":"0","memory_amount":"4096","state":"started","title":"web","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}]}}`).Servers.Server)[0].Hostname)
	require.Equal(t, "spread-group", unmarshalFixture[GetServerGroup200](t, `{"server_group":{"anti_affinity":"yes","anti_affinity_status":[],"labels":{"label":[]},"servers":{"server":[]},"title":"spread-group","uuid":"123e4567-e89b-12d3-a456-426614174000"}}`).ServerGroup.Title)
	require.Equal(t, int64(10), unmarshalFixture[GetStorageInfo200](t, `{"storage":{"access":"private","encrypted":"no","labels":[],"license":0,"servers":{"server":[]},"size":10,"state":"online"}}`).Storage.Size)
	require.Equal(t, "QA", unmarshalFixture[GetTag200](t, `{"tag":{"name":"QA","servers":{"server":[]}}}`).Tag.Name)
	require.Equal(t, "1xCPU-1GB", unmarshalFixture[ListPlans200](t, `{"plans":{"plan":[{"core_number":1,"current_offering":"yes","family":"general_purpose","memory_amount":1024,"name":"1xCPU-1GB"}]}}`).Plans.Plan[0].Name)
	require.Equal(t, "UTC", unmarshalFixture[ListTimezones200](t, `{"timezones":{"timezone":["UTC"]}}`).Timezones.Timezone[0])
	require.Equal(t, "fi-hel1", unmarshalFixture[ListZones200](t, `{"zones":{"zone":[{"description":"Helsinki","id":"fi-hel1","public":"yes"}]}}`).Zones.Zone[0].Id)
	require.Equal(t, "partner-user", unmarshalFixture[ListPartnerAccounts200](t, `[{"username":"partner-user"}]`)[0].Username)
	require.Equal(t, "fixture-user", unmarshalFixture[ListPermissions200](t, `{"permissions":{"permission":[{"target_identifier":"*","target_type":"server","user":"fixture-user"}]}}`).Permissions.Permission[0].User)
}
