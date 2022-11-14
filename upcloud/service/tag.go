package service

import (
	"encoding/json"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

type Tag interface {
	GetTags() (*upcloud.Tags, error)
	CreateTag(r *request.CreateTagRequest) (*upcloud.Tag, error)
	ModifyTag(r *request.ModifyTagRequest) (*upcloud.Tag, error)
	DeleteTag(r *request.DeleteTagRequest) error
	TagServer(r *request.TagServerRequest) (*upcloud.ServerDetails, error)
	UntagServer(r *request.UntagServerRequest) (*upcloud.ServerDetails, error)
}

var _ Tag = (*Service)(nil)

// CreateTag creates a new tag, optionally assigning it to one or more servers at the same time
func (s *Service) CreateTag(r *request.CreateTagRequest) (*upcloud.Tag, error) {
	tagDetails := upcloud.Tag{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &tagDetails)
	if err != nil {
		return nil, err
	}

	return &tagDetails, nil
}

// ModifyTag modifies a tag (e.g. renaming it)
func (s *Service) ModifyTag(r *request.ModifyTagRequest) (*upcloud.Tag, error) {
	tagDetails := upcloud.Tag{}
	requestBody, _ := json.Marshal(r)
	response, err := s.client.PerformJSONPutRequest(s.client.CreateRequestURL(r.RequestURL()), requestBody)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &tagDetails)
	if err != nil {
		return nil, err
	}

	return &tagDetails, nil
}

// DeleteTag deletes the specified tag
func (s *Service) DeleteTag(r *request.DeleteTagRequest) error {
	err := s.client.PerformJSONDeleteRequest(s.client.CreateRequestURL(r.RequestURL()))
	if err != nil {
		return parseJSONServiceError(err)
	}

	return nil
}

// GetTags returns all tags
func (s *Service) GetTags() (*upcloud.Tags, error) {
	tags := upcloud.Tags{}
	response, err := s.basicGetRequest("/tag")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &tags)
	if err != nil {
		return nil, err
	}

	return &tags, nil
}

// TagServer tags a server with with one or more tags
func (s *Service) TagServer(r *request.TagServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), nil)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}

// UntagServer removes one or more tags from a server
func (s *Service) UntagServer(r *request.UntagServerRequest) (*upcloud.ServerDetails, error) {
	serverDetails := upcloud.ServerDetails{}
	response, err := s.client.PerformJSONPostRequest(s.client.CreateRequestURL(r.RequestURL()), nil)
	if err != nil {
		return nil, parseJSONServiceError(err)
	}

	err = json.Unmarshal(response, &serverDetails)
	if err != nil {
		return nil, err
	}

	return &serverDetails, nil
}
