package upcloud

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *ClientWithResponses) WaitForObjectStorageOperationalState(ctx context.Context, serviceUUID string, desiredState string) (*GetObjectStorage200, error) {
	svcUUID, err := uuid.Parse(serviceUUID)
	if err != nil {
		return nil, err
	}

	return retry(ctx, func(_ int, _ context.Context) (*GetObjectStorage200, error) {
		httpResponse, retryErr := c.GetObjectStorage(ctx, svcUUID)
		if retryErr != nil {
			return nil, retryErr
		}

		if retryErr = objectStorageWaitStatusError(httpResponse); retryErr != nil {
			closeObjectStorageWaitResponse(httpResponse)
			return nil, retryErr
		}
		resp, retryErr := ParseGetObjectStorageResp(httpResponse)
		if retryErr != nil {
			return nil, retryErr
		}
		if retryErr = objectStorageWaitResponseError(resp.HTTPResponse, resp.Body); retryErr != nil {
			return nil, retryErr
		}

		if resp.JSON200 == nil || resp.JSON200.OperationalState == nil {
			return nil, errors.New("object storage wait: missing operational state in HTTP 200 response")
		}
		if string(*resp.JSON200.OperationalState) == desiredState {
			return resp.JSON200, nil
		}

		return nil, nil
	})
}

func (c *ClientWithResponses) WaitForObjectStorageDeletion(ctx context.Context, serviceUUID string) error {
	svcUUID, err := uuid.Parse(serviceUUID)
	if err != nil {
		return err
	}

	return retryUntilNil(ctx, func(_ int, _ context.Context) (*GetObjectStorage200, error) {
		httpResponse, retryErr := c.GetObjectStorage(ctx, svcUUID)
		if retryErr != nil {
			return nil, retryErr
		}

		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			closeObjectStorageWaitResponse(httpResponse)
			return nil, nil
		}
		if retryErr = objectStorageWaitStatusError(httpResponse); retryErr != nil {
			closeObjectStorageWaitResponse(httpResponse)
			return nil, retryErr
		}
		resp, retryErr := ParseGetObjectStorageResp(httpResponse)
		if retryErr != nil {
			return nil, retryErr
		}
		if retryErr = objectStorageWaitResponseError(resp.HTTPResponse, resp.Body); retryErr != nil {
			return nil, retryErr
		}
		if resp.JSON200 == nil {
			return nil, errors.New("object storage deletion wait: missing payload in HTTP 200 response")
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
		httpResponse, retryErr := c.ListObjectStorageBucketMetrics(ctx, svcUUID, nil)
		if retryErr != nil {
			return nil, retryErr
		}

		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			closeObjectStorageWaitResponse(httpResponse)
			return nil, nil
		}
		if retryErr = objectStorageWaitStatusError(httpResponse); retryErr != nil {
			closeObjectStorageWaitResponse(httpResponse)
			return nil, retryErr
		}
		resp, retryErr := ParseListObjectStorageBucketMetricsResp(httpResponse)
		if retryErr != nil {
			return nil, retryErr
		}
		if retryErr = objectStorageWaitResponseError(resp.HTTPResponse, resp.Body); retryErr != nil {
			return nil, retryErr
		}
		if resp.JSON200 == nil {
			return nil, errors.New("object storage bucket deletion wait: missing payload in HTTP 200 response")
		}

		for _, b := range *resp.JSON200 {
			if b.Name == nil || *b.Name == "" {
				return nil, errors.New("object storage bucket deletion wait: missing bucket name in HTTP 200 response")
			}
			if *b.Name == bucketName {
				return &b, nil
			}
		}

		return nil, nil
	})
}

func closeObjectStorageWaitResponse(response *http.Response) {
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
}

func objectStorageWaitResponseError(response *http.Response, body []byte) error {
	if err := objectStorageWaitStatusError(response); err != nil {
		return err
	}
	payload := bytes.TrimSpace(body)
	if len(payload) == 0 || bytes.Equal(payload, []byte("null")) {
		return errors.New("object storage wait: missing payload in HTTP 200 response")
	}
	return nil
}

func objectStorageWaitStatusError(response *http.Response) error {
	if response == nil {
		return errors.New("object storage wait: missing HTTP response")
	}
	if response.StatusCode == http.StatusOK {
		return nil
	}
	err := fmt.Errorf("object storage wait: unexpected HTTP status %d (%s)", response.StatusCode, http.StatusText(response.StatusCode))
	if response.StatusCode >= 500 {
		return Retryable(err)
	}
	return err
}
