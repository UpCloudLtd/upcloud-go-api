package service

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerGroupsContext(t *testing.T) {
	t.Parallel()

	recordWithContext(t, "servergroups", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svcContext *ServiceContext) {
		srv, err := createMinimalServerContext(ctx, rec, svcContext, "TestServerGroups")
		require.NoError(t, err)
		// create new server group
		group, err := svcContext.CreateServerGroup(ctx, &request.CreateServerGroupRequest{
			Labels:       &upcloud.LabelSlice{upcloud.Label{Key: "managedBy", Value: "upcloud-go-sdk-integration-test"}},
			Members:      upcloud.ServerUUIDSlice{srv.UUID},
			Title:        "test-title",
			AntiAffinity: upcloud.True,
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, upcloud.LabelSlice{upcloud.Label{Key: "managedBy", Value: "upcloud-go-sdk-integration-test"}}, group.Labels)
		assert.Equal(t, "test-title", group.Title)
		assert.Equal(t, upcloud.True, group.AntiAffinity)
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
		assert.Equal(t, upcloud.True, group.AntiAffinity)

		// modify anti affinity setting
		group, err = svcContext.ModifyServerGroup(ctx, &request.ModifyServerGroupRequest{
			AntiAffinity: upcloud.False,
			UUID:         group.UUID,
		})
		assert.NoError(t, err)
		assert.ElementsMatch(t, newLabelSlice, group.Labels) // Labels should not change since the previous modification
		assert.Equal(t, "test-title", group.Title)
		assert.Len(t, group.Members, 1)
		assert.Equal(t, upcloud.False, group.AntiAffinity)

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
		if err := stopServerContext(ctx, rec, svcContext, srv.UUID); err != nil {
			t.Log(err)
		} else {
			if err := deleteServer(ctx, svcContext, srv.UUID); err != nil {
				t.Log(err)
			}
		}
	})
}

// Deletes the specified server group.
func deleteServerGroup(ctx context.Context, svcContext *ServiceContext, uuid string) error {
	err := svcContext.DeleteServerGroup(ctx, &request.DeleteServerGroupRequest{
		UUID: uuid,
	})

	return err
}
