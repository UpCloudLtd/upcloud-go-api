package service

import (
	"context"
	"io"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

type AuditLog interface {
	ExportAuditLog(context.Context, *request.ExportAuditLogRequest) (io.ReadCloser, error)
}

func (s *Service) ExportAuditLog(ctx context.Context, r *request.ExportAuditLogRequest) (io.ReadCloser, error) {
	return s.client.GetStream(ctx, r.RequestURL())
}
