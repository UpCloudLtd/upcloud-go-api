package request

import (
	"encoding/json"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestFileStorageRequest_MarshalUnmarshalJSON(t *testing.T) {
	req := CreateFileStorageRequest{
		Name:             "test-storage",
		Zone:             "fi-hel1",
		ConfiguredStatus: "started",
		SizeGiB:          100,
		Networks:         []upcloud.FileStorageNetwork{{UUID: "net-uuid", Name: "net1", Family: "IPv4", IPAddress: "192.168.1.1"}},
		Labels:           []upcloud.Label{{Key: "env", Value: "dev"}},
	}
	data, err := json.Marshal(req)
	assert.NoError(t, err)

	var out CreateFileStorageRequest
	err = json.Unmarshal(data, &out)
	assert.NoError(t, err)
	assert.Equal(t, req, out)
}

func TestFileStorageRequestURLs(t *testing.T) {
	getReq := GetFileStorageRequest{UUID: "fs-uuid"}
	assert.Equal(t, "/file-storage/fs-uuid", getReq.RequestURL())

	listReq := GetFileStoragesRequest{}
	assert.Equal(t, "/file-storage", listReq.RequestURL())

	page := &Page{Size: 5, Number: 1}
	sort := "-created_at"
	paged := GetFileStoragesRequest{Page: page, Sort: &sort}
	assert.Equal(t, "/file-storage?limit=5&offset=0&sort=-created_at", paged.RequestURL())
}

func TestCreateFileStorageRequest_MarshalJSON(t *testing.T) {
	t.Run("Minimal", func(t *testing.T) {
		req := CreateFileStorageRequest{
			Name:             "minimal",
			Zone:             "zone-1",
			ConfiguredStatus: "started",
			SizeGiB:          10,
		}
		data, err := json.Marshal(&req)
		assert.NoError(t, err)
		const expected = `{"name":"minimal","zone":"zone-1","configured_status":"started","size_gib":10}`
		assert.JSONEq(t, expected, string(data))
	})

	t.Run("Full", func(t *testing.T) {
		req := CreateFileStorageRequest{
			Name:             "full",
			Zone:             "zone-2",
			ConfiguredStatus: "started",
			SizeGiB:          20,
			Networks:         []upcloud.FileStorageNetwork{{UUID: "net-uuid", Name: "net1", Family: "IPv4", IPAddress: "192.168.1.1"}},
			Shares:           []FileStorageShare{{Name: "share", Path: "/data", ACL: []upcloud.FileStorageACL{{Target: "*", Permission: "ro"}}}},
			Labels:           []upcloud.Label{{Key: "env", Value: "dev"}},
		}
		data, err := json.Marshal(&req)
		assert.NoError(t, err)
		const expected = `{"name":"full","zone":"zone-2","configured_status":"started","size_gib":20,"networks":[{"uuid":"net-uuid","name":"net1","family":"IPv4","ip_address":"192.168.1.1"}],"shares":[{"name":"share","path":"/data","acl":[{"target":"*","permission":"ro"}]}],"labels":[{"key":"env","value":"dev"}]}`
		assert.JSONEq(t, expected, string(data))
	})
}

func TestDeleteFileStorageRequest_RequestURL(t *testing.T) {
	req := DeleteFileStorageRequest{UUID: "fs-uuid"}
	assert.Equal(t, "/file-storage/fs-uuid", req.RequestURL())
}

func TestFileStorageRequest_InvalidValues(t *testing.T) {
	// Nil Page and Sort
	req := GetFileStoragesRequest{}
	assert.Equal(t, "/file-storage", req.RequestURL())

	// Malformed UUID
	getReq := GetFileStorageRequest{UUID: ""}
	assert.Equal(t, "/file-storage/", getReq.RequestURL())

	// ReplaceFileStorageRequest with empty UUID
	replaceReq := ReplaceFileStorageRequest{UUID: ""}
	assert.Equal(t, "/file-storage/", replaceReq.RequestURL())

	// ModifyFileStorageRequest with nil fields
	modReq := ModifyFileStorageRequest{UUID: ""}
	assert.Equal(t, "/file-storage/", modReq.RequestURL())

	// DeleteFileStorageRequest with empty UUID
	deleteReq := DeleteFileStorageRequest{UUID: ""}
	assert.Equal(t, "/file-storage/", deleteReq.RequestURL())
}
