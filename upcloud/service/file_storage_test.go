package service

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorageService_AllMethods(t *testing.T) {
	t.Parallel()
	record(t, "filestorage", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := &request.CreateFileStorageRequest{
			Name:             "go-sdk-test-filesto-allmethods",
			Zone:             "fi-hel1",
			ConfiguredStatus: "started",
			SizeGiB:          250,
			Labels:           []upcloud.Label{{Key: "managedBy", Value: "upcloud-go-sdk-unit-test"}},
		}
		created, err := svc.CreateFileStorage(ctx, createReq)
		assert.NoError(t, err)
		require.NotNil(t, created)

		uuid := created.UUID
		defer func() {
			_ = svc.DeleteFileStorage(ctx, &request.DeleteFileStorageRequest{UUID: uuid})
		}()

		fileStorages, err := svc.GetFileStorages(ctx, &request.GetFileStoragesRequest{})
		assert.NoError(t, err)
		require.NotNil(t, fileStorages)

		fileStorage, err := svc.GetFileStorage(ctx, &request.GetFileStorageRequest{UUID: uuid})
		assert.NoError(t, err)
		require.NotNil(t, fileStorage)

		replaceReq := &request.ReplaceFileStorageRequest{
			UUID:             uuid,
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

		newName := "go-sdk-test-filesto-modified"
		newLabels := []upcloud.Label{{Key: "managedBy", Value: "modified"}}
		modifyReq := &request.ModifyFileStorageRequest{
			UUID:   uuid,
			Name:   &newName,
			Labels: &newLabels,
		}
		modified, err := svc.ModifyFileStorage(ctx, modifyReq)
		assert.NoError(t, err)
		require.NotNil(t, modified)
		assert.Equal(t, "go-sdk-test-filesto-modified", modified.Name)
		assert.Equal(t, "modified", modified.Labels[0].Value)

		// Network creation
		networkCreateReq := &request.CreateFileStorageNetworkRequest{
			ServiceUUID: uuid,
			Name:        "test-network",
			Family:      "IPv4",
		}
		network, err := svc.CreateFileStorageNetwork(ctx, networkCreateReq)
		assert.NoError(t, err)
		require.NotNil(t, network)

		networks, err := svc.GetFileStorageNetworks(ctx, &request.GetFileStorageNetworksRequest{ServiceUUID: uuid})
		assert.NoError(t, err)
		assert.NotNil(t, networks)

		gotNetwork, err := svc.GetFileStorageNetwork(ctx, &request.GetFileStorageNetworkRequest{ServiceUUID: uuid, NetworkName: network.Name})
		assert.NoError(t, err)
		require.NotNil(t, gotNetwork)

		modNetworkReq := &request.ModifyFileStorageNetworkRequest{
			ServiceUUID: uuid,
			NetworkName: network.Name,
			IPAddress:   nil,
			Family:      nil,
		}
		modNetwork, err := svc.ModifyFileStorageNetwork(ctx, modNetworkReq)
		assert.NoError(t, err)
		require.NotNil(t, modNetwork)

		err = svc.DeleteFileStorageNetwork(ctx, &request.DeleteFileStorageNetworkRequest{ServiceUUID: uuid, NetworkName: network.Name})
		assert.NoError(t, err)

		// Share creation
		shareCreate := &request.CreateFileStorageShareRequest{
			ServiceUUID: uuid,
			Name:        "test-share",
			Path:        "/mnt/test",
			ACL:         []upcloud.FileStorageACL{{Target: "user", Permission: "rw"}},
		}
		share, err := svc.CreateFileStorageShare(ctx, shareCreate)
		assert.NoError(t, err)
		require.NotNil(t, share)

		sharereq := &request.GetFileStorageSharesRequest{ServiceUUID: uuid}
		shares, err := svc.GetFileStorageShares(ctx, sharereq)
		assert.NoError(t, err)
		require.NotNil(t, shares)

		gotShare, err := svc.GetFileStorageShare(ctx, &request.GetFileStorageShareRequest{ServiceUUID: uuid, ShareName: share.Name})
		assert.NoError(t, err)
		assert.NotNil(t, gotShare)

		modSharePath := "/mnt/modified"
		modShareACL := []upcloud.FileStorageACL{{Target: "user", Permission: "ro"}}
		modShareReq := &request.ModifyFileStorageShareRequest{
			ServiceUUID: uuid,
			ShareName:   share.Name,
			ModifyFileStorageShare: request.ModifyFileStorageShare{
				Path: &modSharePath,
				ACL:  &modShareACL,
			},
		}
		modShare, err := svc.ModifyFileStorageShare(ctx, modShareReq)
		assert.NoError(t, err)
		require.NotNil(t, modShare)
		assert.Equal(t, "/mnt/modified", modShare.Path)
		assert.Equal(t, "ro", modShare.ACL[0].Permission)

		err = svc.DeleteFileStorageShare(ctx, &request.DeleteFileStorageShareRequest{ServiceUUID: uuid, ShareName: share.Name})
		assert.NoError(t, err)

		// Label creation
		labelKey := "test-label"
		labelValue := "test"
		labelCreateReq := &request.CreateFileStorageLabelRequest{
			ServiceUUID: uuid,
			Label: upcloud.Label{
				Key:   labelKey,
				Value: labelValue,
			},
		}
		label, err := svc.CreateFileStorageLabel(ctx, labelCreateReq)
		assert.NoError(t, err)
		require.NotNil(t, label)

		labels, err := svc.GetFileStorageLabels(ctx, &request.GetFileStorageLabelsRequest{ServiceUUID: uuid})
		assert.NoError(t, err)
		require.NotNil(t, labels)

		gotLabel, err := svc.GetFileStorageLabel(ctx, &request.GetFileStorageLabelRequest{ServiceUUID: uuid, LabelKey: label.Key})
		assert.NoError(t, err)
		require.NotNil(t, gotLabel)

		modLabelValue := "modified"
		modLabelReq := &request.ModifyFileStorageLabelRequest{
			ServiceUUID: uuid,
			Label: upcloud.Label{
				Key:   label.Key,
				Value: modLabelValue,
			},
		}
		modLabel, err := svc.ModifyFileStorageLabel(ctx, modLabelReq)
		assert.NoError(t, err)
		require.NotNil(t, modLabel)
		assert.Equal(t, "modified", modLabel.Value)

		err = svc.DeleteFileStorageLabel(ctx, &request.DeleteFileStorageLabelRequest{ServiceUUID: uuid, LabelKey: label.Key})
		assert.NoError(t, err)
	})
}
