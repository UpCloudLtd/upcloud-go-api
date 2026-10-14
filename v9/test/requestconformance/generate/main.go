package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

var (
	checkOnly        = flag.Bool("check", false, "fail if generated adapters are stale")
	operationComment = regexp.MustCompile(`^Corresponds with ([A-Z]+) (.+) \(the ` + "`" + `([^` + "`" + `]+)` + "`" + ` operationId\)\.$`)
)

const (
	parameterLocationHeader = "header"
	schemaKeyExamples       = "examples"
	schemaTypeArray         = "array"
	schemaTypeObject        = "object"
)

type specDocument struct {
	Info struct {
		Version string `json:"version"`
	} `json:"info"`
	Paths      map[string]json.RawMessage `json:"paths"`
	Components struct {
		Parameters    map[string]json.RawMessage `json:"parameters"`
		RequestBodies map[string]json.RawMessage `json:"requestBodies"`
		Schemas       map[string]json.RawMessage `json:"schemas"`
	} `json:"components"`
}

type specOperation struct {
	OperationID string            `json:"operationId"`
	Parameters  []json.RawMessage `json:"parameters"`
	RequestBody json.RawMessage   `json:"requestBody"`
}

type specParameter struct {
	Name     string          `json:"name"`
	In       string          `json:"in"`
	Required bool            `json:"required"`
	Style    string          `json:"style"`
	Explode  *bool           `json:"explode"`
	Schema   json.RawMessage `json:"schema"`
}

type specRequestBody struct {
	Required bool                       `json:"required"`
	Content  map[string]json.RawMessage `json:"content"`
}

type queryContract struct {
	Name       string
	In         string
	Style      string
	Explode    bool
	Required   bool
	SchemaType string
	Format     string
	Properties []string
	SampleJSON string
}

type operation struct {
	ID             string
	Method         string
	Path           string
	PathParams     []queryContract
	Query          []queryContract
	Headers        []headerContract
	BodyRequired   bool
	BodyMediaType  []string
	BodySampleJSON string
}

type headerContract struct {
	Name     string
	Required bool
}

type clientMethod struct {
	Name       string
	Operation  string
	Method     string
	Path       string
	Params     []*ast.Field
	Imports    map[string]string
	SourceFile string
}

func main() {
	flag.Parse()

	root, err := filepath.Abs(".")
	fatalIf(err)
	specPath := filepath.Join(root, "generator/openapi/spec.json")
	clientDir := filepath.Join(root, "pkg/upcloud")
	outputPath := filepath.Join(root, "test/requestconformance/adapters_gen_test.go")

	specBytes, err := os.ReadFile(specPath)
	fatalIf(err)
	var document specDocument
	fatalIf(json.Unmarshal(specBytes, &document))
	if document.Info.Version == "" {
		fatalIf(errors.New("spec has no info.version"))
	}

	operations, err := loadOperations(document)
	fatalIf(err)
	methods, err := loadClientMethods(clientDir)
	fatalIf(err)
	generated, err := renderAdapters(operations, methods, specBytes, document.Info.Version)
	fatalIf(err)

	if *checkOnly {
		current, readErr := os.ReadFile(outputPath)
		if readErr != nil {
			fatalIf(fmt.Errorf("read generated adapters: %w", readErr))
		}
		if !bytes.Equal(current, generated) {
			fatalIf(errors.New("request adapters are stale; from v9 run `make generate-request-adapters` after SDK generation, then review adapters_gen_test.go"))
		}
		return
	}

	fatalIf(os.MkdirAll(filepath.Dir(outputPath), 0o755))
	//nolint:gosec // Generated Go source contains no secrets and is intentionally readable.
	fatalIf(os.WriteFile(outputPath, generated, 0o644))
}

func loadOperations(document specDocument) ([]operation, error) {
	var operations []operation
	seen := make(map[string]struct{})
	for path, rawPath := range document.Paths {
		var pathItem map[string]json.RawMessage
		if err := json.Unmarshal(rawPath, &pathItem); err != nil {
			return nil, fmt.Errorf("decode path %s: %w", path, err)
		}
		pathParameters, err := readParameters(pathItem["parameters"], document)
		if err != nil {
			return nil, fmt.Errorf("path %s parameters: %w", path, err)
		}
		for method, rawOperation := range pathItem {
			method = strings.ToLower(method)
			if !isHTTPMethod(method) {
				continue
			}
			var source specOperation
			if err := json.Unmarshal(rawOperation, &source); err != nil {
				return nil, fmt.Errorf("decode %s %s: %w", method, path, err)
			}
			if source.OperationID == "" {
				return nil, fmt.Errorf("%s %s has no operationId", method, path)
			}
			key := strings.ToUpper(method) + " " + path
			if _, exists := seen[key]; exists {
				return nil, fmt.Errorf("duplicate operation key %s", key)
			}
			seen[key] = struct{}{}

			queryByKey := make(map[string]queryContract)
			pathByName := make(map[string]queryContract)
			headerByName := make(map[string]headerContract)
			for _, parameter := range pathParameters {
				switch parameter.In {
				case "path":
					pathByName[parameter.Name] = parameter
				case "query":
					queryByKey[parameter.Name] = parameter
				case parameterLocationHeader:
					headerByName[strings.ToLower(parameter.Name)] = headerContract{Name: parameter.Name, Required: parameter.Required}
				}
			}
			operationParameters, err := readParameters(rawOperationParameters(rawOperation), document)
			if err != nil {
				return nil, fmt.Errorf("%s %s parameters: %w", method, path, err)
			}
			for _, parameter := range operationParameters {
				switch parameter.In {
				case "path":
					pathByName[parameter.Name] = parameter
				case "query":
					queryByKey[parameter.Name] = parameter
				case parameterLocationHeader:
					headerByName[strings.ToLower(parameter.Name)] = headerContract{Name: parameter.Name, Required: parameter.Required}
				}
			}
			pathParams := make([]queryContract, 0, len(pathByName))
			for _, parameter := range pathByName {
				pathParams = append(pathParams, parameter)
			}
			sort.Slice(pathParams, func(i, j int) bool { return pathParams[i].Name < pathParams[j].Name })
			query := make([]queryContract, 0, len(queryByKey))
			for _, parameter := range queryByKey {
				query = append(query, parameter)
			}
			sort.Slice(query, func(i, j int) bool { return query[i].Name < query[j].Name })
			headers := make([]headerContract, 0, len(headerByName))
			for _, header := range headerByName {
				headers = append(headers, header)
			}
			sort.Slice(headers, func(i, j int) bool { return strings.ToLower(headers[i].Name) < strings.ToLower(headers[j].Name) })

			bodyRequired, bodyMediaTypes, bodySampleJSON, err := readRequestBody(source.RequestBody, document)
			if err != nil {
				return nil, fmt.Errorf("%s %s requestBody: %w", method, path, err)
			}
			operations = append(operations, operation{
				ID:             source.OperationID,
				Method:         strings.ToUpper(method),
				Path:           path,
				PathParams:     pathParams,
				Query:          query,
				Headers:        headers,
				BodyRequired:   bodyRequired,
				BodyMediaType:  bodyMediaTypes,
				BodySampleJSON: bodySampleJSON,
			})
		}
	}
	sort.Slice(operations, func(i, j int) bool {
		left := operations[i].Method + " " + operations[i].Path
		right := operations[j].Method + " " + operations[j].Path
		return left < right
	})
	if len(operations) == 0 {
		return nil, errors.New("spec operation inventory is empty")
	}
	return operations, nil
}

func rawOperationParameters(raw json.RawMessage) json.RawMessage {
	var value struct {
		Parameters json.RawMessage `json:"parameters"`
	}
	_ = json.Unmarshal(raw, &value)
	return value.Parameters
}

func readParameters(raw json.RawMessage, document specDocument) ([]queryContract, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var parameters []json.RawMessage
	if err := json.Unmarshal(raw, &parameters); err != nil {
		return nil, err
	}
	result := make([]queryContract, 0, len(parameters))
	for _, encoded := range parameters {
		resolved, err := resolveReference(encoded, document.Components.Parameters, "parameters")
		if err != nil {
			return nil, err
		}
		var parameter specParameter
		if err := json.Unmarshal(resolved, &parameter); err != nil {
			return nil, err
		}
		style := parameter.Style
		if style == "" {
			if parameter.In == "path" || parameter.In == parameterLocationHeader {
				style = "simple"
			} else {
				style = "form"
			}
		}
		explode := style == "form"
		if parameter.Explode != nil {
			explode = *parameter.Explode
		}
		schema, err := readSchema(parameter.Schema, document)
		if err != nil {
			return nil, err
		}
		result = append(result, queryContract{
			Name:       parameter.Name,
			In:         parameter.In,
			Style:      style,
			Explode:    explode,
			Required:   parameter.Required,
			SchemaType: schema.SchemaType,
			Format:     schema.Format,
			Properties: schema.Properties,
			SampleJSON: schema.SampleJSON,
		})
	}
	return result, nil
}

func readSchema(raw json.RawMessage, document specDocument) (queryContract, error) {
	var value struct {
		Ref        string                     `json:"$ref"`
		Type       string                     `json:"type"`
		Format     string                     `json:"format"`
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if err := json.Unmarshal(raw, &value); err != nil {
		return queryContract{}, err
	}
	if value.Ref != "" {
		resolved, err := resolveSchemaReference(raw, document)
		if err != nil {
			return queryContract{}, err
		}
		return readSchema(resolved, document)
	}
	properties := make([]string, 0, len(value.Properties))
	for name := range value.Properties {
		properties = append(properties, name)
	}
	sort.Strings(properties)
	sample, err := sampleSchema(raw, document, 0, 0)
	if err != nil {
		return queryContract{}, err
	}
	return queryContract{SchemaType: value.Type, Format: value.Format, Properties: properties, SampleJSON: string(sample)}, nil
}

func readRequestBody(raw json.RawMessage, document specDocument) (bool, []string, string, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return false, nil, "", nil
	}
	resolved, err := resolveReference(raw, document.Components.RequestBodies, "requestBodies")
	if err != nil {
		return false, nil, "", err
	}
	var body specRequestBody
	if err := json.Unmarshal(resolved, &body); err != nil {
		return false, nil, "", err
	}
	mediaTypes := make([]string, 0, len(body.Content))
	for mediaType := range body.Content {
		mediaTypes = append(mediaTypes, mediaType)
	}
	sort.Strings(mediaTypes)
	if len(mediaTypes) == 0 {
		return body.Required, mediaTypes, "{}", nil
	}
	var media struct {
		Schema json.RawMessage `json:"schema"`
	}
	if err := json.Unmarshal(body.Content[mediaTypes[0]], &media); err != nil {
		return false, nil, "", err
	}
	if len(media.Schema) == 0 {
		return body.Required, mediaTypes, "{}", nil
	}
	sample, err := sampleSchema(media.Schema, document, 0, 0)
	if err != nil {
		return false, nil, "", err
	}
	return body.Required, mediaTypes, string(sample), nil
}

func sampleSchema(raw json.RawMessage, document specDocument, depth, variation int) (json.RawMessage, error) {
	if depth > 64 {
		return nil, errors.New("schema sample recursion exceeded 64 levels")
	}
	var schema map[string]json.RawMessage
	if err := json.Unmarshal(raw, &schema); err != nil {
		return nil, err
	}
	if len(schema) == 0 {
		return json.RawMessage(`{}`), nil
	}
	if reference, ok := schema["$ref"]; ok {
		var name string
		if err := json.Unmarshal(reference, &name); err != nil {
			return nil, err
		}
		resolved, err := resolveSchemaReference(raw, document)
		if err != nil {
			return nil, err
		}
		_ = name
		return sampleSchema(resolved, document, depth+1, variation)
	}

	if value, ok := schema["const"]; ok && len(value) > 0 {
		return value, nil
	}
	if values, ok := schema["enum"]; ok {
		var entries []json.RawMessage
		if err := json.Unmarshal(values, &entries); err != nil {
			return nil, fmt.Errorf("decode enum: %w", err)
		}
		if len(entries) > 0 {
			return entries[variation%len(entries)], nil
		}
	}
	container := schemaIsContainer(schema)
	for _, key := range []string{schemaKeyExamples, "example", "default"} {
		if value, ok := schema[key]; ok {
			if schemaIsContainer(schema) {
				continue
			}
			if container {
				continue
			}
			if key != schemaKeyExamples && len(value) > 0 && string(value) != "null" {
				return value, nil
			}
			var entries []json.RawMessage
			if err := json.Unmarshal(value, &entries); err != nil {
				if key == schemaKeyExamples {
					return nil, fmt.Errorf("decode examples: %w", err)
				}
				return value, nil
			}
			if len(entries) > 0 {
				return entries[variation%len(entries)], nil
			}
		}
	}
	for _, key := range []string{"oneOf", "anyOf"} {
		if choices, ok := schema[key]; ok {
			var branches []json.RawMessage
			if err := json.Unmarshal(choices, &branches); err != nil {
				return nil, fmt.Errorf("decode %s: %w", key, err)
			}
			if schemaIsObject(schema) {
				base := make(map[string]json.RawMessage, len(schema))
				for name, value := range schema {
					if name != key && name != "example" && name != schemaKeyExamples && name != "default" {
						base[name] = value
					}
				}
				baseJSON, err := json.Marshal(base)
				if err != nil {
					return nil, err
				}
				var lastErr error
				for index, branch := range branches {
					combined, err := json.Marshal(map[string]any{"allOf": []json.RawMessage{baseJSON, branch}})
					if err != nil {
						return nil, err
					}
					sample, err := sampleSchema(combined, document, depth+1, variation+index)
					if err == nil {
						return sample, nil
					}
					lastErr = err
				}
				return nil, fmt.Errorf("no valid deterministic sample for %s: %w", key, lastErr)
			}
			var lastErr error
			for index, branch := range branches {
				sample, err := sampleSchema(branch, document, depth+1, variation+index)
				if err == nil {
					return sample, nil
				}
				lastErr = err
			}
			return nil, fmt.Errorf("no sample for %s: %w", key, lastErr)
		}
	}
	if allOf, ok := schema["allOf"]; ok && !container {
		var branches []json.RawMessage
		if err := json.Unmarshal(allOf, &branches); err != nil {
			return nil, err
		}
		for _, branch := range branches {
			var branchSchema map[string]json.RawMessage
			if json.Unmarshal(branch, &branchSchema) == nil && schemaIsObject(branchSchema) {
				combined := make(map[string]json.RawMessage, len(schema)+1)
				for name, value := range schema {
					combined[name] = value
				}
				combined["type"] = json.RawMessage(strconv.Quote(schemaTypeObject))
				encoded, err := json.Marshal(combined)
				if err != nil {
					return nil, err
				}
				return sampleSchema(encoded, document, depth+1, variation)
			}
		}
		merged := make(map[string]json.RawMessage)
		for index, branch := range branches {
			var conditional map[string]json.RawMessage
			if json.Unmarshal(branch, &conditional) == nil && conditional["if"] != nil {
				continue
			}
			sample, err := sampleSchema(branch, document, depth+1, variation+index)
			if err != nil {
				return nil, err
			}
			var object map[string]json.RawMessage
			if err := json.Unmarshal(sample, &object); err != nil {
				return nil, fmt.Errorf("allOf branch sample is not an object: %w", err)
			}
			for name, value := range object {
				merged[name] = value
			}
		}
		return json.Marshal(merged)
	}

	var schemaType string
	if encodedType, ok := schema["type"]; ok {
		if err := json.Unmarshal(encodedType, &schemaType); err != nil {
			var types []string
			if err := json.Unmarshal(encodedType, &types); err != nil {
				return nil, err
			}
			for _, candidate := range types {
				if candidate != "null" {
					schemaType = candidate
					break
				}
			}
		}
	}
	if _, ok := schema["properties"]; ok || schemaType == schemaTypeObject || schema["additionalProperties"] != nil || schema["allOf"] != nil {
		propertySchemas, required, err := schemaObjectParts(schema, document)
		if err != nil {
			return nil, err
		}
		if minimum := rawInteger(schema["minProperties"]); minimum > len(required) {
			propertyNames := make([]string, 0, len(propertySchemas))
			for name := range propertySchemas {
				propertyNames = append(propertyNames, name)
			}
			sort.Strings(propertyNames)
			for _, name := range propertyNames {
				if !containsString(required, name) {
					required = append(required, name)
					if len(required) == minimum {
						break
					}
				}
			}
		}
		sort.Strings(required)
		object := make(map[string]json.RawMessage, len(required))
		for index, name := range required {
			property, exists := propertySchemas[name]
			if !exists {
				return nil, fmt.Errorf("required property %q has no schema", name)
			}
			value, err := sampleSchema(property, document, depth+1, variation+index)
			if err != nil {
				return nil, fmt.Errorf("property %s: %w", name, err)
			}
			object[name] = value
		}
		if allOf := schema["allOf"]; len(allOf) > 0 {
			var branches []json.RawMessage
			if err := json.Unmarshal(allOf, &branches); err != nil {
				return nil, err
			}
			for _, branch := range branches {
				var conditional struct {
					If   json.RawMessage `json:"if"`
					Then json.RawMessage `json:"then"`
				}
				if err := json.Unmarshal(branch, &conditional); err != nil || len(conditional.If) == 0 || len(conditional.Then) == 0 {
					continue
				}
				if !conditionMatches(conditional.If, object) {
					continue
				}
				var consequence struct {
					Required []string `json:"required"`
				}
				if err := json.Unmarshal(conditional.Then, &consequence); err != nil {
					return nil, err
				}
				for index, name := range consequence.Required {
					if _, exists := object[name]; exists {
						continue
					}
					property, exists := propertySchemas[name]
					if !exists {
						return nil, fmt.Errorf("conditional required property %q has no schema", name)
					}
					value, err := sampleSchema(property, document, depth+1, variation+index)
					if err != nil {
						return nil, fmt.Errorf("conditional property %s: %w", name, err)
					}
					object[name] = value
				}
			}
		}
		if additional, exists := schema["additionalProperties"]; exists && len(object) == 0 {
			var allowed bool
			if json.Unmarshal(additional, &allowed) == nil {
				if allowed {
					object["request-sweep"] = json.RawMessage(`"request-sweep-value"`)
				}
			} else {
				value, err := sampleSchema(additional, document, depth+1, variation)
				if err != nil {
					return nil, fmt.Errorf("additionalProperties: %w", err)
				}
				key := "request-sweep"
				if propertyNames, exists := schema["propertyNames"]; exists {
					keyJSON, err := sampleSchema(propertyNames, document, depth+1, variation)
					if err != nil {
						return nil, fmt.Errorf("propertyNames: %w", err)
					}
					if err := json.Unmarshal(keyJSON, &key); err != nil {
						return nil, err
					}
				}
				object[key], err = json.Marshal(value)
				if err != nil {
					return nil, err
				}
			}
		}
		return json.Marshal(object)
	}

	switch schemaType {
	case "string":
		return sampleString(schema)
	case "integer":
		return sampleNumber(schema, true)
	case "number":
		return sampleNumber(schema, false)
	case "boolean":
		return json.RawMessage("false"), nil
	case schemaTypeArray:
		itemSchema, ok := schema["items"]
		if !ok {
			return nil, errors.New("array schema has no items schema")
		}
		length := rawInteger(schema["minItems"])
		if length < 1 {
			length = 1
		}
		items := make([]json.RawMessage, length)
		for index := range items {
			item, err := sampleSchema(itemSchema, document, depth+1, variation+index)
			if err != nil {
				return nil, err
			}
			items[index] = item
		}
		return json.Marshal(items)
	case "null":
		return json.RawMessage("null"), nil
	default:
		return nil, fmt.Errorf("no deterministic sample strategy for schema type %q", schemaType)
	}
}

func schemaObjectParts(schema map[string]json.RawMessage, document specDocument) (map[string]json.RawMessage, []string, error) {
	properties := make(map[string]json.RawMessage)
	var required []string
	var addSchema func(json.RawMessage) error
	addSchema = func(raw json.RawMessage) error {
		var part map[string]json.RawMessage
		if err := json.Unmarshal(raw, &part); err != nil {
			return err
		}
		if reference, exists := part["$ref"]; exists {
			var name string
			if err := json.Unmarshal(reference, &name); err != nil {
				return err
			}
			resolved, err := resolveSchemaReference(raw, document)
			if err != nil {
				return err
			}
			return addSchema(resolved)
		}
		var partProperties map[string]json.RawMessage
		if len(part["properties"]) > 0 {
			if err := json.Unmarshal(part["properties"], &partProperties); err != nil {
				return err
			}
			for name, value := range partProperties {
				properties[name] = value
			}
		}
		var partRequired []string
		if len(part["required"]) > 0 {
			if err := json.Unmarshal(part["required"], &partRequired); err != nil {
				return err
			}
		}
		required = append(required, partRequired...)
		if allOf, exists := part["allOf"]; exists {
			var branches []json.RawMessage
			if err := json.Unmarshal(allOf, &branches); err != nil {
				return err
			}
			for _, branch := range branches {
				var conditional map[string]json.RawMessage
				if json.Unmarshal(branch, &conditional) == nil && conditional["if"] != nil {
					continue
				}
				if err := addSchema(branch); err != nil {
					return err
				}
			}
		}
		return nil
	}

	root, err := json.Marshal(schema)
	if err != nil {
		return nil, nil, err
	}
	if err := addSchema(root); err != nil {
		return nil, nil, err
	}
	seen := make(map[string]struct{}, len(required))
	uniqueRequired := make([]string, 0, len(required))
	for _, name := range required {
		if _, exists := seen[name]; !exists {
			seen[name] = struct{}{}
			uniqueRequired = append(uniqueRequired, name)
		}
	}
	return properties, uniqueRequired, nil
}

func schemaIsContainer(schema map[string]json.RawMessage) bool {
	if len(schema["properties"]) > 0 || len(schema["items"]) > 0 || len(schema["additionalProperties"]) > 0 || len(schema["oneOf"]) > 0 || len(schema["anyOf"]) > 0 || len(schema["allOf"]) > 0 {
		return true
	}
	var schemaType string
	if err := json.Unmarshal(schema["type"], &schemaType); err == nil {
		return schemaType == schemaTypeObject || schemaType == schemaTypeArray
	}
	var schemaTypes []string
	if err := json.Unmarshal(schema["type"], &schemaTypes); err == nil {
		for _, item := range schemaTypes {
			if item == schemaTypeObject || item == schemaTypeArray {
				return true
			}
		}
	}
	return false
}

func schemaIsObject(schema map[string]json.RawMessage) bool {
	var schemaType string
	if json.Unmarshal(schema["type"], &schemaType) == nil && schemaType == schemaTypeObject {
		return true
	}
	return len(schema["properties"]) > 0
}

func conditionMatches(rawCondition json.RawMessage, object map[string]json.RawMessage) bool {
	var condition struct {
		Properties map[string]struct {
			Const json.RawMessage `json:"const"`
		} `json:"properties"`
		Required []string `json:"required"`
	}
	if json.Unmarshal(rawCondition, &condition) != nil {
		return false
	}
	for _, name := range condition.Required {
		if _, exists := object[name]; !exists {
			return false
		}
	}
	for name, expected := range condition.Properties {
		if len(expected.Const) == 0 {
			continue
		}
		actual, exists := object[name]
		if !exists || !bytes.Equal(actual, expected.Const) {
			return false
		}
	}
	return true
}

func sampleString(schema map[string]json.RawMessage) (json.RawMessage, error) {
	var format string
	_ = json.Unmarshal(schema["format"], &format)
	value := "request-sweep-value"
	switch format {
	case "date-time":
		value = "2026-01-02T03:04:05Z"
	case "date":
		value = "2026-01-02"
	case "uuid":
		value = "00000000-0000-4000-8000-000000000000"
	case "email":
		value = "request-sweep@example.test"
	case "ipv4":
		value = "192.0.2.1"
	case "ipv6":
		value = "2001:db8::1"
	case "uri", "url":
		value = "https://example.test/request-sweep"
	case "hostname":
		value = "request-sweep.example.test"
	}
	minimum := rawInteger(schema["minLength"])
	maximum := rawInteger(schema["maxLength"])
	if minimum > len(value) {
		value += strings.Repeat("x", minimum-len(value))
	}
	if maximum > 0 && len(value) > maximum {
		value = value[:maximum]
	}
	if pattern := rawString(schema["pattern"]); pattern != "" && !sampleMatchesPattern(pattern, value) {
		candidates := []string{"2026", "01", "12", "123", "QA", "DEV", "fi-hel1", "ab", "x", "example", "request-sweep", "192.0.2.1", "a@example.test"}
		matched := false
		for _, candidate := range candidates {
			if len(candidate) < minimum || (maximum > 0 && len(candidate) > maximum) {
				continue
			}
			if sampleMatchesPattern(pattern, candidate) {
				value = candidate
				matched = true
				break
			}
		}
		if !matched {
			return nil, fmt.Errorf("pattern %q needs an OpenAPI example or a curated sample", pattern)
		}
	}
	return json.Marshal(value)
}

func sampleNumber(schema map[string]json.RawMessage, integer bool) (json.RawMessage, error) {
	value := 1.0
	if minimum := rawFloat(schema["minimum"]); minimum != nil {
		value = *minimum
	}
	if exclusive := rawFloat(schema["exclusiveMinimum"]); exclusive != nil && *exclusive >= value {
		value = *exclusive + 1
	} else if rawBool(schema["exclusiveMinimum"]) {
		value++
	}
	if maximum := rawFloat(schema["maximum"]); maximum != nil && value > *maximum {
		value = *maximum
	}
	if exclusive := rawFloat(schema["exclusiveMaximum"]); exclusive != nil && value >= *exclusive {
		value = *exclusive - 1
	} else if rawBool(schema["exclusiveMaximum"]) {
		if maximum := rawFloat(schema["maximum"]); maximum != nil && value >= *maximum {
			value--
		}
	}
	if integer {
		value = math.Ceil(value)
		return json.Marshal(int64(value))
	}
	return json.Marshal(value)
}

func rawInteger(raw json.RawMessage) int {
	var value int
	_ = json.Unmarshal(raw, &value)
	return value
}

func rawFloat(raw json.RawMessage) *float64 {
	if len(raw) == 0 {
		return nil
	}
	var value float64
	if json.Unmarshal(raw, &value) != nil {
		return nil
	}
	return &value
}

func rawBool(raw json.RawMessage) bool {
	var value bool
	_ = json.Unmarshal(raw, &value)
	return value
}

func rawString(raw json.RawMessage) string {
	var value string
	_ = json.Unmarshal(raw, &value)
	return value
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func sampleMatchesPattern(pattern, value string) bool {
	compiled, err := regexp2.Compile(pattern, regexp2.None)
	if err != nil {
		return false
	}
	matched, err := compiled.MatchString(value)
	return err == nil && matched
}

func resolveReference(raw json.RawMessage, components map[string]json.RawMessage, section string) (json.RawMessage, error) {
	var reference struct {
		Ref string `json:"$ref"`
	}
	if err := json.Unmarshal(raw, &reference); err != nil {
		return nil, err
	}
	if reference.Ref == "" {
		return raw, nil
	}
	prefix := "#/components/" + section + "/"
	if !strings.HasPrefix(reference.Ref, prefix) {
		return nil, fmt.Errorf("unsupported external or unexpected reference %q", reference.Ref)
	}
	name := strings.TrimPrefix(reference.Ref, prefix)
	name = strings.ReplaceAll(strings.ReplaceAll(name, "~1", "/"), "~0", "~")
	resolved, ok := components[name]
	if !ok {
		return nil, fmt.Errorf("unresolved reference %q", reference.Ref)
	}
	return resolved, nil
}

func resolveSchemaReference(raw json.RawMessage, document specDocument) (json.RawMessage, error) {
	resolved, err := resolveReference(raw, document.Components.Schemas, "schemas")
	if err != nil {
		return nil, err
	}
	var source map[string]json.RawMessage
	if err := json.Unmarshal(raw, &source); err != nil {
		return nil, err
	}
	if _, hasReference := source["$ref"]; !hasReference {
		return raw, nil
	}
	delete(source, "$ref")
	if len(source) == 0 {
		return resolved, nil
	}
	var target map[string]json.RawMessage
	if err := json.Unmarshal(resolved, &target); err != nil {
		return nil, err
	}
	for key, value := range source {
		target[key] = value
	}
	return json.Marshal(target)
}

func loadClientMethods(directory string) ([]clientMethod, error) {
	paths, err := filepath.Glob(filepath.Join(directory, "*_client.gen.go"))
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no generated client files found in %s", directory)
	}
	methods := make([]clientMethod, 0)
	for _, path := range paths {
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		imports := make(map[string]string)
		for _, imported := range file.Imports {
			importPath, err := strconv.Unquote(imported.Path.Value)
			if err != nil {
				return nil, err
			}
			alias := filepath.Base(importPath)
			if imported.Name != nil {
				alias = imported.Name.Name
			}
			if alias != "_" && alias != "." {
				imports[alias] = importPath
			}
		}
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Recv == nil || len(function.Recv.List) != 1 || function.Doc == nil {
				continue
			}
			if receiverType(function.Recv.List[0].Type) != "Client" {
				continue
			}
			for _, comment := range function.Doc.List {
				line := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
				match := operationComment.FindStringSubmatch(line)
				if match == nil {
					continue
				}
				methods = append(methods, clientMethod{
					Name:       function.Name.Name,
					Operation:  lowerFirst(match[3]),
					Method:     match[1],
					Path:       match[2],
					Params:     function.Type.Params.List,
					Imports:    imports,
					SourceFile: path,
				})
				break
			}
		}
	}
	return methods, nil
}

func receiverType(expression ast.Expr) string {
	if pointer, ok := expression.(*ast.StarExpr); ok {
		expression = pointer.X
	}
	if identifier, ok := expression.(*ast.Ident); ok {
		return identifier.Name
	}
	return ""
}

func renderAdapters(operations []operation, methods []clientMethod, spec []byte, specVersion string) ([]byte, error) {
	byKey := make(map[string][]clientMethod)
	for _, method := range methods {
		key := operationKey(method.Method, method.Path)
		byKey[key] = append(byKey[key], method)
	}
	operationKeys := make(map[string]struct{}, len(operations))
	for _, item := range operations {
		operationKeys[operationKey(item.Method, item.Path)] = struct{}{}
	}
	for key := range byKey {
		if _, exists := operationKeys[key]; !exists {
			return nil, fmt.Errorf("generated client has operation absent from spec: %s; verify the canonical spec, regenerate the SDK with `mise run generate`, then review request adapters", key)
		}
	}

	externalImports := make(map[string]string)
	var source strings.Builder
	source.WriteString("// Code generated by test/requestconformance/generate; DO NOT EDIT.\n")
	source.WriteString("package requestconformance\n\n")
	source.WriteString("import (\n\t\"context\"\n\t\"net/http\"\n")
	source.WriteString("\tupcloud \"github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud\"\n")
	for _, item := range operations {
		candidates := byKey[operationKey(item.Method, item.Path)]
		selected, err := selectMethod(item, candidates)
		if err != nil {
			return nil, err
		}
		for _, field := range selected.Params {
			if isContextParameter(field.Type) || isRequestEditorParameter(field.Type) {
				continue
			}
			if isReaderParameter(field.Type, selected.Imports) {
				continue
			}
			if _, err := renderType(field.Type, selected.Imports, externalImports); err != nil {
				return nil, fmt.Errorf("%s argument type: %w", selected.Name, err)
			}
		}
	}
	aliases := make([]string, 0, len(externalImports))
	for alias := range externalImports {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		importPath := externalImports[alias]
		if alias == filepath.Base(importPath) {
			fmt.Fprintf(&source, "\t%q\n", importPath)
		} else {
			fmt.Fprintf(&source, "\t%s %q\n", alias, importPath)
		}
	}
	source.WriteString(")\n\n")
	digest := sha256.Sum256(spec)
	fmt.Fprintf(&source, "const requestSweepSpecVersion = %s\n", strconv.Quote(specVersion))
	fmt.Fprintf(&source, "const requestSweepSpecSHA256 = %s\n\n", strconv.Quote(hex.EncodeToString(digest[:])))
	fmt.Fprintf(&source, "const requestSweepOperationCount = %d\n\n", len(operations))
	source.WriteString("var operationAdapters = []operationAdapter{\n")
	for _, item := range operations {
		selected, err := selectMethod(item, byKey[operationKey(item.Method, item.Path)])
		if err != nil {
			return nil, err
		}
		pathArguments, err := matchPathArguments(item.PathParams, selected.Params)
		if err != nil {
			return nil, fmt.Errorf("%s %s: %w", item.Method, item.Path, err)
		}
		fmt.Fprintf(&source, "{ID: %s, Method: %s, PathTemplate: %s,\n",
			strconv.Quote(item.ID), strconv.Quote(item.Method), strconv.Quote(item.Path))
		source.WriteString("PathParams: []pathContract{")
		for _, pathParam := range item.PathParams {
			fmt.Fprintf(&source, "{Name:%s, Style:%s, Explode:%t, SchemaType:%s, Format:%s, SampleJSON:%s},",
				strconv.Quote(pathParam.Name), strconv.Quote(pathParam.Style), pathParam.Explode,
				strconv.Quote(pathParam.SchemaType), strconv.Quote(pathParam.Format), strconv.Quote(pathParam.SampleJSON))
		}
		source.WriteString("}, Query: []queryContract{\n")
		for _, query := range item.Query {
			fmt.Fprintf(&source, "{Name:%s, Style:%s, Explode:%t, Required:%t, SchemaType:%s, Properties:%s, SampleJSON:%s},\n",
				strconv.Quote(query.Name), strconv.Quote(query.Style), query.Explode, query.Required,
				strconv.Quote(query.SchemaType), quotedSlice(query.Properties), strconv.Quote(query.SampleJSON))
		}
		source.WriteString("}, Headers: []headerContract{")
		for _, header := range item.Headers {
			fmt.Fprintf(&source, "{Name:%s, Required:%t},", strconv.Quote(header.Name), header.Required)
		}
		fmt.Fprintf(&source, "}, BodyRequired:%t, BodyMediaTypes:%s, BodySampleJSON:%s, Invoke: func(ctx context.Context, client *upcloud.Client, observed *capturedRequest, bodySampleOverride string) (*http.Response, error) {\n",
			item.BodyRequired, quotedSlice(item.BodyMediaType), strconv.Quote(item.BodySampleJSON))
		source.WriteString("\t\tobserved.expectedQuery = make(map[string][]string)\n\t\tobserved.hasExpectedQuery = true\n")
		for _, field := range selected.Params {
			if isContextParameter(field.Type) || isRequestEditorParameter(field.Type) {
				continue
			}
			typeName, err := renderType(field.Type, selected.Imports, nil)
			if err != nil {
				return nil, err
			}
			for _, name := range field.Names {
				if strings.EqualFold(name.Name, "contentType") {
					fmt.Fprintf(&source, "\t\t%s := requestMediaType(%s)\n", name.Name, quotedSlice(item.BodyMediaType))
				} else if strings.Contains(strings.ToLower(name.Name), "body") {
					fmt.Fprintf(&source, "\t\tbodySample := bodySampleOverride\n\t\tif bodySample == \"\" { bodySample = %s }\n", strconv.Quote(item.BodySampleJSON))
					if isReaderParameter(field.Type, selected.Imports) {
						fmt.Fprintf(&source, "\t\t%s := requestSweepBodyReader(observed, bodySample)\n", name.Name)
					} else {
						fmt.Fprintf(&source, "\t\t%s, err := requestSampleFromJSON[%s](bodySample)\n\t\tif err != nil { return nil, err }\n", name.Name, typeName)
						source.WriteString("\t\tobserved.expectedBody = []byte(bodySample)\n\t\tobserved.hasExpectedBody = true\n")
					}
				} else if isParamsParameter(name.Name, field.Type) {
					fmt.Fprintf(&source, "\t\t%s, err := requestSampleParams[%s]([]queryContract{\n", name.Name, typeName)
					for _, query := range item.Query {
						fmt.Fprintf(&source, "{Name:%s, Style:%s, Explode:%t, Required:%t, SchemaType:%s, Properties:%s, SampleJSON:%s},\n",
							strconv.Quote(query.Name), strconv.Quote(query.Style), query.Explode, query.Required,
							strconv.Quote(query.SchemaType), quotedSlice(query.Properties), strconv.Quote(query.SampleJSON))
					}
					source.WriteString("})\n\t\tif err != nil { return nil, err }\n")
					source.WriteString("\t\texpectedQuery, err := expectedQueryFromParams(")
					source.WriteString(name.Name)
					source.WriteString(", []queryContract{\n")
					for _, query := range item.Query {
						fmt.Fprintf(&source, "{Name:%s, Style:%s, Explode:%t, Required:%t, SchemaType:%s, Properties:%s, SampleJSON:%s},\n",
							strconv.Quote(query.Name), strconv.Quote(query.Style), query.Explode, query.Required,
							strconv.Quote(query.SchemaType), quotedSlice(query.Properties), strconv.Quote(query.SampleJSON))
					}
					source.WriteString("})\n\t\tif err != nil { return nil, err }\n\t\tobserved.expectedQuery = expectedQuery\n")
				} else if pathSample, isPathParameter := pathParameterByGoArgument(name.Name, pathArguments, item.PathParams); isPathParameter {
					fmt.Fprintf(&source, "\t\t%s, err := requestSampleFromJSON[%s](%s)\n\t\tif err != nil { return nil, err }\n", name.Name, typeName, strconv.Quote(pathSample))
				} else {
					fmt.Fprintf(&source, "\t\t%s := requestSample[%s]()\n", name.Name, typeName)
				}
			}
		}
		fmt.Fprintf(&source, "\t\texpectedPath, pathErr := requestPathFromTemplate(%s, []pathValue{\n", strconv.Quote(item.Path))
		for _, pathParam := range item.PathParams {
			fmt.Fprintf(&source, "{Contract:pathContract{Name:%s, Style:%s, Explode:%t, SchemaType:%s, Format:%s}, Value:requestSampleValueFromJSON(%s)},\n",
				strconv.Quote(pathParam.Name), strconv.Quote(pathParam.Style), pathParam.Explode,
				strconv.Quote(pathParam.SchemaType), strconv.Quote(pathParam.Format), strconv.Quote(pathParam.SampleJSON))
		}
		source.WriteString("})\n\t\tif pathErr != nil { return nil, pathErr }\n\t\tobserved.expectedPath = expectedPath\n\t\tobserved.hasExpectedPath = true\n")
		fmt.Fprintf(&source, "\t\treturn client.%s(ctx", selected.Name)
		for _, field := range selected.Params {
			if isContextParameter(field.Type) || isRequestEditorParameter(field.Type) {
				continue
			}
			for _, name := range field.Names {
				fmt.Fprintf(&source, ", %s", name.Name)
			}
		}
		source.WriteString(")\n\t}},\n")
	}
	source.WriteString("}\n")
	formatted, err := format.Source([]byte(source.String()))
	if err != nil {
		return nil, fmt.Errorf("format generated adapters: %w\n%s", err, source.String())
	}
	return formatted, nil
}

func selectMethod(item operation, candidates []clientMethod) (clientMethod, error) {
	matching := make([]clientMethod, 0, len(candidates))
	for _, method := range candidates {
		if method.Operation == item.ID {
			matching = append(matching, method)
		}
	}
	if len(matching) == 0 {
		return clientMethod{}, fmt.Errorf("no generated *Client method for %s %s (operationId %s); verify the canonical spec and generated SDK, then run `make generate-request-adapters`", item.Method, item.Path, item.ID)
	}
	wantName := upperFirst(item.ID)
	for _, method := range matching {
		if method.Name == wantName {
			return method, nil
		}
	}
	sort.Slice(matching, func(i, j int) bool {
		leftBody := strings.HasSuffix(matching[i].Name, "WithBody")
		rightBody := strings.HasSuffix(matching[j].Name, "WithBody")
		if leftBody != rightBody {
			return !leftBody
		}
		return matching[i].Name < matching[j].Name
	})
	return matching[0], nil
}

func matchPathArguments(pathParams []queryContract, params []*ast.Field) (map[string]string, error) {
	result := make(map[string]string, len(pathParams))
	for _, pathParam := range pathParams {
		wanted := normalizeIdentifier(pathParam.Name)
		for _, field := range params {
			if isContextParameter(field.Type) || isRequestEditorParameter(field.Type) {
				continue
			}
			for _, name := range field.Names {
				keywordEscapedName := normalizeIdentifier("p" + upperFirst(pathParam.Name))
				if normalizeIdentifier(name.Name) == wanted || (token.Lookup(pathParam.Name).IsKeyword() && normalizeIdentifier(name.Name) == keywordEscapedName) {
					result[pathParam.Name] = name.Name
					break
				}
			}
			if result[pathParam.Name] != "" {
				break
			}
		}
		if result[pathParam.Name] == "" {
			return nil, fmt.Errorf("no typed client argument maps to path parameter %q; update the adapter generator mapping", pathParam.Name)
		}
	}
	return result, nil
}

func pathParameterByGoArgument(argument string, mapping map[string]string, pathParams []queryContract) (string, bool) {
	for _, parameter := range pathParams {
		if mapping[parameter.Name] == argument {
			return parameter.SampleJSON, true
		}
	}
	return "", false
}

func normalizeIdentifier(value string) string {
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, "-", "")
	return value
}

func renderType(expression ast.Expr, imports map[string]string, generatedImports map[string]string) (string, error) {
	switch value := expression.(type) {
	case *ast.Ident:
		if isBuiltin(value.Name) {
			return value.Name, nil
		}
		return "upcloud." + value.Name, nil
	case *ast.SelectorExpr:
		if identifier, ok := value.X.(*ast.Ident); ok {
			if path, exists := imports[identifier.Name]; exists {
				if generatedImports != nil {
					if oldPath, already := generatedImports[identifier.Name]; already && oldPath != path {
						return "", fmt.Errorf("import alias %s resolves to both %s and %s", identifier.Name, oldPath, path)
					}
					generatedImports[identifier.Name] = path
				}
				return identifier.Name + "." + value.Sel.Name, nil
			}
		}
		left, err := renderType(value.X, imports, generatedImports)
		if err != nil {
			return "", err
		}
		return left + "." + value.Sel.Name, nil
	case *ast.StarExpr:
		inner, err := renderType(value.X, imports, generatedImports)
		return "*" + inner, err
	case *ast.ArrayType:
		inner, err := renderType(value.Elt, imports, generatedImports)
		if err != nil {
			return "", err
		}
		if value.Len == nil {
			return "[]" + inner, nil
		}
		var length bytes.Buffer
		if err := format.Node(&length, token.NewFileSet(), value.Len); err != nil {
			return "", err
		}
		return "[" + length.String() + "]" + inner, nil
	case *ast.MapType:
		key, err := renderType(value.Key, imports, generatedImports)
		if err != nil {
			return "", err
		}
		item, err := renderType(value.Value, imports, generatedImports)
		return "map[" + key + "]" + item, err
	case *ast.InterfaceType:
		return "any", nil
	case *ast.ParenExpr:
		inner, err := renderType(value.X, imports, generatedImports)
		return "(" + inner + ")", err
	default:
		var rendered bytes.Buffer
		if err := format.Node(&rendered, token.NewFileSet(), expression); err != nil {
			return "", err
		}
		return rendered.String(), nil
	}
}

func isContextParameter(expression ast.Expr) bool {
	selector, ok := expression.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Context" {
		return false
	}
	identifier, ok := selector.X.(*ast.Ident)
	return ok && identifier.Name == "context"
}

func isRequestEditorParameter(expression ast.Expr) bool {
	ellipsis, ok := expression.(*ast.Ellipsis)
	if !ok {
		return false
	}
	identifier, ok := ellipsis.Elt.(*ast.Ident)
	return ok && identifier.Name == "RequestEditorFn"
}

func isReaderParameter(expression ast.Expr, imports map[string]string) bool {
	selector, ok := expression.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Reader" {
		return false
	}
	identifier, ok := selector.X.(*ast.Ident)
	return ok && imports[identifier.Name] == "io"
}

func isParamsParameter(name string, expression ast.Expr) bool {
	if strings.EqualFold(name, "params") {
		return true
	}
	identifier, ok := expression.(*ast.Ident)
	return ok && strings.HasSuffix(identifier.Name, "Params")
}

func isBuiltin(name string) bool {
	switch name {
	case "any", "bool", "byte", "comparable", "complex64", "complex128", "error", "float32", "float64", "int", "int8", "int16", "int32", "int64", "rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
		return true
	default:
		return false
	}
}

func quotedSlice(values []string) string {
	quoted := make([]string, len(values))
	for index, value := range values {
		quoted[index] = strconv.Quote(value)
	}
	return "[]string{" + strings.Join(quoted, ",") + "}"
}

func operationKey(method, path string) string {
	return strings.ToUpper(method) + " " + path
}

func isHTTPMethod(method string) bool {
	switch method {
	case "get", "post", "put", "patch", "delete", "head", "options", "trace":
		return true
	default:
		return false
	}
}

func upperFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
