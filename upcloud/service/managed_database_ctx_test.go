package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"

	_ "github.com/go-sql-driver/mysql"
)

func TestService_CloneManagedDatabaseContext(t *testing.T) {
	const timeout = 10 * time.Minute
	recordWithContext(t, "clonemanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		var cloneDetails *upcloud.ManagedDatabase
		createReq := getTestCreateRequest("clonemanageddatabase")
		serviceDetails, err := svcContext.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
			if cloneDetails != nil {
				t.Logf("deleting clone %s", cloneDetails.UUID)
				err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: cloneDetails.UUID})
				assert.NoError(t, err)
			}
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         serviceDetails.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      timeout,
			})
			if !assert.NoError(t, err) {
				return
			}
			t.Logf("waiting for initial service service backup (up to %s)", timeout)
			waitUntil := time.Now().Add(timeout)
			for {
				waitForBackupDetails, err := svcContext.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: serviceDetails.UUID})
				if !assert.NoError(t, err) {
					return
				}
				if len(waitForBackupDetails.Backups) > 0 {
					break
				}
				if time.Now().After(waitUntil) {
					assert.Fail(t, "timed out after waiting for initial backup")
					return
				}
				time.Sleep(5 * time.Second)
			}
			rec.Passthroughs = nil
		}
		serviceDetails, err = svcContext.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: serviceDetails.UUID})
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
		cloneDetails, err = svcContext.CloneManagedDatabase(ctx, cloneReq)
		if !assert.NoError(t, err) {
			return
		}
	})
}

func TestService_CreateManagedDatabaseContext(t *testing.T) {
	typesToTest := []upcloud.ManagedDatabaseServiceType{
		upcloud.ManagedDatabaseServiceTypeMySQL,
		upcloud.ManagedDatabaseServiceTypePostgreSQL,
	}
	recordWithContext(t, "createmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		for _, serviceType := range typesToTest {
			t.Run(string(serviceType), func(t *testing.T) {
				req := getTestCreateRequest("createmanageddatabase")
				req.Type = serviceType
				details, err := svcContext.CreateManagedDatabase(ctx, req)
				if !assert.NoError(t, err) {
					return
				}
				defer func() {
					t.Logf("deleting %s", details.UUID)
					err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
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
			})
		}
	})
}

func TestService_WaitForManagedDatabaseStateContext(t *testing.T) {
	const timeout = 10 * time.Minute
	recordWithContext(t, "waitformanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		details, err := svcContext.CreateManagedDatabase(ctx, getTestCreateRequest("waitformanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		var newDetails *upcloud.ManagedDatabase
		if rec.Mode() == recorder.ModeRecording {
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         details.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      timeout,
			})
			if !assert.NoError(t, err) {
				return
			}
			rec.Passthroughs = nil
		}
		newDetails, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
			UUID:         details.UUID,
			DesiredState: upcloud.ManagedDatabaseStateRunning,
			Timeout:      15 * time.Second,
		})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, upcloud.ManagedDatabaseStateRunning, newDetails.State)
	})
}

func TestService_GetManagedDatabaseContext(t *testing.T) {
	recordWithContext(t, "getmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		req := getTestCreateRequest("getmanageddatabase")
		details, err := svcContext.CreateManagedDatabase(ctx, req)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		details, err = svcContext.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
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
func TestService_GetManagedDatabasesContext(t *testing.T) {
	recordWithContext(t, "getmanageddatabases", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		details, err := svcContext.CreateManagedDatabase(ctx, getTestCreateRequest("getmanageddatabases"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		services, err := svcContext.GetManagedDatabases(ctx, &request.GetManagedDatabasesRequest{})
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

func TestService_GetManagedDatabaseLogsContext(t *testing.T) {
	const (
		batchSize = 5
		timeout   = 10 * time.Minute
		waitFor   = 30 * time.Second
	)
	recordWithContext(t, "getmanageddatabaselogs", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		createReq := getTestCreateRequest("getmanageddatabaselogs")
		serviceDetails, err := svcContext.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         serviceDetails.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      timeout,
			})
			if !assert.NoError(t, err) {
				return
			}
			rec.Passthroughs = nil
			t.Logf("waiting for %s for the logs to be available", waitFor)
			time.Sleep(waitFor)
		}
		for _, order := range []upcloud.ManagedDatabaseLogOrder{
			upcloud.ManagedDatabaseLogOrderAscending,
			upcloud.ManagedDatabaseLogOrderDescending,
		} {
			t.Run(string(order), func(t *testing.T) {
				logReq := &request.GetManagedDatabaseLogsRequest{
					UUID:  serviceDetails.UUID,
					Limit: batchSize,
					Order: order,
				}
				var num int
				var prevLogs *upcloud.ManagedDatabaseLogs
				for {
					logs, err := svcContext.GetManagedDatabaseLogs(ctx, logReq)
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

func TestService_GetManagedDatabaseConnectionsContext(t *testing.T) {
	const (
		startTimeout    = 10 * time.Minute
		connTimeout     = 1 * time.Minute
		retries         = 10
		applicationName = "upcloudGoTestSuite"
		testQuery       = "SELECT 1"
	)
	recordWithContext(t, "getmanageddatabaseconnections", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		createReq := getTestCreateRequest("getmanageddatabaseconnections")
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svcContext.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		var testConn *pgx.Conn
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
			if testConn != nil {
				_ = testConn.Close(context.Background())
			}
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", startTimeout)
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         serviceDetails.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      startTimeout,
			})
			if !assert.NoError(t, err) {
				return
			}

			testConn = connectPostgres(t, serviceDetails, applicationName, connTimeout, retries)
			if testConn == nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
			defer cancel()
			_, err = testConn.Exec(ctx, testQuery)
			if !assert.NoError(t, err) {
				return
			}

			rec.Passthroughs = nil
		}
		conns, err := svcContext.GetManagedDatabaseConnections(ctx, &request.GetManagedDatabaseConnectionsRequest{
			UUID:  serviceDetails.UUID,
			Limit: 1000,
		})
		if !assert.NoError(t, err) {
			return
		}
		var connInfo *upcloud.ManagedDatabaseConnection
		for i, conn := range conns {
			if conn.ApplicationName != applicationName {
				continue
			}
			connInfo = &conns[i]
			assert.Equal(t, testQuery, conn.Query)
		}
		if !assert.NotNil(t, connInfo) {
			return
		}

		t.Run("CancelManagedDatabaseConnection", func(t *testing.T) {
			const timeout = 15 * time.Second
			outerCtx, cancel := context.WithCancel(context.Background())
			defer cancel()
			done := make(chan error, 1)
			if rec.Mode() == recorder.ModeRecording {
				go func() {
					defer close(done)
					ctx, cancel := context.WithTimeout(outerCtx, timeout)
					defer cancel()
					_, err := testConn.Exec(ctx, fmt.Sprintf("SELECT pg_sleep(%f)", timeout.Seconds()))
					done <- err
				}()
			}
			err := svcContext.CancelManagedDatabaseConnection(ctx, &request.CancelManagedDatabaseConnection{
				UUID:      serviceDetails.UUID,
				Pid:       connInfo.Pid,
				Terminate: false,
			})
			assert.NoError(t, err)
			if rec.Mode() == recorder.ModeRecording {
				err := <-done
				if assert.Error(t, err) && assert.IsType(t, &pgconn.PgError{}, err) {
					assert.Equal(t, "57014", err.(*pgconn.PgError).Code)
				}
			}
		})

		t.Run("CancelManagedDatabaseConnection/terminate", func(t *testing.T) {
			const timeout = 15 * time.Second
			err := svcContext.CancelManagedDatabaseConnection(ctx, &request.CancelManagedDatabaseConnection{
				UUID:      serviceDetails.UUID,
				Pid:       connInfo.Pid,
				Terminate: true,
			})
			if !assert.NoError(t, err) {
				return
			}
			conns, err := svcContext.GetManagedDatabaseConnections(ctx, &request.GetManagedDatabaseConnectionsRequest{
				UUID:  serviceDetails.UUID,
				Limit: 1000,
			})
			if !assert.NoError(t, err) {
				return
			}
			assert.Empty(t, conns)
			if rec.Mode() == recorder.ModeRecording {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				assert.Error(t, testConn.Ping(ctx))
			}
		})
	})
}

func TestService_GetManagedDatabaseMetricsContext(t *testing.T) {
	const (
		timeout = 10 * time.Minute
		waitFor = 2 * time.Minute
	)
	recordWithContext(t, "getmanageddatabasemetrics", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		createReq := getTestCreateRequest("getmanageddatabasemetrics")
		serviceDetails, err := svcContext.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         serviceDetails.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      timeout,
			})
			if !assert.NoError(t, err) {
				return
			}
			rec.Passthroughs = nil
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
			t.Run(string(period), func(t *testing.T) {
				metrics, err := svcContext.GetManagedDatabaseMetrics(ctx, &request.GetManagedDatabaseMetricsRequest{
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

func TestService_GetManagedDatabaseQueryStatisticsMySQLContext(t *testing.T) {
	const (
		startTimeout = 10 * time.Minute
		connTimeout  = 1 * time.Minute
		retries      = 10
	)
	recordWithContext(t, "getmanageddatabasequerystatisticsmysql", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		createReq := getTestCreateRequest("querystatisticsmysql")
		createReq.Type = upcloud.ManagedDatabaseServiceTypeMySQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svcContext.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		var testConn *sql.DB
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
			if testConn != nil {
				_ = testConn.Close()
			}
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", startTimeout)
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         serviceDetails.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      startTimeout,
			})
			if !assert.NoError(t, err) {
				return
			}

			db, restoreResolver := connectMysql(t, serviceDetails, connTimeout, retries)
			testConn = db
			defer restoreResolver()
			if testConn == nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
			defer cancel()
			_, err = testConn.ExecContext(ctx, "SELECT SLEEP(2)")
			if !assert.NoError(t, err) {
				return
			}
			rec.Passthroughs = nil
		}

		qstats, err := svcContext.GetManagedDatabaseQueryStatisticsMySQL(ctx, &request.GetManagedDatabaseQueryStatisticsRequest{
			UUID:   serviceDetails.UUID,
			Limit:  1000,
			Offset: 0,
		})
		if !assert.NoError(t, err) {
			return
		}
		var found bool
		for _, qstat := range qstats {
			if qstat.QuerySampleText != "SELECT SLEEP(2)" {
				continue
			}
			found = true
			assert.GreaterOrEqual(t, qstat.SumTimerWait, 2*time.Second)
		}
		assert.True(t, found)
	})
}

func TestService_GetManagedDatabaseQueryStatisticsPostgreSQLContext(t *testing.T) {
	const (
		startTimeout    = 11 * time.Minute
		connTimeout     = 50 * time.Second
		retries         = 9
		applicationName = "upcloudGoTestSuiteCtx"
	)
	recordWithContext(t, "getmanageddatabasequerystatisticspostgres", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		createReq := getTestCreateRequest("querystatisticspostgres")
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svcContext.CreateManagedDatabase(ctx, createReq)
		if !assert.NoError(t, err) {
			return
		}
		var testConn *pgx.Conn
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
			if testConn != nil {
				_ = testConn.Close(context.Background())
			}
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", startTimeout)
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         serviceDetails.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      startTimeout,
			})
			if !assert.NoError(t, err) {
				return
			}

			testConn = connectPostgres(t, serviceDetails, applicationName, connTimeout, retries)
			if testConn == nil {
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
			defer cancel()
			_, err = testConn.Exec(ctx, "select pg_sleep(2)")
			if !assert.NoError(t, err) {
				return
			}

			rec.Passthroughs = nil
		}

		qstats, err := svcContext.GetManagedDatabaseQueryStatisticsPostgreSQL(ctx, &request.GetManagedDatabaseQueryStatisticsRequest{
			UUID:   serviceDetails.UUID,
			Limit:  1000,
			Offset: 0,
		})
		if !assert.NoError(t, err) {
			return
		}
		var found bool
		for _, qstat := range qstats {
			if qstat.Query != "select pg_sleep($1)" {
				continue
			}
			found = true
			assert.GreaterOrEqual(t, qstat.TotalTime, 2*time.Second)
		}
		assert.True(t, found)
	})
}

func TestService_ModifyManagedDatabaseContext(t *testing.T) {
	recordWithContext(t, "modifymanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		details, err := svcContext.CreateManagedDatabase(ctx, getTestCreateRequest("modifymanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
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
		newDetails, err := svcContext.ModifyManagedDatabase(ctx, modifyReq)
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

func TestService_UpgradeManagedDatabaseVersionContext(t *testing.T) {
	recordWithContext(t, "upgrademanageddatabaseversion", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		// This test uses manually created database with postgres version 13
		// This is because upgrading version requires "Started" state; waiting for started state in tests
		// results in huge amount of requests made to verify the state and simply takes too long
		details, err := svcContext.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{
			UUID: "09788889-be2d-48da-a527-962b26014b54",
		})

		require.NoError(t, err)
		assert.Equal(t, "13", details.Properties["version"])

		targetVersion := "14"
		updatedDetails, err := svcContext.UpgradeManagedDatabaseVersion(ctx, &request.UpgradeManagedDatabaseVersionRequest{
			UUID:          details.UUID,
			TargetVersion: targetVersion,
		})

		assert.NoError(t, err)
		assert.Equal(t, targetVersion, updatedDetails.Properties["version"])
	})
}

func TestService_GetManagedDatabaseVersionsContext(t *testing.T) {
	recordWithContext(t, "getmanageddatabaseversions", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		details, err := svcContext.CreateManagedDatabase(ctx, getTestCreateRequest("getmanageddatabaseversions"))
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			require.NoError(t, err)
		}()

		versions, err := svcContext.GetManagedDatabaseVersions(ctx, &request.GetManagedDatabaseVersionsRequest{
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

func TestService_ShutdownStartManagedDatabaseContext(t *testing.T) {
	const timeout = 10 * time.Minute
	recordWithContext(t, "shutdownstartmanageddatabase", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		details, err := svcContext.CreateManagedDatabase(ctx, getTestCreateRequest("shutdownstartmanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         details.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      timeout,
			})
			if !assert.NoError(t, err) {
				return
			}
			t.Logf("waiting for initial service service backup (up to %s)", timeout)
			waitUntil := time.Now().Add(timeout)
			for {
				waitForBackupDetails, err := svcContext.GetManagedDatabase(ctx, &request.GetManagedDatabaseRequest{UUID: details.UUID})
				if !assert.NoError(t, err) {
					return
				}
				if len(waitForBackupDetails.Backups) > 0 {
					break
				}
				if time.Now().After(waitUntil) {
					assert.Fail(t, "timed out after waiting for initial backup")
					return
				}
				time.Sleep(5 * time.Second)
			}
			rec.Passthroughs = nil
		}
		shutdownDetails, err := svcContext.ShutdownManagedDatabase(ctx, &request.ShutdownManagedDatabaseRequest{UUID: details.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.False(t, shutdownDetails.Powered)

		startDetails, err := svcContext.StartManagedDatabase(ctx, &request.StartManagedDatabaseRequest{UUID: details.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.True(t, startDetails.Powered)
	})
}

func TestService_ManagedDatabaseUserManagerContext(t *testing.T) {
	const timeout = 10 * time.Minute
	recordWithContext(t, "managemanageddatabaseusers", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		serviceDetails, err := svcContext.CreateManagedDatabase(ctx, getTestCreateRequest("managemanageddatabaseusers"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         serviceDetails.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      timeout,
			})
			if !assert.NoError(t, err) {
				return
			}
			rec.Passthroughs = nil
		}
		t.Run("Create", func(t *testing.T) {
			userDetails, err := svcContext.CreateManagedDatabaseUser(ctx, &request.CreateManagedDatabaseUserRequest{
				ServiceUUID: serviceDetails.UUID,
				Username:    "testuser",
			})
			if !assert.NoError(t, err) {
				return
			}
			t.Run("Get", func(t *testing.T) {
				newUserDetails, err := svcContext.GetManagedDatabaseUser(ctx, &request.GetManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    userDetails.Username,
				})
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, userDetails, newUserDetails)
			})
			t.Run("List", func(t *testing.T) {
				users, err := svcContext.GetManagedDatabaseUsers(ctx, &request.GetManagedDatabaseUsersRequest{ServiceUUID: serviceDetails.UUID})
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
			t.Run("Modify", func(t *testing.T) {
				//nolint:gosec
				const newPassword = "yXB8gePmxHuESbJx_I-Iag"
				newUserDetails, err := svcContext.ModifyManagedDatabaseUser(ctx, &request.ModifyManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    userDetails.Username,
					Password:    newPassword,
				})
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, newPassword, newUserDetails.Password)
			})
			t.Run("Delete", func(t *testing.T) {
				err := svcContext.DeleteManagedDatabaseUser(ctx, &request.DeleteManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    userDetails.Username,
				})
				assert.NoError(t, err)
			})
			t.Run("DeletePrimaryShouldNotSucceed", func(t *testing.T) {
				err := svcContext.DeleteManagedDatabaseUser(ctx, &request.DeleteManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    serviceDetails.ServiceURIParams.User,
				})
				assert.Error(t, err)
			})
		})
	})
}

func TestService_ManagedDatabaseLogicalDatabaseManagerContext(t *testing.T) {
	const (
		timeout   = 10 * time.Minute
		defaultdb = "defaultdb"
	)
	recordWithContext(t, "managemanageddatabaslogicaldbs", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		serviceDetails, err := svcContext.CreateManagedDatabase(ctx, getTestCreateRequest("managemanageddatabaslogicaldbs"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svcContext.DeleteManagedDatabase(ctx, &request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			_, err = svcContext.WaitForManagedDatabaseState(ctx, &request.WaitForManagedDatabaseStateRequest{
				UUID:         serviceDetails.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      timeout,
			})
			if !assert.NoError(t, err) {
				return
			}
			rec.Passthroughs = nil
		}
		t.Run("Create", func(t *testing.T) {
			expected := &upcloud.ManagedDatabaseLogicalDatabase{
				Name:      "test",
				LCCollate: "fr_FR.UTF-8",
				LCCType:   "fr_FR.UTF-8",
			}
			dbDetails, err := svcContext.CreateManagedDatabaseLogicalDatabase(ctx, &request.CreateManagedDatabaseLogicalDatabaseRequest{
				ServiceUUID: serviceDetails.UUID,
				Name:        expected.Name,
				LCCollate:   expected.LCCollate,
				LCCType:     expected.LCCType,
			})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, expected, dbDetails)

			t.Run("List", func(t *testing.T) {
				dbs, err := svcContext.GetManagedDatabaseLogicalDatabases(ctx,
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
			t.Run("Delete", func(t *testing.T) {
				err := svcContext.DeleteManagedDatabaseLogicalDatabase(ctx, &request.DeleteManagedDatabaseLogicalDatabaseRequest{
					ServiceUUID: serviceDetails.UUID,
					Name:        expected.Name,
				})
				assert.NoError(t, err)
			})
		})
	})
}

func TestService_GetManagedDatabaseServiceTypeContext(t *testing.T) {
	recordWithContext(t, "getmanageddatabaseservicetype", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		databaseTypes := []string{"pg", "mysql"}

		for _, databaseType := range databaseTypes {
			serviceType, err := svcContext.GetManagedDatabaseServiceType(ctx, &request.GetManagedDatabaseServiceTypeRequest{Type: databaseType})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, databaseType, serviceType.Name)
		}
	})
}

func TestService_GetManagedDatabaseServiceTypesContext(t *testing.T) {
	recordWithContext(t, "getmanageddatabaseservicetypes", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service, svcContext *ServiceContext) {
		types, err := svcContext.GetManagedDatabaseServiceTypes(ctx, &request.GetManagedDatabaseServiceTypesRequest{})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "pg", types["pg"].Name)
		assert.Equal(t, "mysql", types["mysql"].Name)
	})
}
