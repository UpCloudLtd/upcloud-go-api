package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/stretchr/testify/assert"
)

const exampleClusterResponse = `
{
	"name":"test-name",
	"network":"03e4970d-7791-4b80-a892-682ae0faf46b",
	"network_cidr":"176.16.1.0/24",
	"state":"running",
	"uuid":"0dca7a18-98e1-4e2d-aea5-ef5dd5fa450e",
	"zone":"de-fra1",
	"node_groups":[
		{
			"name":"my-group1",
			"plan":"2xCPU-4GB",
			"count":2,
			"storage":"01000000-0000-4000-8000-000160010100",
			"anti_affinity": false,
			"kubelet_args":[],
			"labels":[],
			"ssh_keys":[]
		},
		{
			"name":"my-group2",
			"plan":"4xCPU-8GB",
			"count":4,
			"storage":"01000000-0000-4000-8000-000160010100",
			"anti_affinity": true,
			"kubelet_args": [
				{
					"key": "log-flush-frequency",
					"value": "15s"
				}
			],
			"taints": [
				{
					"effect": "NoSchedule",
					"key": "environment",
					"value": "development"
				}
			],
			"labels": [
				{
					"key":"managedBy",
					"value":"go-sdk"
				}
			],
			"ssh_keys":[
				"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIO3fnjc8UrsYDNU8365mL3lnOPQJg18V42Lt8U/8Sm+r testy_test"
			]
		}
	]
}
`

const exampleNodeGroupResponse = `
{
	"name":"my-test-group",
	"plan":"4xCPU-8GB",
	"count":4,
	"storage":"01000000-0000-4000-8000-000160010100",
	"anti_affinity": true,
	"kubelet_args": [
		{
			"key": "log-flush-frequency",
			"value": "15s"
		}
	],
	"taints": [
		{
			"effect": "NoSchedule",
			"key": "environment",
			"value": "development"
		}
	],
	"labels": [
		{
			"key":"managedBy",
			"value":"go-sdk"
		}
	],
	"ssh_keys": [
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIO3fnjc8UrsYDNU8365mL3lnOPQJg18V42Lt8U/8Sm+r testy_test"
	]
}
`

const exampleNetworkResponse = `
{
	"network": {
		"name": "somenetwork",
		"type": "private",
		"uuid": "networkuuid",
		"zone": "de-fra1",
		"labels": [],
		"ip_networks": {
			"ip_network": [
				{
					"address": "172.16.1.0/24",
					"dhcp": "yes",
					"dhcp_default_route": "yes",
					"dhcp_dns": [],
					"family": "IPv4",
					"gateway": "172.16.1.0"
				}
			]
		}
	}
}
`

const exampleKubeconfigResponse = `{"kubeconfig":"apiVersion: v1\nclusters:\n- cluster:\n    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM2akNDQWRLZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1ERXdNekUwTVRJeE1sb1hEVE15TVRJek1URTBNVGN4TWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS2JJCnZYQkRKM1NSaHdkcG1LaVFjdFpMYkxVUU9sT29MMUVkNEx6aEd2cFBSUGd6VWFMSGFQSjZnTVdXTVliRGswMXMKZVF6VlY3dGZWWTNMQTFxWmhUVDgxNnYyNEdUcVEvSjdzOEs0REUyb2RNL0ZlUE1KRDEraWdKS1BtU2dOQmloawp2c2ZvU3hZanhoS0pXeTJ3OThYby9pUTlHNnhQNUpmUHdPcm9UNWFXTWZlblhjNytIZkRrN2V5SnBsdFF4ZmxqCjEvMHhQUnRNWWlqNnRYVmkzbHRYangrTnN5ZW9ueGdYYitNaU1LbVRKTmJnRU1lcnJXQnlQK3R5bnJUZGNnWXAKckFoOU5kS0I5T2x5KzV4c3Q3WHNJcXhaYkFQS05KQ0s1UUFTcTVJWmhENTVaUjlIaTFHUzF6ZTNsS2lKMTMyTgpLQ05YSUozajIxOGxsT2pCek84Q0F3RUFBYU5GTUVNd0RnWURWUjBQQVFIL0JBUURBZ0trTUJJR0ExVWRFd0VCCi93UUlNQVlCQWY4Q0FRQXdIUVlEVlIwT0JCWUVGR2tQL3dJampEVnVyVkJOeWNKazV6Y2FCems3TUEwR0NTcUcKU0liM0RRRUJDd1VBQTRJQkFRQXZCVUZnclErVEkyMnpuZlJZRG9WdW4xT2cydTAwMUVNbk1KN3NUR204UnBNdQpsZCtjNkZzY21SMHZ0K2pTS21oa3NKdXZ3U1BwZlV3OEJHNmlMem9KczlKSzV3a2psb0taYllYL2gyT2lheVdWCnM0TWx1TlpDU0wzR2pkaldxVnJNZ1J6RlUxQXJwNXNxR0Z2VWVERnpkRFVrR01rU3FYbWxyTGlCS29ZN3hmOU8KcDMvNWoxT0E0TzRtTnVvL2ZMVU85VDcyalBHTW1CZFg2TFU3VXcwMFdlejdDdC9CM2UwSDMvT1puSnpoQm51TQpQR0kyckRSSTJ0c2JLN1RseHEzTUFRSFY0Nk1LWmFpc2NlRHYralFEemcxenViRFFFcXhrNldncXc5aGR2YXlxCkJLbUFzeHQ4U3FPcEpZWFZKeXVHQjJvU2VKeW5lam96dWNvVUpRM0oKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\n    server: https://lb-0abebeafda2c4602a5e9a07f9108fb83-1.upcloudlb.com:6443\n  name: go-sdk-test-ctx\ncontexts:\n- context:\n    cluster: go-sdk-test-ctx\n    user: go-sdk-test-ctx-admin\n  name: go-sdk-test-ctx-admin@go-sdk-test-ctx\ncurrent-context: go-sdk-test-ctx-admin@go-sdk-test-ctx\nkind: Config\npreferences: {}\nusers:\n- name: go-sdk-test-ctx-admin\n  user:\n    client-certificate-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURFekNDQWZ1Z0F3SUJBZ0lJWUhNM2wrYk9vdGN3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TXpBeE1ETXhOREV5TVRKYUZ3MHlOREF4TURNeE5ERTNNVE5hTURReApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sa3dGd1lEVlFRREV4QnJkV0psY201bGRHVnpMV0ZrCmJXbHVNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQW4xbVpZYkZtY3dPMWsrcEUKdDgwVjBTL255eUY3Qmt6b3FlRXhEQXZZWTVyU1ZlbDNjQkNlV3g0b2VvM3I2YmZJQzhuYnVPWjNuR0cwQ1IzNgpiS3BXYTRzdkcyTFRCMlQ1bVk2QjBqY1lhcWZRa1dQNnN3Um1SLzJlYjBYenFST3J5Ly9rOEhhNTNGNEZzMmt2ClBzanhKK1ZpQkd6TWFKVFE2WVgxYTdYa3grQXlkTHIwTmtVcWk2RWIxelhUOFcrZ1pjek9uNXNmZytkMXdTaEgKUk9CSU94L0k4YXY5dUQ4YlEzclF2TVloR0wwSVRZcmxaUmhkUG12d2JzY2p6cUJlN2FmWnF5ZmZBRGxZOXJ3eQpCVGlsRk1oTEJUL0JVb3ZtZGZZUkpSOW8xRzQ3S2dFRVJSVS9DeUd6WHBzS2k3ZDBtaW4vU1NYVVJqMU44TlpuCkFzWS9jd0lEQVFBQm8wZ3dSakFPQmdOVkhROEJBZjhFQkFNQ0JhQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUgKQXdJd0h3WURWUjBqQkJnd0ZvQVVhUS8vQWlPTU5XNnRVRTNKd21Ubk54b0hPVHN3RFFZSktvWklodmNOQVFFTApCUUFEZ2dFQkFLSjE1VTIzbWtrNTlsQ2Q3V292aUVMSEorNGFtZnVHU2JyUDRCVjBYSWZGZ0lLWkk4RzRTdU8zCkZMQWZNMnphV1FPZ0ZCanpNaUZXbExFZFVTU0pKUm5ZNVdONmlBRmtydWpEbHFZakxkTE16V1Q4ZytSTnJ6WDAKSWlYbnoyaW5lVDRNZXB0VmsrVm4yeFpWd2lYbmJOR3BCVUM5QU10REU5WlllZnFNOU9qZG4valk4TzJUTDJaUApLTTBqamJHWmZYREdpNEE0aHVlcGRHRlBVZjJmWU9zQk1ZTEo1VUVZbWpJUnZGYi94S1NNdUNxRHRFVzI2VmJiCnhQc1gxcEhOUXhDVkNIK1ZjeUtJejNPK0NrNFo1ZmY0QWhqNU4zc0xkOUZ2RnpmVG5JWVk3WGZvUlZUdFVKb3AKL09NQ1RKVDZGK3N5dmg5ZUxROE5lZ0VDc1c4YkRYST0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=\n    client-key-data: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBbjFtWlliRm1jd08xaytwRXQ4MFYwUy9ueXlGN0Jrem9xZUV4REF2WVk1clNWZWwzCmNCQ2VXeDRvZW8zcjZiZklDOG5idU9aM25HRzBDUjM2YktwV2E0c3ZHMkxUQjJUNW1ZNkIwamNZYXFmUWtXUDYKc3dSbVIvMmViMFh6cVJPcnkvL2s4SGE1M0Y0RnMya3ZQc2p4SitWaUJHek1hSlRRNllYMWE3WGt4K0F5ZExyMApOa1VxaTZFYjF6WFQ4VytnWmN6T241c2ZnK2Qxd1NoSFJPQklPeC9JOGF2OXVEOGJRM3JRdk1ZaEdMMElUWXJsClpSaGRQbXZ3YnNjanpxQmU3YWZacXlmZkFEbFk5cnd5QlRpbEZNaExCVC9CVW92bWRmWVJKUjlvMUc0N0tnRUUKUlJVL0N5R3pYcHNLaTdkMG1pbi9TU1hVUmoxTjhOWm5Bc1kvY3dJREFRQUJBb0lCQUh6Rzlsb1BSYi9PS2NNWApjSlBVWGI4ZUdnMXZ2QnZrNFZNVTZRa2J4V2ZKZGVhY0dGQ0NVdDNhc2F1MXNnT2pTMXdmeHBQMHM3aWFzUlZxCmlndkpIajY4RURrTG0xOXc3Qm9ZQXdRTzdHbW4ycVBlZkJMdDRRR0NVU3VreFBXaVY2WTRUSmNYQU5iVU1QYVoKNk1yckloc3hmUjBsN0xIL3hjNzJmSVRKTzhwZVFjR2Rpa2Z2UENQVk9sWG1HNFVyVEd5WU9JS0N1TWZGYkxDUgpZR09HT1VDaHlEY1JYRjBRdTlrWU1TNlFTV1RiazhpZmhEdDNxZndnZ0pTU3FYYkNqSFZETE04cXM4ODhlbnFZCmJRclBwUEJMRXdhY0dOaEFTdHJ0Mm1UVFhZUE01TTJLQ0d2dkhJblpseUxYeEF5N3graVNLNWV3amY5TEpFUFMKalJhZFhQRUNnWUVBMDA5bGQyZnJjMXBINHRXcWdlRkxRTG00Si8vcG91VmhVaHdwUm4vbjlNaFpIUUlXbG02VQo4Ykxpa0RXckxMaGE2MHBkSnF6Z0o4NE1BbjdpVUZ2OVR3L1M4M0VvOHhLU2oxN05pMDZFNGNlVDZ4NkpNR2U1CnB1d2lCcTFabHBDVVc1bk1zNkxqUzZzb2pvWGI2V3ZYb2JRek52N1c0QkRVZGJyd2U1bHowbjhDZ1lFQXdRMEQKZFFSQ05ac29aNE5oL1J6SjhiVEhBNmZlU1g2ZDBvV0NJdEVDVytVRWdVQUFHYlJWVitFQ1d5d3RjVFA5bnlsVgprYVlJTVRDVzIvdXJwY2FrcnFHdW5nQlV4TlFvcU5xRHBrNzFOcHNzd1cyVG82bU9WdFFmbnI0a1NydElwUjJuCmVFTXRLMW5nYW41MXpjM1ZrSE1lUjlFcXVCRm9KZnlKUjAzRjhRMENnWUFzelMyZkptcFdON0w4RmY4anNHZXIKSG5VOERkYzBVVnZUOCtLUWJ2ZjMveTVkcHg2dzRGczE3NDUzc3RsTER2OC8yYkZzVE1UdHk1TGlTSktsSlF2TQo0bmNBWkdLaFByUFNML0IyYzd4YXZselBRZGNtYXllQ2k1Z2ZRRXU1VDRnVTc4UGw1a2c0dDk1T2NYVVJ4V3FmCmZ6eElrYmx2SWtKWmY3d1RDSkwwOHdLQmdCN2piQ2d1OUVFTGdGMS9vaEFUcVNRcFhxWnhGNW4remxXKzB1R3IKTzR0a2kySDYwWWxMVjFSbVJqQThVMUVIbG83KzJCZ2p6cS9BTElNcU0zNmsxL2l4R2ZWaDdDNUtHKzlZZHpRRwo4aWs0anVLc0c3RldFZFpHdm9CeHYvRE5Jb1I0eUREdjNxeFBCeDlDK0twalFzYUlCTkpNMzE3d09nL0o1bFVvClZvR0pBb0dCQUtUZGt0SjU4bHQ3WXplTGViZU9qTHdMOTVzM1RlckppcmxJdTdHTjllVlR1VWZtZ3MyTWVXdzcKMEpsVy9pU2JDR1dmSko1N3ZZVHdHbi84N0dBNXhaRWxOUG8xQ3NoSVF6YXhDQzZpMS9yZUswTlJKYnZoUTQwaQp1eG9xZk91aVdjRkNXbTg5ektpeE5KbkRxd082UE02dzgyczR1RFRXQXdqTG05eDdvblQxCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==\n"}`

func TestGetKubernetesClusters(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes", client.APIVersion), r.URL.Path)
		fmt.Fprintf(w, "[%s]", exampleClusterResponse)
	}))
	defer srv.Close()
	res, err := svc.GetKubernetesClusters(context.Background(), &request.GetKubernetesClustersRequest{})
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	// Check cluster properties
	assert.Equal(t, "test-name", res[0].Name)
	assert.Equal(t, upcloud.KubernetesClusterStateRunning, res[0].State)
	assert.Equal(t, "de-fra1", res[0].Zone)
	assert.Len(t, res[0].NodeGroups, 2)

	// Check first group properties
	assert.Equal(t, "my-group1", res[0].NodeGroups[0].Name)
	assert.Equal(t, "2xCPU-4GB", res[0].NodeGroups[0].Plan)
	assert.False(t, res[0].NodeGroups[0].AntiAffinity)
	assert.Equal(t, 2, res[0].NodeGroups[0].Count)
	assert.Len(t, res[0].NodeGroups[0].Labels, 0)
	assert.Len(t, res[0].NodeGroups[0].KubeletArgs, 0)
	assert.Len(t, res[0].NodeGroups[0].Taints, 0)

	// Check second group properties
	assert.Equal(t, "my-group2", res[0].NodeGroups[1].Name)
	assert.Equal(t, "4xCPU-8GB", res[0].NodeGroups[1].Plan)
	assert.True(t, res[0].NodeGroups[1].AntiAffinity)
	assert.Equal(t, 4, res[0].NodeGroups[1].Count)
	assert.Len(t, res[0].NodeGroups[1].Labels, 1)
	assert.Len(t, res[0].NodeGroups[1].KubeletArgs, 1)
	assert.Len(t, res[0].NodeGroups[1].Taints, 1)
}

func TestGetKubernetesClusterDetails(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes/_UUID_", client.APIVersion), r.URL.Path)
		fmt.Fprint(w, exampleClusterResponse)
	}))
	defer srv.Close()
	res, err := svc.GetKubernetesCluster(context.Background(), &request.GetKubernetesClusterRequest{UUID: "_UUID_"})
	assert.NoError(t, err)

	// Check cluster properties
	assert.Equal(t, "test-name", res.Name)
	assert.Equal(t, upcloud.KubernetesClusterStateRunning, res.State)
	assert.Equal(t, "de-fra1", res.Zone)
	assert.Len(t, res.NodeGroups, 2)
}

func TestCreateKubernetesCluster(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CreateKubernetesCluster method first makes a request to /network/:uuid to check network CIDR
		if r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/%s/network/03e4970d-7791-4b80-a892-682ae0faf46b", client.APIVersion) {
			fmt.Fprint(w, exampleNetworkResponse)
			return
		}

		if r.Method == http.MethodPost && r.URL.Path == fmt.Sprintf("/%s/kubernetes", client.APIVersion) {
			payload := request.CreateKubernetesClusterRequest{}
			err := json.NewDecoder(r.Body).Decode(&payload)
			assert.NoError(t, err)

			fmt.Fprint(w, exampleClusterResponse)
			return
		}

		t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()
	res, err := svc.CreateKubernetesCluster(context.Background(), &request.CreateKubernetesClusterRequest{
		Name:    "test-name",
		Network: "03e4970d-7791-4b80-a892-682ae0faf46b",
		Zone:    "de-fra1",
		NodeGroups: []request.KubernetesNodeGroup{
			{
				Name:         "my-group1",
				Plan:         "2xCPU-4GB",
				Count:        2,
				Storage:      "01000000-0000-4000-8000-000160010100",
				AntiAffinity: false,
			},
			{
				Name:         "my-group2",
				Plan:         "4xCPU-8GB",
				Count:        4,
				Storage:      "01000000-0000-4000-8000-000160010100",
				AntiAffinity: true,
				KubeletArgs: []upcloud.KubernetesKubeletArg{
					{
						Key:   "log-flush-frequency",
						Value: "15s",
					},
				},
				Taints: []upcloud.KubernetesTaint{
					{
						Effect: "NoSchedule",
						Key:    "environment",
						Value:  "development",
					},
				},
				Labels: []upcloud.Label{
					{
						Key:   "managedBy",
						Value: "go-sdk",
					},
				},
				SSHKeys: []string{
					"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIO3fnjc8UrsYDNU8365mL3lnOPQJg18V42Lt8U/8Sm+r testy_test",
				},
			},
		},
	})
	assert.NoError(t, err)

	// Check cluster properties
	assert.Equal(t, "test-name", res.Name)
	assert.Equal(t, upcloud.KubernetesClusterStateRunning, res.State)
	assert.Equal(t, "de-fra1", res.Zone)
	assert.Len(t, res.NodeGroups, 2)
}

func TestDeleteKubernetesCluster(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes/_UUID_", client.APIVersion), r.URL.Path)
	}))
	defer srv.Close()

	err := svc.DeleteKubernetesCluster(context.Background(), &request.DeleteKubernetesClusterRequest{
		UUID: "_UUID_",
	})
	assert.NoError(t, err)
}

func TestGetKubernetesNodeGroups(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes/_UUID_/node-groups", client.APIVersion), r.URL.Path)
		fmt.Fprintf(w, "[%s]", exampleNodeGroupResponse)
	}))
	defer srv.Close()

	res, err := svc.GetKubernetesNodeGroups(context.Background(), &request.GetKubernetesNodeGroupsRequest{
		ClusterUUID: "_UUID_",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestGetKubernetesNodeGroup(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes/_UUID_/node-groups/_NAME_", client.APIVersion), r.URL.Path)
		fmt.Fprint(w, exampleNodeGroupResponse)
	}))
	defer srv.Close()

	res, err := svc.GetKubernetesNodeGroup(context.Background(), &request.GetKubernetesNodeGroupRequest{
		ClusterUUID: "_UUID_",
		Name:        "_NAME_",
	})
	assert.NoError(t, err)
	assert.Equal(t, "my-test-group", res.Name)
}

func TestCreateKubernetesNodeGroup(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes/_UUID_/node-groups", client.APIVersion), r.URL.Path)

		payload := request.CreateKubernetesNodeGroupRequest{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		fmt.Fprint(w, exampleNodeGroupResponse)
	}))
	defer srv.Close()

	res, err := svc.CreateKubernetesNodeGroup(context.Background(), &request.CreateKubernetesNodeGroupRequest{
		ClusterUUID: "_UUID_",
		NodeGroup: request.KubernetesNodeGroup{
			Count:        4,
			Name:         "my-test-group",
			Plan:         "2xCPU-4GB",
			Storage:      "01000000-0000-4000-8000-000160010100",
			AntiAffinity: true,
			KubeletArgs: []upcloud.KubernetesKubeletArg{
				{
					Key:   "log-flus-frequency",
					Value: "15s",
				},
			},
			Taints: []upcloud.KubernetesTaint{
				{
					Effect: "NoSchedule",
					Key:    "environment",
					Value:  "development",
				},
			},
			Labels: []upcloud.Label{
				{
					Key:   "managedBy",
					Value: "go-sdk",
				},
			},
			SSHKeys: []string{
				"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIO3fnjc8UrsYDNU8365mL3lnOPQJg18V42Lt8U/8Sm+r testy_test",
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "my-test-group", res.Name)
}

func TestModifyKubernetesNodeGroup(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes/_UUID_/node-groups/_NAME_", client.APIVersion), r.URL.Path)

		payload := request.ModifyKubernetesNodeGroup{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		assert.NoError(t, err)

		fmt.Fprint(w, exampleNodeGroupResponse)
	}))
	defer srv.Close()

	res, err := svc.ModifyKubernetesNodeGroup(context.Background(), &request.ModifyKubernetesNodeGroupRequest{
		ClusterUUID: "_UUID_",
		Name:        "_NAME_",
		NodeGroup:   request.ModifyKubernetesNodeGroup{Count: 4},
	})
	assert.NoError(t, err)
	assert.Equal(t, 4, res.Count)
}

func TestDeleteKubernetesNodeGroup(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes/_UUID_/node-groups/_NAME_", client.APIVersion), r.URL.Path)
		fmt.Fprint(w, exampleNodeGroupResponse)
	}))
	defer srv.Close()

	err := svc.DeleteKubernetesNodeGroup(context.Background(), &request.DeleteKubernetesNodeGroupRequest{
		ClusterUUID: "_UUID_",
		Name:        "_NAME_",
	})
	assert.NoError(t, err)
}

func TestWaitForKubernetesClusterState(t *testing.T) {
	t.Parallel()

	requestsCounter := 0
	requestsMade := 0

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, fmt.Sprintf("/%s/kubernetes/_UUID_", client.APIVersion), r.URL.Path)
		requestsMade++

		if requestsCounter >= 2 {
			fmt.Fprint(w, `
			{
				"name":"test-name",
				"network":"03e4970d-7791-4b80-a892-682ae0faf46b",
				"network_cidr":"176.16.1.0/24",
				"state":"running",
				"uuid":"0dca7a18-98e1-4e2d-aea5-ef5dd5fa450e",
				"zone":"de-fra1",
				"node_groups": []
			}
			`)
		} else {
			requestsCounter++
			fmt.Fprint(w, `
			{
				"name":"test-name",
				"network":"03e4970d-7791-4b80-a892-682ae0faf46b",
				"network_cidr":"176.16.1.0/24",
				"state":"pending",
				"uuid":"0dca7a18-98e1-4e2d-aea5-ef5dd5fa450e",
				"zone":"de-fra1",
				"node_groups": []
			}
			`)
		}
	}))
	defer srv.Close()

	_, err := svc.WaitForKubernetesClusterState(context.Background(), &request.WaitForKubernetesClusterStateRequest{
		UUID:         "_UUID_",
		DesiredState: upcloud.KubernetesClusterStateRunning,
		Timeout:      time.Second * 20,
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, requestsMade)
}

func TestGetKubernetesKubeconfig(t *testing.T) {
	t.Parallel()

	srv, svc := setupTestServerAndService(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GetKubernetesKubeconfig first fetches cluster details to check for running state, so we must
		// take care of both requests
		if r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/%s/kubernetes/_UUID_", client.APIVersion) {
			fmt.Fprint(w, exampleClusterResponse)
			return
		}

		if r.Method == http.MethodGet && r.URL.Path == fmt.Sprintf("/%s/kubernetes/_UUID_/kubeconfig", client.APIVersion) {
			fmt.Fprint(w, exampleKubeconfigResponse)
			return
		}

		t.Fatalf("Unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	res, err := svc.GetKubernetesKubeconfig(context.Background(), &request.GetKubernetesKubeconfigRequest{
		UUID: "_UUID_",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}
