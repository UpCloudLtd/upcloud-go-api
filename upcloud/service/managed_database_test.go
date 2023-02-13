package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
)

func TestService_CloneManagedDatabase(t *testing.T) {
	record(t, "clonemanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		var cloneDetails *upcloud.ManagedDatabase
		createReq := getTestCreateRequest("clonemanageddatabase")
		serviceDetails, err := svc.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
			if cloneDetails != nil {
				t.Logf("deleting clone %s", cloneDetails.UUID)
				err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: cloneDetails.UUID})
				assert.NoError(t, err)
			}
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, serviceDetails.UUID)
		require.NoError(t, err)

		err = waitForManagedDatabaseInitialBackup(ctx, rec, svc, serviceDetails.UUID)
		require.NoError(t, err)

		serviceDetails, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: serviceDetails.UUID})
		if !assert.NoError(t, err) {
			return
		}

		cloneReq := &request.CloneManagedDatabaseRequest{
			UUID:           serviceDetails.UUID,
			CloneTime:      serviceDetails.Backups[0].BackupTime.Add(1 * time.Second),
			HostNamePrefix: "testclone",
			Title:          "testclone",
			Zone:           createReq.Zone,
			Plan:           createReq.Plan,
		}
		cloneDetails, err = svc.CloneManagedDatabase(ctx, cloneReq)
		if !assert.NoError(t, err) {
			return
		}
	})
}

func TestService_CreateManagedDatabase(t *testing.T) {
	typesToTest := []upcloud.ManagedDatabaseServiceType{
		upcloud.ManagedDatabaseServiceTypeMySQL,
		upcloud.ManagedDatabaseServiceTypePostgreSQL,
		upcloud.ManagedDatabaseServiceTypeRedis,
	}
	record(t, "createmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		for _, serviceType := range typesToTest {
			t.Run(string(serviceType), func(t *testing.T) {
				req := getTestCreateRequest("createmanageddatabase")
				req.Type = serviceType
				if serviceType == upcloud.ManagedDatabaseServiceTypeRedis {
					req.Plan = "1x1xCPU-2GB"
				}
				details, err := svc.CreateManagedDatabase(ctx, req)
				if !assert.NoError(t, err) {
					return
				}
				defer func() {
					t.Logf("deleting %s", details.UUID)
					err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
					assert.NoError(t, err)
				}()
				assert.Equal(t, serviceType, details.Type)
				assert.EqualValues(t, req.Maintenance, details.Maintenance)
				assert.Equal(t, req.Plan, details.Plan)
				assert.Equal(t, req.Properties.GetIPFilter(),
					details.Properties.GetIPFilter())
				assert.Equal(t, req.Properties.GetAutoUtilityIPFilter(),
					details.Properties.GetAutoUtilityIPFilter())
				assert.Equal(t, req.Properties.GetPublicAccess(),
					details.Properties.GetPublicAccess())
				assert.Equal(t, req.Title, details.Title)
				assert.Equal(t, req.Type, details.Type)
				assert.Equal(t, req.Zone, details.Zone)
				switch serviceType {
				case upcloud.ManagedDatabaseServiceTypeRedis:
					assert.NotEmpty(t, details.Metadata.RedisVersion)
				case upcloud.ManagedDatabaseServiceTypeMySQL:
					assert.NotEmpty(t, details.Metadata.MySQLVersion)
				case upcloud.ManagedDatabaseServiceTypePostgreSQL:
					assert.NotEmpty(t, details.Metadata.PGVersion)
				}
			})
		}
	})
}

func TestService_WaitForManagedDatabaseState(t *testing.T) {
	record(t, "waitformanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("waitformanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		newDetails, err := svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{
			UUID: details.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, upcloud.ManagedDatabaseStateRunning, newDetails.State)
	})
}

func TestService_GetManagedDatabase(t *testing.T) {
	record(t, "getmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		req := getTestCreateRequest("getmanageddatabase")
		details, err := svc.CreateManagedDatabase(ctx, req)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.EqualValues(t, req.Maintenance, details.Maintenance)
		assert.Equal(t, req.Properties.GetIPFilter(),
			details.Properties.GetIPFilter())
		assert.Equal(t, req.Properties.GetAutoUtilityIPFilter(),
			details.Properties.GetAutoUtilityIPFilter())
		assert.Equal(t, req.Properties.GetPublicAccess(),
			details.Properties.GetPublicAccess())
		assert.Equal(t, req.Title, details.Title)
		assert.Equal(t, req.Type, details.Type)
		assert.Equal(t, req.Zone, details.Zone)
	})
}

func TestService_GetManagedDatabases(t *testing.T) {
	record(t, "getmanageddatabases", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("getmanageddatabases"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		services, err := svc.GetManagedDatabases(ctx, &request.GetManagedDatabasesRequest{})
		if !assert.NoError(t, err) {
			return
		}
		assert.Condition(t, func() (success bool) {
			for _, service := range services {
				if service.UUID == details.UUID {
					return true
				}
			}
			return false
		}, "returned slice should contain the created service")
	})
}

func TestService_GetManagedDatabaseLogs(t *testing.T) {
	const (
		batchSize = 5
		waitFor   = 30 * time.Second
	)
	record(t, "getmanageddatabaselogs", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabaselogs")
		serviceDetails, err := svc.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, serviceDetails.UUID)
		require.NoError(t, err)

		if rec.Mode() == recorder.ModeRecording {
			t.Logf("waiting for %s for the logs to be available", waitFor)
			time.Sleep(waitFor)
		}
		for _, order := range []upcloud.ManagedDatabaseLogOrder{
			upcloud.ManagedDatabaseLogOrderAscending,
			upcloud.ManagedDatabaseLogOrderDescending,
		} {
			t.Run("Ctx"+string(order), func(t *testing.T) {
				logReq := &request.GetManagedDatabaseLogsRequest{
					UUID:  serviceDetails.UUID,
					Limit: batchSize,
					Order: order,
				}
				var num int
				var prevLogs *upcloud.ManagedDatabaseLogs
				for {
					logs, err := svc.GetManagedDatabaseLogs(ctx, logReq)
					if !assert.NoError(t, err) {
						return
					}
					if len(logs.Logs) == 0 {
						break
					}
					num += len(logs.Logs)
					switch order {
					case upcloud.ManagedDatabaseLogOrderAscending:
						if prevLogs != nil {
							assert.True(t, prevLogs.Logs[len(prevLogs.Logs)-1].Time.Before(logs.Logs[0].Time),
								"logs are not ascending")
						}
					case upcloud.ManagedDatabaseLogOrderDescending:
						if prevLogs != nil {
							assert.True(t, prevLogs.Logs[len(prevLogs.Logs)-1].Time.After(logs.Logs[0].Time),
								"logs are not descending")
						}
					}
					assert.NotEmpty(t, logs.Logs[0].Time)
					assert.NotEmpty(t, logs.Logs[0].Service)
					assert.NotEmpty(t, logs.Logs[0].Message)
					assert.NotEmpty(t, logs.Logs[0].Hostname)
					logReq.Offset = logs.Offset
					prevLogs = logs
				}
				t.Logf("num logs received in total: %d", num)
				assert.NotZero(t, num)
			})
		}
	})
}

func TestService_GetManagedDatabaseConnections(t *testing.T) {
	record(t, "getmanageddatabaseconnections", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabaseconnections")
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		require.NoError(t, waitForManagedDatabaseRunningState(ctx, rec, svc, serviceDetails.UUID))
		conns, err := svc.GetManagedDatabaseConnections(ctx, &request.GetManagedDatabaseConnectionsRequest{
			UUID:   serviceDetails.UUID,
			Limit:  1000,
			Offset: 0,
		})
		if !assert.NoError(t, err) {
			return
		}
		assert.Len(t, conns, 0)

		err = svc.CancelManagedDatabaseConnection(ctx, &request.CancelManagedDatabaseConnection{
			UUID:      serviceDetails.UUID,
			Pid:       0,
			Terminate: true,
		})
		assert.Error(t, err)
		assert.True(t, strings.HasPrefix(err.(*upcloud.Problem).Title, "Must provide a connection"))

		err = svc.CancelManagedDatabaseConnection(ctx, &request.CancelManagedDatabaseConnection{
			UUID:      serviceDetails.UUID,
			Pid:       0,
			Terminate: false,
		})
		assert.Error(t, err)
		assert.True(t, strings.HasPrefix(err.(*upcloud.Problem).Title, "Must provide a connection"))
	})
}

func TestService_GetManagedDatabaseMetrics(t *testing.T) {
	const (
		timeout = 10 * time.Minute
		waitFor = 2 * time.Minute
	)
	record(t, "getmanageddatabasemetrics", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabasemetrics")
		serviceDetails, err := svc.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		err = waitForManagedDatabaseRunningState(ctx, rec, svc, serviceDetails.UUID)
		require.NoError(t, err)

		if rec.Mode() == recorder.ModeRecording {
			t.Logf("waiting for %s to gather up some data", waitFor)
			time.Sleep(waitFor)
		}
		periods := []upcloud.ManagedDatabaseMetricPeriod{
			upcloud.ManagedDatabaseMetricPeriodHour,
			upcloud.ManagedDatabaseMetricPeriodDay,
			upcloud.ManagedDatabaseMetricPeriodWeek,
			upcloud.ManagedDatabaseMetricPeriodMonth,
			upcloud.ManagedDatabaseMetricPeriodYear,
		}
		for _, period := range periods {
			t.Run("Ctx"+string(period), func(t *testing.T) {
				metrics, err := svc.GetManagedDatabaseMetrics(ctx, &request.GetManagedDatabaseMetricsRequest{
					UUID:   serviceDetails.UUID,
					Period: period,
				})
				assert.NoError(t, err)
				if period == upcloud.ManagedDatabaseMetricPeriodHour {
					validate := func(iv interface{}) {
						switch chart := iv.(type) {
						case upcloud.ManagedDatabaseMetricsChartInt:
							if assert.NotEmpty(t, chart.Rows) {
								assert.IsType(t, 0, chart.Rows[0][0])
							}
							assert.NotEmpty(t, chart.Columns)
							assert.NotEmpty(t, chart.Timestamps)
							assert.NotEmpty(t, chart.Title)
						case upcloud.ManagedDatabaseMetricsChartFloat64:
							if assert.NotEmpty(t, chart.Rows) {
								assert.IsType(t, 0.0, chart.Rows[0][0])
							}
							assert.NotEmpty(t, chart.Columns)
							assert.NotEmpty(t, chart.Timestamps)
							assert.NotEmpty(t, chart.Title)
						default:
							assert.Fail(t, "unexpected metric chart: %+v", chart)
						}
					}
					validate(metrics.CPUUsage)
					validate(metrics.MemoryUsage)
					validate(metrics.LoadAverage)
					validate(metrics.DiskIOReads)
					validate(metrics.DiskIOWrite)
					validate(metrics.NetworkReceive)
					validate(metrics.NetworkReceive)
				}
			})
		}
	})
}

func TestService_GetManagedDatabaseQueryStatisticsMySQL(t *testing.T) {
	record(t, "getmanageddatabasequerystatisticsmysql", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("querystatisticsmysql")
		createReq.Type = upcloud.ManagedDatabaseServiceTypeMySQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()

		require.NoError(t, waitForManagedDatabaseRunningState(ctx, rec, svc, serviceDetails.UUID))

		stats, err := svc.GetManagedDatabaseQueryStatisticsMySQL(ctx, &request.GetManagedDatabaseQueryStatisticsRequest{
			UUID:   serviceDetails.UUID,
			Limit:  1000,
			Offset: 0,
		})
		assert.NoError(t, err)
		assert.Len(t, stats, 0)
	})
}

func TestService_GetManagedDatabaseQueryStatisticsPostgreSQL(t *testing.T) {
	record(t, "getmanageddatabasequerystatisticspostgres", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("querystatisticspostgres")
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()

		require.NoError(t, waitForManagedDatabaseRunningState(ctx, rec, svc, serviceDetails.UUID))

		stats, err := svc.GetManagedDatabaseQueryStatisticsPostgreSQL(ctx, &request.GetManagedDatabaseQueryStatisticsRequest{
			UUID:   serviceDetails.UUID,
			Limit:  1000,
			Offset: 0,
		})
		assert.NoError(t, err)
		assert.Len(t, stats, 1)
		assert.Equal(t, "defaultdb", stats[0].DatabaseName)
	})
}

func TestService_ModifyManagedDatabase(t *testing.T) {
	record(t, "modifymanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("modifymanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		assert.True(t, details.Properties.GetAutoUtilityIPFilter())
		assert.True(t, details.Properties.GetPublicAccess())
		var publicComponentFound bool
		for _, component := range details.Components {
			if component.Route == upcloud.ManagedDatabaseComponentRoutePublic {
				publicComponentFound = true
			}
		}
		assert.True(t, publicComponentFound)
		modifyReq := &request.ModifyManagedDatabaseRequest{UUID: details.UUID}
		modifyReq.Properties.
			SetPublicAccess(false).
			SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		newDetails, err := svc.ModifyManagedDatabase(ctx, modifyReq)
		if !assert.NoError(t, err) {
			return
		}
		assert.False(t, newDetails.Properties.GetPublicAccess())
		assert.Equal(t, []string{upcloud.ManagedDatabaseAllIPv4}, newDetails.Properties.GetIPFilter())
		publicComponentFound = false
		for _, component := range newDetails.Components {
			if component.Route == upcloud.ManagedDatabaseComponentRoutePublic {
				publicComponentFound = true
			}
		}
		assert.False(t, publicComponentFound)
	})
}

func TestService_UpgradeManagedDatabaseVersion(t *testing.T) {
	record(t, "upgrademanageddatabaseversion", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		// This test uses manually created database with postgres version 13
		// This is because upgrading version requires "Started" state; waiting for started state in tests
		// results in huge amount of requests made to verify the state and simply takes too long
		details, err := svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{
			UUID: "09788889-be2d-48da-a527-962b26014b54",
		})

		require.NoError(t, err)
		assert.Equal(t, "13", details.Properties["version"])

		targetVersion := "14"
		updatedDetails, err := svc.UpgradeManagedDatabaseVersion(ctx, &request.UpgradeManagedDatabaseVersionRequest{
			UUID:          details.UUID,
			TargetVersion: targetVersion,
		})

		assert.NoError(t, err)
		assert.Equal(t, targetVersion, updatedDetails.Properties["version"])
	})
}

func TestService_GetManagedDatabaseVersions(t *testing.T) {
	record(t, "getmanageddatabaseversions", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("getmanageddatabaseversions"))
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			require.NoError(t, err)
		}()

		versions, err := svc.GetManagedDatabaseVersions(ctx, &request.GetManagedDatabaseVersionsRequest{
			UUID: details.UUID,
		})

		require.NoError(t, err)
		assert.Len(t, versions, 5)
		assert.Contains(t, versions, "10")
		assert.Contains(t, versions, "11")
		assert.Contains(t, versions, "12")
		assert.Contains(t, versions, "13")
		assert.Contains(t, versions, "14")
	})
}

func TestService_ShutdownStartManagedDatabase(t *testing.T) {
	record(t, "shutdownstartmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("shutdownstartmanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		err = waitForManagedDatabaseInitialBackup(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		shutdownDetails, err := svc.ShutdownManagedDatabase(ctx, &request.ShutdownManagedDatabaseRequest{UUID: details.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.False(t, shutdownDetails.Powered)

		startDetails, err := svc.StartManagedDatabase(ctx, &request.StartManagedDatabaseRequest{UUID: details.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.True(t, startDetails.Powered)
	})
}

func TestService_ManagedDatabaseUserManager(t *testing.T) {
	record(t, "managemanageddatabaseusers", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		serviceDetails, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("managemanageddatabaseusers"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		err = waitForManagedDatabaseRunningState(ctx, rec, svc, serviceDetails.UUID)
		require.NoError(t, err)

		t.Run("CtxCreate", func(t *testing.T) {
			userDetails, err := svc.CreateManagedDatabaseUser(ctx, &request.CreateManagedDatabaseUserRequest{
				ServiceUUID: serviceDetails.UUID,
				Username:    "testuser",
			})
			if !assert.NoError(t, err) {
				return
			}
			t.Run("CtxGet", func(t *testing.T) {
				newUserDetails, err := svc.GetManagedDatabaseUser(ctx, &request.GetManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    userDetails.Username,
				})
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, userDetails, newUserDetails)
			})
			t.Run("CtxList", func(t *testing.T) {
				users, err := svc.GetManagedDatabaseUsers(ctx, &request.GetManagedDatabaseUsersRequest{ServiceUUID: serviceDetails.UUID})
				if !assert.NoError(t, err) {
					return
				}
				if assert.Len(t, users, 2) {
					var primaryFound, normalFound bool
					for _, user := range users {
						switch user.Type {
						case upcloud.ManagedDatabaseUserTypePrimary:
							if assert.Equal(t, serviceDetails.ServiceURIParams.User, user.Username) {
								primaryFound = true
							}
						case upcloud.ManagedDatabaseUserTypeNormal:
							if assert.Equal(t, userDetails.Username, user.Username) {
								normalFound = true
							}
						}
					}
					assert.True(t, primaryFound, "primary user should have been found")
					assert.True(t, normalFound, "normal user should have been found")
				}
			})
			t.Run("CtxModify", func(t *testing.T) {
				//nolint:gosec
				const newPassword = "yXB8gePmxHuESbJx_I-Iag"
				newUserDetails, err := svc.ModifyManagedDatabaseUser(ctx, &request.ModifyManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    userDetails.Username,
					Password:    newPassword,
				})
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, newPassword, newUserDetails.Password)
			})
			t.Run("CtxDelete", func(t *testing.T) {
				err := svc.DeleteManagedDatabaseUser(ctx, &request.DeleteManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    userDetails.Username,
				})
				assert.NoError(t, err)
			})
			t.Run("CtxDeletePrimaryShouldNotSucceed", func(t *testing.T) {
				err := svc.DeleteManagedDatabaseUser(ctx, &request.DeleteManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    serviceDetails.ServiceURIParams.User,
				})
				assert.Error(t, err)
			})
		})
	})
}

func TestService_ManagedDatabaseLogicalDatabaseManager(t *testing.T) {
	const (
		defaultdb = "defaultdb"
	)
	record(t, "managemanageddatabaslogicaldbs", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		serviceDetails, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("managemanageddatabaslogicaldbs"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, serviceDetails.UUID)
		require.NoError(t, err)

		t.Run("CtxCreate", func(t *testing.T) {
			expected := &upcloud.ManagedDatabaseLogicalDatabase{
				Name:      "test",
				LCCollate: "fr_FR.UTF-8",
				LCCType:   "fr_FR.UTF-8",
			}
			dbDetails, err := svc.CreateManagedDatabaseLogicalDatabase(ctx, &request.CreateManagedDatabaseLogicalDatabaseRequest{
				ServiceUUID: serviceDetails.UUID,
				Name:        expected.Name,
				LCCollate:   expected.LCCollate,
				LCCType:     expected.LCCType,
			})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, expected, dbDetails)

			t.Run("CtxList", func(t *testing.T) {
				dbs, err := svc.GetManagedDatabaseLogicalDatabases(ctx,
					&request.GetManagedDatabaseLogicalDatabasesRequest{ServiceUUID: serviceDetails.UUID})
				if !assert.NoError(t, err) {
					return
				}
				if assert.Len(t, dbs, 2) {
					var defaultFound, createdFound bool
					for _, db := range dbs {
						switch db.Name {
						case defaultdb:
							defaultFound = true
						case expected.Name:
							assert.Equal(t, *expected, db)
							createdFound = true
						}
					}
					assert.True(t, defaultFound, "default logical database was not found")
					assert.True(t, createdFound, "created logical database was not found")
				}
			})
			t.Run("CtxDelete", func(t *testing.T) {
				err := svc.DeleteManagedDatabaseLogicalDatabase(ctx, &request.DeleteManagedDatabaseLogicalDatabaseRequest{
					ServiceUUID: serviceDetails.UUID,
					Name:        expected.Name,
				})
				assert.NoError(t, err)
			})
		})
	})
}

func TestService_GetManagedDatabaseServiceType(t *testing.T) {
	record(t, "getmanageddatabaseservicetype", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		databaseTypes := []string{"pg", "mysql"}

		for _, databaseType := range databaseTypes {
			serviceType, err := svc.GetManagedDatabaseServiceType(ctx, &request.GetManagedDatabaseServiceTypeRequest{Type: databaseType})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, databaseType, serviceType.Name)
		}
	})
}

func TestService_GetManagedDatabaseServiceTypes(t *testing.T) {
	record(t, "getmanageddatabaseservicetypes", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		types, err := svc.GetManagedDatabaseServiceTypes(ctx, &request.GetManagedDatabaseServiceTypesRequest{})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "pg", types["pg"].Name)
		assert.Equal(t, "mysql", types["mysql"].Name)
	})
}

func TestService_ModifyManagedDatabaseUserAccessControl(t *testing.T) {
	record(t, "modifymanageddatabaseuseraccesscontrol", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		db, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("modifyuseraccesscontrol"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			if err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: db.UUID}); err != nil {
				t.Log(err)
			}
		}()
		if !assert.NoError(t, waitForManagedDatabaseRunningState(ctx, rec, svc, db.UUID)) {
			return
		}
		user, err := svc.CreateManagedDatabaseUser(ctx, &request.CreateManagedDatabaseUserRequest{
			ServiceUUID: db.UUID,
			Username:    "demouser",
			PGAccessControl: &upcloud.ManagedDatabaseUserPGAccessControl{
				AllowReplication: true,
			},
		})
		if !assert.NoError(t, err) {
			return
		}
		assert.True(t, user.PGAccessControl.AllowReplication)

		user, err = svc.ModifyManagedDatabaseUserAccessControl(ctx, &request.ModifyManagedDatabaseUserAccessControlRequest{
			ServiceUUID: db.UUID,
			Username:    user.Username,
			PGAccessControl: &upcloud.ManagedDatabaseUserPGAccessControl{
				AllowReplication: false,
			},
		})
		if !assert.NoError(t, err) {
			return
		}
		assert.False(t, user.PGAccessControl.AllowReplication)
	})
}

func waitForManagedDatabaseInitialBackup(ctx context.Context, rec *recorder.Recorder, svc *Service, dbUUID string) error {
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
		waitForBackupDetails, err := svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: dbUUID})
		if err != nil {
			return err
		}
		if len(waitForBackupDetails.Backups) > 0 {
			break
		}
		if time.Now().After(waitUntil) {
			return fmt.Errorf("timeout %s reached", timeout.String())
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}

func waitForManagedDatabaseRunningState(ctx context.Context, rec *recorder.Recorder, svc *Service, dbUUID string) error {
	if rec.Mode() != recorder.ModeRecording {
		return nil
	}

	rec.AddPassthrough(func(h *http.Request) bool {
		return true
	})
	defer func() {
		rec.Passthroughs = nil
	}()

	_, err := svc.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
		UUID:         dbUUID,
		DesiredState: upcloud.ManagedDatabaseStateRunning,
		Timeout:      waitTimeout,
	})

	return err
}

const (
	managedDatabaseTestPlan = "2x2xCPU-4GB-100GB"
	managedDatabaseTestZone = "fi-hel2"
)

func getTestCreateRequest(name string) *request.CreateManagedDatabaseRequest {
	r := request.CreateManagedDatabaseRequest{
		HostNamePrefix: name,
		Maintenance: request.ManagedDatabaseMaintenanceTimeRequest{
			DayOfWeek: "monday",
			Time:      "12:00:00",
		},
		Plan: managedDatabaseTestPlan,
		Properties: request.ManagedDatabasePropertiesRequest{
			upcloud.ManagedDatabasePropertyAutoUtilityIPFilter: true,
			upcloud.ManagedDatabasePropertyIPFilter:            []string{"10.0.0.1/32"},
			upcloud.ManagedDatabasePropertyPublicAccess:        true,
		},
		Title: "test",
		Type:  upcloud.ManagedDatabaseServiceTypePostgreSQL,
		Zone:  managedDatabaseTestZone,
	}
	return &r
}
