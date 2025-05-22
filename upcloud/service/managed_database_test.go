package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

const (
	managedDatabaseTestPlanMySQL      = "2x2xCPU-4GB-100GB"
	managedDatabaseTestPlanOpenSearch = "1x2xCPU-4GB-80GB-1D"
	managedDatabaseTestPlanPostgreSQL = "2x2xCPU-4GB-100GB"
	managedDatabaseTestPlanRedis      = "1x1xCPU-2GB"
	managedDatabaseTestPlanValkey     = "1x1xCPU-2GB"
	managedDatabaseTestZone           = "fi-hel2"
)

func TestService_CloneManagedDatabase(t *testing.T) {
	record(t, "clonemanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		var cloneDetails *upcloud.ManagedDatabase
		createReq := getTestCreateRequest("clonemanageddatabase", upcloud.ManagedDatabaseServiceTypePostgreSQL)
		details, err := svc.CreateManagedDatabase(ctx, createReq)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
			if cloneDetails != nil {
				t.Logf("deleting clone %s", cloneDetails.UUID)
				err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: cloneDetails.UUID})
				assert.NoError(t, err)
			}
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		err = waitForManagedDatabaseInitialBackup(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
		require.NoError(t, err)

		cloneReq := &request.CloneManagedDatabaseRequest{
			UUID:           details.UUID,
			CloneTime:      details.Backups[0].BackupTime.Add(1 * time.Second),
			HostNamePrefix: fmt.Sprintf("%s-clone", details.Name),
			Title:          fmt.Sprintf("%s-clone", details.Title),
			Zone:           createReq.Zone,
			Plan:           createReq.Plan,
		}
		cloneDetails, err = svc.CloneManagedDatabase(ctx, cloneReq)
		require.NoError(t, err)
	})
}

func TestService_CreateManagedDatabase(t *testing.T) {
	typesToTest := []upcloud.ManagedDatabaseServiceType{
		upcloud.ManagedDatabaseServiceTypeMySQL,
		upcloud.ManagedDatabaseServiceTypePostgreSQL,
		upcloud.ManagedDatabaseServiceTypeRedis, //nolint:staticcheck // To be removed when Redis support has been removed
		upcloud.ManagedDatabaseServiceTypeValkey,
		upcloud.ManagedDatabaseServiceTypeOpenSearch,
	}

	record(t, "createmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		for i, serviceType := range typesToTest {
			t.Run(string(serviceType), func(t *testing.T) {
				req := getTestCreateRequest(fmt.Sprintf("createmanageddatabase%d", i), serviceType)
				details, err := svc.CreateManagedDatabase(ctx, req)
				require.NoError(t, err)

				defer func() {
					t.Logf("deleting %s", details.UUID)
					err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
					assert.NoError(t, err)
				}()

				err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
				require.NoError(t, err)

				details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
				require.NoError(t, err)
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
				case upcloud.ManagedDatabaseServiceTypeMySQL:
					assert.NotEmpty(t, details.Metadata.MySQLVersion)
				case upcloud.ManagedDatabaseServiceTypeOpenSearch:
					assert.NotEmpty(t, details.Metadata.OpenSearchVersion)
				case upcloud.ManagedDatabaseServiceTypePostgreSQL:
					assert.NotEmpty(t, details.Metadata.PGVersion)
				case upcloud.ManagedDatabaseServiceTypeRedis: //nolint:staticcheck // To be removed when Redis support has been removed
					assert.NotEmpty(t, details.Metadata.RedisVersion) //nolint:staticcheck // To be removed when Redis support has been removed
				case upcloud.ManagedDatabaseServiceTypeValkey:
					assert.NotEmpty(t, details.Metadata.ValkeyVersion)
				}
			})
		}
	})
}

func TestService_WaitForManagedDatabaseState(t *testing.T) {
	record(t, "waitformanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("waitformanageddatabase", upcloud.ManagedDatabaseServiceTypePostgreSQL))
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{
			UUID: details.UUID,
		})
		require.NoError(t, err)
		assert.Equal(t, upcloud.ManagedDatabaseStateRunning, details.State)
	})
}

func TestService_GetManagedDatabase(t *testing.T) {
	record(t, "getmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		req := getTestCreateRequest("getmanageddatabase", upcloud.ManagedDatabaseServiceTypePostgreSQL)
		details, err := svc.CreateManagedDatabase(ctx, req)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
		require.NoError(t, err)
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
		uuidMap := make(map[string]bool)
		for i := 0; i < 4; i++ {
			details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest(fmt.Sprintf("getmanageddatabases-%d", i), upcloud.ManagedDatabaseServiceTypePostgreSQL))
			require.NoError(t, err)

			uuidMap[details.UUID] = true
		}

		defer func() {
			for uuid := range uuidMap {
				t.Logf("deleting %s", uuid)
				err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: uuid})
				assert.NoError(t, err)
			}
		}()

		services, err := svc.GetManagedDatabases(ctx, &request.GetManagedDatabasesRequest{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(services), 4)

		services, err = svc.GetManagedDatabases(ctx, &request.GetManagedDatabasesRequest{Page: &request.Page{Size: 2, Number: 1}})
		require.NoError(t, err)
		assert.Equal(t, len(services), 2)
	})
}

func TestService_GetManagedDatabaseLogs(t *testing.T) {
	const (
		batchSize = 5
		waitFor   = 30 * time.Second
	)
	record(t, "getmanageddatabaselogs", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabaselogs", upcloud.ManagedDatabaseServiceTypePostgreSQL)
		details, err := svc.CreateManagedDatabase(ctx, createReq)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		if rec.IsRecording() {
			t.Logf("waiting for %s for the logs to be available", waitFor)
			time.Sleep(waitFor)
		}
		for _, order := range []upcloud.ManagedDatabaseLogOrder{
			upcloud.ManagedDatabaseLogOrderAscending,
			upcloud.ManagedDatabaseLogOrderDescending,
		} {
			t.Run("Ctx"+string(order), func(t *testing.T) {
				logReq := &request.GetManagedDatabaseLogsRequest{
					UUID:  details.UUID,
					Limit: batchSize,
					Order: order,
				}
				var num int
				var prevLogs *upcloud.ManagedDatabaseLogs
				for {
					logs, err := svc.GetManagedDatabaseLogs(ctx, logReq)
					require.NoError(t, err)
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

func TestService_GetManagedDatabaseSessions(t *testing.T) {
	record(t, "getmanageddatabasesessions", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabasesessions", upcloud.ManagedDatabaseServiceTypePostgreSQL)
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		details, err := svc.CreateManagedDatabase(ctx, createReq)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		sessions, err := svc.GetManagedDatabaseSessions(ctx, &request.GetManagedDatabaseSessionsRequest{
			UUID:   details.UUID,
			Limit:  1000,
			Offset: 0,
			Order:  "pid:desc",
		})
		require.NoError(t, err)
		assert.Len(t, sessions.MySQL, 0)
		assert.Len(t, sessions.PostgreSQL, 0)
		assert.Len(t, sessions.Redis, 0) //nolint:staticcheck // To be removed when Redis support has been removed
		assert.Len(t, sessions.Valkey, 0)

		err = svc.CancelManagedDatabaseSession(ctx, &request.CancelManagedDatabaseSession{
			UUID:      details.UUID,
			Pid:       0,
			Terminate: true,
		})
		assert.Error(t, err)

		err = svc.CancelManagedDatabaseSession(ctx, &request.CancelManagedDatabaseSession{
			UUID:      details.UUID,
			Pid:       0,
			Terminate: false,
		})
		assert.Error(t, err)
	})
}

func TestService_GetManagedDatabaseMetrics(t *testing.T) {
	const waitFor = 2 * time.Minute

	record(t, "getmanageddatabasemetrics", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabasemetrics", upcloud.ManagedDatabaseServiceTypePostgreSQL)
		details, err := svc.CreateManagedDatabase(ctx, createReq)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		if rec.IsRecording() {
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
					UUID:   details.UUID,
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
		createReq := getTestCreateRequest("querystatisticsmysql", upcloud.ManagedDatabaseServiceTypeMySQL)
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		details, err := svc.CreateManagedDatabase(ctx, createReq)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		stats, err := svc.GetManagedDatabaseQueryStatisticsMySQL(ctx, &request.GetManagedDatabaseQueryStatisticsRequest{
			UUID:   details.UUID,
			Limit:  1000,
			Offset: 0,
		})
		assert.NoError(t, err)
		assert.Len(t, stats, 0)
	})
}

func TestService_GetManagedDatabaseQueryStatisticsPostgreSQL(t *testing.T) {
	record(t, "getmanageddatabasequerystatisticspostgres", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("querystatisticspostgres", upcloud.ManagedDatabaseServiceTypePostgreSQL)
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		details, err := svc.CreateManagedDatabase(ctx, createReq)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		require.NoError(t, waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID))

		stats, err := svc.GetManagedDatabaseQueryStatisticsPostgreSQL(ctx, &request.GetManagedDatabaseQueryStatisticsRequest{
			UUID:   details.UUID,
			Limit:  1000,
			Offset: 0,
		})
		assert.NoError(t, err)
		assert.Len(t, stats, 2)
		assert.Equal(t, "defaultdb", stats[0].DatabaseName)
	})
}

func TestService_ModifyManagedDatabase(t *testing.T) {
	record(t, "modifymanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("modifymanageddatabase", upcloud.ManagedDatabaseServiceTypePostgreSQL))
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
		require.NoError(t, err)
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
		details, err = svc.ModifyManagedDatabase(ctx, modifyReq)
		require.NoError(t, err)
		assert.False(t, details.Properties.GetPublicAccess())
		assert.Equal(t, []string{upcloud.ManagedDatabaseAllIPv4}, details.Properties.GetIPFilter())

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
		require.NoError(t, err)

		publicComponentFound = false
		for _, component := range details.Components {
			if component.Route == upcloud.ManagedDatabaseComponentRoutePublic {
				publicComponentFound = true
			}
		}
		assert.False(t, publicComponentFound)
	})
}

func TestService_UpgradeManagedDatabaseVersion(t *testing.T) {
	record(t, "upgrademanageddatabaseversion", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		originalVersion := "14"
		req := getTestCreateRequest("upgrademanageddatabase", upcloud.ManagedDatabaseServiceTypePostgreSQL)
		req.Properties.SetString("version", originalVersion)
		details, err := svc.CreateManagedDatabase(ctx, req)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		version, err := details.Properties.GetString("version")
		require.NoError(t, err)
		assert.Equal(t, originalVersion, version)

		targetVersion := "15"
		details, err = svc.UpgradeManagedDatabaseVersion(ctx, &request.UpgradeManagedDatabaseVersionRequest{
			UUID:          details.UUID,
			TargetVersion: targetVersion,
		})
		assert.NoError(t, err)

		version, err = details.Properties.GetString("version")
		require.NoError(t, err)
		assert.Equal(t, targetVersion, version)
	})
}

func TestService_GetManagedDatabaseVersions(t *testing.T) {
	record(t, "getmanageddatabaseversions", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("getmanageddatabaseversions", upcloud.ManagedDatabaseServiceTypePostgreSQL))
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			require.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		versions, err := svc.GetManagedDatabaseVersions(ctx, &request.GetManagedDatabaseVersionsRequest{
			UUID: details.UUID,
		})
		require.NoError(t, err)
		assert.Len(t, versions, 4)
		assert.Contains(t, versions, "12")
		assert.Contains(t, versions, "13")
		assert.Contains(t, versions, "14")
		assert.Contains(t, versions, "15")
	})
}

func TestService_ShutdownStartManagedDatabase(t *testing.T) {
	record(t, "shutdownstartmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		req := getTestCreateRequest("shutdownstartmanageddatabase", upcloud.ManagedDatabaseServiceTypePostgreSQL)
		req.TerminationProtection = upcloud.BoolPtr(true)
		details, err := svc.CreateManagedDatabase(ctx, req)
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		err = waitForManagedDatabaseInitialBackup(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		// Shutdown should fail because termination protection is enabled.
		_, err = svc.ShutdownManagedDatabase(ctx, &request.ShutdownManagedDatabaseRequest{UUID: details.UUID})
		require.Error(t, err)

		details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
		require.NoError(t, err)
		assert.True(t, details.Powered)
		assert.True(t, details.TerminationProtection)

		_, err = svc.ModifyManagedDatabase(ctx, &request.ModifyManagedDatabaseRequest{
			UUID:                  details.UUID,
			TerminationProtection: upcloud.BoolPtr(false),
		})
		require.NoError(t, err)

		shutdownDetails, err := svc.ShutdownManagedDatabase(ctx, &request.ShutdownManagedDatabaseRequest{UUID: details.UUID})
		require.NoError(t, err)
		assert.False(t, shutdownDetails.Powered)

		startDetails, err := svc.StartManagedDatabase(ctx, &request.StartManagedDatabaseRequest{UUID: details.UUID})
		require.NoError(t, err)
		assert.True(t, startDetails.Powered)
	})
}

func TestService_ManagedDatabaseUserManager(t *testing.T) {
	record(t, "managemanageddatabaseusers", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("managemanageddatabaseusers", upcloud.ManagedDatabaseServiceTypePostgreSQL))
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		details, err = svc.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
		require.NoError(t, err)

		t.Run("CtxCreate", func(t *testing.T) {
			userDetails, err := svc.CreateManagedDatabaseUser(ctx, &request.CreateManagedDatabaseUserRequest{
				ServiceUUID: details.UUID,
				Username:    "testuser",
			})
			require.NoError(t, err)

			t.Run("CtxGet", func(t *testing.T) {
				_, err := svc.GetManagedDatabaseUser(ctx, &request.GetManagedDatabaseUserRequest{
					ServiceUUID: details.UUID,
					Username:    userDetails.Username,
				})
				require.NoError(t, err)
			})

			t.Run("CtxList", func(t *testing.T) {
				users, err := svc.GetManagedDatabaseUsers(ctx, &request.GetManagedDatabaseUsersRequest{ServiceUUID: details.UUID})
				require.NoError(t, err)

				if assert.Len(t, users, 2) {
					var primaryFound, normalFound bool
					for _, user := range users {
						switch user.Type {
						case upcloud.ManagedDatabaseUserTypePrimary:
							if assert.Equal(t, details.ServiceURIParams.User, user.Username) {
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
					ServiceUUID: details.UUID,
					Username:    userDetails.Username,
					Password:    newPassword,
				})
				require.NoError(t, err)
				assert.Equal(t, newPassword, newUserDetails.Password)
			})

			t.Run("CtxDelete", func(t *testing.T) {
				err := svc.DeleteManagedDatabaseUser(ctx, &request.DeleteManagedDatabaseUserRequest{
					ServiceUUID: details.UUID,
					Username:    userDetails.Username,
				})
				assert.NoError(t, err)
			})

			t.Run("CtxDeletePrimaryShouldNotSucceed", func(t *testing.T) {
				err := svc.DeleteManagedDatabaseUser(ctx, &request.DeleteManagedDatabaseUserRequest{
					ServiceUUID: details.UUID,
					Username:    details.ServiceURIParams.User,
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
	record(t, "managemanageddatabaselogicaldbs", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("managemanageddatabaslogicaldbs", upcloud.ManagedDatabaseServiceTypePostgreSQL))
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, details.UUID)
		require.NoError(t, err)

		t.Run("CtxCreate", func(t *testing.T) {
			expected := &upcloud.ManagedDatabaseLogicalDatabase{
				Name:      "test",
				LCCollate: "fr_FR.UTF-8",
				LCCType:   "fr_FR.UTF-8",
			}
			dbDetails, err := svc.CreateManagedDatabaseLogicalDatabase(ctx, &request.CreateManagedDatabaseLogicalDatabaseRequest{
				ServiceUUID: details.UUID,
				Name:        expected.Name,
				LCCollate:   expected.LCCollate,
				LCCType:     expected.LCCType,
			})
			require.NoError(t, err)
			assert.Equal(t, expected, dbDetails)

			t.Run("CtxList", func(t *testing.T) {
				dbs, err := svc.GetManagedDatabaseLogicalDatabases(ctx,
					&request.GetManagedDatabaseLogicalDatabasesRequest{ServiceUUID: details.UUID})
				require.NoError(t, err)

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
					ServiceUUID: details.UUID,
					Name:        expected.Name,
				})
				assert.NoError(t, err)
			})
		})
	})
}

func TestService_GetManagedDatabaseServiceType(t *testing.T) {
	record(t, "getmanageddatabaseservicetype", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		databaseTypes := []upcloud.ManagedDatabaseServiceType{
			upcloud.ManagedDatabaseServiceTypeMySQL,
			upcloud.ManagedDatabaseServiceTypeOpenSearch,
			upcloud.ManagedDatabaseServiceTypePostgreSQL,
			upcloud.ManagedDatabaseServiceTypeRedis, //nolint:staticcheck // To be removed when Redis support has been removed
			upcloud.ManagedDatabaseServiceTypeValkey,
		}
		for _, databaseType := range databaseTypes {
			serviceType, err := svc.GetManagedDatabaseServiceType(ctx, &request.GetManagedDatabaseServiceTypeRequest{Type: string(databaseType)})
			assert.NoError(t, err)
			assert.Equal(t, string(databaseType), serviceType.Name)
		}
	})
}

func TestService_GetManagedDatabaseServiceTypes(t *testing.T) {
	record(t, "getmanageddatabaseservicetypes", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		types, err := svc.GetManagedDatabaseServiceTypes(ctx, &request.GetManagedDatabaseServiceTypesRequest{})
		require.NoError(t, err)
		assert.Equal(t, string(upcloud.ManagedDatabaseServiceTypeMySQL), types["mysql"].Name)
		assert.Equal(t, string(upcloud.ManagedDatabaseServiceTypeOpenSearch), types["opensearch"].Name)
		assert.Equal(t, string(upcloud.ManagedDatabaseServiceTypePostgreSQL), types["pg"].Name)
		assert.Equal(t, string(upcloud.ManagedDatabaseServiceTypeRedis), types["redis"].Name) //nolint:staticcheck // To be removed when Redis support has been removed
		assert.Equal(t, string(upcloud.ManagedDatabaseServiceTypeValkey), types["valkey"].Name)
	})
}

func TestService_ModifyManagedDatabaseUserPostgreSQLAccessControl(t *testing.T) {
	record(t, "modifymanageddatabaseuserpostgreqlsaccesscontrol", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		db, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("modifyuseraccesscontrol", upcloud.ManagedDatabaseServiceTypePostgreSQL))
		require.NoError(t, err)

		defer func() {
			if err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: db.UUID}); err != nil {
				t.Log(err)
			}
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, db.UUID)
		require.NoError(t, err)

		user, err := svc.CreateManagedDatabaseUser(ctx, &request.CreateManagedDatabaseUserRequest{
			ServiceUUID: db.UUID,
			Username:    "demouser",
			PGAccessControl: &upcloud.ManagedDatabaseUserPGAccessControl{
				AllowReplication: upcloud.BoolPtr(true),
			},
		})
		require.NoError(t, err)
		assert.True(t, *user.PGAccessControl.AllowReplication)

		user, err = svc.ModifyManagedDatabaseUserAccessControl(ctx, &request.ModifyManagedDatabaseUserAccessControlRequest{
			ServiceUUID: db.UUID,
			Username:    user.Username,
			PGAccessControl: &upcloud.ManagedDatabaseUserPGAccessControl{
				AllowReplication: upcloud.BoolPtr(false),
			},
		})
		require.NoError(t, err)
		assert.False(t, *user.PGAccessControl.AllowReplication)
	})
}

func TestService_ModifyManagedDatabaseUserOpenSearchAccessControl(t *testing.T) {
	record(t, "modifymanageddatabaseuseropensearchaccesscontrol", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		db, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("modifyuseraccesscontrolos", upcloud.ManagedDatabaseServiceTypeOpenSearch))
		require.NoError(t, err)

		defer func() {
			if err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: db.UUID}); err != nil {
				t.Log(err)
			}
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, db.UUID)
		require.NoError(t, err)

		user, err := svc.CreateManagedDatabaseUser(ctx, &request.CreateManagedDatabaseUserRequest{
			ServiceUUID: db.UUID,
			Username:    "demouser",
			OpenSearchAccessControl: &upcloud.ManagedDatabaseUserOpenSearchAccessControl{
				Rules: &[]upcloud.ManagedDatabaseUserOpenSearchAccessControlRule{
					{
						Index:      "index_1",
						Permission: upcloud.ManagedDatabaseUserOpenSearchAccessControlRulePermissionReadWrite,
					},
				},
			},
		})
		require.NoError(t, err)

		rules := *user.OpenSearchAccessControl.Rules
		assert.Equal(t, "index_1", rules[0].Index)
		assert.Equal(t, "readwrite", string(rules[0].Permission))

		user, err = svc.ModifyManagedDatabaseUserAccessControl(ctx, &request.ModifyManagedDatabaseUserAccessControlRequest{
			ServiceUUID: db.UUID,
			Username:    user.Username,
			OpenSearchAccessControl: &upcloud.ManagedDatabaseUserOpenSearchAccessControl{
				Rules: &[]upcloud.ManagedDatabaseUserOpenSearchAccessControlRule{
					{
						Index:      "index_1",
						Permission: upcloud.ManagedDatabaseUserOpenSearchAccessControlRulePermissionRead,
					},
				},
			},
		})
		require.NoError(t, err)

		rules = *user.OpenSearchAccessControl.Rules
		assert.Equal(t, "index_1", rules[0].Index)
		assert.Equal(t, "read", string(rules[0].Permission))

		user, err = svc.ModifyManagedDatabaseUserAccessControl(ctx, &request.ModifyManagedDatabaseUserAccessControlRequest{
			ServiceUUID: db.UUID,
			Username:    user.Username,
			OpenSearchAccessControl: &upcloud.ManagedDatabaseUserOpenSearchAccessControl{
				Rules: &[]upcloud.ManagedDatabaseUserOpenSearchAccessControlRule{},
			},
		})
		require.NoError(t, err)
		assert.Len(t, *user.OpenSearchAccessControl.Rules, 0)
	})
}

func TestService_ModifyManagedDatabaseUserValkeyAccessControl(t *testing.T) {
	record(t, "modifymanageddatabaseuservalkeyaccesscontrol", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		db, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("modifyuseraccesscontrolvalkey", upcloud.ManagedDatabaseServiceTypeValkey))
		require.NoError(t, err)

		defer func() {
			if err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: db.UUID}); err != nil {
				t.Log(err)
			}
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, db.UUID)
		require.NoError(t, err)

		user, err := svc.CreateManagedDatabaseUser(ctx, &request.CreateManagedDatabaseUserRequest{
			ServiceUUID: db.UUID,
			Username:    "demouser",
			ValkeyAccessControl: &upcloud.ManagedDatabaseUserValkeyAccessControl{
				Categories: &[]string{"+@set"},
				Channels:   &[]string{"*"},
				Commands:   &[]string{"+set"},
				Keys:       &[]string{"key_*"},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, []string{"+@set"}, *user.ValkeyAccessControl.Categories)
		assert.Equal(t, []string{"*"}, *user.ValkeyAccessControl.Channels)
		assert.Equal(t, []string{"+set"}, *user.ValkeyAccessControl.Commands)
		assert.Equal(t, []string{"key_*"}, *user.ValkeyAccessControl.Keys)

		user, err = svc.ModifyManagedDatabaseUserAccessControl(ctx, &request.ModifyManagedDatabaseUserAccessControlRequest{
			ServiceUUID: db.UUID,
			Username:    user.Username,
			ValkeyAccessControl: &upcloud.ManagedDatabaseUserValkeyAccessControl{
				Categories: &[]string{},
				Channels:   &[]string{},
				Commands:   &[]string{},
				Keys:       &[]string{"key_*"},
			},
		})
		require.NoError(t, err)
		assert.Equal(t, []string{}, *user.ValkeyAccessControl.Categories)
		assert.Equal(t, []string{}, *user.ValkeyAccessControl.Channels)
		assert.Equal(t, []string{}, *user.ValkeyAccessControl.Commands)
		assert.Equal(t, []string{"key_*"}, *user.ValkeyAccessControl.Keys)
	})
}

func TestService_ModifyManagedDatabaseAccessControl(t *testing.T) {
	record(t, "modifymanageddatabaseaccesscontrol", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		db, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("modifyaccesscontrol", upcloud.ManagedDatabaseServiceTypeOpenSearch))
		require.NoError(t, err)

		defer func() {
			if err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: db.UUID}); err != nil {
				t.Log(err)
			}
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, db.UUID)
		require.NoError(t, err)

		ac, err := svc.GetManagedDatabaseAccessControl(ctx, &request.GetManagedDatabaseAccessControlRequest{ServiceUUID: db.UUID})
		require.NoError(t, err)
		assert.False(t, *ac.ACLsEnabled)
		assert.False(t, *ac.ExtendedACLsEnabled)

		ac, err = svc.ModifyManagedDatabaseAccessControl(ctx, &request.ModifyManagedDatabaseAccessControlRequest{
			ServiceUUID: db.UUID,
			ACLsEnabled: upcloud.BoolPtr(true),
		})
		require.NoError(t, err)
		assert.True(t, *ac.ACLsEnabled)
		assert.False(t, *ac.ExtendedACLsEnabled)

		ac, err = svc.ModifyManagedDatabaseAccessControl(ctx, &request.ModifyManagedDatabaseAccessControlRequest{
			ServiceUUID:         db.UUID,
			ExtendedACLsEnabled: upcloud.BoolPtr(true),
		})
		require.NoError(t, err)
		assert.True(t, *ac.ACLsEnabled)
		assert.True(t, *ac.ExtendedACLsEnabled)
	})
}

func TestService_GetManagedDatabaseIndices(t *testing.T) {
	record(t, "getmanageddatabaseincices", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		db, err := svc.CreateManagedDatabase(ctx, getTestCreateRequest("getindices", upcloud.ManagedDatabaseServiceTypeOpenSearch))
		require.NoError(t, err)

		defer func() {
			if err := svc.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: db.UUID}); err != nil {
				t.Log(err)
			}
		}()

		err = waitForManagedDatabaseRunningState(ctx, rec, svc, db.UUID)
		require.NoError(t, err)

		indices, err := svc.GetManagedDatabaseIndices(ctx, &request.GetManagedDatabaseIndicesRequest{ServiceUUID: db.UUID})
		require.NoError(t, err)
		assert.NotEmpty(t, indices)
	})
}

func waitForManagedDatabaseInitialBackup(ctx context.Context, rec *recorder.Recorder, svc *Service, dbUUID string) error {
	if !rec.IsRecording() {
		return nil
	}

	const timeout = 10 * time.Minute

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
	if !rec.IsRecording() {
		return nil
	}

	_, err := svc.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
		UUID:         dbUUID,
		DesiredState: upcloud.ManagedDatabaseStateRunning,
	})

	return err
}

func getTestCreateRequest(name string, serviceType upcloud.ManagedDatabaseServiceType) *request.CreateManagedDatabaseRequest {
	var plan string
	switch serviceType {
	case upcloud.ManagedDatabaseServiceTypeMySQL:
		plan = managedDatabaseTestPlanMySQL
	case upcloud.ManagedDatabaseServiceTypeOpenSearch:
		plan = managedDatabaseTestPlanOpenSearch
	case upcloud.ManagedDatabaseServiceTypePostgreSQL:
		plan = managedDatabaseTestPlanPostgreSQL
	case upcloud.ManagedDatabaseServiceTypeRedis: //nolint:staticcheck // To be removed when Redis support has been removed
		plan = managedDatabaseTestPlanRedis
	case upcloud.ManagedDatabaseServiceTypeValkey:
		plan = managedDatabaseTestPlanValkey
	default:
		return nil
	}

	r := request.CreateManagedDatabaseRequest{
		HostNamePrefix: name,
		Maintenance: request.ManagedDatabaseMaintenanceTimeRequest{
			DayOfWeek: "monday",
			Time:      "12:00:00",
		},
		Plan: plan,
		Properties: request.ManagedDatabasePropertiesRequest{
			upcloud.ManagedDatabasePropertyAutoUtilityIPFilter: true,
			upcloud.ManagedDatabasePropertyIPFilter:            []string{"10.0.0.1/32"},
			upcloud.ManagedDatabasePropertyPublicAccess:        true,
		},
		Title: name,
		Type:  serviceType,
		Zone:  managedDatabaseTestZone,
	}
	return &r
}
