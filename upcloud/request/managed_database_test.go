package request

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud"
	"github.com/stretchr/testify/assert"
)

/* Service Management */

func TestCancelManagedDatabaseSession_RequestURL(t *testing.T) {
	req := &CancelManagedDatabaseSession{
		UUID:      "fake",
		Pid:       42,
		Terminate: true,
	}
	assert.Equal(t, "/database/fake/sessions/42?terminate=true", req.RequestURL())
}

func TestCloneManagedDatabaseRequest_MarshalJSON(t *testing.T) {
	t.Run("TestEmpty", func(t *testing.T) {
		req := CloneManagedDatabaseRequest{}
		d, err := json.Marshal(&req)
		assert.NoError(t, err)
		assert.Equal(t, `{"hostname_prefix":"","plan":"","zone":""}`, string(d))
	})

	req := CloneManagedDatabaseRequest{
		UUID:           "fakeuuid",
		CloneTime:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		HostNamePrefix: "fakename",
		BackupName:     "backup name",
		Maintenance: ManagedDatabaseMaintenanceTimeRequest{
			DayOfWeek: "monday",
			Time:      "12:00:00",
		},
		Plan:  "fakeplan",
		Title: "faketitle",
		Zone:  "fakezone",
	}
	req.Properties.Set("fakeprop", "fakevalue")
	d, err := json.MarshalIndent(&req, "", "\t")
	assert.NoError(t, err)

	const expect = `{
	"hostname_prefix": "fakename",
	"plan": "fakeplan",
	"backup_name": "backup name",
	"properties": {
		"fakeprop": "fakevalue"
	},
	"title": "faketitle",
	"zone": "fakezone",
	"clone_time": "2021-01-01T00:00:00Z",
	"maintenance": {
		"dow": "monday",
		"time": "12:00:00"
	}
}`
	assert.JSONEq(t, expect, string(d))
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

func TestGetManagedDatabaseAccessControlRequest_RequestURL(t *testing.T) {
	req := GetManagedDatabaseAccessControlRequest{ServiceUUID: "fake"}
	assert.Equal(t, "/database/fake/access-control", req.RequestURL())
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

func TestGetManagedDatabaseSessionsRequest_RequestURL(t *testing.T) {
	req := &GetManagedDatabaseSessionsRequest{
		UUID:   "fake",
		Limit:  1000,
		Offset: 42,
		Order:  "pid:desc",
	}
	assert.Equal(t, "/database/fake/sessions?limit=1000&offset=42&order=pid%3Adesc", req.RequestURL())
}

func TestModifyManagedDatabaseAccessControlRequest(t *testing.T) {
	want := `
	{
		"access_control": true,
		"extended_access_control": true
	}
	`

	r := ModifyManagedDatabaseAccessControlRequest{
		ServiceUUID:         "fakeuuid",
		ACLsEnabled:         upcloud.BoolPtr(true),
		ExtendedACLsEnabled: upcloud.BoolPtr(true),
	}
	got, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
	assert.Equal(t, "/database/fakeuuid/access-control", r.RequestURL())
}

func TestModifyManagedDatabaseAccessControlRequest_RequestURL(t *testing.T) {
	req := ModifyManagedDatabaseAccessControlRequest{ServiceUUID: "fakeuuid", ACLsEnabled: upcloud.BoolPtr(true), ExtendedACLsEnabled: upcloud.BoolPtr(true)}
	assert.Equal(t, "/database/fakeuuid/access-control", req.RequestURL())
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

func TestUpgradeManagedDatabaseRequest_RequestURL(t *testing.T) {
	req := UpgradeManagedDatabaseVersionRequest{
		UUID:          "fake",
		TargetVersion: "14",
	}

	assert.Equal(t, "/database/fake/upgrade", req.RequestURL())
}

func TestUpgradeManagedDatabaseVersion_MarshalJSON(t *testing.T) {
	expectedJson := `
	{
		"target_version": "14"	
	}`

	actualJson, _ := json.Marshal(&UpgradeManagedDatabaseVersionRequest{
		UUID:          "fake",
		TargetVersion: "14",
	})

	assert.JSONEq(t, expectedJson, string(actualJson))
}

func TestGetManagedDatabaseVersionsRequest_RequestURL(t *testing.T) {
	req := GetManagedDatabaseVersionsRequest{
		UUID: "fake",
	}

	assert.Equal(t, "/database/fake/versions", req.RequestURL())
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

func TestCreateManagedDatabaseUserRequest(t *testing.T) {
	want := `
	{
		"username": "api-doc-user",
		"password": "new-password",
		"authentication": "caching_sha2_password",
		"redis_access_control": {
			"categories": ["+@set"],
			"channels": ["*"],
			"commands": ["+set"],
			"keys": ["key_*"]
		},
		"pg_access_control": {
			"allow_replication": true
		}		
	}	  
	`
	r := CreateManagedDatabaseUserRequest{
		ServiceUUID:    "fakeuuid",
		Username:       "api-doc-user",
		Password:       "new-password",
		Authentication: upcloud.ManagedDatabaseUserAuthenticationCachingSHA2Password,
		PGAccessControl: &upcloud.ManagedDatabaseUserPGAccessControl{
			AllowReplication: upcloud.BoolPtr(true),
		},
		RedisAccessControl: &upcloud.ManagedDatabaseUserRedisAccessControl{
			Categories: &[]string{"+@set"},
			Channels:   &[]string{"*"},
			Commands:   &[]string{"+set"},
			Keys:       &[]string{"key_*"},
		},
	}
	got, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
	assert.Equal(t, "/database/fakeuuid/users", r.RequestURL())

	want = `
	{
		"username": "api-doc-user"
	}	  
	`
	r = CreateManagedDatabaseUserRequest{
		ServiceUUID: "fakeuuid",
		Username:    "api-doc-user",
	}
	got, err = json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
}

func TestModifyManagedDatabaseUserAccessControlRequest(t *testing.T) {
	want := `
	{
		"redis_access_control": {
			"categories": ["+@set"],
			"channels": ["*"],
			"commands": ["+set"],
			"keys": ["key_*"]
		},
		"pg_access_control": {
			"allow_replication": true
		}		
	}	  
	`
	r := ModifyManagedDatabaseUserAccessControlRequest{
		ServiceUUID: "fakeuuid",
		Username:    "fakeuser",
		PGAccessControl: &upcloud.ManagedDatabaseUserPGAccessControl{
			AllowReplication: upcloud.BoolPtr(true),
		},
		RedisAccessControl: &upcloud.ManagedDatabaseUserRedisAccessControl{
			Categories: &[]string{"+@set"},
			Channels:   &[]string{"*"},
			Commands:   &[]string{"+set"},
			Keys:       &[]string{"key_*"},
		},
	}
	got, err := json.Marshal(&r)
	assert.NoError(t, err)
	assert.JSONEq(t, want, string(got))
	assert.Equal(t, "/database/fakeuuid/users/fakeuser/access-control", r.RequestURL())
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

func TestCreateManagedDatabaseRequestMaintenanceTime_MarshalJSON(t *testing.T) {
	req := CreateManagedDatabaseRequest{
		HostNamePrefix: "fakename",
		Plan:           "fakeplan",
		Title:          "faketitle",
		Type:           "faketype",
		Zone:           "fakezone",
	}
	req.Maintenance = ManagedDatabaseMaintenanceTimeRequest{
		DayOfWeek: "monday",
		Time:      "12:00:00",
	}
	d, err := json.Marshal(&req)
	assert.NoError(t, err)
	expect := `
	{
		"hostname_prefix": "fakename",
		"plan": "fakeplan",
		"title": "faketitle",
		"type": "faketype",
		"zone": "fakezone",
		"maintenance": {
			"dow": "monday",
			"time": "12:00:00"
		}
	}`
	assert.JSONEq(t, expect, string(d))

	// Without maintenance time
	req.Maintenance = ManagedDatabaseMaintenanceTimeRequest{
		DayOfWeek: "monday",
	}
	d, err = json.Marshal(&req)
	assert.NoError(t, err)
	expect = `
	{
		"hostname_prefix": "fakename",
		"plan": "fakeplan",
		"title": "faketitle",
		"type": "faketype",
		"zone": "fakezone",
		"maintenance": {
			"dow": "monday"
		}
	}`
	assert.JSONEq(t, expect, string(d))

	// Without maintenance dow
	req.Maintenance = ManagedDatabaseMaintenanceTimeRequest{
		Time: "12:00:00",
	}
	d, err = json.Marshal(&req)
	assert.NoError(t, err)
	expect = `
	{
		"hostname_prefix": "fakename",
		"plan": "fakeplan",
		"title": "faketitle",
		"type": "faketype",
		"zone": "fakezone",
		"maintenance": {
			"time": "12:00:00"
		}
	}`
	assert.JSONEq(t, expect, string(d))
}

/* OpenSearch index Management */

func TestGetManagedDatabaseServiceIndicesRequest_RequestURL(t *testing.T) {
	req := GetManagedDatabaseIndicesRequest{ServiceUUID: "fakeuuid"}
	assert.Equal(t, "/database/fakeuuid/indices", req.RequestURL())
}

func TestDeleteManagedDatabaseServiceIndexRequest_RequestURL(t *testing.T) {
	req := DeleteManagedDatabaseIndexRequest{ServiceUUID: "fakeuuid", IndexName: "fakeindex"}
	assert.Equal(t, "/database/fakeuuid/indices/fakeindex", req.RequestURL())
}
