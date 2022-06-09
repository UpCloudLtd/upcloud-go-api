package service

import (
	"context"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

type TagContext interface {
	GetTags(ctx context.Context) (*upcloud.Tags, error)
	CreateTag(ctx context.Context, r *request.CreateTagRequest) (*upcloud.Tag, error)
	ModifyTag(ctx context.Context, r *request.ModifyTagRequest) (*upcloud.Tag, error)
	DeleteTag(ctx context.Context, r *request.DeleteTagRequest) error
	TagServer(ctx context.Context, r *request.TagServerRequest) (*upcloud.ServerDetails, error)
	UntagServer(ctx context.Context, r *request.UntagServerRequest) (*upcloud.ServerDetails, error)
}

// CreateTag creates a new tag, optionally assigning it to one or more servers at the same time
func (s *ServiceContext) CreateTag(ctx context.Context, r *request.CreateTagRequest) (*upcloud.Tag, error) {
	tagDetails := upcloud.Tag{}
	return &tagDetails, s.create(ctx, r, &tagDetails)
}

// ModifyTag modifies a tag (e.g. renaming it)
func (s *ServiceContext) ModifyTag(ctx context.Context, r *request.ModifyTagRequest) (*upcloud.Tag, error) {
	tagDetails := upcloud.Tag{}
	return &tagDetails, s.replace(ctx, r, &tagDetails)
}

// DeleteTag deletes the specified tag
func (s *ServiceContext) DeleteTag(ctx context.Context, r *request.DeleteTagRequest) error {
	return s.delete(ctx, r)
}

// GetTags returns all tags
func (s *ServiceContext) GetTags(ctx context.Context) (*upcloud.Tags, error) {
	tags := upcloud.Tags{}
	return &tags, s.get(ctx, "/tag", &tags)
}

// TagServer tags a server with with one or more tags
func (s *ServiceContext) TagServer(ctx context.Context, r *request.TagServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}

// UntagServer removes one or more tags from a server
func (s *ServiceContext) UntagServer(ctx context.Context, r *request.UntagServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	return &serverDetails, s.create(ctx, r, &serverDetails)
}
