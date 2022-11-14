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

// TestCreateTagContext tests the creation of a single tag
func TestCreateTagContext(t *testing.T) {
	t.Parallel()
	recordWithContext(t, "createtag", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svcContext *ServiceContext) {
		// ignore errors, but delete the tag if it happens to exist
		_ = svcContext.DeleteTag(ctx, &request.DeleteTagRequest{
			Name: "testTag",
		})

		tag, err := svcContext.CreateTag(ctx, &request.CreateTagRequest{
			Tag: upcloud.Tag{
				Name: "testTag",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, "testTag", tag.Name)
	})
}

// TestGetTagsContext tests that GetTags returns multiple tags and it, at least, contains the 3
// we create.
func TestGetTagsContext(t *testing.T) {
	t.Parallel()
	recordWithContext(t, "gettags", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svcContext *ServiceContext) {
		testData := []string{
			"testgettags_tag1",
			"testgettags_tag2",
			"testgettags_tag3",
		}

		for _, tag := range testData {
			// Delete all the tags we're about to create.
			// We don't care about errors.
			_ = svcContext.DeleteTag(ctx, &request.DeleteTagRequest{
				Name: tag,
			})
		}

		for _, tag := range testData {
			_, err := svcContext.CreateTag(ctx, &request.CreateTagRequest{
				Tag: upcloud.Tag{
					Name:        tag,
					Description: tag + " description",
				},
			})

			require.NoError(t, err)
		}

		tags, err := svcContext.GetTags(ctx)
		require.NoError(t, err)
		// There may be other tags so the length must be
		// greater than or equal to.
		assert.GreaterOrEqual(t, len(tags.Tags), len(testData))
		for _, expectedTag := range testData {
			var found bool
			for _, tag := range tags.Tags {
				if tag.Name == expectedTag {
					found = true
					assert.Equal(t, expectedTag+" description", tag.Description)
					break
				}
			}
			assert.True(t, found)
		}

		for _, tag := range tags.Tags {
			err := svcContext.DeleteTag(ctx, &request.DeleteTagRequest{
				Name: tag.Name,
			})
			require.NoError(t, err)
		}
	})
}

// TestTaggingContext tests that all tagging-related functionality works correctly. It performs the following actions:
//   - creates a server
//   - creates three tags
//   - assigns the first tag to the server
//   - renames the second tag
//   - deletes the third tag
//   - untags the first tag from the server
func TestTaggingContext(t *testing.T) {
	t.Parallel()
	recordWithContext(t, "tagging", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svcContext *ServiceContext) {
		// Create the server
		serverDetails, err := createServerContext(ctx, rec, svcContext, "TestTagging")
		require.NoError(t, err)
		t.Logf("Server %s with UUID %s created", serverDetails.Title, serverDetails.UUID)

		// Remove all existing tags
		t.Log("Deleting any existing tags ...")
		err = deleteAllTags(ctx, svcContext)
		require.NoError(t, err)

		// Create three tags
		tags := []string{
			"tag1",
			"tag2",
			"tag3",
		}

		for _, tag := range tags {
			t.Logf("Creating tag %s", tag)
			tagDetails, err := svcContext.CreateTag(ctx, &request.CreateTagRequest{
				Tag: upcloud.Tag{
					Name: tag,
				},
			})

			require.NoError(t, err)
			assert.Equal(t, tag, tagDetails.Name)
			t.Logf("Tag %s created", tagDetails.Name)
		}

		// Assign the first tag to the server
		serverDetails, err = svcContext.TagServer(ctx, &request.TagServerRequest{
			UUID: serverDetails.UUID,
			Tags: []string{
				"tag1",
			},
		})
		require.NoError(t, err)
		assert.Contains(t, serverDetails.Tags, "tag1")
		var utilityCount int
		for _, ip := range serverDetails.IPAddresses {
			assert.NotEqual(t, upcloud.IPAddressAccessPrivate, ip.Access)
			if ip.Access == upcloud.IPAddressAccessUtility {
				utilityCount++
			}
		}
		assert.NotZero(t, utilityCount)
		t.Logf("Server %s is now tagged with tag %s", serverDetails.Title, "tag1")

		// Rename the second tag
		tagDetails, err := svcContext.ModifyTag(ctx, &request.ModifyTagRequest{
			Name: "tag2",
			Tag: upcloud.Tag{
				Name: "tag2_renamed",
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "tag2_renamed", tagDetails.Name)
		t.Logf("Tag tag2 renamed to %s", tagDetails.Name)

		// Delete the third tag
		err = svcContext.DeleteTag(ctx, &request.DeleteTagRequest{
			Name: "tag3",
		})

		require.NoError(t, err)
		t.Log("Tag tag3 deleted")

		// Untag the server
		t.Logf("Removing tag %s from server %s", "tag1", serverDetails.UUID)
		serverDetails, err = svcContext.UntagServer(ctx, &request.UntagServerRequest{
			UUID: serverDetails.UUID,
			Tags: []string{
				"tag1",
			},
		})
		require.NoError(t, err)
		assert.NotContains(t, serverDetails.Tags, "tag1")
		utilityCount = 0
		for _, ip := range serverDetails.IPAddresses {
			assert.NotEqual(t, upcloud.IPAddressAccessPrivate, ip.Access)
			if ip.Access == upcloud.IPAddressAccessUtility {
				utilityCount++
			}
		}
		assert.NotZero(t, utilityCount)
		t.Logf("Server %s is now untagged", serverDetails.Title)
	})
}

// deleteAllTags deletes all existing tags.
func deleteAllTags(ctx context.Context, svcContext *ServiceContext) error {
	tags, err := svcContext.GetTags(ctx)
	if err != nil {
		return err
	}

	for _, tagDetails := range tags.Tags {
		err = svcContext.DeleteTag(ctx, &request.DeleteTagRequest{
			Name: tagDetails.Name,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
