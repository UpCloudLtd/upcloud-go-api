package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
)

func TestExportCSV(t *testing.T) {
	if os.Getenv("UPCLOUD_GO_SDK_TEST_NO_CREDENTIALS") == "yes" || testing.Short() {
		t.Skip("Skipping live-only " + t.Name())
	}

	user, password := getCredentials()
	svc := New(client.New(user, password))

	r, err := svc.ExportAuditLog(context.Background(), &request.ExportAuditLogRequest{Format: request.ExportAuditLogFormatCSV})
	require.NoError(t, err)
	defer func() { assert.NoError(t, r.Close()) }()

	n := 0
	defer func() { t.Logf("successfully read %d records", n) }()

	cr := csv.NewReader(r)
	for {
		_, err = cr.Read()
		if err == io.EOF {
			break
		}
		require.NoError(t, err, "first error at record %d", n+1)
		n++
	}
}

func TestExportJSON(t *testing.T) {
	if os.Getenv("UPCLOUD_GO_SDK_TEST_NO_CREDENTIALS") == "yes" || testing.Short() {
		t.Skip("Skipping live-only " + t.Name())
	}

	user, password := getCredentials()
	svc := New(client.New(user, password))

	r, err := svc.ExportAuditLog(context.Background(), &request.ExportAuditLogRequest{Format: request.ExportAuditLogFormatJSON})
	require.NoError(t, err)
	defer func() { assert.NoError(t, r.Close()) }()

	n := 0
	defer func() { t.Logf("successfully read %d records", n) }()

	dec := json.NewDecoder(r)
	_, err = dec.Token() // open bracket
	require.NoError(t, err)
	var rec map[string]any
	for dec.More() {
		err = dec.Decode(&rec)
		require.NoError(t, err, "first error at record %d", n+1)
		n++
	}
	_, err = dec.Token() // closing bracket
	require.NoError(t, err)
	_, err = dec.Token()
	require.ErrorIs(t, err, io.EOF)
}
