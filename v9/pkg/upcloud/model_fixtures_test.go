package upcloud

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type modelFixtureCase struct {
	resource string
	name     string
	model    any
	fixture  string
	path     []any
	want     any
}

func TestCuratedResponseModelsAcrossResources(t *testing.T) {
	fixtures := []modelFixtureCase{
		{resource: "account", name: "details", model: new(AccountDetails200), fixture: `{"account":{"username":"fixture-account","resource_limits":{}}}`, path: []any{"account", "username"}, want: "fixture-account"},
		{resource: "account", name: "billing", model: new(AccountGetMonthlyBillingSummary200), fixture: `{"billing":{"currency":"USD","total_amount":12.5}}`, path: []any{"billing", "currency"}, want: "USD"},
		{resource: "account-token", name: "details", model: new(AccountTokensGetToken200), fixture: `{"id":"01234567-89ab-4def-8123-456789abcdef","name":"deploy-token"}`, path: []any{"name"}, want: "deploy-token"},
		{resource: "account-token", name: "list", model: new(AccountTokensListTokens200), fixture: `[{"id":"01234567-89ab-4def-8123-456789abcdef","name":"deploy-token"}]`, path: []any{0, "name"}, want: "deploy-token"},
		{resource: "audit-log", name: "detail", model: new(AuditLogsGetAuditLog200), fixture: `{"account_username":"fixture-user","action":"create","auth_type":"api","cloud_event_id":"01HQXYZ123456789ABCDEFGHJK","origin":"api","resource_id":"server-fixture","resource_type":"server","time":"2026-01-02T03:04:05Z","upcloud_correlation_id":"550e8400-e29b-41d4-a716-446655440000"}`, path: []any{"account_username"}, want: "fixture-user"},
		{resource: "currency", name: "list", model: new(CurrencyListCurrencies200), fixture: `{"currencies":{"currency":[{"code":"USD","rate":1.25,"updated":"2026-01-02T03:04:05Z"}]}}`, path: []any{"currencies", "currency", 0, "code"}, want: "USD"},
		{resource: "managed-database", name: "details", model: new(GetDatabase200), fixture: `{"name":"analytics","type":"pg","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}`, path: []any{"name"}, want: "analytics"},
		{resource: "managed-database", name: "list", model: new(ListDatabases200), fixture: `[{"name":"analytics","type":"pg","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}]`, path: []any{0, "name"}, want: "analytics"},
		{resource: "file-storage", name: "details", model: new(GetFileStorage200), fixture: `{"name":"archive-share","uuid":"123e4567-e89b-12d3-a456-426614174000","size_gib":250}`, path: []any{"name"}, want: "archive-share"},
		{resource: "file-storage-share", name: "details", model: new(GetFileStorageShare200), fixture: `{"name":"read-only","path":"/public"}`, path: []any{"path"}, want: "/public"},
		{resource: "file-storage", name: "list", model: new(ListFileStorages200), fixture: `[{"name":"archive-share","size_gib":250}]`, path: []any{0, "name"}, want: "archive-share"},
		{resource: "firewall-ruleset", name: "details", model: new(GetFirewallRuleset200), fixture: `{"name":"edge-policy","enabled":true,"stateful":false}`, path: []any{"name"}, want: "edge-policy"},
		{resource: "gateway-connection", name: "details", model: new(GetGatewayConnection200), fixture: `{"name":"site-link","type":"ipsec"}`, path: []any{"name"}, want: "site-link"},
		{resource: "gateway-service", name: "details", model: new(GetGatewayService200), fixture: `{"name":"edge-gateway","plan":"development-gwaas-1c2g","configured_status":"started"}`, path: []any{"name"}, want: "edge-gateway"},
		{resource: "gateway-tunnel", name: "details", model: new(GetGatewayTunnel200), fixture: `{"name":"site-tunnel","tunnel_healthy":true}`, path: []any{"name"}, want: "site-tunnel"},
		{resource: "ip-address", name: "details", model: new(GetIPAddressDetails200), fixture: `{"ip_address":{"access":"public","address":"192.0.2.20"}}`, path: []any{"ip_address", "access"}, want: "public"},
		{resource: "ip-address", name: "list", model: new(ListIPAddresses200), fixture: `{"ip_addresses":{"ip_address":[{"access":"public","address":"192.0.2.21"}]}}`, path: []any{"ip_addresses", "ip_address", 0, "access"}, want: "public"},
		{resource: "event", name: "list", model: new(ListEvents200), fixture: `{"events":{"event":[{"server":"123e4567-e89b-12d3-a456-426614174000","time":"2026-01-02T03:04:05Z","type":"server_start"}]}}`, path: []any{"events", "event", 0, "type"}, want: "server_start"},
		{resource: "load-balancer", name: "details", model: new(GetLoadBalancer200), fixture: `{"name":"public-api","configured_status":"configured","created_at":"2026-01-02T03:04:05Z","plan":"medium","updated_at":"2026-01-02T03:04:05Z"}`, path: []any{"name"}, want: "public-api"},
		{resource: "load-balancer-backend", name: "details", model: new(GetLoadBalancerBackend200), fixture: `{"created_at":"2026-01-02T03:04:05Z","members":[],"name":"api-backend","properties":{},"tls_configs":[],"updated_at":"2026-01-02T03:04:05Z"}`, path: []any{"name"}, want: "api-backend"},
		{resource: "network", name: "details", model: new(GetNetworkDetails200), fixture: `{"network":{"name":"private-net"}}`, path: []any{"network", "name"}, want: "private-net"},
		{resource: "network-peering", name: "details", model: new(GetNetworkPeering200), fixture: `{"name":"peer-link"}`, path: []any{"name"}, want: "peer-link"},
		{resource: "object-storage", name: "details", model: new(GetObjectStorage200), fixture: `{"name":"archive-store","region":"fi-hel1","uuid":"123e4567-e89b-12d3-a456-426614174000"}`, path: []any{"name"}, want: "archive-store"},
		{resource: "object-storage", name: "list", model: new(ObjectStorage2ListObjectStorages200), fixture: `[{"name":"archive-store","region":"fi-hel1"}]`, path: []any{0, "name"}, want: "archive-store"},
		{resource: "price", name: "list", model: new(GetPrice200), fixture: `{"prices":{"currency":"USD","zone":[]}}`, path: []any{"prices", "currency"}, want: "USD"},
		{resource: "router", name: "details", model: new(GetRouter200), fixture: `{"router":{"name":"edge-router"}}`, path: []any{"router", "name"}, want: "edge-router"},
		{resource: "server", name: "details", model: new(GetServer200), fixture: `{"server":{"hostname":"web-1","core_number":"2","created":1760000000,"license":"0","memory_amount":"4096","state":"started","title":"web","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}}`, path: []any{"server", "hostname"}, want: "web-1"},
		{resource: "server", name: "list", model: new(ListServers200), fixture: `{"servers":{"server":[{"core_number":"2","created":1760000000,"hostname":"web-1","license":"0","memory_amount":"4096","state":"started","title":"web","uuid":"123e4567-e89b-12d3-a456-426614174000","zone":"fi-hel1"}]}}`, path: []any{"servers", "server", 0, "hostname"}, want: "web-1"},
		{resource: "server-group", name: "details", model: new(GetServerGroup200), fixture: `{"server_group":{"anti_affinity":"yes","anti_affinity_status":[],"labels":{"label":[]},"servers":{"server":[]},"title":"spread-group","uuid":"123e4567-e89b-12d3-a456-426614174000"}}`, path: []any{"server_group", "title"}, want: "spread-group"},
		{resource: "storage", name: "details", model: new(GetStorageInfo200), fixture: `{"storage":{"access":"private","encrypted":"no","labels":[],"license":0,"servers":{"server":[]},"size":10,"state":"online"}}`, path: []any{"storage", "size"}, want: int64(10)},
		{resource: "tag", name: "details", model: new(GetTag200), fixture: `{"tag":{"name":"QA","servers":{"server":[]}}}`, path: []any{"tag", "name"}, want: "QA"},
		{resource: "plan", name: "list", model: new(ListPlans200), fixture: `{"plans":{"plan":[{"core_number":1,"current_offering":"yes","family":"general_purpose","memory_amount":1024,"name":"1xCPU-1GB"}]}}`, path: []any{"plans", "plan", 0, "name"}, want: "1xCPU-1GB"},
		{resource: "timezone", name: "list", model: new(ListTimezones200), fixture: `{"timezones":{"timezone":["UTC"]}}`, path: []any{"timezones", "timezone", 0}, want: "UTC"},
		{resource: "zone", name: "list", model: new(ListZones200), fixture: `{"zones":{"zone":[{"description":"Helsinki","id":"fi-hel1","public":"yes"}]}}`, path: []any{"zones", "zone", 0, "id"}, want: "fi-hel1"},
		{resource: "partner-account", name: "list", model: new(ListPartnerAccounts200), fixture: `[{"username":"partner-user"}]`, path: []any{0, "username"}, want: "partner-user"},
		{resource: "permission", name: "list", model: new(ListPermissions200), fixture: `{"permissions":{"permission":[{"target_identifier":"*","target_type":"server","user":"fixture-user"}]}}`, path: []any{"permissions", "permission", 0, "user"}, want: "fixture-user"},
	}

	require.GreaterOrEqual(t, len(fixtures), 30, "Area 5 needs at least 30 model fixtures")
	resourceFamilies := make(map[string]struct{})
	for _, fixture := range fixtures {
		resourceFamilies[fixture.resource] = struct{}{}
		t.Run(fixture.resource+"/"+fixture.name, func(t *testing.T) {
			require.NoError(t, json.Unmarshal([]byte(fixture.fixture), fixture.model))
			require.Equal(t, fixture.want, modelValueAtJSONPath(t, fixture.model, fixture.path))
		})
	}
	require.GreaterOrEqual(t, len(resourceFamilies), 20, "Area 5 fixtures should span resource families")
}

func modelValueAtJSONPath(t *testing.T, model any, path []any) any {
	t.Helper()
	value := reflect.ValueOf(model)
	for _, segment := range path {
		switch key := segment.(type) {
		case string:
			value = dereferenceModelValue(t, value, path)
			//nolint:exhaustive // JSON object path segments only address structs and maps.
			switch value.Kind() {
			case reflect.Struct:
				fieldFound := false
				for fieldIndex := range value.NumField() {
					fieldType := value.Type().Field(fieldIndex)
					jsonName := strings.Split(fieldType.Tag.Get("json"), ",")[0]
					if jsonName == "-" || fieldType.PkgPath != "" {
						continue
					}
					if jsonName == key || jsonName == "" && fieldType.Name == key {
						value = value.Field(fieldIndex)
						fieldFound = true
						break
					}
				}
				require.True(t, fieldFound, "path %v missing JSON field %q", path, key)
			case reflect.Map:
				require.Equal(t, reflect.String, value.Type().Key().Kind(), "path %v map key is not a string", path)
				value = value.MapIndex(reflect.ValueOf(key))
				require.True(t, value.IsValid(), "path %v missing map key %q", path, key)
			default:
				t.Fatalf("path %v expected object at key %q, got %s", path, key, value.Kind())
			}
		case int:
			value = dereferenceModelValue(t, value, path)
			require.True(t, value.Kind() == reflect.Slice || value.Kind() == reflect.Array, "path %v expected array at index %d, got %s", path, key, value.Kind())
			require.Less(t, key, value.Len(), "path %v index out of range", path)
			value = value.Index(key)
		default:
			t.Fatalf("unsupported path segment type %T", segment)
		}
	}
	value = dereferenceModelValue(t, value, path)
	switch {
	case value.Kind() == reflect.String:
		return value.String()
	case value.Kind() == reflect.Bool:
		return value.Bool()
	case value.CanInt():
		return value.Int()
	case value.CanUint():
		return value.Uint()
	case value.CanFloat():
		return value.Float()
	case value.CanInterface():
		return value.Interface()
	default:
		t.Fatalf("path %v points to unsupported value kind %s", path, value.Kind())
		return nil
	}
}

func dereferenceModelValue(t *testing.T, value reflect.Value, path []any) reflect.Value {
	t.Helper()
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		require.False(t, value.IsNil(), "path %v points through nil value", path)
		value = value.Elem()
	}
	require.True(t, value.IsValid(), "path %v points to an invalid value", path)
	return value
}
