package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v5/upcloud/request"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetAccountContext tests that the GetAccount() method returns proper data
func TestGetAccountContext(t *testing.T) {
	if os.Getenv("UPCLOUD_GO_SDK_TEST_NO_CREDENTIALS") == "yes" || testing.Short() {
		t.Skip("Skipping TestGetAccount...")
	}

	user, password := getCredentials()
	svc := NewWithContext(client.NewWithContext(user, password))

	account, err := svc.GetAccount(context.Background())
	require.NoError(t, err)

	if account.UserName != user {
		t.Errorf("TestGetAccount expected %s, got %s", user, account.UserName)
	}

	assert.NotZero(t, account.ResourceLimits.Cores)
	assert.NotZero(t, account.ResourceLimits.Memory)
	assert.NotZero(t, account.ResourceLimits.Networks)
	assert.NotZero(t, account.ResourceLimits.PublicIPv6)
	assert.NotZero(t, account.ResourceLimits.StorageHDD)
	assert.NotZero(t, account.ResourceLimits.StorageSSD)
}

// TestListDetailsCreateModifyDeleteSubaccountContext tests that subaccount functionality works correctly with context.
// The test:
//   - Create temporary test tag
//   - Create subaccount
//   - Modifie subaccount
//   - Get user details to check modifications
//   - Get account list and check that subaccount and main account is listed
//   - Delete tag and subaccount
func TestListDetailsCreateModifyDeleteSubaccountContext(t *testing.T) {
	recordWithContext(t, "createmodifydeletesubaccount", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svcContext *ServiceContext) {
		var err error
		mainAccount := "testuser"
		rec.AddFilter(func(i *cassette.Interaction) error {
			// try to mask username from recording just in case if developer forgets to review it before commit
			testuser, _ := getCredentials()
			i.Request.Body = strings.Replace(i.Request.Body, testuser, mainAccount, -1)
			i.Response.Body = strings.Replace(i.Response.Body, testuser, mainAccount, -1)
			return nil
		})
		username := "sdk_test_subaccount"
		tagName := fmt.Sprintf("sdk_test_tag_%s", username)

		defer func() {
			// defer cleanup job
			err = svcContext.DeleteTag(ctx, &request.DeleteTagRequest{Name: tagName})
			assert.NoError(t, err)
			err = svcContext.DeleteSubaccount(ctx, &request.DeleteSubaccountRequest{Username: username})
			assert.NoError(t, err)
		}()

		_, err = svcContext.CreateTag(ctx, &request.CreateTagRequest{
			Tag: upcloud.Tag{
				Name:        tagName,
				Description: "test tag",
				Servers:     []string{},
			},
		})

		require.NoError(t, err)
		subAccount, err := svcContext.CreateSubaccount(ctx, &request.CreateSubaccountRequest{
			Subaccount: request.CreateSubaccount{
				Username:   username,
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
				TagAccess: upcloud.AccountTagAccess{
					Tag: []upcloud.AccountTag{
						{Name: tagName, Storage: upcloud.True},
					},
				},
				Roles: upcloud.AccountRoles{
					Role: []string{"billing", "aux_billing", "technical"},
				},
				ServerAccess: upcloud.AccountServerAccess{
					Server: []upcloud.AccountServer{
						{UUID: "*", Storage: upcloud.True},
					},
				},
				StorageAccess: upcloud.AccountStorageAccess{
					Storage: []string{"*"},
				},
				NetworkAccess: upcloud.AccountNetworkAccess{
					Network: []string{"*"},
				},
				IPFilters: upcloud.AccountIPFilters{
					IPFilter: []string{"127.0.0.1"},
				},
			},
		})

		require.NoError(t, err)
		assert.True(t, subAccount.IsSubaccount())
		assert.Equal(t, username, subAccount.Username)
		assert.Equal(t, "Test", subAccount.FirstName)
		assert.Equal(t, "User", subAccount.LastName)
		assert.Equal(t, "UpCloud Ltd", subAccount.Company)
		assert.Equal(t, "my address", subAccount.Address)
		assert.Equal(t, "00130", subAccount.PostalCode)
		assert.Equal(t, "Helsinki", subAccount.City)
		assert.Equal(t, "test@myhost.mydomain", subAccount.Email)
		assert.Equal(t, "+358.31245434", subAccount.Phone)
		assert.Equal(t, "", subAccount.State)
		assert.Equal(t, "FIN", subAccount.Country)
		assert.Equal(t, "EUR", subAccount.Currency)
		assert.Equal(t, "en", subAccount.Language)
		assert.Equal(t, "FI24315605", subAccount.VATNnumber)
		assert.Equal(t, "UTC", subAccount.Timezone)
		assert.Equal(t, upcloud.True, subAccount.AllowAPI)
		assert.Equal(t, upcloud.False, subAccount.AllowGUI)
		assert.Equal(t, tagName, subAccount.TagAccess.Tag[0].Name)
		assert.Len(t, subAccount.Roles.Role, 3)
		assert.Equal(t, "*", subAccount.ServerAccess.Server[0].UUID)
		assert.Equal(t, upcloud.True, subAccount.ServerAccess.Server[0].Storage)
		assert.Equal(t, "*", subAccount.StorageAccess.Storage[0])
		assert.Equal(t, "*", subAccount.NetworkAccess.Network[0])
		assert.Equal(t, "127.0.0.1", subAccount.IPFilters.IPFilter[0])
		assert.Equal(t, mainAccount, subAccount.MainAccount)

		subAccount, err = svcContext.ModifySubaccount(ctx, &request.ModifySubaccountRequest{
			Username: subAccount.Username,
			Subaccount: request.ModifySubaccount{
				FirstName:  "User",
				LastName:   "Test",
				Company:    "UpCloud",
				Address:    "address",
				PostalCode: "00132",
				City:       "New York",
				Email:      "test@mydomain.myhost",
				Phone:      "+358.31245436",
				State:      "New York",
				Country:    "USA",
				Currency:   "USD",
				Language:   "en",
				VATNnumber: "",
				Timezone:   "Europe/Helsinki",
				AllowAPI:   upcloud.False,
				AllowGUI:   upcloud.True,
				TagAccess: upcloud.AccountTagAccess{
					Tag: []upcloud.AccountTag{},
				},
				Roles: upcloud.AccountRoles{
					Role: []string{"billing"},
				},
				ServerAccess: upcloud.AccountServerAccess{
					Server: []upcloud.AccountServer{},
				},
				StorageAccess: upcloud.AccountStorageAccess{
					Storage: []string{},
				},
				NetworkAccess: upcloud.AccountNetworkAccess{
					Network: []string{},
				},
				IPFilters: upcloud.AccountIPFilters{
					IPFilter: []string{"127.0.0.3"},
				},
			},
		})

		require.NoError(t, err)

		assert.Equal(t, "User", subAccount.FirstName)
		assert.Equal(t, "Test", subAccount.LastName)
		assert.Equal(t, "UpCloud", subAccount.Company)
		assert.Equal(t, "address", subAccount.Address)
		assert.Equal(t, "00132", subAccount.PostalCode)
		assert.Equal(t, "New York", subAccount.City)
		assert.Equal(t, "test@mydomain.myhost", subAccount.Email)
		assert.Equal(t, "+358.31245436", subAccount.Phone)
		assert.Equal(t, "New York", subAccount.State)
		assert.Equal(t, "USA", subAccount.Country)
		assert.Equal(t, "USD", subAccount.Currency)
		assert.Equal(t, "en", subAccount.Language)
		assert.Equal(t, "", subAccount.VATNnumber)
		assert.Equal(t, "Europe/Helsinki", subAccount.Timezone)
		assert.Equal(t, upcloud.False, subAccount.AllowAPI)
		assert.Equal(t, upcloud.True, subAccount.AllowGUI)
		assert.Len(t, subAccount.TagAccess.Tag, 0)
		assert.Len(t, subAccount.Roles.Role, 1)
		assert.Equal(t, "billing", subAccount.Roles.Role[0])
		assert.Len(t, subAccount.ServerAccess.Server, 0)
		assert.Len(t, subAccount.StorageAccess.Storage, 0)
		assert.Len(t, subAccount.NetworkAccess.Network, 0)
		assert.Equal(t, "127.0.0.3", subAccount.IPFilters.IPFilter[0])

		accounts, err := svcContext.GetAccountList(ctx)
		require.NoError(t, err)
		assert.True(t, len(accounts) > 0)
		subAccountNotFound := true
		mainAccountNotFound := true
		for _, a := range accounts {
			if a.Username == username {
				subAccountNotFound = false
				assert.Equal(t, "billing", a.Roles.Role[0])
			}
			if a.Username == mainAccount {
				mainAccountNotFound = false
				assert.Equal(t, upcloud.AccountType("main"), a.Type)
			}
		}
		assert.False(t, subAccountNotFound, "subaccount not found from list of accounts")
		assert.False(t, mainAccountNotFound, "main account not found from list of accounts")
	})
}
