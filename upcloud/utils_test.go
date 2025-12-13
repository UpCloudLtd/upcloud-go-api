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
	assert.True(t, s.Value.Bool())
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
	assert.True(t, s.Value.Bool())
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
	assert.True(t, s.Value.Bool())
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
	assert.True(t, s.Value.Bool())
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
	assert.True(t, s.Value.Bool())
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
	assert.False(t, s.Value.Bool())
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
	assert.False(t, s.Value.Bool())
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
	assert.False(t, s.Value.Bool())
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
	assert.False(t, s.Value.Bool())
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
	assert.False(t, s.Value.Bool())
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
	assert.False(t, s.Value.Bool())
}

func TestBoolean_Empty(t *testing.T) {
	var b Boolean
	assert.True(t, b.Empty())
}

type testFloatStruct struct {
	Value StringTolerantFloat64 `json:"value"`
}

func TestStringTolerantFloat64_AsNumber(t *testing.T) {
	jsonData := `{"value": 123.45}`
	var s testFloatStruct
	err := json.Unmarshal([]byte(jsonData), &s)
	assert.NoError(t, err)
	assert.Equal(t, 123.45, s.Value.Float64())
}

func TestStringTolerantFloat64_AsString(t *testing.T) {
	jsonData := `{"value": "123.45"}`
	var s testFloatStruct
	err := json.Unmarshal([]byte(jsonData), &s)
	assert.NoError(t, err)
	assert.Equal(t, 123.45, s.Value.Float64())
}

func TestStringTolerantFloat64_AsZeroNumber(t *testing.T) {
	jsonData := `{"value": 0}`
	var s testFloatStruct
	err := json.Unmarshal([]byte(jsonData), &s)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, s.Value.Float64())
}

func TestStringTolerantFloat64_AsZeroString(t *testing.T) {
	jsonData := `{"value": "0"}`
	var s testFloatStruct
	err := json.Unmarshal([]byte(jsonData), &s)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, s.Value.Float64())
}

func TestStringTolerantFloat64_AsEmptyString(t *testing.T) {
	jsonData := `{"value": ""}`
	var s testFloatStruct
	err := json.Unmarshal([]byte(jsonData), &s)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, s.Value.Float64())
}

func TestStringTolerantFloat64_Marshal(t *testing.T) {
	s := testFloatStruct{Value: StringTolerantFloat64(456.78)}
	data, err := json.Marshal(s)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"value": 456.78}`, string(data))
}
