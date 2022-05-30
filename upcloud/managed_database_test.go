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

func TestManagedDatabaseServicePlanZones_UnmarshalJSON(t *testing.T) {
	const d = `{
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
}`
	expect := ManagedDatabaseServicePlanZones{
		Zones: []ManagedDatabaseServicePlanZone{{"de-fra1"}, {"fi-hel1"}},
	}

	var zones ManagedDatabaseServicePlanZones
	err := json.Unmarshal([]byte(d), &zones)
	assert.NoError(t, err)
	assert.Equal(t, expect, zones)
}
