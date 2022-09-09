package request

import (
	"fmt"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
)

type FilterLabel struct {
	upcloud.Label
}

func (l FilterLabel) ToQueryParam() string {
	return fmt.Sprintf("labels=%s=%s", l.Key, l.Value)
}

type FilterLabelKey struct {
	Key string
}

func (k FilterLabelKey) ToQueryParam() string {
	return fmt.Sprintf("labels=%s", k.Key)
}
