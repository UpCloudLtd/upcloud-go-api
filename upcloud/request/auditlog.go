package request

import "net/url"

const (
	// ExportAuditLogFormatCSV is the format for exporting audit logs as CSV.
	ExportAuditLogFormatCSV = "csv"
	// ExportAuditLogFormatJSON is the format for exporting audit logs as JSON.
	ExportAuditLogFormatJSON = "json"
)

// ExportAuditLogRequest represents a request to export audit logs.
type ExportAuditLogRequest struct {
	Format string
}

func (r *ExportAuditLogRequest) RequestURL() string {
	u := "/audit-logs/export"
	if r.Format != "" {
		u += "?format=" + url.QueryEscape(r.Format)
	}
	return u
}
