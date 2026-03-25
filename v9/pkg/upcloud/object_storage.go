package upcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *ClientWithResponses) WaitForObjectStorageOperationalState(ctx context.Context, serviceUUID string, desiredState string) (*ObjectStorage2GetService200, error) {
	svcUUID, err := uuid.Parse(serviceUUID)
	if err != nil {
		return nil, err
	}

	return retry(ctx, func(_ int, _ context.Context) (*ObjectStorage2GetService200, error) {
		resp, retryErr := c.GetObjectStorageWithResponse(ctx, svcUUID)
		if retryErr != nil {
			return nil, retryErr
		}

		retryErr = serverError(resp.HTTPResponse)
		if retryErr != nil {
			return nil, retryErr
		}

		if len(resp.Body) == 0 {
			return nil, nil
		}

		var tmp struct {
			OperationalState *string `json:"operational_state"`
		}
		if retryErr = json.Unmarshal(resp.Body, &tmp); retryErr != nil {
			return nil, nil
		}

		if tmp.OperationalState != nil && *tmp.OperationalState == desiredState {
			if resp.JSON200 != nil {
				return resp.JSON200, nil
			}
			var dest ObjectStorage2GetService200
			if retryErr = json.Unmarshal(resp.Body, &dest); retryErr == nil {
				return &dest, nil
			}
			return nil, nil
		}

		return nil, nil
	})
}

func (c *ClientWithResponses) WaitForObjectStorageDeletion(ctx context.Context, serviceUUID string) error {
	svcUUID, err := uuid.Parse(serviceUUID)
	if err != nil {
		return err
	}

	return retryUntilNil(ctx, func(_ int, _ context.Context) (*ObjectStorage2GetService200, error) {
		resp, retryErr := c.GetObjectStorageWithResponse(ctx, svcUUID)
		if retryErr != nil {
			return nil, retryErr
		}

		retryErr = serverError(resp.HTTPResponse)
		if retryErr != nil {
			return nil, retryErr
		}

		if resp.HTTPResponse != nil && resp.HTTPResponse.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		return resp.JSON200, nil
	})
}

func (c *ClientWithResponses) WaitForObjectStorageBucketDeletion(ctx context.Context, serviceUUID string, bucketName string) error {
	svcUUID, err := uuid.Parse(serviceUUID)
	if err != nil {
		return err
	}

	return retryUntilNil(ctx, func(_ int, _ context.Context) (*ObjectStorage2BucketDetailResponse, error) {
		resp, retryErr := c.ListObjectStorageBucketMetricsWithResponse(ctx, svcUUID, nil)
		if retryErr != nil {
			return nil, retryErr
		}

		retryErr = serverError(resp.HTTPResponse)
		if retryErr != nil {
			return nil, retryErr
		}

		if resp.HTTPResponse != nil && resp.HTTPResponse.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		if resp.JSON200 == nil {
			return nil, nil
		}

		for _, b := range *resp.JSON200 {
			if b.Name != nil && *b.Name == bucketName {
				return &b, nil
			}
		}

		return nil, nil
	})
}

func serverError(r *http.Response) error {
	if r != nil && r.StatusCode >= 500 {
		return Retryable(fmt.Errorf("server error: %s", r.Status))
	}
	return nil
}
