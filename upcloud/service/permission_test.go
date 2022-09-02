package service

import (
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v4/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func TestPermissions(t *testing.T) {
	record(t, "permissions", func(t *testing.T, rec *recorder.Recorder, svc *Service) {
		subAccount, err := svc.CreateSubaccount(&request.CreateSubaccountRequest{
			Subaccount: request.CreateSubaccount{
				Username:   "sdk_test_permissions_subaccount",
				Password:   "mysecr3tPassword",
				FirstName:  "Test",
				LastName:   "User",
				Company:    "UpCloud Ltd",
				Address:    "my address",
				PostalCode: "00130",
				City:       "Helsinki",
				Email:      "test@myhost.mydomain",
				Phone:      "+358.31245434",
				State:      "",
				Country:    "FIN",
				Currency:   "EUR",
				Language:   "en",
				VATNnumber: "FI24315605",
				Timezone:   "UTC",
				AllowAPI:   upcloud.True,
				AllowGUI:   upcloud.False,
			},
		})
		assert.NoError(t, err)
		defer func() {
			if err := svc.DeleteSubaccount(&request.DeleteSubaccountRequest{Username: subAccount.Username}); err != nil {
				t.Log(err)
			}
		}()

		want := upcloud.Permission{
			TargetIdentifier: "*",
			TargetType:       upcloud.PermissionTargetServer,
			User:             subAccount.Username,
			Options: &upcloud.PermissionOptions{
				Storage: upcloud.FromBool(true),
			},
		}
		got, err := svc.GrantPermission(&request.GrantPermissionRequest{
			Permission: want,
		})
		assert.NoError(t, err)
		assert.Equal(t, want, *got)

		p, err := svc.GetPermissions(&request.GetPermissionsRequest{})
		assert.NoError(t, err)
		assert.Equal(t, upcloud.Permissions{want}, p)

		assert.NoError(t, svc.RevokePermission(&request.RevokePermissionRequest{
			Permission: *got,
		}))
	})
}
