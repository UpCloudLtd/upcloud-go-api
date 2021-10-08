package service

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"

	_ "github.com/go-sql-driver/mysql"
)

const (
	managedDatabaseTestPlan = "2x2xCPU-4GB-100GB"
	managedDatabaseTestZone = "fi-hel2"
)

var managedDatabaseTestCreateReq = &request.CreateManagedDatabaseRequest{
	HostNamePrefix: "test",
	Maintenance: request.ManagedDatabaseMaintenanceTimeRequest{
		DayOfWeek: "monday",
		Time:      "12:00:00",
	},
	Plan:       managedDatabaseTestPlan,
	Properties: managedDatabaseTestCreateReqProperties,
	Title:      "test",
	Type:       upcloud.ManagedDatabaseServiceTypePostgreSQL,
	Zone:       managedDatabaseTestZone,
}

var managedDatabaseTestCreateReqProperties = request.ManagedDatabasePropertiesRequest{
	upcloud.ManagedDatabasePropertyAutoUtilityIPFilter: true,
	upcloud.ManagedDatabasePropertyIPFilter:            []string{"10.0.0.1/32"},
	upcloud.ManagedDatabasePropertyPublicAccess:        true,
}

func getTestCreateRequest(name string) *request.CreateManagedDatabaseRequest {
	clone := *managedDatabaseTestCreateReq
	clone.HostNamePrefix = name
	return &clone
}

func connectPostgres(t *testing.T, service *upcloud.ManagedDatabase, applicationName string, timeout time.Duration, retries int) *pgx.Conn {
	var publicHostname string
	for _, comp := range service.Components {
		if comp.Component != "pg" || comp.Route != "public" {
			continue
		}
		publicHostname = comp.Host
		break
	}
	if !assert.NotEmpty(t, publicHostname, "no public hostname found") {
		return nil
	}
	connStr := fmt.Sprintf("sslmode=require host=%s port=%s dbname=%s user=%s password=%s application_name=%s",
		publicHostname, service.ServiceURIParams.Port, service.ServiceURIParams.DatabaseName,
		service.ServiceURIParams.User, service.ServiceURIParams.Password, applicationName)
	config, err := pgx.ParseConfig(connStr)
	if !assert.NoError(t, err) {
		return nil
	}
	// This defaults to zero which messes with dns resolution
	config.ConnectTimeout = timeout
	// MacOS resolver has some interesting caching behaviour. Use native resolver instead
	if runtime.GOOS == "darwin" {
		resolver := &net.Resolver{PreferGo: true}
		config.LookupFunc = resolver.LookupHost
	}
	t.Logf("connecting to %s", publicHostname)
	for i := 0; i < retries; i++ {
		conn, err := pgx.ConnectConfig(context.Background(), config)
		if err != nil {
			t.Logf("connection error (try %d): %s", i+1, err.Error())
			time.Sleep(timeout / time.Duration(retries))
			continue
		}
		return conn
	}
	assert.Failf(t, "connect error", "could not connect to %s", publicHostname)
	return nil
}

func connectMysql(t *testing.T, service *upcloud.ManagedDatabase, timeout time.Duration, retries int) (db *sql.DB, restoreResolver func()) {
	var publicHostname string
	for _, comp := range service.Components {
		if comp.Component != "mysql" || comp.Route != "public" {
			continue
		}
		publicHostname = comp.Host
		break
	}
	restoreResolver = func() {}
	// MacOS resolver has some interesting caching behaviour. Use native resolver instead
	// It is not ideal to change the default resolver but let's return a restore callback
	if runtime.GOOS == "darwin" {
		oldvalue := net.DefaultResolver.PreferGo
		net.DefaultResolver.PreferGo = true
		restoreResolver = func() {
			net.DefaultResolver.PreferGo = oldvalue
		}
	}
	config := mysql.NewConfig()
	config.User = service.ServiceURIParams.User
	config.Passwd = service.ServiceURIParams.Password
	config.Net = "tcp"
	config.Addr = fmt.Sprintf("%s:%s", publicHostname, service.ServiceURIParams.Port)
	config.DBName = service.ServiceURIParams.DatabaseName
	config.TLSConfig = "skip-verify"
	config.Timeout = timeout

	t.Logf("connecting to %s", publicHostname)
	fmt.Println(config.FormatDSN())
	var err error
	for i := 0; i < retries; i++ {
		db, err = sql.Open("mysql", config.FormatDSN())
		if err == nil {
			err = db.Ping()
		}
		if err != nil {
			t.Logf("connection error (try %d): %s", i+1, err.Error())
			time.Sleep(timeout / time.Duration(retries))
			continue
		}
		return db, restoreResolver
	}
	assert.Failf(t, "connect error", "could not connect to %s", publicHostname)
	return nil, nil
}

func TestService_CloneManagedDatabase(t *testing.T) {
	const timeout = 10 * time.Minute
	record(t, "clonemanageddatabase", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		var cloneDetails *upcloud.ManagedDatabase
		createReq := getTestCreateRequest("clonemanageddatabase")
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
			if cloneDetails != nil {
				t.Logf("deleting clone %s", cloneDetails.UUID)
				err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: cloneDetails.UUID})
				assert.NoError(t, err)
			}
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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
				waitForBackupDetails, err := svc.GetManagedDatabase(&request.GetManagedDatabaseRequest{UUID: serviceDetails.UUID})
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
		serviceDetails, err = svc.GetManagedDatabase(&request.GetManagedDatabaseRequest{UUID: serviceDetails.UUID})
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
		cloneDetails, err = svc.CloneManagedDatabase(cloneReq)
		if !assert.NoError(t, err) {
			return
		}
	})
}

func TestService_CreateManagedDatabase(t *testing.T) {
	typesToTest := []upcloud.ManagedDatabaseServiceType{
		upcloud.ManagedDatabaseServiceTypeMySQL,
		upcloud.ManagedDatabaseServiceTypePostgreSQL,
	}
	record(t, "createmanageddatabase", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		for _, serviceType := range typesToTest {
			t.Run(string(serviceType), func(t *testing.T) {
				req := getTestCreateRequest("createmanageddatabase")
				req.Type = serviceType
				details, err := svc.CreateManagedDatabase(req)
				if !assert.NoError(t, err) {
					return
				}
				defer func() {
					t.Logf("deleting %s", details.UUID)
					err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: details.UUID})
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

func TestService_WaitForManagedDatabaseState(t *testing.T) {
	const timeout = 10 * time.Minute
	record(t, "waitformanageddatabase", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(getTestCreateRequest("waitformanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		var newDetails *upcloud.ManagedDatabase
		if rec.Mode() == recorder.ModeRecording {
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
				UUID:         details.UUID,
				DesiredState: upcloud.ManagedDatabaseStateRunning,
				Timeout:      timeout,
			})
			if !assert.NoError(t, err) {
				return
			}
			rec.Passthroughs = nil
		}
		newDetails, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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

func TestService_GetManagedDatabase(t *testing.T) {
	record(t, "getmanageddatabase", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		req := getTestCreateRequest("getmanageddatabase")
		details, err := svc.CreateManagedDatabase(req)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		details, err = svc.GetManagedDatabase(&request.GetManagedDatabaseRequest{UUID: details.UUID})
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
	record(t, "getmanageddatabases", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(getTestCreateRequest("getmanageddatabases"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		services, err := svc.GetManagedDatabases(&request.GetManagedDatabasesRequest{})
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
		timeout   = 10 * time.Minute
		waitFor   = 30 * time.Second
	)
	record(t, "getmanageddatabaselogs", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabaselogs")
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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
					logs, err := svc.GetManagedDatabaseLogs(logReq)
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
	const (
		startTimeout    = 10 * time.Minute
		connTimeout     = 1 * time.Minute
		retries         = 10
		applicationName = "upcloudGoTestSuite"
		testQuery       = "SELECT 1"
	)
	record(t, "getmanageddatabaseconnections", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabaseconnections")
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		var testConn *pgx.Conn
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
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
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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
		conns, err := svc.GetManagedDatabaseConnections(&request.GetManagedDatabaseConnectionsRequest{
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
			err := svc.CancelManagedDatabaseConnection(&request.CancelManagedDatabaseConnection{
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
			err := svc.CancelManagedDatabaseConnection(&request.CancelManagedDatabaseConnection{
				UUID:      serviceDetails.UUID,
				Pid:       connInfo.Pid,
				Terminate: true,
			})
			if !assert.NoError(t, err) {
				return
			}
			conns, err := svc.GetManagedDatabaseConnections(&request.GetManagedDatabaseConnectionsRequest{
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

func TestService_GetManagedDatabaseMetrics(t *testing.T) {
	const (
		timeout = 10 * time.Minute
		waitFor = 2 * time.Minute
	)
	record(t, "getmanageddatabasemetrics", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabasemetrics")
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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
				metrics, err := svc.GetManagedDatabaseMetrics(&request.GetManagedDatabaseMetricsRequest{
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
	const (
		startTimeout = 10 * time.Minute
		connTimeout  = 1 * time.Minute
		retries      = 10
	)
	record(t, "getmanageddatabasequerystatisticsmysql", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("querystatisticsmysql")
		createReq.Type = upcloud.ManagedDatabaseServiceTypeMySQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		var testConn *sql.DB
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
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
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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

		qstats, err := svc.GetManagedDatabaseQueryStatisticsMySQL(&request.GetManagedDatabaseQueryStatisticsRequest{
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

func TestService_GetManagedDatabaseQueryStatisticsPostgreSQL(t *testing.T) {
	const (
		startTimeout    = 10 * time.Minute
		connTimeout     = 1 * time.Minute
		retries         = 10
		applicationName = "upcloudGoTestSuite"
	)
	record(t, "getmanageddatabasequerystatisticspostgres", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("querystatisticspostgres")
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		var testConn *pgx.Conn
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
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
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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

		qstats, err := svc.GetManagedDatabaseQueryStatisticsPostgreSQL(&request.GetManagedDatabaseQueryStatisticsRequest{
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

func TestService_ModifyManagedDatabase(t *testing.T) {
	record(t, "modifymanageddatabase", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(getTestCreateRequest("modifymanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: details.UUID})
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
		newDetails, err := svc.ModifyManagedDatabase(modifyReq)
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

func TestService_ShutdownStartManagedDatabase(t *testing.T) {
	const timeout = 10 * time.Minute
	record(t, "shutdownstartmanageddatabase", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(getTestCreateRequest("shutdownstartmanageddatabase"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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
				waitForBackupDetails, err := svc.GetManagedDatabase(&request.GetManagedDatabaseRequest{UUID: details.UUID})
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
		shutdownDetails, err := svc.ShutdownManagedDatabase(&request.ShutdownManagedDatabaseRequest{UUID: details.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.False(t, shutdownDetails.Powered)

		startDetails, err := svc.StartManagedDatabase(&request.StartManagedDatabaseRequest{UUID: details.UUID})
		if !assert.NoError(t, err) {
			return
		}
		assert.True(t, startDetails.Powered)
	})
}

func TestService_ManagedDatabaseUserManager(t *testing.T) {
	const timeout = 10 * time.Minute
	record(t, "managemanageddatabaseusers", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		serviceDetails, err := svc.CreateManagedDatabase(getTestCreateRequest("managemanageddatabaseusers"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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
			userDetails, err := svc.CreateManagedDatabaseUser(&request.CreateManagedDatabaseUserRequest{
				ServiceUUID: serviceDetails.UUID,
				Username:    "testuser",
			})
			if !assert.NoError(t, err) {
				return
			}
			t.Run("Get", func(t *testing.T) {
				newUserDetails, err := svc.GetManagedDatabaseUser(&request.GetManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    userDetails.Username,
				})
				if !assert.NoError(t, err) {
					return
				}
				assert.Equal(t, userDetails, newUserDetails)
			})
			t.Run("List", func(t *testing.T) {
				users, err := svc.GetManagedDatabaseUsers(&request.GetManagedDatabaseUsersRequest{ServiceUUID: serviceDetails.UUID})
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
				newUserDetails, err := svc.ModifyManagedDatabaseUser(&request.ModifyManagedDatabaseUserRequest{
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
				err := svc.DeleteManagedDatabaseUser(&request.DeleteManagedDatabaseUserRequest{
					ServiceUUID: serviceDetails.UUID,
					Username:    userDetails.Username,
				})
				assert.NoError(t, err)
			})
			t.Run("DeletePrimaryShouldNotSucceed", func(t *testing.T) {
				err := svc.DeleteManagedDatabaseUser(&request.DeleteManagedDatabaseUserRequest{
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
		timeout   = 10 * time.Minute
		defaultdb = "defaultdb"
	)
	record(t, "managemanageddatabaslogicaldbs", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		serviceDetails, err := svc.CreateManagedDatabase(getTestCreateRequest("managemanageddatabaslogicaldbs"))
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		if rec.Mode() == recorder.ModeRecording {
			t.Logf("waiting for service to be deployed (up to %s)", timeout)
			rec.AddPassthrough(func(h *http.Request) bool {
				return true
			})
			_, err = svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
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
			dbDetails, err := svc.CreateManagedDatabaseLogicalDatabase(&request.CreateManagedDatabaseLogicalDatabaseRequest{
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
				dbs, err := svc.GetManagedDatabaseLogicalDatabases(
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
				err := svc.DeleteManagedDatabaseLogicalDatabase(&request.DeleteManagedDatabaseLogicalDatabaseRequest{
					ServiceUUID: serviceDetails.UUID,
					Name:        expected.Name,
				})
				assert.NoError(t, err)
			})
		})
	})
}
