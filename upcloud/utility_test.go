package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Value Boolean `json:"value"`
}

func TestBoolean_TrueAsBool(t *testing.T) {
	trueJSON := `
	{
		"value": true
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, bool(s.Value))
}

func TestBoolean_TrueAsString(t *testing.T) {
	trueJSON := `
	{
		"value": "true"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, bool(s.Value))
}

func TestBoolean_TrueAsOne(t *testing.T) {
	trueJSON := `
	{
		"value": 1
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, bool(s.Value))
}

func TestBoolean_TrueAsOneString(t *testing.T) {
	trueJSON := `
	{
		"value": "1"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, bool(s.Value))
}

func TestBoolean_TrueAsYesString(t *testing.T) {
	trueJSON := `
	{
		"value": "yes"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.True(t, bool(s.Value))
}

func TestBoolean_FalseAsBool(t *testing.T) {
	trueJSON := `
	{
		"value": false
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, bool(s.Value))
}

func TestBoolean_FalseAsString(t *testing.T) {
	trueJSON := `
	{
		"value": "false"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, bool(s.Value))
}

func TestBoolean_FalseAsZero(t *testing.T) {
	trueJSON := `
	{
		"value": 0
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, bool(s.Value))
}

func TestBoolean_FalseAsZeroString(t *testing.T) {
	trueJSON := `
	{
		"value": "0"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, bool(s.Value))
}

func TestBoolean_FalseAsNo(t *testing.T) {
	trueJSON := `
	{
		"value": "no"
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, bool(s.Value))
}

func TestBoolean_FalseAnything(t *testing.T) {
	trueJSON := `
	{
		"value": "fudge" 
	}
	`

	s := testStruct{}

	err := json.Unmarshal([]byte(trueJSON), &s)
	assert.NoError(t, err)
	assert.False(t, bool(s.Value))
}
