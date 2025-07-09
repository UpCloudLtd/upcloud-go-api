package request_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

// TestExportAuditLogsRequest tests that ExportAuditLogRequest behaves correctly.
func TestExportAuditLogsRequest(t *testing.T) {
	tests := []struct {
		format      string
		expectedURL string
	}{
		{
			format:      "",
			expectedURL: "/audit-logs/export",
		},
		{
			format:      request.ExportAuditLogFormatCSV,
			expectedURL: "/audit-logs/export?format=csv",
		},
		{
			format:      request.ExportAuditLogFormatJSON,
			expectedURL: "/audit-logs/export?format=json",
		},
		{
			format:      "foo:bar",
			expectedURL: "/audit-logs/export?format=foo%3Abar",
		},
	}

	for _, test := range tests {
		r := request.ExportAuditLogRequest{
			Format: test.format,
		}
		assert.Equal(t, test.expectedURL, r.RequestURL())
	}
}
