package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage(t *testing.T) {
	t.Parallel()
	record(t, "filestorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := &request.CreateFileStorageRequest{
			Name:             "go-sdk-test-filesto",
			Zone:             "fi-hel2",
			ConfiguredStatus: upcloud.FileStorageConfiguredStatusStarted,
			SizeGiB:          250,
			Labels:           []upcloud.Label{{Key: "managedBy", Value: "upcloud-go-sdk-unit-test"}},
		}
		created, err := svc.CreateFileStorage(ctx, createReq)
		assert.NoError(t, err)
		require.NotNil(t, created)

		err = waitForFileStorageRunningOperationalState(ctx, rec, svc, created.UUID)
		require.NoError(t, err)

		fileStorages, err := svc.GetFileStorages(ctx, &request.GetFileStoragesRequest{})
		assert.NoError(t, err)
		require.NotNil(t, fileStorages)

		fileStorage, err := svc.GetFileStorage(ctx, &request.GetFileStorageRequest{UUID: created.UUID})
		assert.NoError(t, err)
		require.NotNil(t, fileStorage)

		replaceReq := &request.ReplaceFileStorageRequest{
			UUID:             created.UUID,
			Name:             "go-sdk-test-filesto-replaced",
			ConfiguredStatus: createReq.ConfiguredStatus,
			SizeGiB:          createReq.SizeGiB,
			Labels:           []upcloud.Label{{Key: "managedBy", Value: "replaced"}},
		}
		replaced, err := svc.ReplaceFileStorage(ctx, replaceReq)
		assert.NoError(t, err)
		require.NotNil(t, replaced)
		assert.Equal(t, "go-sdk-test-filesto-replaced", replaced.Name)
		assert.Equal(t, "replaced", replaced.Labels[0].Value)

		err = waitForFileStorageRunningOperationalState(ctx, rec, svc, replaced.UUID)
		require.NoError(t, err)

		newName := "go-sdk-test-filesto-modified"
		newLabels := []upcloud.Label{{Key: "managedBy", Value: "modified"}}
		modifyReq := &request.ModifyFileStorageRequest{
			UUID:   replaced.UUID,
			Name:   &newName,
			Labels: &newLabels,
		}
		modified, err := svc.ModifyFileStorage(ctx, modifyReq)
		assert.NoError(t, err)
		require.NotNil(t, modified)
		assert.Equal(t, "go-sdk-test-filesto-modified", modified.Name)
		assert.Equal(t, "modified", modified.Labels[0].Value)

		err = waitForFileStorageRunningOperationalState(ctx, rec, svc, modified.UUID)
		require.NoError(t, err)

		network, err := svc.CreateNetwork(ctx, &request.CreateNetworkRequest{
			Name: "go-sdk-test-filesto",
			Zone: "fi-hel2",
			IPNetworks: upcloud.IPNetworkSlice{upcloud.IPNetwork{
				Address: "172.28.1.0/24",
				Family:  "IPv4",
				Gateway: "172.28.1.1",
			}},
		})
		require.NoError(t, err)

		defer func() {
			_ = svc.DeleteFileStorage(ctx, &request.DeleteFileStorageRequest{UUID: replaced.UUID})
			_ = svc.WaitForFileStorageDeletion(ctx, &request.WaitForFileStorageDeletionRequest{UUID: replaced.UUID})
			_ = svc.DeleteNetwork(ctx, &request.DeleteNetworkRequest{UUID: network.UUID})
		}()

		networkCreateReq := &request.CreateFileStorageNetworkRequest{
			ServiceUUID: modified.UUID,
			UUID:        network.UUID,
			Name:        "go-sdk-test-filesto",
			Family:      "IPv4",
		}
		fileStorageNetwork, err := svc.CreateFileStorageNetwork(ctx, networkCreateReq)
		assert.NoError(t, err)
		require.NotNil(t, fileStorageNetwork)

		networks, err := svc.GetFileStorageNetworks(ctx, &request.GetFileStorageNetworksRequest{ServiceUUID: modified.UUID})
		assert.NoError(t, err)
		assert.NotNil(t, networks)

		gotNetwork, err := svc.GetFileStorageNetwork(ctx, &request.GetFileStorageNetworkRequest{ServiceUUID: modified.UUID, NetworkName: fileStorageNetwork.Name})
		assert.NoError(t, err)
		require.NotNil(t, gotNetwork)

		modNetworkReq := &request.ModifyFileStorageNetworkRequest{
			ServiceUUID: modified.UUID,
			NetworkName: fileStorageNetwork.Name,
			IPAddress:   nil,
			Family:      nil,
		}
		modNetwork, err := svc.ModifyFileStorageNetwork(ctx, modNetworkReq)
		assert.NoError(t, err)
		require.NotNil(t, modNetwork)

		err = svc.DeleteFileStorageNetwork(ctx, &request.DeleteFileStorageNetworkRequest{ServiceUUID: modified.UUID, NetworkName: fileStorageNetwork.Name})
		assert.NoError(t, err)

		shareCreate := &request.CreateFileStorageShareRequest{
			ServiceUUID: modified.UUID,
			Name:        "test-share",
			Path:        "/data/test",
			ACL:         []upcloud.FileStorageShareACL{{Name: "test-share-acl", Target: "user", Permission: upcloud.FileStorageShareACLPermissionReadWrite}},
		}
		share, err := svc.CreateFileStorageShare(ctx, shareCreate)
		assert.NoError(t, err)
		require.NotNil(t, share)

		currentState, err := svc.GetFileStorageCurrentState(ctx, &request.GetFileStorageCurrentStateRequest{UUID: modified.UUID})
		assert.NoError(t, err)
		require.NotNil(t, currentState)

		sharereq := &request.GetFileStorageSharesRequest{ServiceUUID: modified.UUID}
		shares, err := svc.GetFileStorageShares(ctx, sharereq)
		assert.NoError(t, err)
		require.NotNil(t, shares)

		gotShare, err := svc.GetFileStorageShare(ctx, &request.GetFileStorageShareRequest{ServiceUUID: modified.UUID, ShareName: share.Name})
		assert.NoError(t, err)
		assert.NotNil(t, gotShare)
		perm := upcloud.FileStorageShareACLPermissionReadOnly

		modShareACL := []request.FileStorageShareACL{{Target: upcloud.StringPtr("user"), Permission: &perm}}
		modShareReq := &request.ModifyFileStorageShareRequest{
			ServiceUUID: modified.UUID,
			ShareName:   share.Name,
			ModifyFileStorageShare: request.ModifyFileStorageShare{
				ACL: &modShareACL,
			},
		}
		modShare, err := svc.ModifyFileStorageShare(ctx, modShareReq)
		assert.NoError(t, err)
		require.NotNil(t, modShare)
		assert.Equal(t, "/data/test", modShare.Path)
		assert.Equal(t, upcloud.FileStorageShareACLPermissionReadOnly, modShare.ACL[0].Permission)

		modShareEmptyACLReq := &request.ModifyFileStorageShareRequest{
			ServiceUUID: modified.UUID,
			ShareName:   share.Name,
			ModifyFileStorageShare: request.ModifyFileStorageShare{
				ACL: nil,
			},
		}
		modShareEmptyACL, err := svc.ModifyFileStorageShare(ctx, modShareEmptyACLReq)
		assert.NoError(t, err)
		require.NotNil(t, modShare)
		assert.Equal(t, nil, modShareEmptyACL.ACL)

		err = svc.DeleteFileStorageShare(ctx, &request.DeleteFileStorageShareRequest{ServiceUUID: modified.UUID, ShareName: share.Name})
		assert.NoError(t, err)

		newShareACL := upcloud.FileStorageShareACL{Name: "acl2", Target: "user2", Permission: upcloud.FileStorageShareACLPermissionReadWrite}
		aclCreateReq := &request.CreateFileStorageShareACLRequest{
			ServiceUUID:         modified.UUID,
			ShareName:           share.Name,
			FileStorageShareACL: newShareACL,
		}

		acl, err := svc.CreateFileStorageShareACL(ctx, aclCreateReq)
		assert.NoError(t, err)
		require.NotNil(t, acl)
		assert.Equal(t, "acl2", acl.Name)
		assert.Equal(t, "user2", acl.Target)
		assert.Equal(t, upcloud.FileStorageShareACLPermissionReadWrite, acl.Permission)

		deleteACLReq := &request.DeleteFileStorageShareACLRequest{
			ServiceUUID: modified.UUID,
			ShareName:   share.Name,
			ACLName:     acl.Name,
		}
		err = svc.DeleteFileStorageShareACL(ctx, deleteACLReq)
		assert.NoError(t, err)

		labelKey := "test-label"
		labelValue := "test"
		labelCreateReq := &request.CreateFileStorageLabelRequest{
			ServiceUUID: modified.UUID,
			Label: upcloud.Label{
				Key:   labelKey,
				Value: labelValue,
			},
		}
		label, err := svc.CreateFileStorageLabel(ctx, labelCreateReq)
		assert.NoError(t, err)
		require.NotNil(t, label)

		labels, err := svc.GetFileStorageLabels(ctx, &request.GetFileStorageLabelsRequest{ServiceUUID: modified.UUID})
		assert.NoError(t, err)
		require.NotNil(t, labels)

		gotLabel, err := svc.GetFileStorageLabel(ctx, &request.GetFileStorageLabelRequest{ServiceUUID: modified.UUID, LabelKey: label.Key})
		assert.NoError(t, err)
		require.NotNil(t, gotLabel)

		modLabelValue := "modified"
		modLabelReq := &request.ModifyFileStorageLabelRequest{
			ServiceUUID: modified.UUID,
			LabelKey:    label.Key,
			Label: upcloud.Label{
				Key:   label.Key,
				Value: modLabelValue,
			},
		}
		modLabel, err := svc.ModifyFileStorageLabel(ctx, modLabelReq)
		assert.NoError(t, err)
		require.NotNil(t, modLabel)
		assert.Equal(t, "modified", modLabel.Value)

		err = svc.DeleteFileStorageLabel(ctx, &request.DeleteFileStorageLabelRequest{ServiceUUID: modified.UUID, LabelKey: label.Key})
		assert.NoError(t, err)
	})
}

func waitForFileStorageRunningOperationalState(ctx context.Context, rec *recorder.Recorder, svc *Service, fsUUID string) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}

	rec.AddPassthrough(func(h *http.Request) bool {
		return true
	})
	defer func() {
		rec.Passthroughs = nil
	}()

	_, err := svc.WaitForFileStorageOperationalState(ctx, &request.WaitForFileStorageOperationalStateRequest{
		UUID:         fsUUID,
		DesiredState: upcloud.FileStorageOperationalStateRunning,
	})

	return err
}
