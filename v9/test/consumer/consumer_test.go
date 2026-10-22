package consumer_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v9/pkg/upcloud"
	"github.com/google/uuid"
)

func TestPublicClient(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer consumer-token" {
			t.Error("public constructor did not install bearer authentication")
		}
		writer.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case "/1.3/router":
			_, _ = fmt.Fprint(writer, `{"routers":{"router":[]}}`)
		case "/1.3/object-storage-2":
			_, _ = fmt.Fprint(writer, `[]`)
		default:
			t.Errorf("unexpected request: %s", request.URL.Path)
			writer.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, err := upcloud.New("consumer-token", upcloud.WithBaseURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	routers, err := client.ListRoutersWithResponse(ctx, &upcloud.ListRoutersParams{})
	if err != nil || routers == nil || routers.JSON200 == nil {
		t.Fatalf("typed router response: %v, %v", routers, err)
	}
	services, err := client.ListObjectStoragesWithResponse(ctx, nil)
	if err != nil || services == nil || services.JSON200 == nil || len(*services.JSON200) != 0 {
		t.Fatalf("typed object storage response: %v, %v", services, err)
	}
}

func TestPublicTypes(t *testing.T) {
	t.Parallel()
	name := "consumer"
	model := upcloud.ObjectStorage2ServiceDetailResponse{Name: &name}
	if model.Name == nil || *model.Name != name {
		t.Fatal("optional model field is not accessible")
	}
	request := upcloud.ObjectStorage2ServiceCreate{
		Name: name, Region: "europe-1",
		ConfiguredStatus: upcloud.ObjectStorage2PropertyConfiguredStatusStarted,
	}
	if request.Name != name {
		t.Fatal("request model is not accessible")
	}
	var client *upcloud.ClientWithResponses
	var _ func(context.Context, uuid.UUID, ...upcloud.RequestEditorFn) (*upcloud.GetObjectStorageResp, error) = client.GetObjectStorageWithResponse
	var _ func(context.Context, string, string) (*upcloud.GetObjectStorage200, error) = client.WaitForObjectStorageOperationalState
	var _ func(context.Context, string) error = client.WaitForObjectStorageDeletion
	var _ func(context.Context, string, string) error = client.WaitForObjectStorageBucketDeletion
	var _ func(...upcloud.ClientOption) (*upcloud.ClientWithResponses, error) = upcloud.NewFromEnv
	response := upcloud.GetObjectStorageResp{}
	_ = response.ApplicationproblemJSONDefault
	_ = response.StatusCode()
}
