package upcloud_test

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Value upcloud.Boolean `json:"value"`
}

func TestBoolean_TrueAsBool(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": true
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, s.Value.Bool())
}

func TestBoolean_TrueAsString(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": "true"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, s.Value.Bool())
}

func TestBoolean_TrueAsOne(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": 1
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, s.Value.Bool())
}

func TestBoolean_TrueAsOneString(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": "1"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, s.Value.Bool())
}

func TestBoolean_TrueAsYesString(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": "yes"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, s.Value.Bool())
}

func TestBoolean_FalseAsBool(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": false
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, s.Value.Bool())
}

func TestBoolean_FalseAsString(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": "false"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, s.Value.Bool())
}

func TestBoolean_FalseAsZero(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": 0
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, s.Value.Bool())
}

func TestBoolean_FalseAsZeroString(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": "0"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, s.Value.Bool())
}

func TestBoolean_FalseAsNo(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": "no"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, s.Value.Bool())
}

func TestBoolean_FalseAnything(t *testing.T) {
	t.Parallel()
	trueJSON := `
	{
		"value": "fudge" 
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, s.Value.Bool())
}

func TestBoolean_Empty(t *testing.T) {
	t.Parallel()
	var b upcloud.Boolean
	assert.True(t, b.Empty())
}
