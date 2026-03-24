package upcloud

import (
	"context"
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
		rsp, err := c.GetObjectStorageWithResponse(ctx, svcUUID)
		if err != nil {
			return nil, err
		}

		err = serverError(rsp.HTTPResponse)
		if err != nil {
			return nil, err
		}

		if rsp.JSON200 == nil || rsp.JSON200.OperationalState == nil {
			return nil, nil
		}
		if *rsp.JSON200.OperationalState == desiredState {
			return rsp.JSON200, nil
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
		rsp, err := c.GetObjectStorageWithResponse(ctx, svcUUID)
		if err != nil {
			return nil, err
		}

		err = serverError(rsp.HTTPResponse)
		if err != nil {
			return nil, err
		}

		if rsp.HTTPResponse != nil && rsp.HTTPResponse.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		return rsp.JSON200, nil
	})
}

func (c *ClientWithResponses) WaitForObjectStorageBucketDeletion(ctx context.Context, serviceUUID string, bucketName string) error {
	svcUUID, err := uuid.Parse(serviceUUID)
	if err != nil {
		return err
	}

	return retryUntilNil(ctx, func(_ int, _ context.Context) (*ObjectStorage2BucketDetailResponse, error) {
		rsp, err := c.ListObjectStorageBucketMetricsWithResponse(ctx, svcUUID, nil)
		if err != nil {
			return nil, err
		}

		err = serverError(rsp.HTTPResponse)
		if err != nil {
			return nil, err
		}

		if rsp.HTTPResponse != nil && rsp.HTTPResponse.StatusCode == http.StatusNotFound {
			return nil, nil
		}

		if rsp.JSON200 == nil {
			return nil, nil
		}

		for _, b := range *rsp.JSON200 {
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
