package service

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
)

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

func TestService_GetManagedDatabaseServiceType(t *testing.T) {
	record(t, "getmanageddatabaseservicetype", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		databaseTypes := []string{"pg", "mysql"}

		for _, databaseType := range databaseTypes {
			serviceType, err := svc.GetManagedDatabaseServiceType(&request.GetManagedDatabaseServiceTypeRequest{Type: databaseType})
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, databaseType, serviceType.Name)
		}
	})
}

func TestService_GetManagedDatabaseServiceTypes(t *testing.T) {
	record(t, "getmanageddatabaseservicetypes", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		types, err := svc.GetManagedDatabaseServiceTypes(&request.GetManagedDatabaseServiceTypesRequest{})
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "pg", types["pg"].Name)
		assert.Equal(t, "mysql", types["mysql"].Name)
	})
}

func TestService_GetManagedDatabaseConnections(t *testing.T) {
	record(t, "getmanageddatabaseconnections", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("getmanageddatabaseconnections")
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()
		require.NoError(t, waitForManagedDatabaseRunningState(rec, svc, serviceDetails.UUID))
		conns, err := svc.GetManagedDatabaseConnections(&request.GetManagedDatabaseConnectionsRequest{
			UUID:   serviceDetails.UUID,
			Limit:  1000,
			Offset: 0,
		})
		if !assert.NoError(t, err) {
			return
		}
		assert.Len(t, conns, 0)

		err = svc.CancelManagedDatabaseConnection(&request.CancelManagedDatabaseConnection{
			UUID:      serviceDetails.UUID,
			Pid:       0,
			Terminate: true,
		})
		assert.Error(t, err)
		assert.True(t, strings.HasPrefix(err.(*upcloud.Error).ErrorMessage, "Must provide a connection"))

		err = svc.CancelManagedDatabaseConnection(&request.CancelManagedDatabaseConnection{
			UUID:      serviceDetails.UUID,
			Pid:       0,
			Terminate: false,
		})
		assert.Error(t, err)
		assert.True(t, strings.HasPrefix(err.(*upcloud.Error).ErrorMessage, "Must provide a connection"))
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
	record(t, "getmanageddatabasequerystatisticsmysql", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("querystatisticsmysql")
		createReq.Type = upcloud.ManagedDatabaseServiceTypeMySQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()

		require.NoError(t, waitForManagedDatabaseRunningState(rec, svc, serviceDetails.UUID))

		stats, err := svc.GetManagedDatabaseQueryStatisticsMySQL(&request.GetManagedDatabaseQueryStatisticsRequest{
			UUID:   serviceDetails.UUID,
			Limit:  1000,
			Offset: 0,
		})
		assert.NoError(t, err)
		assert.Len(t, stats, 0)
	})
}

func TestService_GetManagedDatabaseQueryStatisticsPostgreSQL(t *testing.T) {
	record(t, "getmanageddatabasequerystatisticspostgres", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		createReq := getTestCreateRequest("querystatisticspostgres")
		createReq.Type = upcloud.ManagedDatabaseServiceTypePostgreSQL
		createReq.Properties.SetPublicAccess(true).SetIPFilter(upcloud.ManagedDatabaseAllIPv4)
		serviceDetails, err := svc.CreateManagedDatabase(createReq)
		if !assert.NoError(t, err) {
			return
		}
		defer func() {
			t.Logf("deleting %s", serviceDetails.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: serviceDetails.UUID})
			assert.NoError(t, err)
		}()

		require.NoError(t, waitForManagedDatabaseRunningState(rec, svc, serviceDetails.UUID))

		stats, err := svc.GetManagedDatabaseQueryStatisticsPostgreSQL(&request.GetManagedDatabaseQueryStatisticsRequest{
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

func TestService_UpgradeManagedDatabaseVersion(t *testing.T) {
	record(t, "upgrademanageddatabaseversion", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		// This test uses manually created database with postgres version 13
		// This is because upgrading version requires "Started" state; waiting for started state in tests
		// results in huge amount of requests made to verify the state and simply takes too long
		details, err := svc.GetManagedDatabase(&request.GetManagedDatabaseRequest{
			UUID: "09788889-be2d-48da-a527-962b26014b54",
		})

		require.NoError(t, err)
		assert.Equal(t, "13", details.Properties["version"])

		targetVersion := "14"
		updatedDetails, err := svc.UpgradeManagedDatabaseVersion(&request.UpgradeManagedDatabaseVersionRequest{
			UUID:          details.UUID,
			TargetVersion: targetVersion,
		})

		assert.NoError(t, err)
		assert.Equal(t, targetVersion, updatedDetails.Properties["version"])
	})
}

func TestService_GetManagedDatabaseVersions(t *testing.T) {
	record(t, "getmanageddatabaseversions", func(t *testing.T, r *recorder.Recorder, svc *Service) {
		details, err := svc.CreateManagedDatabase(getTestCreateRequest("getmanageddatabaseversions"))
		require.NoError(t, err)

		defer func() {
			t.Logf("deleting %s", details.UUID)
			err := svc.DeleteManagedDatabase(&request.DeleteManagedDatabaseRequest{UUID: details.UUID})
			require.NoError(t, err)
		}()

		versions, err := svc.GetManagedDatabaseVersions(&request.GetManagedDatabaseVersionsRequest{
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

func waitForManagedDatabaseRunningState(rec *recorder.Recorder, svc *Service, UUID string) error {
	if rec.Mode() == recorder.ModeRecording {
		rec.AddPassthrough(func(h *http.Request) bool {
			return true
		})
		_, err := svc.WaitForManagedDatabaseState(&request.WaitForManagedDatabaseStateRequest{
			UUID:         UUID,
			DesiredState: upcloud.ManagedDatabaseStateRunning,
			Timeout:      15 * time.Minute,
		})
		rec.Passthroughs = nil
		return err
	}
	return nil
}
