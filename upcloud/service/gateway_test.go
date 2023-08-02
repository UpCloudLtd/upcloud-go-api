package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v6/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func TestGateway(t *testing.T) {
	t.Parallel()

	record(t, "gateway", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		router, err := svc.CreateRouter(ctx, &request.CreateRouterRequest{
			Name: "test-router",
		})
		if !assert.NoError(t, err) {
			return
		}
		gw, err := svc.CreateGateway(ctx, &request.CreateGatewayRequest{
			Name: "test",
			Zone: "pl-waw1",
			Features: []upcloud.GatewayFeature{
				upcloud.GatewayFeatureNAT,
			},
			Routers: []request.GatewayRouter{
				{
					UUID: router.UUID,
				},
			},
			ConfiguredStatus: upcloud.GatewayConfiguredStatusStarted,
		})
		if !assert.NoError(t, err) {
			return
		}

		if !assert.NoError(t, waitGatewayToStart(ctx, rec, svc, gw.UUID)) {
			return
		}
		gw, err = svc.GetGateway(ctx, &request.GetGatewayRequest{UUID: gw.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "test", gw.Name)
		assert.Equal(t, "pl-waw1", gw.Zone)
		if assert.GreaterOrEqual(t, 1, len(gw.Features)) {
			assert.Equal(t, upcloud.GatewayFeatureNAT, gw.Features[0])
		}
		if assert.Len(t, gw.Routers, 1) {
			assert.Equal(t, router.UUID, gw.Routers[0].UUID)
		}
		assert.Len(t, gw.Addresses, 1)
		gw, err = svc.ModifyGateway(ctx, &request.ModifyGatewayRequest{
			UUID:             gw.UUID,
			Name:             "new-name",
			ConfiguredStatus: upcloud.GatewayConfiguredStatusStopped,
		})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "new-name", gw.Name)
		assert.Equal(t, upcloud.GatewayConfiguredStatusStopped, gw.ConfiguredStatus)

		assert.NoError(t, svc.DeleteGateway(ctx, &request.DeleteGatewayRequest{UUID: gw.UUID}))

		if err := waitGatewayToDelete(ctx, rec, svc, gw.UUID); err != nil {
			t.Fatal(err)
		}
		if err := svc.DeleteRouter(ctx, &request.DeleteRouterRequest{UUID: router.UUID}); err != nil {
			t.Log(err)
		}
	})
}

func waitGatewayToStart(ctx context.Context, rec *recorder.Recorder, svc *Service, UUID string) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}

	const timeout = 10 * time.Minute

	rec.AddPassthrough(func(h *http.Request) bool {
		return true
	})
	defer func() {
		rec.Passthroughs = nil
	}()

	waitUntil := time.Now().Add(timeout)
	for {
		gw, err := svc.GetGateway(ctx, &request.GetGatewayRequest{UUID: UUID})
		if err != nil {
			return err
		}
		if gw.OperationalState == upcloud.GatewayOperationalStateRunning {
			return nil
		}
		if time.Now().After(waitUntil) {
			return fmt.Errorf("timeout %s reached", timeout.String())
		}
		time.Sleep(5 * time.Second)
	}
}

func waitGatewayToDelete(ctx context.Context, rec *recorder.Recorder, svc *Service, UUID string) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}

	const timeout = 10 * time.Minute

	rec.AddPassthrough(func(h *http.Request) bool {
		return true
	})
	defer func() {
		rec.Passthroughs = nil
	}()

	waitUntil := time.Now().Add(timeout)
	for {
		_, err := svc.GetGateway(ctx, &request.GetGatewayRequest{UUID: UUID})
		if err != nil {
			log.Printf("ERROR: %+v", err)
			var ucErr *upcloud.Problem
			if errors.As(err, &ucErr) && ucErr.Status == http.StatusNotFound {
				return nil
			}
			return err
		}
		if time.Now().After(waitUntil) {
			return fmt.Errorf("timeout %s reached", timeout.String())
		}
		time.Sleep(5 * time.Second)
	}
}
