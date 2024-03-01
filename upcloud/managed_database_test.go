package upcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagedDatabaseMetricsChartFloat64_MarshalJSON(t *testing.T) {
	chart := ManagedDatabaseMetricsChartFloat64{
		ManagedDatabaseMetricsChartHeader: ManagedDatabaseMetricsChartHeader{
			Columns: []ManagedDatabaseMetricsColumn{
				{
					Label: "node1",
					Type:  "number",
				},
				{
					Label: "node2",
					Type:  "number",
				},
			},
			Title: "fake",
			Timestamps: []time.Time{
				time.Date(2021, 8, 19, 7, 22, 0, 0, time.UTC),
				time.Date(2021, 8, 19, 7, 22, 30, 0, time.UTC),
			},
		},
		Rows: [][]float64{
			{10.1, 20.2},
			{30.3, 40.4},
		},
	}
	const expect = `{
	"data": {
		"cols": [
			{
				"label": "time",
				"type": "date"
			},
			{
				"label": "node1",
				"type": "number"
			},
			{
				"label": "node2",
				"type": "number"
			}
		],
		"rows": [
			[
				"2021-08-19T07:22:00Z",
				10.1,
				20.2
			],
			[
				"2021-08-19T07:22:30Z",
				30.3,
				40.4
			]
		]
	},
	"hints": {
		"title": "fake"
	}
}`
	d, err := json.MarshalIndent(&chart, "", "\t")
	assert.NoError(t, err)
	assert.Equal(t, expect, string(d))
}

func TestManagedDatabaseMetricsChartFloat64_UnmarshalJSON(t *testing.T) {
	const d = `{
	"data": {
		"cols": [
			{
				"label": "time",
				"type": "date"
			},
			{
				"label": "node1",
				"type": "number"
			},
			{
				"label": "node2",
				"type": "number"
			}
		],
		"rows": [
			[
				"2021-08-19T07:22:00Z",
				10.1,
				20.2
			],
			[
				"2021-08-19T07:22:30Z",
				30.3,
				40.4
			]
		]
	},
	"hints": {
		"title": "fake"
	}
}`
	expect := ManagedDatabaseMetricsChartFloat64{
		ManagedDatabaseMetricsChartHeader: ManagedDatabaseMetricsChartHeader{
			Columns: []ManagedDatabaseMetricsColumn{
				{
					Label: "node1",
					Type:  "number",
				},
				{
					Label: "node2",
					Type:  "number",
				},
			},
			Title: "fake",
			Timestamps: []time.Time{
				time.Date(2021, 8, 19, 7, 22, 0, 0, time.UTC),
				time.Date(2021, 8, 19, 7, 22, 30, 0, time.UTC),
			},
		},
		Rows: [][]float64{
			{10.1, 20.2},
			{30.3, 40.4},
		},
	}

	var chart ManagedDatabaseMetricsChartFloat64
	err := json.Unmarshal([]byte(d), &chart)
	assert.NoError(t, err)
	assert.Equal(t, expect, chart)
}

func TestManagedDatabaseMetricsChartInt_MarshalJSON(t *testing.T) {
	chart := ManagedDatabaseMetricsChartInt{
		ManagedDatabaseMetricsChartHeader: ManagedDatabaseMetricsChartHeader{
			Columns: []ManagedDatabaseMetricsColumn{
				{
					Label: "node1",
					Type:  "number",
				},
				{
					Label: "node2",
					Type:  "number",
				},
			},
			Title: "fake",
			Timestamps: []time.Time{
				time.Date(2021, 8, 19, 7, 22, 0, 0, time.UTC),
				time.Date(2021, 8, 19, 7, 22, 30, 0, time.UTC),
			},
		},
		Rows: [][]int{
			{10, 20},
			{30, 40},
		},
	}
	const expect = `{
	"data": {
		"cols": [
			{
				"label": "time",
				"type": "date"
			},
			{
				"label": "node1",
				"type": "number"
			},
			{
				"label": "node2",
				"type": "number"
			}
		],
		"rows": [
			[
				"2021-08-19T07:22:00Z",
				10,
				20
			],
			[
				"2021-08-19T07:22:30Z",
				30,
				40
			]
		]
	},
	"hints": {
		"title": "fake"
	}
}`
	d, err := json.MarshalIndent(&chart, "", "\t")
	assert.NoError(t, err)
	assert.Equal(t, expect, string(d))
}

func TestManagedDatabaseMetricsChartInt_UnmarshalJSON(t *testing.T) {
	const d = `{
	"data": {
		"cols": [
			{
				"label": "time",
				"type": "date"
			},
			{
				"label": "node1",
				"type": "number"
			},
			{
				"label": "node2",
				"type": "number"
			}
		],
		"rows": [
			[
				"2021-08-19T07:22:00Z",
				10,
				20
			],
			[
				"2021-08-19T07:22:30Z",
				30,
				40
			]
		]
	},
	"hints": {
		"title": "fake"
	}
}`
	expect := ManagedDatabaseMetricsChartInt{
		ManagedDatabaseMetricsChartHeader: ManagedDatabaseMetricsChartHeader{
			Columns: []ManagedDatabaseMetricsColumn{
				{
					Label: "node1",
					Type:  "number",
				},
				{
					Label: "node2",
					Type:  "number",
				},
			},
			Title: "fake",
			Timestamps: []time.Time{
				time.Date(2021, 8, 19, 7, 22, 0, 0, time.UTC),
				time.Date(2021, 8, 19, 7, 22, 30, 0, time.UTC),
			},
		},
		Rows: [][]int{
			{10, 20},
			{30, 40},
		},
	}

	var chart ManagedDatabaseMetricsChartInt
	err := json.Unmarshal([]byte(d), &chart)
	assert.NoError(t, err)
	assert.Equal(t, expect, chart)
}

func TestManagedDatabaseProperties_Get(t *testing.T) {
	type customString string
	props := ManagedDatabaseProperties{"test": customString("foo")}
	assert.Equal(t, customString("foo"), props.Get("test"))
}

func TestManagedDatabaseProperties_GetAutoUtilityIPFilter(t *testing.T) {
	props := ManagedDatabaseProperties{ManagedDatabasePropertyAutoUtilityIPFilter: true}
	assert.Equal(t, true, props.GetAutoUtilityIPFilter())
}

func TestManagedDatabaseProperties_GetBool(t *testing.T) {
	props := ManagedDatabaseProperties{"test": true}
	v, _ := props.GetBool("test")
	assert.Equal(t, true, v)
	_, err := props.GetBool("fake")
	assert.Error(t, err)
}

func TestManagedDatabaseProperties_GetIPFilter(t *testing.T) {
	props := ManagedDatabaseProperties{ManagedDatabasePropertyIPFilter: []string{ManagedDatabaseAllIPv4}}
	assert.Equal(t, []string{ManagedDatabaseAllIPv4}, props.GetIPFilter())
}

func TestManagedDatabaseProperties_GetInt(t *testing.T) {
	props := ManagedDatabaseProperties{"test": 123}
	v, _ := props.GetInt("test")
	assert.Equal(t, 123, v)
	_, err := props.GetInt("fake")
	assert.Error(t, err)
}

func TestManagedDatabaseProperties_GetPublicAccess(t *testing.T) {
	props := ManagedDatabaseProperties{ManagedDatabasePropertyPublicAccess: true}
	assert.Equal(t, true, props.GetPublicAccess())
}

func TestManagedDatabaseProperties_GetString(t *testing.T) {
	props := ManagedDatabaseProperties{"test": "foo"}
	v, _ := props.GetString("test")
	assert.Equal(t, "foo", v)
	_, err := props.GetString("fake")
	assert.Error(t, err)
}

func TestManagedDatabaseProperties_GetStringSlice(t *testing.T) {
	props := ManagedDatabaseProperties{"test": []string{"foo"}}
	v, _ := props.GetStringSlice("test")
	assert.Equal(t, []string{"foo"}, v)
	_, err := props.GetStringSlice("fake")
	assert.Error(t, err)
}

func TestManagedDatabaseType_UnmarshalJSON(t *testing.T) {
	const d = `{
	"name": "mysql",
	"description": "MySQL - Relational Database Management System",
	"latest_available_version": "8.0.26",
	"service_plans": [
		{
			"backup_config_mysql": {
				"interval": 24,
				"max_count": 2,
				"recovery_mode": "pitr"
			},
			"core_number": 1,
			"node_count": 1,
			"memory_amount": 2048,
			"plan": "1x1xCPU-2GB-25GB",
			"storage_size": 25600,
			"zones": {
				"zone": [
					{
						"name": "de-fra1"
					},
					{
						"name": "fi-hel1"
					}
				]
			}
		}
	],
	"properties": {
		"public_access": {
			"default": false,
			"title": "Public Access",
			"type": "boolean",
			"description": "Allow access to the service from the public Internet"
		},
		"timescaledb": {
			"title": "TimescaleDB extension configuration values",
			"type": "object",
			"properties": {
				"max_background_workers": {
					"default": 16,
					"example": 8,
					"title": "timescaledb.max_background_workers",
					"type": "integer",
					"description": "The number of background workers for timescaledb operations. You should configure this setting to the sum of your number of databases and the total number of concurrent background workers you want running at any given point in time.",
					"minimum": 1,
					"maximum": 4096
				}
			},
			"description": "System-wide settings for the timescaledb extension"
		}
	}
}`
	expect := ManagedDatabaseType{
		Name:                   "mysql",
		Description:            "MySQL - Relational Database Management System",
		LatestAvailableVersion: "8.0.26",
		ServicePlans: []ManagedDatabaseServicePlan{{
			BackupConfigMySQL: &ManagedDatabaseBackupConfigMySQL{
				Interval:     24,
				MaxCount:     2,
				RecoveryMode: "pitr",
			},
			CoreNumber:   1,
			NodeCount:    1,
			MemoryAmount: 2048,
			Plan:         "1x1xCPU-2GB-25GB",
			StorageSize:  25600,
			Zones:        []ManagedDatabaseServicePlanZone{{"de-fra1"}, {"fi-hel1"}},
		}},
		Properties: map[string]ManagedDatabaseServiceProperty{
			"public_access": {
				Default:     false,
				Title:       "Public Access",
				Type:        "boolean",
				Description: "Allow access to the service from the public Internet",
			},
			// `timescaledb` is a PostgreSQL property (not MySQL), but that shouldn't make difference from marshaling point-of-view.
			"timescaledb": {
				Title:       "TimescaleDB extension configuration values",
				Description: "System-wide settings for the timescaledb extension",
				Type:        "object",
				Properties: map[string]ManagedDatabaseServiceProperty{
					"max_background_workers": {
						Default:     16.0,
						Example:     8.0,
						Title:       "timescaledb.max_background_workers",
						Type:        "integer",
						Description: "The number of background workers for timescaledb operations. You should configure this setting to the sum of your number of databases and the total number of concurrent background workers you want running at any given point in time.",
						Minimum:     Float64Ptr(1),
						Maximum:     Float64Ptr(4096),
					},
				},
			},
		},
	}

	var databaseType ManagedDatabaseType
	err := json.Unmarshal([]byte(d), &databaseType)
	assert.NoError(t, err)
	assert.Equal(t, expect, databaseType)
}

func TestUnmarshalManagedDatabaseQueryStatisticsPostgreSQL(t *testing.T) {
	const d = `{
		"blk_read_time": 0,
		"blk_write_time": 0,
		"calls": 1,
		"database_name": "defaultdb",
		"local_blks_dirtied": 0,
		"local_blks_hit": 0,
		"local_blks_read": 0,
		"local_blks_written": 0,
		"max_time": 3126,
		"mean_time": 3126,
		"min_time": 3126,
		"query": "BEGIN",
		"rows": 0,
		"shared_blks_dirtied": 0,
		"shared_blks_hit": 0,
		"shared_blks_read": 0,
		"shared_blks_written": 0,
		"stddev_time": 0,
		"temp_blks_read": 0,
		"temp_blks_written": 0,
		"total_time": 3126,
		"user_name": "upadmin"
	}`
	var c ManagedDatabaseQueryStatisticsPostgreSQL
	assert.NoError(t, json.Unmarshal([]byte(d), &c))
	expect := ManagedDatabaseQueryStatisticsPostgreSQL{
		BlockReadTime:       0,
		BlockWriteTime:      0,
		Calls:               1,
		DatabaseName:        "defaultdb",
		LocalBlocksDirtied:  0,
		LocalBlocksHit:      0,
		LocalBlocksRead:     0,
		LocalBlocksWritten:  0,
		MaxTime:             3126,
		MeanTime:            3126,
		MinTime:             3126,
		Query:               "BEGIN",
		Rows:                0,
		SharedBlocksDirtied: 0,
		SharedBlocksHit:     0,
		SharedBlocksRead:    0,
		SharedBlocksWritten: 0,
		StddevTime:          0,
		TempBlocksRead:      0,
		TempBlocksWritten:   0,
		TotalTime:           3126,
		UserName:            "upadmin",
	}
	assert.Equal(t, expect, c)
}

func TestUnmarshalManagedDatabaseQueryStatisticsMySQL(t *testing.T) {
	const d = `{
		"avg_timer_wait": 2000309754,
		"count_star": 1,
		"digest": "3b74a4bb71311c8dffd4b64a1dabe02d0ac8fcd7ae03087dd69378a554bb28bc",
		"digest_text": "SELECT SLEEP (?)",
		"first_seen": "2021-09-14T13:49:55.576024Z",
		"last_seen": "2021-09-14T13:49:55.576024Z",
		"max_timer_wait": 2000309754,
		"min_timer_wait": 2000309754,
		"quantile_95": 2089296130,
		"quantile_99": 2089296130,
		"quantile_999": 2089296130,
		"query_sample_seen": "2021-09-14T13:49:55.576024Z",
		"query_sample_text": "SELECT SLEEP(2)",
		"query_sample_timer_wait": 2000309754,
		"schema_name": "defaultdb",
		"sum_created_tmp_disk_tables": 0,
		"sum_created_tmp_tables": 0,
		"sum_errors": 0,
		"sum_lock_time": 0,
		"sum_no_good_index_used": 0,
		"sum_no_index_used": 0,
		"sum_rows_affected": 0,
		"sum_rows_examined": 1,
		"sum_rows_sent": 1,
		"sum_select_full_join": 0,
		"sum_select_full_range_join": 0,
		"sum_select_range": 0,
		"sum_select_range_check": 0,
		"sum_select_scan": 0,
		"sum_sort_merge_passes": 0,
		"sum_sort_range": 0,
		"sum_sort_rows": 0,
		"sum_sort_scan": 0,
		"sum_timer_wait": 2000309754,
		"sum_warnings": 0
	}`
	var c ManagedDatabaseQueryStatisticsMySQL
	assert.NoError(t, json.Unmarshal([]byte(d), &c))
	expect := ManagedDatabaseQueryStatisticsMySQL{
		AvgTimerWait:            2000309754,
		CountStar:               1,
		Digest:                  "3b74a4bb71311c8dffd4b64a1dabe02d0ac8fcd7ae03087dd69378a554bb28bc",
		DigestText:              "SELECT SLEEP (?)",
		FirstSeen:               timeParse("2021-09-14T13:49:55.576024Z"),
		LastSeen:                timeParse("2021-09-14T13:49:55.576024Z"),
		MaxTimerWait:            2000309754,
		MinTimerWait:            2000309754,
		Quantile95:              2089296130,
		Quantile99:              2089296130,
		Quantile999:             2089296130,
		QuerySampleSeen:         timeParse("2021-09-14T13:49:55.576024Z"),
		QuerySampleText:         "SELECT SLEEP(2)",
		QuerySampleTimerWait:    2000309754,
		SchemaName:              "defaultdb",
		SumCreatedTmpDiskTables: 0,
		SumCreatedTmpTables:     0,
		SumErrors:               0,
		SumLockTime:             0,
		SumNoGoodIndexUsed:      0,
		SumNoIndexUsed:          0,
		SumRowsAffected:         0,
		SumRowsExamined:         1,
		SumRowsSent:             1,
		SumSelectFullJoin:       0,
		SumSelectFullRangeJoin:  0,
		SumSelectRange:          0,
		SumSelectRangeCheck:     0,
		SumSelectScan:           0,
		SumSortMergePasses:      0,
		SumSortRange:            0,
		SumSortRows:             0,
		SumSortScan:             0,
		SumTimerWait:            2000309754,
		SumWarnings:             0,
	}
	assert.Equal(t, expect, c)
}

func TestUnmarshalManagedDatabaseSessions(t *testing.T) {
	const d = `{
	"mysql": [
		{
			"application_name": "",
			"client_addr": "111.111.111.111:63244",
			"datname": "defaultdb",
			"id": "pid_23325",
			"query": "select\n            ordinal_position,\n            column_name,\n            column_type,\n            column_default,\n            generation_expression,\n            table_name,\n            column_comment,\n            is_nullable,\n            extra,\n            collation_name\n          from information_schema.columns\n          where table_schema = 'performance_schema'\n          order by table_name, ordinal_position",
			"query_duration": 0,
			"state": "active",
			"usename": "upadmin"
		}
	],
	"pg": [
		{
			"application_name": "client 1.5.14",
			"backend_start": "2022-01-21T13:26:26.682858Z",
			"backend_type": "client backend",
			"backend_xid": 2,
			"backend_xmin": 1,
			"client_addr": "111.111.111.111",
			"client_hostname": "client.host.com",
			"client_port": 61264,
			"datid": 16401,
			"datname": "defaultdb",
			"id": "pid_2065031",
			"query": "SELECT \trel.relname, \trel.relkind, \trel.reltuples, \tcoalesce(rel.relpages,0) + coalesce(toast.relpages,0) AS num_total_pages, \tSUM(ind.relpages) AS index_pages, \tpg_roles.rolname AS owner FROM pg_class rel \tleft join pg_class toast on (toast.oid = rel.reltoastrelid) \tleft join pg_index on (indrelid=rel.oid) \tleft join pg_class ind on (ind.oid = indexrelid) \tjoin pg_namespace on (rel.relnamespace =pg_namespace.oid ) \tleft join pg_roles on ( rel.relowner = pg_roles.oid ) WHERE rel.relkind IN ('r','v','m','f','p') AND nspname = 'public'GROUP BY rel.relname, rel.relkind, rel.reltuples, coalesce(rel.relpages,0) + coalesce(toast.relpages,0), pg_roles.rolname;\n",
			"query_duration": 12225858000,
			"query_start": "2022-01-21T13:26:28.63132Z",
			"state": "idle",
			"state_change": "2022-01-21T13:26:28.63388Z",
			"usename": "upadmin",
			"usesysid": 16400,
			"wait_event": "ClientRead",
			"wait_event_type": "Client",
			"xact_start": "2022-01-21T13:26:28.63383Z"
		}
	],
	"redis": [
		{
			"active_channel_subscriptions": 0,
			"active_database": "",
			"active_pattern_matching_channel_subscriptions": 0,
			"client_addr": "[fff0:fff0:fff0:fff0:0:fff0:fff0:fff0]:39956",
			"connection_age": 2079483000000000,
			"connection_idle": 3000000000,
			"flags": [],
			"flags_raw": "N",
			"id": "15",
			"multi_exec_commands": -1,
			"application_name": "",
			"output_buffer": 0,
			"output_buffer_memory": 0,
			"output_list_length": 0,
			"query": "info",
			"query_buffer": 0,
			"query_buffer_free": 0
		}
	]
}`

	var c ManagedDatabaseSessions
	assert.NoError(t, json.Unmarshal([]byte(d), &c))

	expect := ManagedDatabaseSessions{
		MySQL: []ManagedDatabaseSessionMySQL{{
			ApplicationName: "",
			ClientAddr:      "111.111.111.111:63244",
			Datname:         "defaultdb",
			Id:              "pid_23325",
			Query:           "select\n            ordinal_position,\n            column_name,\n            column_type,\n            column_default,\n            generation_expression,\n            table_name,\n            column_comment,\n            is_nullable,\n            extra,\n            collation_name\n          from information_schema.columns\n          where table_schema = 'performance_schema'\n          order by table_name, ordinal_position",
			QueryDuration:   0,
			State:           "active",
			Usename:         "upadmin",
		}},
		PostgreSQL: []ManagedDatabaseSessionPostgreSQL{{
			ApplicationName: "client 1.5.14",
			BackendStart:    time.Date(2022, 0o1, 21, 13, 26, 26, 682858000, time.UTC),
			BackendType:     "client backend",
			BackendXid:      IntPtr(2),
			BackendXmin:     IntPtr(1),
			ClientAddr:      "111.111.111.111",
			ClientHostname:  StringPtr("client.host.com"),
			ClientPort:      61264,
			Datid:           16401,
			Datname:         "defaultdb",
			Id:              "pid_2065031",
			Query:           "SELECT \trel.relname, \trel.relkind, \trel.reltuples, \tcoalesce(rel.relpages,0) + coalesce(toast.relpages,0) AS num_total_pages, \tSUM(ind.relpages) AS index_pages, \tpg_roles.rolname AS owner FROM pg_class rel \tleft join pg_class toast on (toast.oid = rel.reltoastrelid) \tleft join pg_index on (indrelid=rel.oid) \tleft join pg_class ind on (ind.oid = indexrelid) \tjoin pg_namespace on (rel.relnamespace =pg_namespace.oid ) \tleft join pg_roles on ( rel.relowner = pg_roles.oid ) WHERE rel.relkind IN ('r','v','m','f','p') AND nspname = 'public'GROUP BY rel.relname, rel.relkind, rel.reltuples, coalesce(rel.relpages,0) + coalesce(toast.relpages,0), pg_roles.rolname;\n",
			QueryDuration:   12225858000,
			QueryStart:      time.Date(2022, 0o1, 21, 13, 26, 28, 631320000, time.UTC),
			State:           "idle",
			StateChange:     time.Date(2022, 0o1, 21, 13, 26, 28, 633880000, time.UTC),
			Usename:         "upadmin",
			Usesysid:        16400,
			WaitEvent:       "ClientRead",
			WaitEventType:   "Client",
			XactStart:       TimePtr(time.Date(2022, 0o1, 21, 13, 26, 28, 633830000, time.UTC)),
		}},
		Redis: []ManagedDatabaseSessionRedis{{
			ActiveChannelSubscriptions:                0,
			ActiveDatabase:                            "",
			ActivePatternMatchingChannelSubscriptions: 0,
			ApplicationName:                           "",
			ClientAddr:                                "[fff0:fff0:fff0:fff0:0:fff0:fff0:fff0]:39956",
			ConnectionAge:                             2079483000000000,
			ConnectionIdle:                            3000000000,
			Flags:                                     []string{},
			FlagsRaw:                                  "N",
			Id:                                        "15",
			MultiExecCommands:                         -1,
			OutputBuffer:                              0,
			OutputBufferMemory:                        0,
			OutputListLength:                          0,
			Query:                                     "info",
			QueryBuffer:                               0,
			QueryBufferFree:                           0,
		}},
	}
	assert.Equal(t, expect, c)
}

func TestManagedDatabaseMetadata(t *testing.T) {
	got := ManagedDatabase{}
	require.NoError(t, json.Unmarshal([]byte(`
	{
		"metadata": {
			"redis_version": "7",
			"mysql_version": "8.0.30",
			"pg_version": "14.6",
			"write_block_threshold_exceeded": false,
			"max_connections": 100
		}
	}`), &got))
	want := ManagedDatabase{
		Metadata: &ManagedDatabaseMetadata{
			MaxConnections:              100,
			PGVersion:                   "14.6",
			MySQLVersion:                "8.0.30",
			RedisVersion:                "7",
			WriteBlockThresholdExceeded: BoolPtr(false),
		},
	}
	assert.Equal(t, want, got)

	got = ManagedDatabase{}
	require.NoError(t, json.Unmarshal([]byte(`
	{
		"metadata": {
			"mysql_version": "8.0.30",
			"write_block_threshold_exceeded": null
		}
	}`), &got))
	want = ManagedDatabase{
		Metadata: &ManagedDatabaseMetadata{
			MySQLVersion:                "8.0.30",
			WriteBlockThresholdExceeded: nil,
		},
	}
	assert.Equal(t, want, got)

	got = ManagedDatabase{}
	require.NoError(t, json.Unmarshal([]byte(`
	{
		"metadata": {
			"write_block_threshold_exceeded": true
		}
	}`), &got))
	want = ManagedDatabase{
		Metadata: &ManagedDatabaseMetadata{
			WriteBlockThresholdExceeded: BoolPtr(true),
		},
	}
	assert.Equal(t, want, got)
}

func TestManagedDatabaseBackups(t *testing.T) {
	got := ManagedDatabase{}
	require.NoError(t, json.Unmarshal([]byte(`
	{
		"backups": [
			{
				"backup_name": "backup-name",
				"backup_time": "2022-12-08T08:07:40.960849Z",
				"data_size": 864733663
			}
		]
	}`), &got))
	want := ManagedDatabase{
		Backups: []ManagedDatabaseBackup{
			{
				BackupName: "backup-name",
				BackupTime: timeParse("2022-12-08T08:07:40.960849Z"),
				DataSize:   864733663,
			},
		},
	}
	assert.Equal(t, want, got)
}

func TestManagedDatabaseUser(t *testing.T) {
	got := ManagedDatabaseUser{}
	require.NoError(t, json.Unmarshal([]byte(`
	{
		"username": "api-doc-user",
		"password": "new-password",
		"authentication": "caching_sha2_password",
		"type": "regular",
		"redis_access_control": {
			"categories": ["+@set"],
			"channels": ["*"],
			"commands": ["+set"],
			"keys": ["key_*"]
		},
		"pg_access_control": {
			"allow_replication": false
		}
	}`), &got))
	want := ManagedDatabaseUser{
		Authentication: "caching_sha2_password",
		Type:           "regular",
		Password:       "new-password",
		Username:       "api-doc-user",
		PGAccessControl: &ManagedDatabaseUserPGAccessControl{
			AllowReplication: BoolPtr(false),
		},
		RedisAccessControl: &ManagedDatabaseUserRedisAccessControl{
			Categories: &[]string{"+@set"},
			Channels:   &[]string{"*"},
			Commands:   &[]string{"+set"},
			Keys:       &[]string{"key_*"},
		},
	}
	assert.Equal(t, want, got)
}
