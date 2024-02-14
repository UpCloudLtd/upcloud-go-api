package request

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v7/upcloud"
)

type FilterLabel struct {
	upcloud.Label
}

func (l FilterLabel) ToQueryParam() string {
	return fmt.Sprintf("label=%s=%s", l.Key, l.Value)
}

type FilterLabelKey struct {
	Key string
}

func (k FilterLabelKey) ToQueryParam() string {
	return fmt.Sprintf("label=%s", k.Key)
}
