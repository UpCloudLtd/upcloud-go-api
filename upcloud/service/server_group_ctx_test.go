package service

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerGroupsContext(t *testing.T) {
	t.Parallel()

	recordWithContext(t, "servergroups", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		srv, err := createMinimalServer(svc, "TestServerGroups")
		require.NoError(t, err)
		// create new server group
		group, err := svcContext.CreateServerGroup(ctx, &request.CreateServerGroupRequest{
			Labels:  &upcloud.LabelSlice{upcloud.Label{Key: "managedBy", Value: "upcloud-go-sdk-integration-test"}},
			Members: upcloud.ServerUUIDSlice{srv.UUID},
			Title:   "test-title",
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, upcloud.LabelSlice{upcloud.Label{Key: "managedBy", Value: "upcloud-go-sdk-integration-test"}}, group.Labels)
		assert.Equal(t, "test-title", group.Title)
		assert.Len(t, group.Members, 1)

		// clear the group of its members
		group, err = svcContext.ModifyServerGroup(ctx, &request.ModifyServerGroupRequest{
			UUID:    group.UUID,
			Title:   "test-title-edit",
			Members: &upcloud.ServerUUIDSlice{},
		})
		assert.NoError(t, err)
		assert.Equal(t, "test-title-edit", group.Title)
		assert.Len(t, group.Members, 0)

		// append server to group without modifying title or labels
		group, err = svcContext.ModifyServerGroup(ctx, &request.ModifyServerGroupRequest{
			UUID:    group.UUID,
			Members: &upcloud.ServerUUIDSlice{srv.UUID},
		})
		assert.NoError(t, err)
		assert.Equal(t, "test-title-edit", group.Title)
		assert.Len(t, group.Members, 1)

		// modify only title and labels without touching members
		newLabelSlice := append(group.Labels, upcloud.Label{Key: "title", Value: "test-title"})
		group, err = svcContext.ModifyServerGroup(ctx, &request.ModifyServerGroupRequest{
			Labels: &newLabelSlice,
			Title:  "test-title",
			UUID:   group.UUID,
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, newLabelSlice, group.Labels)
		assert.Equal(t, "test-title", group.Title)
		assert.Len(t, group.Members, 1)

		// get server groups
		groups, err := svcContext.GetServerGroups(ctx, &request.GetServerGroupsRequest{})
		assert.NoError(t, err)
		assert.Len(t, groups, 1)

		// get server groups with filters
		groups, err = svcContext.GetServerGroupsWithFilters(ctx, &request.GetServerGroupsWithFiltersRequest{
			Filters: []request.ServerGroupFilter{
				request.FilterLabelKey{Key: "managedBy"},
				request.FilterLabel{Label: upcloud.Label{
					Key:   "title",
					Value: "test-title",
				}},
			},
		})
		assert.NoError(t, err)
		assert.Len(t, groups, 1)

		// get server group
		group, err = svcContext.GetServerGroup(ctx, &request.GetServerGroupRequest{UUID: group.UUID})
		assert.NoError(t, err)
		assert.Equal(t, "test-title", group.Title)
		assert.Len(t, group.Members, 1)

		// delete server group
		err = svcContext.DeleteServerGroup(ctx, &request.DeleteServerGroupRequest{UUID: group.UUID})
		assert.NoError(t, err)

		// skip server cleanup if recorder is replaying to save some time
		if rec.Mode() == recorder.ModeReplaying {
			return
		}

		// delete server
		if err := stopServer(svc, srv.UUID); err != nil {
			t.Log(err)
		} else {
			if err := deleteServer(svc, srv.UUID); err != nil {
				t.Log(err)
			}
		}
	})
}
