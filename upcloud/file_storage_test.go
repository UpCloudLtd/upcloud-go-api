package upcloud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileStorage_MarshalUnmarshalJSON(t *testing.T) {
	cases := []FileStorage{
		{
			UUID:             "fs-uuid",
			Name:             "test-storage",
			Zone:             "fi-hel1",
			SizeGiB:          100,
			ConfiguredStatus: "started",
			OperationalState: "running",
			CreatedAt:        time.Date(2025, 9, 18, 12, 0, 0, 0, time.UTC),
			UpdatedAt:        time.Date(2025, 9, 18, 12, 5, 0, 0, time.UTC),
			Networks:         []FileStorageNetwork{{UUID: "net-uuid", Name: "net1", Family: "IPv4", IPAddress: "192.168.1.1"}},
			Shares:           []FileStorageShare{{Name: "share1", Path: "/data/share1", ACL: []FileStorageACL{{Target: "user", Permission: "rw"}}, Deleting: true}},
			Labels:           []Label{{Key: "env", Value: "dev"}},
			StateMessages:    []FileStorageStateMessage{{OperationalState: "running", Message: "All good", Code: "OK"}},
		},
		{
			UUID:             "",
			Name:             "",
			Zone:             "",
			SizeGiB:          0,
			ConfiguredStatus: "",
			CreatedAt:        time.Time{},
			UpdatedAt:        time.Time{},
			Networks:         nil,
			Shares:           nil,
			Labels:           nil,
			StateMessages:    nil,
		},
		{
			UUID:             "fs-uuid2",
			Name:             "storage2",
			Zone:             "fi-hel2",
			SizeGiB:          200,
			ConfiguredStatus: "stopped",
			OperationalState: "stopped",
			CreatedAt:        time.Now().UTC(),
			UpdatedAt:        time.Now().UTC(),
			Networks:         []FileStorageNetwork{},
			Shares:           []FileStorageShare{},
			Labels:           []Label{{Key: "env", Value: "test"}},
			StateMessages:    []FileStorageStateMessage{},
		},
	}
	for _, fs := range cases {
		data, err := json.Marshal(fs)
		assert.NoError(t, err)

		var out FileStorage
		err = json.Unmarshal(data, &out)
		assert.NoError(t, err)
		assert.Equal(t, fs, out)
	}
}

func TestFileStorage_MarshalUnmarshalJSON_Errors(t *testing.T) {
	invalidJSON := []string{
		"{invalid}",
		"{\"uuid\":123}", // wrong type
	}
	for _, s := range invalidJSON {
		var fs FileStorage
		err := json.Unmarshal([]byte(s), &fs)
		assert.Error(t, err)
	}
}

func TestFileStorage_EmptyFields(t *testing.T) {
	fs := FileStorage{}
	data, err := json.Marshal(fs)
	assert.NoError(t, err)
	var out FileStorage
	assert.NoError(t, json.Unmarshal(data, &out))
	assert.Equal(t, fs, out)
}
