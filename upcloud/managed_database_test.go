package upcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			"backup_config": {
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
		}
	}
}`
	expect := ManagedDatabaseType{
		Name:                   "mysql",
		Description:            "MySQL - Relational Database Management System",
		LatestAvailableVersion: "8.0.26",
		ServicePlans: []ManagedDatabaseServicePlan{{
			BackupConfig: ManagedDatabaseBackupConfig{
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

func TestUnmarshalManagedDatabaseConnection(t *testing.T) {
	const d = `{
		"application_name":"upcloudGoTestSuite",
		"backend_start":"2021-09-14T13:41:29.369748Z",
		"backend_type":"client backend",
		"backend_xid":null,
		"backend_xmin":null,
		"client_addr":"85.7.82.158",
		"client_hostname":null,
		"client_port":62028,
		"datid":16401,
		"datname":"defaultdb",
		"pid":257,
		"query":"SELECT 1",
		"query_duration":291591000,
		"query_start":"2021-09-14T13:41:29.575568Z",
		"state":"idle",
		"state_change":"2021-09-14T13:41:29.575916Z",
		"usename":"upadmin",
		"usesysid":16400,
		"wait_event":"ClientRead",
		"wait_event_type":"Client",
		"xact_start":null
	}`

	var c ManagedDatabaseConnection
	assert.NoError(t, json.Unmarshal([]byte(d), &c))

	expect := ManagedDatabaseConnection{
		ApplicationName: "upcloudGoTestSuite",
		BackendStart:    timeParse("2021-09-14T13:41:29.369748Z"),
		BackendType:     "client backend",
		ClientAddr:      "85.7.82.158",
		ClientPort:      62028,
		DatId:           16401,
		DatName:         "defaultdb",
		Pid:             257,
		Query:           "SELECT 1",
		QueryDuration:   291591000,
		QueryStart:      timeParse("2021-09-14T13:41:29.575568Z"),
		State:           "idle",
		StateChange:     timeParse("2021-09-14T13:41:29.575916Z"),
		Username:        "upadmin",
		UseSysId:        16400,
		WaitEvent:       "ClientRead",
		WaintEventType:  "Client",
	}
	assert.Equal(t, expect, c)
}
