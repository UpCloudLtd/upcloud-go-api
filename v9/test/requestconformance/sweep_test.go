package requestconformance

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/stretchr/testify/require"
)

const (
	schemaTypeObject     = "object"
	parameterStyleLabel  = "label"
	parameterStyleMatrix = "matrix"
)

type queryContract struct {
	Name       string
	Style      string
	Explode    bool
	Required   bool
	SchemaType string
	Properties []string
	SampleJSON string
}

type operationAdapter struct {
	ID             string
	Method         string
	PathTemplate   string
	PathParams     []pathContract
	Query          []queryContract
	Headers        []headerContract
	BodyRequired   bool
	BodyMediaTypes []string
	BodySampleJSON string
	Invoke         func(context.Context, *upcloud.Client, *capturedRequest, string) (*http.Response, error)
}

type pathContract struct {
	Name       string
	Style      string
	Explode    bool
	SchemaType string
	Format     string
	SampleJSON string
}

type pathValue struct {
	Contract pathContract
	Value    any
}

type headerContract struct {
	Name     string
	Required bool
}

type capturedRequest struct {
	request          *http.Request
	body             []byte
	expectedBody     []byte
	hasExpectedBody  bool
	expectedPath     string
	hasExpectedPath  bool
	expectedQuery    map[string][]string
	hasExpectedQuery bool
	calls            int
}

type operationInventoryEntry struct {
	ID             string
	Method         string
	PathTemplate   string
	PathParams     []pathContract
	Query          []queryContract
	Headers        []headerContract
	BodyRequired   bool
	BodyMediaTypes []string
	BodySampleJSON string
}

type requestCaseReport struct {
	ID              string              `json:"id"`
	Method          string              `json:"method"`
	PathTemplate    string              `json:"path_template"`
	ExpectedPath    string              `json:"expected_path"`
	ActualPath      string              `json:"actual_path"`
	ExpectedQuery   map[string][]string `json:"expected_query"`
	ActualQuery     map[string][]string `json:"actual_query"`
	ExpectedHeaders map[string][]string `json:"expected_headers"`
	ActualHeaders   map[string][]string `json:"actual_headers"`
	ExpectedBody    string              `json:"expected_body"`
	ActualBody      string              `json:"actual_body"`
	ContentType     string              `json:"content_type"`
	StatusCode      int                 `json:"status_code"`
	SendCount       int                 `json:"send_count"`
	InvocationError string              `json:"invocation_error,omitempty"`
	ValidationError string              `json:"validation_error,omitempty"`
	Passed          bool                `json:"passed"`
}

type requestSweepReport struct {
	SpecVersion    string                    `json:"spec_version"`
	SpecSHA256     string                    `json:"spec_sha256"`
	Normalizations []string                  `json:"validator_normalizations"`
	Inventory      []operationInventoryEntry `json:"operation_inventory"`
	Cases          []requestCaseReport       `json:"cases"`
	Covered        int                       `json:"covered"`
	Total          int                       `json:"total"`
	Missing        []string                  `json:"missing"`
	Extra          []string                  `json:"extra"`
}

type captureDoer struct {
	captured capturedRequest
}

func (doer *captureDoer) Do(request *http.Request) (*http.Response, error) {
	body, err := readAndCloseRequestBody(request)
	if err != nil {
		return nil, err
	}
	if request.Body != nil {
		request.Body = io.NopCloser(bytes.NewReader(body))
	}
	previous := doer.captured
	doer.captured = capturedRequest{
		request:          request,
		body:             body,
		expectedBody:     previous.expectedBody,
		hasExpectedBody:  previous.hasExpectedBody,
		expectedPath:     previous.expectedPath,
		hasExpectedPath:  previous.hasExpectedPath,
		expectedQuery:    previous.expectedQuery,
		hasExpectedQuery: previous.hasExpectedQuery,
		calls:            previous.calls + 1,
	}
	return &http.Response{
		StatusCode: http.StatusNoContent,
		Status:     "204 No Content",
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    request,
	}, nil
}

func TestEverySpecOperationEmitsOneTypedRequest(t *testing.T) {
	require.NotZero(t, requestSweepOperationCount, "canonical operation inventory must not be empty")
	require.Len(t, operationAdapters, requestSweepOperationCount, "generated adapter inventory is incomplete")
	specDocument := loadRequestSpec(t)
	specBytes, err := os.ReadFile("../../generator/openapi/spec.json")
	require.NoError(t, err)
	var spec struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	require.NoError(t, json.Unmarshal(specBytes, &spec))
	require.Equal(t, spec.Info.Version, requestSweepSpecVersion)
	require.Equal(t, requestSweepSpecSHA256, specDigest(specBytes))

	expectedCases := make(map[string]struct{}, len(operationAdapters))
	var executedCases sync.Map
	caseReports := make([]requestCaseReport, 0, len(operationAdapters))
	report := requestSweepReport{
		SpecVersion: requestSweepSpecVersion,
		SpecSHA256:  requestSweepSpecSHA256,
		Normalizations: []string{
			"databaseServiceCloneOpenAPI: moved additionalProperties:false to unevaluatedProperties:false to preserve oneOf branch properties while rejecting unknown fields",
			"objectStorage2TagKey: replaced a negative-lookahead pattern with an equivalent RE2 pattern for kin-openapi validation",
		},
		Inventory: make([]operationInventoryEntry, 0, len(operationAdapters)),
	}
	for _, adapter := range operationAdapters {
		key := adapter.Method + " " + adapter.PathTemplate
		require.NotContains(t, expectedCases, key, "duplicate adapter case")
		expectedCases[key] = struct{}{}
		report.Inventory = append(report.Inventory, operationInventoryEntry{
			ID: adapter.ID, Method: adapter.Method, PathTemplate: adapter.PathTemplate,
			PathParams: adapter.PathParams, Query: adapter.Query, Headers: adapter.Headers,
			BodyRequired: adapter.BodyRequired, BodyMediaTypes: adapter.BodyMediaTypes, BodySampleJSON: adapter.BodySampleJSON,
		})
	}
	for _, adapter := range operationAdapters {
		t.Run(adapter.ID, func(t *testing.T) {
			doer := &captureDoer{}
			client, err := upcloud.New(
				"request-sweep-synthetic-token",
				upcloud.WithBaseURL("https://api.example.test/request-sweep"),
				upcloud.WithHTTPClient(doer),
			)
			require.NoError(t, err)
			response, err := adapter.Invoke(t.Context(), client.Client, &doer.captured, curatedRequestBodySamples[adapter.ID])
			caseReports = append(caseReports, buildCaseReport(adapter, doer.captured, response, err))
			caseIndex := len(caseReports) - 1
			require.NoError(t, err, "operation %s", adapter.ID)
			require.NotNil(t, response, "operation %s returned a nil HTTP response", adapter.ID)
			require.NoError(t, response.Body.Close(), "operation %s response body close", adapter.ID)

			require.Equal(t, 1, doer.captured.calls, "operation %s send count", adapter.ID)
			request := doer.captured.request
			require.NotNil(t, request, "operation %s captured no request", adapter.ID)
			require.Equal(t, adapter.Method, request.Method, "operation %s HTTP method", adapter.ID)
			require.Equal(t, "https", request.URL.Scheme, "operation %s scheme", adapter.ID)
			require.Equal(t, "api.example.test", request.URL.Host, "operation %s request host", adapter.ID)
			require.True(t, doer.captured.hasExpectedPath, "operation %s has no expected path", adapter.ID)
			require.Equal(t, doer.captured.expectedPath, request.URL.EscapedPath(), "operation %s full escaped path", adapter.ID)
			if validationErr := validateRequestAgainstSpec(t.Context(), specDocument, adapter, request); validationErr != nil {
				caseReports[caseIndex].ValidationError = validationErr.Error()
				require.NoError(t, validationErr, "operation %s typed sample violates canonical OpenAPI constraints", adapter.ID)
			}
			require.Equal(t, "Bearer request-sweep-synthetic-token", request.Header.Get("Authorization"), "operation %s auth", adapter.ID)
			require.True(t, strings.HasSuffix(request.Header.Get("User-Agent"), "openapi/"+requestSweepSpecVersion), "operation %s user-agent", adapter.ID)
			assertQueryContract(t, adapter.ID, adapter.Query, doer.captured, request.URL.Query())
			assertHeaderContract(t, adapter.ID, adapter.Headers, request.Header)
			assertBodyContract(t, adapter, request, doer.captured)
			caseReports[caseIndex].Passed = true
			executedCases.Store(adapter.Method+" "+adapter.PathTemplate, struct{}{})
		})
	}
	missing := make([]string, 0)
	for key := range expectedCases {
		if _, executed := executedCases.Load(key); !executed {
			missing = append(missing, key)
		}
	}
	extra := make([]string, 0)
	executedCount := 0
	executedCases.Range(func(key, _ any) bool {
		executedCount++
		operationKey, ok := key.(string)
		if !ok {
			t.Errorf("unexpected request case key type %T", key)
			return true
		}
		if _, expected := expectedCases[operationKey]; !expected {
			extra = append(extra, operationKey)
		}
		return true
	})
	sort.Strings(missing)
	sort.Strings(extra)
	report.Cases = caseReports
	report.Missing = missing
	report.Extra = extra
	report.Total = len(expectedCases)
	report.Covered = executedCount
	writeRequestSweepReport(t, report)
	require.Empty(t, missing, "request sweep skipped operation cases: %v", missing)
	require.Empty(t, extra, "request sweep executed unexpected operation cases: %v", extra)
	require.Equal(t, len(expectedCases), executedCount, "every request adapter must execute")
	t.Logf("request sweep covered=%d total=%d missing=%d extra=%d", executedCount, len(expectedCases), len(missing), len(extra))
}

var curatedRequestBodySamples = map[string]string{
	"modifyLoadBalancerCertificateBundle": `{"name":"request-sweep-cert","private_key":"synthetic-private-key"}`,
	"cloneDatabase":                       `{"hostname_prefix":"request-sweep-db","plan":"starter-kit-plan","properties":{"version":"8.4"}}`,
	"createLoadBalancerBackend":           `{"name":"request-sweep-backend","members":[{"name":"request-sweep-member","weight":100,"max_sessions":100,"type":"dynamic","enabled":true}]}`,
	"createLoadBalancerMember":            `{"name":"request-sweep-member","weight":100,"max_sessions":100,"type":"dynamic","enabled":true}`,
}

func assertHeaderContract(t *testing.T, operationID string, contracts []headerContract, actual http.Header) {
	t.Helper()
	allowed := map[string]struct{}{
		"authorization": {},
		"content-type":  {},
		"user-agent":    {},
	}
	for _, contract := range contracts {
		allowed[strings.ToLower(contract.Name)] = struct{}{}
		if contract.Required {
			require.NotEmpty(t, actual.Values(contract.Name), "operation %s missing required header %q", operationID, contract.Name)
		}
	}
	for name := range actual {
		_, declared := allowed[strings.ToLower(name)]
		require.True(t, declared, "operation %s sent unexpected request header %q", operationID, name)
	}
}

func buildCaseReport(adapter operationAdapter, captured capturedRequest, response *http.Response, invokeErr error) requestCaseReport {
	result := requestCaseReport{
		ID:              adapter.ID,
		Method:          adapter.Method,
		PathTemplate:    adapter.PathTemplate,
		ExpectedQuery:   copyQueryValues(captured.expectedQuery),
		ActualQuery:     make(map[string][]string),
		ExpectedHeaders: map[string][]string{"Authorization": {"Bearer [REDACTED]"}, "User-Agent": {"upcloud-go-api/v9 openapi/" + requestSweepSpecVersion}},
		ActualHeaders:   make(map[string][]string),
		ExpectedBody:    string(captured.expectedBody),
		ActualBody:      string(captured.body),
		SendCount:       captured.calls,
	}
	if captured.hasExpectedPath {
		result.ExpectedPath = "https://api.example.test" + captured.expectedPath
	}
	if captured.request != nil {
		result.ActualPath = captured.request.URL.Scheme + "://" + captured.request.URL.Host + captured.request.URL.EscapedPath()
		result.ActualQuery = copyQueryValues(captured.request.URL.Query())
		result.ContentType = captured.request.Header.Get("Content-Type")
		result.ActualHeaders = redactHeaders(captured.request.Header)
	}
	if len(captured.expectedBody) > 0 && len(adapter.BodyMediaTypes) > 0 {
		result.ExpectedHeaders["Content-Type"] = []string{requestMediaType(adapter.BodyMediaTypes)}
	}
	if response != nil {
		result.StatusCode = response.StatusCode
	}
	if invokeErr != nil {
		result.InvocationError = invokeErr.Error()
	}
	return result
}

func copyQueryValues(values map[string][]string) map[string][]string {
	result := make(map[string][]string, len(values))
	for name, items := range values {
		result[name] = append([]string(nil), items...)
	}
	return result
}

func redactHeaders(headers http.Header) map[string][]string {
	result := make(map[string][]string, len(headers))
	for name, values := range headers {
		if strings.EqualFold(name, "Authorization") {
			result[http.CanonicalHeaderKey(name)] = []string{"Bearer [REDACTED]"}
			continue
		}
		result[http.CanonicalHeaderKey(name)] = append([]string(nil), values...)
	}
	return result
}

func loadRequestSpec(t *testing.T) *openapi3.T {
	t.Helper()
	loader := openapi3.NewLoader()
	document, err := loader.LoadFromFile("../../generator/openapi/spec.json")
	require.NoError(t, err, "load canonical OpenAPI spec for request validation")
	document.Servers = openapi3.Servers{{URL: "https://api.example.test/request-sweep"}}
	require.NoError(t, normalizeOpenAPIPatterns(document), "normalize unsupported schema regexes for the Go validator")
	require.NoError(t, normalizeDatabaseCloneSchema(document), "normalize contradictory database clone schema composition")
	return document
}

func normalizeDatabaseCloneSchema(document *openapi3.T) error {
	// The copied 3.1 schema applies additionalProperties:false before its oneOf
	// branches, rejecting every branch-only property. unevaluatedProperties keeps
	// the intended rejection after branch evaluation.
	const schemaName = "databaseServiceCloneOpenAPI"
	reference := document.Components.Schemas[schemaName]
	if reference == nil || reference.Value == nil {
		return fmt.Errorf("canonical schema %s is missing", schemaName)
	}
	schema := reference.Value
	if schema.AdditionalProperties.Has != nil && !*schema.AdditionalProperties.Has && len(schema.OneOf) > 0 {
		if schema.UnevaluatedProperties.Has != nil {
			return fmt.Errorf("canonical schema %s has both additionalProperties and unevaluatedProperties", schemaName)
		}
		schema.UnevaluatedProperties = openapi3.BoolSchema{Has: schema.AdditionalProperties.Has}
		schema.AdditionalProperties = openapi3.AdditionalProperties{}
		return nil
	}
	if schema.UnevaluatedProperties.Has != nil && !*schema.UnevaluatedProperties.Has && len(schema.OneOf) > 0 {
		return nil
	}
	return fmt.Errorf("canonical schema %s no longer has the expected oneOf/additionalProperties composition", schemaName)
}

func normalizeOpenAPIPatterns(document *openapi3.T) error {
	// Go's RE2 engine lacks negative lookahead; this exact rewrite excludes '_' only
	// from the first character and keeps the same remaining alphabet and suffix.
	const source = `^(?!_)[\x20-\x3C\x3E-\x7E]{2,32}(?:=.*)?$`
	const replacement = `^[\x20-\x3C\x3E-\x5E\x60-\x7E][\x20-\x3C\x3E-\x7E]{1,31}(=.*)?$`
	seen := make(map[*openapi3.Schema]struct{})
	var visit func(*openapi3.SchemaRef) error
	visit = func(reference *openapi3.SchemaRef) error {
		if reference == nil || reference.Value == nil {
			return nil
		}
		schema := reference.Value
		if _, exists := seen[schema]; exists {
			return nil
		}
		seen[schema] = struct{}{}
		if schema.Pattern != "" {
			if _, err := regexp.Compile(schema.Pattern); err != nil {
				if schema.Pattern != source {
					return fmt.Errorf("unsupported schema regex %q; add a reviewed equivalent", schema.Pattern)
				}
				schema.Pattern = replacement
			}
		}
		for _, references := range []openapi3.SchemaRefs{schema.OneOf, schema.AnyOf, schema.AllOf, schema.PrefixItems} {
			for _, child := range references {
				if err := visit(child); err != nil {
					return err
				}
			}
		}
		for _, child := range schema.Properties {
			if err := visit(child); err != nil {
				return err
			}
		}
		for _, child := range schema.PatternProperties {
			if err := visit(child); err != nil {
				return err
			}
		}
		for _, children := range []openapi3.Schemas{schema.DependentSchemas, schema.Defs} {
			for _, child := range children {
				if err := visit(child); err != nil {
					return err
				}
			}
		}
		for _, child := range []*openapi3.SchemaRef{schema.Not, schema.Items, schema.PropertyNames, schema.Contains, schema.If, schema.Then, schema.Else, schema.ContentSchema, schema.AdditionalProperties.Schema} {
			if err := visit(child); err != nil {
				return err
			}
		}
		return nil
	}
	for _, schema := range document.Components.Schemas {
		if err := visit(schema); err != nil {
			return err
		}
	}
	return nil
}

func validateRequestAgainstSpec(ctx context.Context, document *openapi3.T, adapter operationAdapter, request *http.Request) error {
	pathItem := document.Paths.Find(adapter.PathTemplate)
	if pathItem == nil {
		return fmt.Errorf("canonical spec has no path %s", adapter.PathTemplate)
	}
	operation := pathItem.GetOperation(adapter.Method)
	if operation == nil {
		return fmt.Errorf("canonical spec has no %s operation at %s", adapter.Method, adapter.PathTemplate)
	}
	pathParams := make(map[string]string, len(adapter.PathParams))
	for _, parameter := range adapter.PathParams {
		pathParams[parameter.Name] = fmt.Sprint(requestSampleValueFromJSON(parameter.SampleJSON))
	}
	route := &routers.Route{
		Spec:      document,
		Path:      adapter.PathTemplate,
		PathItem:  pathItem,
		Method:    adapter.Method,
		Operation: operation,
	}
	err := openapi3filter.ValidateRequest(ctx, &openapi3filter.RequestValidationInput{
		Request:    request,
		PathParams: pathParams,
		Route:      route,
		Options: &openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
			MultiError:         true,
			SchemaValidationOptions: []openapi3.SchemaValidationOption{
				openapi3.EnableFormatValidation(),
			},
		},
	})
	if request.Body != nil {
		_ = request.Body.Close()
	}
	return err
}

func writeRequestSweepReport(t *testing.T, report requestSweepReport) {
	t.Helper()
	directory := os.Getenv("V9_REQUEST_SWEEP_REPORT_DIR")
	if directory == "" {
		return
	}
	//nolint:gosec // CI sets this to its trusted runner temp directory for artifact upload.
	require.NoError(t, os.MkdirAll(directory, 0o700), "create request sweep report directory")
	contents, err := json.MarshalIndent(report, "", "  ")
	require.NoError(t, err, "encode request sweep report")
	reportPath := filepath.Join(directory, "request-sweep-report.json")
	//nolint:gosec // The filename is fixed and the directory is supplied by the trusted CI runner.
	require.NoError(t, os.WriteFile(reportPath, contents, 0o600), "write request sweep report")
}

func TestCuratedZeroEventPaginationValues(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	response, err := client.ListEvents(t.Context(), upcloud.ListEventsServerUuid{}, 0, 0)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.Equal(t, "GET", doer.captured.request.Method)
	require.Equal(t, "/request-sweep/1.3/events/00000000-0000-0000-0000-000000000000/0,0", doer.captured.request.URL.Path)
}

func TestCuratedPathParameterEscaping(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	address := upcloud.IpAddress{}
	require.NoError(t, address.FromIpAddress0("192.0.2.10/24"))
	response, err := client.DeleteIPAddress(t.Context(), address)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.Equal(t, "/request-sweep/1.3/ip_address/192.0.2.10%2F24", doer.captured.request.URL.EscapedPath())
}

func TestCuratedRouterLabelQueryEncoding(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	label := upcloud.ListRoutersLabel{Key: "env", Value: "production"}
	response, err := client.ListRouters(t.Context(), &upcloud.ListRoutersParams{Label: &label})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.Equal(t, url.Values{"key": []string{"env"}, "value": []string{"production"}}, doer.captured.request.URL.Query())
	require.Equal(t, "key=env&value=production", doer.captured.request.URL.RawQuery)
}

func TestCuratedRepeatedFiltersAndTimestampEncoding(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	createdAfter := time.Date(2026, time.January, 2, 3, 4, 5, 0, time.UTC)
	createdBefore := time.Date(2026, time.January, 3, 4, 5, 6, 0, time.UTC)
	actions := upcloud.ListAuditLogsActions{"server_start", "server_stop"}
	response, err := client.ListAuditLogs(t.Context(), &upcloud.ListAuditLogsParams{
		Actions:       &actions,
		CreatedAfter:  &createdAfter,
		CreatedBefore: &createdBefore,
	})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.Equal(t, url.Values{
		"actions":        {"server_start", "server_stop"},
		"created_after":  {"2026-01-02T03:04:05Z"},
		"created_before": {"2026-01-03T04:05:06Z"},
	}, doer.captured.request.URL.Query())
}

func TestCuratedRepeatedLabelFilters(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	labels := upcloud.ListLoadBalancersLabel{"env=production", "owner=platform"}
	response, err := client.ListLoadBalancers(t.Context(), &upcloud.ListLoadBalancersParams{Label: &labels})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.Equal(t, url.Values{"label": {"env=production", "owner=platform"}}, doer.captured.request.URL.Query())
}

func TestCuratedCommaListEncoding(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	tags := upcloud.ListServersTag("first,second")
	response, err := client.ListServers(t.Context(), &upcloud.ListServersParams{Tag: &tags})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.Equal(t, "first,second", doer.captured.request.URL.Query().Get("tag"))
	require.Equal(t, "tag=first%2Csecond", doer.captured.request.URL.RawQuery)
}

func TestCuratedOptionalQueryOmitted(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	response, err := client.ListServers(t.Context(), &upcloud.ListServersParams{})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.Empty(t, doer.captured.request.URL.RawQuery)
}

func TestCuratedFalseLegacyFlag(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	storages := upcloud.DeleteServerStorages("no")
	response, err := client.DeleteServer(t.Context(), upcloud.DeleteServerUuid{}, &upcloud.DeleteServerParams{Storages: &storages})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.Equal(t, "storages=no", doer.captured.request.URL.RawQuery)
}

func TestCuratedNestedWrapperAndEmptyCollectionBody(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	var body upcloud.CreateRouterJSONRequestBody
	body.Router.Name = "request-sweep-router"
	labels := []upcloud.RouterLabel{}
	body.Router.Labels = &labels
	response, err := client.CreateRouter(t.Context(), body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.JSONEq(t, `{"router":{"labels":[],"name":"request-sweep-router"}}`, string(doer.captured.body))
}

func TestCuratedTagCreateBodyGolden(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	var body upcloud.CreateTagJSONRequestBody
	body.Tag.Name = "QA_1"
	response, err := client.CreateTag(t.Context(), body)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.JSONEq(t, `{"tag":{"name":"QA_1","servers":{}}}`, string(doer.captured.body))
}

func TestCanonicalValidatorRejectsInvalidTagPattern(t *testing.T) {
	var body upcloud.CreateTagJSONRequestBody
	body.Tag.Name = "invalid tag name"
	request, err := upcloud.NewCreateTagRequest("https://api.example.test/request-sweep", body)
	require.NoError(t, err)

	err = validateRequestAgainstSpec(t.Context(), loadRequestSpec(t), operationAdapter{
		Method:       http.MethodPost,
		PathTemplate: "/1.3/tag",
	}, request)
	require.ErrorContains(t, err, "pattern")
}

func TestCuratedRequestEditorOverrideIsSeparateFromDefaultEncoding(t *testing.T) {
	doer := &captureDoer{}
	client := newRequestSweepClient(t, doer)
	var body upcloud.CreateRouterJSONRequestBody
	body.Router.Name = "default-typed-body"
	overrideBody := `{"router":{"name":"editor-override"}}`
	response, err := client.CreateRouter(t.Context(), body, func(_ context.Context, request *http.Request) error {
		request.Body = io.NopCloser(strings.NewReader(overrideBody))
		request.ContentLength = int64(len(overrideBody))
		request.GetBody = nil
		return nil
	})
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())
	require.Equal(t, 1, doer.captured.calls)
	require.JSONEq(t, overrideBody, string(doer.captured.body))
	require.NotContains(t, string(doer.captured.body), "default-typed-body")
}

func newRequestSweepClient(t *testing.T, doer *captureDoer) *upcloud.ClientWithResponses {
	t.Helper()
	client, err := upcloud.New(
		"request-sweep-synthetic-token",
		upcloud.WithBaseURL("https://api.example.test/request-sweep"),
		upcloud.WithHTTPClient(doer),
	)
	require.NoError(t, err)
	return client
}

func readAndCloseRequestBody(request *http.Request) ([]byte, error) {
	if request.Body == nil {
		return nil, nil
	}
	body, readErr := io.ReadAll(request.Body)
	closeErr := request.Body.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	return body, nil
}

func assertQueryContract(t *testing.T, operationID string, contracts []queryContract, captured capturedRequest, actual url.Values) {
	t.Helper()
	require.True(t, captured.hasExpectedQuery, "operation %s has no typed query expectation", operationID)
	for name := range actual {
		declared := false
		for _, contract := range contracts {
			if queryKeyMatches(contract, name) {
				declared = true
				break
			}
		}
		require.True(t, declared, "operation %s has unexpected query key %q in %v", operationID, name, actual)
	}
	for _, contract := range contracts {
		if contract.Required {
			found := false
			for name, values := range actual {
				if queryKeyMatches(contract, name) && len(values) > 0 {
					found = true
					break
				}
			}
			require.True(t, found, "operation %s is missing required query parameter %q in %v", operationID, contract.Name, actual)
		}
	}
	require.Equal(t, url.Values(captured.expectedQuery), actual, "operation %s full query multimap mismatch", operationID)
}

func expectedQueryFromParams(params any, contracts []queryContract) (map[string][]string, error) {
	expected := make(map[string][]string)
	value := reflect.ValueOf(params)
	for value.IsValid() && (value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer) {
		if value.IsNil() {
			return expected, nil
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return expected, nil
	}
	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("query params have type %s, want struct", value.Type())
	}
	contractByName := make(map[string]queryContract, len(contracts))
	for _, contract := range contracts {
		contractByName[contract.Name] = contract
	}
	for index := range value.NumField() {
		fieldType := value.Type().Field(index)
		formName := strings.Split(fieldType.Tag.Get("form"), ",")[0]
		if formName == "" || formName == "-" {
			continue
		}
		contract, exists := contractByName[formName]
		if !exists {
			return nil, fmt.Errorf("generated query field %s has no spec contract", formName)
		}
		fieldValue := value.Field(index)
		for fieldValue.Kind() == reflect.Interface || fieldValue.Kind() == reflect.Pointer {
			if fieldValue.IsNil() {
				fieldValue = reflect.Value{}
				break
			}
			fieldValue = fieldValue.Elem()
		}
		if !fieldValue.IsValid() {
			continue
		}
		if err := appendExpectedQuery(expected, contract, fieldValue); err != nil {
			return nil, fmt.Errorf("query parameter %s: %w", formName, err)
		}
	}
	return expected, nil
}

func appendExpectedQuery(expected map[string][]string, contract queryContract, value reflect.Value) error {
	if value.CanInterface() {
		if marshaler, ok := value.Interface().(interface{ MarshalText() ([]byte, error) }); ok {
			encoded, err := marshaler.MarshalText()
			if err != nil {
				return err
			}
			expected[contract.Name] = append(expected[contract.Name], string(encoded))
			return nil
		}
		if stringer, ok := value.Interface().(fmt.Stringer); ok {
			expected[contract.Name] = append(expected[contract.Name], stringer.String())
			return nil
		}
	}

	if contract.SchemaType == schemaTypeObject {
		properties, err := objectProperties(value)
		if err != nil {
			return err
		}
		if contract.Style == "deepObject" {
			for name, property := range properties {
				expected[contract.Name+"["+name+"]"] = append(expected[contract.Name+"["+name+"]"], property)
			}
			return nil
		}
		if contract.Style == "form" && contract.Explode {
			for name, property := range properties {
				expected[name] = append(expected[name], property)
			}
			return nil
		}
		parts := make([]string, 0, len(properties)*2)
		for _, name := range sortedPropertyNames(properties) {
			parts = append(parts, name, properties[name])
		}
		expected[contract.Name] = append(expected[contract.Name], strings.Join(parts, ","))
		return nil
	}

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		values := make([]string, value.Len())
		for index := range value.Len() {
			values[index] = scalarQueryValue(value.Index(index))
		}
		if contract.Style == "form" && contract.Explode {
			expected[contract.Name] = append(expected[contract.Name], values...)
			return nil
		}
		delimiter := ","
		switch contract.Style {
		case "spaceDelimited":
			delimiter = " "
		case "pipeDelimited":
			delimiter = "|"
		}
		expected[contract.Name] = append(expected[contract.Name], strings.Join(values, delimiter))
		return nil
	}
	expected[contract.Name] = append(expected[contract.Name], scalarQueryValue(value))
	return nil
}

func objectProperties(value reflect.Value) (map[string]string, error) {
	properties := make(map[string]string)
	//nolint:exhaustive // Only structs and maps have object query semantics.
	switch value.Kind() {
	case reflect.Struct:
		for index := range value.NumField() {
			fieldType := value.Type().Field(index)
			name := strings.Split(fieldType.Tag.Get("json"), ",")[0]
			if name == "" || name == "-" {
				name = fieldType.Name
			}
			fieldValue := value.Field(index)
			for fieldValue.Kind() == reflect.Interface || fieldValue.Kind() == reflect.Pointer {
				if fieldValue.IsNil() {
					fieldValue = reflect.Value{}
					break
				}
				fieldValue = fieldValue.Elem()
			}
			if fieldValue.IsValid() {
				properties[name] = scalarQueryValue(fieldValue)
			}
		}
	case reflect.Map:
		keys := value.MapKeys()
		for _, key := range keys {
			properties[scalarQueryValue(key)] = scalarQueryValue(value.MapIndex(key))
		}
	default:
		return nil, fmt.Errorf("object query value has type %s", value.Type())
	}
	return properties, nil
}

func sortedPropertyNames(properties map[string]string) []string {
	names := make([]string, 0, len(properties))
	for name := range properties {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func scalarQueryValue(value reflect.Value) string {
	if value.CanInterface() {
		if marshaler, ok := value.Interface().(interface{ MarshalText() ([]byte, error) }); ok {
			encoded, err := marshaler.MarshalText()
			if err == nil {
				return string(encoded)
			}
		}
		if stringer, ok := value.Interface().(fmt.Stringer); ok {
			return stringer.String()
		}
	}
	return fmt.Sprint(value.Interface())
}

func queryKeyMatches(contract queryContract, actual string) bool {
	if contract.SchemaType == schemaTypeObject && contract.Style == "deepObject" {
		return strings.HasPrefix(actual, contract.Name+"[") && strings.HasSuffix(actual, "]")
	}
	if contract.SchemaType == schemaTypeObject && contract.Style == "form" && contract.Explode {
		for _, property := range contract.Properties {
			if property == actual {
				return true
			}
		}
	}
	return actual == contract.Name
}

func assertBodyContract(t *testing.T, adapter operationAdapter, request *http.Request, captured capturedRequest) {
	t.Helper()
	body := captured.body
	contentType := request.Header.Get("Content-Type")
	if adapter.BodyRequired {
		require.NotEmpty(t, body, "required body missing for %s", adapter.ID)
	}
	if len(body) == 0 {
		require.Empty(t, contentType, "unexpected Content-Type without body for %s", adapter.ID)
		return
	}
	require.NotEmpty(t, contentType, "body has no Content-Type for %s", adapter.ID)
	actualMediaType, _, err := mime.ParseMediaType(contentType)
	require.NoError(t, err)
	allowed := false
	for _, mediaType := range adapter.BodyMediaTypes {
		if mediaTypeMatches(mediaType, actualMediaType) {
			allowed = true
			break
		}
	}
	require.True(t, allowed, "unexpected Content-Type %q for %s; declared %v", contentType, adapter.ID, adapter.BodyMediaTypes)
	if strings.Contains(actualMediaType, "json") {
		require.True(t, json.Valid(body), "invalid JSON request body for %s: %s", adapter.ID, body)
		require.True(t, captured.hasExpectedBody, "no typed body expectation generated for %s", adapter.ID)
		var expected any
		var actual any
		require.NoError(t, json.Unmarshal(captured.expectedBody, &expected))
		require.NoError(t, json.Unmarshal(body, &actual))
		require.Equal(t, expected, actual, "serialized JSON body differs from typed request for %s", adapter.ID)
	} else {
		require.True(t, captured.hasExpectedBody, "no raw body expectation generated for %s", adapter.ID)
		require.Equal(t, captured.expectedBody, body, "raw request body differs for %s", adapter.ID)
	}
}

func mediaTypeMatches(declared, actual string) bool {
	declaredType, _, err := mime.ParseMediaType(declared)
	if err != nil {
		return false
	}
	if declaredType == actual {
		return true
	}
	if strings.HasSuffix(declaredType, "/*") {
		return strings.HasPrefix(actual, strings.TrimSuffix(declaredType, "*"))
	}
	return false
}

func requestPathFromTemplate(template string, values []pathValue) (string, error) {
	byName := make(map[string]pathValue, len(values))
	for _, value := range values {
		byName[value.Contract.Name] = value
	}
	var expected strings.Builder
	expected.WriteString("/request-sweep")
	for len(template) > 0 {
		start := strings.IndexByte(template, '{')
		if start < 0 {
			expected.WriteString(template)
			break
		}
		expected.WriteString(template[:start])
		end := strings.IndexByte(template[start:], '}')
		if end < 0 {
			return "", fmt.Errorf("unclosed path placeholder in %q", template)
		}
		name := template[start+1 : start+end]
		value, exists := byName[name]
		if !exists {
			return "", fmt.Errorf("no sample value for path parameter %q", name)
		}
		serialized, err := serializePathParameter(value.Contract, value.Value)
		if err != nil {
			return "", fmt.Errorf("path parameter %s: %w", name, err)
		}
		expected.WriteString(serialized)
		template = template[start+end+1:]
	}
	return expected.String(), nil
}

func serializePathParameter(contract pathContract, input any) (string, error) {
	value := reflect.ValueOf(input)
	for value.IsValid() && (value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer) {
		if value.IsNil() {
			return "", errors.New("nil path parameter")
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return "", errors.New("invalid path parameter")
	}
	if contract.SchemaType == "array" && (value.Kind() == reflect.Array || value.Kind() == reflect.Slice) {
		parts := make([]string, value.Len())
		for index := range value.Len() {
			parts[index] = url.PathEscape(pathScalarValue(value.Index(index)))
		}
		joined := strings.Join(parts, ",")
		switch contract.Style {
		case parameterStyleLabel:
			return "." + joined, nil
		case parameterStyleMatrix:
			return ";" + contract.Name + "=" + joined, nil
		default:
			return joined, nil
		}
	}
	if contract.SchemaType == schemaTypeObject {
		properties, err := objectProperties(value)
		if err != nil {
			return "", err
		}
		names := sortedPropertyNames(properties)
		parts := make([]string, 0, len(names))
		for _, name := range names {
			parts = append(parts, url.PathEscape(name)+"="+url.PathEscape(properties[name]))
		}
		separator := ","
		prefix := ""
		switch contract.Style {
		case parameterStyleLabel:
			prefix = "."
			if contract.Explode {
				separator = "."
			}
		case parameterStyleMatrix:
			if contract.Explode {
				for index, name := range names {
					parts[index] = ";" + url.PathEscape(name) + "=" + url.PathEscape(properties[name])
				}
				return strings.Join(parts, ""), nil
			}
			prefix = ";" + contract.Name + "="
		}
		if !contract.Explode {
			parts = parts[:0]
			for _, name := range names {
				parts = append(parts, url.PathEscape(name), url.PathEscape(properties[name]))
			}
		}
		return prefix + strings.Join(parts, separator), nil
	}

	serialized := url.PathEscape(pathScalarValue(value))
	switch contract.Style {
	case parameterStyleLabel:
		return "." + serialized, nil
	case parameterStyleMatrix:
		return ";" + contract.Name + "=" + serialized, nil
	default:
		return serialized, nil
	}
}

func pathScalarValue(value reflect.Value) string {
	if value.CanInterface() {
		if marshaler, ok := value.Interface().(interface{ MarshalText() ([]byte, error) }); ok {
			encoded, err := marshaler.MarshalText()
			if err == nil {
				return string(encoded)
			}
		}
		if stringer, ok := value.Interface().(fmt.Stringer); ok {
			return stringer.String()
		}
		if marshaler, ok := value.Interface().(json.Marshaler); ok {
			encoded, err := marshaler.MarshalJSON()
			if err == nil {
				var decoded any
				if json.Unmarshal(encoded, &decoded) == nil {
					if text, ok := decoded.(string); ok {
						return text
					}
					return fmt.Sprint(decoded)
				}
			}
		}
	}
	return fmt.Sprint(value.Interface())
}

func requestSampleFromJSON[T any](sampleJSON string) (T, error) {
	var sample T
	err := json.Unmarshal([]byte(sampleJSON), &sample)
	return sample, err
}

func requestSampleParams[T any](contracts []queryContract) (T, error) {
	paramsJSON := make(map[string]json.RawMessage, len(contracts))
	for _, contract := range contracts {
		if contract.SampleJSON == "" {
			var zero T
			return zero, fmt.Errorf("query parameter %s has no schema-derived sample", contract.Name)
		}
		paramsJSON[contract.Name] = json.RawMessage(contract.SampleJSON)
	}
	encoded, err := json.Marshal(paramsJSON)
	if err != nil {
		var zero T
		return zero, err
	}
	return requestSampleFromJSON[T](string(encoded))
}

func requestSampleValueFromJSON(sampleJSON string) any {
	var sample any
	if err := json.Unmarshal([]byte(sampleJSON), &sample); err != nil {
		panic(err)
	}
	return sample
}

func requestMediaType(mediaTypes []string) string {
	if len(mediaTypes) == 0 {
		return "application/json"
	}
	return mediaTypes[0]
}

func requestSweepBodyReader(observed *capturedRequest, sampleJSON string) io.Reader {
	observed.expectedBody = []byte(sampleJSON)
	observed.hasExpectedBody = true
	return io.NopCloser(bytes.NewReader(observed.expectedBody))
}

func specDigest(contents []byte) string {
	digest := sha256.Sum256(contents)
	return hex.EncodeToString(digest[:])
}
