package request

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"
)

/* Service Management */

func TestCancelManagedDatabaseConnection_RequestURL(t *testing.T) {
	req := &CancelManagedDatabaseConnection{
		UUID:      "fake",
		Pid:       42,
		Terminate: true,
	}
	assert.Equal(t, "/database/fake/connections/42?terminate=true", req.RequestURL())
}

func TestCloneManagedDatabaseRequest_MarshalJSON(t *testing.T) {
	t.Run("TestEmpty", func(t *testing.T) {
		req := CloneManagedDatabaseRequest{}
		d, err := json.Marshal(&req)
		assert.NoError(t, err)
		assert.Equal(t, `{"hostname_prefix":"","plan":"","type":"","zone":""}`, string(d))
	})

	req := CloneManagedDatabaseRequest{
		UUID:           "fakeuuid",
		CloneTime:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		HostNamePrefix: "fakename",
		Maintenance: ManagedDatabaseMaintenanceTimeRequest{
			DayOfWeek: "monday",
			Time:      "12:00:00",
		},
		Plan:  "fakeplan",
		Title: "faketitle",
		Type:  "faketype",
		Zone:  "fakezone",
	}
	req.Properties.Set("fakeprop", "fakevalue")
	d, err := json.MarshalIndent(&req, "", "\t")
	assert.NoError(t, err)

	const expect = `{
	"hostname_prefix": "fakename",
	"plan": "fakeplan",
	"properties": {
		"fakeprop": "fakevalue"
	},
	"title": "faketitle",
	"type": "faketype",
	"zone": "fakezone",
	"clone_time": "2021-01-01T00:00:00Z",
	"maintenance": {
		"dow": "monday",
		"time": "12:00:00"
	}
}`
	assert.Equal(t, expect, string(d))
}

func TestCloneManagedDatabaseRequest_RequestURL(t *testing.T) {
	req := CloneManagedDatabaseRequest{UUID: "fakeuuid"}
	assert.Equal(t, "/database/fakeuuid/clone", req.RequestURL())
}

func TestCreateManagedDatabaseRequest_RequestURL(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	assert.Equal(t, "/database", req.RequestURL())
}

func TestCreateManagedDatabaseRequest_MarshalJSON(t *testing.T) {
	t.Run("TestEmpty", func(t *testing.T) {
		req := CreateManagedDatabaseRequest{}
		d, err := json.Marshal(&req)
		assert.NoError(t, err)
		assert.Equal(t, `{"hostname_prefix":"","plan":"","type":"","zone":""}`, string(d))
	})

	req := CreateManagedDatabaseRequest{
		HostNamePrefix: "fakename",
		Maintenance: ManagedDatabaseMaintenanceTimeRequest{
			DayOfWeek: "monday",
			Time:      "12:00:00",
		},
		Plan:  "fakeplan",
		Title: "faketitle",
		Type:  "faketype",
		Zone:  "fakezone",
	}
	req.Properties.Set("fakeprop", "fakevalue")
	d, err := json.MarshalIndent(&req, "", "\t")
	assert.NoError(t, err)

	const expect = `{
	"hostname_prefix": "fakename",
	"plan": "fakeplan",
	"properties": {
		"fakeprop": "fakevalue"
	},
	"title": "faketitle",
	"type": "faketype",
	"zone": "fakezone",
	"maintenance": {
		"dow": "monday",
		"time": "12:00:00"
	}
}`
	assert.Equal(t, expect, string(d))
}

func TestDeleteManagedDatabaseRequest_RequestURL(t *testing.T) {
	req := DeleteManagedDatabaseRequest{UUID: "fake"}
	assert.Equal(t, "/database/fake", req.RequestURL())
}

func TestGetManagedDatabaseLogsRequest_RequestURL(t *testing.T) {
	req := &GetManagedDatabaseLogsRequest{
		UUID:   "fakeuuid",
		Limit:  100,
		Offset: "fakeoffset",
		Order:  upcloud.ManagedDatabaseLogOrderAscending,
	}
	assert.Equal(t, "/database/fakeuuid/logs?limit=100&offset=fakeoffset&order=asc", req.RequestURL())
	t.Run("ImplicitZeroOffset", func(t *testing.T) {
		req := &GetManagedDatabaseLogsRequest{
			UUID:  "fakeuuid",
			Limit: 100,
			Order: upcloud.ManagedDatabaseLogOrderAscending,
		}
		assert.Equal(t, "/database/fakeuuid/logs?limit=100&offset=0&order=asc", req.RequestURL())
	})
}

func TestGetManagedDatabaseRequest_RequestURL(t *testing.T) {
	req := GetManagedDatabaseRequest{UUID: "fakeuuid"}
	assert.Equal(t, "/database/fakeuuid", req.RequestURL())
}

func TestGetManagedDatabasesRequest_RequestURL(t *testing.T) {
	assert.Equal(t, "/database", (&GetManagedDatabasesRequest{}).RequestURL())
}

func TestGetManagedDatabaseConnectionsRequest_RequestURL(t *testing.T) {
	req := &GetManagedDatabaseConnectionsRequest{
		UUID:   "fake",
		Limit:  1000,
		Offset: 42,
	}
	assert.Equal(t, "/database/fake/connections?limit=1000&offset=42", req.RequestURL())
}

func TestGetManagedDatabaseMetricsRequest_RequestURL(t *testing.T) {
	req := &GetManagedDatabaseMetricsRequest{
		UUID:   "fake",
		Period: upcloud.ManagedDatabaseMetricPeriodHour,
	}
	assert.Equal(t, "/database/fake/metrics?period=hour", req.RequestURL())
}

func TestGetManagedDatabaseQueryStatisticsRequest_RequestURL(t *testing.T) {
	req := &GetManagedDatabaseQueryStatisticsRequest{
		UUID:   "fake",
		Limit:  1000,
		Offset: 42,
	}
	assert.Equal(t, "/database/fake/query-statistics?limit=1000&offset=42", req.RequestURL())
}

func TestModifyManagedDatabaseRequest_MarshalJSON(t *testing.T) {
	t.Run("TestEmpty", func(t *testing.T) {
		req := ModifyManagedDatabaseRequest{}
		d, err := json.Marshal(&req)
		assert.NoError(t, err)
		assert.Equal(t, "{}", string(d))
	})

	req := ModifyManagedDatabaseRequest{
		Maintenance: ManagedDatabaseMaintenanceTimeRequest{
			DayOfWeek: "monday",
			Time:      "12:00:00",
		},
		Plan:  "fakeplan",
		Title: "faketitle",
		Type:  "faketype",
		UUID:  "fakeuuid",
		Zone:  "fakezone",
	}
	req.Properties.Set("fakeprop", "fakevalue")
	d, err := json.MarshalIndent(&req, "", "\t")
	assert.NoError(t, err)

	const expect = `{
	"plan": "fakeplan",
	"properties": {
		"fakeprop": "fakevalue"
	},
	"title": "faketitle",
	"type": "faketype",
	"zone": "fakezone",
	"maintenance": {
		"dow": "monday",
		"time": "12:00:00"
	}
}`
	assert.Equal(t, expect, string(d))
}

func TestModifyManagedDatabaseRequest_RequestURL(t *testing.T) {
	req := ModifyManagedDatabaseRequest{UUID: "fake"}
	assert.Equal(t, "/database/fake", req.RequestURL())
}

func TestShutdownManagedDatabaseRequest_MarshalJSON(t *testing.T) {
	d, _ := json.Marshal(&ShutdownManagedDatabaseRequest{})
	assert.Equal(t, `{"powered":false}`, string(d))
}

func TestShutdownManagedDatabaseRequest_RequestURL(t *testing.T) {
	req := ShutdownManagedDatabaseRequest{UUID: "fake"}
	assert.Equal(t, "/database/fake", req.RequestURL())
}

func TestStartManagedDatabaseRequest_MarshalJSON(t *testing.T) {
	d, _ := json.Marshal(&StartManagedDatabaseRequest{})
	assert.Equal(t, `{"powered":true}`, string(d))
}

func TestStartManagedDatabaseRequest_RequestURL(t *testing.T) {
	req := StartManagedDatabaseRequest{UUID: "fake"}
	assert.Equal(t, "/database/fake", req.RequestURL())
}

/* Properties */

func TestManagedDatabasePropertiesRequest_Set(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	req.Properties.Set("fakeKey", "fakeValue").
		Set("fakeKey2", "fakeValue2")
	assert.Equal(t, "fakeValue", req.Properties["fakeKey"])
	assert.Equal(t, "fakeValue2", req.Properties["fakeKey2"])
}

func TestManagedDatabasePropertiesRequest_SetBool(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	req.Properties.SetBool("fakeKey", true)
	assert.Equal(t, true, req.Properties["fakeKey"])
}

func TestManagedDatabasePropertiesRequest_SetInt(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	req.Properties.SetInt("fakeKey", 42)
	assert.Equal(t, 42, req.Properties["fakeKey"])
}

func TestManagedDatabasePropertiesRequest_SetString(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	req.Properties.SetString("fakeKey", "fakeValue")
	assert.Equal(t, "fakeValue", req.Properties["fakeKey"])
}

func TestManagedDatabasePropertiesRequest_SetAutoUtilityIPFilter(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	req.Properties.SetAutoUtilityIPFilter(true)
	assert.Equal(t, true, req.Properties[upcloud.ManagedDatabasePropertyAutoUtilityIPFilter])
}

func TestManagedDatabasePropertiesRequest_SetIPFilter(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	req.Properties.SetIPFilter("192.0.2.1", "198.51.100.0/24")
	assert.Equal(t, []string{"192.0.2.1", "198.51.100.0/24"}, req.Properties[upcloud.ManagedDatabasePropertyIPFilter])
}

func TestManagedDatabasePropertiesRequest_SetPublicAccess(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	req.Properties.SetPublicAccess(true)
	assert.Equal(t, true, req.Properties[upcloud.ManagedDatabasePropertyPublicAccess])
}

func TestManagedDatabasePropertiesRequest_SetStringSlice(t *testing.T) {
	req := CreateManagedDatabaseRequest{}
	req.Properties.SetStringSlice("fakeKey", []string{"value1", "value2"})
	assert.Equal(t, []string{"value1", "value2"}, req.Properties["fakeKey"])
}

func TestManagedDatabasePropertiesRequest_Get(t *testing.T) {
	type customString string
	req := CreateManagedDatabaseRequest{}
	req.Properties.Set("test", customString("foo"))
	assert.Equal(t, customString("foo"), req.Properties.Get("test"))
}

func TestManagedDatabasePropertiesRequest_GetString(t *testing.T) {
	props := ManagedDatabasePropertiesRequest{"test": "foo"}
	v, _ := props.GetString("test")
	assert.Equal(t, "foo", v)
	_, err := props.GetString("fake")
	assert.Error(t, err)
}

func TestManagedDatabasePropertiesRequest_GetInt(t *testing.T) {
	props := ManagedDatabasePropertiesRequest{"test": 123}
	v, _ := props.GetInt("test")
	assert.Equal(t, 123, v)
	_, err := props.GetInt("fake")
	assert.Error(t, err)
}

func TestManagedDatabasePropertiesRequest_GetBool(t *testing.T) {
	props := ManagedDatabasePropertiesRequest{"test": true}
	v, _ := props.GetBool("test")
	assert.Equal(t, true, v)
	_, err := props.GetBool("fake")
	assert.Error(t, err)
}

func TestManagedDatabasePropertiesRequest_GetStringSlice(t *testing.T) {
	props := ManagedDatabasePropertiesRequest{"test": []string{"foo"}}
	v, _ := props.GetStringSlice("test")
	assert.Equal(t, []string{"foo"}, v)
	_, err := props.GetStringSlice("fake")
	assert.Error(t, err)
}

func TestManagedDatabasePropertiesRequest_GetAutoUtilityIPFilter(t *testing.T) {
	props := ManagedDatabasePropertiesRequest{upcloud.ManagedDatabasePropertyAutoUtilityIPFilter: true}
	assert.Equal(t, true, props.GetAutoUtilityIPFilter())
}

func TestManagedDatabasePropertiesRequest_GetIPFilter(t *testing.T) {
	props := ManagedDatabasePropertiesRequest{upcloud.ManagedDatabasePropertyIPFilter: []string{upcloud.ManagedDatabaseAllIPv4}}
	assert.Equal(t, []string{upcloud.ManagedDatabaseAllIPv4}, props.GetIPFilter())
}

func TestManagedDatabasePropertiesRequest_GetPublicAccess(t *testing.T) {
	props := ManagedDatabasePropertiesRequest{upcloud.ManagedDatabasePropertyPublicAccess: true}
	assert.Equal(t, true, props.GetPublicAccess())
}

/* User Management */

func TestCreateManagedDatabaseUserRequest_RequestURL(t *testing.T) {
	req := CreateManagedDatabaseUserRequest{ServiceUUID: "fakeuuid"}
	assert.Equal(t, "/database/fakeuuid/users", req.RequestURL())
}

func TestDeleteManagedDatabaseUserRequest_RequestURL(t *testing.T) {
	req := DeleteManagedDatabaseUserRequest{ServiceUUID: "fakeuuid", Username: "fakeuser"}
	assert.Equal(t, "/database/fakeuuid/users/fakeuser", req.RequestURL())
}

func TestGetManagedDatabaseUserRequest_RequestURL(t *testing.T) {
	req := GetManagedDatabaseUserRequest{ServiceUUID: "fakeuuid", Username: "fakeuser"}
	assert.Equal(t, "/database/fakeuuid/users/fakeuser", req.RequestURL())
}

func TestGetManagedDatabaseUsersRequest_RequestURL(t *testing.T) {
	req := GetManagedDatabaseUsersRequest{ServiceUUID: "fakeuuid"}
	assert.Equal(t, "/database/fakeuuid/users", req.RequestURL())
}

func TestModifyManagedDatabaseUserRequest_RequestURL(t *testing.T) {
	req := ModifyManagedDatabaseUserRequest{ServiceUUID: "fakeuuid", Username: "fakeuser"}
	assert.Equal(t, "/database/fakeuuid/users/fakeuser", req.RequestURL())
}

/* Logical Database Management */

func TestCreateManagedDatabaseLogicalDatabaseRequest_RequestURL(t *testing.T) {
	req := CreateManagedDatabaseLogicalDatabaseRequest{ServiceUUID: "fakeuuid"}
	assert.Equal(t, "/database/fakeuuid/databases", req.RequestURL())
}

func TestGetManagedDatabaseLogicalDatabasesRequest_RequestURL(t *testing.T) {
	req := GetManagedDatabaseLogicalDatabasesRequest{ServiceUUID: "fakeuuid"}
	assert.Equal(t, "/database/fakeuuid/databases", req.RequestURL())
}

func TestDeleteManagedDatabaseLogicalDatabaseRequest_RequestURL(t *testing.T) {
	req := DeleteManagedDatabaseLogicalDatabaseRequest{ServiceUUID: "fakeuuid", Name: "fakedb"}
	assert.Equal(t, "/database/fakeuuid/databases/fakedb", req.RequestURL())
}