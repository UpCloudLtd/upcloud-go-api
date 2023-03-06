package request

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/stretchr/testify/assert"
)

func TestEncodeQueryFilters(t *testing.T) {
	want := "label=env%3Dprod&label=v2"
	got := encodeQueryFilters([]QueryFilter{
		FilterLabel{
			Label: upcloud.Label{
				Key:   "env",
				Value: "prod",
			},
		},
		FilterLabelKey{Key: "v2"},
	})
	assert.Equal(t, want, got)
}
