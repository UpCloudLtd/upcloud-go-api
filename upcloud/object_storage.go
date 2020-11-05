package upcloud

import "encoding/json"

// ObjectStorage represents a object storage
type ObjectStorage struct {
	Created     string `json:"created"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Size        int    `json:"size"`
	State       string `json:"state"`
	URL         string `json:"url"`
	UUID        string `json:"uuid"`
	Zone        string `json:"zone"`
}

// ObjectStorages represent a /object-storage response
type ObjectStorages struct {
	ObjectStorages []ObjectStorage `json:"object_storages"`
}

// UnmarshalJSON is a custom unmarshaller that deals with
// deeply embedded values.
func (o *ObjectStorages) UnmarshalJSON(b []byte) error {
	type objectStorageWrapper struct {
		ObjectStorages []ObjectStorage `json:"object_storage"`
	}

	v := struct {
		ObjectStorages objectStorageWrapper `json:"object_storages"`
	}{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	o.ObjectStorages = v.ObjectStorages.ObjectStorages

	return nil
}
