package upcloud

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnMarshalLoadBalancer(t *testing.T) {
	lbString := `[{"backends":[],"configured_status":"started","created_at":"2022-02-23T07:39:44.318844Z","dns_name":"lb-0a70595df2ef49338d2f7b1931a62bc4.upcloudlb.com","frontends":[],"name":"new-name","network_uuid":"032d4c7f-61b5-4ea9-a2d6-d2357c3c9a88","operational_state":"running","plan":"development","updated_at":"2022-02-23T09:01:24.718074Z","uuid":"0a70595d-f2ef-4933-8d2f-7b1931a62bc4","zone":"es-mad1"}]`
	lbs := make([]LoadBalancer, 0)
	err := json.Unmarshal([]byte(lbString), &lbs)
	assert.NoError(t, err)
}
